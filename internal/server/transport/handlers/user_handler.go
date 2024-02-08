package handlers

import (
	"GoPass/internal/server"
	mw "GoPass/internal/server/transport/middleware"
	"GoPass/internal/server/users"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type UserHandler struct {
	UseCase users.UseCase
}

type Auth struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

func NewUserHandler(useCase users.UseCase) *UserHandler {
	return &UserHandler{
		UseCase: useCase,
	}
}

func (u *UserHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	var req Auth

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		mw.LogError(w, r, err)
		http.Error(w, "Failed to decode JSON request", http.StatusBadRequest)
		return
	}

	username := req.Login
	password := req.Password
	fmt.Println(username)

	if username == "" || password == "" {
		mw.LogError(w, r, server.ErrEmptyField)
		http.Error(w, "Username and password must not be empty", http.StatusBadRequest)
		return
	}

	user, err := u.UseCase.AuthenticateUser(r.Context(), username, password)
	if err != nil {
		mw.LogError(w, r, err)
		http.Error(w, "Authentication failed", http.StatusUnauthorized)
		return
	}

	token, err := createAuthToken(user)
	if err != nil {
		mw.LogError(w, r, err)
		http.Error(w, "Failed to create token", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:  "auth_token",
		Value: token,
	})

	w.WriteHeader(http.StatusOK)
}

func (u *UserHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var req Auth

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		mw.LogError(w, r, err)
		http.Error(w, "Failed to decode JSON request", http.StatusBadRequest)
		return
	}

	username := req.Login
	password := req.Password

	if username == "" || password == "" {
		mw.LogError(w, r, server.ErrEmptyField)
		http.Error(w, "Username and password must not be empty", http.StatusBadRequest)
		return
	}

	user, err := u.UseCase.RegisterUser(r.Context(), username, password)
	if err != nil {
		if errors.Is(err, server.ErrUserExists) {
			mw.LogError(w, r, err)
			http.Error(w, "Username is already taken", http.StatusConflict)
			return
		}
		http.Error(w, "Registration failed", http.StatusBadRequest)
		return
	}

	token, err := createAuthToken(user)
	if err != nil {
		mw.LogError(w, r, err)
		http.Error(w, "Failed to create token", http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:  "auth_token",
		Value: token,
	})

	w.WriteHeader(http.StatusOK)
}
