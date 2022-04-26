{{- define "k.result.secret" -}}
apiVersion: v1
kind: Secret
metadata:
  name: {{ .Release.Name }}-result
  annotations:
    iter8.tools/revision: {{ .Release.Revision }}
stringData:
  experiment.yaml: |
    startTime:         {{ now }},
    numCompletedTasks: 0,
    failure:           false,
    iter8Version:      {{ .Values.iter8lib.majorMinor }},
{{- end }}
