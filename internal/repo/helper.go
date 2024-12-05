package repo

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/go-sql-driver/mysql"
)

// Reads config from environment variables and returns dsn string.
func GetConfig(datasource string) string {

	switch datasource {
	case "psql":
		return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_USER"), os.Getenv("DB_PASS"), os.Getenv("DB_NAME"))
	default:
		cfg := mysql.Config{
			User:   os.Getenv("DB_USER"),
			Passwd: os.Getenv("DB_PASS"),
			Net:    "tcp",
			Addr:   os.Getenv("DB_ADDR"),
			DBName: os.Getenv("DB_NAME"),
			Params: map[string]string{
				"parseTime": "true",
			},
		}
		return cfg.FormatDSN()
	}

}

// Creates database connection.
func Connect(datasource, dsn string) (*sql.DB, error) {

	db, err := sql.Open(datasource, dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
