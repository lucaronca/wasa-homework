package routes

import (
	"net/http"

	"github.com/gofrs/uuid"
	"github.com/julienschmidt/httprouter"
	"github.com/lucaronca/wasa-homework/service/api/reqcontext"
	"github.com/sirupsen/logrus"
)

// WithReqCtx parses the request and adds a reqcontext.RequestContext instance related to the request.
func NewReqCtxMiddleware(logger *logrus.FieldLogger) func(fn Handler) httprouter.Handle {
	return func(fn Handler) httprouter.Handle {
		return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
			reqUUID, err := uuid.NewV4()
			if err != nil {
				(*logger).WithError(err).Error("can't generate a request UUID")
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			var ctx = reqcontext.RequestContext{
				ReqUUID: reqUUID,
			}

			// Create a request-specific logger
			ctx.Logger = (*logger).WithFields(logrus.Fields{
				"reqid":     ctx.ReqUUID.String(),
				"remote-ip": r.RemoteAddr,
			})

			// Call the next handler in chain
			fn(w, r, ps, ctx)
		}
	}
}
