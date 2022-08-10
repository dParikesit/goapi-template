package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/cenkalti/backoff/v4"
	"github.com/cockroachdb/cockroach-go/v2/crdb"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	_ "github.com/lib/pq"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	db, err := initStore()
	if err != nil {
		log.Fatalf("failed to initialise the store: %s", err)
	}
	defer db.Close()

	if (len(os.Args) == 2 && (os.Args[1] == "migrate")) {
		err := migrateDb(db)
		if err != nil {
			log.Fatalf("Error trying to migrate database: %s", err)
			os.Exit(-1)
		}
		os.Exit(0)
	}

	e.GET("/", func(c echo.Context) error {
		return rootHandler(db, c)
	})

	e.GET("/ping", func(c echo.Context) error {
		return c.JSON(http.StatusOK, struct{ Status string }{Status: "OK"})
	})

	e.POST("/send", func(c echo.Context) error {
		return sendHandler(db, c)
	})

	httpPort := os.Getenv("PORT")
	if httpPort == "" {
		httpPort = "8080"
	}

	e.Logger.Fatal(e.Start(":" + httpPort))
}

type Message struct {
	Value string `json:"value"`
}

func initStore() (*sql.DB, error) {

	coherence_dev, coherence_present := os.LookupEnv("COHERENCE_DEV")
	dbhost, dbhost_present := os.LookupEnv("DB_HOST")
	dbname := os.Getenv("DB_NAME")
	dbuser := os.Getenv("DB_USER")
	dbpass := os.Getenv("DB_PASSWORD")

	pgConnString := ""
	dbsocket := ""
	dbendpoint := ""
	dbport := ""

	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)

		if strings.HasSuffix(pair[0], fmt.Sprintf("DB1_SOCKET")) {
			dbsocket = pair[1]
		}
		if strings.HasSuffix(pair[0], fmt.Sprintf("DB1_ENDPOINT")) {
			dbendpoint = pair[1]
		}
		if strings.HasSuffix(pair[0], fmt.Sprintf("DB1_PORT")) {
			dbport = pair[1]
		}
	}

	if (coherence_present && coherence_dev == "true") {
		if (dbhost_present && dbhost != "") {
			pgConnString = fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=disable", dbhost, dbport, dbname, dbuser, dbpass)
		} else {
			pgConnString = fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=disable", "localhost", dbport, dbname, dbuser, dbpass)
		}
	} else {
		if (dbendpoint != "") {
			pgConnString = fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=disable", dbendpoint, dbport, dbname, dbuser, dbpass)
		} else {
			pgConnString = fmt.Sprintf("host=%s dbname=%s user=%s password=%s sslmode=disable", dbsocket, dbname, dbuser, dbpass)
		}
	}

	var (
		db  *sql.DB
		err error
	)
	openDB := func() error {
		db, err = sql.Open("postgres", pgConnString)
		return err
	}

	err = backoff.Retry(openDB, backoff.NewExponentialBackOff())
	if err != nil {
		return nil, err
	}

	return db, nil
}


func migrateDb(db *sql.DB) (error) {
	fmt.Println("Migrating database")
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}
	m, err := migrate.NewWithDatabaseInstance(
		"file:///app/db/migrations",
		"postgres", driver)
	if err != nil {
		return err
	}
	m.Up() // or m.Step(2) if you want to explicitly set the number of migrations to run
	return(err)
}

func rootHandler(db *sql.DB, c echo.Context) error {
	r, err := countRecords(db)
	if err != nil {
		return c.HTML(http.StatusInternalServerError, err.Error())
	}
	return c.HTML(http.StatusOK, fmt.Sprintf("Hello, Docker! (%d)\n", r))
}

func sendHandler(db *sql.DB, c echo.Context) error {

	m := &Message{}

	if err := c.Bind(m); err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	err := crdb.ExecuteTx(context.Background(), db, nil,
		func(tx *sql.Tx) error {
			_, err := tx.Exec(
				"INSERT INTO message (value) VALUES ($1) ON CONFLICT (value) DO UPDATE SET value = excluded.value",
				m.Value,
			)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, err)
			}
			return nil
		})

	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, m)
}

func countRecords(db *sql.DB) (int, error) {

	rows, err := db.Query("SELECT COUNT(*) FROM message")
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		if err := rows.Scan(&count); err != nil {
			return 0, err
		}
		rows.Close()
	}

	return count, nil
}
