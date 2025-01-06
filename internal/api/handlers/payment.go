package handlers

import (
	"net/http"
	"time"

	"oxo-game-api/internal/models"
	"oxo-game-api/pkg/utils/response"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PaymentHandler struct {
	db *gorm.DB
}

func NewPaymentHandler(db *gorm.DB) *PaymentHandler {
	return &PaymentHandler{db: db}
}

// ProcessPayment godoc
// @Summary Process a payment
// @Description Process a payment with the specified method and amount
// @Tags payments
// @Accept json
// @Produce json
// @Param payment body models.Payment true "Payment information"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 402 {object} response.PaymentError
// @Failure 500 {object} response.Response
// @Router /payments [post]
func (h *PaymentHandler) ProcessPayment(c *gin.Context) {
	var payment models.Payment
	if err := c.ShouldBindJSON(&payment); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	var transactionID string
	var status string
	var errorMessage string

	status = models.StatusPending

	switch payment.Method {
	case models.MethodCreditCard:
		transactionID, status, errorMessage = simulateCreditCardPayment(payment)
	case models.MethodBankTransfer:
		transactionID, status, errorMessage = simulateBankTransfer(payment)
	case models.MethodThirdParty:
		transactionID, status, errorMessage = simulateThirdPartyPayment(payment)
	case models.MethodBlockchain:
		transactionID, status, errorMessage = simulateBlockchainPayment(payment)
	default:
		response.Error(c, http.StatusBadRequest, "Invalid payment method")
		return
	}

	payment.TransactionID = transactionID
	payment.Status = status
	payment.ErrorMessage = errorMessage
	payment.CreatedAt = time.Now()

	if payment.Status != models.StatusSuccess && payment.Status != models.StatusFail && payment.Status != models.StatusPending {
		response.Error(c, http.StatusBadRequest, "Invalid payment status")
		return
	}

	if err := h.db.Create(&payment).Error; err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to record payment")
		return
	}

	if payment.Status == models.StatusFail {
		response.PaymentErrorResponse(c, http.StatusPaymentRequired, transactionID, status, errorMessage)
		return
	}

	response.Success(c, gin.H{
		"transaction_id": transactionID,
		"status":         status,
	})
}

// GetPayment godoc
// @Summary Get payment details
// @Description Get details of a specific payment by ID
// @Tags payments
// @Produce json
// @Param id path int true "Payment ID"
// @Success 200 {object} models.Payment
// @Failure 404 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /payments/{id} [get]
func (h *PaymentHandler) GetPayment(c *gin.Context) {
	id := c.Param("id")

	var payment models.Payment
	if err := h.db.First(&payment, id).Error; err != nil {
		response.Error(c, http.StatusNotFound, "Payment not found")
		return
	}

	response.Success(c, payment)
}

// mock diff payments
func simulateCreditCardPayment(payment models.Payment) (string, string, string) {
	return "CC123456789", models.StatusSuccess, ""
}

func simulateBankTransfer(payment models.Payment) (string, string, string) {
	return "BT987654321", models.StatusPending, ""
}

func simulateThirdPartyPayment(payment models.Payment) (string, string, string) {
	return "TP456789123", models.StatusFail, "Third-party payment failed due to insufficient funds."
}

func simulateBlockchainPayment(payment models.Payment) (string, string, string) {
	return "BC321654987", models.StatusSuccess, ""
}
