package transactions_test

import (
	"testing"
	"time"

	"github.com/govalues/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"stori/model"
	trs "stori/transactions"
)

func TestCountTransactionsByMonth_WhenNoTransactions(t *testing.T) {
	t.Parallel()

	// Arrange
	var transactions []model.Transaction

	// Act
	results, err := trs.Process(transactions)

	// Assert
	require.NoError(t, err)
	for i := time.January; i <= time.December; i++ {
		assert.Equal(t, 0, results.TransactionsPerMonth[i])
	}
}

func TestCountTransactionsByMonth_SeveralTransactions(t *testing.T) {
	t.Parallel()

	// Arrange
	transactions := []model.Transaction{
		{
			ID:     1,
			Date:   time.Date(2024, time.January, 1, 0, 0, 0, 0, time.UTC),
			Amount: decimal.MustNew(100, 0),
		},
		{
			ID:     2,
			Date:   time.Date(2024, time.January, 1, 12, 0, 0, 0, time.UTC),
			Amount: decimal.MustNew(100, 0),
		},
		{
			ID:     3,
			Date:   time.Date(2024, time.February, 11, 17, 3, 0, 0, time.UTC),
			Amount: decimal.MustNew(100, 0),
		},
		{
			ID:     4,
			Date:   time.Date(2024, time.November, 11, 17, 12, 40, 0, time.UTC),
			Amount: decimal.MustNew(100, 0),
		},
		{
			ID:     5,
			Date:   time.Date(2024, time.November, 30, 20, 31, 0, 0, time.UTC),
			Amount: decimal.MustNew(100, 0),
		},
	}

	// Act
	results, err := trs.Process(transactions)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, 2, results.TransactionsPerMonth[time.January])
	assert.Equal(t, 1, results.TransactionsPerMonth[time.February])
	assert.Equal(t, 0, results.TransactionsPerMonth[time.March])
	assert.Equal(t, 0, results.TransactionsPerMonth[time.April])
	assert.Equal(t, 0, results.TransactionsPerMonth[time.May])
	assert.Equal(t, 0, results.TransactionsPerMonth[time.June])
	assert.Equal(t, 0, results.TransactionsPerMonth[time.July])
	assert.Equal(t, 0, results.TransactionsPerMonth[time.August])
	assert.Equal(t, 0, results.TransactionsPerMonth[time.September])
	assert.Equal(t, 0, results.TransactionsPerMonth[time.October])
	assert.Equal(t, 2, results.TransactionsPerMonth[time.November])
	assert.Equal(t, 0, results.TransactionsPerMonth[time.December])
}
