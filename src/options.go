package src

import (
	i "github.com/stamm/dep_radar/interfaces"
)

// MapRecomended map with option for libraries
type MapRecomended map[i.Pkg]Option

// Option for library
type Option struct {
	// Put here restriction for version, for example `>=0.13`
	Recomended string `json:"recomended"`
	// Is it library is mandatory and must be in an app
	Mandatory bool `json:"mandatory"`
	// This library must be absent in an app
	Exclude bool `json:"exclude"`
	// You mustn't use commits hash, but version
	NeedVersion bool `json:"need_version"`
	// UpdatedAt   time.Time
	// UpdateAfter time.Time
}
