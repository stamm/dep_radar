package app

import (
	i "github.com/stamm/dep_radar/interfaces"
)

var _ i.IApp = &App{}

type App struct {
	pkg         i.Pkg
	provider    i.IProvider
	deps        i.AppDeps
	depDetector i.IDepDetector
}

func (a *App) Package() i.Pkg {
	return a.pkg
}

func (a *App) Provider() i.IProvider {
	return a.provider
}

func (a *App) Deps() (i.AppDeps, error) {
	return a.depDetector.Deps(a)
}

func New(pkg i.Pkg, detector i.IDetector, depDetector i.IDepDetector) (*App, error) {
	provider, err := detector.Detect(pkg)
	if err != nil {
		return nil, err
	}

	return &App{
		pkg:         pkg,
		provider:    provider,
		depDetector: depDetector,
	}, nil
}
