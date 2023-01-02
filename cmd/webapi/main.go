/*
This API allows using WASA Photo with a variety of functionalities like image upload, user profile photos,
stream and information etc.
Webapi is the executable for the main web server.
It builds a web server around APIs from `service/api`.
Webapi connects to external resources needed (database) and starts two web servers: the API web server, and the debug.
Everything is served via the API web server, except debug variables (/debug/vars) and profiler infos (pprof).

Usage:

	webapi [flags]

Flags and configurations are handled automatically by the code in `load-configuration.go`.

Return values (exit codes):

	0
		The program ended successfully (no errors, stopped by signal)

	> 0
		The program ended due to an error

Note that this program will update the schema of the database to the latest version available (embedded in the
executable during the build).
*/
package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/ardanlabs/conf"
	"github.com/lucaronca/wasa-homework/service/api"
	"github.com/lucaronca/wasa-homework/service/api/controllers"
	"github.com/lucaronca/wasa-homework/service/api/repositories"
	"github.com/lucaronca/wasa-homework/service/api/routes"
	"github.com/lucaronca/wasa-homework/service/api/services"
	"github.com/lucaronca/wasa-homework/service/database"
	"github.com/lucaronca/wasa-homework/service/globaltime"
	_ "github.com/mattn/go-sqlite3"
	"github.com/sirupsen/logrus"
)

// main is the program entry point. The only purpose of this function is to call run() and set the exit code if there is
// any error
func main() {
	if err := run(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "error: ", err)
		os.Exit(1)
	}
}

/*
newHandler returns a new http.Handler for an api Router. Given an api Router, Handler dependencies
and configuration, this function defines the necessary api repositories, services, middlewares, controllers,
and makes sure that the router is configured with the HTTP endpoints and handlers defined in the controllers.
Finally it returns an http.Handler ready to be invoked by an http.Server.
*/
func newHandler(router api.Router, db database.AppDatabase, assetsCfg Assets) http.Handler {
	// Liveness checker
	livenessChecker := api.NewLivenessChecker(db.Ping)

	// Instantiate repositories
	authRepository, _ := repositories.NewAuthRepository(db)
	usersRepository, _ := repositories.NewUsersRepository(db)
	bansRepository, _ := repositories.NewBansRepository(db)
	followsRepository, _ := repositories.NewFollowsRepository(db)
	photosRepository, _ := repositories.NewPhotosRepository(db)
	likesRepository, _ := repositories.NewLikesRepository(db)
	commentsRepository, _ := repositories.NewCommentsRepository(db)

	// Instantiate services
	authService := services.NewAuthService(authRepository, usersRepository)
	bansService := services.NewBansService(usersRepository, bansRepository, followsRepository)
	followsService := services.NewFollowsService(usersRepository, bansRepository, followsRepository)
	photosService := services.NewPhotosService(
		assetsCfg.PhotosDirectory,
		assetsCfg.PhotosUrlPath,
		usersRepository,
		bansRepository,
		photosRepository,
		likesRepository,
		commentsRepository,
		followsRepository,
	)
	usersService := services.NewUsersService(
		usersRepository,
		bansRepository,
		followsRepository,
		photosRepository,
	)
	likesService := services.NewLikesService(
		usersRepository,
		bansRepository,
		likesRepository,
		photosRepository,
	)
	commentsService := services.NewCommentsService(
		usersRepository,
		bansRepository,
		commentsRepository,
		photosRepository,
	)

	// Instantiate middlewares
	tokenAuthMiddleware := routes.NewTokenAuthMiddleware(authService)

	// Instantiate controllers
	loginController := controllers.NewLoginController(authService)
	bansController := controllers.NewBansController(bansService)
	followsController := controllers.NewFollowsController(followsService)
	photosController := controllers.NewPhotosController(photosService)
	usersController := controllers.NewUsersController(usersService)
	likesController := controllers.NewLikesController(likesService)
	commentsController := controllers.NewCommentsController(commentsService)

	// Handler Configuration
	handlerCfg := api.HandlerConfig{
		Photos: api.HandlerConfigPhotos{
			PhotosDirectory: assetsCfg.PhotosDirectory,
			PhotosUrlPath:   assetsCfg.PhotosUrlPath,
		},
		Deps: api.HandlerConfigDependencies{
			LivenessChecker:     livenessChecker,
			TokenAuthMiddleware: tokenAuthMiddleware,
		},
	}

	// Instantiate routes
	handler := router.Handler(
		handlerCfg,
		loginController,
		followsController,
		bansController,
		photosController,
		usersController,
		likesController,
		commentsController,
	)
	return handler
}

