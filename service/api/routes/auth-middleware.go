package routes

import (
	"errors"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/lucaronca/wasa-homework/service/api/reqcontext"
	"github.com/lucaronca/wasa-homework/service/api/services"
)

// wrap parses the request and adds a reqcontext.RequestContext instance related to the request.
func NewTokenAuthMiddleware(authService services.AuthService) Middleware {
	return func(fn Handler) Handler {
		return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
			// Get the Basic Authentication credentials
			authHeader := strings.Split(r.Header.Get("Authorization"), " ")
			var bearer string
			if len(authHeader) < 2 {
				bearer = ""
			} else {
				bearer = authHeader[1]
			}

			if len(bearer) > 0 {
				if user, err := authService.Authorize(bearer); err != nil {
					if errors.Is(err, services.ErrNoUser) {
						w.Header().Set("WWW-Authenticate", "Bearer")
						http.Error(
							w,
							"Cannot find a user associated with the given access token",
							http.StatusUnauthorized,
						)
					} else {
						http.Error(w, err.Error(), http.StatusInternalServerError)
					}
				} else {
					// Set user info in reqcontext
					ctx.User = reqcontext.User{Id: user.Id, Username: user.Username}
					// Delegate request to the given handle
					fn(w, r, ps, ctx)
				}
			} else {
				// Request Authentication otherwise
				w.Header().Set("WWW-Authenticate", "Bearer")
				http.Error(w, "Access token is missing or invalid", http.StatusUnauthorized)
			}
		}
	}
}
