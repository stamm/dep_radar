package html

import (
	"fmt"
	"sort"
	"strings"
	"time"

	i "github.com/stamm/dep_radar/interfaces"
	"github.com/stamm/dep_radar/src/fill"
	"github.com/stamm/dep_radar/src/providers"
	versionpkg "github.com/stamm/dep_radar/src/version"
)

const (
	classMandatoryOk  = "table-success"
	classMandatoryBad = "bg-danger"
	classExclude      = "table-warning"
	classNeedVersion  = "table-warning"
	// color= ""
	// color= ""
)

type templateStruct struct {
	SortedPkgsKeys []string
	AppsMap        i.AppListWithDeps
	LibsMap        i.LibMapWithTags
}

// MapRecomended map with option for libraries
type MapRecomended map[i.Pkg]Option

// Option for library
type Option struct {
	Recomended  string
	Mandatory   bool
	Exclude     bool
	NeedVersion bool
	UpdatedAt   time.Time
	UpdateAfter time.Time
}

// AppsHTML return html with table. In the head apps, on the left side - libs
func AppsHTML(apps []i.IApp, detector *providers.Detector, rec MapRecomended) (string, error) {
	appsMap, libsMap := fill.GetTags(apps, detector)
	// fmt.Printf("libsMap = %+v\n", libsMap)
	appsMap = fill.AddVersionLibToApp(appsMap, libsMap)

	pkgsKey := make(sort.StringSlice, 0, len(libsMap))
	for pkgKey := range libsMap {
		pkgsKey = append(pkgsKey, string(pkgKey))
	}
	for pkg := range rec {
		if _, ok := libsMap[pkg]; !ok {
			pkgsKey = append(pkgsKey, string(pkg))
		}
	}
	sort.Strings(pkgsKey) //sort by key
	appsKey := make(sort.StringSlice, 0, len(apps))
	for _, appPkg := range apps {
		appsKey = append(appsKey, string(appPkg.Package()))
	}
	sort.Strings(appsKey)
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

	result := `<table class="table table-striped table-hover table-sm">`
	result += `<thead class="thead-dark"><tr><th>apps</th>`
	for _, pkg := range appsKey {
		result += fmt.Sprintf("<th>%s</th>", pkg)
	}
	result += "</thead></tr><tbody>"
	for _, pkg := range pkgsKey {
		var hints []string
		if recLib, ok := rec[i.Pkg(pkg)]; ok {
			if recLib.Recomended != "" {
				hints = append(hints, `<font style="background-color: `+colorMandatoryOk+`">`+recLib.Recomended+"</font>")
			}
			if recLib.Mandatory {
				hints = append(hints, "mandatory")
			}
			if recLib.Exclude {
				hints = append(hints, "exclude")
			}
			if recLib.NeedVersion {
				hints = append(hints, "only version")
			}
		}
		result += fmt.Sprintf("<tr><td><b>%s</b> <small>%s</small></td>", pkg, strings.Join(hints, ", "))
		for _, appKey := range appsKey {
			libs := appsMap[i.Pkg(appKey)]
			libPkg := i.Pkg(pkg)
			opt, okOpt := rec[libPkg]

			dep, ok := libs[libPkg]
			if !ok {
				extra := ""
				if okOpt {
					if opt.Mandatory {
						extra = fmt.Sprintf(` class="%s" data-toggle="tooltip" title="%s"`, classMandatoryBad, "Mandatory to use")
					}
					if opt.Exclude {
						extra = fmt.Sprintf(` class="%s"`, classMandatoryOk)
					}
				}
				result += fmt.Sprintf("<td%s>—</td>", extra)
				continue
			}
			version := dep.Version
			if version == "" {
				version = string(libs[libPkg].Hash)[0:8]
			}
			extra := ""
			if okOpt {
				if opt.Exclude {
					extra = fmt.Sprintf(` class="%s" data-toggle="tooltip" title="%s"`, classExclude, "Need to exclude this library")
				} else {
					goodVersion, _ := versionpkg.Compare(opt.Recomended, dep.Version)
					if goodVersion {
						extra = fmt.Sprintf(` class="%s"`, classMandatoryOk)
					} else {
						extra = fmt.Sprintf(` class="%s" data-toggle="tooltip" title="%s"`, classMandatoryBad, "Need to change version")
					}
				}
			}
			if okOpt && opt.NeedVersion && dep.Version == "" {
				extra = fmt.Sprintf(` class="%s" data-toggle="tooltip" title="%s"`, classNeedVersion, "Use version instead of revision")
			}
			result += fmt.Sprintf("<td%s>%s</td>", extra, version)
		}
		result += "</tr>"
	}
	result += "</tbody></table>"
	return `<head>
	<link rel="stylesheet" href="https://maxcdn.bootstrapcdn.com/bootstrap/4.0.0/css/bootstrap.min.css" integrity="sha384-Gn5384xqQ1aoWXA+058RXPxPg6fy4IWvTNh0E263XmFcJlSAwiGgFAW/dAiS6JXm" crossorigin="anonymous"></head><body>` + result + `
<script src="https://code.jquery.com/jquery-3.2.1.slim.min.js" integrity="sha384-KJ3o2DKtIkvYIK3UENzmM7KCkRr/rE9/Qpg6aAZGJwFDMVNA/GpGFF93hXpG5KkN" crossorigin="anonymous"></script>
<script src="https://cdnjs.cloudflare.com/ajax/libs/popper.js/1.12.9/umd/popper.min.js" integrity="sha384-ApNbgh9B+Y1QKtv3Rn7W3mgPxhU9K/ScQsAP7hUibX39j7fakFPskvXusvfa0b4Q" crossorigin="anonymous"></script>
<script src="https://maxcdn.bootstrapcdn.com/bootstrap/4.0.0/js/bootstrap.min.js" integrity="sha384-JZR6Spejh4U02d8jOt6vLEHfe/JQGiRRSQQxSfFWpi1MquVdAyjUar5+76PVCmYl" crossorigin="anonymous"></script>
<script>
$(function () {
  $('[data-toggle="tooltip"]').tooltip()
})
</script>

</body>`, nil
}
