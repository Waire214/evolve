package interfaces

import (
	"encoding/json"
	"evolve/services/transaction/domain"
	"log"
	"net/http"
)

func (handler *HTTPHandler) CreateTransfer(w http.ResponseWriter, r *http.Request) {
	transfer := domain.Transfer{}
	err := json.NewDecoder(r.Body).Decode(&transfer)
	if err != nil {
		log.Println("unable to decode empty struct")
		return
	}
	resp, err := handler.TransactionApplication.CreateTransfer(r.Context(), transfer)
	if err != nil {
		encodeResult(w, err.Error(), http.StatusOK)
		return
	}
	resp.Message = "successful transfer"
	encodeResult(w, resp, http.StatusOK)
}
