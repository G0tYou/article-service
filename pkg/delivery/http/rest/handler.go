package rest

import (
	"net/http"

	"article/config"
	"article/pkg/adding"

	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
)

var validate *validator.Validate

func Handler(as adding.Service) http.Handler {

	r := mux.NewRouter()
	r = r.PathPrefix("/" + config.AppName + "/v1").Subrouter()

	// Adding
	r.HandleFunc("/article", addArticle(as)).Methods("POST")

	return r
}
