package routes

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/lucaronca/wasa-homework/service/api/reqcontext"
)

// Handler is the signature for functions that accepts a reqcontext.RequestContext in addition to those
// required by the httprouter package.
type Handler func(http.ResponseWriter, *http.Request, httprouter.Params, reqcontext.RequestContext)

type Middleware func(Handler) Handler

// A Route defines the parameters for an api endpoint
type Route struct {
	Name         string
	Method       string
	Path         string
	AuthRequired bool
	HandlerFunc  Handler
}

// Routes are a collection of defined api endpoints
type Routes []Route
