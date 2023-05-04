package domain

import (
	"github.com/gofrs/uuid"
)

type Transactions struct {
	ID                            uuid.UUID `json:"id" sql:"id"`
	LinkID                        uuid.UUID `json:"link_id" sl:"link_id"`
	SourceUserID                  uuid.UUID `json:"source_user_id" sql:"source_user_id"`
	SourceAccountID               uuid.UUID `json:"source_account_id" sql:"source_account_id"`
	SourceAccountNumber           string    `json:"source_account_number" sql:"source_account_number"`
	TargetUserID                  uuid.UUID `json:"target_user_id" sql:"target_user_id"`
	TargetAccountID               uuid.UUID `json:"target_account_id" sql:"target_account_id"`
	TargetAccountNumber           string    `json:"target_account_number" sql:"target_account_number"`
	TransactionType               string    `json:"transaction_type" sql:"transaction_type"`
	SourceTransactionAmount       float64   `json:"source_transaction_amount" sql:"source_transaction_amount"`
	TargetTransactionAmount       float64   `json:"target_transaction_amount" sql:"target_transaction_amount"`
	SourceBalanceAfterTransaction float64   `json:"source_balance_after_transaction" sql:"source_balance_after_transaction"`
	TargetBalanceAfterTransaction float64   `json:"target_balance_after_transaction" sql:"target_balance_after_transaction"`
	Status                        string    `json:"status" sql:"status"`
	TransactionDate               string    `json:"transaction_date" sql:"transaction_date"`
	UpdatedAt                     string    `json:"updated_at" sql:"updated_at"`
}

type Transfer struct {
	SourceUserID        uuid.UUID `json:"source_user_id" sql:"source_user_id"`
	SourceAccountID     uuid.UUID `json:"source_account_id" sql:"source_account_id"`
	SourceAccountNumber string    `json:"source_account_number" sql:"source_account_number"`
	TargetUserID        uuid.UUID `json:"target_user_id" sql:"target_user_id"`
	TargetAccountID     uuid.UUID `json:"target_account_id" sql:"target_account_id"`
	TargetAccountNumber string    `json:"target_account_number" sql:"target_account_number"`
	Amount              float64   `json:"amount" sql:"amount"`
	TransactionPin      string    `json:"transaction_pin" sql:"transaction_pin"`
}
