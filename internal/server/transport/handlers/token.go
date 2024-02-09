package handlers

import (
	"GoPass/internal/server/users"
	"errors"
	"net/http"

	"github.com/dgrijalva/jwt-go"
)

var jwtSecret = []byte("super-secret-key")

func createAuthToken(user *users.User) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = user.ID
	claims["username"] = user.Login

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func getTokenFromRequest(r *http.Request) (string, error) {
	cookie, err := r.Cookie("auth_token")
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

func getCurrentUser(r *http.Request) (int, error) {
	tokenString, err := getTokenFromRequest(r)
	if err != nil {
		return 0, err
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("wrong token method")
		}
		return jwtSecret, nil
	})
	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, errors.New("claims error")
	}

	userID, ok := claims["user_id"].(float64)
	if !ok {
		return 0, errors.New("error getting user id")
	}

	return int(userID), nil
}
