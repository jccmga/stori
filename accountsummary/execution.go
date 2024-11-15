package accountsummary

import (
	"stori/model"
)

type Execution struct {
	AccountSummary model.AccountSummary
	Transactions   []model.Transaction
	FilePath       string
}
