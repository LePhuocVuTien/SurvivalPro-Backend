package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/LePhuocVuTien/SurvivalPro-Backend/internal/db"
	"github.com/LePhuocVuTien/SurvivalPro-Backend/internal/utils"
	"github.com/LePhuocVuTien/SurvivalPro-Backend/models"
	"github.com/LePhuocVuTien/SurvivalPro-Backend/redis"
)

func GetAllGuides(w http.ResponseWriter, r *http.Request) {
	category := r.URL.Query().Get("category")
	search := r.URL.Query().Get("search")

	cacheKey := fmt.Sprintf("guide:%s:%s", category, search)

	cached, err := redis.CacheGet(cacheKey)

	if err == nil {
		var guides []models.SurvivalGuide
		if err := json.Unmarshal([]byte(cached), &guides); err == nil {
			utils.ResponseJSON(w, http.StatusOK, guides)
			return
		}
	}

	query := `
		SELECT id, title, category, difficulty, icon, content, image_url, views, created_at
		FROM survival_guides
		WHERE 1=1
	`
	args := []interface{}{}
	argIdx := 1

	if category != "" && category != "Tất cả" {
		query += fmt.Sprintf(" AND category = $%d", argIdx)
		args = append(args, category)
		argIdx++
	}

	if search != "" {
		query += fmt.Sprintf(" AND title ILIKE $%d", argIdx)
		args = append(args, "%"+search+"%")
		argIdx++
	}

	query += " ORDER BY views DESC, created_at DESC"

	rows, err := db.DB.Query(query, args...)
	if err != nil {
		utils.ResponseJSON(w, http.StatusInternalServerError, map[string]string{"error": "Error fetching guides"})
		return
	}

	guides := []models.SurvivalGuide{}

	for rows.Next() {
		var guide models.SurvivalGuide
		if err := rows.Scan(
			&guide.ID,
			&guide.Title,
			&guide.Category,
			&guide.Difficulty,
			&guide.Icon,
			&guide.Content,
			&guide.ImageURL,
			&guide.Views,
			&guide.CreatedAt,
		); err != nil {
			continue
		}
		guides = append(guides, guide)
	}

	redis.CacheSet(cacheKey, guides, 5*time.Minute)
	utils.ResponseJSON(w, http.StatusOK, guides)
}
