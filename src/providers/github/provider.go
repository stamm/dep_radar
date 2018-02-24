package github

import (
	"context"
	"encoding/json"
	"fmt"
	urlpkg "net/url"
	"regexp"
	"strings"

	i "github.com/stamm/dep_radar/interfaces"
	"github.com/stamm/dep_radar/src/http"
)

const (
	// Prefix for github
	Prefix = "github.com"
)

var (
	_ i.ITagGetter = &Provider{}
	_ i.IProvider  = &Provider{}
	_ i.IWebClient = &HTTPWrapper{}
)

type tag struct {
	Ref    string `json:"name"`
	Commit commit `json:"commit"`
}

type commit struct {
	SHA string `json:"sha"`
}

type Provider struct {
	token  string
	client i.IWebClient
}

type HTTPWrapper struct {
	token  string
	client i.IWebClient
}

func NewHTTPWrapper(token string, limit int) *HTTPWrapper {
	return &HTTPWrapper{
		token:  token,
		client: http.NewClient(http.Options{}, limit),
	}
}

func (c *HTTPWrapper) Get(ctx context.Context, url string) ([]byte, error) {
	newURL, err := c.getURL(url)
	if err != nil {
		return []byte{}, err
	}
	return c.client.Get(ctx, newURL)
}

func (c *HTTPWrapper) getURL(url string) (string, error) {
	if c.token == "" {
		return url, nil
	}
	urlObj, err := urlpkg.Parse(url)
	if err != nil {
		return "", err
	}
	if len(urlObj.Query()) == 0 {
		url += "?access_token=" + c.token
	} else {
		url += "&access_token=" + c.token
	}
	return url, nil
}

func New(client i.IWebClient) *Provider {
	return &Provider{
		client: client,
	}
}

func (g Provider) Tags(ctx context.Context, pkg i.Pkg) ([]i.Tag, error) {
	return g.tagsHttp(ctx, pkg)
}

func (g Provider) File(ctx context.Context, pkg i.Pkg, branch, name string) ([]byte, error) {
	return g.client.Get(ctx, g.makeURL(pkg, branch, name))
}

func (g Provider) makeURL(pkg i.Pkg, branch, name string) string {
	pkgName := strings.Trim(string(pkg), "/")
	re := regexp.MustCompile("^" + regexp.QuoteMeta(Prefix) + "/")
	repo := re.ReplaceAllString(pkgName, "")
	parts := strings.SplitN(repo, "/", 3)
	url := fmt.Sprintf("%s/%s/", strings.Join(parts[:2], "/"), branch)
	if len(parts) > 2 {
		url += parts[2] + "/"
	}
	url += strings.Trim(name, "/")
	return "https://raw.githubusercontent.com/" + url
}

func (g Provider) GoGetUrl() string {
	return Prefix
}

func (g Provider) tagsHttp(ctx context.Context, pkg i.Pkg) ([]i.Tag, error) {
	url := "https://api.github.com/repos/" + getPkgName(pkg) + "/tags?per_page=100"
	content, err := g.client.Get(ctx, url)
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
