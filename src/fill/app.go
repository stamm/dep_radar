package fill

import (
	"log"
	_ "net/http/pprof"
	"sync"

	i "github.com/stamm/dep_radar/interfaces"
	"github.com/stamm/dep_radar/src/deps"
	"github.com/stamm/dep_radar/src/providers"
)

type Lib struct {
	Dep    i.Dep
	Pkg    i.Pkg
	AppPkg i.Pkg
}

func GetTags(apps []i.IApp, detector *providers.Detector) (i.AppListWithDeps, i.LibMapWithTags) {
	var (
		muRes sync.Mutex
		wg    sync.WaitGroup
	)
	libsMap := make(map[i.Pkg]struct{})
	res := make(map[i.Pkg][]i.Tag)
	appList := make(i.AppListWithDeps, len(apps))
	for lib := range depsChan(apps) {
		// fmt.Printf("lib = %+v\n", lib)
		if _, ok := appList[lib.AppPkg]; !ok {
			appList[lib.AppPkg] = make(map[i.Pkg]i.Dep)
		}
		appList[lib.AppPkg][lib.Pkg] = lib.Dep

		if _, ok := libsMap[lib.Pkg]; !ok {
			wg.Add(1)
			libsMap[lib.Pkg] = struct{}{}
			// go get list of tags
			go func(libPkg i.Pkg) {
				defer wg.Done()
				tagList, err := getTagsForLib(libPkg, detector)
				if err != nil {
					if err != providers.ErrNoProvider {
						log.Printf("Error on getting tags for lib %s: %s", libPkg, err)
					}
					return
				}
				muRes.Lock()
				res[libPkg] = tagList
				muRes.Unlock()
			}(lib.Pkg)
		}
	}
	wg.Wait()
	return appList, res
}

func depsChan(apps []i.IApp) chan Lib {
	ch := make(chan Lib, 100)
	depTools := deps.DefaultTools()
	go func() {
		var wg sync.WaitGroup
		for _, app := range apps {
			wg.Add(1)
			go func(app i.IApp) {
				defer wg.Done()
				deps, err := depTools.Deps(app)
				if err != nil {
					log.Printf("err for app %s: %+v\n", app.Package(), err)
					return
				}
				for pkg, dep := range deps.Deps {
					ch <- Lib{
						Pkg:    pkg,
						Dep:    dep,
						AppPkg: app.Package(),
					}
				}
			}(app)
		}
		wg.Wait()
		close(ch)
	}()
	return ch
}

func getTagsForLib(pkg i.Pkg, detector *providers.Detector) ([]i.Tag, error) {
	tagsGetter, err := detector.Detect(pkg)
	if err != nil {
		if err != providers.ErrNoProvider {
			log.Printf("err for pkg %q from route: %s", pkg, err)
		}
		return nil, err
	}
	tagList, err := tagsGetter.Tags()
	if err != nil {
		log.Printf("err for pkg %q from tag getter: %s", pkg, err)
		return nil, err
	}
	return tagList, nil
}

func AddVersionLibToApp(apps i.AppListWithDeps, libs i.LibMapWithTags) i.AppListWithDeps {
	for appPkg, appLibs := range apps {
		for libPkg, appLib := range appLibs {
			if tags, ok := libs[libPkg]; ok {
				for _, tag := range tags {
					if tag.Hash == appLib.Hash {
						appLib.Version = tag.Version
						apps[appPkg][libPkg] = appLib
					}
				}
			}
		}
	}
	return apps
}
