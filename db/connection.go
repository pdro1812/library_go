package db


import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
)

var DB *sql.DB

func ConnectPostgres(user, password, dbname, host string, port int) error {
	psqlInfo := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname,
	)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		return err
	}
	if err = db.Ping(); err != nil {
		return err
	}
	DB = db
	return nil
}