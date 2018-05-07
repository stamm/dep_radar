package dep_radar

import "context"

// IDepDetector interface for detector of dep
type IDepDetector interface {
	Deps(context.Context, IApp) (AppDeps, error)
	AddTool(IDepTool) IDepDetector
}

//go:generate mockery -name=IDepTool -case=underscore

// IDepTool interface for dep tool
type IDepTool interface {
	Deps(context.Context, IApp) (AppDeps, error)
	Name() string
}
