package main

import (
	"database/sql"
	"fmt"
	"github.com/incazteca/services"
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

	mux := http.NewServeMux()

	homeDir, _ := os.UserHomeDir()
	mux.Handle(
		"/storage/",
		http.StripPrefix("/storage/", mux.FileServer(http.Dir(homeDir+"/storage"))),
	)
	mux.Handle(
		"/static/",
		http.StripPrefix("/static/", mux.FileServer(http.Dir("static"))),
	)

	http.HandleFunc("/", env.getPhotos)
	http.HandleFunc("/photo", env.handlePhoto)

	photoService := services.NewPhotoService(db)
	photoHandler := routers.NewPhotoHandler(mux, photoService)

	fmt.Println("Start server")
	log.Fatal(mux.ListenAndServe(":8080", nil))
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
