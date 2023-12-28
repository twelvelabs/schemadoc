{{ define "ConstTpl" -}}
{{ if .String }}<p>Allowed value:</p><ul><li><code>{{ .String }}</code></li></ul>{{ end -}}
{{ end -}}

{{ define "EnumTpl" -}}
{{ if . }}<p>Allowed values:</p>{{ . | toHTML }}{{ end -}}
{{ end -}}

{{ define "ExamplesTpl" -}}
{{ if . }}<p>Examples:</p>{{ . | toHTML }}{{ end -}}
{{ end -}}

{{ define "PropertiesTpl" -}}
{{ if . -}}
| Property | Type | Required | Default | Description |
| -------- | ---- | -------- | ------- | ----------- |
{{ $propParent := . -}}
{{ range $key, $prop := $propParent.Properties -}}
| `{{ $key }}` | {{ $prop.TypeInfoMarkdown }} | {{ if $propParent.RequiredKey $key }}✅{{ end }} | {{ $prop.Default }} | {{ $prop.DescriptionMarkdown | toHTML }}{{ template "ConstTpl" $prop.Const }}{{ template "EnumTpl" $prop.EnumMarkdown }}{{ template "ExamplesTpl" $prop.ExamplesMarkdown }} |
{{ end -}}
{{ end -}}
{{ end -}}

# {{ .EntityName }}

{{ .DescriptionMarkdown }}

## {{ .EntityName }} Properties

{{ template "PropertiesTpl" . }}
{{ range $key, $def := .Definitions -}}
{{ if $def.Enum }}{{ continue }}{{ end -}}

## {{ $def.EntityName }}

{{ $def.DescriptionMarkdown }}

{{ if $def.Properties -}}

### {{ $def.EntityName }} Properties

{{ template "PropertiesTpl" $def }}
{{ end -}}

{{ if $def.OneOf -}}

### {{ $def.EntityName }} Variants

{{ range $key, $subSchema := $def.OneOf -}}

- [{{ $subSchema.EntityName }}]({{ $subSchema.EntityLink }})
{{ end }}
{{ end -}}
{{ end -}}
