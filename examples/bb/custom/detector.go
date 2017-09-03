package custom

import (
	"os"

	i "github.com/stamm/dep_radar/interfaces"
	"github.com/stamm/dep_radar/src/http"
	"github.com/stamm/dep_radar/src/providers"
	bbprivate "github.com/stamm/dep_radar/src/providers/bitbucketprivate"
)

var (
	bbClient   *http.Client
	bbGitUrl   string
	bbGoGetUrl string
)

func init() {
	bbGitUrl = os.Getenv("BB_GIT_URL")
	bbGoGetUrl = os.Getenv("BB_GO_GET_URL")
	bbClient = http.NewClient(
		http.Options{
			URL:      "https://" + bbGitUrl,
			User:     os.Getenv("BB_USER"),
			Password: os.Getenv("BB_PASSWORD"),
		}, 10)
}

func Detector() *providers.Detector {
	// detector := providers.DefaultDetector()
	detector := providers.NewDetector()
	detector.AddProvider(providers.Provider{
		GoGetUrl: bbGoGetUrl,
		Fn: func(pkg i.Pkg) (i.IProvider, error) {
			return bbprivate.New(pkg, bbClient, bbGitUrl)
		},
	})
	return detector
}
