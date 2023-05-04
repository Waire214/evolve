package infrastructure

import (
	"context"
	"database/sql"
	"evolve/services/transaction/domain"
	"fmt"
	"go.uber.org/zap"
)

type transactionStoreHandler struct {
	db     *sql.DB
	logger *zap.Logger
}

func NewTransactionPersistence(db *sql.DB, logger *zap.Logger) domain.TransactionRepository {
	return &transactionStoreHandler{db: db, logger: logger}
}

func (t transactionStoreHandler) SaveTransaction(ctx context.Context, request []domain.Transactions) error {
	sqlStr := `INSERT INTO transactions (link_id, source_user_id, source_account_id, source_account_number, target_user_id, target_account_id, target_account_number, transaction_type, source_transaction_amount, target_transaction_amount, source_balance_after_transaction, target_balance_after_transaction, status, transaction_date) VALUES `
	var values []interface{}

	for i, row := range request {
		p1 := i * 14

		sqlStr += fmt.Sprintf("($%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d,$%d)", p1+1, p1+2, p1+3, p1+4, p1+5, p1+6, p1+7, p1+8, p1+9, p1+10, p1+11, p1+12, p1+13, p1+14)
		if i < len(request)-1 {
			sqlStr += ","
		}
		values = append(values, row.LinkID, row.SourceUserID, row.SourceAccountID, row.SourceAccountNumber, row.TargetUserID, row.TargetAccountID, row.TargetAccountNumber, row.TransactionType, row.SourceTransactionAmount, row.TargetTransactionAmount, row.SourceBalanceAfterTransaction, row.TargetBalanceAfterTransaction, row.Status, row.TransactionDate)
	}
	tx, err := t.db.BeginTx(ctx, nil)
	if err != nil {
		t.logger.Error("SaveTransaction", zap.String("error starting transaction", ""), zap.String("error", err.Error()), zap.String("query", sqlStr))
		return err
	}
	defer tx.Rollback()

	stmt, err := t.db.PrepareContext(ctx, sqlStr)
	if err != nil {
		t.logger.Error("SaveTransaction", zap.String("error preparing statement", ""), zap.Error(err), zap.String("query", sqlStr))
		return err
	}

	result, err := stmt.ExecContext(ctx, values...)
	if err != nil {
		t.logger.Error("SaveTransaction", zap.String("error executing query", ""), zap.Error(err), zap.String("query", sqlStr))
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if affected <= 0 {
		t.logger.Error("SaveTransaction", zap.String("error no rows created", ""), zap.Error(err), zap.String("query", sqlStr))
		return fmt.Errorf("error no rows created")
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
