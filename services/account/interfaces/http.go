package interfaces

import (
	"encoding/json"
	"evolve/services/account/application"
	"net/http"
)

type HTTPHandler struct {
	AccountApplication application.AccountApplication
}

func NewAccountHTTPHandler(AccountApplication application.AccountApplication) *HTTPHandler {
	return &HTTPHandler{
		AccountApplication: AccountApplication,
	}

}

func encodeResult(w http.ResponseWriter, result interface{}, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	data := struct {
		Data interface{} `json:"data"`
	}{
		Data: result,
	}

	err := json.NewEncoder(w).Encode(&data)
	if err != nil {
		return
	}
}
