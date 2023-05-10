package http_app

import (
	"fmt"
	"net/http"
)

func (f *HTTPApp) fileHandler(writer http.ResponseWriter, request *http.Request) {
	http.ServeFile(writer, request, fmt.Sprintf("./%v", request.URL.Path))
}
