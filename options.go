package dep_radar

// MapRecommended map with option for libraries
type MapRecommended map[Pkg]Option

// Option for library
type Option struct {
	// Put here restriction for version, for example `>=0.13`
	Recommended string `json:"recommended"`
	// Is it library is mandatory and must be in an app
	Mandatory bool `json:"mandatory"`
	// This library must be absent in an app
	Exclude bool `json:"exclude"`
	// You mustn't use commits hash, but version
	NeedVersion bool `json:"need_version"`
	// UpdatedAt   time.Time
	// UpdateAfter time.Time
}