// run executes the program. The body of this function should perform the following steps:
// * reads the configuration
// * creates and configure the logger
// * connects to any external resources (like databases, authenticators, etc.)
// * creates an instance of the service/api package
// * starts the principal web server (using the service/api.Router.Handler() for HTTP handlers)
// * waits for any termination event: SIGTERM signal (UNIX), non-recoverable server error, etc.
// * closes the principal web server
func run() error {
	rand.Seed(globaltime.Now().UnixNano())
	// Load Configuration and defaults
	cfg, err := loadConfiguration()
	if err != nil {
		if errors.Is(err, conf.ErrHelpWanted) {
			return nil
		}
		return err
	}

	// Init logging
	logger := logrus.New()
	logger.SetOutput(os.Stdout)
	if cfg.Debug {
		logger.SetLevel(logrus.DebugLevel)
	} else {
		logger.SetLevel(logrus.InfoLevel)
	}

	logger.Infof("application initializing")

	// Start Database
	logger.Println("initializing database support")
	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	dbPath := os.Getenv("DB_PATH")

	if dbPath == "" {
		dbPath = filepath.Join(pwd, "/data")
	}

	dbconn, err := sql.Open(
		"sqlite3",
		fmt.Sprintf(
			"file:%s?_foreign_keys=on",
			filepath.Join(dbPath, cfg.DB.Filename),
		),
	)
	if err != nil {
		logger.WithError(err).Error("error opening SQLite DB")
		return fmt.Errorf("opening SQLite: %w", err)
	}
	defer func() {
		logger.Debug("database stopping")
		_ = dbconn.Close()
	}()
	db, err := database.New(dbconn)
	if err != nil {
		logger.WithError(err).Error("error creating AppDatabase")
		return fmt.Errorf("creating AppDatabase: %w", err)
	}

	// Start (main) API server
	logger.Info("initializing API server")

	// Make a channel to listen for an interrupt or terminate signal from the OS.
	// Use a buffered channel because the signal package requires it.
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// Make a channel to listen for errors coming from the listener. Use a
	// buffered channel so the goroutine can exit if we don't collect this error.
	serverErrors := make(chan error, 1)

	// Create the API router
	apirouter, err := api.New(api.RouterConfig{
		Logger: logger,
	})
	if err != nil {
		logger.WithError(err).Error("error creating the API server instance")
		return fmt.Errorf("creating the API server instance: %w", err)
	}

	assetsCfg := Assets{
		PhotosDirectory: filepath.Join(pwd, cfg.Assets.PhotosDirectory),
		PhotosUrlPath:   cfg.Assets.PhotosUrlPath,
	}

	handler := newHandler(
		apirouter,
		db,
		assetsCfg,
	)

	handler, err = registerWebUI(handler)
	if err != nil {
		logger.WithError(err).Error("error registering web UI handler")
		return fmt.Errorf("registering web UI handler: %w", err)
	}

	// Apply CORS policy
	handler = applyCORSHandler(handler)

	// Create the API server
	apiserver := http.Server{
		Addr:              cfg.Web.APIHost,
		Handler:           handler,
		ReadTimeout:       cfg.Web.ReadTimeout,
		ReadHeaderTimeout: cfg.Web.ReadTimeout,
		WriteTimeout:      cfg.Web.WriteTimeout,
	}

	// Start the service listening for requests in a separate goroutine
	go func() {
		logger.Infof("API listening on %s", apiserver.Addr)
		serverErrors <- apiserver.ListenAndServe()
		logger.Infof("stopping API server")
	}()

	// Waiting for shutdown signal or POSIX signals
	select {
	case err := <-serverErrors:
		// Non-recoverable server error
		return fmt.Errorf("server error: %w", err)

	case sig := <-shutdown:
		logger.Infof("signal %v received, start shutdown", sig)

		// Asking API server to shut down and load shed.
		err := apirouter.Close()
		if err != nil {
			logger.WithError(err).Warning("graceful shutdown of apirouter error")
		}

		// Give outstanding requests a deadline for completion.
		ctx, cancel := context.WithTimeout(context.Background(), cfg.Web.ShutdownTimeout)
		defer cancel()

		// Asking listener to shut down and load shed.
		err = apiserver.Shutdown(ctx)
		if err != nil {
			logger.WithError(err).Warning("error during graceful shutdown of HTTP server")
			err = apiserver.Close()
		}

		// Log the status of this shutdown.
		switch {
		case sig == syscall.SIGSTOP:
			return errors.New("integrity issue caused shutdown")
		case err != nil:
			return fmt.Errorf("could not stop server gracefully: %w", err)
		}
	}

	return nil
}
