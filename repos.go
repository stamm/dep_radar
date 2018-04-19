package dep_radar

import "context"

// Version is an alias for string.
type Version string

// Hash is an alias for string.
type Hash string

// Pkg is an alias for string.
type Pkg string

// Tag contains version and hash.
type Tag struct {
	Version string
	Hash    Hash
}

// Provider

// IFileGetter shows can provider get file.
type IFileGetter interface {
	File(ctx context.Context, pkg Pkg, branch, filename string) ([]byte, error)
}

// ITagGetter shows can provider get tags.
type ITagGetter interface {
	Tags(context.Context, Pkg) ([]Tag, error)
}

//go:generate mockery -name=IProvider -case=underscore

// IProvider describe provider like github or bitbucket.
type IProvider interface {
	IFileGetter
	ITagGetter
	GoGetURL() string
}

// IProviderDetector interface for detector of providers
type IProviderDetector interface {
	Detect(context.Context, Pkg) (IProvider, error)
}
