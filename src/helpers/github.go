package helpers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	i "github.com/stamm/dep_radar/interfaces"
	"github.com/stamm/dep_radar/src"
	"github.com/stamm/dep_radar/src/app"
	"github.com/stamm/dep_radar/src/deps"
	"github.com/stamm/dep_radar/src/html"
	"github.com/stamm/dep_radar/src/providers"
	"github.com/stamm/dep_radar/src/providers/github"
)

// GithubOrg starts an http server for this organisation
func GithubOrg(token, orgName, listen string, recom src.MapRecomended) {
	log.SetFlags(log.Lmicroseconds)
	http.HandleFunc("/", wrapOrgHandler(token, orgName, recom))
	http.HandleFunc("/favicon.ico", http.NotFound)
	log.Println("Started: http://localhost:8081/")
	http.ListenAndServe(listen, nil)
}

// GithubPkgs starts an http server for particular list of applications
func GithubPkgs(token, listen string, recom src.MapRecomended, pkgs ...string) {
	log.SetFlags(log.Lmicroseconds)
	http.HandleFunc("/", wrapHandler(token, recom, pkgs...))
	http.HandleFunc("/favicon.ico", http.NotFound)
	log.Println("Started: http://localhost:8081/")
	http.ListenAndServe(listen, nil)
}

func wrapOrgHandler(token, orgName string, recom src.MapRecomended) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		githubProv := github.New(github.NewHTTPWrapper(token, 10))
		provDetector := providers.NewDetector().AddProvider(githubProv)
		depDetector := deps.DefaultDetector()

		apps := make(chan i.IApp, 10)
		go func() {
			pkgs, err := githubProv.GetAllOrgRepos(context.Background(), orgName)
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

		htmlResult, err := html.AppsHTML(apps, provDetector, recom)
		// htmlResult, err := html.LibsHTML(apps, provDetector, mapRec)
		if err != nil {
			fmt.Fprintf(w, "Error: %s", err)
		} else {
			w.Write(htmlResult)
		}
		fmt.Fprintf(w, "took %s", time.Since(start))
	}
}

func wrapHandler(token string, recom src.MapRecomended, pkgs ...string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		githubProv := github.New(github.NewHTTPWrapper(token, 10))
		provDetector := providers.NewDetector().AddProvider(githubProv)
		depDetector := deps.DefaultDetector()

		apps := make(chan i.IApp, 10)
		go func() {
			for _, pkg := range pkgs {
				apiApp, err := app.New(i.Pkg(pkg), "master", provDetector, depDetector)
				if err != nil {
					log.Printf("cant create app %s, got err: %s\n", pkg, err)
				}
				apps <- apiApp
			}
			close(apps)
		}()

		htmlResult, err := html.AppsHTML(apps, provDetector, recom)
		// htmlResult, err := html.LibsHTML(apps, provDetector, mapRec)
		if err != nil {
			fmt.Fprintf(w, "Error: %s", err)
		} else {
			w.Write(htmlResult)
		}
		fmt.Fprintf(w, "took %s", time.Since(start))
	}
}
