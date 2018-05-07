package dep

import (
	"context"
	"fmt"

	"github.com/BurntSushi/toml"
	"github.com/pkg/errors"
	"github.com/stamm/dep_radar"
)

const (
	file = "Gopkg.lock"
)

var (
	_ dep_radar.IDepTool = &Tool{}
)

// Tool dep
type Tool struct{}

type rawLock struct {
	Projects []rawLockedProject `toml:"projects"`
}

type rawLockedProject struct {
	Name     string   `toml:"name"`
	Branch   string   `toml:"branch,omitempty"`
	Revision string   `toml:"revision"`
	Version  string   `toml:"version,omitempty"`
	Source   string   `toml:"source,omitempty"`
	Packages []string `toml:"packages"`
}

// New creates new instance of tool
func New() *Tool {
	return &Tool{}
}

// Name gets the name for dep
func (t *Tool) Name() string {
	return "dep"
}

// Deps returns deps
func (t *Tool) Deps(ctx context.Context, a dep_radar.IApp) (dep_radar.AppDeps, error) {
	res := dep_radar.AppDeps{}
	content, err := a.Provider().File(ctx, a.Package(), a.Branch(), file)
	if err != nil {
		return res, err
	}
	if len(content) == 0 {
		return res, fmt.Errorf("File %s is empty", file)
	}
	// fmt.Printf("content = \n%+v\n---------------------\n", string(content))
	raw := rawLock{}
	err = toml.Unmarshal(content, &raw)
	if err != nil {
		return res, errors.Wrap(err, "Unable to parse the lock as TOML")
	}
	// fmt.Printf("lockFile = %+v\n", lockFile.Imports)
	res = make(map[dep_radar.Pkg]dep_radar.Dep, len(raw.Projects))
	for _, imp := range raw.Projects {
		// fmt.Printf("imp.Name = %+v\n", imp.Name)
		res[dep_radar.Pkg(imp.Name)] = dep_radar.Dep{
			Package: dep_radar.Pkg(imp.Name),
			Hash:    dep_radar.Hash(imp.Revision),
			Version: imp.Version,
		}
	}
	// fmt.Printf("deps = %+v\n", deps)
	return res, nil
}
