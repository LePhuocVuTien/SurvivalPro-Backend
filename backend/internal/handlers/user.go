package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/LePhuocVuTien/SurvivalPro-Backend/internal/db"
	"github.com/LePhuocVuTien/SurvivalPro-Backend/internal/models"
	"github.com/LePhuocVuTien/SurvivalPro-Backend/internal/utils"
	"golang.org/x/crypto/bcrypt"
)

func RegisterUser(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		utils.ResponseJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request"})
		return
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	query := `INSERT INTO users (email, password, name) VALUES ($1, $2, $3) RETURNING id, created_at`
	err := db.DB.QueryRow(query, user.Email, string(hashedPassword), user.Name).Scan(&user.ID, &user.CreatedAt)

	if err != nil {
		utils.ResponseJSON(w, http.StatusInternalServerError, map[string]string{"error": "Email already exists"})
		return
	}

	token, _ := utils.GenerateJWT(user.ID, user.Email)
	user.Password = ""
	utils.ResponseJSON(w, http.StatusCreated, map[string]interface{}{
		"user":  user,
		"token": token,
	})
}

func LoginUser(w http.ResponseWriter, r *http.Request) {
	var credentials struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	json.NewDecoder(r.Body).Decode(&credentials)

	var user models.User
	var hashedPassword string

	query := `SELECT id, email, password, name, created_at FROM users WHERE email = $1`
	err := db.DB.QueryRow(query, credentials.Email).Scan(&user.ID, &user.Email, &hashedPassword, &user.Name, &user.CreatedAt)
	if err != nil {
		utils.ResponseJSON(w, http.StatusUnauthorized, map[string]string{"error": "Invalid email or password"})
		return
	}

	// Compare passwords
	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(credentials.Password))
	if err != nil {
		utils.ResponseJSON(w, http.StatusUnauthorized, map[string]string{"error": "Invalid credentials"})
		return
	}

	// Generate token
	token, _ := utils.GenerateJWT(user.ID, user.Email)

	user.Password = ""
	utils.ResponseJSON(w, http.StatusOK, map[string]interface{}{
		"user":  user,
		"token": token,
	})
}

func UpdateUserProfile(w http.ResponseWriter, r *http.Request) {
	UserID := r.Context().Value("userID").(int)
	var user models.User
	json.NewDecoder(r.Body).Decode(&user)
	query := `UPDATE users SET name=$1, blood_type=$2, allergies=$3, emergency=$4 WHERE id=$5`
	db.DB.Exec(query, user.Name, user.BloodType, user.Allergies, user.EmergencyNum, UserID)

	utils.ResponseJSON(w, http.StatusOK, map[string]string{"message": "Profile updated"})
}
