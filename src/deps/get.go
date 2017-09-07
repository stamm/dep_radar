package deps

import (
	"errors"
	"strings"

	i "github.com/stamm/dep_radar/interfaces"
	"github.com/stamm/dep_radar/src/deps/dep"
	"github.com/stamm/dep_radar/src/deps/glide"
)

var (
	_ i.IDepDetector = &Detector{}
)

type Detector struct {
	Tools []i.IDepTool
}

func (d *Detector) AddTool(tool i.IDepTool) {
	d.Tools = append(d.Tools, tool)
}

func (d *Detector) Deps(app i.IApp) (i.AppDeps, error) {
	var errs []string
	for _, tool := range d.Tools {
		deps, err := tool.Deps(app)
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

func NewDetector() *Detector {
	return &Detector{}
}

func DefaultDetector() *Detector {
	detector := &Detector{}
	detector.AddTool(dep.New())
	detector.AddTool(glide.New())
	return detector
}
