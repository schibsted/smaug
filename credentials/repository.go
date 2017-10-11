package credentials

import (
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/service/sts/stsiface"
	log "github.com/sirupsen/logrus"
	"time"
)

type CredentialsRepository interface {
	FindCredentialsByRoleArn(string) (*SmaugCredentials, error)
}

// Default Credentials Repository
func NewDefaultCredentialsRepository(client stsiface.STSAPI) *DefaultCredentialsRepository {
	credentialsProviders := make(map[string]*credentials.Credentials)
	return &DefaultCredentialsRepository{
		client,
		credentialsProviders,
	}
}

type DefaultCredentialsRepository struct {
	client              stsiface.STSAPI
	credentialsProvider map[string]*credentials.Credentials
}

func (r *DefaultCredentialsRepository) FindCredentialsByRoleArn(roleArn string) (*SmaugCredentials, error) {
	provider := &stscreds.AssumeRoleProvider{
		Client:       r.client,
		RoleARN:      roleArn,
		Duration:     1 * time.Hour,
		ExpiryWindow: 10 * time.Second,
	}

	if _, ok := r.credentialsProvider[roleArn]; !ok {
		r.credentialsProvider[roleArn] = credentials.NewCredentials(provider)
	}

	creds := r.credentialsProvider[roleArn]

	jobCredentials, err := convertToValidCredentials(provider, creds)

	return jobCredentials, err
}
func convertToValidCredentials(provider *stscreds.AssumeRoleProvider, creds *credentials.Credentials) (*SmaugCredentials, error) {
	roleArn := provider.RoleARN

	credsValue, err := creds.Get()

	if err != nil {
		log.Error(err)
		return nil, err
	}

	expirationTime := time.Now().Add(provider.Duration).Add(-provider.ExpiryWindow)
	expiry := expirationTime.Format("2006-01-02T15:04:05Z")

	smaugCredentials := &SmaugCredentials{
		roleArn,
		credsValue.AccessKeyID,
		credsValue.SecretAccessKey,
		credsValue.SessionToken,
		expiry,
	}

	return smaugCredentials, nil
}
