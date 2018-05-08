package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/stamm/dep_radar"
	"github.com/stamm/dep_radar/app"
	"github.com/stamm/dep_radar/deps"
	"github.com/stamm/dep_radar/html"
	"github.com/stamm/dep_radar/html/templates"
	"github.com/stamm/dep_radar/providers"
	"github.com/stamm/dep_radar/providers/github"
)

const (
	port = "8081"
)

func main() {
	http.HandleFunc("/", handlerPage)
	http.HandleFunc("/api/", handlerAPI)
	http.HandleFunc("/favicon.ico", http.NotFound)
	log.Printf("Started: http://localhost:%s/\n", port)
	http.ListenAndServe(":"+port, nil)
}

func handlerPage(w http.ResponseWriter, r *http.Request) {
	raw, err := templates.Asset("html/templates/main.html")
	if err != nil {
		http.Error(w, fmt.Sprintf("err: %s", err), http.StatusBadRequest)
		return
	}
	w.Write(raw)
}

func handlerAPI(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	githubProv := github.New(github.NewHTTPWrapper(os.Getenv("GITHUB_TOKEN"), 10))
	provDetector := providers.NewDetector().AddProvider(githubProv)
	depDetector := deps.DefaultDetector()

	apps := make(chan dep_radar.IApp, 10)
	name := r.URL.Query().Get("name")
	if name == "" {
		http.Error(w, "empty name", http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	var (
		pkgs []dep_radar.Pkg
		err  error
	)
	if strings.Contains(name, "/") {
		pkgs = append(pkgs, dep_radar.Pkg(github.Prefix+"/"+name))
	} else {
		if githubProv.UserExists(ctx, name) {
			pkgs, err = githubProv.GetAllUserRepos(ctx, name)
			if err != nil {
				log.Printf("error while getting repos for user: %s\n", err)
				http.Error(w, fmt.Sprintf("error while getting repos for user: %s\n", err), http.StatusBadRequest)
				return
			}
		} else {
			pkgs, err = githubProv.GetAllOrgRepos(ctx, name)
			if err != nil {
				log.Printf("error for getting repos: %s\n", err)
				http.Error(w, fmt.Sprintf("error while getting repos: %s\n", err), http.StatusBadRequest)
				return
			}
		}
	}
	go func() {
		defer close(apps)
		for _, pkg := range pkgs {
			apiApp, err := app.New(ctx, pkg, "master", provDetector, depDetector)
			if err != nil {
				log.Printf("cant create app %s, got err: %s\n", pkg, err)
			}
			apps <- apiApp
		}
	}()

	var recom dep_radar.MapRecommended
	// 	err = json.Unmarshal(raw, &recom)
	// 	if err != nil {
	// 		fmt.Printf("error on unmarshal json: %s\n", err.Error())
	// 		os.Exit(1)
	// 	}
	data := html.Prepare(ctx, apps, provDetector, recom)
	result, err := json.Marshal(data)
	if err != nil {
		http.Error(w, fmt.Sprintf("error %s", err), http.StatusInternalServerError)
		return
	}
	w.Write(result)
}
