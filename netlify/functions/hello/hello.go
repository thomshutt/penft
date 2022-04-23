package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/grokify/go-awslambda"
	"io"
	"os"
)

var TATUM_API_KEY = os.Getenv("TATUM_API_KEY")

func handler(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	responseBody := "Hello, World! "

	getwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	os.Setenv("TMPDIR", getwd)

	r, err := awslambda.NewReaderMultipart(req)
	if err != nil {
		return nil, err
	}

	var content []byte
	for part, err := r.NextPart(); part != nil && err == nil; {
		partBytes, err := io.ReadAll(part)
		if err != nil {
			return nil, err
		}
		content = append(content, partBytes...)
	}
	if err != nil && err != io.EOF {
		return nil, err
	}

	responseBody += fmt.Sprintf("Buffer Length: %d", len(content))

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
