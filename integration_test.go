package depstatus

import (
	"fmt"
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

	provDetector := providers.NewDetector()
	provDetector.AddProvider(githubProv)

	depDetector := deps.NewDetector()
	depDetector.AddTool(dep.New())
	depDetector.AddTool(glide.New())

	appsPkgs := []i.Pkg{"github.com/Masterminds/glide"}
	apps := make([]i.IApp, 0, len(appsPkgs))
	for _, appPkg := range appsPkgs {
		app, err := app.New(appPkg, provDetector, depDetector)
		require.NoError(err)
		apps = append(apps, app)
	}

	libs := make(map[i.Pkg]i.Dep)
	for _, app := range apps {
		deps, err := app.Deps()
		require.NoError(err)
		for _, dep := range deps.Deps {
			if _, ok := libs[dep.Package]; !ok {
				libs[dep.Package] = dep
			}
		}
	}

	// tags := make(i.LibMapWithTags, len(libs))
	// for pkg := range libs {
	// 	lib, err := lib.New(pkg, provDetector)
	// 	require.NoError(err)
	// 	lig.Tags()
	// }

	deps, err := apps[0].Deps()
	require.NoError(err)

	fmt.Printf("len(deps) = %+v\n", len(deps.Deps))
	require.Len(deps.Deps, 5)
	// require.Equal(appDeps[0].Package, "test")
	// require.Equal(appDeps[0].Hash, i.Hash("hash1"))

}
