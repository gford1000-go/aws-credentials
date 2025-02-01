package credentials

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	c "github.com/aws/aws-sdk-go-v2/credentials"
)

// AWSCredentials holds details of the AWS key and secret key
type AWSCredentials struct {
	id              CredentialsID
	accessKeyID     string
	secretAccessKey string
}

// CredentialsID specifies a given set of credentials
type CredentialsID string

// NewAWSCredentials creates an instance of AWSCredentials from supplied details
// with no verification of their correctness
func NewAWSCredentials(id CredentialsID, accessKeyID, secretAccessKey string) *AWSCredentials {
	return &AWSCredentials{
		id:              id,
		accessKeyID:     accessKeyID,
		secretAccessKey: secretAccessKey,
	}
}

type awsCredentials string

var singletonAWSCredentialsKey awsCredentials = "theAWSKeys"

type credentialsMap struct {
	defaultID CredentialsID
	cred      map[CredentialsID]*AWSCredentials
}

// ErrContextMustNotBeNil raisewd when ContextWithAWSCredentials is called with a nil context
var ErrContextMustNotBeNil = errors.New("no context provided")

// ErrCredentialsMustNotBeNil raised when a nil credential pointer is passed to ContextWithAWSCredentials
var ErrCredentialsMustNotBeNil = errors.New("no AWS credentials provided")

// ErrInvalidCredentialsProvided raised when the AWSCredentials do not contain sufficient data
var ErrInvalidCredentialsProvided = errors.New("provided credentials are invalid")

// ContextWithAWSCredentials returns a new context based on the supplied, which stores the credentials.
// Multiple credentials can be stored within the same context, facilitating granular access controls.
// If this is the first set of credentials added to the context, they are treated as the default set.
func ContextWithAWSCredentials(ctx context.Context, credentials *AWSCredentials) (context.Context, error) {
	if ctx == nil {
		return ctx, ErrContextMustNotBeNil
	}
	if credentials == nil {
		return ctx, ErrCredentialsMustNotBeNil
	}
	if credentials.accessKeyID == "" || credentials.secretAccessKey == "" {
		return ctx, ErrInvalidCredentialsProvided
	}

	m, err := credentialsFromContext(ctx)
	if err != nil {
		return ctx, nil
	}
	if m == nil {

		c := &credentialsMap{
			defaultID: credentials.id,
			cred: map[CredentialsID]*AWSCredentials{
				credentials.id: {
					accessKeyID:     credentials.accessKeyID,
					secretAccessKey: credentials.secretAccessKey,
					id:              credentials.id,
				},
			},
		}

		return context.WithValue(ctx, singletonAWSCredentialsKey, c), nil
	}

	m.cred[credentials.id] = &AWSCredentials{
		accessKeyID:     credentials.accessKeyID,
		secretAccessKey: credentials.secretAccessKey,
		id:              credentials.id,
	}

	return ctx, nil
}

// ErrUnexpectedErrorRetrievingCredentials raised if a context holds a key to the AWSCredentials, but the object
// stored against that key is the wrong type
var ErrUnexpectedErrorRetrievingCredentials = errors.New("invalid object stored as AWSCredentials")

// credentialsFromContext returns the AWS credentials if they are present in the context
func credentialsFromContext(ctx context.Context) (*credentialsMap, error) {
	v := ctx.Value(singletonAWSCredentialsKey)
	if v == nil {
		return nil, nil // Not being present is not an error
	}
	m, ok := v.(*credentialsMap)
	if !ok {
		return nil, ErrUnexpectedErrorRetrievingCredentials
	}
	return m, nil
}

// ErrUnknownCredentialsID raised if the credentials for the specified id cannot be found
var ErrUnknownCredentialsID = errors.New("credentials for the specified id are unavailable")

// GetCredentialsProvider returns an aws.CredentialsProvider if the context contains AWSCredentials
func GetCredentialsProvider(ctx context.Context, id CredentialsID) (aws.CredentialsProvider, error) {

	m, err := credentialsFromContext(ctx)
	if err != nil {
		return nil, err
	}

	a, ok := m.cred[id]
	if !ok {
		return nil, ErrUnknownCredentialsID
	}

	return c.NewStaticCredentialsProvider(
		a.accessKeyID,
		a.secretAccessKey,
		""), nil
}

// GetDefaultCredentialsProvider returns an aws.CredentialsProvider using the default set of credentials
func GetDefaultCredentialsProvider(ctx context.Context) (aws.CredentialsProvider, error) {
	m, err := credentialsFromContext(ctx)
	if err != nil {
		return nil, err
	}

	a, ok := m.cred[m.defaultID]
	if !ok {
		return nil, ErrUnknownCredentialsID
	}

	return c.NewStaticCredentialsProvider(
		a.accessKeyID,
		a.secretAccessKey,
		""), nil

}
