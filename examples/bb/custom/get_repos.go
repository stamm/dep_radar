package custom

import (
	"context"
	"encoding/json"
	"fmt"

	i "github.com/stamm/dep_radar/interfaces"
)

type ReposResponse struct {
	IsLastPage    bool        `json:"isLastPage"`
	Values        []RepoValue `json:"values"`
	NextPageStart int         `json:"nextPageStart"`
}

// RepoValue sub response from bitbucket for repo
type RepoValue struct {
	Slug string `json:"slug"`
}

// GetRepos get repos
func GetRepos(ctx context.Context, start int) (ReposResponse, error) {
	var (
		repos ReposResponse
	)
	url := fmt.Sprintf("rest/api/1.0/projects/%s/repos?start=%d", "GO", start)
	reposResponse, err := bbClient.Get(url)
	if err != nil {
		return repos, err
	}
	err = json.Unmarshal(reposResponse, &repos)

	return repos, err
}

func GetAllRepos(ctx context.Context) ([]i.Pkg, error) {
	var (
		resultRepos []i.Pkg
		start       int
		isLastPage  bool
	)

	for !isLastPage {
		fmt.Println("Do!")
		repos, err := GetRepos(ctx, start)
		if err != nil {
			return resultRepos, err
		}
		for _, repo := range repos.Values {
			resultRepos = append(resultRepos, i.Pkg(bbGoGetUrl+"/"+repo.Slug))
		}
		isLastPage = repos.IsLastPage
		start = repos.NextPageStart
	}
	return resultRepos, nil
}
