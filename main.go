package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type EventBody struct {
	Src struct {
		Image   string `json:"image"`
		Account string `json:"account"`
		Region  string `json:"region"`
		Role    string `json:"role"`
		Type    string `json:"type"`
	} `json:"src"`
	Dst struct {
		Image   string `json:"image"`
		Account string `json:"account"`
		Region  string `json:"region"`
		Role    string `json:"role"`
		Type    string `json:"type"`
	} `json:"dst"`
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

	var requestBody EventBody

	if err := json.Unmarshal([]byte(event.Body), &requestBody); err != nil {
		return events.LambdaFunctionURLResponse{
			StatusCode: 400,
			Body:       "Invalid request body",
		}, nil
	}
	srcRepository, dstRepository := setup_registries(requestBody)
	if err := copy_image(srcRepository, dstRepository); err != nil {
		// cancel()
		fmt.Printf("log.Logger:")
		return events.LambdaFunctionURLResponse{
			StatusCode: 400,
			Body:       err.Error(),
		}, nil
	}
	return events.LambdaFunctionURLResponse{
		StatusCode: 200,
		Body:       fmt.Sprintf("COPY EVENT: Copy %s to %s", srcRepository, dstRepository),
	}, nil
}

func sqs_handler(ctx context.Context, sqsEvent events.SQSEvent) error {
	var sqsBody EventBody
	for _, message := range sqsEvent.Records {
		fmt.Printf("The message %s for event source %s = %s \n", message.MessageId, message.EventSource, message.Body)
		if err := json.Unmarshal([]byte(message.Body), &sqsBody); err != nil {
			return err
		}
		srcRepository, dstRepository := setup_registries(sqsBody)

		if err := copy_image(srcRepository, dstRepository); err != nil {
			fmt.Printf("log.Logger:")
			return err
		}
	}
	return nil
}


func setup_registries(images EventBody) (string, string) {
	var srcRegistry string
	var srcRepository string
	var dstRegistry string
	var dstRepository string
	

	if images.Src.Type == "ecr" {
		srcToken, _ := EcrGetToken(images.Src.Region, images.Src.Role)
		srcRegistry = fmt.Sprintf("%s.dkr.ecr.%s.amazonaws.com", images.Src.Account, images.Src.Region)
		dockerConfig.AddRegistry(srcRegistry, srcToken)
		srcRepository = fmt.Sprintf("%s/%s", srcRegistry, images.Src.Image)
	}

	if images.Src.Type == "ecr" {
		dstToken, _ := EcrGetToken(images.Dst.Region, images.Dst.Role)
		dstRegistry = fmt.Sprintf("%s.dkr.ecr.%s.amazonaws.com", images.Dst.Account, images.Dst.Region)
		dockerConfig.AddRegistry(dstRegistry, dstToken)
		dstRepository = fmt.Sprintf("%s/%s", dstRegistry, images.Dst.Image)
	}

	dockerFilePath, _ := GenerateAuth()
	fmt.Println(dockerFilePath)
	os.Setenv("DOCKER_CONFIG", dockerFilePath)

	return srcRepository, dstRepository
}

func test() error {
	log.Println("test")

	jsonFile, err := os.Open("./tests/sqs_event.json")
	if err != nil {
		log.Fatal(err)
	}
	defer jsonFile.Close()

	byteValue, _ := io.ReadAll(jsonFile)

	var result EventBody
	json.Unmarshal([]byte(byteValue), &result)

	srcRepository, dstRepository := setup_registries(result)

	

	if err := copy_image(srcRepository, dstRepository); err != nil {
		fmt.Println(err.Error())
		return err
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
	case "TEST":
		test()
	default:
		log.Println("No trigger handler selected!")
	}

}
