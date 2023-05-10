package http_app

import (
	"io"
	"net/http"
)

func (f *HTTPApp) pingPongHandler(writer http.ResponseWriter, request *http.Request) {
	io.WriteString(writer, "pong")
}
