package bitbucketprivate

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"

	i "github.com/stamm/dep_radar/interfaces"
)

var (
	_ i.ITagGetter = &Provider{}
	_ i.IProvider  = &Provider{}
)

// TagsResponse response from bitbucket for tags
type TagsResponse struct {
	IsLastPage    bool        `json:"isLastPage"`
	Values        []TagsValue `json:"values"`
	NextPageStart int         `json:"nextPageStart"`
}

// TagsValue sub response from bibucket for tags
type TagsValue struct {
	ID           string `json:"id"`
	LatestCommit string `json:"latestCommit"`
}

// Repo struct for repo
type Repo struct {
	Project     string
	Name        string
	Slug        string
	HashVersion string
	Version     string
	Deps        map[string]Repo
	Tags        []string
	Branch      string
}

// FileOpts options for getting file
type FileOpts struct {
	Project    string
	RepoName   string
	Filename   string
	BranchName string
}

type Provider struct {
	project    string
	httpClient i.IWebClient
	gitDomain  string
	goGetUrl   string
	apiUrl     string
	muMap      sync.RWMutex
	mapProject map[i.Pkg]string
}

type Options struct {
	URL      string
	User     string
	Password string
}

func New(httpClient i.IWebClient, gitDomain, goGetUrl, apiUrl string) *Provider {
	return &Provider{
		httpClient: httpClient,
		goGetUrl:   goGetUrl,
		gitDomain:  gitDomain,
		apiUrl:     apiUrl,
		mapProject: make(map[i.Pkg]string),
	}
}

func (p *Provider) File(pkg i.Pkg, branch, name string) ([]byte, error) {
	project, err := p.getProject(pkg)
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf("%s/projects/%s/repos/%s/raw/%s?at=refs%%2Fheads%%2F%s", p.apiUrl, project, p.repoName(pkg), name, branch)
	return p.httpClient.Get(url)
}

func (p *Provider) GoGetUrl() string {
	return p.goGetUrl
}

// todo cache result in map
func (p *Provider) getProject(pkg i.Pkg) (string, error) {
	p.muMap.RLock()
	if project, ok := p.mapProject[pkg]; ok {
		p.muMap.RUnlock()
		return project, nil
	}
	p.muMap.RUnlock()
	// if p.project != "" {
	// 	return nil
	// }
	project, err := GetProject(p.httpClient, pkg, p.gitDomain)
	if err != nil {
		return "", err
	}
	p.muMap.Lock()
	p.mapProject[pkg] = project
	p.muMap.Unlock()
	// fmt.Printf("project = %+v\n", project)
	// p.project = project
	return project, nil
}

// Tags get tags from bitbucket
func (p *Provider) Tags(pkg i.Pkg) ([]i.Tag, error) {
	project, err := p.getProject(pkg)
	if err != nil {
		return nil, err
	}
	var (
		tagsResult []i.Tag
		start      int
		isLastPage bool
	)
	for !isLastPage {
		tags, err := p.tags(pkg, project, start)
		if err != nil {
			return tagsResult, fmt.Errorf("Error on getting tags: %s", err)
		}
		for _, tag := range tags.Values {
			tagVersion := strings.Replace(tag.ID, "refs/tags/", "", 1)
			tagsResult = append(tagsResult, i.Tag{Version: tagVersion, Hash: i.Hash(tag.LatestCommit)})
		}
		isLastPage = tags.IsLastPage
		start = tags.NextPageStart
	}
	return tagsResult, nil
}

func (p *Provider) tags(pkg i.Pkg, project string, start int) (TagsResponse, error) {
	var tags TagsResponse
	url := fmt.Sprintf("%s/rest/api/1.0/projects/%s/repos/%s/tags?start=%d", p.apiUrl, project, p.repoName(pkg), start)
	reposResponse, err := p.httpClient.Get(url)
	if err != nil {
		return tags, err
	}
	err = json.Unmarshal(reposResponse, &tags)

	return tags, err
}

func (p *Provider) repoName(pkg i.Pkg) string {
	strpkg := string(pkg)
	pos := strings.Index(strpkg, "/")
	if pos == -1 {
		return strpkg
	}
	return strpkg[pos+1:]
}
