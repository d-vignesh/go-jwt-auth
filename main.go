package main 

import (
	"net/http"
	"time"
	"os"
	"os/signal"
	"context"

	"github.com/gorilla/mux"
	"github.com/hashicorp/go-hclog"
	"github.com/d-vignesh/go-jwt-auth/handlers"
	"github.com/d-vignesh/go-jwt-auth/data"
	"github.com/d-vignesh/go-jwt-auth/service"
	"github.com/d-vignesh/go-jwt-auth/utils"
)

const schema = `
		create table if not exists users (
			id varchar(36) not null,
			email varchar(225) not null unique,
			username varchar(225),
			password varchar(225) not null,
			primary key (id)
		);
`

func main() {

	logger := hclog.New(&hclog.LoggerOptions{
			Name: "user-auth",
			Level: hclog.LevelFromString("DEBUG"),
	})

	configs := utils.NewConfigurations(logger)

	// validator contains all the methods that are need to validate the user json in request
	validator := data.NewValidation()

	db, err := data.NewConnection(configs, logger)
	if err != nil {
		logger.Error("unable to connect to db", "error", err)
		panic(err)
	}
	defer db.Close()
	
	db.MustExec(schema)

	// repository contains all the methods that interact with DB to perform CURD operations for user.
	repository := data.NewPostgresRepository(db, logger)

	// authService contains all methods that help in authorizing a user request
	authService := service.NewAuthService(logger, configs)

	// UserHandler encapsulates all the services related to user
	uh := handlers.NewUserHandler(logger, validator, repository, authService)

	// create a serve mux
	sm := mux.NewRouter()

	// register handlers
	postR := sm.Methods(http.MethodPost).Subrouter()
	postR.HandleFunc("/signup", uh.Signup)
	postR.HandleFunc("/login", uh.Login)
	postR.Use(uh.MiddlewareValidateUser)

	// used the PathPrefix as workaround for scenarios where all the 
	// get requests my use the ValidateAccessToken middleware except 
	// the /refresh-token request which has to use ValidateRefreshToken middleware
	refToken := sm.PathPrefix("/refresh-token").Subrouter()
	refToken.HandleFunc("", uh.RefreshToken)
	refToken.Use(uh.MiddlewareValidateRefreshToken)

	getR := sm.Methods(http.MethodGet).Subrouter()
	getR.HandleFunc("/greet", uh.Greet)
	getR.Use(uh.MiddlewareValidateAccessToken)

	// create a server
	svr := http.Server{
		Addr:	 	  configs.ServerPort,
		Handler: 	  sm,
		ErrorLog:	  logger.StandardLogger(&hclog.StandardLoggerOptions{}),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// start the server
	go func() {
		logger.Info("starting the server at port", configs.ServerPort)

		err := svr.ListenAndServe()
		if err != nil {
			logger.Error("could not start the server", "error", err)
			os.Exit(1)
		}
	}()

	// look for interrupts for graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, os.Kill)

	sig := <-c
	logger.Info("shutting down the server", "received signal", sig)

	//gracefully shutdown the server, waiting max 30 seconds for current operations to complete
	ctx, _ := context.WithTimeout(context.Background(), 30 * time.Second)
	svr.Shutdown(ctx)
}