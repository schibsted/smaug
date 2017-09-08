package credentials

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/aws/aws-sdk-go/service/sts/stsiface"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestDefaultCredentialsRepositoryFindCredentialsReturnsCredentialsFromAssumeRoleCredentialsProvider(t *testing.T) {
	roleArn := "arn:aws:iam::111111111:myrole/role"

	expiry := time.Now().Add(60 * time.Minute)
	expectedCredentials := &sts.Credentials{
		// Just reflect the role arn to the provider.
		AccessKeyId:     aws.String(roleArn),
		SecretAccessKey: aws.String("Secret"),
		SessionToken:    aws.String("token"),
		Expiration:      &expiry,
	}

	stub := &MockSTSClient{}
	stub.SetCredentials(expectedCredentials)

	repo := NewDefaultCredentialsRepository(stub)

	creds, err := repo.FindCredentialsByRoleArn(roleArn)
	assert.Nil(t, err)

	assert.Equal(t, *expectedCredentials.AccessKeyId, creds.AccessKeyID)
	assert.Equal(t, *expectedCredentials.SecretAccessKey, creds.SecretAccessKey)
	assert.Equal(t, *expectedCredentials.SessionToken, creds.SessionToken)
}

type MockSTSClient struct {
	stsiface.STSAPI
	creds *sts.Credentials
}

func (m *MockSTSClient) SetCredentials(creds *sts.Credentials) {
	m.creds = creds
}
func (m *MockSTSClient) AssumeRole(input *sts.AssumeRoleInput) (*sts.AssumeRoleOutput, error) {
	return &sts.AssumeRoleOutput{
		Credentials: m.creds,
	}, nil
}
