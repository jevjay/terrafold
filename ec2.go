package main

import (
	"bytes"
	"html/template"
)

type EC2 struct {
	Name                              string   `yaml:"name"`
	Ami                               string   `yaml:"ami"`
	AvailabilityZone                  string   `yaml:"availability_zone"`
	PlacementGroup                    string   `yaml:"placement_group"`
	Tenancy                           string   `yaml:"tenancy"`
	EbsOptimized                      bool     `yaml:"ebs_optimized"`
	DisableAPITermination             bool     `yaml:"disable_api_termination"`
	InstanceInitiatedShutdownBehavior string   `yaml:"instance_initiated_shutdown_behavior"`
	InstanceType                      string   `yaml:"instance_type"`
	KeyName                           string   `yaml:"key_name"`
	Monitoring                        bool     `yaml:"monitoring"`
	SecurityGroups                    []string `yaml:"security_groups"`
	VpcSecurityGroupIds               []string `yaml:"vpc_security_group_ids"`
	SubnetID                          string   `yaml:"subnet_id"`
	AssociatePublicIPAddress          bool     `yaml:"associate_public_ip_address"`
	PrivateIP                         string   `yaml:"private_ip"`
	SourceDestCheck                   bool     `yaml:"source_dest_check"`
	UserData                          string   `yaml:"user_data"`
	UserDataBase64                    string   `yaml:"user_data_base64"`
	IamInstanceProfile                string   `yaml:"iam_instance_profile"`
	IPv6AddressCount                  int      `yaml:"ipv6_address_count"`
	IPv6Addresses                     []string `yaml:"ipv6_addresses"`
	VolumeTags                        string   `yaml:"volume_tags"`

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
	DeleteOnTermination string `yaml:"delete_on_termination"`
}

type EbsBlockDevice struct {
	DeviceName          string `yaml:"device_name"`
	SnapshotID          string `yaml:"snapshot_id"`
	VolumeType          string `yaml:"volume_type"`
	VolumeSize          string `yaml:"volume_size"`
	Iops                string `yaml:"iops"`
	DeleteOnTermination string `yaml:"delete_on_termination"`
	Encrypted           string `yaml:"encrypted"`
}

type EphemeralBlockDevice struct {
	DeviceName  string `yaml:"encrypted"`
	VirtualName string `yaml:"encrypted"`
	NoDevice    string `yaml:"encrypted"`
}

type NetworkInterface struct {
	DeviceIndex         string `yaml:"device_index"`
	NetworkInterfaceId  string `yaml:"network_interface_id"`
	DeleteOnTermination string `yaml:"delete_on_termination"`
}

const tmpl = `
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

func (e EC2) generateContent() string {
	t := template.New("Ec2 template")
	t, err := t.Parse(tmpl)

	check(err)
	var tpl bytes.Buffer
	err = t.Execute(&tpl, e)
	check(err)
	return tpl.String()
}
