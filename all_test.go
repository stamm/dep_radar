package depstatus

import (
	"testing"

	i "github.com/stamm/dep_radar/interfaces"
	"github.com/stamm/dep_radar/interfaces/mocks"
	"github.com/stamm/dep_radar/src/deps"
	"github.com/stretchr/testify/require"
)

func TestAllFlow(t *testing.T) {
	t.Parallel()
	require := require.New(t)
	mApp := &mocks.IApp{}
	mApp.On("Package").Return(i.Pkg("app1"))

	getRepos := &mocks.IReposGetter{}
	getRepos.On("Apps").Return([]i.Pkg{
		"app1",
	}, nil)

	apps, err := getRepos.Apps()
	require.NoError(err)
	require.Len(apps, 1)
	require.Equal(i.Pkg("app1"), apps[0])

	mapDeps := i.AppDeps{
		Manager: i.Manager(-1),
		Deps: map[i.Pkg]i.Dep{
			"test": {
				Package: "test",
				Hash:    "hash1",
			},
		},
	}
	mTool := &mocks.IDepTool{}
	mTool.On("Deps", mApp).Return(mapDeps, nil)

	depDetector := deps.NewDetector()
	depDetector.AddTool(mTool)

	appDeps, err := depDetector.Deps(mApp)
	require.Nil(err)
	dependens := appDeps.Deps
	require.Len(dependens, 1)
	require.Equal(i.Pkg("test"), dependens["test"].Package)
	require.Equal(i.Hash("hash1"), dependens["test"].Hash)

}
