package src_test

import (
	"context"
	"testing"

	"github.com/stamm/dep_radar"
	"github.com/stamm/dep_radar/deps"
	"github.com/stamm/dep_radar/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestAllFlow(t *testing.T) {
	t.Parallel()
	require := require.New(t)
	mApp := &mocks.IApp{}
	mApp.On("Package").Return(dep_radar.Pkg("app1"))

	mapDeps := dep_radar.AppDeps{
		"test": {
			Package: "test",
			Hash:    "hash1",
		},
	}
	mTool := &mocks.IDepTool{}
	mTool.On("Deps", mock.Anything, mApp).Return(mapDeps, nil)

	depDetector := deps.NewDetector()
	depDetector.AddTool(mTool)

	appDeps, err := depDetector.Deps(context.Background(), mApp)
	require.Nil(err)
	require.Len(appDeps, 1)
	require.Equal(dep_radar.Pkg("test"), appDeps["test"].Package)
	require.Equal(dep_radar.Hash("hash1"), appDeps["test"].Hash)

}
