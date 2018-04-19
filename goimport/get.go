package goimport

import (
	"bytes"
	"context"

	"github.com/stamm/dep_radar"
)

// MetaImport represents the parsed <meta name="go-import"
// content="prefix vcs reporoot" /> tags from HTML files.
type MetaImport struct {
	Prefix, VCS, RepoRoot string
}

// GetImports returns meta imports
func GetImports(ctx context.Context, client dep_radar.IWebClient, url string) ([]MetaImport, error) {
	fullURL := "https://" + url + "?go-get=1"
	resp, err := client.Get(ctx, fullURL)
	if err != nil {
		return nil, err
	}
	buf := bytes.NewBuffer(resp)
	imports, err := parseMetaGoImports(buf)
	if err != nil {
		return nil, err
	}
	return imports, nil
}
