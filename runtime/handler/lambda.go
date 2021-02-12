package handler

import (
	"context"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/foldsh/fold/logging"
)

type HTTPRequestDoer interface {
	DoRequest(http.ResponseWriter, *http.Request)
}

func NewLambda(logger logging.Logger, doer HTTPRequestDoer) *LambdaHandler {
	return &LambdaHandler{logger, doer}
}

type LambdaHandler struct {
	logger logging.Logger
	doer   HTTPRequestDoer
}

func (lh *LambdaHandler) Handle(
	ctx context.Context,
	e events.APIGatewayProxyRequest,
) (events.APIGatewayProxyResponse, error) {
	url, err := url.Parse(e.Path)
	if err != nil {
		// Pretty sure this can never happen but...
		lh.logger.Errorf("failed to parse path: %s", e.Path)
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       `{"title":"Failed to parse path"}`,
		}, err
	}
	req, err := http.NewRequest(
		e.HTTPMethod,
		url.String(),
		ioutil.NopCloser(strings.NewReader(e.Body)),
	)
	if err != nil {
		lh.logger.Errorf("failed to translate gateway request")
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       `{"title":"Failed to translate AWS APIGateway Request"}`,
		}, err
	}
	req.Header = e.MultiValueHeaders
	req.ContentLength = int64(len(e.Body))
	req.Close = false
	req.Host = e.Headers["Host"]
	res := newResponseWriter()
	lh.doer.DoRequest(res, req)
	return res.toAPIGatewayResponse(), nil
}

func (lh *LambdaHandler) Serve() {
	lambda.Start(lh.Handle)
}

type responseWriter struct {
	statusCode int
	headers    http.Header
	body       []byte
}

func newResponseWriter() *responseWriter {
	return &responseWriter{headers: make(map[string][]string)}
}

func (rw *responseWriter) toAPIGatewayResponse() events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		StatusCode:        rw.statusCode,
		MultiValueHeaders: rw.headers,
		Body:              string(rw.body),
	}
}

func (rw *responseWriter) Header() http.Header {
	return rw.headers
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	// TODO this isn't perfect really, and is coupled with what the router is doing.
	rw.body = b
	return len(b), nil
}

func (rw *responseWriter) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
}
