package credentials

import (
	"github.com/go-errors/errors"
	"github.com/schibsted/smaug/role"
)

type CredentialsProvider interface {
	GetCredentialsForJob(string) (*SmaugCredentials, error)
}

// InMemory Credentials Provider
func NewInMemoryCredentialsProvider() *InMemoryCredentialsProvider {
	return &InMemoryCredentialsProvider{}
}

type InMemoryCredentialsProvider struct {
	credentials map[string]*SmaugCredentials
}

func (r *InMemoryCredentialsProvider) AddCredentials(jobId string, creds *SmaugCredentials) {
	if r.credentials == nil {
		r.credentials = make(map[string]*SmaugCredentials)
	}
	r.credentials[jobId] = creds
}
func (r *InMemoryCredentialsProvider) GetCredentialsForJob(jobId string) (*SmaugCredentials, error) {
	if creds, ok := r.credentials[jobId]; ok {
		return creds, nil
	}

	return nil, errors.Errorf("Couldn't find credentials for job: %s", jobId)
}

// Default Credentials Provider
func NewDefaultCredentialsProvider(roleRepository role.RoleRepository, credentialsRepository CredentialsRepository) *DefaultCredentialsProvider {
	return &DefaultCredentialsProvider{
		roleRepository,
		credentialsRepository,
	}
}

type DefaultCredentialsProvider struct {
	roleRepository        role.RoleRepository
	credentialsRepository CredentialsRepository
}

func (provider *DefaultCredentialsProvider) GetCredentialsForJob(jobId string) (*SmaugCredentials, error) {
	roleArn, err := provider.roleRepository.FindRoleByJobId(jobId)

	if err != nil {
		return nil, errors.Errorf("Could not get role for job: %s", jobId)
	}

	creds, err := provider.credentialsRepository.FindCredentialsByRoleArn(roleArn)

	if err != nil {
		return nil, errors.Errorf("Could not get credentials for role: %s", roleArn)
	}

	return creds, nil
}
