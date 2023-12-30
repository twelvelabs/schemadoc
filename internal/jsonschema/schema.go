package jsonschema

import (
	"embed"
	"encoding/json"
	"fmt"
	"net/url"
	"path"
	"strings"

	"github.com/gobuffalo/flect"
	"github.com/twelvelabs/termite/render"
	"gopkg.in/yaml.v3"
)

//go:embed templates/*
var Templates embed.FS

// Schema represents a JSON schema document.
type Schema struct {
	AdditionalItems       *Schema            `json:"additionalItems,omitempty"`
	AdditionalProperties  any                `json:"additionalProperties,omitempty"`
	AllOf                 []*Schema          `json:"allOf,omitempty"`
	AnyOf                 []*Schema          `json:"anyOf,omitempty"`
	Comment               string             `json:"$comment,omitempty"`
	Const                 Any                `json:"const,omitempty"`
	Contains              *Schema            `json:"contains,omitempty"`
	ContentEncoding       string             `json:"contentEncoding,omitempty"`
	ContentMediaType      string             `json:"contentMediaType,omitempty"`
	Default               Any                `json:"default,omitempty"`
	Definitions           map[string]*Schema `json:"definitions,omitempty"`
	Deprecated            bool               `json:"deprecated,omitempty"`
	Description           string             `json:"description,omitempty"`
	Else                  *Schema            `json:"else,omitempty"`
	Enum                  []Any              `json:"enum,omitempty"`
	EnumDescriptions      []string           `json:"enumDescriptions,omitempty"`
	Examples              []Any              `json:"examples,omitempty"`
	ExclusiveMaximum      float64            `json:"exclusiveMaximum,omitempty"`
	ExclusiveMinimum      float64            `json:"exclusiveMinimum,omitempty"`
	Format                string             `json:"format,omitempty"`
	ID                    string             `json:"$id,omitempty"`
	If                    *Schema            `json:"if,omitempty"`
	Items                 *Schema            `json:"items,omitempty"`
	MarkdownDescription   string             `json:"markdownDescription,omitempty"`
	MaxContains           int                `json:"maxContains,omitempty"`
	Maximum               float64            `json:"maximum,omitempty"`
	MaxItems              int                `json:"maxItems,omitempty"`
	MaxLength             int                `json:"maxLength,omitempty"`
	MaxProperties         int                `json:"maxProperties,omitempty"`
	MinContains           int                `json:"minContains,omitempty"`
	Minimum               float64            `json:"minimum,omitempty"`
	MinItems              int                `json:"minItems,omitempty"`
	MinLength             int                `json:"minLength,omitempty"`
	MinProperties         int                `json:"minProperties,omitempty"`
	MultipleOf            float64            `json:"multipleOf,omitempty"`
	Not                   []*Schema          `json:"not,omitempty"`
	OneOf                 []*Schema          `json:"oneOf,omitempty"`
	Pattern               string             `json:"pattern,omitempty"`
	PatternProperties     map[string]*Schema `json:"patternProperties,omitempty"`
	PrefixItems           []*Schema          `json:"prefixItems,omitempty"`
	Properties            map[string]*Schema `json:"properties,omitempty"`
	PropertyNames         *Schema            `json:"propertyNames,omitempty"`
	ReadOnly              bool               `json:"readOnly,omitempty"`
	Ref                   string             `json:"$ref,omitempty"`
	Required              []string           `json:"required,omitempty"`
	Schema                string             `json:"$schema,omitempty"`
	Then                  *Schema            `json:"then,omitempty"`
	Title                 string             `json:"title,omitempty"`
	Types                 Types              `json:"type,omitempty"`
	UnevaluatedProperties bool               `json:"unevaluatedProperties,omitempty"`
	UniqueItems           bool               `json:"uniqueItems,omitempty"`
	WriteOnly             bool               `json:"writeOnly,omitempty"`

	/**
	* Internal, non-schema fields.
	**/

	Context      *Context        `json:"-"`
	Document     any             `json:"-"`
	GenPathTpl   render.Template `json:"-"`
	Key          string          `json:"key,omitempty"`
	Parent       *Schema         `json:"-"`
	Resolved     bool            `json:"resolved,omitempty"`
	RetrievalURI string          `json:"retrievalURI,omitempty"`
}

// BaseURI returns the resolved base URI for the schema.
// The base URI is the schema $id attribute resolved against
// the retrieval URI.
func (s *Schema) BaseURI() *url.URL {
	retURI, err := url.Parse(s.Root().RetrievalURI)
	if err != nil {
		panic(err)
	}

	if s.ID == "" {
		// If the $id attribute is missing, then the base URI is assumed
		// to be the same as the retrieval URI.
		// See: https://json-schema.org/understanding-json-schema/structuring.html#retrieval-uri
		return retURI
	}

	baseURI, err := url.Parse(s.Root().ID)
	if err != nil {
		panic(err)
	}
	// The $id attribute MAY be relative to the retrieval URI.
	// See: https://json-schema.org/understanding-json-schema/structuring.html#id
	return retURI.ResolveReference(baseURI)
}

