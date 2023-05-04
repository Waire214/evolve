package application

import (
	"context"
	"database/sql"
	"errors"
	"evolve/services"
	accountDomain "evolve/services/account/domain"
	"evolve/services/transaction/domain"
	"fmt"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"
	"time"
)

var (
	universalPassword          string = "evolve"
	universalEncryptedPassword string
)

type transactionAppHandler struct {
	tokenHandler  services.TokenHandler
	encryptor     services.EncryptorManager
	allRepository services.Repositories
	logger        *zap.Logger
}

type TransactionApplication interface {
	CreateTransfer(ctx context.Context, request domain.Transfer) (*services.DefaultResponse, error)
}

func NewTransactionApplication(tokenHandler services.TokenHandler, allRepository services.Repositories, logger *zap.Logger) TransactionApplication {
	encryptor := services.NewEncryptor()

	// random password
	encryptPassword(encryptor)

	return &transactionAppHandler{
		tokenHandler:  tokenHandler,
		encryptor:     encryptor,
		logger:        logger,
		allRepository: allRepository,
	}
}

func (t transactionAppHandler) CreateTransfer(ctx context.Context, request domain.Transfer) (*services.DefaultResponse, error) {
	claims, err := t.tokenHandler.GetClaimsFromCtx(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to get claims")
	}
	t.logger.Info("repo")
	if t.validatePin(request.TransactionPin) != nil {
		t.logger.Info("validatePin")

		return nil, fmt.Errorf("%s", "invalid pin credential")
	}

	if request.SourceAccountID == request.TargetAccountID {
		t.logger.Info("accountCheck")

		return nil, fmt.Errorf("unable to send money to self")
	}
	ctx = context.Background()
	sourceAccount, err := t.allRepository.AccountRepository.GetAccountByAccountID(ctx, request.SourceAccountID)
	if err != nil {
		t.logger.Info("GetAccountByAccountID")

		return nil, err
	}

	t.logger.Info("sourceAccount", zap.Any("account", sourceAccount))

	if request.Amount > sourceAccount.Balance {
		t.logger.Info("balance")
		t.logger.Error("accountBalance", zap.String("balance", "insufficient balance"))
		return nil, fmt.Errorf("insufficient balance")
	}

	if err := t.performOutwardTransfer(ctx, claims.UserID, request, sourceAccount); err != nil {
		t.logger.Info("performOutward")
		return nil, err
	}

	return &services.DefaultResponse{Success: true, Message: "transfer successful"}, nil
}

func (t transactionAppHandler) performOutwardTransfer(ctx context.Context, userId uuid.UUID, request domain.Transfer, sourceAccount *accountDomain.Account) error {
	var (
		targetAccount      *accountDomain.Account
		isEvolve           bool
		targetAfterBalance float64
		sourceAfterBalance float64
	)
	targetAccount, isEvolve = t.isEvolve(ctx, request.TargetAccountNumber)
	if !isEvolve {
		//	non-Evolve account not covered
		return errors.New("payment not payable to non-evolve accounts")
	}
	if isEvolve {
		t.logger.Info("targetAccount", zap.Any("account", targetAccount))
		t.logger.Info("targetAccount", zap.String("account", "account belongs to evolve"))

		targetAfterBalance = targetAccount.Balance + request.Amount

		err := t.allRepository.AccountRepository.UpdateAccountBalances(ctx, &accountDomain.UpdateAccountBalance{ID: targetAccount.ID, Balance: targetAfterBalance})
		if err != nil {
			return err
		}
	}
	sourceAfterBalance = sourceAccount.Balance - request.Amount

	err := t.allRepository.AccountRepository.UpdateAccountBalances(ctx, &accountDomain.UpdateAccountBalance{ID: sourceAccount.ID, Balance: sourceAfterBalance})
	if err != nil {
		return err
	}

	// assumes that target account belongs to evolve
	transaction := domain.Transactions{
		SourceUserID:                  userId,
		SourceAccountID:               sourceAccount.ID,
		SourceAccountNumber:           sourceAccount.AccountNumber,
		TargetUserID:                  targetAccount.UserID,
		TargetAccountID:               targetAccount.ID,
		TargetAccountNumber:           targetAccount.AccountNumber,
		SourceTransactionAmount:       request.Amount,
		TargetTransactionAmount:       request.Amount,
		SourceBalanceAfterTransaction: sourceAfterBalance,
		TargetBalanceAfterTransaction: targetAfterBalance,
		Status:                        "paid",
		TransactionDate:               time.Now().Format(time.DateTime),
	}

	var (
		transactions []domain.Transactions
		withdrawal   = transaction
		deposit      = transaction
	)
	linkedID, err := uuid.NewV4()
	if err != nil {
		return err
	}

	withdrawal.TransactionType = "withdrawal"
	withdrawal.LinkID = linkedID
	deposit.TransactionType = "deposit"
	deposit.LinkID = linkedID
	transactions = append(transactions, withdrawal, deposit)

	err = t.allRepository.TransactionRepository.SaveTransaction(ctx, transactions)
	if err != nil {
		return err
	}
	return err
}

func (t transactionAppHandler) isEvolve(ctx context.Context, targetAccountNumber string) (*accountDomain.Account, bool) {
	targetAccount, err := t.allRepository.AccountRepository.GetAccountByAccountNumber(ctx, targetAccountNumber)
	if err == sql.ErrNoRows {
		return nil, false
	}

	if targetAccount != nil {
		return targetAccount, true
	}

	return nil, false
}
