package controllers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/lucaronca/wasa-homework/service/api/reqcontext"
	"github.com/lucaronca/wasa-homework/service/api/routes"
	"github.com/lucaronca/wasa-homework/service/api/services"
)

// usersController binds http requests to an api service and writes the service results to the http response
type usersController struct {
	service      services.UsersService
	errorHandler ErrorHandler
}

// NewUsersController creates a default api controller
func NewUsersController(s services.UsersService) Controller {
	controller := &usersController{
		service:      s,
		errorHandler: errorHandler,
	}

	return controller
}

// Routes returns all the api routes for the UsersController
func (c *usersController) Routes() routes.Routes {
	return routes.Routes{
		{
			Name:         "GetUserProfile",
			Method:       http.MethodGet,
			Path:         "/users/:userId",
			AuthRequired: true,
			HandlerFunc:  c.GetUserProfile,
		},
		{
			Name:         "GetUsers",
			Method:       http.MethodGet,
			Path:         "/users",
			AuthRequired: true,
			HandlerFunc:  c.GetUsers,
		},
		{
			Name:         "SetMyUserName",
			Method:       http.MethodPatch,
			Path:         "/users/me",
			AuthRequired: true,
			HandlerFunc:  c.SetMyUserName,
		},
	}
}

func (c *usersController) GetUserProfile(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	userIdParam := ps.ByName("userId")
	var parsedIdParam int
	if userIdParam == "me" {
		parsedIdParam = ctx.User.Id
	} else {
		parsed, err := parseIntParameter(userIdParam, true)
		if err != nil {
			c.errorHandler(w, r, &ParsingError{errors.New("userId should be a valid int number")}, ctx)
			return
		}
		parsedIdParam = parsed
	}

	result, err := c.service.GetUser(ctx.User.Id, parsedIdParam)
	// If an error occurred, encode the error with the status code
	if errors.Is(err, services.ErrNoUser) {
		c.errorHandler(w, r, &NotFoundError{"User"}, ctx)
		return
	} else if err != nil {
		c.errorHandler(w, r, err, ctx)
		return
	}
	// If no error, encode the result and the result code
	encodeJSONResponse(result, http.StatusOK, w, ctx)
}

func (c *usersController) GetUsers(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	query := r.URL.Query()

	username := query.Get("username")

	// if username == "" {
	// 	c.errorHandler(w, r, &ParsingError{errors.New("username is required")}, ctx)
	// 	return
	// }

	result, err := c.service.GetUsers(username)
	// If an error occurred, encode the error with the status code
	if errors.Is(err, services.ErrNoUser) {
		c.errorHandler(w, r, &NotFoundError{"User"}, ctx)
		return
	} else if err != nil {
		c.errorHandler(w, r, err, ctx)
		return
	}
	// If no error, encode the result and the result code
	encodeJSONResponse(result, http.StatusOK, w, ctx)
}

func (c *usersController) SetMyUserName(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	contentType := r.Header.Get("Content-Type")
	if contentType != "application/json-patch+json" {
		c.errorHandler(w, r, &ParsingError{errors.New("Wrong `Content-Type` header")}, ctx)
		return
	}

	setMyUserNameRequestParam := []SetMyUserNameRequest{}
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if err := d.Decode(&setMyUserNameRequestParam); err != nil {
		c.errorHandler(w, r, &ParsingError{errors.New("Payload not valid")}, ctx)
		return
	}

	if len(setMyUserNameRequestParam) != 1 {
		c.errorHandler(w, r, &ParsingError{errors.New("Wrong patch length")}, ctx)
		return
	}

	if err := assertSetMyUserNameRequestValid(setMyUserNameRequestParam[0]); err != nil {
		switch {
		case errors.Is(err, ErrSeMyUserNameOpIsZero):
			c.errorHandler(w, r, &RequiredError{"op"}, ctx)
		case errors.Is(err, ErrSeMyUserNamePathIsZero):
			c.errorHandler(w, r, &RequiredError{"path"}, ctx)
		case errors.Is(err, ErrSetMyUserNameValueIsZero):
			c.errorHandler(w, r, &RequiredError{"value"}, ctx)
		case
			errors.Is(err, ErrSeMyUserNameOpIsNotValid),
			errors.Is(err, ErrSeMyUserNamePathIsNotValid),
			errors.Is(err, ErrSetMyUserNameValueIsNotValid):
			c.errorHandler(w, r, &ParsingError{err}, ctx)
		}
		return
	}

	result, err := c.service.UpdateUsername(ctx.User.Id, setMyUserNameRequestParam[0].Value)
	if err != nil {
		c.errorHandler(w, r, err, ctx)
		return
	}
	encodeJSONResponse(result, http.StatusOK, w, ctx)
}
