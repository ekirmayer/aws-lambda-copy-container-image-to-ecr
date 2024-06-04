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

type SqsBody struct {
	Src struct {
		Image   string `json:"image"`
		Account string `json:"account"`
	} `json:"src"`
	Dst struct {
		Image   string `json:"image"`
		Account string `json:"account"`
	} `json:"dst"`
}

func init() {
	logs.Warn.SetOutput(os.Stderr)
	logs.Progress.SetOutput(os.Stderr)
}

func copy_image(src string, dst string) (error) {
	log.Printf("COPY EVENT: Copy %s to %s", src, dst)

	ecrHelper := ecr.NewECRHelper(ecr.WithClientFactory(api.DefaultClientFactory{}))

	if err := crane.Copy(src, dst, crane.WithAuthFromKeychain(authn.NewMultiKeychain(authn.DefaultKeychain, authn.NewKeychainFromHelper(ecrHelper)))); err != nil {
		
		// cancel()
		fmt.Printf("log.Logger:")
		return err
	}
	return nil
}

func function_handler(ctx context.Context, event events.LambdaFunctionURLRequest) (events.LambdaFunctionURLResponse, error) {
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


	if err := copy_image(requestBody.Src, requestBody.Dest); err != nil {
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

func sqs_handler(ctx context.Context, sqsEvent events.SQSEvent) (error) {
	var sqsBody SqsBody
	for _, message := range sqsEvent.Records {
		fmt.Printf("The message %s for event source %s = %s \n", message.MessageId, message.EventSource, message.Body)
		if err := json.Unmarshal([]byte(message.Body), &sqsBody); err != nil {
			return err
		}

		if err := copy_image(sqsBody.Src.Image, sqsBody.Dst.Image); err != nil {
			fmt.Printf("log.Logger:")
			return err
		}
	}
	return nil
}

func main() {
	var trigger = os.Getenv("TRIGGER")
	switch trigger {
		case "Function":
			lambda.Start(function_handler)
		case "SQS":
			lambda.Start(sqs_handler)
	}

}
