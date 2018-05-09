package github

import (
	"context"
	"encoding/json"

	"github.com/stamm/dep_radar"
)

type reposResponse struct {
	FullName string `json:"full_name"`
}

func (g *Provider) getOrgReposURL(org string) string {
	return makeAPIURL("orgs/" + org + "/repos")
}

// GetAllOrgRepos get repos for organization
func (g *Provider) GetAllOrgRepos(ctx context.Context, org string) ([]dep_radar.Pkg, error) {
	url := g.getOrgReposURL(org)
	return g.repos(ctx, url)
}

func (g *Provider) getUserURL(user string) string {
	return makeAPIURL("users/" + user)
}

// UserExists check if user exists
func (g *Provider) UserExists(ctx context.Context, user string) bool {
	url := g.getUserURL(user)
	_, err := g.client.Get(ctx, url)
	return err == nil
}

func (g *Provider) getUserReposURL(user string) string {
	return makeAPIURL("users/" + user + "/repos")
}

// GetAllUserRepos get repos for username
func (g *Provider) GetAllUserRepos(ctx context.Context, user string) ([]dep_radar.Pkg, error) {
	url := g.getUserReposURL(user)
	return g.repos(ctx, url)
}

func (g *Provider) repos(ctx context.Context, url string) ([]dep_radar.Pkg, error) {
	repos, err := g.reposFromRequest(ctx, url)
	if err != nil {
		return nil, err
	}
	resultRepos := make([]dep_radar.Pkg, 0, len(repos))
	for _, repo := range repos {
		resultRepos = append(resultRepos, makePkgName(repo.FullName))
	}
	return resultRepos, nil
}

func (g *Provider) reposFromRequest(ctx context.Context, url string) ([]reposResponse, error) {
	resp, err := g.client.Get(ctx, url)
	if err != nil {
		return nil, err
	}
	var repos []reposResponse
	err = json.Unmarshal(resp, &repos)

	return repos, err
}
