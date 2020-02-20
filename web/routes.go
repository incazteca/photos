package web

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"os"
)

const tenMBLimit int64 = (1000 * 1000) + 1

// Env Hold environment related data, TODO: also have it hold storage path
type Env struct {
	db *sql.DB
}

// SuccessCreate has the response for succesful creates
type SuccessCreate struct {
	RecordID int `json:"record_id"`
}

var t *template.Template

// StartServer begin web server
func StartServer(db *sql.DB) {
	env := &Env{db}
	env.routes()

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func (env *Env) routes() {
	homeDir, _ := os.UserHomeDir()
	http.Handle(
		"/storage/",
		http.StripPrefix("/storage/", http.FileServer(http.Dir(homeDir+"/storage"))),
	)
	http.Handle(
		"/static/",
		http.StripPrefix("/static/", http.FileServer(http.Dir("static"))),
	)

	http.HandleFunc("/", env.getPhotos)
	http.HandleFunc("/photo", env.handlePhoto)
}
