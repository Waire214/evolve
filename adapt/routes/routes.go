package routes

import (
	"evolve/services"
	accountInterface "evolve/services/account/interfaces"
	transactionInterface "evolve/services/transaction/interfaces"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"net/http"
)

type AllHTTPHandlers struct {
	Transaction *transactionInterface.HTTPHandler
	Account     *accountInterface.HTTPHandler
}

func AllInterfaces(interfaces *AllHTTPHandlers) *AllHTTPHandlers {
	return &AllHTTPHandlers{Transaction: interfaces.Transaction, Account: interfaces.Account}
}

func SetupRouter(tokenHandler *services.TokenHandler, interfaces *AllHTTPHandlers) (*chi.Mux, *services.TokenHandler, error) {
	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))
	router.Use(middleware.Logger)

	transactionRouter := buildTransactionEndpoints(*interfaces.Transaction, tokenHandler)
	accountRouter := buildAccountEndpoints(*interfaces.Account, tokenHandler)

	router.Route("/evolve", func(r chi.Router) {
		r.Mount("/transaction", transactionRouter)
		r.Mount("/account", accountRouter)
	})

	return router, tokenHandler, nil
}

func buildTransactionEndpoints(transaction transactionInterface.HTTPHandler, tokenHandler *services.TokenHandler) http.Handler {
	router := chi.NewRouter()
	router.Use(tokenHandler.ValidateMiddleware)
	router.Post("/make_transfer", transaction.CreateTransfer)
	return router
}

func buildAccountEndpoints(account accountInterface.HTTPHandler, tokenHandler *services.TokenHandler) http.Handler {
	router := chi.NewRouter()
	router.Use(tokenHandler.ValidateMiddleware)
	router.Post("/create_cache", account.CreateCashCache)
	return router
}
