# Changelog

{{ if .Versions -}}
{{ range .Versions }}

<!-- markdownlint-disable MD033 -->
<a name="{{ .Tag.Name }}"></a>
<!-- markdownlint-enable MD033 -->

## {{ if .Tag.Previous }}[{{ .Tag.Name }}]({{ $.Info.RepositoryURL }}/compare/{{ .Tag.Previous.Name }}...{{ .Tag.Name }}){{ else }}{{ .Tag.Name }}{{ end }} ({{ datetime "2006-01-02" .Tag.Date }})

{{ range .CommitGroups -}}

### {{ .Title }}

{{ range .Commits -}}

* {{ if .Scope }}**{{ .Scope }}:** {{ end }}{{ .Subject }}{{ if .Refs }} ({{ range $i, $ref := .Refs }}{{ if gt $i 0 }}, {{ end }}{{ $ref }}{{ end }}){{ end }}
{{ end }}

{{ end -}}

{{- if .NoteGroups -}}
{{ range .NoteGroups -}}

<!-- markdownlint-disable MD024 -->
### {{ .Title }}
<!-- markdownlint-enable MD024 -->

{{ range .Notes }}
{{ .Body }}
{{ end }}

{{ end -}}
{{ end -}}
{{ end -}}
{{ end -}}
