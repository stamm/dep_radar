package glide

import (
	"context"
	"fmt"

	"github.com/Masterminds/glide/cfg"
	"github.com/Masterminds/glide/path"
	i "github.com/stamm/dep_radar/interfaces"
)

var (
	_ i.IDepTool = &Tool{}
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
func (t *Tool) Deps(ctx context.Context, a i.IApp) (i.AppDeps, error) {
	res := i.AppDeps{}
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
