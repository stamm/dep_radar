package src

import (
	i "github.com/stamm/dep_radar/interfaces"
)

// MapRecomended map with option for libraries
type MapRecomended map[i.Pkg]Option

// Option for library
type Option struct {
	Recomended  string
	Mandatory   bool
	Exclude     bool
	NeedVersion bool
	// UpdatedAt   time.Time
	// UpdateAfter time.Time
}
