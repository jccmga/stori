package filereader_test

import (
	"testing"
	"time"

	"github.com/govalues/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"stori/adapters/filereader"
	"stori/test"
)

func TestReadTransactions_WhenFileDoesNotExist_Error(t *testing.T) {
	test.IntegrationTest(t)
	t.Parallel()

	// Arrange
	sut := filereader.NewLocalReader("testdata/non-existent-file.csv")

	// Act
	_, err := sut.ReadTransactions()

	// Assert
	require.Error(t, err)
	assert.ErrorIs(t, err, filereader.ErrFileNotFound)
}

func TestReadTransactions_WhenSingleTransaction_Success(t *testing.T) {
	test.IntegrationTest(t)
	t.Parallel()

	// Arrange
	sut := filereader.NewLocalReader("testdata/single_transaction.csv")

	// Act
	transactions, err := sut.ReadTransactions()

	// Assert
	require.NoError(t, err)
	assert.Len(t, transactions, 1)
	assert.Equal(t, 0, transactions[0].ID)
	assert.Equal(t, time.July, transactions[0].Date.Month())
	assert.Equal(t, decimal.MustParse("60.5"), transactions[0].Amount)
}

func TestReadTransactions_WhenSeveralTransactions_Success(t *testing.T) {
	test.IntegrationTest(t)
	t.Parallel()

	// Arrange
	sut := filereader.NewLocalReader("testdata/several_transactions.csv")

	// Act
	transactions, err := sut.ReadTransactions()

	// Assert
	require.NoError(t, err)
	assert.Len(t, transactions, 4)

	assert.Equal(t, 0, transactions[0].ID)
	assert.Equal(t, time.July, transactions[0].Date.Month())
	assert.Equal(t, decimal.MustParse("60.5"), transactions[0].Amount)

	assert.Equal(t, 1, transactions[1].ID)
	assert.Equal(t, time.July, transactions[1].Date.Month())
	assert.Equal(t, decimal.MustParse("-10.3"), transactions[1].Amount)

	assert.Equal(t, 2, transactions[2].ID)
	assert.Equal(t, time.August, transactions[2].Date.Month())
	assert.Equal(t, decimal.MustParse("-20.46"), transactions[2].Amount)

	assert.Equal(t, 3, transactions[3].ID)
	assert.Equal(t, time.August, transactions[3].Date.Month())
	assert.Equal(t, decimal.MustParse("+10"), transactions[3].Amount)
}

func TestReadTransactions_WhenFileIsInvalid_Error(t *testing.T) {
	test.IntegrationTest(t)
	t.Parallel()

	testCases := []struct {
		name          string
		filename      string
		expectedError error
	}{
		{
			name:          "When file is empty",
			filename:      "testdata/empty_file.csv",
			expectedError: filereader.ErrFileIsEmpty,
		},
		{
			name:          "When file has invalid header",
			filename:      "testdata/invalid_header.csv",
			expectedError: filereader.ErrInvalidHeader,
		},
		{
			name:          "When file has invalid amount",
			filename:      "testdata/invalid_amount.csv",
			expectedError: filereader.ErrInvalidAmount,
		},
		{
			name:          "When file has invalid date format",
			filename:      "testdata/invalid_date.csv",
			expectedError: filereader.ErrInvalidDateFormat,
		},
		{
			name:          "When file has invalid columns",
			filename:      "testdata/invalid_columns.csv",
			expectedError: filereader.ErrInvalidFile,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			// Arrange
			sut := filereader.NewLocalReader(tc.filename)

			// Act
			_, err := sut.ReadTransactions()

			// Assert
			require.Error(t, err)
			assert.ErrorIs(t, err, tc.expectedError)
		})
	}
}
