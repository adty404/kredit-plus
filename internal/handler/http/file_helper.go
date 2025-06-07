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

	// Buat nama file yang unik untuk menghindari konflik.
	// Contoh: 3271011505900001-ktp-1622818800.jpg
	ext := filepath.Ext(file.Filename)
	filename := fmt.Sprintf("%s-%s-%d%s", nik, fileType, time.Now().Unix(), ext)

	// Tentukan path penyimpanan. Pastikan folder 'uploads' sudah dibuat.
	uploadPath := filepath.Join("uploads", filename)

	// Simpan file ke path tersebut.
	if err := c.SaveUploadedFile(file, uploadPath); err != nil {
		return "", err
	}

	// Kembalikan path yang bisa disimpan ke database.
	return uploadPath, nil
}
