package deps

import (
	"context"
	"errors"
	"strings"
	"sync"

	i "github.com/stamm/dep_radar/interfaces"
	"github.com/stamm/dep_radar/src/deps/dep"
	"github.com/stamm/dep_radar/src/deps/glide"
)

var (
	_ i.IDepDetector = &Detector{}
)

// Detector for deps
type Detector struct {
	Tools []i.IDepTool
}

// AddTool adds dep tool
func (d *Detector) AddTool(tool i.IDepTool) i.IDepDetector {
	d.Tools = append(d.Tools, tool)
	return d
}

// Deps returns deps for app
func (d *Detector) Deps(ctx context.Context, app i.IApp) (i.AppDeps, error) {
	var (
		errs []string
		wg   sync.WaitGroup
	)
	depResult := make(chan i.AppDeps)
	done := make(chan struct{})

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	for _, tool := range d.Tools {
		wg.Add(1)
		go func(tool i.IDepTool) {
			defer wg.Done()
			deps, err := tool.Deps(ctx, app)
			if err == nil {
				depResult <- deps
				return
				// return deps, nil
			}
			errs = append(errs, err.Error())
		}(tool)
	}

	go func() {
		wg.Wait()
		done <- struct{}{}
	}()

	select {
	case <-ctx.Done():
		return i.AppDeps{}, ctx.Err()
	// TODO: check that done will be called later that result
	case result := <-depResult:
		cancel()
		return result, nil
	case <-done:
	}

	if len(errs) > 0 {
		str := strings.Join(errs, "; ")
		return i.AppDeps{}, errors.New(str)
	}
	return i.AppDeps{}, errors.New("Bad")
}

// NewDetector returns empty detector
func NewDetector() *Detector {
	return &Detector{}
}

// DefaultDetector returns detector with all dependency systems
func DefaultDetector() *Detector {
	detector := &Detector{}
	detector.AddTool(dep.New())
	detector.AddTool(glide.New())
	return detector
}
