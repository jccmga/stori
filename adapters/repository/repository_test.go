package repository_test

import (
	"database/sql"
	"fmt"
	"log"
	"testing"
	"time"

	_ "github.com/lib/pq"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/pressly/goose/v3"
)

var DB *sql.DB //nolint:gochecknoglobals // For tests splitting

func TestMain(m *testing.M) {
	pool := buildPool()
	resource, databaseURL := buildPostgresContainer(pool)
	pingDatabase(pool, databaseURL)

	defer func() {
		if err := pool.Purge(resource); err != nil {
			log.Fatalf("Could not purge resource: %s", err)
		}
	}()

	runMigrations()

	runTests(m)
}

func buildPool() *dockertest.Pool {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not construct pool: %s", err)
	}

	err = pool.Client.Ping()
	if err != nil {
		log.Fatalf("Could not connect to Docker: %s", err)
	}

	return pool
}

func buildPostgresContainer(pool *dockertest.Pool) (*dockertest.Resource, string) {
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "11",
		Env: []string{
			"POSTGRES_PASSWORD=secret",
			"POSTGRES_USER=user_name",
			"POSTGRES_DB=dbname",
			"listen_addresses = '*'",
		},
	}, func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	hostAndPort := resource.GetHostPort("5432/tcp")
	databaseURL := fmt.Sprintf("postgres://user_name:secret@%s/dbname?sslmode=disable", hostAndPort)

	if errExp := resource.Expire(120); errExp != nil {
		log.Fatalf("Could not set expiry to db container: %s", errExp)
	}

	return resource, databaseURL
}

func pingDatabase(pool *dockertest.Pool, databaseURL string) {
	var err error
	pool.MaxWait = 120 * time.Second
	if errR := pool.Retry(func() error {
		DB, err = sql.Open("postgres", databaseURL)
		if err != nil {
			return err
		}
		return DB.Ping()
	}); errR != nil {
		log.Fatalf("Could not connect to docker: %s", errR)
	}
}

func runMigrations() {
	if err := goose.SetDialect("postgres"); err != nil {
		panic(err)
	}

	if err := goose.Up(DB, "../../migrations"); err != nil {
		panic(err)
	}
}

func runTests(m *testing.M) {
	m.Run()
}
