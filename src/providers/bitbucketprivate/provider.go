package bitbucketprivate

import (
	"encoding/json"
	"fmt"
	"strings"

	i "github.com/stamm/dep_radar/interfaces"
)

var (
	_ i.IApp       = &BitBucketPrivate{}
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
	pkg        i.Pkg
	httpClient i.IWebClient
	gitDomain  string
}

type Options struct {
	URL      string
	User     string
	Password string
}

func New(pkg i.Pkg, httpClient i.IWebClient, gitDomain string) (*BitBucketPrivate, error) {
	return &BitBucketPrivate{
		pkg:        pkg,
		httpClient: httpClient,
		gitDomain:  gitDomain,
	}, nil
}

func (a *BitBucketPrivate) Package() i.Pkg {
	return a.pkg
}

func (a *BitBucketPrivate) File(name string) ([]byte, error) {
	a.setProject()
	url := fmt.Sprintf("projects/%s/repos/%s/raw/%s", a.project, a.repoName(), name)
	fmt.Printf("url = %+v\n", url)
	return a.httpClient.Get(url)
}

func (a *BitBucketPrivate) setProject() {
	if a.project != "" {
		return
	}
	project, err := GetProject(a, a.gitDomain)
	if err != nil {
		return
	}
	// fmt.Printf("project = %+v\n", project)
	a.project = project
}

// Tags get tags from bitbucket
func (a *BitBucketPrivate) Tags() ([]i.Tag, error) {
	a.setProject()
	var (
		tagsResult []i.Tag
		start      int
		isLastPage bool
	)
	for !isLastPage {
		tags, err := a.tags(start)
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

func (a *BitBucketPrivate) tags(start int) (TagsResponse, error) {
	var tags TagsResponse
	url := fmt.Sprintf("rest/api/1.0/projects/%s/repos/%s/tags?start=%d", a.project, a.repoName(), start)
	reposResponse, err := a.httpClient.Get(url)
	if err != nil {
		return tags, err
	}
	err = json.Unmarshal(reposResponse, &tags)

	return tags, err
}

func (a *BitBucketPrivate) repoName() string {
	strpkg := string(a.pkg)
	pos := strings.Index(strpkg, "/")
	if pos == -1 {
		return strpkg
	}
	return strpkg[pos+1:]
}
