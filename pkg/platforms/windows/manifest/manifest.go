package manifest

import (
	"bytes"
	"embed"
	"html"
	"text/template"
)

//go:embed manifest.xml.tmpl
var IconFS embed.FS

func GetManifest(properties Properties) ([]byte, error) {
	if err := properties.Validate(); err != nil {
		return nil, err
	}
	properties.Identifier = html.EscapeString(properties.Identifier)

	tpl, err := template.New("manifest.xml.tmpl").ParseFS(IconFS, "manifest.xml.tmpl")
	if err != nil {
		return nil, err
	}

	var buff bytes.Buffer
	if err := tpl.Execute(&buff, properties); err != nil {
		return nil, err
	}

	return buff.Bytes(), nil
}
