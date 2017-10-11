package credentials

import (
	"fmt"
	"github.com/go-errors/errors"
	"github.com/schibsted/smaug/role"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInMemoryProviderFindCredentialsByJobName(t *testing.T) {
	jobName := "mytestjob"

	expectedAccessKey := "Key"
	expectedSecretKey := "Secret"
	expectedToken := "token"

	credentialsProvider := NewInMemoryCredentialsProvider()
	credentialsProvider.AddCredentials(jobName, GetCredentials("arn:aws:iam::111111111:myrole/role", expectedAccessKey, expectedSecretKey, expectedToken))

	returnedCredentials, err := credentialsProvider.GetCredentialsForJob(jobName)
	assert.Nil(t, err)
	assert.Equal(t, expectedAccessKey, returnedCredentials.AccessKeyID)
	assert.Equal(t, expectedSecretKey, returnedCredentials.SecretAccessKey)
	assert.Equal(t, expectedToken, returnedCredentials.SessionToken)
}

func TestInMemoryProviderFindCredentialsByJobNameReturnsError(t *testing.T) {
	jobName := "mytestjob"
	credentialsProvider := NewInMemoryCredentialsProvider()

	credentialsValues, err := credentialsProvider.GetCredentialsForJob(jobName)

	if assert.Error(t, err, "An error was expected") {
		assert.Equal(t, err.Error(), fmt.Sprintf("Couldn't find credentials for job: %s", jobName))
	}
	assert.Nil(t, credentialsValues)
}

func TestComposableProviderFindCredentialsByJobName(t *testing.T) {
	roleArn := "arn:aws:iam::111111111:myrole/role"
	jobName := "mytestjob"

	expectedAccessKey := "Key"
	expectedSecretKey := "Secret"
	expectedToken := "token"

	roleRepository := role.NewInMemoryRoleRepository()
	roleRepository.AddRole(jobName, roleArn)

	credentialsRepository := &mockCredentialsRepository{
		GetCredentials("arn:aws:iam::111111111:myrole/role", expectedAccessKey, expectedSecretKey, expectedToken),
	}
	credentialsProvider := NewDefaultCredentialsProvider(roleRepository, credentialsRepository)

	returnedCredentials, err := credentialsProvider.GetCredentialsForJob(jobName)
	assert.Nil(t, err)
	assert.Equal(t, expectedAccessKey, returnedCredentials.AccessKeyID)
	assert.Equal(t, expectedSecretKey, returnedCredentials.SecretAccessKey)
	assert.Equal(t, expectedToken, returnedCredentials.SessionToken)
}

func TestComposableProviderFindCredentialsByJobNameReturnsErrorIfJobHasNoRole(t *testing.T) {
	jobName := "mytestjob"

	roleRepository := role.NewInMemoryRoleRepository()
	credentialsRepository := &mockCredentialsRepository{}
	credentialsProvider := NewDefaultCredentialsProvider(roleRepository, credentialsRepository)

	credentialsValues, err := credentialsProvider.GetCredentialsForJob(jobName)

	if assert.Error(t, err, "An error was expected") {
		assert.Equal(t, err.Error(), fmt.Sprintf("Could not get role for job: %s", jobName))
	}

	assert.Nil(t, credentialsValues)
}

func TestComposableProviderFindCredentialsByJobNameReturnsErrorIfRoleHasNoCredentials(t *testing.T) {
	roleArn := "arn:aws:iam::111111111:myrole/role"
	jobName := "mytestjob"

	roleRepository := role.NewInMemoryRoleRepository()
	roleRepository.AddRole(jobName, roleArn)
	credentialsRepository := &mockCredentialsRepository{}
	credentialsProvider := NewDefaultCredentialsProvider(roleRepository, credentialsRepository)

	credentialsValues, err := credentialsProvider.GetCredentialsForJob(jobName)

	if assert.Error(t, err, "An error was expected") {
		assert.Equal(t, err.Error(), fmt.Sprintf("Could not get credentials for role: %s", roleArn))
	}

	assert.Nil(t, credentialsValues)
}

func GetCredentials(roleArn string, accessKey string, secretKey string, token string) *SmaugCredentials {
	creds := &SmaugCredentials{
		// Just reflect the role arn to the provider.
		roleArn,
		accessKey,
		secretKey,
		token,
		"",
	}

	return creds
}

type mockCredentialsRepository struct {
	creds *SmaugCredentials
}

func (r *mockCredentialsRepository) FindCredentialsByRoleArn(roleArn string) (*SmaugCredentials, error) {
	if r.creds == nil || r.creds.RoleArn != roleArn {
		return nil, errors.Errorf("No credentials")
	}
	return r.creds, nil
}
