{{/* vim: set filetype=mustache: */}}
{{/*
Expand the name of the chart.
*/}}
{{- define "certmanager-civo.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "certmanager-civo.fullname" -}}
{{- if .Values.fullnameOverride -}}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- $name := default .Chart.Name .Values.nameOverride -}}
{{- if contains $name .Release.Name -}}
{{- .Release.Name | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" -}}
{{- end -}}
{{- end -}}
{{- end -}}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "certmanager-civo.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" -}}
{{- end -}}

{{/*
Common labels
*/}}
{{- define "certmanager-civo.labels" -}}
helm.sh/chart: {{ include "certmanager-civo.chart" . }}
{{ include "certmanager-civo.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end -}}

{{/*
Selector labels
*/}}
{{- define "certmanager-civo.selectorLabels" -}}
app.kubernetes.io/name: {{ include "certmanager-civo.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end -}}

{{/*
Create the name of the service account to use
*/}}
{{- define "certmanager-civo.serviceAccountName" -}}
{{- if .Values.serviceAccount.create -}}
    {{ default (include "certmanager-civo.fullname" .) .Values.serviceAccount.name }}
{{- else -}}
    {{ default "default" .Values.serviceAccount.name }}
{{- end -}}
{{- end -}}

{{- define "certmanager-civo.rootCAIssuer" -}}
{{ printf "%s-ca" (include "certmanager-civo.fullname" .) }}
{{- end -}}

{{- define "certmanager-civo.rootCACertificate" -}}
{{ printf "%s-ca" (include "certmanager-civo.fullname" .) }}
{{- end -}}

{{- define "certmanager-civo.servingCertificate" -}}
{{ printf "%s-webhook-tls" (include "certmanager-civo.fullname" .) }}
{{- end -}}

{{- define "certmanager-civo.selfSignedIssuer" -}}
{{ printf "%s-selfsign" (include "certmanager-civo.fullname" .) }}
{{- end -}}

