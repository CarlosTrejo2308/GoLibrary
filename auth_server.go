package main

import (
	"crypto/sha1"
	"fmt"
	"net/http"
	"time"
)

var Tokens = make(map[string]time.Time)

func main() {
	http.Handle("/", http.HandlerFunc(Auth))

	http.ListenAndServe(":8085", nil)
}

func Auth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case http.MethodGet:
		getToken(w, r)
	case http.MethodPost:
		postToken(w, r)
	}

}

func generateToken(uid string) string {
	sha := sha1.New()
	sha.Write([]byte(uid + time.Now().String()))

	token := fmt.Sprintf("%x", sha.Sum(nil))
	return token
}

func postToken(w http.ResponseWriter, r *http.Request) {
	uid := r.Header.Get("X_UID")
	pwd := r.Header.Get("X_PWD")

	if uid == "12" && pwd == "secret" {
		token := generateToken(uid)
		Tokens[token] = time.Now()
		w.Write([]byte(token))
	} else {
		w.WriteHeader(http.StatusUnauthorized)
	}
}

func getToken(w http.ResponseWriter, r *http.Request) {
	token := r.Header.Get("X_TOKEN")

	if token == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	ts, ok := Tokens[token]
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if time.Now().Sub(ts) > (3 * time.Minute) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	} else {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}
}
