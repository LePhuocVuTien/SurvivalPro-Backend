package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/LePhuocVuTien/SurvivalPro-Backend/internal/db"
	"github.com/LePhuocVuTien/SurvivalPro-Backend/internal/models"
	"github.com/LePhuocVuTien/SurvivalPro-Backend/internal/utils"
	"github.com/gorilla/mux"
)

func GetUserChecklist(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int)

	rows, err := db.DB.Query(`SELECT id, title, category, is_checked, created_at FROM checklist_items WHERE user_id=$1`, userID)

	if err != nil {
		utils.ResponseJSON(w, http.StatusInternalServerError, map[string]string{
			"error": "Database error",
		})
		return
	}

	items := []models.CheckListItem{}
	for rows.Next() {
		var item models.CheckListItem

		if err := rows.Scan(
			&item.ID,
			&item.Title,
			&item.Category,
			&item.IsChecked,
			&item.CreatedAt,
		); err != nil {
			log.Println("Scan error:", err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		item.UserID = userID
		items = append(items, item)
	}

	utils.ResponseJSON(w, http.StatusOK, items)
}

func CreateChecklistItem(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int)
	var item models.CheckListItem
	json.NewDecoder(r.Body).Decode(&item)

	db.DB.QueryRow(`INSERT INTO checklist_items(user_id, title, category, is_checked) VALUES($1, $2, $3, $4) RETURNING id, created_at`,
		userID, item.Title, item.Category, item.IsChecked).Scan(&item.ID, &item.CreatedAt)

	item.UserID = userID
	utils.ResponseJSON(w, http.StatusOK, item)
}

func UpdateCheckListItem(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int)
	vars := mux.Vars(r)
	itemID := vars["id"]

	var item models.CheckListItem
	json.NewDecoder(r.Body).Decode(&item)

	db.DB.Exec(`UPDATE checklist_items SET title=$1, category=$2, is_checked=$3 WHERE id=$4 AND user_id=$5`,
		item.Title, item.Category, item.IsChecked, itemID, userID)

	utils.ResponseJSON(w, http.StatusOK, map[string]string{"message": "Item update"})
}

func DeleteChecklistItem(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int)
	id := mux.Vars(r)["id"]
	db.DB.Exec(`DELETE FROM checklist_items WHERE id=$1 AND user_id=$2`, id, userID)
	utils.ResponseJSON(w, http.StatusOK, map[string]string{"message": "Item deleted"})
}
