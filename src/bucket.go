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

	ServerSideEncryptionConfiguration *ServerSideEncryptionConfiguration `yaml:"server_side_encryption_configuration,omitempty"`
	ReplicationConfiguration          *ReplicationConfiguration          `yaml:"replication_configuration,omitempty"`
	LifecycleRule                     *LifecycleRule                     `yaml:"lifecycle_rule,omitempty"`
	Logging                           *Logging                           `yaml:"logging,omitempty"`
	Versioning                        *Versioning                        `yaml:"versioning,omitempty"`
	CorsRule                          *CorsRule                          `yaml:"cors_rule,omitempty"`
	Website                           *Website                           `yaml:"website,omitempty"`
	Tags                              *Tags                              `yaml:"tags,omitempty"`
}

type LifecycleRule struct {
	ID                                 string `yaml:"id,omitempty"`
	Prefix                             string `yaml:"prefix,omitempty"`
	Enabled                            bool   `yaml:"enabled,omitempty"`
	AbortIncompleteMultipartUploadDays int    `yaml:"abort_incomplete_multipart_upload_days,omitempty"`

	Expiration                  `yaml:"expiration,omitempty"`
	Transition                  `yaml:"transition,omitempty"`
	NoncurrentVersionExpiration `yaml:"noncurrent_version_expiration,omitempty"`
	NoncurrentVersionTransition `yaml:"noncurrent_version_transition,omitempty"`
	Tags                        `yaml:"tags,omitempty"`
}

type Expiration struct {
	Date                      string `yaml:"date,omitempty"`
	Days                      int    `yaml:"days,omitempty"`
	ExpiredObjectDeleteMarker bool   `yaml:"expired_object_delete_marker,omitempty"`
}

type Transition struct {
	Date         string `yaml:"tags,omitempty"`
	Days         int    `yaml:"days,omitempty"`
	StorageClass string `yaml:"storage_class,omitempty"`
}

type NoncurrentVersionExpiration struct {
	Days int `yaml:"days,omitempty"`
}

type NoncurrentVersionTransition struct {
	Days         int    `yaml:"days,omitempty"`
	StorageClass string `yaml:"storage_class,omitempty"`
}

type Logging struct {
	TargetBucket string `yaml:"target_bucket,omitempty"`
	TargetPrefix string `yaml:"target_prefix,omitempty"`
}

type Versioning struct {
	Enabled   bool `yaml:"enabled,omitempty"`
	MfaDelete bool `yaml:"mfa_delete,omitempty"`
}

type CorsRule struct {
	AllowedHeaders []string `yaml:"allowed_headers,omitempty"`
	AllowedMethods []string `yaml:"allowed_methods,omitempty"`
	AllowedOrigins []string `yaml:"allowed_origins,omitempty"`
	ExposeHeaders  []string `yaml:"expose_headers,omitempty"`
	MaxAgeSeconds  int      `yaml:"max_age_seconds,omitempty"`
}

type ReplicationConfiguration struct {
	Role string `yaml:"role,omitempty"`

	Rules `yaml:"rules,omitempty"`
}

type Rules struct {
	ID     string `yaml:"id,omitempty"`
	Prefix string `yaml:"prefix,omitempty"`
	Status string `yaml:"status,omitempty"`

	Destination             `yaml:"destination,omitempty"`
	SourceSelectionCriteria `yaml:"source_selection_criteria,omitempty"`
}

type Destination struct {
	Bucket          string `yaml:"id,omitempty"`
	StorageClass    string `yaml:"storage_class,omitempty"`
	ReplicaKmsKeyID string `yaml:"replica_kms_key_id,omitempty"`
}

type SourceSelectionCriteria struct {
	SseKmsEncryptedObjects `yaml:"sse_kms_encrypted_objects,omitempty"`
}

type SseKmsEncryptedObjects struct {
	Enabled bool `yaml:"enabled,omitempty"`
}

type ServerSideEncryptionConfiguration struct {
	Rule `yaml:"rule,omitempty"`
}

type Rule struct {
	ApplyServerSideEncryptionByDefault `yaml:"apply_server_side_encryption_by_default,omitempty"`
}

type ApplyServerSideEncryptionByDefault struct {
	SseAlgorithm   string `yaml:"sse_algorithm,omitempty"`
	KmsMasterKeyID string `yaml:"kms_master_key_id,omitempty"`
}

type Website struct {
	IndexDocument         string `yaml:"index_document,omitempty"`
	ErrorDocument         string `yaml:"error_document,omitempty"`
	RedirectAllRequestsTo string `yaml:"redirect_all_requests_to,omitempty"`
	RoutingRules          string `yaml:"routing_rules,omitempty"`
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

	if err != nil {
		panic(err)
	}

	var tpl bytes.Buffer
	err = t.Execute(&tpl, b)

	if err != nil {
		panic(err)
	}

	return tpl.String()
}
