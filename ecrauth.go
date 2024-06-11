package main

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

func EcrGetToken(region string, role string ) (string, error) {
	var ecrClient *ecr.Client
	ctx := context.TODO()
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(region),
	)
	if err != nil {
		log.Fatal(err)
	}

	if role != "" {
		stsClient := sts.NewFromConfig(cfg)
		provider := stscreds.NewAssumeRoleProvider(stsClient, role)
		cfg.Credentials = aws.NewCredentialsCache(provider)
	} 
	
	ecrClient = ecr.NewFromConfig(cfg)
	authTokenOutput, err := ecrClient.GetAuthorizationToken(ctx, &ecr.GetAuthorizationTokenInput{})
	if err != nil {
		log.Fatalf("unable to get authorization token, %v", err)
		return "", err
	}

	return *authTokenOutput.AuthorizationData[0].AuthorizationToken, nil

}
