package interfaces

type Version string
type Hash string
type Pkg string

//go:generate mockery -name=IReposGetter -case=underscore
type IReposGetter interface {
	Apps() []IApp
	// Libs() []ILib
}

//go:generate mockery -name=IApp -case=underscore
type IApp interface {
	IRepo
}

type IRepo interface {
	Package() Pkg
	File(filename string) ([]byte, error)
}

type ILib interface {
	Versions() []Version
}

type IDeps interface {
	Deps() (AppDeps, error)
}

// TODO maybe move
type AppDeps struct {
	Manager Manager
	Deps    map[Pkg]Dep
}

type Dep struct {
	Package string
	Hash    Hash
	Version string
}

type AppListWithDeps map[Pkg]map[Pkg]Dep

type LibMapWithTags map[Pkg][]Tag

type IDepStrategy func(IApp) (AppDeps, error)

type ITagGetter interface {
	Tags() ([]Tag, error)
}

type Tag struct {
	Version string
	Hash    Hash
}

//go:generate mockery -name=IWebClient -case=underscore
type IWebClient interface {
	Get(string) ([]byte, error)
}

type IProvider interface {
	IRepo
	ITagGetter
}
