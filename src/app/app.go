package app

import (
	"context"

	i "github.com/stamm/dep_radar/interfaces"
)

var _ i.IApp = &App{}

// App struct for an app
type App struct {
	pkg         i.Pkg
	branch      string
	provider    i.IProvider
	deps        i.AppDeps
	depDetector i.IDepDetector
}

// Package returns package name
func (a *App) Package() i.Pkg {
	return a.pkg
}

// Provider returns provider for the app
func (a *App) Provider() i.IProvider {
	return a.provider
}

// Deps returns deps
func (a *App) Deps(ctx context.Context) (i.AppDeps, error) {
	return a.depDetector.Deps(ctx, a)
}

// Branch returns branch
func (a *App) Branch() string {
	return a.branch
}

// New creates app
func New(ctx context.Context, pkg i.Pkg, branch string, detector i.IProviderDetector, depDetector i.IDepDetector) (*App, error) {
	provider, err := detector.Detect(ctx, pkg)
	if err != nil {
		return nil, err
	}

	return &App{
		pkg:         pkg,
		branch:      branch,
		provider:    provider,
		depDetector: depDetector,
	}, nil
}
