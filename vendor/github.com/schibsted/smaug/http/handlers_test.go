package http_test

import (
	"fmt"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/aws/aws-sdk-go/service/sts/stsiface"
	"github.com/schibsted/smaug/credentials"
	http_pkg "github.com/schibsted/smaug/http"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetJobNameFromRequestReturnsJobNameIfUrlIsCorrect(t *testing.T) {
	req, _ := http.NewRequest("GET", "/credentials/dd9404e6-d09b-4124-85bf-98018695b05d", nil)

	jobName, _ := http_pkg.GetJobIdFromRequest(req)

	assert.Equal(t, "dd9404e6-d09b-4124-85bf-98018695b05d", jobName)

}

func TestGetJobNameFromRequestReturnsErrorIfUrlIsIncorrect(t *testing.T) {
	testUrl := "/invalid-url"
	req, _ := http.NewRequest("GET", testUrl, nil)

	jobName, err := http_pkg.GetJobIdFromRequest(req)

	if assert.Error(t, err, "An error was expected") {
		assert.Equal(t, err.Error(), fmt.Sprintf("Couldn't get Job Id from request url: %s", testUrl))
	}

	assert.Empty(t, jobName)
}

func TestSecurityProviderHandlerReturnsErrorIfUrlIsIncorrect(t *testing.T) {
	testUrl := "/credentials"
	req, _ := http.NewRequest("GET", testUrl, nil)

	credentialsProvider := credentials.NewInMemoryCredentialsProvider()
	handler := http_pkg.NewCredentialsProviderHandler(credentialsProvider)

	writer := httptest.NewRecorder()
	handler.ServeHTTP(writer, req)

	body, _ := ioutil.ReadAll(writer.Body)
	assert.Equal(t, 404, writer.Code)
	assert.Equal(t, fmt.Sprintf("Couldn't get Job Id from request url: %s", testUrl), string(body))
}

func TestSecurityProviderHandlerReturnNoCredentialsIfCredentialsForJobDoNotExist(t *testing.T) {
	testJobName := "dd9404e6-d09b-4124-85bf-98018695b05d"
	req, _ := http.NewRequest("GET", "/credentials/"+testJobName, nil)

	credentialsProvider := credentials.NewInMemoryCredentialsProvider()
	handler := http_pkg.NewCredentialsProviderHandler(credentialsProvider)

	writer := httptest.NewRecorder()
	handler.ServeHTTP(writer, req)

	body, _ := ioutil.ReadAll(writer.Body)
	assert.Equal(t, 404, writer.Code)
	assert.Equal(t, "Couldn't find credentials for job: dd9404e6-d09b-4124-85bf-98018695b05d", string(body))
}

func TestSecurityProviderHandlerReturnCredentialsValueIfCredentialsForJobExists(t *testing.T) {
	jobId := "dd9404e6-d09b-4124-85bf-98018695b05d"
	req, _ := http.NewRequest("GET", "/credentials/"+jobId, nil)

	credentialsProvider := credentials.NewInMemoryCredentialsProvider()
	cr := GetCredentials("arn:aws:iam::111111111:myrole/role")
	fmt.Println("test creds", cr)
	credentialsProvider.AddCredentials(jobId, cr)
	handler := http_pkg.NewCredentialsProviderHandler(credentialsProvider)

	writer := httptest.NewRecorder()
	handler.ServeHTTP(writer, req)

	body, _ := ioutil.ReadAll(writer.Body)
	assert.Equal(t, 200, writer.Code)
	expectedResponseBody := "{\"RoleArn\":\"myKey\",\"AccessKeyId\":\"mySecret\",\"SecretAccessKey\":\"MyToken\",\"Token\":\"MyProvider\",\"Expiration\":\"2017-04-11T21:49:00Z\"}"
	assert.Equal(t, expectedResponseBody, string(body))
}

func GetCredentials(roleArn string) *credentials.SmaugCredentials {
	creds := &credentials.SmaugCredentials{
		"myKey",
		"mySecret",
		"MyToken",
		"MyProvider",
		"2017-04-11T21:49:00Z",
	}

	return creds
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
