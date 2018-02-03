package custom

import (
	"os"

	"github.com/stamm/dep_radar/src/http"
	"github.com/stamm/dep_radar/src/providers"
	bbprivate "github.com/stamm/dep_radar/src/providers/bitbucketprivate"
	"github.com/stamm/dep_radar/src/providers/github"
)

var (
	bbClient   *http.Client
	bbGitURL   string
	bbAPIURL   string
	bbGoGetURL string
)

func init() {
	bbGitURL = os.Getenv("BB_GIT_URL")
	bbGoGetURL = os.Getenv("BB_GO_GET_URL")
	bbAPIURL = "https://" + bbGitURL
	bbClient = http.NewClient(
		http.Options{
			User:     os.Getenv("BB_USER"),
			Password: os.Getenv("BB_PASSWORD"),
		}, 10)
}

// Detector returns bitbucket provider and detector
func Detector() (*bbprivate.Provider, *providers.Detector) {
	bbProv := bbprivate.New(bbClient, bbGitURL, bbGoGetURL, bbAPIURL)
	githubProv := github.New(github.NewHTTPWrapper(os.Getenv("GITHUB_TOKEN"), 10))
	detector := providers.NewDetector().
		AddProvider(bbProv).
		AddProvider(githubProv)
	return bbProv, detector
}
