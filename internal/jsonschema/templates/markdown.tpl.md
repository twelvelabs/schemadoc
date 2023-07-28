{{ define "ConstTpl" -}}
{{ if .String }}<br><br>Allowed value:<br>• {{ .String }}{{ end -}}
{{ end -}}

{{ define "EnumTpl" -}}
{{ if . }}<br><br>Allowed values:<br>
{{- $items := list -}}
{{- range $e := . }}{{ $items = cat "•" $e.String | append $items }}{{ end -}}
{{- $items | join "<br>" -}}
{{ end -}}
{{ end -}}

{{ define "PropertiesTpl" -}}
{{ if . -}}
| Property | Type | Required | Default | Description |
| -------- | ---- | -------- | ------- | ----------- |
{{ $propParent := . -}}
{{ range $key, $prop := $propParent.Properties -}}
| `{{ $key }}` | {{ $prop.TypeInfoMarkdown }} | {{ if $propParent.RequiredKey $key }}✅{{ end }} | {{ $prop.Default }} | {{ $prop.Description }}{{ template "ConstTpl" $prop.Const }}{{ template "EnumTpl" $prop.Enum }} |
{{ end }}
{{ end -}}
{{ end -}}

# {{ .EntityName }}

{{ .Description }}

## {{ .EntityName }} Properties

{{ template "PropertiesTpl" . -}}

{{ range $key, $def := .Definitions -}}
{{ if $def.Enum }}{{ continue }}{{ end -}}

## {{ $def.EntityName }}

{{ $def.Description }}

{{ if $def.Properties -}}

### {{ $def.EntityName }} Properties

{{ template "PropertiesTpl" $def -}}

{{ end -}}

{{ if $def.OneOf -}}

### {{ $def.EntityName }} Variants

{{ range $key, $subSchema := $def.OneOf -}}

- [{{ $subSchema.EntityName }}]({{ $subSchema.EntityLink }})
{{ end }}
{{ end -}}
{{ end -}}
