package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/LePhuocVuTien/SurvivalPro-Backend/internal/db"
	"github.com/LePhuocVuTien/SurvivalPro-Backend/internal/models"
	"github.com/LePhuocVuTien/SurvivalPro-Backend/internal/redis"
	"github.com/LePhuocVuTien/SurvivalPro-Backend/internal/utils"
)

func SaveUserLocation(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int)

	var location models.UserLocation
	if err := json.NewDecoder(r.Body).Decode(&location); err != nil {
		utils.ResponseJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request"})
		return
	}

	db.DB.QueryRow(`INSERT INTO user_locations(user_id, latitude, longtitude) VALUES($1, $2, $3) RETURNING id, timestamp`,
		userID, location.Latitude, location.Longtitude).Scan(&location.ID, &location.Timestamp)

	location.UserID = userID
	cacheKey := fmt.Sprintf("location:%d", userID)
	redis.CacheSet(cacheKey, location, 5*time.Minute)
	utils.ResponseJSON(w, http.StatusOK, location)

}

func GetUserLocation(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int)

	cacheKey := fmt.Sprintf("location:%d", userID)
	cached, errCache := redis.CacheGet(cacheKey)

	if errCache == nil {
		var cachedLocation models.UserLocation
		json.Unmarshal([]byte(cached), &cachedLocation)
		utils.ResponseJSON(w, http.StatusOK, cachedLocation)
		return
	}

	var loc models.UserLocation
	err := db.DB.QueryRow(`
	SELECT id, user_id, latitude, longtitude, timestamp
	FROM user_locations
	WHERE user_id=$1
	ORDER BY timestamp DESC
	LIMIT 1
	`, userID).Scan(
		&loc.ID,
		&loc.UserID,
		&loc.Latitude,
		&loc.Longtitude,
		&loc.Timestamp,
	)

	if err == sql.ErrNoRows {
		utils.ResponseJSON(w, http.StatusNotFound, map[string]string{"error": "Location not found"})
		return
	}

	redis.CacheSet(cacheKey, loc, 5*time.Minute)

	utils.ResponseJSON(w, http.StatusOK, loc)
}
