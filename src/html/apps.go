package html

import (
	"fmt"
	"sort"
	"time"

	i "github.com/stamm/dep_radar/interfaces"
	"github.com/stamm/dep_radar/src/fill"
	"github.com/stamm/dep_radar/src/providers"
	versionpkg "github.com/stamm/dep_radar/src/version"
)

const (
	colorMandatoryBad = "#F60018"
	colorMandatoryOk  = "#80EA69"
	colorNeedVersion  = "#FFAD73"
	colorExclude      = "#FB717E"
	// color= ""
	// color= ""
)

type templateStruct struct {
	SortedPkgsKeys []string
	AppsMap        i.AppListWithDeps
	LibsMap        i.LibMapWithTags
}

type MapRecomended map[i.Pkg]Option

type Option struct {
	Recomended  string
	Mandatory   bool
	Exclude     bool
	NeedVersion bool
	UpdatedAt   time.Time
	UpdateAfter time.Time
}

func AppsHtml(apps []i.IApp, detector *providers.Detector, rec MapRecomended) (string, error) {
	appsMap, libsMap := fill.GetTags(apps, detector)
	// fmt.Printf("libsMap = %+v\n", libsMap)
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
	for pkg := range rec {
		if _, ok := libsMap[pkg]; !ok {
			pkgsKey = append(pkgsKey, string(pkg))
		}
	}

	result := "<table>"
	result += "<thead><tr><th>apps</th>"
	for pkg := range appsMap {
		result += fmt.Sprintf("<th>%s</th>", pkg)
	}
	result += "</thead></tr><tbody>"
	for _, pkg := range pkgsKey {
		recomended := ""
		if recLib, ok := rec[i.Pkg(pkg)]; ok {
			if recLib.Recomended != "" {
				recomended = ` (<font style="background-color: ` + colorMandatoryOk + `">` + recLib.Recomended + "</font>)"
			}
		}
		result += fmt.Sprintf("<tr><td>%s%s</td>", pkg, recomended)
		for _, libs := range appsMap {
			libPkg := i.Pkg(pkg)
			opt, okOpt := rec[libPkg]

			dep, ok := libs[libPkg]
			if !ok {
				extra := ""
				if okOpt {
					if opt.Mandatory {
						extra = ` bgcolor="` + colorMandatoryBad + `"`
					}
					if opt.Exclude {
						extra = ` bgcolor="` + colorMandatoryOk + `"`
					}
				}
				result += fmt.Sprintf("<td%s>--</td>", extra)
				continue
			}
			version := dep.Version
			if version == "" {
				version = string(libs[libPkg].Hash)[0:8]
			}
			extra := ""
			if okOpt {
				if opt.Exclude {
					extra = ` bgcolor="` + colorExclude + `"`
				} else {
					goodVersion, _ := versionpkg.Compare(opt.Recomended, dep.Version)
					if goodVersion {
						extra = ` bgcolor="` + colorMandatoryOk + `"`
					} else {
						extra = ` bgcolor="` + colorMandatoryBad + `"`
					}
				}
			}
			if okOpt && opt.NeedVersion && dep.Version == "" {
				extra = ` bgcolor="` + colorNeedVersion + `"`
			}
			result += fmt.Sprintf("<td%s>%s</td>", extra, version)
		}
		result += "</tr>"
	}
	result += "</tbody></table>"
	result += `<br/><br/>Legend:<br/>
	<font style="background-color: ` + colorMandatoryBad + `">Don't fit mandatory</font><br/>
	<font style="background-color: ` + colorMandatoryOk + `">Fit mandatory</font><br/>
	<font style="background-color: ` + colorNeedVersion + `">Need version</font><br/>
	<font style="background-color: ` + colorExclude + `">Need to exclude</font><br/>
	<br/>`
	return "<body>" + result + "</body>", nil
}
