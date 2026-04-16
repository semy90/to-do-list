package transport

import (
	"auth/internal/database"
	"auth/internal/models"
	"auth/jwt"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

type Auth struct {
	Cache   database.RefreshCache
	Storage database.UserStorage
}

// todo валидацию email + проверка на повторную регу
func (auth *Auth) Register(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	var user models.User
	json.Unmarshal(body, &user)
	//user.HashPassword not crypted
	hash, err := bcrypt.GenerateFromPassword([]byte(user.HashPassword), 12)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	id, err := auth.Storage.AddUser(user.Name, user.Email, string(hash))
	if err != nil {
		fmt.Println(err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	idStr := strconv.Itoa(id)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(idStr))
}

func (auth *Auth) Login(w http.ResponseWriter, r *http.Request) {
	var tmpuser models.User
	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	err = json.Unmarshal(body, &tmpuser)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	user, err := auth.Storage.GetUserByEmail(tmpuser.Email)
	if err != nil {
		http.Error(w, "user not found", http.StatusNotFound)
		return
	}
	if bcrypt.CompareHashAndPassword([]byte(tmpuser.HashPassword), []byte(user.HashPassword)) != nil {
		http.Error(w, "wrong password", http.StatusUnauthorized)
		return
	}
	access, refresh, err := jwt.NewPairOfTokens(user.Id)
	if err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	auth.Cache.Add(user.Id, refresh)
	cookie := &http.Cookie{Name: "Authorization", Value: "Bearer " + access, HttpOnly: true}

	http.SetCookie(w, cookie)
	w.Write([]byte("you logged succesfuly"))
	w.WriteHeader(http.StatusAccepted)
}
func (auth *Auth) Logout(w http.ResponseWriter, r *http.Request) {}
