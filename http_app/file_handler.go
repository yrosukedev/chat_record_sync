package http_app

import "net/http"

func (f *HTTPApp) fileHandler(writer http.ResponseWriter, request *http.Request) {
	http.ServeFile(writer, request, request.URL.Path)
}
