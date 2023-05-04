package adapt

import (
	"context"
	"database/sql"
	"evolve/adapt/routes"
	"evolve/services"
	accountApplications "evolve/services/account/application"
	accountInfrastructures "evolve/services/account/infrastructure"
	accountInterfaces "evolve/services/account/interfaces"
	transactionApplications "evolve/services/transaction/application"
	transactionInfrastructures "evolve/services/transaction/infrastructure"
	transactionInterfaces "evolve/services/transaction/interfaces"
	"fmt"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"log"
	"net/http"
)

type application struct {
	logger       *zap.Logger
	config       *ServerConfig
	db           *sql.DB
	router       *chi.Mux
	repositories services.Repositories
}

// New instances a new application
// The application contains all the related components that allow the execution of the service
func New(logger *zap.Logger) (*application, error) {
	var app application
	var err error

	app.logger = logger
	app.config, err = app.buildConfig()

	if err != nil {
		return nil, err
	}
	//build application clients
	app.db = app.buildSqlClient()

	if err := app.db.PingContext(context.Background()); err != nil {
		app.logger.Info("msg", zap.String("msg", "failed to ping to database"))
		log.Fatal(err)
	}

	tokenHandler, err := services.NewMiddlewares()
	if err != nil {
		return nil, err
	}

	allInterfaces := app.buildApplicationConnection(*tokenHandler)

	router, tokenHandler, err := routes.SetupRouter(tokenHandler, allInterfaces)
	if err != nil {
		return nil, err
	}
	app.router = router

	return &app, nil
}

// Run executes the application
func (app *application) Run() error {
	defer app.db.Close()

	app.router.Get("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte("Welcome to the evolve credit Server.."))
	})

	svr := http.Server{
		Addr:    fmt.Sprintf(":%d", app.config.HTTPPort),
		Handler: app.router,
	}
	err := svr.ListenAndServe()
	if err != nil {
		return err
	}
	return nil
}

func (app *application) buildConfig() (*ServerConfig, error) {
	return Read(*app.logger)
}

func (app *application) buildSqlClient() *sql.DB {
	db := Database{Config: app.config, Log: app.logger}

	dbConn, err := db.ConnectDB()

	if err != nil {
		log.Fatal(err)
	}

	if err := db.RunMigration(dbConn); err != nil {
		log.Fatal(err)
	}

	return dbConn
}

func (app *application) buildApplicationConnection(tokenHandler services.TokenHandler) *routes.AllHTTPHandlers {
	transactionPersistence := transactionInfrastructures.NewTransactionPersistence(app.db, app.logger)
	accountPersistence := accountInfrastructures.NewAccountPersistence(app.db, app.logger)
	allRepositories := services.Repositories{
		TransactionRepository: transactionPersistence,
		AccountRepository:     accountPersistence,
	}
	app.repositories = allRepositories

	transactionApplication := transactionApplications.NewTransactionApplication(tokenHandler, allRepositories, app.logger)
	accountApplication := accountApplications.NewAccountApplication(tokenHandler, allRepositories)

	transactionInterface := transactionInterfaces.NewTransactionHTTPHandler(transactionApplication)
	accountInterface := accountInterfaces.NewAccountHTTPHandler(accountApplication)

	allInterfaces := routes.AllHTTPHandlers{
		Transaction: transactionInterface,
		Account:     accountInterface,
	}
	return routes.AllInterfaces(&allInterfaces)
}
