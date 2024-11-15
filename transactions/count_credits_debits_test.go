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

func TestCountCreditsAndDebits_WhenNoTransactions(t *testing.T) {
	t.Parallel()

	// Arrange
	var transactions []model.Transaction

	// Act
	results, err := trs.Process(transactions)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, 0, results.CreditTransactionsCount)
	assert.Equal(t, 0, results.DebitTransactionsCount)
	assert.True(t, results.TotalCreditAmount.IsZero())
	assert.True(t, results.TotalDebitAmount.IsZero())
}

func TestCountCreditsAndDebits_WhenSingleTransactionFromEach(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name                 string
		amounts              []decimal.Decimal
		expectedCreditAmount decimal.Decimal
		expectedDebitAmount  decimal.Decimal
		expectedCreditCount  int
		expectedDebitCount   int
	}{
		{
			name: "Even credit and debit without decimals",
			amounts: []decimal.Decimal{
				decimal.MustNew(100, 0),
				decimal.MustNew(-100, 0),
			},
			expectedCreditAmount: decimal.MustNew(100, 0),
			expectedDebitAmount:  decimal.MustNew(-100, 0),
			expectedCreditCount:  1,
			expectedDebitCount:   1,
		},
		{
			name: "Credit and debit with 2 decimals",
			amounts: []decimal.Decimal{
				decimal.MustNew(1124065, 2),
				decimal.MustNew(-998323, 2),
			},
			expectedCreditAmount: decimal.MustNew(1124065, 2),
			expectedDebitAmount:  decimal.MustNew(-998323, 2),
			expectedCreditCount:  1,
			expectedDebitCount:   1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// Arrange
			var transactions []model.Transaction
			for i, amount := range tc.amounts {
				transactions = append(transactions, model.Transaction{
					ID:     i + 1,
					Date:   time.Now(),
					Amount: amount,
				})
			}

			// Act
			results, err := trs.Process(transactions)

			// Assert
			require.NoError(t, err)
			assert.Equal(t, tc.expectedCreditCount, results.CreditTransactionsCount)
			assert.Equal(t, tc.expectedDebitCount, results.DebitTransactionsCount)
			assert.Equal(t, tc.expectedCreditAmount, results.TotalCreditAmount)
			assert.Equal(t, tc.expectedDebitAmount, results.TotalDebitAmount)
		})
	}
}

func TestCountCreditsAndDebits_SeveralTransactions(t *testing.T) {
	t.Parallel()

	// Arrange
	amounts := []decimal.Decimal{
		decimal.MustNew(147221686, 2),
		decimal.MustNew(133802022, 2),
		decimal.MustNew(-133802022, 2),
		decimal.MustNew(139519835, 2),
		decimal.MustNew(-191088052, 2),
		decimal.MustNew(-114266700, 2),
		decimal.MustNew(141848380, 2),
		decimal.MustNew(109259987, 2),
	}
	expectedCreditAmount := decimal.MustNew(671651910, 2)
	expectedDebitAmount := decimal.MustNew(-439156774, 2)

	var transactions []model.Transaction
	for i, amount := range amounts {
		transactions = append(transactions, model.Transaction{
			ID:     i + 1,
			Date:   time.Now(),
			Amount: amount,
		})
	}

	// Act
	results, err := trs.Process(transactions)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, 5, results.CreditTransactionsCount)
	assert.Equal(t, 3, results.DebitTransactionsCount)
	assert.Equal(t, expectedCreditAmount, results.TotalCreditAmount)
	assert.Equal(t, expectedDebitAmount, results.TotalDebitAmount)
}
