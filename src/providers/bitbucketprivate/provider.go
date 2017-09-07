package bitbucketprivate

import (
	"encoding/json"
	"fmt"
	"strings"

	i "github.com/stamm/dep_radar/interfaces"
)

var (
	_ i.ITagGetter = &BitBucketPrivate{}
	_ i.IProvider  = &BitBucketPrivate{}
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

type BitBucketPrivate struct {
	project    string
	httpClient i.IWebClient
	gitDomain  string
	goGetUrl   string
	apiUrl     string
}

type Options struct {
	URL      string
	User     string
	Password string
}

func New(httpClient i.IWebClient, gitDomain, goGetUrl, apiUrl string) *BitBucketPrivate {
	return &BitBucketPrivate{
		httpClient: httpClient,
		goGetUrl:   goGetUrl,
		gitDomain:  gitDomain,
		apiUrl:     apiUrl,
	}
}

func (a *BitBucketPrivate) File(pkg i.Pkg, name string) ([]byte, error) {
	project, err := a.getProject(pkg)
	if err != nil {
		return nil, err
	}
	url := fmt.Sprintf("%s/projects/%s/repos/%s/raw/%s", a.apiUrl, project, a.repoName(pkg), name)
	return a.httpClient.Get(url)
}

func (a *BitBucketPrivate) GoGetUrl() string {
	return a.goGetUrl
}

// todo cache result in map
func (a *BitBucketPrivate) getProject(pkg i.Pkg) (string, error) {
	// if a.project != "" {
	// 	return nil
	// }
	project, err := GetProject(a.httpClient, pkg, a.gitDomain)
	if err != nil {
		return "", err
	}
	// fmt.Printf("project = %+v\n", project)
	// a.project = project
	return project, nil
}

// Tags get tags from bitbucket
func (a *BitBucketPrivate) Tags(pkg i.Pkg) ([]i.Tag, error) {
	project, err := a.getProject(pkg)
	if err != nil {
		return nil, err
	}
	var (
		tagsResult []i.Tag
		start      int
		isLastPage bool
	)
	for !isLastPage {
		tags, err := a.tags(pkg, project, start)
		if err != nil {
			return tagsResult, fmt.Errorf("Error on getting tags %s\n", err)
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

func (a *BitBucketPrivate) tags(pkg i.Pkg, project string, start int) (TagsResponse, error) {
	var tags TagsResponse
	url := fmt.Sprintf("%s/rest/api/1.0/projects/%s/repos/%s/tags?start=%d", a.apiUrl, project, a.repoName(pkg), start)
	reposResponse, err := a.httpClient.Get(url)
	if err != nil {
		return tags, err
	}
	err = json.Unmarshal(reposResponse, &tags)

	return tags, err
}

func (a *BitBucketPrivate) repoName(pkg i.Pkg) string {
	strpkg := string(pkg)
	pos := strings.Index(strpkg, "/")
	if pos == -1 {
		return strpkg
	}
	return strpkg[pos+1:]
}
