package models

import (
	"time"
)

const (
	StatusSuccess = "success"
	StatusFail    = "fail"
	StatusPending = "pending"

	MethodCreditCard   = "credit_card"
	MethodBankTransfer = "bank_transfer"
	MethodThirdParty   = "third_party"
	MethodBlockchain   = "blockchain"
)

type Payment struct {
	ID            uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Method        string    `json:"method" gorm:"not null"`
	Amount        float64   `json:"amount" gorm:"not null"`
	Details       string    `json:"details" gorm:"type:text"`
	Status        string    `json:"status" gorm:"not null"`
	TransactionID string    `json:"transaction_id" gorm:"not null"`
	ErrorMessage  string    `json:"error_message" gorm:"type:text"`
	CreatedAt     time.Time `json:"created_at"`
}
