package main

import (
	"context"
	"log"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/lambdacontext"
)

// Response is of type APIGatewayProxyResponse since we're leveraging the
// AWS Lambda Proxy Request functionality (default behavior)
//
// https://serverless.com/framework/docs/providers/aws/events/apigateway/#lambda-proxy-integration
type Response events.APIGatewayProxyResponse

// Request is of type APIGatewayProxyRequest
type Request events.APIGatewayProxyRequest

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(ctx context.Context, request Request) (Response, error) {
	responseHeader := map[string]string{"Content-Type": "application/json", "X-BaseAPI-V1-Reply": "health-handler"}

	// not actually using this here - just setting up a patten as to what will be implemented
	lc, ok := lambdacontext.FromContext(ctx)
	if !ok {
		log.Printf("Could not process lc: %v", lc)
		return Response{StatusCode: http.StatusNotFound, Headers: responseHeader}, nil
	}

	log.Printf("Request received  AWS.RequestID: %v", lc.AwsRequestID)
	return Response{StatusCode: http.StatusNoContent, Headers: responseHeader}, nil
}

func main() {
	lambda.Start(Handler)
}
