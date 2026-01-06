package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/LePhuocVuTien/SurvivalPro-Backend/internal/db"
	"github.com/LePhuocVuTien/SurvivalPro-Backend/internal/utils"
	"github.com/LePhuocVuTien/SurvivalPro-Backend/models"
)

func RegisterPushToken(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int)

	var data struct {
		PushToken string ` json:"push_token"`
	}

	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		utils.ResponseJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request"})
		return
	}

	if data.PushToken == "" {
		utils.ResponseJSON(w, http.StatusBadRequest, map[string]string{"error": "Push token is required"})
		return
	}

	_, err := db.DB.Exec(
		"UPDATE users SET push_token = $1 WHERE id = $2",
		data.PushToken,
		userID,
	)

	if err != nil {
		utils.ResponseJSON(w, http.StatusInternalServerError, map[string]string{"error": " Error registering token"})
		return
	}

	utils.ResponseJSON(w, http.StatusOK, map[string]string{"message": "Push token register successfully"})

}

func GetUserNotifications(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int)
	rows, err := db.DB.Query(`SELECT id, title, body, data, sent, created_at FROM notifications WHERE user_id=$1 ORDER BY created_at DESC LIMIT 50`, userID)

	if err != nil {
		utils.ResponseJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "Database error",
		})
		return
	}

	notifications := []models.PushNotification{}
	for rows.Next() {
		var n models.PushNotification
		rows.Scan(&n.ID, &n.Body, &n.Data, &n.Sent, &n.CreatedAt)
		n.UserID = userID
		notifications = append(notifications, n)
	}
	utils.ResponseJSON(w, http.StatusOK, notifications)
}
