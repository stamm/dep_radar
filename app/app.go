package app

import (
	"context"

	"github.com/stamm/dep_radar"
)

var _ dep_radar.IApp = &App{}

// App struct for an app.
type App struct {
	pkg         dep_radar.Pkg
	branch      string
	provider    dep_radar.IProvider
	deps        dep_radar.AppDeps
	depDetector dep_radar.IDepDetector
}

// Package returns package name.
func (a *App) Package() dep_radar.Pkg {
	return a.pkg
}

// Provider returns provider for the app.
func (a *App) Provider() dep_radar.IProvider {
	return a.provider
}

// Deps returns deps.
func (a *App) Deps(ctx context.Context) (dep_radar.AppDeps, error) {
	return a.depDetector.Deps(ctx, a)
}

// Branch returns branch.
func (a *App) Branch() string {
	return a.branch
}

// New creates app.
func New(ctx context.Context, pkg dep_radar.Pkg, branch string, detector dep_radar.IProviderDetector, depDetector dep_radar.IDepDetector) (*App, error) {
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
