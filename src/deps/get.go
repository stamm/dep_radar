package deps

import (
	"errors"
	"strings"

	i "github.com/stamm/dep_radar/interfaces"
	"github.com/stamm/dep_radar/src/deps/dep"
	"github.com/stamm/dep_radar/src/deps/glide"
)

type Tools struct {
	Tools []Tool
}

type Tool struct {
	Name string
	Fn   func(i.IApp) (i.AppDeps, error)
}

func (d *Tools) AddTool(prov Tool) {
	d.Tools = append(d.Tools, prov)
}

func (d *Tools) Deps(app i.IApp) (i.AppDeps, error) {
	var errs []string
	for _, tool := range d.Tools {
		deps, err := tool.Fn(app)
		if err == nil {
			// fmt.Printf("deps = %+v\n", deps)
			return deps, nil
		}
		errs = append(errs, err.Error())
	}
	if len(errs) > 0 {
		str := strings.Join(errs, "; ")
		return i.AppDeps{}, errors.New(str)
	}
	return i.AppDeps{}, errors.New("Bad")
}

func NewTools() *Tools {
	return &Tools{}
}

func DefaultTools() *Tools {
	depTool := &Tools{}
	depTool.AddTool(Tool{
		Name: "dep",
		Fn:   dep.Tool,
	})
	depTool.AddTool(Tool{
		Name: "glide",
		Fn:   glide.Tool,
	})
	return depTool
}
