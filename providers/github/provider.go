package github

import (
	"context"
	"encoding/json"
	"fmt"
	urlpkg "net/url"
	"regexp"
	"strings"

	"github.com/stamm/dep_radar"
	"github.com/stamm/dep_radar/http"
)

const (
	// Prefix for github
	Prefix = "github.com"
	apiURL = "https://api.github.com/"
)

var (
	_ dep_radar.ITagGetter = &Provider{}
	_ dep_radar.IProvider  = &Provider{}
	_ dep_radar.IWebClient = &HTTPWrapper{}
)

type tag struct {
	Ref    string `json:"name"`
	Commit commit `json:"commit"`
}

type commit struct {
	SHA string `json:"sha"`
}

// Provider for github
type Provider struct {
	token  string
	client dep_radar.IWebClient
}

// HTTPWrapper for github
type HTTPWrapper struct {
	token  string
	client dep_radar.IWebClient
}

// NewHTTPWrapper returns http wrapper
func NewHTTPWrapper(token string, limit int) *HTTPWrapper {
	return &HTTPWrapper{
		token:  token,
		client: http.NewClient(http.Options{}, limit),
	}
}

// Get returns response for get url
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

// New creates new instance of provider
func New(client dep_radar.IWebClient) *Provider {
	return &Provider{
		client: client,
	}
}

// Tags get tags from github
func (g Provider) Tags(ctx context.Context, pkg dep_radar.Pkg) ([]dep_radar.Tag, error) {
	return g.tagsFromGithub(ctx, pkg)
}

// File gets file from github
func (g Provider) File(ctx context.Context, pkg dep_radar.Pkg, branch, name string) ([]byte, error) {
	return g.client.Get(ctx, g.makeURL(pkg, branch, name))
}

func (g Provider) makeURL(pkg dep_radar.Pkg, branch, name string) string {
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

// GoGetURL gets url for go get
func (g Provider) GoGetURL() string {
	return Prefix
}

func (g Provider) tagsFromGithub(ctx context.Context, pkg dep_radar.Pkg) ([]dep_radar.Tag, error) {
	url := makeAPIURL("repos/" + getPkgName(pkg) + "/tags?per_page=100")
	content, err := g.client.Get(ctx, url)
	if err != nil {
		return nil, err
	}

	var tags []tag
	err = json.Unmarshal(content, &tags)
	if err != nil {
		return nil, err
	}

	result := make([]dep_radar.Tag, 0, len(tags))
	for _, t := range tags {
		result = append(result, dep_radar.Tag{
			Version: t.Ref,
			Hash:    dep_radar.Hash(t.Commit.SHA),
		})
	}

	return result, nil
}

func getPkgName(pkg dep_radar.Pkg) string {
	name := strings.Trim(string(pkg), "/")
	re := regexp.MustCompile("^github\\.com/")
	return re.ReplaceAllString(name, "")
}

func makePkgName(fullname string) dep_radar.Pkg {
	return dep_radar.Pkg(Prefix + "/" + fullname)
}

func makeAPIURL(uri string) string {
	return fmt.Sprintf("%s%s", apiURL, uri)
}
