package glide

import (
	"context"
	"fmt"

	"github.com/Masterminds/glide/cfg"
	"github.com/Masterminds/glide/path"
	"github.com/stamm/dep_radar"
)

var (
	_ dep_radar.IDepTool = &Tool{}
)

// Tool glide
type Tool struct{}

// New creates new instance of tool
func New() *Tool {
	return &Tool{}
}

// Name gets the name for glide
func (t *Tool) Name() string {
	return "glide"
}

// Deps returns deps
func (t *Tool) Deps(ctx context.Context, a dep_radar.IApp) (dep_radar.AppDeps, error) {
	res := dep_radar.AppDeps{}
	content, err := a.Provider().File(ctx, a.Package(), a.Branch(), path.LockFile)
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
	res = make(map[dep_radar.Pkg]dep_radar.Dep, len(lockFile.Imports))
	for _, imp := range lockFile.Imports {
		// fmt.Printf("imp.Name = %+v\n", imp.Name)
		res[dep_radar.Pkg(imp.Name)] = dep_radar.Dep{
			Package: dep_radar.Pkg(imp.Name),
			Hash:    dep_radar.Hash(imp.Version),
		}
	}
	// fmt.Printf("deps = %+v\n", deps)
	return res, nil
}
