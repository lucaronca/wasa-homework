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

// LoginController binds http requests to an api service and writes the service results to the http response
type loginController struct {
	service      services.AuthService
	errorHandler ErrorHandler
}

// NewLoginController creates a default api controller
func NewLoginController(s services.AuthService) Controller {
	controller := &loginController{
		service:      s,
		errorHandler: errorHandler,
	}

	return controller
}

// Routes returns all the api routes for the LoginController
func (c *loginController) Routes() routes.Routes {
	return routes.Routes{
		{
			Name:         "DoLogin",
			Method:       http.MethodPost,
			Path:         "/session",
			AuthRequired: false,
			HandlerFunc:  c.DoLogin,
		},
	}
}

// DoLogin - Logs in the user
func (c *loginController) DoLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	doLoginRequestParam := DoLoginRequest{}
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if err := d.Decode(&doLoginRequestParam); err != nil {
		c.errorHandler(w, r, &ParsingError{errors.New("Payload not valid")}, ctx)
		return
	}

	if err := assertDoLoginRequestValid(doLoginRequestParam); err != nil {
		if errors.Is(err, ErrLoginNameIsZero) {
			c.errorHandler(w, r, &RequiredError{"name"}, ctx)
		} else if errors.Is(err, ErrLoginNameIsNotValid) {
			c.errorHandler(w, r, &ParsingError{err}, ctx)
		}
		return
	}

	token, isNewUser, err := c.service.DoLogin(doLoginRequestParam.Name)
	// If an error occurred, encode the error with the status code
	if err != nil {
		c.errorHandler(w, r, err, ctx)
		return
	}
	// If no error, and new user return 201 else 200
	if isNewUser {
		encodeJSONResponse(DoLoginResponse{Identifier: token}, http.StatusCreated, w, ctx)
	} else {
		encodeJSONResponse(DoLoginResponse{Identifier: token}, http.StatusOK, w, ctx)
	}
}
