package main

import (
	"bytes"
	"context"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"io"
	"math"
	"mime/multipart"
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
		return nil, fmt.Errorf("Error %d: %s", 0, err)
	}

	// some headers may be important, let get all of them, just in case
	for name, value := range request.Headers {
		httpReq.Header.Add(name, value)
	}

	//file, _, err := httpReq.FormFile("file")
	httpReq.ParseMultipartForm(math.MaxInt)
	if err != nil {
		return nil, fmt.Errorf("Error %d: %s", 1, err)
	}

	var file multipart.File
	if httpReq.MultipartForm != nil && httpReq.MultipartForm.File != nil {
		if fhs := httpReq.MultipartForm.File["file"]; len(fhs) > 0 {
			file, err = fhs[0].Open()
			if err != nil {
				return nil, fmt.Errorf("Error %d: %s", 2, err)
			}
			defer file.Close()
		}
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, file)
	if err != nil {
		return nil, fmt.Errorf("Error %d: %s", 3, err)
	}

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
