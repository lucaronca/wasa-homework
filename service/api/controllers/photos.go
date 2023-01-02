package controllers

import (
	"errors"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/lucaronca/wasa-homework/service/api/reqcontext"
	"github.com/lucaronca/wasa-homework/service/api/routes"
	"github.com/lucaronca/wasa-homework/service/api/services"
)

// photosController binds http requests to an api service and writes the service results to the http response
type photosController struct {
	service      services.PhotosService
	errorHandler ErrorHandler
}

// NewPhotosController creates a default api controller
func NewPhotosController(s services.PhotosService) Controller {
	controller := &photosController{
		service:      s,
		errorHandler: errorHandler,
	}

	return controller
}

// Routes returns all the api routes for the BansController
func (c *photosController) Routes() routes.Routes {
	return routes.Routes{
		{
			Name:         "DeletePhoto",
			Method:       http.MethodDelete,
			Path:         "/photos/:photoId",
			AuthRequired: true,
			HandlerFunc:  c.DeletePhoto,
		},
		{
			Name:         "UploadPhoto",
			Method:       http.MethodPost,
			Path:         "/photos",
			AuthRequired: true,
			HandlerFunc:  c.UploadPhoto,
		},
		{
			Name:         "GetPhotos",
			Method:       http.MethodGet,
			Path:         "/users/:userId/photos",
			AuthRequired: true,
			HandlerFunc:  c.GetPhotos,
		},
		{
			Name:         "GetMyStream",
			Method:       http.MethodGet,
			Path:         "/users/:userId/stream",
			AuthRequired: true,
			HandlerFunc:  c.GetMyStream,
		},
	}
}

// GetPhotos - Get user photos
func (c *photosController) GetPhotos(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
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

	query := r.URL.Query()

	offsetParam, err := parseIntParameter(query.Get("offset"), false)
	if err != nil {
		c.errorHandler(w, r, &ParsingError{err}, ctx)
		return
	}
	limit := query.Get("limit")
	if limit == "" {
		limit = "20"
	}
	limitParam, err := parseIntParameter(limit, false)
	if err != nil {
		c.errorHandler(w, r, &ParsingError{err}, ctx)
		return
	}
	result, err := c.service.GetUserPhotos(ctx.User.Id, parsedIdParam, offsetParam, limitParam)
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

// UploadPhoto - Upload a photo
func (c *photosController) UploadPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	defer r.Body.Close()

	newPhoto, err := c.service.CreatePhoto(ctx.User.Id, r.Body)
	if errors.Is(err, services.ErrNoUser) {
		c.errorHandler(w, r, &NotFoundError{"User"}, ctx)
		return
	} else if errors.Is(err, services.ErrPhotoFormatNotSupported) {
		c.errorHandler(w, r, &ParsingError{err}, ctx)
		return
	} else if err != nil {
		c.errorHandler(w, r, err, ctx)
		return
	}

	encodeJSONResponse(newPhoto, http.StatusCreated, w, ctx)
}

// DeletePhoto - Delete a photo
func (c *photosController) DeletePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	targetPhotoIdParam, err := parseIntParameter(ps.ByName("photoId"), true)
	if err != nil {
		c.errorHandler(w, r, &ParsingError{err}, ctx)
		return
	}

	err = c.service.DeletePhoto(ctx.User.Id, targetPhotoIdParam)
	if errors.Is(err, services.ErrNoUser) {
		c.errorHandler(w, r, &NotFoundError{"User"}, ctx)
		return
	} else if errors.Is(err, services.ErrNoPhoto) {
		c.errorHandler(w, r, &NotFoundError{"Photo"}, ctx)
		return
	} else if errors.Is(err, services.ErrUserForbidden) {
		c.errorHandler(w, r, &ForbiddenError{err}, ctx)
		return
	} else if err != nil {
		c.errorHandler(w, r, err, ctx)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (c *photosController) GetMyStream(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	userIdParam := ps.ByName("userId")
	if userIdParam != "me" {
		c.errorHandler(w, r, &ParsingError{errors.New("Invalid user param")}, ctx)
		return
	}

	query := r.URL.Query()

	offsetParam, err := parseIntParameter(query.Get("offset"), false)
	if err != nil {
		c.errorHandler(w, r, &ParsingError{err}, ctx)
		return
	}
	limit := query.Get("limit")
	if limit == "" {
		limit = "20"
	}
	limitParam, err := parseIntParameter(limit, false)
	if err != nil {
		c.errorHandler(w, r, &ParsingError{err}, ctx)
		return
	}
	result, err := c.service.GetStream(ctx.User.Id, offsetParam, limitParam)
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
