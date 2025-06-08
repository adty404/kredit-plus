package http

import (
	"github.com/adty404/kredit-plus/internal/domain"
	"github.com/adty404/kredit-plus/internal/usecase"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// TransactionHandler sekarang memiliki dependensi ke consumerRepo untuk validasi akses.
type TransactionHandler struct {
	uc           usecase.TransactionUsecase
	consumerRepo domain.ConsumerRepository // Dependensi baru
}

// NewTransactionHandler diperbarui untuk menerima dependensi baru.
func NewTransactionHandler(uc usecase.TransactionUsecase, consumerRepo domain.ConsumerRepository) *TransactionHandler {
	return &TransactionHandler{
		uc:           uc,
		consumerRepo: consumerRepo,
	}
}

// CreateTransaction sekarang memiliki validasi kontrol akses.
func (h *TransactionHandler) CreateTransaction(c *gin.Context) {
	// Ambil consumer ID dari parameter URL.
	idStr := c.Param("id")
	consumerIDFromURL, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid consumer ID format"})
		return
	}

	// Ambil data pengguna yang login dari context (di-set oleh middleware).
	loggedInUserID := c.GetUint("userID")
	loggedInUserRole := c.GetString("userRole")

	// --- VALIDASI KONTROL AKSES ---
	// Jika yang login bukan admin, pastikan dia hanya membuat transaksi untuk dirinya sendiri.
	if loggedInUserRole != "admin" {
		// Cari profil consumer yang terhubung dengan user yang sedang login.
		consumer, err := h.consumerRepo.FindByUserID(loggedInUserID)
		if err != nil || consumer.ID != uint(consumerIDFromURL) {
			c.JSON(
				http.StatusForbidden,
				gin.H{"error": "You are not authorized to create a transaction for this consumer"},
			)
			return
		}
	}
	// -----------------------------

	// Jika validasi lolos, lanjutkan proses...
	var input usecase.CreateTransactionInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data", "details": err.Error()})
		return
	}

	transaction, err := h.uc.CreateTransaction(uint(consumerIDFromURL), input)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Transaction created successfully", "data": transaction})
}

// GetTransactionsByConsumerID juga memerlukan validasi kontrol akses.
func (h *TransactionHandler) GetTransactionsByConsumerID(c *gin.Context) {
	idStr := c.Param("id")
	consumerIDFromURL, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid consumer ID format"})
		return
	}

	loggedInUserID := c.GetUint("userID")
	loggedInUserRole := c.GetString("userRole")

	if loggedInUserRole != "admin" {
		consumer, err := h.consumerRepo.FindByUserID(loggedInUserID)
		if err != nil || consumer.ID != uint(consumerIDFromURL) {
			c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to view these transactions"})
			return
		}
	}

	transactions, err := h.uc.GetTransactionsByConsumerID(uint(consumerIDFromURL))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": transactions})
}
