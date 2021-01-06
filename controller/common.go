package controller

import (
	"io"
	"net/http"
)

func sendMsgHandler(w http.ResponseWriter, r *http.Request, message string, httpStatusCode int) {
	w.WriteHeader(httpStatusCode)
	io.WriteString(w, message)
}
