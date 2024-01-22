package http

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/arkreddy21/eligos"
	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strings"
	"time"
)

func (s *Server) authRoutes(r chi.Router) {
	r.Post("/login", s.handleLogin)
	r.Post("/register", s.handleRegister)
}

func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("unable to parse form"))
		return
	}
	email := r.Form.Get("email")
	password := r.Form.Get("password")
	if email == "" || password == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("provide all input fields"))
		return
	}
	user, err := s.UserService.GetUser(email)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("user not found"))
		return
	}
	if !CheckPasswordHash(password, user.Password) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("password incorrect"))
	}

	token, err := createToken(user.Id.String(), s.jwtKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("could not login"))
		return
	}
	response, err := json.Marshal(map[string]any{
		"message": "success",
		"jwt":     token,
	})
	w.Header().Set("Content-Type", "application/json")
	w.Write(response)
}

func (s *Server) handleRegister(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("unable to parse form"))
		return
	}
	name := r.Form.Get("name")
	email := r.Form.Get("email")
	password := r.Form.Get("password")
	if name == "" || email == "" || password == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("provide all input fields"))
		return
	}
	hashedPassword, _ := HashPassword(password)
	user := &eligos.User{
		Name:     name,
		Email:    email,
		Password: hashedPassword,
	}
	err = s.UserService.CreateUser(user)
	if err != nil {
		fmt.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Unable to create user. Check if email is already registered"))
		return
	}
	w.Write([]byte("register successful"))
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func createToken(issuer string, key []byte) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    issuer,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24 * 100)),
	})
	token, err := claims.SignedString(key)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (s *Server) validateJwt(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		const BearerSchema = "Bearer "
		authHeader := r.Header.Get("Authorization")
		if !strings.HasPrefix(authHeader, BearerSchema) {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("user unauthorized"))
			return
		}
		authToken := authHeader[len(BearerSchema):]
		token, err := jwt.ParseWithClaims(authToken, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
			return s.jwtKey, nil
		})
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("user unauthorized"))
			return
		}
		claims := token.Claims.(*jwt.RegisteredClaims)
		ctx := context.WithValue(r.Context(), "userId", claims.Issuer)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
