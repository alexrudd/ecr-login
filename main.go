package main

import (
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws/ec2metadata"
	"github.com/aws/aws-sdk-go-v2/aws/external"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
)

func main() {
	// Using the SDK's default configuration, loading additional config
	// and credentials values from the environment variables, shared
	// credentials, and shared configuration files
	cfg, err := external.LoadDefaultAWSConfig()
	if err != nil {
		panic("unable to load SDK config, " + err.Error())
	}

	if cfg.Region == "" {
		// Get ec2 metadata client
		mdc := ec2metadata.New(cfg)

		// Get current region
		region, err := mdc.Region()
		if err != nil {
			panic("failed to determine AWS region, " + err.Error())
		}
		cfg.Region = region
	}

	// Using the Config value, create the ECR client
	svc := ecr.New(cfg)

	// Build the request with its input parameters
	req := svc.GetAuthorizationTokenRequest(&ecr.GetAuthorizationTokenInput{
		RegistryIds: nil,
	})

	// Send the request, and get the response or error back
	resp, err := req.Send()
	if err != nil {
		panic("failed to get ECR authorization token, " + err.Error())
	}

	// Decode the ecr authorization token, should be string in format user:password
	decoded, err := base64.StdEncoding.DecodeString(*resp.AuthorizationData[0].AuthorizationToken)
	if err != nil {
		panic("failed to decode ECR authorization token, " + err.Error())
	}

	// Split on colon
	up := strings.Split(string(decoded), ":")
	if len(up) != 2 {
		panic("ECR authorization token was not in expected 'user:password' format, " + string(decoded))
	}

	// Print as docker login command
	fmt.Printf("docker login -u %s -p %s %s", up[0], up[1], *resp.AuthorizationData[0].ProxyEndpoint)
}
