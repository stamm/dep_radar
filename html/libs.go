package html

import (
	"bytes"
	"context"
	"html/template"
	"io/ioutil"

	"github.com/stamm/dep_radar"
	"github.com/stamm/dep_radar/html/templates"
	"github.com/stamm/dep_radar/providers"
)

// LibsHTML return html with table. In the head libs, on the left side - apps
func LibsHTML(ctx context.Context, apps <-chan dep_radar.IApp, detector *providers.Detector, rec dep_radar.MapRecommended) ([]byte, error) {
	var buf bytes.Buffer
	raw, err := templates.Asset("src/html/templates/libs.html")
	if err != nil {
		return buf.Bytes(), err
	}
	tmpl, err := template.New("libs").Parse(string(raw))
	if err != nil {
		return buf.Bytes(), err
	}
	data := Prepare(ctx, apps, detector, rec)
	err = tmpl.Execute(&buf, data)
	if err != nil {
		return buf.Bytes(), err
	}
	html, err := ioutil.ReadAll(&buf)
	if err != nil {
		return buf.Bytes(), err
	}
	return html, nil
}
