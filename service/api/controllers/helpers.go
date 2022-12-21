package controllers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"

	"github.com/lucaronca/wasa-homework/service/api/reqcontext"
)

func encodeTextResponse(resp string, status int, w http.ResponseWriter, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	w.WriteHeader(status)

	_, err := io.WriteString(w, resp)
	if err != nil {
		ctx.Logger.WithError(err).Error("Error writing string response")
	}
}

// EncodeJSONResponse uses the json encoder to write an interface to the http response with an optional status code
func encodeJSONResponse(i interface{}, status int, w http.ResponseWriter, ctx reqcontext.RequestContext) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)

	err := json.NewEncoder(w).Encode(i)
	if err != nil {
		ctx.Logger.WithError(err).Error("Error writing JSON response")
	}
}

const errMsgRequiredMissing = "required parameter is missing"

// parseIntParameter parses a string parameter to an int.
func parseIntParameter(param string, required bool) (int, error) {
	if param == "" {
		if required {
			return 0, errors.New(errMsgRequiredMissing)
		}

		return 0, nil
	}

	val, err := strconv.ParseInt(param, 10, 32)
	if err != nil {
		return 0, err
	}

	return int(val), nil
}
