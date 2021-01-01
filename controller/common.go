package controller

import (
	"io"
	"net/http"
)

func errorHandler(w http.ResponseWriter, r *http.Request, message string, httpStatusCode int) {
	w.WriteHeader(httpStatusCode)
	io.WriteString(w, message)
}
