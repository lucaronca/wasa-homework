package api

import (
	"net/http"

	"github.com/lucaronca/wasa-homework/service/api/controllers"
	"github.com/lucaronca/wasa-homework/service/api/routes"
)

/*
Handler returns an instance of http.Handler that handle APIs registered by the controllers.
You need to pass:
- An `HandlerConfig` struct used to configure where photos assets will be saved,
to specify the path which will be used to serve the photos, and the needed Handler dependecies
- 1..n controllers to handle the API endpoints.
*/
func (rt *_router) Handler(
	cfg HandlerConfig,
	controllers ...controllers.Controller,
) http.Handler {
	tokenAuthMiddleware := cfg.Deps.TokenAuthMiddleware
	reqCtxMiddleware := routes.NewReqCtxMiddleware(&rt.baseLogger)

	// Register routes defined in the controllers
	for _, controller := range controllers {
		for _, route := range controller.Routes() {
			handler := route.HandlerFunc
			// If AuthRequired, wrap the handler with token authentication middleware
			if route.AuthRequired {
				handler = tokenAuthMiddleware(route.HandlerFunc)
			}
			rt.router.Handle(route.Method, route.Path, reqCtxMiddleware(handler))
		}
	}

	// Special routes
	liveness := cfg.Deps.LivenessChecker
	rt.router.GET("/liveness", liveness)

	// Static files
	rt.router.ServeFiles(cfg.Photos.PhotosUrlPath+"/*filepath", http.Dir(cfg.Photos.PhotosDirectory))

	return rt.router
}
