package main

import (
	"bytes"
	"html/template"
)

type DigitaloceanLoadBalancer struct {
	Name                string   `yaml:"name"`
	Region              string   `yaml:"region"`
	Algorithm           string   `yaml:"algorithm"`
	RedirectHTTPToHTTPS bool     `yaml:"redirect_http_to_https"`
	DropletIDs          []string `yaml:"droplet_ids"`
	DropletTag          string   `yaml:"droplet_tag"`

	ForwardingRule *ForwardingRule `yaml:"forwarding_rule"`
	Healthcheck    *Healthcheck    `yaml:"healthcheck"`
	StickySessions *StickySessions `yaml:"sticky_sessions"`
}

type ForwardingRule struct {
	EntryProtocol  string `yaml:"entry_protocol"`
	EntryPort      int    `yaml:"entry_port"`
	TargetProtocol string `yaml:"target_protocol"`
	TargetPort     int    `yaml:"target_port"`
	CertificateID  string `yaml:"certificate_id"`
	TLSPassthrough bool   `yaml:"tls_passthrough"`
}

type Healthcheck struct {
	Protocol               string `yaml:"protocol"`
	Port                   int    `yaml:"port"`
	Path                   string `yaml:"path"`
	CheckIntervalSeconds   int    `yaml:"check_interval_seconds"`
	ResponseTimeoutSeconds int    `yaml:"response_timeout_seconds"`
	UnhealthyThreshold     int    `yaml:"unhealthy_threshold"`
	HealthyThreshold       int    `yaml:"healthy_threshold"`
}

type StickySessions struct {
	Type             string `yaml:"type"`
	CookieName       string `yaml:"cookie_name"`
	CookieTTLSeconds int    `yaml:"cookie_ttl_seconds"`
}

const dolbtmpl = `
resource "digitalocean_loadbalancer" "{{.Name}}" {
	name = "{{.Name}}"
	region = "{{.Region}}"
	{{- if .Algorithm}}
	algorithm = {{.Algorithm}}
	{{- end}}
	{{- if .RedirectHTTPToHTTPS}}
	redirect_http_to_https = {{.RedirectHTTPToHTTPS}}
	{{- end}}
	{{- if .DropletIDs}}
	droplet_ids = [{{- range .DropletIDs}}"{{.}}",{{- end}}]
	{{- end}}
	{{- if .DropletTag}}
	droplet_tag = "{{.DropletTag}}"
	{{- end}}

	forwarding_rule {
	  entry_port = {{.ForwardingRule.EntryPort}}
	  entry_protocol = "{{.ForwardingRule.ForwardingProtocol}}"
	  target_port ={{.ForwardingRule.TargetPort}}
	  target_protocol = "{{.ForwardingRule.TargetProtocol}}"
	  {{- if .ForwardingRule.CertificateID}}
	  certificate_id = {{.ForwardingRule.CertificateID}}
	  {{- end}}
	  {{- if .ForwardingRule.TLSPassthrough}}
	  tls_passthrough = {{.ForwardingRule.TLSPassthrough}}
	  {{- end}}
	}
	
	{{- if .HealthCheck}}
	healthcheck {
	  protocol = "{{.HealthCheck.Protocol}}"
	  {{- if .HealthCheck.Port}}
	  port = {{.HealthCheck.Port}}
	  {{- end}}
	  {{- if .HealthCheck.Path}}
	  path = "{{.HealthCheck.Path}}"
	  {{- end}}
	  {{- if .HealthCheck.CheckIntervalSeconds}}
	  check_interval_seconds = {{.HealthCheck.CheckIntervalSeconds}}
	  {{- end}}
	  {{- if .HealthCheck.ResponseTimeoutSeconds}}
	  response_timeout_seconds = {{.HealthCheck.ResponseTimeoutSeconds}}
	  {{- end}}
	  {{- if .HealthCheck.UnhealthyThreshold}}
	  unhealthy_threshold = {{.HealthCheck.UnhealthyThreshold}}
	  {{- end}}
	  {{- if .HealthCheck.HealthyThreshold}}
	  healthy_threshold = {{.HealthCheck.HealthyThreshold}}
	  {{- end}}
	}
	{{- end}}

	{{- if .StickySessions}}
	sticky_sessions {
		type = {{.StickySessions.Type}}
		{{- if .StickySessions.CookieName}}
		cookie_name = {{.StickySessions.CookieName}}
		{{- end}}
		{{- if .StickySessions.CookieTTLSeconds}}
		cookie_ttl_seconds = {{.StickySessions.CookieTTLSeconds}}
		{{- end}}
	}
	{{- end}}
  
	droplet_ids = ["${digitalocean_droplet.web.id}"]
  }
`

func (dc *DigitaloceanLoadBalancer) generateContent() string {
	if dc == nil {
		return ""
	}

	t := template.New("Digital Ocean Loadbalancer template")
	t, err := t.Parse(dolbtmpl)

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
