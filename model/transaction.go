package model

import (
	"time"
)

type Transaction struct {
	ID         int64     `json:"id"`
	ProductID  int64     `json:"product_id"`
	UserID     int64     `json:"user_id"`
	Amount     float32   `json:"amount"`
	Status     string    `json:"status"`
	PaymentURL string    `json:"payment_url"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type CreateTransactionRequest struct {
	ProductID int64 `json:"product_id" validate:"required"`
	User      User
}

type UpdateTransactionRequest struct {
	TransactionStatus string `json:"transaction_status"`
	OrderID           string `json:"order_id"`
	PaymentType       string `json:"payment_type"`
	FraudStatus       string `json:"fraud_status"`
}
