package repository_test

import (
	"testing"
	"time"

	"github.com/govalues/decimal"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"stori/accountsummary"
	"stori/adapters/repository"
	"stori/model"
	"stori/test"
)

func TestCreate(t *testing.T) {
	test.IntegrationTest(t)
	t.Parallel()

	// Arrange
	dbx := sqlx.NewDb(DB, "postgres")
	sut := repository.New(dbx)
	summary := model.AccountSummary{
		Email:               "john.doe@example.com",
		TotalBalance:        decimal.MustParse("100.0"),
		AverageDebitAmount:  decimal.MustParse("50.0"),
		AverageCreditAmount: decimal.MustParse("150.0"),
		TransactionsPerMonth: map[time.Month]int{ //nolint:exhaustive
			time.January:  1,
			time.February: 1,
		},
	}
	execution := accountsummary.Execution{
		AccountSummary: summary,
		Transactions: []model.Transaction{
			{
				ID:     1,
				Date:   time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC),
				Amount: decimal.MustNew(150, 0),
			},
			{
				ID:     2,
				Date:   time.Date(2024, time.February, 1, 0, 0, 0, 0, time.UTC),
				Amount: decimal.MustNew(50, 0),
			},
		},
		FilePath: "test.csv",
	}

	// Act
	err := sut.Create(execution)

	// Assert
	require.NoError(t, err)

	respSummaries := []repository.AccountSummary{}
	err = dbx.Select(&respSummaries,
		`select email, total_balance,
				average_debit_amount, average_credit_amount, transactions_per_month, file_path
				from account_summary`)

	require.NoError(t, err)
	require.Len(t, respSummaries, 1)

	res := respSummaries[0]
	assert.Equal(t, "john.doe@example.com", res.Email)
	assert.Equal(t, decimal.MustParse("100.0"), res.TotalBalance)
	assert.Equal(t, decimal.MustParse("150.0"), res.AverageCreditAmount)
	assert.Equal(t, decimal.MustParse("50.0"), res.AverageDebitAmount)
	assert.Equal(t, "test.csv", res.FilePath)

	require.NotEmpty(t, res.TransactionsPerMonth)
	assert.Equal(t, 1, res.TransactionsPerMonth[time.January])
	assert.Equal(t, 1, res.TransactionsPerMonth[time.February])
}
