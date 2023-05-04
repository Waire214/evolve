package domain

import (
	"github.com/gofrs/uuid"
	"time"
)

type Account struct {
	ID            uuid.UUID `json:"id" sql:"id"`
	UserID        uuid.UUID `json:"user_id" sql:"user_id"`
	AccountType   string    `json:"account_type" sql:"account_type"`
	BankCode      string    `json:"bank_code" sql:"bank_code"`
	AccountNumber string    `json:"account_number" sql:"account_number"`
	Balance       float64   `json:"balance" sql:"balance"`
	Status        string    `json:"status" sql:"status"`
	CreatedAt     string    `json:"created_at" sql:"created_at"`
	UpdatedAt     string    `json:"updated_at" sql:"updated_at"`
}

type UpdateAccountBalance struct {
	ID      uuid.UUID `json:"id" sql:"id"`
	Balance float64   `json:"balance" sql:"balance"`
}

type CashCache struct {
	ID            uuid.UUID     `json:"id" sql:"id"`
	UserID        uuid.UUID     `json:"user_id" sql:"user_id"`
	AccountNumber string        `json:"account_number" sql:"account_number"`
	Amount        float64       `json:"amount" sql:"cached_amount"`
	Duration      time.Duration `json:"duration" sql:"cache_duration"`
	CreatedAt     string        `json:"created_at" sql:"cached_at"`
}
