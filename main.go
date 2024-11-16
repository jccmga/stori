package main

import (
	"embed"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/jmoiron/sqlx"

	"stori/accountsummary"
	"stori/adapters/emailsender"
	"stori/adapters/filereader"
	"stori/adapters/repository"

	_ "github.com/lib/pq"
	"github.com/pressly/goose/v3"
)

const (
	s3FilePathPrefix = "s3://"

	EmailHost     = "EMAIL_HOST"
	EmailPort     = "EMAIL_PORT"
	EmailUsername = "EMAIL_USERNAME"
	EmailPassword = "EMAIL_PASSWORD"

	migrationsDir = "migrations"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

func main() {
	var email, filepath string
	flag.StringVar(&email, "email", "", "Email address where the results of the process will be sent")
	flag.StringVar(&filepath, "filepath", "", "CSV filepath with transactions to be processed")
	flag.Parse()

	reader := buildTransactionsReader(filepath)

	emailSender, err := buildEmailSender()
	if err != nil {
		panic(err)
	}

	db, errSetup := setupDB()
	if errSetup != nil {
		log.Printf("DB not setup, persistence disabled: %s", errSetup.Error())
	}

	if db != nil {
		defer db.Close()
	}

	repo := repository.New(db)

	application := accountsummary.New(accountsummary.Config{
		Email:              email,
		FilePath:           filepath,
		TransactionsReader: reader,
		EmailSender:        emailSender,
		Repository:         repo,
	})

	if errRun := application.Run(); errRun != nil {
		panic(errRun)
	}
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

func setupDB() (*sqlx.DB, error) {
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")
	dbHost := "localhost"
	if dh := os.Getenv("DB_HOST"); dh != "" {
		dbHost = dh
	}

	dsn := fmt.Sprintf("dbname=%s user=%s password=%s host=%s sslmode=disable", dbName, dbUser, dbPass, dbHost)

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, err
	}

	if err = runMigrations(db); err != nil {
		return nil, err
	}

	return db, nil
}

func runMigrations(db *sqlx.DB) error {
	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	if err := goose.Up(db.DB, migrationsDir); err != nil {
		return err
	}

	return nil
}
