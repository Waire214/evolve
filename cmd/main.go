package main

import (
	"evolve/adapt"
	"evolve/services/logger"
	"fmt"
	"os"
)

func main() {
	appLogger := logger.New()

	app, err := adapt.New(appLogger)
	if err != nil {
		appLogger.Error(fmt.Sprintf("Fatal error creating application: %v", err))
		os.Exit(1)
	}

	if err := app.Run(); err != nil {
		appLogger.Error(fmt.Sprintf("Fatal error running application: %v", err))
		os.Exit(1)
	}
}
