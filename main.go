package main

import (
	"database/sql"
	"fmt"
	"github.com/incazteca/photos/web"
	_ "github.com/lib/pq"
	"os"
	"strconv"
)

const (
	conn = "host=%s port=%d user=%s password=%s dbname=photos"
)

func main() {
	fmt.Println("Connecting to DB...")
	db := connectToDb()
	defer db.Close()

	fmt.Println("Start server")
	web.StartServer(db)
}

// TODO, best thing to do here is probably keep the connection string and then
// connect where needed?
func connectToDb() *sql.DB {
	dbPassword := os.Getenv("DB_PASSWORD")
	rawDbPort := os.Getenv("DB_PORT")
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")

	if dbHost == "" {
		dbHost = "127.0.0.1"
	}

	dbPort, err := strconv.Atoi(rawDbPort)
	if err != nil {
		dbPort = 5432
	}

	if dbUser == "" {
		dbUser = "photos_user"
	}

	connString := fmt.Sprintf(conn, dbHost, dbPort, dbUser, dbPassword)
	db, err := sql.Open("postgres", connString)
	if err != nil {
		panic(err)
	}

	fmt.Println("Connected to DB")
	return db
}
