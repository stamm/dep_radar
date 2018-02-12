package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	_ "net/http/pprof"
	"time"

	"github.com/stamm/dep_radar/examples/bb/custom"
	i "github.com/stamm/dep_radar/interfaces"
	"github.com/stamm/dep_radar/src/app"
	"github.com/stamm/dep_radar/src/deps"
	"github.com/stamm/dep_radar/src/html"
)

func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/favicon.ico", http.NotFound)
	log.Println("Started: http://localhost:8081/")
	http.ListenAndServe(":8081", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
	// get all applications packages
	prov, provDetector := custom.Detector()
	depDetector := deps.DefaultDetector()

	// Create a little wrapper with custom logic for detect
	apps := make(chan i.IApp, 100)
	go func() {
		defer close(apps)
		pkgs, err := prov.GetAllRepos(context.Background(), "GO")
		if err != nil {
			fmt.Fprintf(w, "Error: %s", err)
			return
		}
		for _, pkg := range pkgs {
			apiApp, err := app.New(pkg, "master", provDetector, depDetector)
			if err != nil {
				log.Printf("cant create app %s, got err: %s\n", pkg, err)
			}
			apps <- apiApp
		}
	}()

	htmlResult, err := html.LibsHTML(apps, provDetector, nil)
	if err != nil {
		log.Fatal(err)
	}
	w.Write(htmlResult)
	fmt.Fprintf(w, "took %s", time.Since(start))
}
