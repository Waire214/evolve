package infrastructure

import (
	"context"
	"database/sql"
	domain "evolve/services/account/domain"
	"fmt"
	"github.com/gofrs/uuid"
	"go.uber.org/zap"
)

type accountStoreHandler struct {
	db     *sql.DB
	logger *zap.Logger
}

func NewAccountPersistence(db *sql.DB, logger *zap.Logger) domain.AccountRepository {
	return &accountStoreHandler{db: db, logger: logger}
}

func (a accountStoreHandler) GetAccountByAccountID(ctx context.Context, request uuid.UUID) (*domain.Account, error) {
	sqlQuery := `SELECT * FROM accounts WHERE id=$1`
	stmt, err := a.db.PrepareContext(ctx, sqlQuery)
	if err != nil {
		a.logger.Error("GetAccountByAccountID", zap.String("error preparing statement", ""), zap.Error(err), zap.String("query", sqlQuery))
		return nil, err
	}
	var account domain.Account
	err = stmt.QueryRowContext(ctx, request).Scan(
		&account.ID,
		&account.UserID,
		&account.AccountType,
		&account.BankCode,
		&account.AccountNumber,
		&account.Balance,
		&account.Status,
		&account.CreatedAt,
		&account.UpdatedAt,
	)
	if err != nil {
		a.logger.Error("GetAccountByAccountID", zap.String("error scanning row", ""), zap.Error(err), zap.String("query", sqlQuery))
		return nil, fmt.Errorf("not found")

	}

	return &account, nil
}

func (a accountStoreHandler) GetAccountByAccountNumber(ctx context.Context, request string) (*domain.Account, error) {
	sqlQuery := `SELECT * FROM accounts WHERE account_number=$1`
	stmt, err := a.db.PrepareContext(ctx, sqlQuery)
	if err != nil {
		a.logger.Error("GetAccountByAccountNumber", zap.String("error preparing statement", ""), zap.Error(err), zap.String("query", sqlQuery))
		return nil, err
	}
	var account domain.Account
	err = stmt.QueryRowContext(ctx, request).Scan(
		&account.ID,
		&account.UserID,
		&account.AccountType,
		&account.BankCode,
		&account.AccountNumber,
		&account.Balance,
		&account.Status,
		&account.CreatedAt,
		&account.UpdatedAt,
	)
	if err != nil {
		a.logger.Error("GetAccountByAccountNumber", zap.String("error scanning row", ""), zap.Error(err), zap.String("query", sqlQuery))
		return nil, fmt.Errorf("not found")
	}

	return &account, nil
}

func (a accountStoreHandler) UpdateAccountBalances(ctx context.Context, request *domain.UpdateAccountBalance) error {
	var medID string
	sqlQuery := `UPDATE accounts SET balance=$2 WHERE id=$1 RETURNING id`
	stmt, err := a.db.PrepareContext(ctx, sqlQuery)
	if err != nil {
		a.logger.Error("UpdateAccountBalances", zap.String("error preparing statement", err.Error()), zap.String("sqlQuery : ", sqlQuery))

		return err
	}
	row := stmt.QueryRowContext(ctx, request.ID, request.Balance)
	if err := row.Scan(&medID); err != nil {
		a.logger.Error("UpdateAccountBalances", zap.String("error scanning row", err.Error()))
		return err
	}
	return nil

}

func (a accountStoreHandler) CreateCashCache(ctx context.Context, request *domain.CashCache) (*domain.CashCache, error) {
	const SQL = `INSERT INTO cashCache (user_id, account_number, cached_amount, cache_duration, cached_at) VALUES ($1, $2, $3, $4, $5) RETURNING id`

	tx, err := a.db.BeginTx(ctx, nil)
	if err != nil {
		a.logger.Error("CreateCashCache", zap.String("error starting transaction", ""), zap.String("error", err.Error()), zap.String("query", SQL))
		return nil, err
	}
	defer tx.Rollback()

	tmpSmt, err := tx.PrepareContext(ctx, SQL)
	if err != nil {
		a.logger.Error("CreateCashCache", zap.String("error preparing statement", ""), zap.Error(err), zap.String("query", SQL))
		return nil, err
	}

	var createdCacheID uuid.UUID

	err = tmpSmt.QueryRowContext(ctx,
		request.UserID,
		request.AccountNumber,
		request.Amount,
		request.Duration,
		request.CreatedAt,
	).Scan(&createdCacheID)
	if err != nil {
		a.logger.Error("error", zap.String("error", err.Error()), zap.String("query", SQL))
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	request.ID = createdCacheID
	return request, nil
}
