package main

import (
	"bytes"
	"html/template"
)

type DigitaloceanDroplet struct {
	Name              string   `yaml:"name"`
	Image             string   `yaml:"image"`
	Region            string   `yaml:"region"`
	Size              int      `yaml:"size"`
	Backups           bool     `yaml:"backups"`
	Monitoring        bool     `yaml:"monitoring"`
	Ipv6              bool     `yaml:"ipv6"`
	PrivateNetworking bool     `yaml:"private_networking"`
	SSHKeys           []string `yaml:"ssh_keys"`
	ResizeDisk        bool     `yaml:"resize_disk"`
	UserData          string   `yaml:"user_data"`
	VolumeIds         []string `yaml:"volume_ids"`
}

const dodroptmpl = `
resource "digitalocean_droplet" "{{.Name}}" {
	{{- if .Name}}
	name = "{{.Name}}"
	{{- end}}
	{{- if .Image}}
	image = "{{.Image}}"
	{{- end}}
	{{- if .Region}}
	region = "{{.Region}}"
	{{- end}}
	{{- if .Size}}
	size = {{.Size}}
	{{- end}}
	{{- if .Backups}}
	backups = {{.Backups}}
	{{- end}}
	{{- if .Monitoring}}
	monitoring = {{.Monitoring}}
	{{- end}}
	{{- if .Ipv6}}
	ipv6 = {{.Ipv6}}
	{{- end}}
	{{- if .PrivateNetworking}}
	private_networking = [{{- range .PrivateNetworking}}"{{.}}",{{- end}}]
	{{- end}}
	{{- if .SSHKeys}}
	ssh_keys = {{.SSHKeys}}
	{{- end}}
	{{- if .ResizeDisk}}
	resize_disk = {{.ResizeDisk}}
	{{- end}}
	{{- if .UserData}}
	user_data = {{.UserData}}
	{{- end}}
	{{- if .VolumeIds}}
	volume_ids = [{{- range .PrivateNetworking}}"{{.}}",{{- end}}]
	{{- end}}
}
`

func (dc *DigitaloceanDroplet) generateContent() string {
	if dc == nil {
		return ""
	}

	t := template.New("Ec2 template")
	t, err := t.Parse(dodroptmpl)

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
