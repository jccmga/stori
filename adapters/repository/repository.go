package repository

import (
	"github.com/jmoiron/sqlx"

	"stori/model"
)

type (
	Repository struct {
		Config
	}

	Config struct {
		DB *sqlx.DB
	}
)

func New(db *sqlx.DB) *Repository {
	return &Repository{
		Config: Config{
			DB: db,
		},
	}
}

func (repo Repository) AccountSummaryFromModel(summary model.AccountSummary, filePath string) AccountSummary {
	return AccountSummary{
		Email:                summary.Email,
		TotalBalance:         summary.TotalBalance,
		AverageDebitAmount:   summary.AverageDebitAmount,
		AverageCreditAmount:  summary.AverageCreditAmount,
		TransactionsPerMonth: summary.TransactionsPerMonth,
		FilePath:             filePath,
	}
}

func (repo Repository) TransactionFromModel(transaction model.Transaction, filePath string) Transaction {
	return Transaction{
		Date:     transaction.Date,
		Amount:   transaction.Amount,
		FilePath: filePath,
	}
}
