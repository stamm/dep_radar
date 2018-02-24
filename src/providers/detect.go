package providers

import (
	"context"
	"errors"
	"os"
	"strings"

	i "github.com/stamm/dep_radar/interfaces"
	"github.com/stamm/dep_radar/src/providers/github"
)

const (
	defaultGithubLimit = 20
)

var (
	// ErrNoProvider says that no provider was found
	ErrNoProvider                     = errors.New("No provider")
	_             i.IProviderDetector = &Detector{}
)

// Detector for provider
type Detector struct {
	Providers []i.IProvider
}

// NewDetector creates new detector
func NewDetector() *Detector {
	return &Detector{}
}

// AddProvider add provider to detector
func (d *Detector) AddProvider(prov i.IProvider) *Detector {
	d.Providers = append(d.Providers, prov)
	return d
}

// Detect right provider
func (d *Detector) Detect(ctx context.Context, pkg i.Pkg) (i.IProvider, error) {
	url := string(pkg)
	for _, prov := range d.Providers {
		if strings.HasPrefix(url, prov.GoGetURL()) {
			return prov, nil
		}
	}
	return nil, ErrNoProvider
}

// DefaultDetector return detector that support only github
func DefaultDetector() *Detector {
	githubClient := github.NewHTTPWrapper(os.Getenv("GITHUB_TOKEN"), defaultGithubLimit)

	defaultDetector := NewDetector()
	defaultDetector.AddProvider(github.New(githubClient))
	return defaultDetector
}
