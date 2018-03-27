package main

import (
	"bytes"
	"html/template"
)

type DigitaloceanFirewall struct {
	Name         string          `yaml:"name"`
	DropletIDs   []string        `yaml:"droplet_ids"`
	InboundRule  *[]InboundRule  `yaml:"inbound_rule"`
	OutboundRule *[]OutboundRule `yaml:"outbound_rule"`
}

type InboundRule struct {
	Protocol               string   `yaml:"protocol"`
	PortRange              string   `yaml:"port_range"`
	SourceAddresses        []string `yaml:"source_addresses"`
	SourceDropletIDs       []string `yaml:"source_droplet_ids"`
	SourceTags             []string `yaml:"source_tags"`
	SourceLoadBalancerUIDs []string `yaml:"source_load_balancer_uids"`
}

type OutboundRule struct {
	Protocol                    string   `yaml:"protocol"`
	PortRange                   string   `yaml:"port_range"`
	DestinationAddresses        []string `yaml:"destination_addresses"`
	DestinationDropletIDs       []string `yaml:"destination_droplet_ids"`
	DestinationTags             []string `yaml:"destination_tags"`
	DestinationLoadBalancerUIDs []string `yaml:"destination_load_balancer_uids"`
}

const dofwltmpl = `
resource "digitalocean_firewall" "{{.Name}}" {
	{{- if .Name}}
	name = "{{.Name}}"
	{{- end}}
	{{- if .DropletIDs}}
	image = [{{- range .DropletIDs}}"{{.}}",{{- end}}]
	{{- end}}
	{{- if .InboundRule}}
	inbound_rule [{{- range .InboundRule}}
	{
		{{- if .Protocol}}
		protocol = "{{.Protocol}}"
		{{- end}}
		{{- if .PortRange}}
		port_range = "{{.PortRange}}"
		{{- end}}
		{{- if .DestinationAddresses}}
		protocol = [{{- range .DestinationAddresses}}"{{.}}",{{- end}}]
		{{- end}}
		{{- if .SourceDropletIDs}}
		protocol = [{{- range .SourceDropletIDs}}"{{.}}",{{- end}}]
		{{- end}}
		{{- if .SourceTags}}
		protocol = [{{- range .SourceTags}}"{{.}}",{{- end}}]
		{{- end}}
		{{- if .SourceLoadBalancerUIDs}}
		protocol = [{{- range .SourceLoadBalancerUIDs}}"{{.}}",{{- end}}]
		{{- end}}
	}
	{{- end}}]
	{{- end}}
	{{- if .OutboundRule}}
	outbound_rule {
		{{- if .Protocol}}
		protocol = "{{.Protocol}}"
		{{- end}}
		{{- if .PortRange}}
		port_range = "{{.PortRange}}"
		{{- end}}
		{{- if .DestinationAddresses}}
		protocol = [{{- range .DestinationAddresses}}"{{.}}",{{- end}}]
		{{- end}}
		{{- if .DestinationDropletIDs}}
		protocol = [{{- range .DestinationDropletIDs}}"{{.}}",{{- end}}]
		{{- end}}
		{{- if .DestinationTags}}
		protocol = [{{- range .DestinationTags}}"{{.}}",{{- end}}]
		{{- end}}
		{{- if .DestinationLoadBalancerUIDs}}
		protocol = [{{- range .DestinationLoadBalancerUIDs}}"{{.}}",{{- end}}]
		{{- end}}
	}
	{{- end}}
}
`

func (dc *DigitaloceanFirewall) generateContent() string {
	if dc == nil {
		return ""
	}

	t := template.New("Digital Ocean Firewall template")
	t, err := t.Parse(dofwltmpl)

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
