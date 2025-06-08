package http

import (
	"errors"
	"github.com/gin-gonic/gin/binding"
	"net/http"
	"strconv"

	"github.com/adty404/kredit-plus/internal/usecase"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ConsumerHandler struct {
	consumerUsecase usecase.ConsumerUsecase
}

func NewConsumerHandler(uc usecase.ConsumerUsecase) *ConsumerHandler {
	return &ConsumerHandler{consumerUsecase: uc}
}

func (h *ConsumerHandler) CreateConsumer(c *gin.Context) {
	var input usecase.CreateConsumerFormInput
	if err := c.ShouldBind(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed", "details": err.Error()})
		return
	}

	// Simpan file yang di-upload.
	fotoKtpPath, err := SaveUploadedFile(c, input.FotoKtp, input.Nik, "ktp")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not save foto_ktp file", "details": err.Error()})
		return
	}
	fotoSelfiePath, err := SaveUploadedFile(c, input.FotoSelfie, input.Nik, "selfie")
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "Could not save foto_selfie file", "details": err.Error()},
		)
		return
	}

	// Konversi nilai string dari form ke float64.
	gaji, _ := strconv.ParseFloat(input.Gaji, 64)
	overallCreditLimit, _ := strconv.ParseFloat(input.OverallCreditLimit, 64)

	// Siapkan input untuk usecase.
	usecaseInput := usecase.CreateConsumerInput{
		Nik:                input.Nik,
		FullName:           input.FullName,
		LegalName:          input.LegalName,
		Email:              input.Email,
		Password:           input.Password,
		TempatLahir:        input.TempatLahir,
		TanggalLahir:       input.TanggalLahir,
		Gaji:               gaji,
		OverallCreditLimit: overallCreditLimit,
		FotoKtpPath:        fotoKtpPath,
		FotoSelfiePath:     fotoSelfiePath,
	}

	// Panggil usecase.
	consumer, err := h.consumerUsecase.CreateConsumer(usecaseInput)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Consumer and user created successfully", "data": consumer})
}

// GetAllConsumers menangani permintaan untuk mendapatkan semua konsumen.
func (h *ConsumerHandler) GetAllConsumers(c *gin.Context) {
	consumers, err := h.consumerUsecase.GetAllConsumers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve consumers"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": consumers})
}

func (h *ConsumerHandler) GetConsumerByID(c *gin.Context) {
	// Ambil ID dari URL
	idStr := c.Param("id")
	id, _ := strconv.ParseUint(idStr, 10, 32)

	// Ambil data user yang login dari context
	loggedInUserID := c.GetUint("userID")
	loggedInUserRole := c.GetString("userRole")

	// --- LOGIKA KONTROL AKSES ---
	if loggedInUserRole != "admin" {
		consumer, err := h.consumerUsecase.GetConsumerByUserID(loggedInUserID)
		if err != nil || consumer.ID != uint(id) {
			c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to view this consumer"})
			return
		}
	}

	consumer, err := h.consumerUsecase.GetConsumerByID(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Consumer not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve consumer"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": consumer})
}

func (h *ConsumerHandler) UpdateConsumer(c *gin.Context) {
	// ambil ID dari URL parameter
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid consumer ID format"})
		return
	}

	// Ambil data user yang login dari context
	loggedInUserID := c.GetUint("userID")
	loggedInUserRole := c.GetString("userRole")

	// --- LOGIKA KONTROL AKSES ---
	if loggedInUserRole != "admin" {
		// Cari data consumer berdasarkan user ID yang login
		consumer, err := h.consumerUsecase.GetConsumerByUserID(loggedInUserID)
		if err != nil || consumer.ID != uint(id) {
			c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to view this consumer"})
			return
		}
	}

	// Bind input JSON ke struct UpdateConsumerInput
	var input usecase.UpdateConsumerInput
	if err := c.ShouldBindWith(&input, binding.JSON); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Validation failed", "details": err.Error()})
		return
	}

	// Update Consumer
	consumer, err := h.consumerUsecase.UpdateConsumer(uint(id), input)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Consumer not found"})
			return
		}
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Consumer updated successfully", "data": consumer})
}

// DeleteConsumer menangani penghapusan konsumen.
func (h *ConsumerHandler) DeleteConsumer(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid consumer ID format"})
		return
	}

	err = h.consumerUsecase.DeleteConsumer(uint(id))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Consumer not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete consumer"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Consumer deleted successfully"})
}
