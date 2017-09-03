package html

import (
	"fmt"
	"sort"

	i "github.com/stamm/dep_radar/interfaces"
	"github.com/stamm/dep_radar/src/fill"
	"github.com/stamm/dep_radar/src/providers"
)

func LibsHtml(apps []i.IApp, detector *providers.Detector) (string, error) {
	appsMap, libsMap := fill.GetTags(apps, detector)
	appsMap = fill.AddVersionLibToApp(appsMap, libsMap)

	pkgsKey := make([]string, 0, len(libsMap))
	for pkgKey := range libsMap {
		pkgsKey = append(pkgsKey, string(pkgKey))
	}
	sort.Strings(pkgsKey) //sort by key

	// TODO: need to fix
	result := "<table>"
	result += "<tr><td>apps</td>"
	for pkg := range libsMap {
		result += fmt.Sprintf("<td>%s</td>", pkg)
	}
	result += "</tr>"
	for _, pkg := range pkgsKey {
		result += fmt.Sprintf("<tr><td>%s</td>", pkg)
		for _, libs := range appsMap {
			dep, ok := libs[i.Pkg(pkg)]
			if !ok {
				result += "<td>-</td>"
				continue
			}
			version := dep.Version
			if version == "" {
				version = string(libs[i.Pkg(pkg)].Hash)
			}
			result += fmt.Sprintf("<td>%s</td>", version)
		}
		result += "</tr>"
	}
	result += "</table>"
	return "<body>" + result + "</body>", nil
}
