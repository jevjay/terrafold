package main

import (
	"bytes"
	"html/template"
)

type EC2 struct {
	Name                              string   `yaml:"name,omitempty"`
	Ami                               string   `yaml:"ami,omitempty"`
	AvailabilityZone                  string   `yaml:"availability_zone,omitempty"`
	PlacementGroup                    string   `yaml:"placement_group,omitempty"`
	Tenancy                           string   `yaml:"tenancy,omitempty"`
	EbsOptimized                      bool     `yaml:"ebs_optimized,omitempty"`
	DisableAPITermination             bool     `yaml:"disable_api_termination,omitempty"`
	InstanceInitiatedShutdownBehavior string   `yaml:"instance_initiated_shutdown_behavior,omitempty"`
	InstanceType                      string   `yaml:"instance_type,omitempty"`
	KeyName                           string   `yaml:"key_name,omitempty"`
	Monitoring                        bool     `yaml:"monitoring,omitempty"`
	SecurityGroups                    []string `yaml:"security_groups,omitempty"`
	VpcSecurityGroupIds               []string `yaml:"vpc_security_group_ids,omitempty"`
	SubnetID                          string   `yaml:"subnet_id,omitempty"`
	AssociatePublicIPAddress          bool     `yaml:"associate_public_ip_address,omitempty"`
	PrivateIP                         string   `yaml:"private_ip,omitempty"`
	SourceDestCheck                   bool     `yaml:"source_dest_check,omitempty"`
	UserData                          string   `yaml:"user_data,omitempty"`
	UserDataBase64                    string   `yaml:"user_data_base64,omitempty"`
	IamInstanceProfile                string   `yaml:"iam_instance_profile,omitempty"`
	IPv6AddressCount                  int      `yaml:"ipv6_address_count,omitempty"`
	IPv6Addresses                     []string `yaml:"ipv6_addresses,omitempty"`
	VolumeTags                        string   `yaml:"volume_tags,omitempty"`

	RootBlockDevice      `yaml:"root_block_device"`
	EbsBlockDevice       `yaml:"ebs_block_device"`
	EphemeralBlockDevice `yaml:"ephemeral_block_device"`
	NetworkInterface     `yaml:"network_interface"`
	Timeouts             `yaml:"timeouts"`
	Tags                 `yaml:"tags"`
}

type Timeouts struct {
	Create string `yaml:"create"`
	Update string `yaml:"update"`
	Delete string `yaml:"delete"`
}

type RootBlockDevice struct {
	VolumeType          string `yaml:"volume_type"`
	VolumeSize          string `yaml:"volume_size"`
	Iops                string `yaml:"iops"`
	DeleteOnTermination bool   `yaml:"delete_on_termination"`
}

type EbsBlockDevice struct {
	DeviceName          string `yaml:"device_name"`
	SnapshotID          string `yaml:"snapshot_id"`
	VolumeType          string `yaml:"volume_type"`
	VolumeSize          int    `yaml:"volume_size"`
	Iops                string `yaml:"iops"`
	DeleteOnTermination bool   `yaml:"delete_on_termination"`
	Encrypted           bool   `yaml:"encrypted"`
}

type EphemeralBlockDevice struct {
	DeviceName  string `yaml:"device_name"`
	VirtualName string `yaml:"virtual_name"`
	NoDevice    bool   `yaml:"no_device"`
}

type NetworkInterface struct {
	DeviceIndex         string `yaml:"device_index"`
	NetworkInterfaceID  string `yaml:"network_interface_id"`
	DeleteOnTermination string `yaml:"delete_on_termination"`
}

const ec2tmpl = `
resource "aws_instance" "{{.Name}}" {
	{{- if .Ami }}
	ami = "{{.Ami}}"
	{{- end}}
	{{- if .InstanceType}}
	instance_type = "{{.InstanceType}}"
	{{- end}}
	{{- if .AvailabilityZone}}
	availability_zone = "{{.AvailabilityZone}}"
	{{- end}}
	{{- if .PlacementGroup}}
	placement_group = "{{.PlacementGroup}}"
	{{- end}}
	{{- if .Tenancy}}
	tenancy = "{{.Tenancy}}"
	{{- end}}
	{{- if .EbsOptimized}}
	ebs_optimized = {{.EbsOptimized}}
	{{- end}}
	{{- if .DisableAPITermination}}
	disable_api_termination = {{.DisableAPITermination}}
	{{- end}}
	{{- if .InstanceInitiatedShutdownBehavior}}
	instance_initiated_shutdown_behavior = "{{.InstanceInitiatedShutdownBehavior}}"
	{{- end}}
	{{- if .KeyName}}
	key_name = "{{.KeyName}}"
	{{- end}}
	{{- if .Monitoring}}
	monitoring = {{.Monitoring}}
	{{- end}}
	{{- if .SecurityGroups}}
	security_groups = [{{- range .SecurityGroups}}"{{.}}",{{- end}}]
	{{- end}}
	{{- if .VpcSecurityGroupIds}}
	vpc_security_group_ids = [{{- range .VpcSecurityGroupIds}}"{{.}}",{{- end}}]
	{{- end}}
	{{- if .SubnetID}}
	subnet_id = "{{.SubnetID}}"
	{{- end}}
	{{- if .AssociatePublicIPAddress}}
	associate_public_ip_address = "{{.AssociatePublicIPAddress}}"
	{{- end}}
	{{- if .PrivateIP}}
	private_ip = "{{.PrivateIP}}"
	{{- end}}
	{{- if .SourceDestCheck}}
	source_dest_check = "{{.SourceDestCheck}}"
	{{- end}}
	{{- if .UserData}}
	user_data = "{{.UserData}}"
	{{- end}}
	{{- if .UserDataBase64}}
	user_data_base64 = "{{.UserDataBase64}}"
	{{- end}}
	{{- if .IamInstanceProfile}}
	iam_instance_profile = "{{.IamInstanceProfile}}"
	{{- end}}
	{{- if .IPv6AddressCount}}
	ipv6_address_count = {{.IPv6AddressCount}}
	{{- end}}
	{{- if .IPv6Addresses}}
	ipv6_addresses = [{{- range .IPv6Addresses}}"{{.}}",{{- end}}]
	{{- end}}
	{{- if .VolumeTags}}
	volume_tags = {{.VolumeTags}}
	{{- end}}
}
`

func (e *EC2) generateContent() string {
	if e == nil {
		return ""
	}

	t := template.New("Ec2 template")
	t, err := t.Parse(ec2tmpl)

	check(err)
	var tpl bytes.Buffer
	err = t.Execute(&tpl, e)
	check(err)
	return tpl.String()
}
