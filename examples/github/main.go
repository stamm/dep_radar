package main

import (
	"os"

	"github.com/stamm/dep_radar"
	"github.com/stamm/dep_radar/helpers"
)

func main() {
	recom := dep_radar.MapRecommended{
		"github.com/pkg/errors": dep_radar.Option{
			Recommended: ">=0.8.0",
			Mandatory:   true,
		},
		"github.com/pkg/sftp": dep_radar.Option{
			Recommended: ">=1.3.0",
		},
		"github.com/kr/fs": dep_radar.Option{
			Exclude: true,
		},
		"github.com/pkg/bson": dep_radar.Option{
			Mandatory: true,
		},
		"github.com/pkg/singlefile": dep_radar.Option{
			Exclude: true,
		},
		"github.com/pkg/profile": dep_radar.Option{
			NeedVersion: true,
		},
	}
	helpers.GithubOrg(os.Getenv("GITHUB_TOKEN"), "dep-radar", ":8081", recom)
}
