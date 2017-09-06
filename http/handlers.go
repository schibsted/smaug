package http

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/go-errors/errors"
	"github.com/schibsted/smaug/credentials"
	"net/http"
	"regexp"
)

var (
	UrlRegexExpression = "^/credentials/(.*)$"
)

func NewCredentialsProviderHandler(provider credentials.CredentialsProvider) *CredentialsProviderHandler {
	return &CredentialsProviderHandler{provider}
}

type CredentialsProviderHandler struct {
	credentialsProvider credentials.CredentialsProvider
}

func (h *CredentialsProviderHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	JobId, err := GetJobIdFromRequest(r)
	log.Debug("JobId: ", JobId)
	if err != nil {
		writeErrorResponse(err.Error(), 404, w)
		return
	}

	smaugCredentials, err := h.credentialsProvider.GetCredentialsForJob(JobId)

	if err != nil {
		writeErrorResponse(err.Error(), 404, w)
		return
	}

	encoded, err := json.Marshal(smaugCredentials)
	if err != nil {
		writeErrorResponse(err.Error(), 500, w)
		log.Error("Couldn't encode credentials")
		return
	}

	w.Header().Add("Content-Type", "application/json")
	fmt.Fprintf(w, string(encoded))
}

func writeErrorResponse(errorMessage string, returnCode int, writer http.ResponseWriter) {
	log.Error(errorMessage)
	writer.WriteHeader(returnCode)
	writer.Write([]byte(errorMessage))
}

func GetJobIdFromRequest(r *http.Request) (string, error) {
	UrlRegex := regexp.MustCompile(UrlRegexExpression)

	match := UrlRegex.FindStringSubmatch(r.URL.EscapedPath())

	if match != nil {
		return match[1], nil
	}

	return "", errors.Errorf("Couldn't get Job Id from request url: %s", r.URL)
}
