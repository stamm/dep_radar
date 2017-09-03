package providers

import (
	"testing"

	i "github.com/stamm/dep_radar/interfaces"
	"github.com/stamm/dep_radar/src/providers/github"
	"github.com/stretchr/testify/assert"
)

func TestDetect_ExpectGithub(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	detect := DefaultDetector()
	app, err := detect.Detect(i.Pkg("github.com/golang/dep"))
	assert.NoError(err)
	assert.IsType(github.Github{}, app)
	assert.Implements((*i.IProvider)(nil), app)
}

func TestDetect_ExpectError(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	detect := DefaultDetector()
	app, err := detect.Detect(i.Pkg("gopkg.in/yaml.v2"))
	assert.EqualError(err, "No provider")
	assert.Nil(app)
}
