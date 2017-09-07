package deps

import (
	"errors"
	"testing"

	i "github.com/stamm/dep_radar/interfaces"
	"github.com/stamm/dep_radar/interfaces/mocks"
	"github.com/stretchr/testify/require"
)

func Test_Ok(t *testing.T) {
	t.Parallel()
	require := require.New(t)

	mApp := &mocks.IApp{}
	mTool := &mocks.IDepTool{}
	mTool.On("Deps", mApp).Return(mapDep(), nil)

	detector := NewDetector()
	detector.AddTool(mTool)

	appDeps, err := detector.Deps(mApp)
	require.NoError(err)
	require.Equal(i.Manager(-1), appDeps.Manager)
	deps := appDeps.Deps
	require.Len(deps, 1)
	require.Equal(i.Pkg("test"), deps["test"].Package)
	require.Equal(i.Hash("hash1"), deps["test"].Hash)
}

func Test_Error(t *testing.T) {
	t.Parallel()
	require := require.New(t)

	mApp := &mocks.IApp{}

	mTool := &mocks.IDepTool{}
	mTool.On("Deps", mApp).Return(i.AppDeps{}, errors.New("error"))

	detector := NewDetector()
	detector.AddTool(mTool)

	appDeps, err := detector.Deps(mApp)
	require.EqualError(err, "error")
	require.Len(appDeps.Deps, 0)
}

func Test_FirstOk(t *testing.T) {
	t.Parallel()
	require := require.New(t)

	mApp := &mocks.IApp{}

	mTool1 := &mocks.IDepTool{}
	mTool1.On("Deps", mApp).Return(mapDep(), nil)

	mTool2 := &mocks.IDepTool{}
	mTool2.On("Deps", mApp).Return(i.AppDeps{}, errors.New("error"))

	detector := NewDetector()
	detector.AddTool(mTool1)
	detector.AddTool(mTool2)

	appDeps, err := detector.Deps(mApp)
	require.NoError(err)
	deps := appDeps.Deps
	require.Len(deps, 1)
	require.Equal(i.Pkg("test"), deps["test"].Package)
	require.Equal(i.Hash("hash1"), deps["test"].Hash)
}

func Test_FirstBad(t *testing.T) {
	t.Parallel()
	require := require.New(t)

	mApp := &mocks.IApp{}

	mTool1 := &mocks.IDepTool{}
	mTool1.On("Deps", mApp).Return(i.AppDeps{}, errors.New("error"))
	mTool2 := &mocks.IDepTool{}
	mTool2.On("Deps", mApp).Return(mapDep(), nil)

	detector := NewDetector()
	detector.AddTool(mTool1)
	detector.AddTool(mTool2)

	appDeps, err := detector.Deps(mApp)
	require.NoError(err)
	deps := appDeps.Deps
	require.Len(deps, 1)
	require.Equal(i.Pkg("test"), deps["test"].Package)
	require.Equal(i.Hash("hash1"), deps["test"].Hash)
}

func Test_No(t *testing.T) {
	t.Parallel()
	require := require.New(t)

	mApp := &mocks.IApp{}
	detector := NewDetector()

	appDeps, err := detector.Deps(mApp)
	require.EqualError(err, "Bad")
	require.Len(appDeps.Deps, 0)
}

func Test_DefaultDetector(t *testing.T) {
	t.Parallel()
	require := require.New(t)

	detector := DefaultDetector()

	require.Len(detector.Tools, 2)
	require.Equal("dep", detector.Tools[0].Name())
	require.Equal("glide", detector.Tools[1].Name())
}

func mapDep() i.AppDeps {
	return i.AppDeps{
		Manager: i.Manager(-1),
		Deps: map[i.Pkg]i.Dep{
			"test": {
				Package: "test",
				Hash:    "hash1",
			},
		},
	}
}
