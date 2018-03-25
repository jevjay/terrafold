package main

import (
	"bytes"
	"html/template"
)

type ELB struct {
	Name                      string   `yaml:"name,omitempty"`
	AvailabilityZones         []string `yaml:"availability_zones,omitempty"`
	Instances                 []string `yaml:"instances,omitempty"`
	CrossZoneLoadBalancing    bool     `yaml:"cross_zone_load_balancing,omitempty"`
	IdleTimeout               int      `yaml:"idle_timeout,omitempty"`
	ConnectionDraining        bool     `yaml:"connection_draining,omitempty"`
	ConnectionDrainingTimeout int      `yaml:"connection_draining_timeout,omitempty"`

	AccessLogs  *AccessLogs  `yaml:"access_logs,omitempty"`
	Listener    *Listener    `yaml:"listener,omitempty"`
	HealthCheck *HealthCheck `yaml:"health_check,omitempty"`
	Tags        *Tags        `yaml:"tags,omitempty"`
}

type AccessLogs struct {
	Bucket       string `yaml:"bucket,omitempty"`
	BucketPrefix string `yaml:"bucket_prefix,omitempty"`
	Interval     int    `yaml:"interval,omitempty"`
}

type Listener struct {
	InstancePort     int    `yaml:"instance_port,omitempty"`
	InstanceProtocol string `yaml:"instance_protocol,omitempty"`
	LbPort           int    `yaml:"lb_port,omitempty"`
	LbProtocol       string `yaml:"lb_protocol,omitempty"`
}

type HealthCheck struct {
	HealthyThreshold   int `yaml:"healthy_threshold,omitempty"`
	UnhealthyThreshold int `yaml:"unhealthy_threshold,omitempty"`
	Timeout            int `yaml:"timeout,omitempty"`
	Target             int `yaml:"target,omitempty"`
	Interval           int `yaml:"interval,omitempty"`
}

const elbtmpl = `
resource "aws_elb" "{{.Name}}" {
	{{- if .Name}}
	name               = "{{.Name}}"
	{{- end}}
	{{- if .AvailabilityZones}}
	availability_zones = [{{- range .AvailabilityZones}}"{{.}}",{{- end}}]
	{{- end}}
	{{- if .Instances}}
	instances = [{{- range .Instances}}"{{.}}",{{- end}}]
	{{- end}}
	{{- if .CrossZoneLoadBalancing}}
	cross_zone_load_balancing = {{.CrossZoneLoadBalancing}}
	{{- end}}
	{{- if .IdleTimeout}}
	idle_timeout = {{.IdleTimeout}}
	{{- end}}
	{{- if .ConnectionDraining}}
	connection_draining = {{.ConnectionDraining}}
	{{- end}}
	{{- if .ConnectionDrainingTimeout}}
	connection_draining_timeout = {{.ConnectionDrainingTimeout}}
	{{- end}}
	{{- if .AccessLogs}}
	access_logs {
	  bucket        = "{{.Bucket}}"
	  bucket_prefix = "{{.BucketPrefix}}"
	  interval      = {{.Interval}}
	}
	{{- end}}
	{{- if .Listener}}
	listener {
	  instance_port     = {{.InstancePort}}
	  instance_protocol = "{{.InstanceProtocol}}"
	  lb_port           = {{.LbPort}}
	  lb_protocol       = "{{.LbProtocol}}"
	}
	{{- end}}
	{{- if .HealthCheck}}
	health_check {
	  healthy_threshold   = {{.HealthyThreshold}}
	  unhealthy_threshold = {{.UnhealthyThreshold}}
	  timeout             = {{.Timeout}}
	  target              = "{{.Target}}"
	  interval            = {{.Interval}}
	}
	{{- end}}
	{{- if .Tags}}
	tags {
	  Name = "{{.Name}}"
	}
	{{- end}}
  } 
`

func (elb *ELB) generateContent() string {
	if elb == nil {
		return ""
	}

	t := template.New("ELB template")
	t, err := t.Parse(elbtmpl)

	if err != nil {
		panic(err)
	}

	var tpl bytes.Buffer
	err = t.Execute(&tpl, elb)

	if err != nil {
		panic(err)
	}

	return tpl.String()
}
