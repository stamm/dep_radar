package providers

import (
	"context"
	"testing"

	"github.com/stamm/dep_radar"
	"github.com/stamm/dep_radar/providers/github"
	"github.com/stretchr/testify/assert"
)

func TestDetect_ExpectGithub(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	detect := DefaultDetector()
	prov, err := detect.Detect(context.Background(), dep_radar.Pkg("github.com/golang/dep"))
	assert.NoError(err)
	assert.IsType(&github.Provider{}, prov)
	assert.Implements((*dep_radar.IProvider)(nil), prov)
}

func TestDetect_ExpectError(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	detect := DefaultDetector()
	app, err := detect.Detect(context.Background(), dep_radar.Pkg("gopkg.in/yaml.v2"))
	assert.EqualError(err, "No provider")
	assert.Nil(app)
}
