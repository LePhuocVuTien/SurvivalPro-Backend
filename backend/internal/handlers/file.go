package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/LePhuocVuTien/SurvivalPro-Backend/internal/db"
	"github.com/LePhuocVuTien/SurvivalPro-Backend/internal/utils"
)

const uploadDir = "internal/uploads"

func saveFile(file io.Reader, filename string) (string, error) {
	ext := filepath.Ext(filename)
	newFilename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	filePath := filepath.Join(uploadDir, newFilename)

	dst, err := os.Create(filePath)

	if err != nil {
		return "", err
	}
	defer dst.Close()
	io.Copy(dst, file)
	return "/uploads/" + newFilename, nil
}

func UploadAvatar(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int)
	r.ParseMultipartForm(5 << 20)
	file, header, err := r.FormFile("avatar")
	if err != nil {
		utils.ResponseJSON(w, http.StatusBadRequest, map[string]string{"error": "No file upload"})
		return
	}
	defer file.Close()

	contentType := header.Header.Get("Content-Type")
	if !strings.HasPrefix(contentType, "image/") {
		utils.ResponseJSON(w, http.StatusBadRequest, map[string]string{"error": "Only images allowed"})
		return
	}
	url, _ := saveFile(file, header.Filename)
	db.DB.Exec(`UPDATE users SET avatar_url=$1 WHERE id=$2`, url, userID)
	utils.ResponseJSON(w, http.StatusOK, map[string]string{"avatar_url": url})
}

func UploadGuideImage(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20)
	file, header, _ := r.FormFile("image")
	defer file.Close()
	url, _ := saveFile(file, header.Filename)
	utils.ResponseJSON(w, http.StatusOK, map[string]string{"image_url": url})
}
