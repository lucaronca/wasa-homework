package controllers

import (
	"errors"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/lucaronca/wasa-homework/service/api/reqcontext"
	"github.com/lucaronca/wasa-homework/service/api/routes"
	"github.com/lucaronca/wasa-homework/service/api/services"
)

// followsController binds http requests to an api service and writes the service results to the http response
type followsController struct {
	service      services.FollowsService
	errorHandler ErrorHandler
}

// NewFollowsController creates a default api controller
func NewFollowsController(s services.FollowsService) Controller {
	controller := &followsController{
		service:      s,
		errorHandler: errorHandler,
	}

	return controller
}

// Routes returns all the api routes for the followsController
func (c *followsController) Routes() routes.Routes {
	return routes.Routes{
		{
			Name:         "FollowUser",
			Method:       http.MethodPut,
			Path:         "/users/me/followings/:targetUserId",
			AuthRequired: true,
			HandlerFunc:  c.FollowUser,
		},
		{
			Name:         "GetUserFollowers",
			Method:       http.MethodGet,
			Path:         "/users/:userId/followers",
			AuthRequired: true,
			HandlerFunc:  c.GetUserFollowers,
		},
		{
			Name:         "GetUserFollowings",
			Method:       http.MethodGet,
			Path:         "/users/:userId/followings",
			AuthRequired: true,
			HandlerFunc:  c.GetUserFollowings,
		},
		{
			Name:         "UnfollowUser",
			Method:       http.MethodDelete,
			Path:         "/users/me/followings/:targetUserId",
			AuthRequired: true,
			HandlerFunc:  c.UnfollowUser,
		},
	}
}

// FollowUser - Follow a user
func (c *followsController) FollowUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	targetUserIdParam, err := parseIntParameter(ps.ByName("targetUserId"), true)
	if err != nil {
		c.errorHandler(w, r, &ParsingError{errors.New("targetUserId should be a valid int number")}, ctx)
		return
	}
	if targetUserIdParam == ctx.User.Id {
		c.errorHandler(w, r, &ParsingError{errors.New("You can't follow yourself")}, ctx)
		return
	}

	err = c.service.FollowUser(ctx.User.Id, targetUserIdParam)
	if errors.Is(err, services.ErrNoUser) {
		c.errorHandler(w, r, &NotFoundError{"User"}, ctx)
		return
	} else if err != nil {
		c.errorHandler(w, r, err, ctx)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetUserFollowers - Get user followers
func (c *followsController) GetUserFollowers(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	userIdParam := ps.ByName("userId")
	var parsedIdParam int
	if userIdParam == "me" {
		parsedIdParam = ctx.User.Id
	} else {
		parsed, err := parseIntParameter(userIdParam, true)
		if err != nil {
			c.errorHandler(w, r, &ParsingError{errors.New("targetUserId should be a valid int number")}, ctx)
			return
		}
		parsedIdParam = parsed
	}

	result, err := c.service.GetUserFollowers(ctx.User.Id, parsedIdParam)
	// If an error occurred, encode the error with the status code
	if errors.Is(err, services.ErrNoUser) {
		c.errorHandler(w, r, &NotFoundError{"User"}, ctx)
		return
	} else if err != nil {
		c.errorHandler(w, r, err, ctx)
		return
	}
	// If no error, encode the body and the result code
	encodeJSONResponse(result, http.StatusOK, w, ctx)
}

// GetUserFollowings - Get user followings
func (c *followsController) GetUserFollowings(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	userIdParam := ps.ByName("userId")
	var parsedIdParam int
	if userIdParam == "me" {
		parsedIdParam = ctx.User.Id
	} else {
		parsed, err := parseIntParameter(userIdParam, true)
		if err != nil {
			c.errorHandler(w, r, &ParsingError{errors.New("targetUserId should be a valid int number")}, ctx)
			return
		}
		parsedIdParam = parsed
	}

	result, err := c.service.GetUserFollowings(ctx.User.Id, parsedIdParam)
	// If an error occurred, encode the error with the status code
	if errors.Is(err, services.ErrNoUser) {
		c.errorHandler(w, r, &NotFoundError{"User"}, ctx)
		return
	} else if err != nil {
		c.errorHandler(w, r, err, ctx)
		return
	}
	// If no error, encode the body and the result code
	encodeJSONResponse(result, http.StatusOK, w, ctx)
}

// UnfollowUser - Unfollow a user
func (c *followsController) UnfollowUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	targetUserIdParam, err := parseIntParameter(ps.ByName("targetUserId"), true)
	if err != nil {
		c.errorHandler(w, r, &ParsingError{errors.New("targetUserId should be a valid int number")}, ctx)
		return
	}

	err = c.service.UnfollowUser(ctx.User.Id, targetUserIdParam)
	if errors.Is(err, services.ErrNoUser) {
		c.errorHandler(w, r, &NotFoundError{"User"}, ctx)
		return
	} else if err != nil {
		c.errorHandler(w, r, err, ctx)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
