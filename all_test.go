package depstatus

import (
	"context"
	"testing"

	i "github.com/stamm/dep_radar/interfaces"
	"github.com/stamm/dep_radar/interfaces/mocks"
	"github.com/stamm/dep_radar/src/deps"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestAllFlow(t *testing.T) {
	t.Parallel()
	require := require.New(t)
	mApp := &mocks.IApp{}
	mApp.On("Package").Return(i.Pkg("app1"))

	getRepos := &mocks.IReposGetter{}
	getRepos.On("Apps", mock.Anything).Return([]i.Pkg{
		"app1",
	}, nil)

	apps, err := getRepos.Apps(context.Background())
	require.NoError(err)
	require.Len(apps, 1)
	require.Equal(i.Pkg("app1"), apps[0])

	mapDeps := i.AppDeps{
		Deps: map[i.Pkg]i.Dep{
			"test": {
				Package: "test",
				Hash:    "hash1",
			},
		},
	}
	mTool := &mocks.IDepTool{}
	mTool.On("Deps", mock.Anything, mApp).Return(mapDeps, nil)

	depDetector := deps.NewDetector()
	depDetector.AddTool(mTool)

	appDeps, err := depDetector.Deps(context.Background(), mApp)
	require.Nil(err)
	dependens := appDeps.Deps
	require.Len(dependens, 1)
	require.Equal(i.Pkg("test"), dependens["test"].Package)
	require.Equal(i.Hash("hash1"), dependens["test"].Hash)

}
