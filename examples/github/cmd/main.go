package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
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
	log.SetFlags(log.Lmicroseconds)
	http.HandleFunc("/", handler)
	http.HandleFunc("/favicon.ico", http.NotFound)
	log.Println("Started: http://localhost:8081/")
	http.ListenAndServe(":8081", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
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
	start := time.Now()
	githubProv := github.New(github.NewHTTPWrapper(os.Getenv("GITHUB_TOKEN"), 10))
	provDetector := providers.NewDetector().AddProvider(githubProv)
	depDetector := deps.DefaultDetector()

	apps := make(chan i.IApp, 10)
	go func() {
		pkgs, err := githubProv.GetAllOrgRepos(context.Background(), "dep-radar")
		if err != nil {
			log.Fatal(err)
		}
		for _, pkg := range pkgs {
			apiApp, err := app.New(pkg, "master", provDetector, depDetector)
			if err != nil {
				log.Printf("cant create app %s, got err: %s\n", pkg, err)
			}
			apps <- apiApp
		}
		close(apps)
	}()

	htmlResult, err := html.AppsHTML(apps, provDetector, mapRec)
	// htmlResult, err := html.LibsHTML(apps, provDetector, mapRec)
	if err != nil {
		fmt.Fprintf(w, "Error: %s", err)
	} else {
		w.Write(htmlResult)
	}
	fmt.Fprintf(w, "took %s", time.Since(start))
}
