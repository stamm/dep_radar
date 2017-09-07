package bitbucketprivate

import (
	"context"
	"encoding/json"
	"fmt"

	i "github.com/stamm/dep_radar/interfaces"
)

type reposResponse struct {
	IsLastPage    bool        `json:"isLastPage"`
	Values        []RepoValue `json:"values"`
	NextPageStart int         `json:"nextPageStart"`
}

// RepoValue sub response from bitbucket for repo
type RepoValue struct {
	Slug string `json:"slug"`
}

// GetRepos get repos
func (a *BitBucketPrivate) getRepos(ctx context.Context, project string, start int) (reposResponse, error) {
	var repos reposResponse
	url := fmt.Sprintf("https://%s/rest/api/1.0/projects/%s/repos?start=%d", a.gitDomain, project, start)
	reposResponse, err := a.httpClient.Get(url)
	if err != nil {
		return repos, err
	}
	err = json.Unmarshal(reposResponse, &repos)

	return repos, err
}

func (a *BitBucketPrivate) GetAllRepos(ctx context.Context, project string) ([]i.Pkg, error) {
	var (
		resultRepos []i.Pkg
		start       int
		isLastPage  bool
	)

	for !isLastPage {
		repos, err := a.getRepos(ctx, project, start)
		if err != nil {
			return resultRepos, err
		}
		for _, repo := range repos.Values {
			resultRepos = append(resultRepos, i.Pkg(a.goGetUrl+"/"+repo.Slug))
		}
		isLastPage = repos.IsLastPage
		start = repos.NextPageStart
	}
	return resultRepos, nil
}
