package utils

import (
	"fmt"
	"net/http"
	"strconv"
)

// return from, to, error
func ParseQuery(r *http.Request) (int, int, error) {
	fromStr := r.URL.Query().Get("from")
	toStr := r.URL.Query().Get("to")
	if fromStr == "" || toStr == "" {
		return -1, -1, fmt.Errorf("empty from or to")
	}
	from, err := strconv.Atoi(fromStr)
	to, err := strconv.Atoi(toStr)
	if err != nil {
		return -1, -1, fmt.Errorf("not a number in variable(from or to)")
	}
	if from > to {
		return -1, -1, fmt.Errorf("from more than to")
	}
	return from, to, nil
}
