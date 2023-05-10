package tencent_faas_adapter

import (
	"context"
	"github.com/tencentyun/scf-go-lib/cloudevents/scf"
	"github.com/yrosukedev/chat_record_sync/logger"
	"net/http"
	"net/url"
	"strings"
)

type HTTPAdapter struct {
	handler http.Handler
	logger  logger.Logger
}

func NewHTTPAdapter(handler http.Handler, logger logger.Logger) *HTTPAdapter {
	return &HTTPAdapter{
		handler: handler,
		logger:  logger,
	}
}

func (h *HTTPAdapter) HTTPHandler(ctx context.Context, req scf.APIGatewayProxyRequest) (resp scf.APIGatewayProxyResponse, err error) {
	h.logger.Info(ctx, "[tencent faas http adapter] fass handler begin, request: %#v", req)
	defer h.logger.Info(ctx, "[tencent faas http adapter] fass handler end, resq: %#v, err: %v", resp, err)

	httpReq, err := h.httpRequestFrom(req)
	if err != nil {
		h.logger.Error(ctx, "[tencent faas http adapter] fails to transform faas request to http request, faas request: %#v, err: %v", req, err)
		return scf.APIGatewayProxyResponse{}, err
	}

	responseWriter := NewAPIGatewayProxyResponseBuilder()

	h.handler.ServeHTTP(responseWriter, httpReq)

	return responseWriter.Build(), nil
}

func (h *HTTPAdapter) httpRequestFrom(req scf.APIGatewayProxyRequest) (*http.Request, error) {
	httpReq, err := http.NewRequest(req.HTTPMethod, h.urlFrom(req), strings.NewReader(req.Body))
	if err != nil {
		return nil, err
	}
	httpReq.Header = h.httpHeaderFrom(req)
	return httpReq, nil
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
