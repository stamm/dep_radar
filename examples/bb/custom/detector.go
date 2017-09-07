package custom

import (
	"os"

	"github.com/stamm/dep_radar/src/http"
	"github.com/stamm/dep_radar/src/providers"
	bbprivate "github.com/stamm/dep_radar/src/providers/bitbucketprivate"
)

var (
	bbClient   *http.Client
	bbGitUrl   string
	bbApiUrl   string
	bbGoGetUrl string
)

func init() {
	bbGitUrl = os.Getenv("BB_GIT_URL")
	bbGoGetUrl = os.Getenv("BB_GO_GET_URL")
	bbApiUrl = "https://" + bbGitUrl
	bbClient = http.NewClient(
		http.Options{
			User:     os.Getenv("BB_USER"),
			Password: os.Getenv("BB_PASSWORD"),
		}, 10)
}

func Detector() (*bbprivate.BitBucketPrivate, *providers.Detector) {
	// detector := providers.DefaultDetector()
	bbProv := bbprivate.New(bbClient, bbGitUrl, bbGoGetUrl, bbApiUrl)
	detector := providers.NewDetector()
	detector.AddProvider(bbProv)
	return bbProv, detector
}
