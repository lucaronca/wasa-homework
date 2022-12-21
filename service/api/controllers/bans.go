package controllers

import (
	"errors"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/lucaronca/wasa-homework/service/api/reqcontext"
	"github.com/lucaronca/wasa-homework/service/api/routes"
	"github.com/lucaronca/wasa-homework/service/api/services"
)

// bansController binds http requests to an api service and writes the service results to the http response
type bansController struct {
	service      services.BansService
	errorHandler ErrorHandler
}

// NewBansController creates a default api controller
func NewBansController(s services.BansService) Controller {
	controller := &bansController{
		service:      s,
		errorHandler: errorHandler,
	}

	return controller
}

// Routes returns all the api routes for the BansController
func (c *bansController) Routes() routes.Routes {
	return routes.Routes{
		{
			Name:         "BanUser",
			Method:       http.MethodPut,
			Path:         "/users/me/bans/:targetUserId",
			AuthRequired: true,
			HandlerFunc:  c.BanUser,
		},
		{
			Name:         "UnbanUser",
			Method:       http.MethodDelete,
			Path:         "/users/me/bans/:targetUserId",
			AuthRequired: true,
			HandlerFunc:  c.UnbanUser,
		},
	}
}

// BanUser - Ban a user
func (c *bansController) BanUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	targetUserIdParam, err := parseIntParameter(ps.ByName("targetUserId"), true)
	if err != nil {
		c.errorHandler(w, r, &ParsingError{err}, ctx)
		return
	}

	err = c.service.BanUser(ctx.User.Id, targetUserIdParam)
	if errors.Is(err, services.ErrNoUser) {
		c.errorHandler(w, r, &NotFoundError{"User"}, ctx)
		return
	} else if err != nil {
		c.errorHandler(w, r, err, ctx)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// UnbanUser - Unban a user
func (c *bansController) UnbanUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	targetUserIdParam, err := parseIntParameter(ps.ByName("targetUserId"), true)
	if err != nil {
		c.errorHandler(w, r, &ParsingError{err}, ctx)
		return
	}

	err = c.service.UnbanUser(ctx.User.Id, targetUserIdParam)
	if errors.Is(err, services.ErrNoUser) {
		c.errorHandler(w, r, &NotFoundError{"User"}, ctx)
		return
	} else if err != nil {
		c.errorHandler(w, r, err, ctx)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
