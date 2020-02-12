// +build !ignore_autogenerated

// This file was autogenerated by openapi-gen. Do not edit it manually!

package v1alpha1

import (
	spec "github.com/go-openapi/spec"
	common "k8s.io/kube-openapi/pkg/common"
)

func GetOpenAPIDefinitions(ref common.ReferenceCallback) map[string]common.OpenAPIDefinition {
	return map[string]common.OpenAPIDefinition{
		"github.com/briangallagher/integreatly-operator/pkg/apis/monitoring/v1alpha1.ApplicationMonitoring": schema_pkg_apis_monitoring_v1alpha1_ApplicationMonitoring(ref),
		"github.com/briangallagher/integreatly-operator/pkg/apis/monitoring/v1alpha1.BlackboxTarget":        schema_pkg_apis_monitoring_v1alpha1_BlackboxTarget(ref),
	}
}

func schema_pkg_apis_monitoring_v1alpha1_ApplicationMonitoring(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "ApplicationMonitoring is the Schema for the applicationmonitorings API",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"kind": {
						SchemaProps: spec.SchemaProps{
							Description: "Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#types-kinds",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"apiVersion": {
						SchemaProps: spec.SchemaProps{
							Description: "APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#resources",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"metadata": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta"),
						},
					},
					"spec": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("github.com/briangallagher/integreatly-operator/pkg/apis/monitoring/v1alpha1.ApplicationMonitoringSpec"),
						},
					},
					"status": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("github.com/briangallagher/integreatly-operator/pkg/apis/monitoring/v1alpha1.ApplicationMonitoringStatus"),
						},
					},
				},
			},
		},
		Dependencies: []string{
			"github.com/briangallagher/integreatly-operator/pkg/apis/monitoring/v1alpha1.ApplicationMonitoringSpec", "github.com/briangallagher/integreatly-operator/pkg/apis/monitoring/v1alpha1.ApplicationMonitoringStatus", "k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta"},
	}
}

func schema_pkg_apis_monitoring_v1alpha1_BlackboxTarget(ref common.ReferenceCallback) common.OpenAPIDefinition {
	return common.OpenAPIDefinition{
		Schema: spec.Schema{
			SchemaProps: spec.SchemaProps{
				Description: "BlackboxTarget is the Schema for the blackboxtargets API",
				Type:        []string{"object"},
				Properties: map[string]spec.Schema{
					"kind": {
						SchemaProps: spec.SchemaProps{
							Description: "Kind is a string value representing the REST resource this object represents. Servers may infer this from the endpoint the client submits requests to. Cannot be updated. In CamelCase. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#types-kinds",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"apiVersion": {
						SchemaProps: spec.SchemaProps{
							Description: "APIVersion defines the versioned schema of this representation of an object. Servers should convert recognized schemas to the latest internal value, and may reject unrecognized values. More info: https://git.k8s.io/community/contributors/devel/api-conventions.md#resources",
							Type:        []string{"string"},
							Format:      "",
						},
					},
					"metadata": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta"),
						},
					},
					"spec": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("github.com/briangallagher/integreatly-operator/pkg/apis/monitoring/v1alpha1.BlackboxTargetSpec"),
						},
					},
					"status": {
						SchemaProps: spec.SchemaProps{
							Ref: ref("github.com/briangallagher/integreatly-operator/pkg/apis/monitoring/v1alpha1.BlackboxTargetStatus"),
						},
					},
				},
			},
		},
		Dependencies: []string{
			"github.com/briangallagher/integreatly-operator/pkg/apis/monitoring/v1alpha1.BlackboxTargetSpec", "github.com/briangallagher/integreatly-operator/pkg/apis/monitoring/v1alpha1.BlackboxTargetStatus", "k8s.io/apimachinery/pkg/apis/meta/v1.ObjectMeta"},
	}
}