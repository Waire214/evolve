package application

import (
	"context"
	"evolve/services"
	"evolve/services/account/domain"
	domain2 "evolve/services/transaction/domain"
	"fmt"
	"github.com/gofrs/uuid"
	"time"
)

type accountAppHandler struct {
	tokenHandler  services.TokenHandler
	encryptor     services.EncryptorManager
	allRepository services.Repositories
}

type AccountApplication interface {
	GetAccountByAccountID(ctx context.Context, request uuid.UUID) (*domain.Account, error)
	GetAccountByAccountNumber(ctx context.Context, request string) (*domain.Account, error)
	UpdateAccountBalances(ctx context.Context, request *domain.UpdateAccountBalance) error
	CreateCashCache(ctx context.Context, request *domain.CashCache) (*domain.CashCache, error)
}

func NewAccountApplication(tokenHandler services.TokenHandler, allRepository services.Repositories) AccountApplication {
	encryptor := services.NewEncryptor()

	return &accountAppHandler{
		tokenHandler:  tokenHandler,
		encryptor:     encryptor,
		allRepository: allRepository,
	}
}

func (a accountAppHandler) GetAccountByAccountID(ctx context.Context, request uuid.UUID) (*domain.Account, error) {
	return nil, nil
}

func (a accountAppHandler) GetAccountByAccountNumber(ctx context.Context, request string) (*domain.Account, error) {
	return nil, nil
}

func (a accountAppHandler) UpdateAccountBalances(ctx context.Context, request *domain.UpdateAccountBalance) error {
	return nil
}

func (a accountAppHandler) CreateCashCache(ctx context.Context, request *domain.CashCache) (*domain.CashCache, error) {
	claims, err := a.tokenHandler.GetClaimsFromCtx(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to get claims")
	}
	// get user's account details and check if active
	accountDetails, err := a.allRepository.AccountRepository.GetAccountByAccountNumber(ctx, request.AccountNumber)
	if err != nil {
		return nil, err
	}

	if accountDetails.Status != "active" {
		return nil, fmt.Errorf("unable to create cache because your savings account is inactive")
	}
	if accountDetails.Balance < request.Amount {
		return nil, fmt.Errorf("you do not have sufficient balance to open a cache")
	}

	// initiate cache creation
	now := time.Now()
	oneDayCache := now.Add(24 * time.Hour).Sub(now)
	cache := domain.CashCache{
		UserID:        claims.UserID,
		AccountNumber: request.AccountNumber, // assumes the generation of random account numbers
		Amount:        request.Amount,
		Duration:      oneDayCache,
		CreatedAt:     time.Now().Format(time.DateTime),
	}
	openedCache, err := a.allRepository.AccountRepository.CreateCashCache(ctx, &cache)
	if err != nil {
		return nil, err
	}

	var balance = accountDetails.Balance - request.Amount
	transaction := domain2.Transactions{
		SourceUserID:                  accountDetails.UserID,
		SourceAccountID:               accountDetails.ID,
		SourceAccountNumber:           accountDetails.AccountNumber,
		TransactionType:               "cache funding",
		SourceTransactionAmount:       request.Amount,
		SourceBalanceAfterTransaction: balance,
		TransactionDate:               time.Now().Format(time.DateTime),
	}
	var transactions []domain2.Transactions

	transactions = append(transactions, transaction)

	err = a.allRepository.TransactionRepository.SaveTransaction(ctx, transactions)
	if err != nil {
		return nil, err
	}

	err = a.allRepository.AccountRepository.UpdateAccountBalances(ctx, &domain.UpdateAccountBalance{ID: accountDetails.ID, Balance: balance})
	if err != nil {
		return nil, err
	}
	return openedCache, nil
}
