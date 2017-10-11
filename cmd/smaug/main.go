package main

import (
	"flag"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/schibsted/smaug/credentials"
	http_pkg "github.com/schibsted/smaug/http"
	"github.com/schibsted/smaug/role"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
)

var (
	serverAddr                string
	DEFAULT_SERVER_ADDRESS    = ":8080"
	verbose                   bool
	credentialsRepositoryFile string
)

func main() {
	parseFlags()
	setLogLevel()
	validateParameters()

	roleRepository, err := role.NewFileRoleRepository(credentialsRepositoryFile)
	if err != nil {
		log.Error(err)
		os.Exit(1)
	}
	stsClient := createStsClient()
	credentialsRepo := credentials.NewDefaultCredentialsRepository(stsClient)
	credentialsProvider := credentials.NewDefaultCredentialsProvider(roleRepository, credentialsRepo)
	credentialsRequestHandler := http_pkg.NewCredentialsProviderHandler(credentialsProvider)
	http.Handle("/credentials/", credentialsRequestHandler)
	http.HandleFunc("/health-check/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Ok"))
	})

	log.Info("Listening on ", serverAddr)
	log.Panic(http.ListenAndServe(serverAddr, nil))
}

func setLogLevel() {
	log.SetLevel(log.InfoLevel)

	if verbose {
		log.SetLevel(log.DebugLevel)
	}
}

func validateParameters() {
	if credentialsRepositoryFile == "" {
		log.Error("credentials-repository-file is required")
		os.Exit(1)
	}
}

func parseFlags() {
	flag.BoolVar(&verbose, "verbose", false, "Enable verbosity")
	flag.StringVar(&serverAddr, "server-address", DEFAULT_SERVER_ADDRESS, "Server address")
	flag.StringVar(&credentialsRepositoryFile, "credentials-repository-file", "", "Credentials Repository False")

	flag.Parse()
}

func createStsClient() *sts.STS {
	config := &aws.Config{
		Region: aws.String("eu-west-1"),
	}
	sess := session.Must(session.NewSession(config))
	stsClient := sts.New(sess)
	return stsClient
}
