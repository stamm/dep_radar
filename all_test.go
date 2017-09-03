package depstatus

import (
	"testing"

	i "github.com/stamm/dep_radar/interfaces"
	"github.com/stamm/dep_radar/interfaces/mocks"
	"github.com/stamm/dep_radar/src/deps"
	"github.com/stretchr/testify/assert"
)

func TestAllFlow(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)
	mApp := &mocks.IApp{}
	mApp.On("Package").Return(i.Pkg("app1"))
	getRepos := &mocks.IReposGetter{}
	getRepos.On("Apps").Return([]i.IApp{
		mApp,
	})
	apps := getRepos.Apps()
	assert.Equal(len(apps), 1)
	assert.Equal(apps[0].Package(), i.Pkg("app1"))

	strategy := func(i.IApp) (i.AppDeps, error) {
		return i.AppDeps{
			Deps: map[i.Pkg]i.Dep{
				"test": {
					Package: "test",
					Hash:    "hash1",
				},
			},
		}, nil
	}

	depTool := deps.NewTools()
	depTool.AddTool(deps.Tool{
		Fn: strategy,
	})
	appDeps, err := depTool.Deps(mApp)
	assert.Nil(err)
	deps := appDeps.Deps
	assert.Equal(len(deps), 1)
	assert.Equal(deps["test"].Package, "test")
	assert.Equal(deps["test"].Hash, i.Hash("hash1"))

}
