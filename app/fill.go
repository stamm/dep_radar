package app

import (
	"context"
	"fmt"
	"log"
	"sync"

	"github.com/stamm/dep_radar"
	"github.com/stamm/dep_radar/providers"
)

// lib struct for lib.
type lib struct {
	Dep    dep_radar.Dep
	Pkg    dep_radar.Pkg
	AppPkg dep_radar.Pkg
}

// GetTags return list of apps with dependencies and libraries with tags.
func GetTags(ctx context.Context, appsCh <-chan dep_radar.IApp, detector *providers.Detector) (dep_radar.AppListWithDeps, dep_radar.LibMapWithTags) {
	var (
		muRes sync.RWMutex
		wg    sync.WaitGroup
	)
	libExists := make(map[dep_radar.Pkg]struct{})
	res := make(map[dep_radar.Pkg][]dep_radar.Tag)
	apps := make(dep_radar.AppListWithDeps, 100)

	for lib := range extractLibs(ctx, appsCh) {
		if _, ok := apps[lib.AppPkg]; !ok {
			apps[lib.AppPkg] = make(map[dep_radar.Pkg]dep_radar.Dep)
		}
		apps[lib.AppPkg][lib.Pkg] = lib.Dep
		if _, ok := libExists[lib.Pkg]; ok {
			continue
		}
		libExists[lib.Pkg] = struct{}{}
		wg.Add(1)
		// go get list of tags
		go func(libPkg dep_radar.Pkg) {
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
	wg.Wait()
	return apps, res
}

// extractLibs get libs from channel of apps
func extractLibs(ctx context.Context, apps <-chan dep_radar.IApp) chan lib {
	libs := make(chan lib, 100)
	go func() {
		var wg sync.WaitGroup
		for app := range apps {
			wg.Add(1)
			go func(app dep_radar.IApp) {
				defer wg.Done()
				deps, err := app.Deps(ctx)
				if err != nil {
					log.Printf("err for getting deps for app %s: %+v\n", app.Package(), err)
					return
				}
				for pkg, dep := range deps {
					libs <- lib{
						Pkg:    pkg,
						Dep:    dep,
						AppPkg: app.Package(),
					}
				}
			}(app)
		}
		wg.Wait()
		close(libs)
	}()
	return libs
}

func getTagsForLib(ctx context.Context, pkg dep_radar.Pkg, detector *providers.Detector) ([]dep_radar.Tag, error) {
	tagsGetter, err := detector.Detect(ctx, pkg)
	if err != nil {
		fmt.Printf("err = %+v\n", err)
		if err != providers.ErrNoProvider {
			log.Printf("err for pkg %q from route: %s", pkg, err)
		}
		return nil, err
	}
	tagList, err := tagsGetter.Tags(ctx, pkg)
	if err != nil {
		return nil, err
	}
	return tagList, nil
}

// AddVersionLibToApp set version of libs inside each app
func AddVersionLibToApp(apps dep_radar.AppListWithDeps, libs dep_radar.LibMapWithTags) dep_radar.AppListWithDeps {
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
