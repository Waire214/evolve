package services

import (
	accountDomain "evolve/services/account/domain"
	transactionDomain "evolve/services/transaction/domain"
)

type Repositories struct {
	TransactionRepository transactionDomain.TransactionRepository
	AccountRepository     accountDomain.AccountRepository
}

type DefaultResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}
