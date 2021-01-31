package handler

import (
	"context"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/aws/aws-lambda-go/events"

	"github.com/foldsh/fold/internal/testutils"
	"github.com/foldsh/fold/logging"
)

// The purpose of the lambda handler is to take api gateway events and turn them
// into http requests which can be handled by the regular routing stack.
func TestLambdaHandler(t *testing.T) {
	logger := logging.NewTestLogger()
	doer := &mockHTTPDoer{t}
	lambda := &lambdaHandler{logger, doer}

	req := events.APIGatewayProxyRequest{
		HTTPMethod:        "DELETE",
		Path:              "/foo/bar/baz",
		Body:              `{"statusCode":1234,"headers":{"Test-Header":["foo","bar","baz"]}}`,
		MultiValueHeaders: map[string][]string{"Content-Type": []string{"application/json"}},
	}

	res, err := lambda.Handle(context.Background(), req)
	if err != nil {
		t.Fatalf("failed to process API Gateway request: %v", err)
	}
	expectation := events.APIGatewayProxyResponse{
		StatusCode:        1234,
		MultiValueHeaders: map[string][]string{"Test-Header": []string{"foo", "bar", "baz"}},
		Body:              `{"method":"DELETE","path":"/foo/bar/baz"}`,
	}

	testutils.Diff(t, expectation, res, "Body did not match expectation")
}

type mockHTTPDoer struct {
	t *testing.T
}

func (m *mockHTTPDoer) DoRequest(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		m.t.Fatalf("failed to read body")
	}
	json := testutils.UnmarshalJSON(m.t, body)
	w.WriteHeader(int(json["statusCode"].(float64)))
	headers := w.Header()
	for key, values := range castHeaders(json["headers"].(map[string]interface{})) {
		for _, value := range values {
			headers.Add(key, value)
		}
	}
	responseBody := map[string]interface{}{"method": r.Method, "path": r.URL.String()}
	w.Write(testutils.MarshalJSON(m.t, responseBody))
}

func castHeaders(raw map[string]interface{}) http.Header {
	results := make(map[string][]string)
	for key, value := range raw {
		results[key] = []string{}
		for _, i := range value.([]interface{}) {
			results[key] = append(results[key], i.(string))
		}
	}
	return results
}
