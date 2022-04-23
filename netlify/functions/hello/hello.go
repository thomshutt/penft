package main

import (
	"bytes"
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"io"
	"net/http"
	"os"
	"strings"
)

var TATUM_API_KEY = os.Getenv("TATUM_API_KEY")

func handler(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	responseBody := "Hello, World! "

	getwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	os.Setenv("TMPDIR", getwd)

	httpReq, err := http.NewRequestWithContext(context.Background(), request.HTTPMethod, request.Path, strings.NewReader(request.Body))
	if err != nil {
		return nil, err
	}

	// some headers may be important, let get all of them, just in case
	for name, value := range request.Headers {
		httpReq.Header.Add(name, value)
	}

	file, _, err := httpReq.FormFile("file")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var buf bytes.Buffer
	io.Copy(&buf, file)

	responseBody += fmt.Sprintf("Buffer Length: %d", buf.Len())

	return &events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers:    map[string]string{"Content-Type": "text/plain"},
		//MultiValueHeaders: http.Header{"Set-Cookie": {"Ding", "Ping"}},
		Body:            responseBody,
		IsBase64Encoded: false,
	}, nil
}

func main() {
	// Make the handler available for Remote Procedure Call
	lambda.Start(handler)
}
