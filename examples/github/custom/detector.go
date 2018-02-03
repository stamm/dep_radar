package custom

import (
	"os"

	"github.com/stamm/dep_radar/src/providers"
	"github.com/stamm/dep_radar/src/providers/github"
)

// Detector returns github detector
func Detector() *providers.Detector {
	githubProv := github.New(github.NewHTTPWrapper(os.Getenv("GITHUB_TOKEN"), 10))
	return providers.NewDetector().AddProvider(githubProv)
}
