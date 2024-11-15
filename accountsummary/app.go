package accountsummary

import (
	"errors"
	"fmt"
	"net/mail"

	"github.com/govalues/decimal"

	"stori/model"
	trans "stori/transactions"
)

var (
	ErrNoTransactions = errors.New("no transactions")
	ErrInvalidEmail   = errors.New("invalid email")
)

type (
	Config struct {
		Email              string
		FilePath           string
		TransactionsReader TransactionsReader
		EmailSender        EmailSender
		Repository         Repository
	}

	App struct {
		Config
	}

	AverageAmounts struct {
		credit decimal.Decimal
		debit  decimal.Decimal
	}
)

func New(config Config) App {
	return App{Config: config}
}

func (app App) Run() error {
	fail := func(err error) error {
		return fmt.Errorf("app: App: Run: %w", err)
	}

	if isInvalidEmail(app.Email) {
		return fail(ErrInvalidEmail)
	}

	transactions, err := app.TransactionsReader.ReadTransactions()
	if err != nil {
		return fail(err)
	}

	if validTransactions(transactions) {
		if err = app.processTransactions(transactions); err != nil {
			return fail(err)
		}
	} else {
		return fail(ErrNoTransactions)
	}

	return nil
}

func isInvalidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err != nil
}

func validTransactions(transactions []model.Transaction) bool {
	return len(transactions) > 0
}

func (app App) processTransactions(transactions []model.Transaction) error {
	fail := func(err error) error {
		return fmt.Errorf("app: App: processTransactions: %w", err)
	}

	results, err := trans.Process(transactions)
	if err != nil {
		return fail(err)
	}

	averages, err := calculateAverageAmounts(results)
	if err != nil {
		return fail(err)
	}

	summary := model.AccountSummary{
		Email:                app.Email,
		TotalBalance:         results.TotalBalance,
		AverageDebitAmount:   averages.debit,
		AverageCreditAmount:  averages.credit,
		TransactionsPerMonth: results.TransactionsPerMonth,
	}

	if err = app.Repository.Create(Execution{
		FilePath:       app.FilePath,
		AccountSummary: summary,
		Transactions:   transactions,
	}); err != nil {
		return fail(err)
	}

	if err = app.sendEmail(summary); err != nil {
		return fail(err)
	}

	return nil
}

func calculateAverageAmounts(results trans.ExecutionResults) (AverageAmounts, error) {
	fail := func(err error) (AverageAmounts, error) {
		return AverageAmounts{}, fmt.Errorf("app: App: calculateAverageAmounts: %w", err)
	}

	avgDebit := decimal.Zero
	avgCredit := decimal.Zero
	var err error

	if results.DebitTransactionsCount > 0 {
		avgDebit, err = results.TotalDebitAmount.Quo(
			decimal.MustNew(int64(results.DebitTransactionsCount), 0))
		if err != nil {
			return fail(err)
		}
	}

	if results.CreditTransactionsCount > 0 {
		avgCredit, err = results.TotalCreditAmount.Quo(
			decimal.MustNew(int64(results.CreditTransactionsCount), 0))
		if err != nil {
			return fail(err)
		}
	}

	return AverageAmounts{
		credit: avgCredit,
		debit:  avgDebit,
	}, nil
}

func (app App) sendEmail(summary model.AccountSummary) error {
	if err := app.EmailSender.Send(summary); err != nil {
		return fmt.Errorf("app: App: sendEmail: %w", err)
	}

	return nil
}
