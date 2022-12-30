package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/caseyrwebb/go-jwt-auth/app"
	"github.com/caseyrwebb/go-jwt-auth/app/utils"
	"go.uber.org/zap"
)

func main() {

	filename := "logs.log"
	logger := utils.CustomLogger(filename)

	configs := utils.NewConfigurations(logger)

	app := app.New(*logger, configs)
	app.DB.SetDBLogger(logger)
	err := app.DB.Open(configs)
	if err != nil {
		logger.Error(fmt.Sprintf("%s %s %s", "could not connect to the database", "error", err))
		os.Exit(1)
	}

	defer app.DB.Close()

	svr := http.Server{
		Addr:         configs.ServerAddress,
		Handler:      app.Router,
		ErrorLog:     zap.NewStdLog(logger),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		logger.Info(fmt.Sprintf("%s %s", "starting the server at port", configs.ServerAddress))

		err := svr.ListenAndServe()
		if err != nil {
			logger.Error(fmt.Sprintf("%s %s %s", "could not start the server", "error", err))
			os.Exit(1)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)

	sig := <-c
	logger.Info(fmt.Sprintf("%s %s %s", "shutting down the server", "received signal", sig))

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	svr.Shutdown(ctx)

}
