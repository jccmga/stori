package filereader

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/govalues/decimal"

	"stori/model"
)

func readTransactions(file *os.File) ([]model.Transaction, error) {
	fail := func(err error) ([]model.Transaction, error) {
		return nil, fmt.Errorf("filereader: Local: readTransactions: %w", err)
	}

	csvReader := csv.NewReader(file)
	csvReader.FieldsPerRecord = numberOfColumns

	data, err := csvReader.ReadAll()
	if err != nil {
		return fail(fmt.Errorf("%w: %w", ErrInvalidFile, err))
	}

	if isNotValidHeader(data[0]) {
		return fail(ErrInvalidHeader)
	}

	transactions := make([]model.Transaction, 0)
	for _, row := range data[1:] {
		transaction, trError := buildTransaction(row)
		if trError != nil {
			return fail(trError)
		}

		transactions = append(transactions, transaction)
	}

	return transactions, nil
}

func isNotValidHeader(header []string) bool {
	return !isValidHeader(header)
}

func isValidHeader(header []string) bool {
	return strings.ToLower(header[0]) == "id" &&
		strings.ToLower(header[1]) == "date" &&
		strings.ToLower(header[2]) == "transaction"
}

func buildTransaction(rowData []string) (model.Transaction, error) {
	fail := func(err error) (model.Transaction, error) {
		return model.Transaction{}, fmt.Errorf("filereader: buildTransaction: %w", err)
	}

	id, err := strconv.Atoi(rowData[0])
	if err != nil {
		return fail(err)
	}

	date, err := parseDate(rowData[1])
	if err != nil {
		return fail(err)
	}

	amount, err := decimal.Parse(rowData[2])
	if err != nil {
		return fail(ErrInvalidAmount)
	}

	return model.Transaction{
		ID:     id,
		Date:   date,
		Amount: amount,
	}, nil
}

func parseDate(datestr string) (time.Time, error) {
	fail := func(err error) (time.Time, error) {
		return time.Time{}, fmt.Errorf("filereader: parseDate %s: %w", datestr, err)
	}

	validDateLayouts := []string{"2006/01/02", "01/02", "1/02", "1/2"}

	for _, layout := range validDateLayouts {
		if date, err := time.Parse(layout, datestr); err == nil {
			return date, nil
		}
	}

	return fail(ErrInvalidDateFormat)
}
