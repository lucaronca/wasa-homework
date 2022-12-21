/*
Package api exposes the main API engine.

To use this package, you should create a new instance with New() passing a valid Config. The resulting Router will have
the Router.Handler() function that returns a handler that can be used in a http.Server (or in other middlewares).

Example:

	// Create the API router
	apirouter, err := api.New(api.Config{
		Logger:   logger
	})
	if err != nil {
		logger.WithError(err).Error("error creating the API server instance")
		return fmt.Errorf("error creating the API server instance: %w", err)
	}
	handler := apirouter.Handler(
		handlerCfg,
		controllers...,
	)

	// Create the API server
	apiserver := http.Server{
		Addr:              cfg.Web.APIHost,
		Handler:           handler,
		ReadTimeout:       cfg.Web.ReadTimeout,
		ReadHeaderTimeout: cfg.Web.ReadTimeout,
		WriteTimeout:      cfg.Web.WriteTimeout,
	}

	// Start the service listening for requests in a separate goroutine
	apiserver.ListenAndServe()

See the `main.go` file inside the `cmd/webapi` for a full usage example.
*/
package api

import (
	"errors"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/lucaronca/wasa-homework/service/api/controllers"
	"github.com/sirupsen/logrus"
)

// Router is the package API interface representing an API handler builder
type Router interface {
	// Handler returns an HTTP handler for APIs provided in this package
	Handler(
		cfg HandlerConfig,
		controllers ...controllers.Controller,
	) http.Handler

	// Close terminates any resource used in the package
	Close() error
}

// New returns a new Router instance
func New(cfg RouterConfig) (Router, error) {
	// Check if the configuration is correct
	if cfg.Logger == nil {
		return nil, errors.New("logger is required")
	}

	// Create a new router where we will register HTTP endpoints. The server will pass requests to this router to be
	// handled.
	router := httprouter.New()
	router.RedirectTrailingSlash = false
	router.RedirectFixedPath = false

	return &_router{
		router:     router,
		baseLogger: cfg.Logger,
	}, nil
}

type _router struct {
	router *httprouter.Router

	// baseLogger is a logger for non-requests contexts, like goroutines or background tasks not started by a request.
	// Use context logger if available (e.g., in requests) instead of this logger.
	baseLogger logrus.FieldLogger
}
