package deps

import (
	"errors"
	"testing"

	i "github.com/stamm/dep_radar/interfaces"
	"github.com/stamm/dep_radar/interfaces/mocks"
	"github.com/stretchr/testify/assert"
)

func TestDep_Ok(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	mApp := &mocks.IApp{}

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

	depTool := NewTools()
	depTool.AddTool(Tool{
		Fn: strategy,
	})
	appDeps, err := depTool.Deps(mApp)
	assert.Nil(err)
	deps := appDeps.Deps
	assert.Equal(len(deps), 1)
	assert.Equal(deps["test"].Package, "test")
	assert.Equal(deps["test"].Hash, i.Hash("hash1"))
}

func TestDep_Error(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	mApp := &mocks.IApp{}

	strategy := func(i.IApp) (i.AppDeps, error) {
		return i.AppDeps{}, errors.New("error")
	}

	depTool := NewTools()
	depTool.AddTool(Tool{
		Fn: strategy,
	})
	appDeps, err := depTool.Deps(mApp)
	assert.EqualError(err, "error")
	assert.Equal(len(appDeps.Deps), 0)
}

func TestDep_FirstOk(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	mApp := &mocks.IApp{}

	strategies := []i.IDepStrategy{
		func(i.IApp) (i.AppDeps, error) {
			return i.AppDeps{
				Deps: map[i.Pkg]i.Dep{
					"test": {
						Package: "test",
						Hash:    "hash1",
					},
				},
			}, nil
		},
		func(i.IApp) (i.AppDeps, error) {
			return i.AppDeps{}, errors.New("error")
		},
	}

	depTool := NewTools()
	for _, strategy := range strategies {
		depTool.AddTool(Tool{
			Fn: strategy,
		})
	}
	appDeps, err := depTool.Deps(mApp)
	assert.Nil(err)
	deps := appDeps.Deps
	assert.Equal(len(deps), 1)
	assert.Equal(deps["test"].Package, "test")
	assert.Equal(deps["test"].Hash, i.Hash("hash1"))
}

func TestDep_FirstBad(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	mApp := &mocks.IApp{}

	strategies := []i.IDepStrategy{
		func(i.IApp) (i.AppDeps, error) {
			return i.AppDeps{}, errors.New("error")
		},
		func(i.IApp) (i.AppDeps, error) {
			return i.AppDeps{
				Deps: map[i.Pkg]i.Dep{
					"test": {
						Package: "test",
						Hash:    "hash1",
					},
				},
			}, nil
		},
	}

	depTool := NewTools()
	for _, strategy := range strategies {
		depTool.AddTool(Tool{
			Fn: strategy,
		})
	}
	appDeps, err := depTool.Deps(mApp)
	assert.Nil(err)
	deps := appDeps.Deps
	assert.Equal(len(deps), 1)
	assert.Equal(deps["test"].Package, "test")
	assert.Equal(deps["test"].Hash, i.Hash("hash1"))
}

func TestDep_No(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	mApp := &mocks.IApp{}

	depTool := NewTools()
	appDeps, err := depTool.Deps(mApp)
	assert.EqualError(err, "Bad")
	assert.Equal(len(appDeps.Deps), 0)
}
