package awscredentials

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
)

func ExampleContextWithAWSCredentials() {

	// These are provided at runtime via secure mechanism - do not hardcode values
	var accessKeyID = "A"
	var secretAccessKey = "B"

	// Store AWS credentials into context
	ctx, err := ContextWithAWSCredentials(context.Background(), NewAWSCredentials("SomeID", accessKeyID, secretAccessKey))
	if err != nil {
		panic(err)
	}

	// Create provider to inject into AWS config
	provider, err := GetDefaultCredentialsProvider(ctx)
	if err != nil {
		panic(err)
	}

	// Create config
	cfg, err := config.LoadDefaultConfig(ctx, config.WithCredentialsProvider(provider))
	if err != nil {
		panic(err)
	}

	// Recover credentials
	c, err := cfg.Credentials.Retrieve(ctx)
	if err != nil {
		panic(err)
	}

	fmt.Println(c.AccessKeyID == accessKeyID && c.SecretAccessKey == secretAccessKey)
	// Output: true
}
