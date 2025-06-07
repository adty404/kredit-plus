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

// ConsumerHandler menangani permintaan HTTP yang berhubungan dengan konsumen.
type ConsumerHandler struct {
	consumerUsecase usecase.ConsumerUsecase
}

// NewConsumerHandler adalah factory function untuk membuat instance ConsumerHandler.
func NewConsumerHandler(uc usecase.ConsumerUsecase) *ConsumerHandler {
	return &ConsumerHandler{consumerUsecase: uc}
}

// CreateConsumer menangani pembuatan konsumen baru dengan upload file.
// @Summary      Create a new consumer with file uploads
// @Description  Membuat konsumen baru dengan data form dan upload file KTP serta Selfie.
// @Tags         consumers
// @Accept       multipart/form-data
// @Produce      json
// @Param        nik           formData  string                true  "Nomor Induk Kependudukan (NIK)"
// @Param        full_name     formData  string                true  "Nama Lengkap"
// @Param        legal_name    formData  string                true  "Nama Sesuai KTP"
// @Param        tempat_lahir  formData  string                true  "Tempat Lahir"
// @Param        tanggal_lahir formData  string                true  "Tanggal Lahir (Format: YYYY-MM-DD)"
// @Param        gaji          formData  number                true  "Gaji per bulan"
// @Param        foto_ktp      formData  file                  false "File Foto KTP"
// @Param        foto_selfie   formData  file                  false "File Foto Selfie"
// @Success      201           {object}  map[string]interface{}
// @Failure      400           {object}  map[string]interface{}
// @Failure      500           {object}  map[string]interface{}
// @Router       /consumers [post]
func (h *ConsumerHandler) CreateConsumer(c *gin.Context) {
	// Dapatkan data teks dari form-data
	nik := c.PostForm("nik")
	fullName := c.PostForm("full_name")
	legalName := c.PostForm("legal_name")
	tempatLahir := c.PostForm("tempat_lahir")
	tanggalLahir := c.PostForm("tanggal_lahir")
	gajiStr := c.PostForm("gaji")

	// Validasi input sederhana
	if nik == "" || fullName == "" || legalName == "" || tempatLahir == "" || tanggalLahir == "" || gajiStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing required text fields"})
		return
	}

	// Dapatkan file dari form-data (opsional)
	fotoKtp, _ := c.FormFile("foto_ktp")
	fotoSelfie, _ := c.FormFile("foto_selfie")

	// Simpan file yang di-upload dan dapatkan path-nya
	fotoKtpPath, err := SaveUploadedFile(c, fotoKtp, nik, "ktp")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not save foto_ktp file", "details": err.Error()})
		return
	}

	fotoSelfiePath, err := SaveUploadedFile(c, fotoSelfie, nik, "selfie")
	if err != nil {
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "Could not save foto_selfie file", "details": err.Error()},
		)
		return
	}

	// Konversi gaji dari string ke float64
	gaji, err := strconv.ParseFloat(gajiStr, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid format for gaji, must be a number"})
		return
	}

	// Siapkan input untuk usecase
	input := usecase.CreateConsumerInput{
		Nik:            nik,
		FullName:       fullName,
		LegalName:      legalName,
		TempatLahir:    tempatLahir,
		TanggalLahir:   tanggalLahir,
		Gaji:           gaji,
		FotoKtpPath:    fotoKtpPath,
		FotoSelfiePath: fotoSelfiePath,
	}

	// Panggil usecase
	consumer, err := h.consumerUsecase.CreateConsumer(input)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Consumer created successfully", "data": consumer})
}

// GetAllConsumers menangani permintaan untuk mendapatkan semua konsumen.
// @Summary      Get all consumers
// @Description  Mengambil daftar semua konsumen beserta data relasinya.
// @Tags         consumers
// @Produce      json
// @Success      200  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /consumers [get]
func (h *ConsumerHandler) GetAllConsumers(c *gin.Context) {
	consumers, err := h.consumerUsecase.GetAllConsumers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve consumers"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": consumers})
}

// GetConsumerByID menangani permintaan untuk mendapatkan satu konsumen berdasarkan ID.
// @Summary      Get a consumer by ID
// @Description  Mengambil detail satu konsumen berdasarkan ID.
// @Tags         consumers
// @Produce      json
// @Param        id   path      int  true  "Consumer ID"
// @Success      200  {object}  map[string]interface{}
// @Failure      400  {object}  map[string]interface{}
// @Failure      404  {object}  map[string]interface{}
// @Router       /consumers/{id} [get]
func (h *ConsumerHandler) GetConsumerByID(c *gin.Context) {
	// Mengambil ID dari URL parameter.
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid consumer ID format"})
		return
	}

	consumer, err := h.consumerUsecase.GetConsumerByID(uint(id))
	if err != nil {
		// Cek apakah error karena data tidak ditemukan.
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Consumer not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve consumer"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": consumer})
}

// UpdateConsumer menangani pembaruan data konsumen.
// @Summary      Update a consumer
// @Description  Memperbarui data konsumen yang ada berdasarkan ID.
// @Tags         consumers
// @Accept       json
// @Produce      json
// @Param        id        path      int                      true  "Consumer ID"
// @Param        consumer  body      usecase.UpdateConsumerInput  true  "Data Pembaruan Konsumen"
// @Success      200       {object}  map[string]interface{}
// @Failure      400       {object}  map[string]interface{}
// @Failure      404       {object}  map[string]interface{}
// @Failure      422       {object}  map[string]interface{}
// @Router       /consumers/{id} [put]
func (h *ConsumerHandler) UpdateConsumer(c *gin.Context) {
	// ambil ID dari URL parameter
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid consumer ID format"})
		return
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
// @Summary      Delete a consumer
// @Description  Menghapus data konsumen berdasarkan ID.
// @Tags         consumers
// @Produce      json
// @Param        id   path      int  true  "Consumer ID"
// @Success      200  {object}  map[string]string
// @Failure      404  {object}  map[string]interface{}
// @Failure      500  {object}  map[string]interface{}
// @Router       /consumers/{id} [delete]
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
