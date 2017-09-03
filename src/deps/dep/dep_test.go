package dep

import (
	"errors"
	"testing"

	i "github.com/stamm/dep_radar/interfaces"
	"github.com/stamm/dep_radar/interfaces/mocks"
	"github.com/stretchr/testify/assert"
)

func TestDep(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	content := []byte(`[[projects]]
  name = "pkg1"
  revision = "hash1"
`)

	app := &mocks.IApp{}
	app.On("File", "Gopkg.lock").Return(content, nil)

	appDeps, err := Tool(app)
	assert.Nil(err)
	assert.Equal(appDeps.Manager, i.DepManager)
	deps := appDeps.Deps
	assert.Equal(len(deps), 1, "Expect 1 dependency")
	assert.Equal(deps["pkg1"].Package, "pkg1")
	assert.Equal(deps["pkg1"].Hash, i.Hash("hash1"))
}

func TestDep_ErrorOnGettingFile(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	app := &mocks.IApp{}
	app.On("File", "Gopkg.lock").Return(nil, errors.New("error"))

	deps, err := Tool(app)
	assert.EqualError(err, "error")
	assert.Equal(len(deps.Deps), 0, "Expect 0 dependency")
}

func TestDep_EmptyFile(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	content := []byte(``)

	app := &mocks.IApp{}
	app.On("File", "Gopkg.lock").Return(content, nil)

	deps, err := Tool(app)
	assert.Error(err)
	assert.Equal(len(deps.Deps), 0, "Expect 0 dependency")
}

func TestDep_BadFile(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	content := []byte(`-`)

	app := &mocks.IApp{}
	app.On("File", "Gopkg.lock").Return(content, nil)

	deps, err := Tool(app)
	assert.Error(err)
	assert.Equal(len(deps.Deps), 0, "Expect 0 dependency")
}
