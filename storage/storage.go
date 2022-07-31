package storage

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/Jibaru/go-db/pkg/invoice"
	"github.com/Jibaru/go-db/pkg/invoiceheader"
	"github.com/Jibaru/go-db/pkg/invoiceitem"
	"github.com/Jibaru/go-db/pkg/product"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var (
	db   *sql.DB
	once sync.Once
)

// Driver of storage
type Driver string

// Drivers
const (
	MySQL    Driver = "MYSQL"
	Postgres Driver = "POSTGRES"
)

func New(d Driver) {
	switch d {
	case MySQL:
		newMySQLDB()
	case Postgres:
		newPostgresDB()
	}
}

// DAOProduct factory of product storage
func DAOProduct(driver Driver) (product.Storage, error) {
	switch driver {
	case Postgres:
		return newPsqlProduct(db), nil
	case MySQL:
		return newMySQLProduct(db), nil
	default:
		return nil, fmt.Errorf("Driver not found")
	}
}

// DAOInvoiceHeader factory of invoiceheader storage
func DAOInvoiceHeader(driver Driver) (invoiceheader.Storage, error) {
	switch driver {
	case Postgres:
		return newPsqlInvoiceHeader(db), nil
	case MySQL:
		return newMySQLInvoiceHeader(db), nil
	default:
		return nil, fmt.Errorf("Driver not found")
	}
}

// DAOInvoiceItem factory of invoiceitem storage
func DAOInvoiceItem(driver Driver) (invoiceitem.Storage, error) {
	switch driver {
	case Postgres:
		return newPsqlInvoiceItem(db), nil
	case MySQL:
		return newMySQLInvoiceItem(db), nil
	default:
		return nil, fmt.Errorf("Driver not found")
	}
}

// DAOInvoice factory of invoice storage
func DAOInvoice(
	driver Driver,
	headerStorage invoiceheader.Storage,
	itemStorage invoiceitem.Storage,
) (invoice.Storage, error) {
	switch driver {
	case Postgres:
		return newPsqlInvoice(
			db,
			headerStorage,
			itemStorage,
		), nil
	case MySQL:
		return newMySQLInvoice(
			db,
			headerStorage,
			itemStorage,
		), nil
	default:
		return nil, fmt.Errorf("Driver not found")
	}
}

func newPostgresDB() {
	once.Do(func() {
		var err error

		var user, password, port, dbName string = goDotEnvVariable("POSTGRES_USER"),
			goDotEnvVariable("POSTGRES_PASSWORD"),
			goDotEnvVariable("POSTGRES_PORT"),
			goDotEnvVariable("POSTGRES_DB")

		db, err = sql.Open(
			"postgres",
			fmt.Sprintf(
				"postgres://%s:%s@localhost:%s/%s?sslmode=disable",
				user,
				password,
				port,
				dbName,
			),
		)

		if err != nil {
			log.Fatalf("can't open db: %v", err)
		}

		// defer db.Close()

		if err = db.Ping(); err != nil {
			log.Fatalf("can't do ping: %v", err)
		}

		fmt.Println("connected to postgres")
	})
}

func newMySQLDB() {
	once.Do(func() {
		var err error

		var user, password, port, dbName string = goDotEnvVariable("MYSQL_USER"),
			goDotEnvVariable("MYSQL_PASSWORD"),
			goDotEnvVariable("MYSQL_PORT"),
			goDotEnvVariable("MYSQL_DB")

		db, err = sql.Open(
			"mysql",
			fmt.Sprintf(
				"%s:%s@tcp(localhost:%s)/%s?parseTime=true",
				user,
				password,
				port,
				dbName,
			),
		)

		if err != nil {
			log.Fatalf("can't open db: %v", err)
		}

		// defer db.Close()

		if err = db.Ping(); err != nil {
			log.Fatalf("can't do ping: %v", err)
		}

		fmt.Println("connected to mysql")
	})
}

// Pool return a unique instance of db
func Pool() *sql.DB {
	return db
}

// stringToNull
func stringToNull(s string) sql.NullString {
	null := sql.NullString{
		String: s,
	}

	if null.String != "" {
		null.Valid = true
	}

	return null
}

// stringToNull
func timeToNull(t time.Time) sql.NullTime {
	null := sql.NullTime{
		Time: t,
	}

	if !null.Time.IsZero() {
		null.Valid = true
	}

	return null
}

// use godot package to load/read the .env file and
// return the value of the key
func goDotEnvVariable(key string) string {

	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	return os.Getenv(key)
}
