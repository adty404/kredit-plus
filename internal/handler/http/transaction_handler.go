package http

import (
	"github.com/adty404/kredit-plus/internal/usecase"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type TransactionHandler struct {
	uc usecase.TransactionUsecase
}

func NewTransactionHandler(uc usecase.TransactionUsecase) *TransactionHandler {
	return &TransactionHandler{uc: uc}
}

// CreateTransaction menangani pembuatan transaksi baru untuk seorang konsumen.
// @Summary      Create a new transaction for a consumer
// @Description  Membuat transaksi kredit baru berdasarkan OTR, tenor, dan detail lainnya.
// @Tags         transactions
// @Accept       json
// @Produce      json
// @Param        consumer_id  path      int                          true  "Consumer ID"
// @Param        transaction  body      usecase.CreateTransactionInput  true  "Detail Transaksi"
// @Success      201          {object}  map[string]interface{}
// @Failure      400          {object}  map[string]interface{}
// @Failure      404          {object}  map[string]interface{}
// @Failure      422          {object}  map[string]interface{}
// @Router       /consumers/{consumer_id}/transactions [post]
func (h *TransactionHandler) CreateTransaction(c *gin.Context) {
	idStr := c.Param("id")
	consumerID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid consumer ID format"})
		return
	}

	var input usecase.CreateTransactionInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data", "details": err.Error()})
		return
	}

	transaction, err := h.uc.CreateTransaction(uint(consumerID), input)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Transaction created successfully", "data": transaction})
}

// GetTransactionsByConsumerID mengambil semua transaksi dari seorang konsumen.
// @Summary      Get all transactions for a consumer
// @Description  Mengambil daftar semua transaksi yang pernah dilakukan oleh seorang konsumen.
// @Tags         transactions
// @Produce      json
// @Param        consumer_id  path      int  true  "Consumer ID"
// @Success      200          {object}  map[string]interface{}
// @Failure      404          {object}  map[string]interface{}
// @Router       /consumers/{consumer_id}/transactions [get]
func (h *TransactionHandler) GetTransactionsByConsumerID(c *gin.Context) {
	idStr := c.Param("id")
	consumerID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid consumer ID format"})
		return
	}

	transactions, err := h.uc.GetTransactionsByConsumerID(uint(consumerID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": transactions})
}
