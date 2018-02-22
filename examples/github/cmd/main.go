package main

import (
	"os"

	"github.com/stamm/dep_radar/src"
	"github.com/stamm/dep_radar/src/helpers"
)

func main() {
	recom := src.MapRecommended{
		"github.com/pkg/errors": src.Option{
			Recommended: ">=0.8.0",
			Mandatory:   true,
		},
		"github.com/pkg/sftp": src.Option{
			Recommended: ">=1.3.0",
		},
		"github.com/kr/fs": src.Option{
			Exclude: true,
		},
		"github.com/pkg/bson": src.Option{
			Mandatory: true,
		},
		"github.com/pkg/singlefile": src.Option{
			Exclude: true,
		},
		"github.com/pkg/profile": src.Option{
			NeedVersion: true,
		},
	}
	helpers.GithubOrg(os.Getenv("GITHUB_TOKEN"), "dep-radar", ":8081", recom)
}
