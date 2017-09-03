package providers

import (
	"errors"
	"os"
	"strings"

	i "github.com/stamm/dep_radar/interfaces"
	"github.com/stamm/dep_radar/src/providers/github"
)

var (
	ErrNoProvider = errors.New("No provider")
)

type Provider struct {
	GitUrl   string
	GoGetUrl string
	Fn       func(i.Pkg) (i.IProvider, error)
}

type Detector struct {
	Providers []Provider
}

func NewDetector() *Detector {
	return &Detector{}
}

func (d *Detector) AddProvider(prov Provider) {
	d.Providers = append(d.Providers, prov)
}

func (d *Detector) Detect(pkg i.Pkg) (i.IProvider, error) {
	url := string(pkg)
	for _, prov := range d.Providers {
		if strings.HasPrefix(url, prov.GoGetUrl) {
			return prov.Fn(pkg)
		}
	}
	return nil, ErrNoProvider
}

func DefaultDetector() *Detector {
	githubClient := github.NewHTTPWrapper(os.Getenv("GITHUB_TOKEN"), 2)

	defaultDetector := NewDetector()
	defaultDetector.AddProvider(Provider{
		GitUrl:   "github.com",
		GoGetUrl: "github.com",
		Fn: func(pkg i.Pkg) (i.IProvider, error) {
			return github.New(pkg, githubClient)
		},
	})
	return defaultDetector
}
