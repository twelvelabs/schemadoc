{{ define "ConstTpl" -}}
{{ if . -}}

Allowed Value:

- `{{ . }}`

{{ end -}}
{{ end -}}

{{ define "DescriptionTpl" -}}
{{ if . -}}

{{ . }}

{{ end -}}
{{ end -}}

{{ define "EnumTpl" -}}
{{ if . -}}

Allowed Values:

{{ . }}

{{ end -}}
{{ end -}}

{{ define "ExamplesTpl" -}}
{{ if . -}}

Examples:

{{ range $example := . -}}

```yaml
{{ $example }}
```

{{ end -}}
{{ end -}}
{{ end -}}

{{ define "PropertiesTpl" -}}
{{ if . -}}
| Property | Type | Required | Enum | Default | Description |
| -------- | ---- | -------- | ---- | ------- | ----------- |
{{ $propParent := . -}}
{{ range $key, $prop := $propParent.Properties -}}
| [`{{ $key }}`](#{{ $key }}) | {{ $prop.TypeInfoMarkdown }} | {{ if $propParent.RequiredKey $key }}✅{{ else }}➖{{ end }} | {{ if $prop.Enum }}✅{{ else }}➖{{ end }} | {{ $prop.Default.JSONString | wrapCode | default "➖" }} | {{ $prop.DescriptionMarkdown | stripMarkdown | firstSentence | toHTML }} |
{{ end -}}
{{ end -}}
{{ end -}}

# {{ .EntityName }}

{{ template "DescriptionTpl" .DescriptionMarkdown }}
{{ template "ConstTpl" .Const.String }}
{{ template "EnumTpl" .EnumMarkdown }}
{{ template "ExamplesTpl" .Examples }}

{{ if .OneOf -}}

## Variants

{{ range $key, $schema := .OneOf -}}

- [{{ $schema.EntityName }}]({{ $schema.EntityLink }})
{{ end -}}
{{ end -}}
{{ if .Properties -}}

## Properties

{{ template "PropertiesTpl" . }}

{{ end -}}
{{ $root := . -}}
{{ range $key, $prop := .Properties -}}

### `{{ $key }}`

| Type | Required | Enum | Default |
| ---- | -------- | ---- | ------- |
| {{ $prop.TypeInfoMarkdown }} | {{ if $root.RequiredKey $key }}✅{{ else }}➖{{ end }} | {{ if $prop.Enum }}✅{{ else }}➖{{ end }} | {{ $prop.Default.JSONString | wrapCode | default "➖" }} |

{{ template "DescriptionTpl" $prop.DescriptionMarkdown }}
{{ template "ExamplesTpl" $prop.YAMLExamples }}

{{ end -}}
