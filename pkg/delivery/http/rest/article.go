package rest

import (
	"encoding/json"
	"net/http"

	"article/pkg/adding"
	"article/pkg/delivery/http/response"

	"github.com/go-playground/validator"
)

// addArticle returns a handler for POST /article requests
func addArticle(as adding.Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var aa adding.Article

		if err := json.NewDecoder(r.Body).Decode(&aa); err != nil {
			response.JSON(w, http.StatusBadRequest, response.Payload{IsSuccess: false, Message: http.StatusText(http.StatusBadRequest), Error: err.Error()})
			return
		}

		validate = validator.New()
		if err := validate.Struct(aa); err != nil {
			response.JSON(w, http.StatusBadRequest, response.Payload{IsSuccess: false, Message: http.StatusText(http.StatusBadRequest), Error: err.Error()})
			return
		}

		_, err := as.AddArticle(r.Context(), aa)
		if err != nil {
			response.JSON(w, http.StatusInternalServerError, response.Payload{IsSuccess: false, Message: http.StatusText(http.StatusInternalServerError), Error: err.Error()})

			return
		}

		response.JSON(w, http.StatusCreated, response.Payload{IsSuccess: true, Message: http.StatusText(http.StatusCreated)})
	}
}
