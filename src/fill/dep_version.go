package fill

import (
	i "github.com/stamm/dep_radar/interfaces"
	"github.com/stamm/dep_radar/src/tags"
)

func DepVersion(deps i.AppDeps, tagList map[i.Pkg][]i.Tag) i.AppDeps {
	for pkg, pkgTags := range tagList {
		findTag, err := tags.FindByHash(deps.Deps[pkg].Hash, pkgTags)
		if err == nil {
			tmp := deps.Deps[pkg]
			tmp.Version = findTag.Version
			deps.Deps[pkg] = tmp
		}
	}
	return deps
}

func DepWithOnlyVersion(deps i.AppDeps, tagList map[i.Pkg][]i.Tag) i.AppDeps {
	result := i.AppDeps{
		Deps: make(map[i.Pkg]i.Dep, 0),
	}
	for pkg, pkgTags := range tagList {
		findTag, err := tags.FindByHash(deps.Deps[pkg].Hash, pkgTags)
		if err == nil {
			tmp := deps.Deps[pkg]
			tmp.Version = findTag.Version
			result.Deps[pkg] = tmp
		}
	}
	return result
}
