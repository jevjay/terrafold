package main

import (
	"bytes"
	"html/template"
)

type RDS struct {
	Name               string `yaml:"name,omitempty"`
	AllocatedStorage   int    `yaml:"allocated_storage,omitempty"`
	StorageType        string `yaml:"storage_type,omitempty"`
	Engine             string `yaml:"engine,omitempty"`
	EngineVersion      string `yaml:"engine_version,omitempty"`
	InstanceClass      string `yaml:"instance_class,omitempty"`
	Username           string `yaml:"username,omitempty"`
	Password           string `yaml:"password,omitempty"`
	DbSubnetGroupName  string `yaml:"db_subnet_group_name,omitempty"`
	ParameterGroupName string `yaml:"parameter_group_name,omitempty"`
}

const rdstmpl = `
resource "aws_db_instance" "{{.Name}}-rds" {
	{{- if .AllocatedStorage}}
	allocated_storage = {{.AllocatedStorage}}
	{{- end}}
	{{- if .StorageType}}
	storage_type = "{{.StorageType}}"
	{{- end}}
	{{- if .Engine}}
	engine = "{{.Engine}}"
	{{- end}}
	{{- if .EngineVersion}}
	engine_version = "{{.EngineVersion}}"
	{{- end}}
	{{- if .InstanceClass}}
	instance_class = "{{.InstanceClass}}"
	{{- end}}
	{{- if .Name}}
	name = "{{.Name}}"
	{{- end}}
	{{- if .Username}}
	username = "{{.Username}}"
	{{- end}}
	{{- if .Password}}
	password = "{{.Password}}"
	{{- end}}
	{{- if .DbSubnetGroupName}}
	db_subnet_group_name = "{{.DbSubnetGroupName}}"
	{{- end}}
	{{- if .ParameterGroupName}}
	parameter_group_name = "{{.ParameterGroupName}}"
	{{- end}}
}
`

func (r *RDS) generateContent() string {
	if r == nil {
		return ""
	}

	t := template.New("Ec2 template")
	t, err := t.Parse(ec2tmpl)

	check(err)
	var tpl bytes.Buffer
	err = t.Execute(&tpl, r)
	check(err)
	return tpl.String()
}
