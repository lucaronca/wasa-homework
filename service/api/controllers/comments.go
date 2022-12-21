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

// CommentsController binds http requests to an api service and writes the service results to the http response
type commentsController struct {
	service      services.CommentsService
	errorHandler ErrorHandler
}

// Newc creates a default api controller
func NewCommentsController(s services.CommentsService) Controller {
	controller := &commentsController{
		service:      s,
		errorHandler: errorHandler,
	}

	return controller
}

// Routes returns all the api routes for the commentsController
func (c *commentsController) Routes() routes.Routes {
	return routes.Routes{
		{
			Name:         "GetPhotoComments",
			Method:       http.MethodGet,
			Path:         "/photos/:photoId/comments",
			AuthRequired: true,
			HandlerFunc:  c.GetPhotoComments,
		},
		{
			Name:         "CommentPhoto",
			Method:       http.MethodPost,
			Path:         "/photos/:photoId/comments",
			AuthRequired: true,
			HandlerFunc:  c.CommentPhoto,
		},
		{
			Name:         "UncommentPhoto",
			Method:       http.MethodDelete,
			Path:         "/photos/:photoId/comments/:commentId",
			AuthRequired: true,
			HandlerFunc:  c.UncommentPhoto,
		},
	}
}

// GetPhotoComments - Get the comments of a photo
func (c *commentsController) GetPhotoComments(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	photoIdParam, err := parseIntParameter(ps.ByName("photoId"), true)
	if err != nil {
		c.errorHandler(w, r, &ParsingError{err}, ctx)
		return
	}

	result, err := c.service.GetPhotoComments(photoIdParam, ctx.User.Id)
	// If an error occurred, encode the error with the status code
	if errors.Is(err, services.ErrNoPhoto) {
		c.errorHandler(w, r, &NotFoundError{"Photo"}, ctx)
		return
	} else if errors.Is(err, services.ErrNoComment) {
		c.errorHandler(w, r, &NotFoundError{"Comment"}, ctx)
		return
	} else if err != nil {
		c.errorHandler(w, r, err, ctx)
		return
	}
	// If no error, encode the body and the result code
	encodeJSONResponse(result, http.StatusOK, w, ctx)
}

// CommentPhoto - Post a comment to a photo
func (c *commentsController) CommentPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	photoIdParam, err := parseIntParameter(ps.ByName("photoId"), true)
	if err != nil {
		c.errorHandler(w, r, &ParsingError{err}, ctx)
		return
	}

	commentParam := CommentPhoto{}
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields()
	if err := d.Decode(&commentParam); err != nil {
		c.errorHandler(w, r, &ParsingError{errors.New("Payload not valid")}, ctx)
		return
	}
	if err := assertCommentPhotoValid(commentParam); err != nil {
		if errors.Is(err, ErrLoginNameIsZero) {
			c.errorHandler(w, r, &RequiredError{"content"}, ctx)
		} else if errors.Is(err, ErrLoginNameIsNotValid) {
			c.errorHandler(w, r, &ParsingError{err}, ctx)
		}
	}

	newComment, err := c.service.CommentPhoto(photoIdParam, ctx.User.Id, commentParam.Content)
	// If an error occurred, encode the error with the status code
	if errors.Is(err, services.ErrNoPhoto) {
		c.errorHandler(w, r, &NotFoundError{"Photo"}, ctx)
		return
	} else if errors.Is(err, services.ErrNoComment) {
		c.errorHandler(w, r, &NotFoundError{"Comment"}, ctx)
		return
	} else if err != nil {
		c.errorHandler(w, r, err, ctx)
		return
	}

	encodeJSONResponse(newComment, http.StatusCreated, w, ctx)
}

// UncommentPhoto - Remove a comment to a photo
func (c *commentsController) UncommentPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	photoIdParam, err := parseIntParameter(ps.ByName("photoId"), true)
	if err != nil {
		c.errorHandler(w, r, &ParsingError{err}, ctx)
		return
	}

	commentIdParam, err := parseIntParameter(ps.ByName("commentId"), true)
	if err != nil {
		c.errorHandler(w, r, &ParsingError{err}, ctx)
		return
	}

	err = c.service.UncommentPhoto(photoIdParam, commentIdParam, ctx.User.Id)
	// If an error occurred, encode the error with the status code
	if errors.Is(err, services.ErrNoPhoto) {
		c.errorHandler(w, r, &NotFoundError{"Photo"}, ctx)
		return
	} else if errors.Is(err, services.ErrNoComment) {
		c.errorHandler(w, r, &NotFoundError{"Comment"}, ctx)
		return
	} else if errors.Is(err, services.ErrDeleteNotAllowed) {
		c.errorHandler(w, r, &ForbiddenError{err}, ctx)
		return
	} else if err != nil {
		c.errorHandler(w, r, err, ctx)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
