package domain

import "context"

type TransactionRepository interface {
	SaveTransaction(ctx context.Context, request []Transactions) error
}
