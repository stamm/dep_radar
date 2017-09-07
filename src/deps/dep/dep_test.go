package dep

import (
	"errors"
	"testing"

	i "github.com/stamm/dep_radar/interfaces"
	"github.com/stamm/dep_radar/interfaces/mocks"
	"github.com/stretchr/testify/require"
)

func TestDep(t *testing.T) {
	t.Parallel()
	require := require.New(t)

	content := []byte(`[[projects]]
  name = "pkg1"
  revision = "hash1"
`)

	appDeps, err := New().Deps(appMock(content, nil))
	require.NoError(err)
	require.Equal(i.DepManager, appDeps.Manager)
	deps := appDeps.Deps
	require.Len(deps, 1, "Expect 1 dependency")
	require.Equal(i.Pkg("pkg1"), deps["pkg1"].Package)
	require.Equal(i.Hash("hash1"), deps["pkg1"].Hash)
}

func TestDep_ErrorOnGettingFile(t *testing.T) {
	t.Parallel()
	require := require.New(t)

	appDeps, err := New().Deps(appMock(nil, errors.New("error")))
	require.EqualError(err, "error")
	require.Len(appDeps.Deps, 0, "Expect 0 dependency")
}

func TestDep_EmptyFile(t *testing.T) {
	t.Parallel()
	require := require.New(t)

	appDeps, err := New().Deps(appMock([]byte(``), nil))
	require.Error(err)
	require.Equal(len(appDeps.Deps), 0, "Expect 0 dependency")
}

func TestDep_BadFile(t *testing.T) {
	t.Parallel()
	require := require.New(t)

	appDeps, err := New().Deps(appMock([]byte(`-`), nil))
	require.Error(err)
	require.Len(appDeps.Deps, 0, "Expect 0 dependency")
}

func TestDep_Name(t *testing.T) {
	t.Parallel()
	require := require.New(t)

	tool := New()
	require.Equal("dep", tool.Name())
}

func appMock(content []byte, err error) *mocks.IApp {
	prov := &mocks.IProvider{}
	prov.On("File", i.Pkg("app_pkg"), "Gopkg.lock").Return(content, err)
	app := &mocks.IApp{}
	app.On("Package").Return(i.Pkg("app_pkg"))
	app.On("Provider").Return(prov)
	return app
}
