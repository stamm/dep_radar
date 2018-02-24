package interfaces

import "context"

type Version string
type Hash string
type Pkg string

//go:generate mockery -name=IApp -case=underscore
type IApp interface {
	Package() Pkg
	Provider() IProvider
	Deps(context.Context) (AppDeps, error)
	Branch() string
}

type ILib interface {
	Versions() []Version
}

type IDeps interface {
	Deps(context.Context) (AppDeps, error)
}

// TODO maybe move
type AppDeps struct {
	Manager Manager
	Deps    map[Pkg]Dep
}

type Dep struct {
	Package Pkg
	Hash    Hash
	Version string
}

type AppListWithDeps map[Pkg]map[Pkg]Dep

type LibMapWithTags map[Pkg][]Tag

type IDepStrategy func(IApp) (AppDeps, error)

type Tag struct {
	Version string
	Hash    Hash
}

//go:generate mockery -name=IWebClient -case=underscore
type IWebClient interface {
	Get(context.Context, string) ([]byte, error)
}

// Provider

type IFileGetter interface {
	File(ctx context.Context, pkg Pkg, branch, filename string) ([]byte, error)
}

type ITagGetter interface {
	Tags(context.Context, Pkg) ([]Tag, error)
}

//go:generate mockery -name=IProvider -case=underscore
type IProvider interface {
	IFileGetter
	ITagGetter
	GoGetUrl() string
}

//go:generate mockery -name=IReposGetter -case=underscore
type IReposGetter interface {
	Apps(context.Context) ([]Pkg, error)
	// Libs() []ILib
}

type IProviderDetector interface {
	Detect(context.Context, Pkg) (IProvider, error)
}

type IDepDetector interface {
	Deps(context.Context, IApp) (AppDeps, error)
	AddTool(IDepTool) IDepDetector
}

//go:generate mockery -name=IDepTool -case=underscore
type IDepTool interface {
	Deps(context.Context, IApp) (AppDeps, error)
	Name() string
}
