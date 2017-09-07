package glide

import (
	"fmt"

	"github.com/Masterminds/glide/cfg"
	"github.com/Masterminds/glide/path"
	i "github.com/stamm/dep_radar/interfaces"
)

var (
	_ i.IDepTool = &Tool{}
)

type Tool struct{}

func (t *Tool) Name() string {
	return "glide"
}

func (t *Tool) Deps(a i.IApp) (i.AppDeps, error) {
	res := i.AppDeps{
		Manager: i.GlideManager,
	}
	content, err := a.Provider().File(a.Package(), path.LockFile)
	if err != nil {
		return res, err
	}
	if len(content) == 0 {
		return res, fmt.Errorf("File %s is empty", path.LockFile)
	}
	// fmt.Printf("content = \n%+v\n---------------------\n", string(content))

	lockFile, err := cfg.LockfileFromYaml(content)
	if err != nil {
		return res, err
	}
	// fmt.Printf("lockFile = %+v\n", lockFile.Imports)
	res.Deps = make(map[i.Pkg]i.Dep, len(lockFile.Imports))
	for _, imp := range lockFile.Imports {
		// fmt.Printf("imp.Name = %+v\n", imp.Name)
		res.Deps[i.Pkg(imp.Name)] = i.Dep{
			Package: i.Pkg(imp.Name),
			Hash:    i.Hash(imp.Version),
		}
	}
	// fmt.Printf("deps = %+v\n", deps)
	return res, nil
}

func New() *Tool {
	return &Tool{}
}
