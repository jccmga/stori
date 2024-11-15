package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"

	"stori/accountsummary"
)

func (repo Repository) Create(execution accountsummary.Execution) error {
	fail := func(err error) error {
		return fmt.Errorf("repository: Repository: Create: %w", err)
	}

	tx, errX := repo.DB.Beginx()
	if errX != nil {
		return fail(fmt.Errorf("failed creating a db transaction: %w", errX))
	}

	defer func(tx *sqlx.Tx) {
		if err := tx.Rollback(); err != nil && !errors.Is(err, sql.ErrTxDone) {
			log.Fatal(fail(err))
		}
	}(tx)

	accountDB := repo.AccountSummaryFromModel(execution.AccountSummary, execution.FilePath)
	if err := repo.createAccountSummary(accountDB); err != nil {
		return fail(fmt.Errorf("failed creating an account summary: %w", err))
	}

	transactionsDB := make([]Transaction, len(execution.Transactions))
	for i, transaction := range execution.Transactions {
		transactionsDB[i] = repo.TransactionFromModel(transaction, execution.FilePath)
	}

	if err := repo.createTransactions(transactionsDB); err != nil {
		return fail(fmt.Errorf("failed creating transactions: %w", err))
	}

	if err := tx.Commit(); err != nil {
		return fail(fmt.Errorf("failed committing the db transaction: %w", err))
	}

	return nil
}

func (repo Repository) createAccountSummary(summary AccountSummary) error {
	fail := func(err error) error {
		return fmt.Errorf("repository: Repository: createAccountSummary: %w", err)
	}

	query := `insert into account_summary(
email, total_balance, average_debit_amount, average_credit_amount, transactions_per_month, file_path)
values (:email, :total_balance, :average_debit_amount, :average_credit_amount, :transactions_per_month, :file_path)`

	_, err := repo.DB.NamedExec(query, summary)
	if err != nil {
		return fail(err)
	}

	return nil
}

func (repo Repository) createTransactions(transactions []Transaction) error {
	fail := func(err error) error {
		return fmt.Errorf("repository: Repository: createTransactions: %w", err)
	}

	query := `insert into transaction(date, amount, file_path)
values (:date, :amount, :file_path)`

	_, err := repo.DB.NamedExec(query, transactions)
	if err != nil {
		return fail(err)
	}

	return nil
}
