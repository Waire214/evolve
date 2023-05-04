package interfaces

import (
	"encoding/json"
	"evolve/services/account/domain"
	"log"
	"net/http"
)

func (handler *HTTPHandler) CreateCashCache(w http.ResponseWriter, r *http.Request) {
	cashCache := domain.CashCache{}
	err := json.NewDecoder(r.Body).Decode(&cashCache)
	if err != nil {
		log.Println("unable to decode empty struct")
		return
	}
	resp, err := handler.AccountApplication.CreateCashCache(r.Context(), &cashCache)
	if err != nil {
		encodeResult(w, err.Error(), http.StatusOK)
		return
	}
	encodeResult(w, resp, http.StatusOK)
}
