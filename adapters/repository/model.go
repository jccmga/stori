package repository

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/govalues/decimal"
)

type (
	Transaction struct {
		Date     time.Time       `db:"date"`
		Amount   decimal.Decimal `db:"amount"`
		FilePath string          `db:"file_path"`
	}

	AccountSummary struct {
		Email                string          `db:"email"`
		TotalBalance         decimal.Decimal `db:"total_balance"`
		AverageDebitAmount   decimal.Decimal `db:"average_debit_amount"`
		AverageCreditAmount  decimal.Decimal `db:"average_credit_amount"`
		TransactionsPerMonth TransPerMonth   `db:"transactions_per_month"`
		FilePath             string          `db:"file_path"`
	}

	TransPerMonth map[time.Month]int
)

func (tpm TransPerMonth) Value() (driver.Value, error) {
	return json.Marshal(tpm)
}

func (tpm *TransPerMonth) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &tpm)
}
