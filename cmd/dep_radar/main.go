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

var (
	githubProv   *github.Provider
	provDetector *providers.Detector
	depDetector  *deps.Detector
)

func init() {
	githubProv = github.New(github.NewHTTPWrapper(os.Getenv("GITHUB_TOKEN"), 10))
	provDetector = providers.NewDetector().AddProvider(githubProv)
	depDetector = deps.DefaultDetector()
}

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

	apps := make(chan dep_radar.IApp, 10)
	name := r.URL.Query().Get("name")
	if name == "" {
		http.Error(w, "empty name", http.StatusBadRequest)
		return
	}

	var mainErr error
	ctx := context.Background()
	pkgs, errs := getRepos(ctx, name)
	go func() {
		defer close(apps)
		select {
		case pkg, ok := <-pkgs:
			if !ok {
				return
			}
			if mainErr != nil {
				break
			}
			apiApp, err := app.New(ctx, pkg, "master", provDetector, depDetector)
			if err != nil {
				log.Printf("cant create app %s, got err: %s\n", pkg, err)
				break
			}
			apps <- apiApp
		case err, ok := <-errs:
			if !ok {
				return
			}
			mainErr = err
		}
	}()

	rules := r.URL.Query().Get("rules")
	var recom dep_radar.MapRecommended
	if rules != "" {
		err := json.Unmarshal([]byte(rules), &recom)
		if err != nil {
			http.Error(w, fmt.Sprintf("error on unmarshal json: %s\n", err.Error()), http.StatusBadRequest)
			return
		}
	}

	data := make(chan html.TemplateStruct)
	go func() {
		defer close(data)
		data <- html.Prepare(ctx, apps, provDetector, recom)
	}()
	for {
		select {
		case val := <-data:
			result, err := json.Marshal(val)
			if err != nil {
				http.Error(w, fmt.Sprintf("error %s", err), http.StatusInternalServerError)
				return
			}
			w.Write(result)
			return
		case err, ok := <-errs:
			if !ok && mainErr != nil {
				err = mainErr
			}
			if err != nil {
				log.Println(err)
				http.Error(w, fmt.Sprintf("error %s", err), http.StatusInternalServerError)
				return
			}
		}
	}
}

func getRepos(ctx context.Context, name string) (<-chan dep_radar.Pkg, <-chan error) {
	pkgs := make(chan dep_radar.Pkg, 10)
	errsCh := make(chan error)
	go func() {
		defer func() {
			close(pkgs)
			close(errsCh)
		}()
		if strings.Contains(name, "/") {
			pkgs <- dep_radar.Pkg(github.Prefix + "/" + name)
			return
		}

		if githubProv.UserExists(ctx, name) {
			pkgsUser, errs := githubProv.GetAllUserRepos(ctx, name)
			go func() {
				for pkg := range pkgsUser {
					pkgs <- pkg
				}
			}()
			err, ok := <-errs
			if ok && err != nil {
				errsCh <- fmt.Errorf("error while getting repos for user: %s", err)
				return
			}
			return
		}

		pkgsOrg, errs := githubProv.GetAllOrgRepos(ctx, name)
		go func() {
			for pkg := range pkgsOrg {
				pkgs <- pkg
			}
		}()
		err, ok := <-errs
		if ok && err != nil {
			errsCh <- fmt.Errorf("error while getting repos: %s", err)
			return
		}
	}()
	return pkgs, errsCh
}
