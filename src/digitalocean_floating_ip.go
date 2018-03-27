package main

import (
	"bytes"
	"html/template"
)

type DigitaloceanFloatingIP struct {
	Name                   string `yaml:"name"`
	DigitaloceanFloatingIP string `yaml:"digitalocean_floating_ip"`
	DigitaloceanDroplet    string `yaml:"digitalocean_droplet"`
}

const dofloatiptmpl = `
resource "digitalocean_floating_ip" "{{.Name}}" {
	{{- if .DigitaloceanFloatingIP}}
	digitalocean_floating_ip = "{{.DigitaloceanFloatingIP}}"
	{{- end}}
	{{- if .DigitaloceanDroplet}}
	digitalocean_droplet = "{{.DigitaloceanDroplet}}"
	{{- end}}
}
`

func (dc *DigitaloceanFloatingIP) generateContent() string {
	if dc == nil {
		return ""
	}

	t := template.New("Digital Ocean Floating IP template")
	t, err := t.Parse(dofloatiptmpl)

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
