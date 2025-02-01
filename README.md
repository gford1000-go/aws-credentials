[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](https://en.wikipedia.org/wiki/MIT_License)
[![Documentation](https://img.shields.io/badge/Documentation-GoDoc-green.svg)](https://godoc.org/github.com/gford1000-go/aws/credentials)

# Credentials | Simplified AWS credential management

Store AWS Key and Secret within a `context` and then access via an `aws.CredentialsProvider` when required.

```go
func main() {
    // These are provided at runtime via secure mechanism - do not hardcode values
    var accessKeyID = "A"
    var secretAccessKey = "B"

    // Store AWS credentials into context
    ctx, _err_ := ContextWithAWSCredentials(context.Background(), NewAWSCredentials(accessKeyID, secretAccessKey))

    // Create provider when needed, to inject into AWS config
    provider, _ := GetCredentialsProvider(ctx)

    // Create config using values
    cfg, _ := config.LoadDefaultConfig(ctx, config.WithCredentialsProvider(provider))
}
```
