package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	ecr "github.com/awslabs/amazon-ecr-credential-helper/ecr-login"
	"github.com/awslabs/amazon-ecr-credential-helper/ecr-login/api"
	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/crane"
	"github.com/google/go-containerregistry/pkg/logs"
)

type RequestBody struct {
	Src  string `json:"src"`
	Dest string `json:"dest"`
}

func init() {
	logs.Warn.SetOutput(os.Stderr)
	logs.Progress.SetOutput(os.Stderr)
}

func handler(ctx context.Context, event events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {
	fmt.Printf("Processing request data for request %s.\n", event.RequestContext.RequestID)
	fmt.Printf("Request Body: %s\n", event.Body)

	// Check if the body is empty
	if event.Body == "" {
		return events.LambdaFunctionURLResponse{
			StatusCode: 400,
			Body:       "Request body is empty",
		}, nil
	}

	var requestBody RequestBody

	if err := json.Unmarshal([]byte(event.Body), &requestBody); err != nil {
		return events.LambdaFunctionURLResponse{
			StatusCode: 400,
			Body:       "Invalid request body",
		}, nil
	}
	log.Printf("COPY EVENT: Copy %s to %s", requestBody.Src, requestBody.Dest)

	ecrHelper := ecr.NewECRHelper(ecr.WithClientFactory(api.DefaultClientFactory{}))

	if err := crane.Copy(requestBody.Src, requestBody.Dest, crane.WithAuthFromKeychain(authn.NewMultiKeychain(authn.DefaultKeychain, authn.NewKeychainFromHelper(ecrHelper)))); err != nil {
		// cancel()
		fmt.Printf("log.Logger:")
		return events.LambdaFunctionURLResponse{
			StatusCode: 400,
			Body:       err.Error(),
		}, nil
	}
	return events.LambdaFunctionURLResponse{
		StatusCode: 200,
		Body:       fmt.Sprintf("COPY EVENT: Copy %s to %s", requestBody.Src, requestBody.Dest),
	}, nil
}

func main() {
	lambda.Start(handler)
}
