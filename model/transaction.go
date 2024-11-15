package model

import (
	"time"

	"github.com/govalues/decimal"
)

type Transaction struct {
	ID     int
	Date   time.Time
	Amount decimal.Decimal
}
