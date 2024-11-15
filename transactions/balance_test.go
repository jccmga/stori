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

func TestCalculateBalance_WhenNoTransactions_ShouldReturnZero(t *testing.T) {
	t.Parallel()

	// Arrange
	var transactions []model.Transaction

	// Act
	results, err := trs.Process(transactions)

	// Assert
	require.NoError(t, err)
	assert.Equal(t, decimal.Zero, results.TotalBalance)
}

func TestCalculateBalance_WhenSingleTransaction(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name            string
		amount          decimal.Decimal
		expectedBalance decimal.Decimal
	}{
		{
			name:            "Single credit without decimals",
			amount:          decimal.MustNew(101, 0),
			expectedBalance: decimal.MustNew(101, 0),
		},
		{
			name:            "Single credit with 2 decimals",
			amount:          decimal.MustNew(10037, 2),
			expectedBalance: decimal.MustNew(10037, 2),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// Arrange
			transactions := []model.Transaction{
				{
					ID:     1,
					Date:   time.Now(),
					Amount: tc.amount,
				},
			}

			// Act
			results, err := trs.Process(transactions)

			// Assert
			require.NoError(t, err)
			assert.Equal(t, tc.expectedBalance, results.TotalBalance)
		})
	}
}

func TestCalculateBalance_WhenMultipleTransactionsAndNoProfit(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name    string
		amounts []decimal.Decimal
	}{
		{
			name: "Even credit and debit",
			amounts: []decimal.Decimal{
				decimal.MustNew(100, 0),
				decimal.MustNew(-100, 0),
			},
		},
		{
			name: "Even credit and debits with decimals",
			amounts: []decimal.Decimal{
				decimal.MustNew(100, 0),
				decimal.MustNew(-9999, 2),
				decimal.MustNew(-1, 2),
			},
		},
		{
			name: "Several credits and debits with decimals",
			amounts: []decimal.Decimal{
				decimal.MustNew(214216912, 2),
				decimal.MustNew(-78095601, 2),
				decimal.MustNew(-33941701, 2),
				decimal.MustNew(19098288, 2),
				decimal.MustNew(-31101253, 2),
				decimal.MustNew(-5406011, 2),
				decimal.MustNew(83409871, 2),
				decimal.MustNew(-85600515, 2),
				decimal.MustNew(-75723452, 2),
				decimal.MustNew(-27698238, 2),
				decimal.MustNew(39226411, 2),
				decimal.MustNew(22842284, 2),
				decimal.MustNew(-41226995, 2),
			},
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
			assert.True(t, results.TotalBalance.IsZero(), results.TotalBalance)
		})
	}
}

func TestCalculateBalance_WhenMultipleTransactionsAndNonZeroProfit(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name            string
		amounts         []decimal.Decimal
		expectedBalance decimal.Decimal
	}{
		{
			name: "Credit and debit with profit",
			amounts: []decimal.Decimal{
				decimal.MustNew(100, 0),
				decimal.MustNew(-99, 0),
			},
			expectedBalance: decimal.MustNew(1, 0),
		},
		{
			name: "Credit and debit with loss",
			amounts: []decimal.Decimal{
				decimal.MustNew(100, 0),
				decimal.MustNew(-101, 0),
			},
			expectedBalance: decimal.MustNew(-1, 0),
		},
		{
			name: "Credits and debits with profit",
			amounts: []decimal.Decimal{
				decimal.MustNew(147221686, 2),
				decimal.MustNew(133802022, 2),
				decimal.MustNew(-133802022, 2),
				decimal.MustNew(139519835, 2),
				decimal.MustNew(-191088052, 2),
				decimal.MustNew(-114266700, 2),
				decimal.MustNew(141848380, 2),
				decimal.MustNew(109259987, 2),
			},
			expectedBalance: decimal.MustNew(232495136, 2),
		},
		{
			name: "Credits and debits with no profit",
			amounts: []decimal.Decimal{
				decimal.MustNew(127419250, 2),
				decimal.MustNew(-140520427, 2),
				decimal.MustNew(153383766, 2),
				decimal.MustNew(-191931547, 2),
				decimal.MustNew(-199292563, 2),
			},
			expectedBalance: decimal.MustNew(-250941521, 2),
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
			assert.Equal(t, tc.expectedBalance, results.TotalBalance)
		})
	}
}