// RefURI resolves the given ref against the base URI.
func (s *Schema) RefURI(ref string) *url.URL {
	refURI, err := url.Parse(ref)
	if err != nil {
		panic(err)
	}
	return s.BaseURI().ResolveReference(refURI)
}

// Clone returns a deep-copied clone of the receiver.
func (s *Schema) Clone() (*Schema, error) {
	doc, err := json.Marshal(s)
	if err != nil {
		return nil, err
	}
	return s.Context.parse(doc, s.Parent)
}

// DescriptionMarkdown returns the schema description formatted as Markdown,
// Will prioritize the non-standard `markdownDescription` attribute if present,
// otherwise uses `description`.
func (s *Schema) DescriptionMarkdown() string {
	if s.MarkdownDescription != "" {
		return s.MarkdownDescription
	}
	return s.Description
}

func (s *Schema) EntityName() string {
	if s.Title != "" {
		return flect.Pascalize(s.Title)
	}
	if s.Key != "" {
		return flect.Pascalize(s.Key)
	}
	if s.Ref != "" {
		return flect.Pascalize(path.Base(s.Ref))
	}
	return ""
}

func (s *Schema) EntityLink() string {
	anchor := strings.ToLower(s.EntityName())
	anchor = strings.ReplaceAll(anchor, " ", "-")
	return s.GenPath() + "#" + anchor
}

func (s *Schema) GenPath() string {
	path, err := s.Root().GenPathTpl.Render(s)
	if err != nil {
		panic(err)
	}
	return path
}

func (s *Schema) EnsureDocument() {
	if s.Document == nil {
		buf, err := json.Marshal(s)
		if err != nil {
			panic(err)
		}
		err = json.Unmarshal(buf, &s.Document)
		if err != nil {
			panic(err)
		}
	}
}

// EnumMarkdown returns the enum and enum descriptions
// formatted as a markdown list.
func (s *Schema) EnumMarkdown() string {
	items := []string{}

	for idx, enum := range s.Enum {
		desc := ""
		if idx < len(s.EnumDescriptions) {
			desc = s.EnumDescriptions[idx]
		}
		item := "- `" + enum.String() + "`"
		if desc != "" {
			item += ": " + desc
		}
		items = append(items, item)
	}

	if len(items) == 0 {
		return ""
	}
	return strings.Join(items, "\n")
}

// ExamplesMarkdown returns the examples as a markdown list.
func (s *Schema) ExamplesMarkdown() string {
	items := []string{}

	for _, example := range s.Examples {
		items = append(items, fmt.Sprintf(" * `%s`", example.YAMLString()))
	}

	if len(items) == 0 {
		return ""
	}
	return strings.Join(items, "\n")
}

// Merge merges fields from other into the receiver.
func (s *Schema) Merge(other *Schema) {
	s.EnsureDocument()
	other.EnsureDocument()
	doc, ok := other.Document.(map[string]any)
	if !ok {
		// Empty or invalid doc, nothing to do.
		return
	}
	for key := range doc {
		switch key {
		case "additionalItems":
			s.AdditionalItems = other.AdditionalItems
		case "additionalProperties":
			s.AdditionalProperties = other.AdditionalProperties
		case "allOf":
			s.AllOf = other.AllOf
		case "anyOf":
			s.AnyOf = other.AnyOf
		case "$comment":
			s.Comment = other.Comment
		case "const":
			s.Const = other.Const
		case "contains":
			s.Contains = other.Contains
		case "contentEncoding":
			s.ContentEncoding = other.ContentEncoding
		case "contentMediaType":
			s.ContentMediaType = other.ContentMediaType
		case "default":
			s.Default = other.Default
		case "definitions":
			s.Definitions = other.Definitions
		case "deprecated":
			s.Deprecated = other.Deprecated
		case "description":
			s.Description = other.Description
		case "else":
			s.Else = other.Else
		case "enum":
			s.Enum = other.Enum
		case "enumDescriptions":
			s.EnumDescriptions = other.EnumDescriptions
		case "examples":
			s.Examples = other.Examples
		case "exclusiveMaximum":
			s.ExclusiveMaximum = other.ExclusiveMaximum
		case "exclusiveMinimum":
			s.ExclusiveMinimum = other.ExclusiveMinimum
		case "format":
			s.Format = other.Format
		case "$id":
			s.ID = other.ID
		case "if":
			s.If = other.If
		case "items":
			s.Items = other.Items
		case "markdownDescription":
			s.MarkdownDescription = other.MarkdownDescription
		case "maxContains":
			s.MaxContains = other.MaxContains
		case "maximum":
			s.Maximum = other.Maximum
		case "maxItems":
			s.MaxItems = other.MaxItems
		case "maxLength":
			s.MaxLength = other.MaxLength
		case "maxProperties":
			s.MaxProperties = other.MaxProperties
		case "minContains":
			s.MinContains = other.MinContains
		case "minimum":
			s.Minimum = other.Minimum
		case "minItems":
			s.MinItems = other.MinItems
		case "minLength":
			s.MinLength = other.MinLength
		case "minProperties":
			s.MinProperties = other.MinProperties
		case "multipleOf":
			s.MultipleOf = other.MultipleOf
		case "not":
			s.Not = other.Not
		case "oneOf":
			s.OneOf = other.OneOf
		case "pattern":
			s.Pattern = other.Pattern
		case "patternProperties":
			s.PatternProperties = other.PatternProperties
		case "prefixItems":
			s.PrefixItems = other.PrefixItems
		case "properties":
			s.Properties = other.Properties
		case "propertyNames":
			s.PropertyNames = other.PropertyNames
		case "readOnly":
			s.ReadOnly = other.ReadOnly
		case "$ref":
			s.Ref = other.Ref
		case "required":
			s.Required = other.Required
		case "$schema":
			s.Schema = other.Schema
		case "then":
			s.Then = other.Then
		case "title":
			s.Title = other.Title
		case "type":
			s.Types = other.Types
		case "unevaluatedProperties":
			s.UnevaluatedProperties = other.UnevaluatedProperties
		case "uniqueItems":
			s.UniqueItems = other.UniqueItems
		case "writeOnly":
			s.WriteOnly = other.WriteOnly
		default:
		}
	}
}

