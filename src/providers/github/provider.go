package github

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strings"

	i "github.com/stamm/dep_radar/interfaces"
	"github.com/stamm/dep_radar/src/http"
)

const (
	Prefix = "github.com"
)

var (
	_ i.ITagGetter = &Github{}
	_ i.IApp       = &Github{}
	_ i.IProvider  = &Github{}
	_ i.IWebClient = &HTTPWrapper{}
)

type tag struct {
	Ref    string `json:"name"`
	Commit commit `json:"commit"`
}

type commit struct {
	SHA string `json:"sha"`
}

type Github struct {
	pkg    i.Pkg
	token  string
	client i.IWebClient
}

type HTTPWrapper struct {
	token  string
	client i.IWebClient
}

func NewHTTPWrapper(token string, limit int) *HTTPWrapper {
	return &HTTPWrapper{
		token: token,
		client: http.NewClient(http.Options{
			URL: "https://raw.githubusercontent.com/",
		}, 10),
	}
}

func (c *HTTPWrapper) Get(url string) ([]byte, error) {
	if c.token != "" {
		url += "?access_token=" + c.token
	}
	return c.client.Get(url)
}

func New(pkg i.Pkg, client i.IWebClient) (Github, error) {
	if strings.Index(string(pkg), Prefix) != 0 {
		return Github{}, fmt.Errorf("package %s is not for github", pkg)
	}
	return Github{
		pkg:    pkg,
		client: client,
	}, nil
}

func (g Github) Tags() ([]i.Tag, error) {
	return g.tagsHttp()
}

func (g Github) Package() i.Pkg {
	return g.pkg
}

func (g Github) File(name string) ([]byte, error) {
	return g.client.Get(g.makeURL(name))
}

func (g Github) makeURL(name string) string {
	pkg := strings.Trim(string(g.pkg), "/")
	re := regexp.MustCompile("^" + regexp.QuoteMeta(Prefix) + "/")
	repo := re.ReplaceAllString(pkg, "")
	parts := strings.SplitN(repo, "/", 3)
	url := strings.Join(parts[:2], "/") + "/master/"
	if len(parts) > 2 {
		url += parts[2] + "/"
	}
	url += strings.Trim(name, "/")
	return url
}

func (g Github) tagsHttp() ([]i.Tag, error) {
	url := "https://api.github.com/repos/" + getPkgName(g.pkg) + "/tags"
	content, err := g.client.Get(url)
	if err != nil {
		return nil, err
	}

	// fmt.Printf("string(content) = %+v\n", string(content))

	var tags []tag
	err = json.Unmarshal(content, &tags)
	if err != nil {
		return nil, err
	}

	result := make([]i.Tag, 0, len(tags))
	for _, t := range tags {
		result = append(result, i.Tag{
			Version: t.Ref,
			Hash:    i.Hash(t.Commit.SHA),
		})
	}

	return result, nil
}

func getPkgName(pkg i.Pkg) string {
	name := strings.Trim(string(pkg), "/")
	re := regexp.MustCompile("^github\\.com/")
	return re.ReplaceAllString(name, "")
}
