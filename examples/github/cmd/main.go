package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"time"

	"github.com/stamm/dep_radar/examples/github/custom"
	i "github.com/stamm/dep_radar/interfaces"
	"github.com/stamm/dep_radar/src/app"
	"github.com/stamm/dep_radar/src/deps"
	"github.com/stamm/dep_radar/src/html"
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
	provDetector := custom.Detector()
	depDetector := deps.DefaultDetector()

	// pkgs := []i.Pkg{"github.com/dep-radar/test_app"}

	githubProv := github.New(github.NewHTTPWrapper(os.Getenv("GITHUB_TOKEN"), 10))
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
	}
	htmlResult, err := html.AppsHtml(apps, provDetector, mapRec)
	// htmlResult, err := html.LibsHtml(apps, provDetector)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprint(w, htmlResult)
	fmt.Fprintf(w, "took %s", time.Since(start))
}