// RequiredKey returns true when key is a required property.
func (s *Schema) RequiredKey(key string) bool {
	for _, k := range s.Required {
		if k == key {
			return true
		}
	}
	return false
}

// Root returns the root schema.
func (s *Schema) Root() *Schema {
	if s.Parent != nil {
		return s.Parent.Root()
	}
	return s
}

func (s *Schema) TypeInfo() []TypeInfo {
	result := []TypeInfo{}
	for _, t := range s.Types {
		switch t {
		case TypeArray:
			result = append(result, TypeInfo{
				Type:   t,
				Schema: s.Items,
			})
		case TypeObject:
			result = append(result, TypeInfo{
				Type:   t,
				Schema: s,
			})
		default:
			result = append(result, TypeInfo{
				Type:   t,
				Schema: nil,
			})
		}
	}
	return result
}

func (s *Schema) TypeInfoMarkdown() string {
	segments := []string{}
	for _, ti := range s.TypeInfo() {
		segments = append(segments, ti.Markdown())
	}
	return strings.Join(segments, " &#124; ")
}

type Type string

const (
	TypeArray   Type = "array"
	TypeBoolean Type = "boolean"
	TypeInteger Type = "integer"
	TypeNull    Type = "null"
	TypeNumber  Type = "number"
	TypeObject  Type = "object"
	TypeString  Type = "string"
)

type Types []Type

// Contains returns true if the receiver contains val.
func (ts *Types) Contains(val Type) bool {
	for _, t := range *ts {
		if t == val {
			return true
		}
	}
	return false
}

func (ts *Types) UnmarshalJSON(data []byte) error {
	var value any
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	// JSON Schema supports `type` being either a string or an array
	// of strings, so we need to normalize into a slice.
	switch val := value.(type) {
	case string:
		*ts = []Type{Type(val)}
		return nil
	case []any:
		var types []Type
		for _, t := range val {
			s, ok := t.(string)
			if !ok {
				return fmt.Errorf("unsupported type: %T for %v", t, t)
			}
			types = append(types, Type(s))
		}
		*ts = types
	default:
		return fmt.Errorf("unsupported type: %T for %v", val, val)
	}

	return nil
}

type TypeInfo struct {
	Type   Type
	Schema *Schema
}

func (ti TypeInfo) Markdown() string {
	if ti.Schema == nil {
		return string(ti.Type)
	}
	switch ti.Type {
	case TypeArray:
		switch len(ti.Schema.Types) {
		case 0:
			// array
			return string(ti.Type)
		case 1:
			// string[]
			// Widget[]
			// etc...
			return fmt.Sprintf("%s[]", ti.Schema.TypeInfoMarkdown())
		default:
			// (string | number | Widget)[]
			return fmt.Sprintf("(%s)[]", ti.Schema.TypeInfoMarkdown())
		}
	case TypeObject:
		if ti.Schema.EntityName() == "" {
			return string(ti.Type)
		}
		return fmt.Sprintf("[%s](%s)", ti.Schema.EntityName(), ti.Schema.EntityLink())
	default:
		return string(ti.Type)
	}
}

type Any struct {
	value any
}

func (a Any) MarshalJSON() ([]byte, error) {
	return json.Marshal(a.value)
}

func (a *Any) UnmarshalJSON(data []byte) error {
	return json.Unmarshal(data, &a.value)
}

func (a Any) String() string {
	if a.value == nil {
		return ""
	}
	return fmt.Sprintf("%v", a.value)
}

func (a Any) JSONString() string {
	if a.value == nil {
		return ""
	}
	serialized, err := json.Marshal(a.value)
	if err != nil {
		return err.Error()
	}
	return strings.TrimSpace(string(serialized))
}

func (a Any) YAMLString() string {
	if a.value == nil {
		return ""
	}
	serialized, err := yaml.Marshal(a.value)
	if err != nil {
		return err.Error()
	}
	return strings.TrimSpace(string(serialized))
}
