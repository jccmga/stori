package accountsummary

import "stori/model"

type TransactionsReader interface {
	ReadTransactions() ([]model.Transaction, error)
}

type EmailSender interface {
	Send(summary model.AccountSummary) error
}

type Repository interface {
	Create(execution Execution) error
}
