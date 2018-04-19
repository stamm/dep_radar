package deps

import (
	"context"
	"errors"
	"testing"

	"github.com/stamm/dep_radar"
	"github.com/stamm/dep_radar/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func Test_Ok(t *testing.T) {
	t.Parallel()
	require := require.New(t)

	mApp := &mocks.IApp{}
	mTool := &mocks.IDepTool{}
	mTool.On("Deps", mock.Anything, mApp).Return(mapDep(), nil)

	detector := NewDetector()
	detector.AddTool(mTool)

	appDeps, err := detector.Deps(context.Background(), mApp)
	require.NoError(err)
	require.Len(appDeps, 1)
	require.Equal(dep_radar.Pkg("test"), appDeps["test"].Package)
	require.Equal(dep_radar.Hash("hash1"), appDeps["test"].Hash)
}

func Test_Error(t *testing.T) {
	t.Parallel()
	require := require.New(t)

	mApp := &mocks.IApp{}

	mTool := &mocks.IDepTool{}
	mTool.On("Deps", mock.Anything, mApp).Return(dep_radar.AppDeps{}, errors.New("error"))

	detector := NewDetector()
	detector.AddTool(mTool)

	appDeps, err := detector.Deps(context.Background(), mApp)
	require.EqualError(err, "error")
	require.Len(appDeps, 0)
}

func Test_FirstOk(t *testing.T) {
	t.Parallel()
	require := require.New(t)

	mApp := &mocks.IApp{}

	mTool1 := &mocks.IDepTool{}
	mTool1.On("Deps", mock.Anything, mApp).Return(mapDep(), nil)

	mTool2 := &mocks.IDepTool{}
	mTool2.On("Deps", mock.Anything, mApp).Return(dep_radar.AppDeps{}, errors.New("error"))

	detector := NewDetector()
	detector.AddTool(mTool1)
	detector.AddTool(mTool2)

	appDeps, err := detector.Deps(context.Background(), mApp)
	require.NoError(err)
	require.Len(appDeps, 1)
	require.Equal(dep_radar.Pkg("test"), appDeps["test"].Package)
	require.Equal(dep_radar.Hash("hash1"), appDeps["test"].Hash)
}

func Test_FirstBad(t *testing.T) {
	t.Parallel()
	require := require.New(t)

	mApp := &mocks.IApp{}

	mTool1 := &mocks.IDepTool{}
	mTool1.On("Deps", mock.Anything, mApp).Return(dep_radar.AppDeps{}, errors.New("error"))
	mTool2 := &mocks.IDepTool{}
	mTool2.On("Deps", mock.Anything, mApp).Return(mapDep(), nil)

	detector := NewDetector()
	detector.AddTool(mTool1)
	detector.AddTool(mTool2)

	appDeps, err := detector.Deps(context.Background(), mApp)
	require.NoError(err)
	require.Len(appDeps, 1)
	require.Equal(dep_radar.Pkg("test"), appDeps["test"].Package)
	require.Equal(dep_radar.Hash("hash1"), appDeps["test"].Hash)
}

func Test_No(t *testing.T) {
	t.Parallel()
	require := require.New(t)

	mApp := &mocks.IApp{}
	detector := NewDetector()

	appDeps, err := detector.Deps(context.Background(), mApp)
	require.EqualError(err, "Bad")
	require.Len(appDeps, 0)
}

func Test_DefaultDetector(t *testing.T) {
	t.Parallel()
	require := require.New(t)

	detector := DefaultDetector()

	require.Len(detector.Tools, 2)
	require.Equal("dep", detector.Tools[0].Name())
	require.Equal("glide", detector.Tools[1].Name())
}

func mapDep() dep_radar.AppDeps {
	return dep_radar.AppDeps{
		"test": {
			Package: "test",
			Hash:    "hash1",
		},
	}
}
