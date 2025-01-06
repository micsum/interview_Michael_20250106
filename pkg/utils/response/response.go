package response

import "github.com/gin-gonic/gin"

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type PaymentError struct {
	Code          int    `json:"code"`
	TransactionID string `json:"transaction_id"`
	Status        string `json:"status"`
	ErrorMessage  string `json:"error_message"`
}

type JoinResponse struct {
	Success 	bool   `json:"success_join"`
	ChallengeID uint    `json:"challenge_id"`
	Message		string `json:"message"`
}

type LevelCreateResponse struct{
	LevelID 	uint 	`json:"level_id"`
}

type LogCreateResponse struct{
	LogID 	uint 	`json:"log_id"`
}

type PlayerCreateResponse struct{
	PlayerID uint	`json:"player_id"`
}

type ReservCreateResponse struct{
	ReservID uint	`json:"reservation_id"`
}

type RoomCreateResponse struct{
	RoomID uint	`json:"room_id"`
}

func Success(c *gin.Context, data interface{}) {
	c.JSON(200, Response{
		Code:    200,
		Message: "success",
		Data:    data,
	})
}

func Error(c *gin.Context, code int, message string) {
	c.JSON(code, Response{
		Code:    code,
		Message: message,
		Data:    nil,
	})
}

func PaymentErrorResponse(c *gin.Context, code int, transactionID, status, errorMessage string) {
	c.JSON(code, PaymentError{
		Code:          code,
		TransactionID: transactionID,
		Status:        status,
		ErrorMessage:  errorMessage,
	})
}
