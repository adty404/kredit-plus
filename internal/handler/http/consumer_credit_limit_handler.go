package http

import (
	"github.com/adty404/kredit-plus/internal/usecase"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type ConsumerCreditLimitHandler struct {
	usecase         usecase.ConsumerCreditLimitUsecase
	consumerUsecase usecase.ConsumerUsecase
}

func NewConsumerCreditLimitHandler(
	usecase usecase.ConsumerCreditLimitUsecase,
	consumerUsecase usecase.ConsumerUsecase,
) ConsumerCreditLimitHandler {
	return ConsumerCreditLimitHandler{
		usecase:         usecase,
		consumerUsecase: consumerUsecase,
	}
}

func (h *ConsumerCreditLimitHandler) CreateLimitForConsumer(c *gin.Context) {
	idStr := c.Param("consumer_id")
	consumerID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid consumer ID format"})
		return
	}

	var input usecase.CreateConsumerCreditLimitInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data", "details": err.Error()})
		return
	}

	limit, err := h.usecase.CreateConsumerCreditLimit(uint(consumerID), input)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Credit limit created successfully", "data": limit})
}
