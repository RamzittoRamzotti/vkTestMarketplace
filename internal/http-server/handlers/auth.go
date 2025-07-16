package handlers

import (
	"encoding/json"
	"net/http"
	"vkTestMarketplace/internal/lib/jwt"
	"vkTestMarketplace/internal/models"
	"vkTestMarketplace/internal/storage"
	utilsauth "vkTestMarketplace/internal/utils/auth"
	"vkTestMarketplace/internal/utils/logger"
)

type UserStorage interface {
	CreateUser(user *models.User) (int64, error)
	GetUserByUsername(username string) (*models.User, error)
	GetUserByID(id int) (*models.User, error)
}

type AuthHandlers struct {
	Users     UserStorage
	JWTSecret string
}

func (h *AuthHandlers) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Login    string `json:"login"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	if err := utilsauth.ValidateUsername(req.Login); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := utilsauth.ValidatePassword(req.Password); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	hash, err := utilsauth.HashPassword(req.Password)
	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	user := &models.User{
		Login:        req.Login,
		PasswordHash: hash,
	}
	var id int64
	if id, err = h.Users.CreateUser(user); err != nil {
		if err == storage.ErrUserExists {
			logger.Warn("Registration failed: user exists: %s", req.Login)
			http.Error(w, "user already exists", http.StatusConflict)
			return
		}
		logger.Error("Registration failed: %v", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	logger.Info("User registered: %s", user.Login)
	user.PasswordHash = ""
	user.ID = int(id)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func (h *AuthHandlers) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"login"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}
	user, err := h.Users.GetUserByUsername(req.Username)
	if err != nil {
		if err == storage.ErrUserNotFound {
			logger.Warn("Login failed: user not found: %s", req.Username)
			http.Error(w, "invalid credentials", http.StatusUnauthorized)
			return
		}
		logger.Error("Login failed: %v", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	if err := utilsauth.CheckPasswordHash(req.Password, user.PasswordHash); err != nil {
		logger.Warn("Login failed: wrong password for user: %s", req.Username)
		http.Error(w, "invalid credentials", http.StatusUnauthorized)
		return
	}
	token, err := jwt.GenerateToken(user.ID, h.JWTSecret)
	if err != nil {
		logger.Error("Token generation failed: %v", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}
	logger.Info("User logged in: %s", user.Login)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}
