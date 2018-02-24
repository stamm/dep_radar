package interfaces

import "context"

// Version is an alias for string
type Version string

// Hash is an alias for string
type Hash string

// Pkg is an alias for string
type Pkg string

//go:generate mockery -name=IApp -case=underscore

// IApp interface for an app
type IApp interface {
	Package() Pkg
	Provider() IProvider
	Deps(context.Context) (AppDeps, error)
	Branch() string
}

// AppDeps shows deps in an app
// TODO maybe move
type AppDeps struct {
	Deps map[Pkg]Dep
}

// Dep is a struct for particular library
type Dep struct {
	Package Pkg
	Hash    Hash
	Version string
}

// AppListWithDeps is an alias for apps with partucular libs
type AppListWithDeps map[Pkg]map[Pkg]Dep

// LibMapWithTags is an alias for libs with all tags
type LibMapWithTags map[Pkg][]Tag

// Tag contains version and hash
type Tag struct {
	Version string
	Hash    Hash
}

//go:generate mockery -name=IWebClient -case=underscore

// IWebClient is an interface for getting file
type IWebClient interface {
	Get(context.Context, string) ([]byte, error)
}

// Provider

// IFileGetter shows can provider get file
type IFileGetter interface {
	File(ctx context.Context, pkg Pkg, branch, filename string) ([]byte, error)
}

// ITagGetter shows can provider get tags
type ITagGetter interface {
	Tags(context.Context, Pkg) ([]Tag, error)
}

//go:generate mockery -name=IProvider -case=underscore

// IProvider describe provider like github or bitbucket
type IProvider interface {
	IFileGetter
	ITagGetter
	GoGetURL() string
}

//go:generate mockery -name=IReposGetter -case=underscore

// IReposGetter shows can provider get repos
type IReposGetter interface {
	Apps(context.Context) ([]Pkg, error)
}

// IProviderDetector interface for detector of providers
type IProviderDetector interface {
	Detect(context.Context, Pkg) (IProvider, error)
}

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
