package goimport

import (
	"bytes"
	"context"

	i "github.com/stamm/dep_radar/interfaces"
)

// metaImport represents the parsed <meta name="go-import"
// content="prefix vcs reporoot" /> tags from HTML files.
type metaImport struct {
	Prefix, VCS, RepoRoot string
}

type metaSource struct {
	Prefix, Main, Dir, File string
}

// GetImports returns meta imports
func GetImports(ctx context.Context, client i.IWebClient, url string) ([]metaImport, error) {
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
