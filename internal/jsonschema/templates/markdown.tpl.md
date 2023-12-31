{{ define "DescriptionTpl" -}}
{{ if . -}}

{{ . }}

{{ end -}}
{{ end -}}

{{ define "EnumTpl" -}}
{{ if . -}}

Allowed Values:

{{ range $enum := . -}}

- {{ $enum }}
{{ end -}}
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
{{ range $key, $prop := .Properties -}}
| [`{{ $prop.Key }}`](#{{ $prop.Key }}) | {{ $prop.TypeInfoMarkdown }} | {{ if $prop.Parent.RequiredKey $prop.Key }}✅{{ else }}➖{{ end }} | {{ if $prop.EnumMarkdownItems }}✅{{ else }}➖{{ end }} | {{ $prop.Default.JSONString | wrapCode | default "➖" }} | {{ $prop.DescriptionMarkdown | stripMarkdown | firstSentence | toHTML }} |
{{ end -}}
{{ end -}}
{{ end -}}

# {{ .EntityName }}

{{ template "DescriptionTpl" .DescriptionMarkdown }}
{{ template "EnumTpl" .EnumMarkdownItems }}
{{ template "ExamplesTpl" .YAMLExamples }}

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
{{ range $key, $prop := .Properties -}}

### `{{ $prop.Key }}`

| Type | Required | Enum | Default |
| ---- | -------- | ---- | ------- |
| {{ $prop.TypeInfoMarkdown }} | {{ if $prop.Parent.RequiredKey $prop.Key }}✅{{ else }}➖{{ end }} | {{ if $prop.EnumMarkdownItems }}✅{{ else }}➖{{ end }} | {{ $prop.Default.JSONString | wrapCode | default "➖" }} |

{{ template "DescriptionTpl" $prop.DescriptionMarkdown }}
{{ template "EnumTpl" $prop.EnumMarkdownItems }}
{{ template "ExamplesTpl" $prop.YAMLExamples }}

{{ end -}}
