package tencent_faas_adapter

import (
	"github.com/tencentyun/scf-go-lib/cloudevents/scf"
	"net/http"
	"strings"
)

type APIGatewayProxyResponseBuilder struct {
	statusCode int
	header     http.Header
	body       []byte
}

func NewAPIGatewayProxyResponseBuilder() *APIGatewayProxyResponseBuilder {
	return &APIGatewayProxyResponseBuilder{
		header: make(http.Header),
	}
}

func (b *APIGatewayProxyResponseBuilder) Header() http.Header {
	return b.header
}

func (b *APIGatewayProxyResponseBuilder) Write(bytes []byte) (int, error) {
	b.body = append(b.body, bytes...)
	return len(bytes), nil
}

func (b *APIGatewayProxyResponseBuilder) WriteHeader(statusCode int) {
	b.statusCode = statusCode
}

func (b *APIGatewayProxyResponseBuilder) Build() scf.APIGatewayProxyResponse {
	return scf.APIGatewayProxyResponse{
		StatusCode:      b.statusCode,
		Headers:         b.fromHTTPHeader(b.header),
		Body:            string(b.body),
		IsBase64Encoded: false,
	}
}

func (b *APIGatewayProxyResponseBuilder) fromHTTPHeader(header http.Header) map[string]string {
	headers := make(map[string]string)

	for k, v := range header {
		headers[k] = strings.Join(v, ", ")
	}

	return headers
}
