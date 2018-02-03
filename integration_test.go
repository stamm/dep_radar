package depstatus

import (
	"testing"

	i "github.com/stamm/dep_radar/interfaces"
	"github.com/stamm/dep_radar/src/app"
	"github.com/stamm/dep_radar/src/deps"
	"github.com/stamm/dep_radar/src/deps/dep"
	"github.com/stamm/dep_radar/src/deps/glide"
	"github.com/stamm/dep_radar/src/providers"
	"github.com/stamm/dep_radar/src/providers/github"
	"github.com/stretchr/testify/require"
)

func TestIntegration_Github(t *testing.T) {
	t.Parallel()
	require := require.New(t)

	githubProv := github.New(github.NewHTTPWrapper("", 2))

	provDetector := providers.NewDetector().
		AddProvider(githubProv)

	depDetector := deps.NewDetector().
		AddTool(dep.New()).
		AddTool(glide.New())

	appPkg := i.Pkg("github.com/dep-radar/test_app")
	app, err := app.New(appPkg, "master", provDetector, depDetector)
	require.NoError(err)

	// libs := make(map[i.Pkg]i.Dep)
	// for _, app := range apps {
	// 	deps, err := app.Deps()
	// 	require.NoError(err)
	// 	for _, dep := range deps.Deps {
	// 		if _, ok := libs[dep.Package]; !ok {
	// 			libs[dep.Package] = dep
	// 		}
	// 	}
	// }

	// tags := make(i.LibMapWithTags, len(libs))
	// for pkg := range libs {
	// 	lib, err := lib.New(pkg, provDetector)
	// 	require.NoError(err)
	// 	lig.Tags()
	// }

	appDeps, err := app.Deps()
	require.NoError(err)
	require.Len(appDeps.Deps, 1)
	require.Contains(appDeps.Deps, i.Pkg("github.com/pkg/errors"))
	require.Equal(i.Hash("645ef00459ed84a119197bfb8d8205042c6df63d"), appDeps.Deps[i.Pkg("github.com/pkg/errors")].Hash)
	require.Equal("v0.8.0", appDeps.Deps[i.Pkg("github.com/pkg/errors")].Version)

}
