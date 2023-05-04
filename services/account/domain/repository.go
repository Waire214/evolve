package domain

import (
	"context"
	"github.com/gofrs/uuid"
)

type AccountRepository interface {
	GetAccountByAccountID(ctx context.Context, request uuid.UUID) (*Account, error)
	GetAccountByAccountNumber(ctx context.Context, request string) (*Account, error)
	UpdateAccountBalances(ctx context.Context, request *UpdateAccountBalance) error
	CreateCashCache(ctx context.Context, request *CashCache) (*CashCache, error)
}
