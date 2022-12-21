package api

import (
	"github.com/julienschmidt/httprouter"
	"github.com/lucaronca/wasa-homework/service/api/routes"
	"github.com/sirupsen/logrus"
)

// RouterConfig is used to provide dependencies and configuration to the New function.
type RouterConfig struct {
	// Logger where log entries are sent
	Logger logrus.FieldLogger
}

// HandlerConfig is used to provide dependencies and configuration to the Handler function.
type HandlerConfigPhotos struct {
	PhotosDirectory string
	PhotosUrlPath   string
}

type HandlerConfigDependencies struct {
	LivenessChecker     httprouter.Handle
	TokenAuthMiddleware routes.Middleware
}

type HandlerConfig struct {
	Photos HandlerConfigPhotos
	Deps   HandlerConfigDependencies
}
