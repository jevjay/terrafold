package main

import (
	"bytes"
	"html/template"
)

type Bucket struct {
	Name               string `yaml:"name"`
	Bucket             string `yaml:"bucket"`
	BucketPrefix       string `yaml:"bucket_prefix"`
	ACL                string `yaml:"acl"`
	Policy             string `yaml:"policy"`
	ForceDestroy       bool   `yaml:"force_destroy"`
	AccelerationStatus string `yaml:"acceleration_status"`
	Region             string `yaml:"region"`
	RequestPayer       string `yaml:"request_payer"`

	ServerSideEncryptionConfiguration `yaml:"server_side_encryption_configuration"`
	ReplicationConfiguration          `yaml:"replication_configuration"`
	LifecycleRule                     `yaml:"lifecycle_rule"`
	Logging                           `yaml:"logging"`
	Versioning                        `yaml:"versioning"`
	CorsRule                          `yaml:"cors_rule"`
	Website                           `yaml:"website"`
	Tags                              `yaml:"tags"`
}

type LifecycleRule struct {
	ID                                 string `yaml:"id"`
	Prefix                             string `yaml:"prefix"`
	Enabled                            bool   `yaml:"enabled"`
	AbortIncompleteMultipartUploadDays int    `yaml:"abort_incomplete_multipart_upload_days"`

	Expiration                  `yaml:"expiration"`
	Transition                  `yaml:"transition"`
	NoncurrentVersionExpiration `yaml:"noncurrent_version_expiration"`
	NoncurrentVersionTransition `yaml:"noncurrent_version_transition"`
	Tags                        `yaml:"tags"`
}

type Expiration struct {
	Date                      string `yaml:"date"`
	Days                      string `yaml:"days"`
	ExpiredObjectDeleteMarker string `yaml:"expired_object_delete_marker"`
}

type Transition struct {
	Date         string `yaml:"tags"`
	Days         int    `yaml:"days"`
	StorageClass string `yaml:"storage_class"`
}

type NoncurrentVersionExpiration struct {
	Days int `yaml:"days"`
}

type NoncurrentVersionTransition struct {
	Days         int    `yaml:"days"`
	StorageClass string `yaml:"storage_class"`
}

type Logging struct {
	TargetBucket string `yaml:"target_bucket"`
	TargetPrefix string `yaml:"target_prefix"`
}

type Versioning struct {
	enabled   bool `yaml:"enabled"`
	MfaDelete bool `yaml:"mfa_delete"`
}

type CorsRule struct {
	AllowedHeaders []string `yaml:"allowed_headers"`
	AllowedMethods []string `yaml:"allowed_methods"`
	AllowedOrigins []string `yaml:"allowed_origins"`
	ExposeHeaders  []string `yaml:"expose_headers"`
	MaxAgeSeconds  int      `yaml:"max_age_seconds"`
}

type ReplicationConfiguration struct {
	Role string `yaml:"role"`

	Rules `yaml:"rules"`
}

type Rules struct {
	ID     string `yaml:"id"`
	Prefix string `yaml:"prefix"`
	Status string `yaml:"status"`

	Destination             `yaml:"destination"`
	SourceSelectionCriteria `yaml:"source_selection_criteria"`
}

type Destination struct {
	Bucket          string `yaml:"id"`
	StorageClass    string `yaml:"storage_class"`
	ReplicaKmsKeyID string `yaml:"replica_kms_key_id"`
}

type SourceSelectionCriteria struct {
	SseKmsEncryptedObjects `yaml:"sse_kms_encrypted_objects"`
}

type SseKmsEncryptedObjects struct {
	Enabled bool `yaml:"enabled"`
}

type ServerSideEncryptionConfiguration struct {
	Rule `yaml:"rule"`
}

type Rule struct {
	ApplyServerSideEncryptionByDefault `yaml:"apply_server_side_encryption_by_default"`
}

type ApplyServerSideEncryptionByDefault struct {
	SseAlgorithm   string `yaml:"sse_algorithm"`
	KmsMasterKeyID string `yaml:"kms_master_key_id"`
}

type Website struct {
	IndexDocument         string `yaml:"index_document"`
	ErrorDocument         string `yaml:"error_document"`
	RedirectAllRequestsTo string `yaml:"redirect_all_requests_to"`
	RoutingRules          string `yaml:"routing_rules"`
}

const s3tmpl = `
resource "aws_instance" "{{.Name}}" {
	{{- if .Bucket }}
	bucket = "{{.Bucket}}"
	{{- end}}
	{{- if .BucketPrefix}}
	bucket_prefix = "{{.BucketPrefix}}"
	{{- end}}
	{{- if .ACL}}
	acl = "{{.ACL}}"
	{{- end}}
	{{- if .Policy}}
	policy = "{{.Policy}}"
	{{- end}}
	{{- if .ForceDestroy}}
	force_destroy = "{{.ForceDestroy}}"
	{{- end}}
	{{- if .AccelerationStatus}}
	acceleration_status = {{.AccelerationStatus}}
	{{- end}}
	{{- if .Region}}
	region = {{.Region}}
	{{- end}}
	{{- if .RequestPayer}}
	request_payer = {{.RequestPayer}}
	{{- end}}
	{{- if .ServerSideEncryptionConfiguration}}
	instance_initiated_shutdown_behavior {
		rule {
			apply_server_side_encryption_by_default {
				kms_master_key_id = "{{.InstanceInitiatedShutdownBehavior.Rule.ApplyServerSideEncryptionByDefault.KmsMasterKeyID}}"
				sse_algorithm = "{{.InstanceInitiatedShutdownBehavior.Rule.ApplyServerSideEncryptionByDefault.SseAlgorithm}}"
			}
		}
	}
	{{- end}}
	{{- if .ReplicationConfiguration}}
	replication_configuration {
		role = "{{.ReplicationConfiguration.Role}}"
		rules {
			id = "{{.ReplicationConfiguration.Role.ID}}"
			prefix = "{{.ReplicationConfiguration.Role.Prefix}}"
			status = "{{.ReplicationConfiguration.Role.Status}}"

			destination {
				bucket = "{{.ReplicationConfiguration.Role.Destination.Bucket}}"
				storage_class = "{{.ReplicationConfiguration.Role.Destination.StorageClass}}"
				replica_kms_key_id = "{{.ReplicationConfiguration.Role.Destination.ReplicaKmsKeyID}}"
			}

			source_selection_criteria {
				sse_kms_encrypted_objects {
					enabled = "{{.ReplicationConfiguration.Role.SourceSelectionCriteria.SseKmsEncryptedObjects.Enabled}}"
				}
			}
		} ""
	}
	{{- end}}
}
`

func (b Bucket) generateContent() string {
	t := template.New("Ec2 template")
	t, err := t.Parse(s3tmpl)

	check(err)
	var tpl bytes.Buffer
	err = t.Execute(&tpl, b)
	check(err)
	return tpl.String()
}
