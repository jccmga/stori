package transactions

import (
	"fmt"
	"time"

	"github.com/govalues/decimal"

	"stori/model"
)

type ExecutionResults struct {
	TotalBalance decimal.Decimal

	TransactionsPerMonth map[time.Month]int

	DebitTransactionsCount  int
	CreditTransactionsCount int
	TotalDebitAmount        decimal.Decimal
	TotalCreditAmount       decimal.Decimal

	MinimumBalance decimal.Decimal
}

func newExecutionResults() ExecutionResults {
	transPerMonth := make(map[time.Month]int)
	for i := time.January; i <= time.December; i++ {
		transPerMonth[i] = 0
	}

	return ExecutionResults{
		TotalBalance:         decimal.Zero,
		TransactionsPerMonth: make(map[time.Month]int),
		TotalDebitAmount:     decimal.Zero,
		TotalCreditAmount:    decimal.Zero,
		MinimumBalance:       decimal.Zero,
	}
}

func Process(transactions []model.Transaction) (ExecutionResults, error) {
	fail := func(err error) (ExecutionResults, error) {
		return ExecutionResults{}, fmt.Errorf("transactions: Process: %w", err)
	}

	res := newExecutionResults()
	var err error

	for _, transaction := range transactions {
		if res, err = accountTransaction(res, transaction); err != nil {
			return fail(err)
		}

		res.TransactionsPerMonth[transaction.Date.Month()]++
	}

	return res, nil
}

func accountTransaction(res ExecutionResults, transaction model.Transaction) (ExecutionResults, error) {
	fail := func(err error) (ExecutionResults, error) {
		return ExecutionResults{}, fmt.Errorf("transactions: accountTransaction: %w", err)
	}

	amount := transaction.Amount
	var err error

	if isCredit(amount) {
		res.CreditTransactionsCount++
		if res.TotalCreditAmount, err = res.TotalCreditAmount.Add(amount); err != nil {
			return fail(err)
		}
	} else {
		res.DebitTransactionsCount++
		if res.TotalDebitAmount, err = res.TotalDebitAmount.Add(amount); err != nil {
			return fail(err)
		}
	}

	if res.TotalBalance, err = res.TotalBalance.Add(amount); err != nil {
		return fail(err)
	}

	return res, nil
}

func isCredit(amount decimal.Decimal) bool {
	return amount.Sign() > 0
}
