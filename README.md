# Stori - Challenge

## Installation
### Locally
To run locally, make sure to have [Golang v1.23.2+](https://golang.org/) installed.
Also make sure to have some environment variables set:

```
export EMAIL_HOST='smtp.gmail.com'
export EMAIL_PORT=587
export EMAIL_USERNAME='sender's email address
export EMAIL_PASSWORD='sender's email password for the app'
export DB_USER=stori
export DB_PASS=storipwd
export DB_NAME=storidb
```
Then, run:
```
make build
./bin/stori -email 'jcamilo.36@gmail.com' -filepath ./data/several_transactions.csv
./bin/stori -email 'jcamilo.36@gmail.com' -filepath s3://jcamilostori/several_transactions.csv
```

### Docker

Make sure to have [Docker](https://www.docker.com/) installed.
Also make sure to have your `./build/.env` file set with proper values.

Then, run:

```
./run.sh
```

### Tests

Distinction between unit tests and integration tests follow definition from Khorikov
[Unit Testing Principles, Patterns and Practices](https://www.manning.com/books/unit-testing),
in particular the use of shared resources or not.

To run all tests (unit tests and integration tests), execute:

```
make test-all
```

To run unit tests only, execute:

```
make test
```

### Linting

To run linting, execute:

```
make lint
```

Configuration for linting is in `.golangci.yml`.

## Frameworks and libs used

### Decimals
As no decimals library out of the box is provided by the language, I opt to use a third-party
library: [govalues](https://github.com/govalues/decimal). More about this decision [here](./docs/architecture/decisions/0004-handling-decimals.md).

### Emails
I opt to use [go-mail](https://github.com/wneessen/go-mail). As I don't have a domain for the challenge, I used my Gmail account.
To generate a new password, go to [App passwords](https://support.google.com/mail/answer/185833?hl=en#:~:text=Create%20and%20manage%20your%20app%20passwords).
Replace the `EMAIL_PASSWORD` environment variable with the generated password.

### Managed dependencies
To test the database accesses, we used [Dockertest](https://github.com/ory/dockertest) because of its ease of use in
this particular case. More about this decision [here](./docs/architecture/decisions/0007-testing-the-database.md).

go get github.com/aws/aws-sdk-go/aws
go get github.com/aws/aws-sdk-go/aws/session
go get github.com/aws/aws-sdk-go/service/s3

go install github.com/pressly/goose/v3/cmd/goose@latest
```
go install github.com/vektra/mockery/v2@v2.46.3
go get -u github.com/ory/dockertest/v3
go get github.com/jmoiron/sqlx
```

go get github.com/aws/aws-lambda-go/lambda