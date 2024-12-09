package main

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/mysqlstore"
	"github.com/alexedwards/scs/v2"
	"github.com/bernhardson/stub/internal/log"
	"github.com/bernhardson/stub/internal/repo"
	"github.com/bernhardson/stub/internal/web"
	"github.com/joho/godotenv"
)

func main() {
	if err := run(os.Getenv); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run(getEnv func(string) string) error {

	srv, err := startServer(getEnv)
	if err != nil {
		return err
	}

	certPath := getEnv("ROOT_DIR") + getEnv("TLS_CERT_PATH")
	keyPath := getEnv("ROOT_DIR") + getEnv("TLS_KEY_PATH")
	err = srv.ListenAndServeTLS(certPath, keyPath)
	if err != nil {
		return err
	}

	return nil
}

func startServer(getEnv func(string) string) (*http.Server, error) {
	//load environment variables from .env ino runtime environment
	err := godotenv.Overload()
	if err != nil {
		return nil, err
	}
	//setup custom logger
	logger, err := log.NewCustomLogger(os.Stdout, getEnv("LOG_LEVEL"))
	if err != nil {
		return nil, err
	}
	//establish database connection
	datasource := getEnv("DATA_SOURCE")
	logger.Info().Msgf("connecting to db type %s", datasource)
	dsn := repo.GetConfig(datasource)
	db, err := repo.Connect(datasource, dsn)
	if err != nil {
		return nil, err
	}
	userRepo, err := repo.UserRepoFactory(datasource, db)
	if err != nil {
		return nil, err
	}

	//create session manager
	sessionManager := scs.New()
	sessionManager.Store = mysqlstore.New(db)
	sessionManager.Lifetime = 12 * time.Hour

	//create application struct
	app := web.Application{
		UserRepo:       userRepo,
		Logger:         logger,
		SessionManager: sessionManager,
	}

	//obviously create  start the http server
	addr := getEnv("SERVER")
	logger.Info().Msgf("starting server on %s", addr)
	tlsConfig := &tls.Config{CurvePreferences: []tls.CurveID{tls.X25519, tls.CurveP256}}
	srv := &http.Server{
		Addr:           addr,
		Handler:        app.Routes(),
		IdleTimeout:    time.Minute,
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 524288,
		TLSConfig:      tlsConfig,
	}
	return srv, nil
}
