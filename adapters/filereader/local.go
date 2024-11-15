package filereader

import (
	"fmt"
	"os"

	"stori/model"
)

const numberOfColumns = 3

type Local struct {
	filePath string
}

func NewLocalReader(filePath string) Local {
	return Local{filePath: filePath}
}

//nolint:unparam
func (reader Local) ReadTransactions() ([]model.Transaction, error) {
	fail := func(err error) ([]model.Transaction, error) {
		return nil, fmt.Errorf("filereader: Local: ReadTransactions: %w", err)
	}

	file, err := openCVSFile(reader.filePath)
	if err != nil {
		return fail(err)
	}

	defer file.Close()

	transactions, err := readTransactions(file)
	if err != nil {
		return fail(err)
	}

	return transactions, nil
}

func openCVSFile(filePath string) (*os.File, error) {
	fail := func(err error) (*os.File, error) {
		return nil, fmt.Errorf("filereader: Local: openCVSFile: %w", err)
	}

	file, err := os.Open(filePath)
	if err != nil {
		return fail(ErrFileNotFound)
	}

	fileInfo, err := file.Stat()
	if err != nil {
		return fail(err)
	}

	if fileInfo.Size() == 0 {
		return fail(ErrFileIsEmpty)
	}

	return file, nil
}
