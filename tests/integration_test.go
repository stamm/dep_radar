package src_test

import (
	"context"
	"testing"

	"github.com/stamm/dep_radar"
	"github.com/stamm/dep_radar/app"
	"github.com/stamm/dep_radar/deps"
	"github.com/stamm/dep_radar/deps/dep"
	"github.com/stamm/dep_radar/deps/glide"
	"github.com/stamm/dep_radar/providers"
	"github.com/stamm/dep_radar/providers/github"
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

	appPkg := dep_radar.Pkg("github.com/dep-radar/test_app")
	app, err := app.New(context.Background(), appPkg, "master", provDetector, depDetector)
	require.NoError(err)

	// libs := make(map[dep_radar.Pkg]dep_radar.Dep)
	// for _, app := range apps {
	// 	deps, err := app.Deps()
	// 	require.NoError(err)
	// 	for _, dep := range deps.Deps {
	// 		if _, ok := libs[dep.Package]; !ok {
	// 			libs[dep.Package] = dep
	// 		}
	// 	}
	// }

	// tags := make(dep_radar.LibMapWithTags, len(libs))
	// for pkg := range libs {
	// 	lib, err := lib.New(pkg, provDetector)
	// 	require.NoError(err)
	// 	lig.Tags()
	// }

	appDeps, err := app.Deps(context.Background())
	require.NoError(err)
	require.Len(appDeps, 5)
	require.Contains(appDeps, dep_radar.Pkg("github.com/pkg/errors"))
	require.Equal(dep_radar.Hash("645ef00459ed84a119197bfb8d8205042c6df63d"), appDeps[dep_radar.Pkg("github.com/pkg/errors")].Hash)
	require.Equal("v0.8.0", appDeps[dep_radar.Pkg("github.com/pkg/errors")].Version)

}
