package role

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInMemoryRepository_FindRoleByJobName(t *testing.T) {
	expectedRoleArn := "arn:aws:iam::111111111:myrole/role"
	jobId := "myjob"

	inMemoryRoleRepository := NewInMemoryRoleRepository()
	inMemoryRoleRepository.AddRole(jobId, expectedRoleArn)

	arnRole, err := inMemoryRoleRepository.FindRoleByJobId(jobId)

	assert.Nil(t, err)
	assert.Equal(t, expectedRoleArn, arnRole)
}

func TestInMemoryRepository_FindRoleByJobNameReturnsError(t *testing.T) {
	jobId := "myNonExistentJob"
	inMemoryRoleRepository := NewInMemoryRoleRepository()

	arnRole, err := inMemoryRoleRepository.FindRoleByJobId(jobId)

	if assert.Error(t, err, "An error was expected") {
		assert.Equal(t, err.Error(), fmt.Sprintf("Role for job %s do not exist", jobId))
	}

	assert.Empty(t, arnRole)
}

func TestFileRoleRepository_NewFileRoleRepositoryReturnsErrorIfFileDoesNotExist(t *testing.T) {
	filePath := "invalid path"

	_, err := NewFileRoleRepository(filePath)

	assert.NotNil(t, err, "Invalid path should return an error")
}

func TestFileRoleRepository_FindRoleByJobName(t *testing.T) {
	filePath := "fixtures/roles.ini"
	repository, err := NewFileRoleRepository(filePath)
	assert.Nil(t, err, "Valid path should not return an error")

	jobId := "myjob"
	role, err := repository.FindRoleByJobId(jobId)

	assert.Nil(t, err)
	assert.Equal(t, "arn:aws:iam::111111111:myrole/role", role)

}
