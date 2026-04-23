package service

import (
	"auth/internal/database"
	"auth/internal/models"
	"auth/jwt_utils"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/mail"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type Auth struct {
	Cache   database.RefreshCache
	Storage database.UserStorage
	Ctx     context.Context
}

func (auth *Auth) Register(w http.ResponseWriter, r *http.Request) {
	const op = "auth.Register"
	logger, _ := auth.Ctx.Value(("logger")).(*zap.Logger)

	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		logger.Error("bad request", zap.String("path", op), zap.Error(err))
		return
	}
	var user models.User
	json.Unmarshal(body, &user)

	_, err = mail.ParseAddress(user.Email)
	if err != nil {
		logger.Error("wrong email", zap.String("email", user.Email), zap.String("path", op), zap.Error(err))
		http.Error(w, "wrong email", http.StatusBadRequest)
		return
	}
	if user, _ := auth.Storage.GetUserByEmail(auth.Ctx, user.Email); user != nil {
		logger.Error("user already registred", zap.String("email", user.Email), zap.String("path", op), zap.Error(err))
		http.Error(w, "already registred", http.StatusBadRequest)
		return
	}
	//user.HashPassword not crypted
	hash, err := bcrypt.GenerateFromPassword([]byte(user.HashPassword), 12)
	if err != nil {
		logger.Error("err with crypt", zap.String("path", op), zap.Error(err))
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	id, err := auth.Storage.AddUser(auth.Ctx, user.Name, user.Email, string(hash))
	if err != nil {
		logger.Error("database add err", zap.String("path", op), zap.Error(err))
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	idStr := strconv.Itoa(id)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(idStr))
	logger.Info("succesfully register", zap.String("email", user.Email), zap.String("path", op))
}

func (auth *Auth) Login(w http.ResponseWriter, r *http.Request) {
	const op = "auth.Login"
	logger, _ := auth.Ctx.Value(("logger")).(*zap.Logger)

	var tmpuser models.User
	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		logger.Error("read body err", zap.String("path", op), zap.Error(err))
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	err = json.Unmarshal(body, &tmpuser)

	if err != nil {
		logger.Error("json unmarshal err", zap.String("path", op), zap.Error(err))
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	_, err = mail.ParseAddress(tmpuser.Email)
	if err != nil {
		logger.Error("wrong email", zap.String("email", tmpuser.Email), zap.String("path", op), zap.Error(err))
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	user, err := auth.Storage.GetUserByEmail(auth.Ctx, tmpuser.Email)
	if err != nil {
		logger.Info("user not found", zap.String("email", tmpuser.Email), zap.String("path", op), zap.Error(err))
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}
	if bcrypt.CompareHashAndPassword([]byte(user.HashPassword), []byte(tmpuser.HashPassword)) != nil {
		logger.Info("wrong password", zap.String("email", tmpuser.Email), zap.String("path", op))
		http.Error(w, "wrong password", http.StatusUnauthorized)
		return
	}
	access, refresh, err := jwt_utils.NewPairOfTokens(user.Id)
	if err != nil {
		logger.Error("access and refresh tokens err", zap.String("path", op), zap.Error(err))
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	auth.Cache.Add(auth.Ctx, user.Id, refresh)
	cookie := &http.Cookie{Name: "Authorization", Value: "Bearer " + access, HttpOnly: true}

	http.SetCookie(w, cookie)
	w.WriteHeader(http.StatusAccepted)
	w.Write([]byte("login succesfully"))
	logger.Info("succesfuly logged", zap.String("email", tmpuser.Email), zap.String("path", op))

}
func (auth *Auth) Logout(w http.ResponseWriter, r *http.Request) {
	const op = "auth.Logout"
	logger, _ := auth.Ctx.Value(("logger")).(*zap.Logger)
	cookie, err := r.Cookie("Authorization")
	if err != nil {
		http.Error(w, "wrong cookie", http.StatusUnauthorized)
		logger.Info("wrong cookie", zap.String("path", op))
		return
	}
	parts := strings.Split(cookie.Value, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		http.Error(w, "go to login", http.StatusUnauthorized)
		logger.Info("wrong authheader", zap.String("path", op))
		return
	}
	id, err := jwt_utils.GetIdFromToken(parts[1])
	if err != nil {
		http.Error(w, "wrong authtoken", http.StatusUnauthorized)
		logger.Info("wrong authtoken", zap.String("path", op))
		return
	}

	auth.Cache.Del(auth.Ctx, id)
	cookie = &http.Cookie{Name: "Authorization", Value: "", HttpOnly: true}
	http.SetCookie(w, cookie)
	w.Write([]byte("logout succesfully"))
	logger.Info("logout succesfully", zap.Int("id", id), zap.String("path", op))
}
func (auth *Auth) CheckAuth(w http.ResponseWriter, r *http.Request) {
	const op = "auth.CheckAuth"
	logger, _ := auth.Ctx.Value(("logger")).(*zap.Logger)
	cookie, err := r.Cookie("Authorization")
	if err != nil {
		http.Error(w, "wrong cookie", http.StatusUnauthorized)
		logger.Info("wrong cookie", zap.String("path", op))
		return
	}
	parts := strings.Split(cookie.Value, " ")
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		http.Error(w, "go to login", http.StatusUnauthorized)
		logger.Info("wrong authheader", zap.String("path", op))
		return
	}

	id, err := jwt_utils.GetIdFromToken(parts[1])
	idString := strconv.Itoa(id)
	refreshToken, _ := auth.Cache.Get(auth.Ctx, id)

	if err == jwt.ErrTokenExpired && refreshToken != "" {
		access, _ := jwt_utils.NewAccessToken(id)
		cookie := &http.Cookie{Name: "Authorization", Value: "Bearer " + access, HttpOnly: true}
		http.SetCookie(w, cookie)
	} else if err != nil {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		logger.Info("redirect to login", zap.String("path", op))
		return
	}
	w.Header().Set("X-User-Id", idString)
	w.Write([]byte("successful"))
	logger.Info("validate", zap.Int("id", id), zap.String("path", op))
}
