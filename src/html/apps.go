package html

import (
	"bytes"
	"html/template"
	"io/ioutil"

	i "github.com/stamm/dep_radar/interfaces"
	"github.com/stamm/dep_radar/src/providers"
)

// AppsHTML return html with table. In the head apps, on the left side - libs
func AppsHTML(apps <-chan i.IApp, detector *providers.Detector, rec MapRecomended) ([]byte, error) {
	var buf bytes.Buffer
	tmpl, err := template.ParseFiles("src/html/apps.html")
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
