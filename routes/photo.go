package routes

import (
	"encoding/json"
	"fmt"
	"github.com/incazteca/photos/services"
	"html/template"
	"net/http"
)

type PhotoHandler struct {
	photoService services.PhotoService
}

//NewPhotoHandler return a new handler for photos
func NewPhotoHandler(mux *http.ServeMux, ps services.PhotoService) PhotoHandler {
	handler := PhotoHandler{ps}

	mux.HandleFunc("/", handler.getPhotos)
	mux.HandleFunc("/photo", handler.handlePhoto)
	return handler
}

func (handler *PhotoHandler) getPhotos(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	if r.Method == "GET" {
		photos, err := handler.photoService.FetchAll()

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		t := template.Must(
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

func (handler *PhotoHandler) handlePhoto(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/photo" {
		http.NotFound(w, r)
		return
	}

	if r.Method == "POST" {
		recordID, err := handler.photoService.StorePhoto(r.Body)

		if err != nil {
			fmt.Println(err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]int{"record_id": recordID})
		return
	}
}
