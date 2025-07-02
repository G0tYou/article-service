package response

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Payload struct {
	IsSuccess bool        `json:"is_success"`
	Code      int         `json:"code"`
	Data      interface{} `json:"data,omitempty"`
	Error     string      `json:"error,omitempty"`
	Message   string      `json:"message,omitempty"`
}

func JSON(w http.ResponseWriter, httpStatus int, p Payload) {
	p.Code = httpStatus

	if !p.IsSuccess {
		p.Data = nil
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, GET")
	w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)

	err := json.NewEncoder(w).Encode(p)
	if err != nil {
		fmt.Fprintf(w, "%s", err.Error())
	}

}
