package controllers

import (
	"errors"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/lucaronca/wasa-homework/service/api/reqcontext"
	"github.com/lucaronca/wasa-homework/service/api/routes"
	"github.com/lucaronca/wasa-homework/service/api/services"
)

// likesController binds http requests to an api service and writes the service results to the http response
type likesController struct {
	service      services.LikesService
	errorHandler ErrorHandler
}

// NewLikesController creates a default api controller
func NewLikesController(s services.LikesService) Controller {
	controller := &likesController{
		service:      s,
		errorHandler: errorHandler,
	}

	return controller
}

// Routes returns all the api routes for the likesController
func (c *likesController) Routes() routes.Routes {
	return routes.Routes{
		{
			Name:         "GetPhotoLikes",
			Method:       http.MethodGet,
			Path:         "/photos/:photoId/likes",
			AuthRequired: true,
			HandlerFunc:  c.GetPhotoLikes,
		},
		{
			Name:         "LikePhoto",
			Method:       http.MethodPut,
			Path:         "/photos/:photoId/likes/me",
			AuthRequired: true,
			HandlerFunc:  c.LikePhoto,
		},
		{
			Name:         "UnlikePhoto",
			Method:       http.MethodDelete,
			Path:         "/photos/:photoId/likes/me",
			AuthRequired: true,
			HandlerFunc:  c.UnlikePhoto,
		},
	}
}

// GetPhotoLikes - Get the likes of a photo
func (c *likesController) GetPhotoLikes(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	photoIdParam, err := parseIntParameter(ps.ByName("photoId"), true)
	if err != nil {
		c.errorHandler(w, r, &ParsingError{err}, ctx)
		return
	}

	result, err := c.service.GetPhotoLikes(photoIdParam, ctx.User.Id)
	// If an error occurred, encode the error with the status code
	if errors.Is(err, services.ErrNoPhoto) {
		c.errorHandler(w, r, &NotFoundError{"Photo"}, ctx)
		return
	} else if errors.Is(err, services.ErrNoLike) {
		c.errorHandler(w, r, &NotFoundError{"Like"}, ctx)
		return
	} else if err != nil {
		c.errorHandler(w, r, err, ctx)
		return
	}
	// If no error, encode the body and the result code
	encodeJSONResponse(result, http.StatusOK, w, ctx)
}

// LikePhoto - Put a like to a photo
func (c *likesController) LikePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	photoIdParam, err := parseIntParameter(ps.ByName("photoId"), true)
	if err != nil {
		c.errorHandler(w, r, &ParsingError{err}, ctx)
		return
	}

	err = c.service.LikePhoto(photoIdParam, ctx.User.Id)
	// If an error occurred, encode the error with the status code
	if errors.Is(err, services.ErrNoPhoto) {
		c.errorHandler(w, r, &NotFoundError{"Photo"}, ctx)
		return
	} else if errors.Is(err, services.ErrNoLike) {
		c.errorHandler(w, r, &NotFoundError{"Like"}, ctx)
		return
	} else if err != nil {
		c.errorHandler(w, r, err, ctx)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// UnlikePhoto - Remove a like to a photo
func (c *likesController) UnlikePhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	photoIdParam, err := parseIntParameter(ps.ByName("photoId"), true)
	if err != nil {
		c.errorHandler(w, r, &ParsingError{err}, ctx)
		return
	}

	err = c.service.UnlikePhoto(photoIdParam, ctx.User.Id)
	// If an error occurred, encode the error with the status code
	if errors.Is(err, services.ErrNoPhoto) {
		c.errorHandler(w, r, &NotFoundError{"Photo"}, ctx)
		return
	} else if errors.Is(err, services.ErrNoLike) {
		c.errorHandler(w, r, &NotFoundError{"Like"}, ctx)
		return
	} else if err != nil {
		c.errorHandler(w, r, err, ctx)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
