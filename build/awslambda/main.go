package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/aws/aws-lambda-go/events"
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

	filePathKey = "filepath"
	emailKey    = "email"
)

func handleRequest(_ context.Context, event *events.APIGatewayV2HTTPRequest) (*string, error) {
	body := map[string]string{}

	err := json.Unmarshal([]byte(event.Body), &body)
	if err != nil {
		panic(err)
	}

	filePath, ok := body[filePathKey]
	if !ok {
		panic(errors.New("filepath not found"))
	}

	email, ok := body[emailKey]
	if !ok {
		panic(errors.New("email not found"))
	}

	emailSender, err := buildEmailSender()
	if err != nil {
		panic(err)
	}

	reader := buildTransactionsReader(filePath)
	noopRepo := repository.New(nil)

	application := accountsummary.New(accountsummary.Config{
		Email:              email,
		FilePath:           filePath,
		TransactionsReader: reader,
		EmailSender:        emailSender,
		Repository:         noopRepo,
	})

	if errRun := application.Run(); errRun != nil {
		panic(errRun)
	}

	successMessage := fmt.Sprintf("Process executed successfully for email: %s and file: %s", email, filePath)

	return &successMessage, nil
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
