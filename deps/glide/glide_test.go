package glide

import (
	"context"
	"errors"
	"testing"

	"github.com/stamm/dep_radar"
	"github.com/stamm/dep_radar/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestGlide(t *testing.T) {
	t.Parallel()
	require := require.New(t)

	content := []byte(`imports:
- name: pkg1
  version: hash1`)

	appDeps, err := New().Deps(context.Background(), appMock(content, nil))
	require.Nil(err)
	require.Len(appDeps, 1, "Expect 1 dependency")
	require.Equal(dep_radar.Pkg("pkg1"), appDeps["pkg1"].Package)
	require.Equal(dep_radar.Hash("hash1"), appDeps["pkg1"].Hash)
}

func TestGlide_ErrorOnGettingFile(t *testing.T) {
	t.Parallel()
	require := require.New(t)

	appDeps, err := New().Deps(context.Background(), appMock(nil, errors.New("error")))
	require.EqualError(err, "error")
	require.Len(appDeps, 0, "Expect 0 dependency")
}

func TestGlide_EmptyFile(t *testing.T) {
	t.Parallel()
	require := require.New(t)

	appDeps, err := New().Deps(context.Background(), appMock([]byte(``), nil))
	require.Error(err)
	require.Len(appDeps, 0, "Expect 0 dependency")
}

func TestGlide_BadFile(t *testing.T) {
	t.Parallel()
	require := require.New(t)

	appDeps, err := New().Deps(context.Background(), appMock([]byte(`-`), nil))
	require.Error(err)
	require.Len(appDeps, 0, "Expect 0 dependency")
}

func TestGlide_Name(t *testing.T) {
	t.Parallel()
	require := require.New(t)

	tool := New()
	require.Equal("glide", tool.Name())
}

func appMock(content []byte, err error) *mocks.IApp {
	prov := &mocks.IProvider{}
	prov.On("File", mock.Anything, dep_radar.Pkg("app_pkg"), "master", "glide.lock").Return(content, err)
	app := &mocks.IApp{}
	app.On("Package").Return(dep_radar.Pkg("app_pkg"))
	app.On("Branch").Return("master")
	app.On("Provider").Return(prov)
	return app
}
