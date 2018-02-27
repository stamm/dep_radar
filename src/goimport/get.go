package goimport

import (
	"bytes"
	"context"

	i "github.com/stamm/dep_radar/src/interfaces"
)

// MetaImport represents the parsed <meta name="go-import"
// content="prefix vcs reporoot" /> tags from HTML files.
type MetaImport struct {
	Prefix, VCS, RepoRoot string
}

// GetImports returns meta imports
func GetImports(ctx context.Context, client i.IWebClient, url string) ([]MetaImport, error) {
	fullURL := "https://" + url + "?go-get=1"
	// fmt.Printf("fullUrl = %+v\n", fullUrl)
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
