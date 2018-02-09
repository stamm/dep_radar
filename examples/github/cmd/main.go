package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"time"

	i "github.com/stamm/dep_radar/interfaces"
	"github.com/stamm/dep_radar/src/app"
	"github.com/stamm/dep_radar/src/deps"
	"github.com/stamm/dep_radar/src/html"
	"github.com/stamm/dep_radar/src/providers"
	"github.com/stamm/dep_radar/src/providers/github"
)

func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/favicon.ico", http.NotFound)
	log.Println("Started: http://localhost:8081/")
	http.ListenAndServe(":8081", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	githubProv := github.New(github.NewHTTPWrapper(os.Getenv("GITHUB_TOKEN"), 10))
	provDetector := providers.NewDetector().AddProvider(githubProv)
	depDetector := deps.DefaultDetector()

	// pkgs := []i.Pkg{"github.com/dep-radar/test_app"}

	pkgs, err := githubProv.GetAllOrgRepos(context.Background(), "dep-radar")
	if err != nil {
		log.Fatal(err)
	}

	// Create a little wrapper with custom logic for detect
	apps := make([]i.IApp, 0, len(pkgs))
	for _, pkg := range pkgs {
		apiApp, err := app.New(pkg, "master", provDetector, depDetector)
		if err != nil {
			log.Printf("cant create app %s, got err: %s\n", pkg, err)
		}
		apps = append(apps, apiApp)
	}

	mapRec := html.MapRecomended{
		"github.com/pkg/errors": html.Option{
			Recomended: ">=0.8.0",
			Mandatory:  true,
		},
		"github.com/pkg/sftp": html.Option{
			Recomended: ">=1.3.0",
		},
		"github.com/kr/fs": html.Option{
			Exclude: true,
		},
		"github.com/pkg/bson": html.Option{
			Mandatory: true,
		},
		"github.com/pkg/singlefile": html.Option{
			Exclude: true,
		},
		"github.com/pkg/profile": html.Option{
			NeedVersion: true,
		},
	}
	htmlResult, err := html.AppsHTML(apps, provDetector, mapRec)
	// htmlResult, err := html.LibsHtml(apps, provDetector)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprint(w, htmlResult)
	fmt.Fprintf(w, "took %s", time.Since(start))
}
