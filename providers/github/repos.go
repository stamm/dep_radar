package github

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/stamm/dep_radar"
)

type reposResponse struct {
	FullName string `json:"full_name"`
}

// GetAllOrgRepos get repos for organization
func (g *Provider) GetAllOrgRepos(ctx context.Context, org string) ([]dep_radar.Pkg, error) {
	var (
		resultRepos []dep_radar.Pkg
	)

	url := g.getOrgReposURL(org)
	repos, err := g.getRepos(ctx, url)
	if err != nil {
		return resultRepos, err
	}
	for _, repo := range repos {
		resultRepos = append(resultRepos, dep_radar.Pkg(Prefix+"/"+repo.FullName))
	}
	return resultRepos, nil
}

func (g *Provider) getRepos(ctx context.Context, url string) ([]reposResponse, error) {
	var repos []reposResponse
	reposResponse, err := g.client.Get(ctx, url)
	if err != nil {
		return repos, err
	}
	err = json.Unmarshal(reposResponse, &repos)

	return repos, err
}

func (g *Provider) getOrgReposURL(org string) string {
	return fmt.Sprintf("https://api.github.com/orgs/%s/repos", org)
}

// UserExists check if user exists
func (g *Provider) UserExists(ctx context.Context, user string) bool {
	url := g.getUserURL(user)
	_, err := g.client.Get(ctx, url)
	return err == nil
}

func (g *Provider) getUserURL(user string) string {
	return fmt.Sprintf("https://api.github.com/users/%s", user)
}

// GetAllUserRepos get repos for username
func (g *Provider) GetAllUserRepos(ctx context.Context, user string) ([]dep_radar.Pkg, error) {
	var (
		resultRepos []dep_radar.Pkg
	)

	url := g.getUserReposURL(user)
	repos, err := g.getRepos(ctx, url)
	if err != nil {
		return resultRepos, err
	}
	for _, repo := range repos {
		resultRepos = append(resultRepos, dep_radar.Pkg(Prefix+"/"+repo.FullName))
	}
	return resultRepos, nil
}

func (g *Provider) getUserReposURL(org string) string {
	return fmt.Sprintf("https://api.github.com/users/%s/repos", org)
}
