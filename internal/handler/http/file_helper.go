package http

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"mime/multipart"
	"path/filepath"
	"time"
)

func SaveUploadedFile(c *gin.Context, file *multipart.FileHeader, nik, fileType string) (string, error) {
	if file == nil {
		return "", nil // Tidak ada file yang di-upload, ini bukan error.
	}

	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%s-%s-%d%s", nik, fileType, time.Now().Unix(), ext)

	uploadPath := filepath.Join("uploads", filename)

	if err := c.SaveUploadedFile(file, uploadPath); err != nil {
		return "", err
	}

	return uploadPath, nil
}
