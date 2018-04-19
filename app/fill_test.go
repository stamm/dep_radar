package app

import (
	"context"
	"errors"
	"testing"

	"github.com/stamm/dep_radar"
	"github.com/stamm/dep_radar/mocks"
	"github.com/stamm/dep_radar/providers"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var (
	mapDeps = dep_radar.AppDeps{
		"github.com/test": {
			Package: "github.com/test",
			Hash:    "hash1",
		},
	}
)

func TestFill(t *testing.T) {
	t.Parallel()
	require := require.New(t)

	apps := make(chan dep_radar.IApp)
	app := createApp("github.com/app1", mapDeps)
	go func() {
		apps <- app
		close(apps)
	}()

	prov := createProv()
	provDetector := providers.NewDetector().AddProvider(prov)

	appsWithDeps, libsWithTags := GetTags(context.Background(), apps, provDetector)

	require.Len(appsWithDeps, 1)
	require.Contains(appsWithDeps, dep_radar.Pkg("github.com/app1"))
	require.Len(libsWithTags, 1)
	require.Contains(libsWithTags, dep_radar.Pkg("github.com/test"))
	prov.AssertExpectations(t)
	app.AssertExpectations(t)
}

func TestFillSameDepLib(t *testing.T) {
	t.Parallel()
	require := require.New(t)

	apps := make(chan dep_radar.IApp)
	app := createApp("github.com/app1", mapDeps)
	app2 := createApp("github.com/app2", mapDeps)
	go func() {
		apps <- app
		apps <- app2
		close(apps)
	}()

	prov := createProv()
	provDetector := providers.NewDetector().AddProvider(prov)

	appsWithDeps, libsWithTags := GetTags(context.Background(), apps, provDetector)

	require.Len(appsWithDeps, 2)
	require.Contains(appsWithDeps, dep_radar.Pkg("github.com/app1"))
	require.Contains(appsWithDeps, dep_radar.Pkg("github.com/app2"))
	require.Len(libsWithTags, 1)
	require.Contains(libsWithTags, dep_radar.Pkg("github.com/test"))
	prov.AssertExpectations(t)
	app.AssertExpectations(t)
	app2.AssertExpectations(t)
}

func TestFillDetectorErr(t *testing.T) {
	t.Parallel()
	require := require.New(t)

	apps := make(chan dep_radar.IApp)
	app := createApp("github.com/app1", mapDeps)
	go func() {
		apps <- app
		close(apps)
	}()

	provDetector := providers.NewDetector()

	appsWithDeps, libsWithTags := GetTags(context.Background(), apps, provDetector)

	require.Len(appsWithDeps, 1)
	require.Contains(appsWithDeps, dep_radar.Pkg("github.com/app1"))
	require.Len(libsWithTags, 0)
	app.AssertExpectations(t)
}

func TestFillTagsErr(t *testing.T) {
	t.Parallel()
	require := require.New(t)

	apps := make(chan dep_radar.IApp)
	app := createApp("github.com/app1", mapDeps)
	go func() {
		apps <- app
		close(apps)
	}()

	prov := &mocks.IProvider{}
	prov.On("GoGetURL").Return("github.com").Times(1)
	prov.On("Tags", mock.Anything, dep_radar.Pkg("github.com/test")).Return(nil, errors.New("err")).Times(1)
	provDetector := providers.NewDetector().AddProvider(prov)

	appsWithDeps, libsWithTags := GetTags(context.Background(), apps, provDetector)

	require.Len(appsWithDeps, 1)
	require.Contains(appsWithDeps, dep_radar.Pkg("github.com/app1"))
	require.Len(libsWithTags, 0)
	prov.AssertExpectations(t)
	app.AssertExpectations(t)
}

func TestFillDepsErr(t *testing.T) {
	t.Parallel()
	require := require.New(t)

	apps := make(chan dep_radar.IApp)
	app := &mocks.IApp{}
	go func() {
		app.On("Package").Return(dep_radar.Pkg("github.com/app1")).Times(1)
		app.On("Deps", mock.Anything).Return(dep_radar.AppDeps{}, errors.New("dep error")).Times(1)
		apps <- app
		close(apps)
	}()

	prov := &mocks.IProvider{}
	provDetector := providers.NewDetector().AddProvider(prov)

	appsWithDeps, libsWithTags := GetTags(context.Background(), apps, provDetector)

	require.Len(appsWithDeps, 0)
	require.Len(libsWithTags, 0)
	prov.AssertExpectations(t)
	app.AssertExpectations(t)
}

func createApp(pkg string, deps dep_radar.AppDeps) *mocks.IApp {
	app := &mocks.IApp{}
	app.On("Package").Return(dep_radar.Pkg(pkg)).Times(1)
	app.On("Deps", mock.Anything).Return(deps, nil).Times(1)
	return app
}

func createProv() *mocks.IProvider {
	prov := &mocks.IProvider{}
	prov.On("GoGetURL").Return("github.com").Times(1)
	prov.On("Tags", mock.Anything, dep_radar.Pkg("github.com/test")).Return([]dep_radar.Tag{
		{
			Version: "1.0.0",
			Hash:    dep_radar.Hash("hash1"),
		},
	}, nil).Times(1)
	return prov
}
