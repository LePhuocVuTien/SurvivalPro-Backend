package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/LePhuocVuTien/SurvivalPro-Backend/internal/models"
	"github.com/LePhuocVuTien/SurvivalPro-Backend/internal/redis"
	"github.com/LePhuocVuTien/SurvivalPro-Backend/internal/utils"
)

func GetCurrentWeather(w http.ResponseWriter, r *http.Request) {
	latStr := r.URL.Query().Get("lat")
	lonStr := r.URL.Query().Get("lon")
	if latStr == "" || lonStr == "" {
		utils.ResponseJSON(w, http.StatusBadRequest, map[string]string{"error": "lat and lon required"})
		return
	}

	lat, _ := strconv.ParseFloat(latStr, 64)
	lon, _ := strconv.ParseFloat(lonStr, 64)

	cacheKey := fmt.Sprintf("weather:%f:%f", lat, lon)
	cached, err := redis.CacheGet(cacheKey)

	if err == nil {
		var weather models.WeatherResponse
		json.Unmarshal([]byte(cached), &weather)
		utils.ResponseJSON(w, http.StatusOK, weather)
		return
	}

	apiKey := os.Getenv("OPENWEATHER_API_KEY")
	if apiKey == "" {
		utils.ResponseJSON(w, http.StatusBadRequest, map[string]string{"error": "Key weather invalid"})
		return
	}

	url := fmt.Sprintf("https://api.openweathermap.org/data/2.5/weather?lat=%f&lon=%f&appid=%s&units=metric&lang=vi", lat, lon, apiKey)
	resp, err := http.Get(url)
	if err != nil {
		log.Println("Error fetching URL:", err)
	}
	defer resp.Body.Close()

	var weather models.WeatherResponse
	json.NewDecoder(resp.Body).Decode(&weather)
	redis.CacheSet(cacheKey, weather, 10*60*1e9)

	utils.ResponseJSON(w, http.StatusOK, weather)
}
