package interfaces

import (
	"encoding/json"
	"evolve/services/transaction/application"
	"net/http"
)

type HTTPHandler struct {
	TransactionApplication application.TransactionApplication
}

func NewTransactionHTTPHandler(transactionApplication application.TransactionApplication) *HTTPHandler {
	return &HTTPHandler{
		TransactionApplication: transactionApplication,
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
