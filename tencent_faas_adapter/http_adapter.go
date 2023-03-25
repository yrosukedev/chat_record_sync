package tencent_faas_adapter

import (
	"context"
	"github.com/tencentyun/scf-go-lib/cloudevents/scf"
	"net/http"
	"net/url"
	"strings"
)

type HTTPAdapter struct {
	handler http.Handler
}

func NewHTTPAdapter(handler http.Handler) *HTTPAdapter {
	return &HTTPAdapter{
		handler: handler,
	}
}

func (h *HTTPAdapter) HTTPHandler(ctx context.Context, req scf.APIGatewayProxyRequest) (resp scf.APIGatewayProxyResponse, err error) {
	httpReq, err := http.NewRequest(req.HTTPMethod, h.urlFrom(req), strings.NewReader(req.Body))
	if err != nil {
		return scf.APIGatewayProxyResponse{}, err
	}
	httpReq.Header = h.httpHeaderFrom(req)

	responseWriter := NewAPIGatewayProxyResponseBuilder()

	h.handler.ServeHTTP(responseWriter, httpReq)

	return responseWriter.Build(), nil
}

func (h *HTTPAdapter) httpHeaderFrom(req scf.APIGatewayProxyRequest) http.Header {
	httpHeader := make(http.Header)

	for k, v := range req.Headers {
		values := strings.Split(v, ",")

		for _, v := range values {
			httpHeader.Add(k, strings.TrimSpace(v))
		}
	}

	return httpHeader
}

func (h *HTTPAdapter) urlFrom(req scf.APIGatewayProxyRequest) string {
	host, ok := req.Headers["host"]
	if !ok {
		host = "localhost"
	}

	theURL := url.URL{
		Scheme: "https",
		Host:   host,
		Path:   req.Path,
	}

	query := theURL.Query()
	for k, v := range req.QueryString {
		query.Add(k, v)
	}
	theURL.RawQuery = query.Encode()

	return theURL.String()
}
