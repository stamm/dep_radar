package glide

import (
	"errors"
	"testing"

	i "github.com/stamm/dep_radar/interfaces"
	"github.com/stamm/dep_radar/interfaces/mocks"
	"github.com/stretchr/testify/assert"
)

func TestGlide(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	content := []byte(`imports:
- name: pkg1
  version: hash1`)

	app := &mocks.IApp{}
	app.On("File", "glide.lock").Return(content, nil)

	appDeps, err := Tool(app)
	assert.Nil(err)
	assert.Equal(appDeps.Manager, i.GlideManager)
	deps := appDeps.Deps
	assert.Equal(len(deps), 1, "Expect 1 dependency")
	assert.Equal(deps["pkg1"].Package, "pkg1")
	assert.Equal(deps["pkg1"].Hash, i.Hash("hash1"))
}

func TestGlide_ErrorOnGettingFile(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	app := &mocks.IApp{}
	app.On("File", "glide.lock").Return(nil, errors.New("error"))

	deps, err := Tool(app)
	assert.EqualError(err, "error")
	assert.Equal(len(deps.Deps), 0, "Expect 0 dependency")
}

func TestGlide_EmptyFile(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	content := []byte(``)

	app := &mocks.IApp{}
	app.On("File", "glide.lock").Return(content, nil)

	deps, err := Tool(app)
	assert.Error(err)
	assert.Equal(len(deps.Deps), 0, "Expect 0 dependency")
}

func TestGlide_BadFile(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	content := []byte(`-`)

	app := &mocks.IApp{}
	app.On("File", "glide.lock").Return(content, nil)

	deps, err := Tool(app)
	assert.Error(err)
	assert.Equal(len(deps.Deps), 0, "Expect 0 dependency")
}
