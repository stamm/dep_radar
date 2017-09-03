package bitbucketprivate

import (
	"regexp"
	"strings"

	i "github.com/stamm/dep_radar/interfaces"
	"github.com/stamm/dep_radar/src/goimport"
)

func GetProject(app i.IApp, prefix string) (string, error) {
	prefix = strings.Trim(prefix, "/")
	prefix = regexp.QuoteMeta(prefix)
	re := regexp.MustCompile(prefix + `/([^/]+)`)
	url := string(app.Package())
	sources, err := goimport.GetImports(url)
	if err != nil {
		return "", err
	}
	for _, source := range sources {
		matched := re.FindAllStringSubmatch(source.RepoRoot, -1)
		if len(matched) > 0 {
			return matched[0][1], nil
		}
	}
	return "", nil
}
