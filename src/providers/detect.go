package providers

import (
	"errors"
	"os"
	"strings"

	i "github.com/stamm/dep_radar/interfaces"
	"github.com/stamm/dep_radar/src/providers/github"
)

var (
	ErrNoProvider             = errors.New("No provider")
	_             i.IDetector = &Detector{}
)

type Detector struct {
	Providers []i.IProvider
	DepTools  []i.IDepTool
}

func NewDetector() *Detector {
	return &Detector{}
}

func (d *Detector) AddProvider(prov i.IProvider) {
	d.Providers = append(d.Providers, prov)
}

func (d *Detector) AddDepTool(tool i.IDepTool) {
	d.DepTools = append(d.DepTools, tool)
}

func (d *Detector) Detect(pkg i.Pkg) (i.IProvider, error) {
	url := string(pkg)
	for _, prov := range d.Providers {
		if strings.HasPrefix(url, prov.GoGetUrl()) {
			return prov, nil
		}
	}
	return nil, ErrNoProvider
}

func DefaultDetector() *Detector {
	githubClient := github.NewHTTPWrapper(os.Getenv("GITHUB_TOKEN"), 2)

	defaultDetector := NewDetector()
	defaultDetector.AddProvider(github.New(githubClient))
	return defaultDetector
}
