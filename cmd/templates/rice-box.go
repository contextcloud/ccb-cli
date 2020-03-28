package templates

import (
	"time"

	"github.com/GeertJohan/go.rice/embedded"
)

func init() {

	// define files
	file3 := &embedded.EmbeddedFile{
		Filename:    "deployment/deployment.yaml",
		FileModTime: time.Unix(1585403630, 0),

		Content: string("apiVersion: apps/v1beta1\nkind: Deployment\nmetadata:\n  name: {{ .Name }}\n  namespace: {{ .Namespace }}\n  labels: \n    release: {{ .Name }}\nspec:\n{{- if $.Replicas }}\n  replicas: {{ .Replicas }}\n{{- end }}\n  revisionHistoryLimit: 10\n  template:\n    metadata:\n      name: {{ .Name }}\n      labels:\n        release: {{ .Name }}\n        {{- with .Labels }}\n        {{- toYaml . | nindent 8 }}\n        {{- end }}\n{{- if $.Annotations }}\n      annotations:\n        {{- toYaml .Annotations | nindent 8 }}        \n{{- end }}\n    spec:\n      containers:\n      - name: {{ .Name }}\n        image: {{ .Image }}\n{{- if $.Environment }}\n        env:\n        {{- range $key, $value := .Environment }}\n        - name: {{ $key }}\n          value: {{ $value }}\n        {{- end }}\n{{- end }}\n        ports:\n        - name: http\n          containerPort: 8080\n        - name: metrics\n          containerPort: 8081\n        - name: health\n          containerPort: 8082\n{{- if $.LivenessProbe.Enabled }}\n        livenessProbe:\n          httpGet:\n            path: {{ $.LivenessProbe.Path }}\n            port: health\n            scheme: HTTP\n          initialDelaySeconds: {{ $.LivenessProbe.InitialDelaySeconds }}\n          timeoutSeconds: {{ $.LivenessProbe.TimeoutSeconds }}\n          periodSeconds: {{ $.LivenessProbe.PeriodSeconds }}\n{{- end }}\n{{- if $.ReadinessProbe.Enabled }}\n        readinessProbe:\n          httpGet:\n            path: {{ $.ReadinessProbe.Path }}\n            port: health\n            scheme: HTTP\n          initialDelaySeconds: {{ $.ReadinessProbe.InitialDelaySeconds }}\n          timeoutSeconds: {{ $.ReadinessProbe.TimeoutSeconds }}\n          periodSeconds: {{ $.ReadinessProbe.PeriodSeconds }}\n{{- end }}\n{{- if or $.Limits $.Requests }}\n        resources:\n{{- if $.Limits }}\n          limits:\n            {{- range $key, $value := .Limits }}\n            {{ $key }}: {{ $value }}\n            {{- end }}\n{{- end }}\n{{- if $.Requests }}\n          requests:\n            {{- range $key, $value := .Requests }}\n            {{ $key }}: {{ $value }}\n            {{- end }}\n{{- end }}\n{{- end }}\n{{- if $.NodeSelector }}\n      nodeSelector:\n        {{- range $key, $value := .NodeSelector }}\n        {{ $key }}: {{ $value }}\n        {{- end }}\n{{- end }}\n        securityContext:\n          readOnlyRootFilesystem: {{ $.ReadOnlyRootFilesystem }}\n        volumeMounts:\n        - mountPath: /tmp\n          name: temp\n{{- if $.Secrets }}\n      {{- range $key, $value := .Secrets }}\n        - mountPath: /var/cag/secrets\n          name: {{ $value.Name }}\n          readOnly: true\n      {{- end }}\n{{- end }}\n      volumes:\n      - emptyDir: {}\n        name: temp\n{{- if $.Secrets }}\n    {{- range $key, $value := .Secrets }}\n      - name: {{ $value.Name }}\n        projected:\n          defaultMode: 420\n          sources:\n          - secret:\n              name: {{ $value.Name }}\n    {{- end }}\n{{- end }}"),
	}
	file4 := &embedded.EmbeddedFile{
		Filename:    "deployment/service.yaml",
		FileModTime: time.Unix(1585403037, 0),

		Content: string("apiVersion: v1\nkind: Service\nmetadata:\n  name: {{ .Name }}\n  namespace: {{ .Namespace }}\n  labels: \n    release: {{ .Name }}\nspec:\n  ports:\n    - name: http\n      port: 8080\n    - name: metrics\n      port: 8081\n  selector:\n    release: {{ .Name }}"),
	}
	file5 := &embedded.EmbeddedFile{
		Filename:    "rice-box.go",
		FileModTime: time.Unix(1585405910, 0),

		Content: string(""),
	}
	file6 := &embedded.EmbeddedFile{
		Filename:    "templates.go",
		FileModTime: time.Unix(1585405908, 0),

		Content: string("//go:generate rice embed-go\n\npackage templates\n\nimport rice \"github.com/GeertJohan/go.rice\"\n\nfunc NewBox() *rice.Box {\n\tconf := rice.Config{\n\t\tLocateOrder: []rice.LocateMethod{rice.LocateEmbedded, rice.LocateAppended, rice.LocateFS},\n\t}\n\treturn conf.MustFindBox(\".\")\n}\n"),
	}

	// define dirs
	dir1 := &embedded.EmbeddedDir{
		Filename:   "",
		DirModTime: time.Unix(1585405910, 0),
		ChildFiles: []*embedded.EmbeddedFile{
			file5, // "rice-box.go"
			file6, // "templates.go"

		},
	}
	dir2 := &embedded.EmbeddedDir{
		Filename:   "deployment",
		DirModTime: time.Unix(1585399140, 0),
		ChildFiles: []*embedded.EmbeddedFile{
			file3, // "deployment/deployment.yaml"
			file4, // "deployment/service.yaml"

		},
	}

	// link ChildDirs
	dir1.ChildDirs = []*embedded.EmbeddedDir{
		dir2, // "deployment"

	}
	dir2.ChildDirs = []*embedded.EmbeddedDir{}

	// register embeddedBox
	embedded.RegisterEmbeddedBox(`.`, &embedded.EmbeddedBox{
		Name: `.`,
		Time: time.Unix(1585405910, 0),
		Dirs: map[string]*embedded.EmbeddedDir{
			"":           dir1,
			"deployment": dir2,
		},
		Files: map[string]*embedded.EmbeddedFile{
			"deployment/deployment.yaml": file3,
			"deployment/service.yaml":    file4,
			"rice-box.go":                file5,
			"templates.go":               file6,
		},
	})
}