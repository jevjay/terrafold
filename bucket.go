package main

import (
	"bytes"
	"html/template"
)

type Bucket struct {
	Name               string `yaml:"name,omitempty"`
	Bucket             string `yaml:"bucket,omitempty"`
	BucketPrefix       string `yaml:"bucket_prefix,omitempty"`
	ACL                string `yaml:"acl,omitempty"`
	Policy             string `yaml:"policy,omitempty"`
	ForceDestroy       bool   `yaml:"force_destroy,omitempty"`
	AccelerationStatus string `yaml:"acceleration_status,omitempty"`
	Region             string `yaml:"region,omitempty"`
	RequestPayer       string `yaml:"request_payer,omitempty"`

	ServerSideEncryptionConfiguration `yaml:"server_side_encryption_configuration,omitempty"`
	ReplicationConfiguration          `yaml:"replication_configuration,omitempty"`
	LifecycleRule                     `yaml:"lifecycle_rule,omitempty"`
	Logging                           `yaml:"logging,omitempty"`
	Versioning                        `yaml:"versioning,omitempty"`
	CorsRule                          `yaml:"cors_rule,omitempty"`
	Website                           `yaml:"website,omitempty"`
	Tags                              `yaml:"tags,omitempty"`
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
	Days                      int    `yaml:"days"`
	ExpiredObjectDeleteMarker bool   `yaml:"expired_object_delete_marker"`
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
	Enabled   bool `yaml:"enabled"`
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
resource "aws_s3_bucket" "{{.Name}}" {
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
	server_side_encryption_configuration {
		rule {
			apply_server_side_encryption_by_default {
				kms_master_key_id = "{{.ServerSideEncryptionConfiguration.Rule.ApplyServerSideEncryptionByDefault.KmsMasterKeyID}}"
				sse_algorithm = "{{.ServerSideEncryptionConfiguration.Rule.ApplyServerSideEncryptionByDefault.SseAlgorithm}}"
			}
		}
	}
	{{- end}}

	{{- if .ReplicationConfiguration}}
	replication_configuration {
		role = "{{.ReplicationConfiguration.Role}}"
		rules {
			{{- if .ReplicationConfiguration.Rules.ID}}
			id = "{{.ReplicationConfiguration.Rules.ID}}"
			{{- end}}
			prefix = "{{.ReplicationConfiguration.Rules.Prefix}}"
			status = "{{.ReplicationConfiguration.Rules.Status}}"

			destination {
				bucket = "{{.ReplicationConfiguration.Rules.Destination.Bucket}}"
				storage_class = "{{.ReplicationConfiguration.Rules.Destination.StorageClass}}"
				replica_kms_key_id = "{{.ReplicationConfiguration.Rules.Destination.ReplicaKmsKeyID}}"
			}
			{{- if .ReplicationConfiguration.Rules.SourceSelectionCriteria}}
			source_selection_criteria {
				sse_kms_encrypted_objects {
					enabled = "{{.ReplicationConfiguration.Rules.SourceSelectionCriteria.SseKmsEncryptedObjects.Enabled}}"
				}
			}
			{{- end}}
		}
	}
	{{- end}}

	{{- if .LifecycleRule}}
	lifecycle_rule {
		enabled = "{{.LifecycleRule.Enabled}}"
		{{- if .LifecycleRule.ID}}
		id = "{{.LifecycleRule.ID}}"
		{{- end}}
		{{- if .LifecycleRule.Prefix}}
		prefix = {{.LifecycleRule.Prefix}}
		{{- end}}
		{{- if .LifecycleRule.AbortIncompleteMultipartUploadDays}}
		abort_incomplete_multipart_upload_days = {{.LifecycleRule.AbortIncompleteMultipartUploadDays}}
		{{- end}}

		{{- if .LifecycleRule.Expiration}}
		expiration {
			{{- if .LifecycleRule.Expiration.Date}}
			date = "{{.LifecycleRule.Expiration.Date}}"
			{{- end}}
			{{- if .LifecycleRule.Expiration.Days}}
			days = {{.LifecycleRule.Expiration.Days}}
			{{- end}}
			{{- if .LifecycleRule.Expiration.ExpiredObjectDeleteMarker}}
			expired_object_delete_marker = {{.LifecycleRule.Expiration.ExpiredObjectDeleteMarker}}
			{{- end}}
		}
		{{- end}}

		{{- if .LifecycleRule.Transition}}
		transition {
			{{- if .LifecycleRule.Transition.Days}}
			days = {{.LifecycleRule.Transition.Days}}
			{{- end}}
			{{- if .LifecycleRule.Transition.StorageClass}}
			storage_class = "{{.LifecycleRule.Transition.StorageClass}}"
			{{- end}}
		}
		{{- end}}

		{{- if .LifecycleRule.NoncurrentVersionTransition}}
		noncurrent_version_transition {
			{{- if .LifecycleRule.NoncurrentVersionTransition.Days}}
			days = {{.LifecycleRule.NoncurrentVersionTransition.Days}}
			{{- end}}
			{{- if .LifecycleRule.NoncurrentVersionTransition.StorageClass}}
			storage_class = "{{.LifecycleRule.NoncurrentVersionTransition.StorageClass}}"
			{{- end}}
		  }
		{{- end}}

		{{- if .LifecycleRule.NoncurrentVersionExpiration}}
		noncurrent_version_expiration {
			{{- if .LifecycleRule.NoncurrentVersionExpiration.Days}}
			days = {{.LifecycleRule.NoncurrentVersionExpiration.Days}}
			{{- end}}
		}
		{{- end}}
	}
	{{- end}}

	{{- if .Logging}}
	logging {
		target_bucket = "{{.Logging.TargetBucket}}"
		{{- if .Logging.TargetPrefix}}
		target_prefix = "{{.Logging.TargetPrefix}}"
		{{- end}}
	}
	{{- end}}

	{{- if .Versioning}}
	versioning {
		enabled = {{.Versioning.Enabled}}
		mfa_delete = {{.Versioning.MfaDelete}}
	}
	{{- end}}

	{{- if .CorsRule}}
	cors_rule {
		{{- if .CorsRule.AllowedHeaders}}
		allowed_headers = [{{- range .CorsRule.AllowedHeaders}}"{{.}}",{{- end}}]
		{{- end}}
		allowed_methods = [{{- range .CorsRule.AllowedMethods}}"{{.}}",{{- end}}]
		allowed_origins = [{{- range .CorsRule.AllowedOrigins}}"{{.}}",{{- end}}]
		{{- if .CorsRule.ExposeHeaders}}
		expose_headers = [{{- range .CorsRule.ExposeHeaders}}"{{.}}",{{- end}}]
		{{- end}}
		{{- if .CorsRule.MaxAgeSeconds}}
		max_age_seconds = {{.CorsRule.MaxAgeSeconds}}
		{{- end}}
	}
	{{- end}}

	{{- if .Website}}
	website {
		index_document = "{{.Website.IndexDocument}}"
		{{- if .Website.ErrorDocument}}
		error_document = "{{.Website.ErrorDocument}}"
		{{- end}}
		{{- if .Website.RedirectAllRequestsTo}}
		redirect_all_requests_to = "{{.Website.RedirectAllRequestsTo}}"
		{{- end}}
		{{- if .Website.RoutingRules}}
		routing_rules = {{.Website.RoutingRules}}
		{{- end}}
	}
	{{- end}}
}
`

func (b *Bucket) generateContent() string {
	if b == nil {
		return ""
	}

	t := template.New("S3 Bucket template")
	t, err := t.Parse(s3tmpl)
	check(err)
	var tpl bytes.Buffer
	err = t.Execute(&tpl, b)
	check(err)
	return tpl.String()
}
