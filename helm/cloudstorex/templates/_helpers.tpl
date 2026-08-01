{{/*
Expand the name of the chart.
*/}}
{{- define "cloudstorex.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
*/}}
{{- define "cloudstorex.fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- $name := default .Chart.Name .Values.nameOverride }}
{{- if contains $name .Release.Name }}
{{- .Release.Name | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}
{{- end }}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "cloudstorex.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "cloudstorex.labels" -}}
helm.sh/chart: {{ include "cloudstorex.chart" . }}
{{ include "cloudstorex.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "cloudstorex.selectorLabels" -}}
app.kubernetes.io/name: {{ include "cloudstorex.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "cloudstorex.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "cloudstorex.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}

{{/*
Resolve Database Host
*/}}
{{- define "cloudstorex.dbHost" -}}
{{- if .Values.postgresql.enabled -}}
{{- printf "%s-postgres" (include "cloudstorex.fullname" .) -}}
{{- else -}}
{{- .Values.externalDatabase.host -}}
{{- end -}}
{{- end -}}

{{/*
Resolve Database Name
*/}}
{{- define "cloudstorex.dbName" -}}
{{- if .Values.postgresql.enabled -}}
{{- .Values.postgresql.auth.database -}}
{{- else -}}
{{- .Values.externalDatabase.database -}}
{{- end -}}
{{- end -}}

{{/*
Resolve Database User
*/}}
{{- define "cloudstorex.dbUser" -}}
{{- if .Values.postgresql.enabled -}}
{{- .Values.postgresql.auth.username -}}
{{- else -}}
{{- .Values.externalDatabase.username -}}
{{- end -}}
{{- end -}}

{{/*
Resolve Database Password
*/}}
{{- define "cloudstorex.dbPassword" -}}
{{- if .Values.postgresql.enabled -}}
{{- .Values.postgresql.auth.password -}}
{{- else -}}
{{- .Values.externalDatabase.password -}}
{{- end -}}
{{- end -}}

{{/*
Resolve Redis Host
*/}}
{{- define "cloudstorex.redisHost" -}}
{{- if .Values.redis.enabled -}}
{{- printf "%s-redis" (include "cloudstorex.fullname" .) -}}
{{- else -}}
{{- .Values.externalRedis.host -}}
{{- end -}}
{{- end -}}

{{/*
Resolve MinIO/Storage Endpoint
*/}}
{{- define "cloudstorex.storageEndpoint" -}}
{{- if .Values.minio.enabled -}}
{{- printf "%s-minio:%d" (include "cloudstorex.fullname" .) (int .Values.minio.service.apiPort) -}}
{{- else -}}
{{- .Values.externalStorage.endpoint -}}
{{- end -}}
{{- end -}}

{{/*
Resolve Storage Bucket
*/}}
{{- define "cloudstorex.storageBucket" -}}
{{- if .Values.minio.enabled -}}
{{- .Values.minio.defaultBucket -}}
{{- else -}}
{{- .Values.externalStorage.bucket -}}
{{- end -}}
{{- end -}}

{{/*
Resolve Storage Access Key
*/}}
{{- define "cloudstorex.storageAccessKey" -}}
{{- if .Values.minio.enabled -}}
{{- .Values.minio.auth.rootUser -}}
{{- else -}}
{{- .Values.externalStorage.accessKey -}}
{{- end -}}
{{- end -}}

{{/*
Resolve Storage Secret Key
*/}}
{{- define "cloudstorex.storageSecretKey" -}}
{{- if .Values.minio.enabled -}}
{{- .Values.minio.auth.rootPassword -}}
{{- else -}}
{{- .Values.externalStorage.secretKey -}}
{{- end -}}
{{- end -}}
