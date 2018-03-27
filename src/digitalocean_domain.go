package main

import (
	"bytes"
	"html/template"
)

type DigitaloceanDomain struct {
	Name      string `yaml:"name"`
	IPAddress string `yaml:"ip_address"`
}

const dodomaintmpl = `
resource "digitalocean_domain" "{{.Name}}" {
	{{- if .Name}}
	name = "{{.Name}}"
	{{- end}}
	{{- if .IPAddress}}
	ip_address = "{{.IPAddress}}"
	{{- end}}
}
`

func (dc *DigitaloceanDomain) generateContent() string {
	if dc == nil {
		return ""
	}

	t := template.New("Digital Ocean Domain template")
	t, err := t.Parse(dodomaintmpl)

	if err != nil {
		panic(err)
	}

	var tpl bytes.Buffer
	err = t.Execute(&tpl, dc)

	if err != nil {
		panic(err)
	}

	return tpl.String()
}
