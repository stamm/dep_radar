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
	prov, detector := custom.Detector()
	depDetector := deps.DefaultDetector()

	pkgs, err := prov.GetAllRepos(context.Background(), "GO")
	if err != nil {
		fmt.Fprintf(w, "Error: %s", err)
		return
	}
	fmt.Printf("get all repos %s \n", time.Since(start))

	// Create a little wrapper with custom logic for detect
	apps := make([]i.IApp, 0, len(pkgs))
	for _, pkg := range pkgs {
		app, err := app.New(pkg, detector, depDetector)
		if err != nil {
			log.Fatal(err)
		}
		apps = append(apps, app)
	}
	fmt.Printf("get all apps %s \n", time.Since(start))
	htmlResult, err := html.AppsHtml(apps, detector)
	// htmlResult, err := html.LibsHtml(apps)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprint(w, htmlResult)
	fmt.Fprintf(w, "took %s", time.Since(start))
}
