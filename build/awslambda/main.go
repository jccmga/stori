package main

import (
	"context"
	"encoding/json"
	"os"
	"strconv"
	"strings"

	"github.com/aws/aws-lambda-go/lambda"

	"stori/accountsummary"
	"stori/adapters/emailsender"
	"stori/adapters/filereader"
	"stori/adapters/repository"
)

const (
	s3FilePathPrefix = "s3://"

	EmailHost     = "EMAIL_HOST"
	EmailPort     = "EMAIL_PORT"
	EmailUsername = "EMAIL_USERNAME"
	EmailPassword = "EMAIL_PASSWORD"
)

type Params struct {
	FilePath string `json:"filepath"`
	Email    string `json:"email"`
}

func handleRequest(_ context.Context, event json.RawMessage) error {
	var params Params
	if err := json.Unmarshal(event, &params); err != nil {
		return err
	}

	emailSender, err := buildEmailSender()
	if err != nil {
		panic(err)
	}

	reader := buildTransactionsReader(params.FilePath)
	noopRepo := repository.New(nil)

	application := accountsummary.New(accountsummary.Config{
		Email:              params.Email,
		FilePath:           params.FilePath,
		TransactionsReader: reader,
		EmailSender:        emailSender,
		Repository:         noopRepo,
	})

	if errRun := application.Run(); errRun != nil {
		panic(errRun)
	}
	return nil
}

func buildTransactionsReader(filepath string) accountsummary.TransactionsReader {
	if strings.HasPrefix(filepath, s3FilePathPrefix) {
		return filereader.NewS3Reader(filepath)
	}

	return filereader.NewLocalReader(filepath)
}

func buildEmailSender() (accountsummary.EmailSender, error) {
	host := os.Getenv(EmailHost)
	port, err := strconv.Atoi(os.Getenv(EmailPort))
	if err != nil {
		return nil, err
	}

	username := os.Getenv(EmailUsername)
	password := os.Getenv(EmailPassword)

	return emailsender.New(emailsender.Config{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
	}), nil
}

func main() {
	lambda.Start(handleRequest)
}
