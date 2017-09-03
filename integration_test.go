package depstatus

import (
	"fmt"
	"testing"

	"github.com/stamm/dep_radar/src/deps"
	"github.com/stamm/dep_radar/src/providers/github"
	"github.com/stretchr/testify/assert"
)

func TestIntegration_Github(t *testing.T) {
	t.Parallel()
	assert := assert.New(t)

	app, err := github.New("github.com/Masterminds/glide", github.NewHTTPWrapper("", 2))
	assert.NoError(err)

	depTools := deps.DefaultTools()

	appDeps, err := depTools.Deps(app)
	assert.NoError(err)
	deps := appDeps.Deps
	fmt.Printf("len(deps) = %+v\n", len(deps))
	assert.Equal(len(deps), 5)
	// assert.Equal(appDeps[0].Package, "test")
	// assert.Equal(appDeps[0].Hash, i.Hash("hash1"))

}
