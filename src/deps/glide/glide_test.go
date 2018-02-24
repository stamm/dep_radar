package glide

import (
	"context"
	"errors"
	"testing"

	i "github.com/stamm/dep_radar/interfaces"
	"github.com/stamm/dep_radar/interfaces/mocks"
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
	require.Equal(appDeps.Manager, i.GlideManager)
	deps := appDeps.Deps
	require.Len(deps, 1, "Expect 1 dependency")
	require.Equal(i.Pkg("pkg1"), deps["pkg1"].Package)
	require.Equal(i.Hash("hash1"), deps["pkg1"].Hash)
}

func TestGlide_ErrorOnGettingFile(t *testing.T) {
	t.Parallel()
	require := require.New(t)

	appDeps, err := New().Deps(context.Background(), appMock(nil, errors.New("error")))
	require.EqualError(err, "error")
	require.Len(appDeps.Deps, 0, "Expect 0 dependency")
}

func TestGlide_EmptyFile(t *testing.T) {
	t.Parallel()
	require := require.New(t)

	appDeps, err := New().Deps(context.Background(), appMock([]byte(``), nil))
	require.Error(err)
	require.Len(appDeps.Deps, 0, "Expect 0 dependency")
}

func TestGlide_BadFile(t *testing.T) {
	t.Parallel()
	require := require.New(t)

	appDeps, err := New().Deps(context.Background(), appMock([]byte(`-`), nil))
	require.Error(err)
	require.Len(appDeps.Deps, 0, "Expect 0 dependency")
}

func TestGlide_Name(t *testing.T) {
	t.Parallel()
	require := require.New(t)

	tool := New()
	require.Equal("glide", tool.Name())
}

func appMock(content []byte, err error) *mocks.IApp {
	prov := &mocks.IProvider{}
	prov.On("File", mock.Anything, i.Pkg("app_pkg"), "master", "glide.lock").Return(content, err)
	app := &mocks.IApp{}
	app.On("Package").Return(i.Pkg("app_pkg"))
	app.On("Branch").Return("master")
	app.On("Provider").Return(prov)
	return app
}
