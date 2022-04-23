package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"io/ioutil"
	"net/http"
)

const TATUM_API_KEY = "3e69a90d-1bd7-4994-ae9f-0aefe391b452"
const imageFile = "public/images/image.jpg"

func handler(request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	files, err := ioutil.ReadDir("..")
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		fmt.Println("File: ", file.Name(), file.IsDir())
	}

	_, err = ioutil.ReadFile(imageFile)
	if err != nil {
		return nil, err
	}

	responseBody := "Hello, World!"

	return &events.APIGatewayProxyResponse{
		StatusCode:        200,
		Headers:           map[string]string{"Content-Type": "text/plain"},
		MultiValueHeaders: http.Header{"Set-Cookie": {"Ding", "Ping"}},
		Body:              responseBody,
		IsBase64Encoded:   false,
	}, nil
}

func main() {
	// Make the handler available for Remote Procedure Call
	lambda.Start(handler)
}
