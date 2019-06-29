package httprouter

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/aws/aws-lambda-go/events"
)

func LambdaRedirect(ctx context.Context, req events.APIGatewayProxyRequest, newUrl string, code int) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: code,
		Headers: map[string]string{
			"Location": newUrl,
		},
	}, nil
}

func LambdaAllow(ctx context.Context, req events.APIGatewayProxyRequest, allow string) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers: map[string]string{
			"Allow": allow,
		},
	}, nil
}

func LambdaNotAllowed(ctx context.Context, req events.APIGatewayProxyRequest, allow string) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: 405,
		Headers: map[string]string{
			"Allow": allow,
		},
		Body: `{"error": "Method Not Allowed"}`,
	}, nil
}

func LambdaNotFound(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: 404,
		Body:       `{"error": "Not Found"}`,
	}, nil
}

func RequestToLambda(req *http.Request) (events.APIGatewayProxyRequest, error) {
	e := events.APIGatewayProxyRequest{
		HTTPMethod:            req.Method,
		Path:                  strings.Split(req.URL.RequestURI(), "?")[0],
		Headers:               map[string]string{},
		QueryStringParameters: map[string]string{},
		PathParameters:        map[string]string{},
		StageVariables:        map[string]string{},
		// Resource:              params.Path,
	}
	// e.RequestContext.RequestID = utils.UUID()
	// e.RequestContext.ResourcePath = params.Path
	e.RequestContext.HTTPMethod = req.Method
	for i := range req.URL.Query() {

		e.QueryStringParameters[i] = req.URL.Query().Get(i)
	}
	for i := range req.Header {
		if strings.HasPrefix(i, "Stagevariable_") {
			e.StageVariables[strings.ReplaceAll(i, "Stagevariable_", "")] = req.Header.Get(i)
			continue
		}
		e.Headers[i] = req.Header.Get(i)
	}
	// e.Headers["X-Forwarded-For"] = GetForwarded(req)
	b, _ := ioutil.ReadAll(req.Body)
	e.Body = fmt.Sprintf("%s", b)
	return e, nil
}

func ResToHttp(w http.ResponseWriter, req *http.Request, res events.APIGatewayProxyResponse) {
	w.WriteHeader(res.StatusCode)
	for key := range res.Headers {
		w.Header().Set(key, res.Headers[key])
	}
	w.Write([]byte(res.Body))
}

func HttpAddParams(event events.APIGatewayProxyRequest, params Params) events.APIGatewayProxyRequest {
	for i := range params {
		if params[i].Key == "__stage__" {
			event.RequestContext.Stage = params[i].Value
			continue
		}
		event.PathParameters[params[i].Key] = params[i].Value
	}
	return event
}
