package dep_radar

import "context"

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
type AppDeps map[Pkg]Dep

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
