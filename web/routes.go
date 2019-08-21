package web

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/incazteca/photos/services/photos"
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

func (env *Env) getPhotos(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	if r.Method == "GET" {
		photos, err := photos.FetchAll(env.db)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		t = template.Must(
			template.ParseFiles(
				"templates/common/head.tmpl",
				"templates/photos/index.tmpl",
				"templates/common/footer.tmpl",
			),
		)
		err = t.ExecuteTemplate(w, "index.tmpl", photos)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

func (env *Env) handlePhoto(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/photo" {
		http.NotFound(w, r)
		return
	}

	if r.Method == "POST" {
		recordID, err := photos.Create(env.db, r.Body)

		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(SuccessCreate{recordID})
		return
	}
}
