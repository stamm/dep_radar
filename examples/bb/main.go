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
	myhttp "github.com/stamm/dep_radar/http"
	"github.com/stamm/dep_radar/providers"
	bbprivate "github.com/stamm/dep_radar/providers/bitbucketprivate"
	"github.com/stamm/dep_radar/providers/github"

	"github.com/stamm/dep_radar/app"
	"github.com/stamm/dep_radar/deps"
	"github.com/stamm/dep_radar/html"
)

var (
	bbClient   *myhttp.Client
	bbGitURL   string
	bbAPIURL   string
	bbGoGetURL string
)

func init() {
	bbGitURL = os.Getenv("BB_GIT_URL")
	bbGoGetURL = os.Getenv("BB_GO_GET_URL")
	bbAPIURL = "https://" + bbGitURL
	bbClient = myhttp.NewClient(
		myhttp.Options{
			User:     os.Getenv("BB_USER"),
			Password: os.Getenv("BB_PASSWORD"),
		}, 10)
}

func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/favicon.ico", http.NotFound)
	log.Println("Started: http://localhost:8081/")
	http.ListenAndServe(":8081", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	// get all applications packages
	prov, provDetector := Detector()
	depDetector := deps.DefaultDetector()

	// Create a little wrapper with custom logic for detect
	apps := make(chan dep_radar.IApp, 100)
	go func() {
		defer close(apps)
		pkgs, err := prov.GetAllRepos(context.Background(), "GO")
		if err != nil {
			fmt.Fprintf(w, "Error: %s", err)
			return
		}
		for _, pkg := range pkgs {
			apiApp, err := app.New(context.Background(), pkg, "master", provDetector, depDetector)
			if err != nil {
				log.Printf("cant create app %s, got err: %s\n", pkg, err)
			}
			apps <- apiApp
		}
	}()

	htmlResult, err := html.LibsHTML(context.Background(), apps, provDetector, nil)
	if err != nil {
		log.Fatal(err)
	}
	w.Write(htmlResult)
	fmt.Fprintf(w, "took %s", time.Since(start))
}

// Detector returns bitbucket provider and detector
func Detector() (*bbprivate.Provider, *providers.Detector) {
	bbProv := bbprivate.New(bbClient, bbGitURL, bbGoGetURL, bbAPIURL)
	githubProv := github.New(github.NewHTTPWrapper(os.Getenv("GITHUB_TOKEN"), 10))
	detector := providers.NewDetector().
		AddProvider(bbProv).
		AddProvider(githubProv)
	return bbProv, detector
}
