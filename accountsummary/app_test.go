package accountsummary_test

import (
	"testing"
	"time"

	"github.com/govalues/decimal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	accsum "stori/accountsummary"
	mocks "stori/mocks/stori/accountsummary"
	"stori/model"
)

func TestAppRun_NoTransactions(t *testing.T) {
	t.Parallel()

	// Arrange
	readerStub := mocks.NewMockTransactionsReader(t)
	readerStub.EXPECT().ReadTransactions().Return([]model.Transaction{}, nil)

	sut := accsum.New(accsum.Config{
		Email:              "john.doe@stori.com",
		TransactionsReader: readerStub,
	})

	// Act
	err := sut.Run()

	// Assert
	require.Error(t, err)
	assert.ErrorIs(t, err, accsum.ErrNoTransactions)
}

func TestAppRun_WhenSingleTransaction_SuccessEmailSenderCalled(t *testing.T) {
	t.Parallel()

	// Arrange && Assert
	email := "john.doe@stori.com"
	transactions := []model.Transaction{buildTransaction(1, time.November, "100")}
	expectedSummary := model.AccountSummary{
		Email:                email,
		TotalBalance:         decimal.MustParse("100"),
		AverageDebitAmount:   decimal.Zero,
		AverageCreditAmount:  decimal.MustParse("100"),
		TransactionsPerMonth: buildTransactionPerMonth([]time.Month{time.November}, []int{1}),
	}

	readerStub := mocks.NewMockTransactionsReader(t)
	emailSenderMock := mocks.NewMockEmailSender(t)
	repositoryMock := mocks.NewMockRepository(t)

	readerStub.EXPECT().ReadTransactions().Return(transactions, nil)
	emailSenderMock.EXPECT().Send(expectedSummary).Return(nil)
	repositoryMock.EXPECT().Create(mock.Anything).Return(nil)

	sut := accsum.New(accsum.Config{
		Email:              email,
		TransactionsReader: readerStub,
		EmailSender:        emailSenderMock,
		Repository:         repositoryMock,
	})

	// Act
	err := sut.Run()

	// Assert
	require.NoError(t, err)
}

func TestAppRun_WhenSeveralTransactions_SuccessEmailSenderCalled(t *testing.T) {
	t.Parallel()

	// Arrange && Assert
	email := "john.doe@stori.com"
	transactions := []model.Transaction{
		buildTransaction(1, time.February, "1472216.86"),
		buildTransaction(2, time.February, "1338020.22"),
		buildTransaction(3, time.April, "-1338020.22"),
		buildTransaction(4, time.November, "1395198.35"),
		buildTransaction(5, time.November, "-1910880.52"),
		buildTransaction(6, time.November, "-1142667.00"),
		buildTransaction(7, time.November, "1418483.80"),
		buildTransaction(8, time.December, "1092599.87"),
	}
	expectedSummary := model.AccountSummary{
		Email:               email,
		TotalBalance:        decimal.MustNew(232495136, 2),
		AverageDebitAmount:  decimal.MustNew(-1463855913333333333, 12),
		AverageCreditAmount: decimal.MustNew(134330382, 2),
		TransactionsPerMonth: buildTransactionPerMonth(
			[]time.Month{time.February, time.April, time.November, time.December}, []int{2, 1, 4, 1}),
	}

	readerStub := mocks.NewMockTransactionsReader(t)
	emailSenderMock := mocks.NewMockEmailSender(t)
	repositoryMock := mocks.NewMockRepository(t)

	readerStub.EXPECT().ReadTransactions().Return(transactions, nil)
	emailSenderMock.EXPECT().Send(expectedSummary).Return(nil).Once()
	repositoryMock.EXPECT().Create(mock.Anything).Return(nil)

	sut := accsum.New(accsum.Config{
		Email:              email,
		TransactionsReader: readerStub,
		EmailSender:        emailSenderMock,
		Repository:         repositoryMock,
	})

	// Act
	err := sut.Run()

	// Assert
	require.NoError(t, err)
}

func TestAppRun_WhenInvalidEmail_Error(t *testing.T) {
	t.Parallel()

	email := "invalid-email"

	readerStub := mocks.NewMockTransactionsReader(t)
	emailSenderMock := mocks.NewMockEmailSender(t)

	sut := accsum.New(accsum.Config{
		Email:              email,
		TransactionsReader: readerStub,
		EmailSender:        emailSenderMock,
	})

	// Act
	err := sut.Run()

	// Assert
	require.Error(t, err)
	assert.ErrorIs(t, err, accsum.ErrInvalidEmail)
}

func buildTransaction(id int, month time.Month, value string) model.Transaction {
	return model.Transaction{
		ID:     id,
		Date:   time.Date(2024, month, 1, 0, 0, 0, 0, time.UTC),
		Amount: decimal.MustParse(value),
	}
}

func buildTransactionPerMonth(nonZeroMonths []time.Month, countsPerMonth []int) map[time.Month]int {
	transactionsPerMonth := make(map[time.Month]int)

	for i, month := range nonZeroMonths {
		transactionsPerMonth[month] = countsPerMonth[i]
	}

	return transactionsPerMonth
}
