package connectDB

import (
	"fmt"
	"database/sql"
	_ "github.com/lib/pq"
)

const (
    host     = "localhost"
    port     = 5432
    user     = "postgres"
    password = "12345"
    dbname   = "librarian"
)

func GetDB() (*sql.DB, error) {
	connection := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
	host, port, user, password, dbname)
	
	db, err := sql.Open("postgres", connection)

	fmt.Println("Connected with database")

	return db, err
}