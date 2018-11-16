package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"time"

	"github.com/stamm/dep_radar"
	"github.com/stamm/dep_radar/app"
	"github.com/stamm/dep_radar/deps"
	"github.com/stamm/dep_radar/html"
	"github.com/stamm/dep_radar/providers"
	"github.com/stamm/dep_radar/providers/github"
	gitlabProv "github.com/stamm/dep_radar/providers/gitlab"
	"github.com/xanzy/go-gitlab"
)

var (
	gitlabClient   *gitlab.Client
	gitlabGoGetURL string
)

func init() {
	gitlabAPIURL := os.Getenv("GITLAB_API_URL")     // use this if you want to connect to private gitlab instance for e.g http://gitlab.company.org/api/v4
	gitlabGoGetURL = os.Getenv("GITLAB_GO_GET_URL") // gitlab.com or set your private go get url ( required )
	gitlabToken := os.Getenv("GITLAB_TOKEN")        // one of ways for auth in your gitlab ( see github.com/xanzy/go-gitlab for a lot of auth methods )
	gitlabClient = gitlab.NewClient(nil, gitlabToken)
	gitlabClient.SetBaseURL(gitlabAPIURL) // if you use gitlab.com you can skip this assigment
}

func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/favicon.ico", http.NotFound)
	log.Println("Started: http://localhost:8082/")
	http.ListenAndServe(":8082", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	provDetector := Detector()
	depDetector := deps.DefaultDetector()

	// Create a little wrapper with custom logic for detect
	apps := make(chan dep_radar.IApp, 100)

	go func() {
		defer close(apps)

		pkg := dep_radar.Pkg("gitlab.com/wxcsdb88/go")

		apiApp, err := app.New(context.Background(), pkg, "master", provDetector, depDetector)
		if err != nil {
			log.Printf("cant create app %s, got err: %s\n", pkg, err)
		}
		apps <- apiApp
	}()

	htmlResult, err := html.LibsHTML(context.Background(), apps, provDetector, nil)
	if err != nil {
		log.Fatal(err)
	}
	w.Write(htmlResult)
	fmt.Fprintf(w, "took %s", time.Since(start))
}

func Detector() *providers.Detector {
	prov := gitlabProv.New(gitlabClient, gitlabGoGetURL)
	githubProv := github.New(github.NewHTTPWrapper(os.Getenv("GITHUB_TOKEN"), 10))
	detector := providers.NewDetector().
		AddProvider(prov).
		AddProvider(githubProv)
	return detector
}
