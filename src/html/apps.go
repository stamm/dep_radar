package html

import (
	"bytes"
	"html/template"
	"io/ioutil"

	i "github.com/stamm/dep_radar/interfaces"
	"github.com/stamm/dep_radar/src"
	"github.com/stamm/dep_radar/src/html/templates"
	"github.com/stamm/dep_radar/src/providers"
)

// AppsHTML return html with table. In the head apps, on the left side - libs
func AppsHTML(apps <-chan i.IApp, detector *providers.Detector, rec src.MapRecommended) ([]byte, error) {
	var buf bytes.Buffer
	raw, err := templates.Asset("src/html/templates/apps.html")
	if err != nil {
		return buf.Bytes(), err
	}
	tmpl, err := template.New("apps").Parse(string(raw))
	if err != nil {
		return buf.Bytes(), err
	}
	data := prepare(apps, detector, rec)
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
