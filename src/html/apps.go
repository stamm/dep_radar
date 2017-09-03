package html

import (
	"fmt"
	"sort"

	i "github.com/stamm/dep_radar/interfaces"
	"github.com/stamm/dep_radar/src/fill"
	"github.com/stamm/dep_radar/src/providers"
)

type templateStruct struct {
	SortedPkgsKeys []string
	AppsMap        i.AppListWithDeps
	LibsMap        i.LibMapWithTags
}

func AppsHtml(apps []i.IApp, detector *providers.Detector) (string, error) {
	appsMap, libsMap := fill.GetTags(apps, detector)
	appsMap = fill.AddVersionLibToApp(appsMap, libsMap)

	pkgsKey := make([]string, 0, len(libsMap))
	for pkgKey := range libsMap {
		pkgsKey = append(pkgsKey, string(pkgKey))
	}
	sort.Strings(pkgsKey) //sort by key
	// fmt.Printf("pkgsKey = %+v\n", pkgsKey)

	// tmpl, err := template.ParseFiles("src/html/apps.html")
	// if err != nil {
	// 	return "", err
	// }
	// data := templateStruct{
	// 	SortedPkgsKeys: pkgsKey,
	// 	AppsMap:        appsMap,
	// 	LibsMap:        libsMap,
	// }
	// var buf bytes.Buffer
	// err = tmpl.Execute(&buf, data)
	// if err != nil {
	// 	return "", err
	// }
	// html, err := ioutil.ReadAll(&buf)
	// if err != nil {
	// 	return "", err
	// }
	// return string(html), nil

	result := "<table>"
	result += "<tr><td>apps</td>"
	for pkg := range appsMap {
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
