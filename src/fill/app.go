package fill

import (
	"context"
	"log"
	"sync"

	i "github.com/stamm/dep_radar/interfaces"
	"github.com/stamm/dep_radar/src/providers"
)

// Lib struct for lib
type Lib struct {
	Dep    i.Dep
	Pkg    i.Pkg
	AppPkg i.Pkg
}

// GetTags return list of apps with dependencies and libraries with tags
func GetTags(ctx context.Context, apps <-chan i.IApp, detector *providers.Detector) (i.AppListWithDeps, i.LibMapWithTags) {
	var (
		muRes sync.RWMutex
		wg    sync.WaitGroup
	)
	libsMap := make(map[i.Pkg]struct{})
	res := make(map[i.Pkg][]i.Tag)
	appList := make(i.AppListWithDeps, 100)
	for lib := range depsChan(ctx, apps) {
		if _, ok := appList[lib.AppPkg]; !ok {
			appList[lib.AppPkg] = make(map[i.Pkg]i.Dep)
		}
		appList[lib.AppPkg][lib.Pkg] = lib.Dep
		muRes.RLock()
		_, ok := libsMap[lib.Pkg]
		muRes.RUnlock()
		if !ok {
			wg.Add(1)
			libsMap[lib.Pkg] = struct{}{}
			// go get list of tags
			go func(libPkg i.Pkg) {
				defer wg.Done()
				tagList, err := getTagsForLib(ctx, libPkg, detector)
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

func depsChan(ctx context.Context, apps <-chan i.IApp) chan Lib {
	ch := make(chan Lib, 100)
	go func() {
		var wg sync.WaitGroup
		for app := range apps {
			wg.Add(1)
			go func(app i.IApp) {
				defer wg.Done()
				deps, err := app.Deps(ctx)
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

func getTagsForLib(ctx context.Context, pkg i.Pkg, detector *providers.Detector) ([]i.Tag, error) {
	tagsGetter, err := detector.Detect(ctx, pkg)
	if err != nil {
		if err != providers.ErrNoProvider {
			log.Printf("err for pkg %q from route: %s", pkg, err)
		}
		return nil, err
	}
	tagList, err := tagsGetter.Tags(ctx, pkg)
	if err != nil {
		log.Printf("err for pkg %q from tag getter: %s", pkg, err)
		return nil, err
	}
	return tagList, nil
}

// AddVersionLibToApp set version of libs inside each app
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
