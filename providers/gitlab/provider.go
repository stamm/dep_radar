package gitlab

import (
	"context"
	"strings"

	"github.com/stamm/dep_radar"
	"github.com/xanzy/go-gitlab"
)

var (
	_ dep_radar.ITagGetter = &Provider{}
	_ dep_radar.IProvider  = &Provider{}
)

type Provider struct {
	goGetURL string

	client *gitlab.Client
}

// New creates new instance of provider
func New(client *gitlab.Client, goGetURL string) *Provider {
	return &Provider{
		client:   client,
		goGetURL: goGetURL,
	}
}

func (p *Provider) cleanID(pkg dep_radar.Pkg) string {
	pack := strings.TrimPrefix(string(pkg), p.goGetURL)
	pack = strings.Trim(pack, "/")
	return pack
}

// File gets file from github
func (p *Provider) File(ctx context.Context, pkg dep_radar.Pkg, branch, filename string) ([]byte, error) {
	opts := &gitlab.GetRawFileOptions{
		Ref: &branch,
	}
	data, _, err := p.client.RepositoryFiles.GetRawFile(p.cleanID(pkg), filename, opts)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// GoGetURL gets url for go get
func (p *Provider) GoGetURL() string {
	return p.goGetURL
}

// Tags get tags from github
func (p *Provider) Tags(ctx context.Context, pkg dep_radar.Pkg) ([]dep_radar.Tag, error) {
	opts := gitlab.ListTagsOptions{}
	tagsResp, _, err := p.client.Tags.ListTags(p.cleanID(pkg), &opts)
	if err != nil {
		return nil, err
	}

	result := make([]dep_radar.Tag, len(tagsResp))

	for idx, tag := range tagsResp {
		result[idx] = dep_radar.Tag{
			Version: tag.Name,
			Hash:    dep_radar.Hash(tag.Commit.ID),
		}
	}

	return result, nil
}
