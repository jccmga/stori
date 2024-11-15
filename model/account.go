package model

import (
	"time"

	"github.com/govalues/decimal"
)

type AccountSummary struct {
	Email                string
	TotalBalance         decimal.Decimal
	AverageDebitAmount   decimal.Decimal
	AverageCreditAmount  decimal.Decimal
	TransactionsPerMonth map[time.Month]int
}
