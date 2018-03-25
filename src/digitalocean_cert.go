package main

import (
	"bytes"
	"html/template"
)

type DigitaloceanCertificate struct {
	Name             string `yaml:"name"`
	PrivateKey       string `yaml:"private_key_path"`
	LeafCertificate  string `yaml:"leaf_certificate_path"`
	CertificateChain string `yaml:"certificate_chain_path"`
}

const docerttmpl = `
resource "digitalocean_certificate" "{{.Name}}" {
	{{- if .Name}}
	name = "{{.Name}}"
	{{- end}}
	{{- if .PrivateKeyName}}
	private_key = "{{.PrivateKeyName}}"
	{{- end}}
	{{- if .LeafCertificate}}
	leaf_certificate = "{{.LeafCertificate}}"
	{{- end}}
	{{- if .CertificateChain}}
	certificate_chain = "{{.CertificateChain}}"
	{{- end}}
}
`

func (dc *DigitaloceanCertificate) generateContent() string {
	if dc == nil {
		return ""
	}

	t := template.New("Ec2 template")
	t, err := t.Parse(docerttmpl)

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
