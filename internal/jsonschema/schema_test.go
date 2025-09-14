package jsonschema

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/twelvelabs/termite/render"
)

func TestSchema_BaseURI(t *testing.T) {
	require := require.New(t)

	schema := Schema{}
	require.Equal("", schema.BaseURI().String())

	schema = Schema{
		RetrievalURI: "https://example.com/schema.json",
	}
	require.Equal("https://example.com/schema.json", schema.BaseURI().String())

	schema = Schema{
		RetrievalURI: "https://example.com/schema.json",
		ID:           "/foo/bar.json",
	}
	require.Equal("https://example.com/foo/bar.json", schema.BaseURI().String())

	schema = Schema{
		RetrievalURI: "https://example.com/schema.json",
		ID:           "https://example.com/foo/bar.json",
	}
	require.Equal("https://example.com/foo/bar.json", schema.BaseURI().String())

	require.Panics(func() {
		schema = Schema{
			RetrievalURI: "\n",
		}
		schema.BaseURI()
	})

	require.Panics(func() {
		schema = Schema{
			RetrievalURI: "https://example.com/schema.json",
			ID:           "\n",
		}
		schema.BaseURI()
	})
}

func TestSchema_RefURI(t *testing.T) {
	require := require.New(t)

	var schema Schema

	schema = Schema{}
	require.Equal("", schema.RefURI("").String())

	schema = Schema{}
	require.Equal(
		"#/definitions/foo",
		schema.RefURI("#/definitions/foo").String(),
	)

	schema = Schema{
		ID: "https://example.com/schema.json",
	}
	require.Equal(
		"https://example.com/schema.json#/definitions/foo",
		schema.RefURI("#/definitions/foo").String(),
	)

	schema = Schema{
		ID: "https://example.com/schema.json",
	}
	require.Equal(
		"https://example.com/other.json",
		schema.RefURI("/other.json").String(),
	)

	schema = Schema{
		ID: "https://example.com/schema.json",
	}
	require.Equal(
		"https://other.com/other.json",
		schema.RefURI("https://other.com/other.json").String(),
	)

	schema = Schema{
		ID: "https://example.com/schema.json",
	}
	require.Equal(
		"https://other.com/other.json#/something",
		schema.RefURI("https://other.com/other.json#/something").String(),
	)
}

func TestSchema_DescriptionMarkdown(t *testing.T) {
	require := require.New(t)

	// Defaults to empty string.
	schema := Schema{}
	require.Equal("", schema.DescriptionMarkdown())

	// Uses `.Description` if present.
	schema = Schema{
		Description: "This _could_ be Markdown",
	}
	require.Equal("This _could_ be Markdown", schema.DescriptionMarkdown())

	// Uses `.MarkdownDescription` if present.
	schema = Schema{
		Description:         "This _could_ be Markdown",
		MarkdownDescription: "This **is** Markdown",
	}
	require.Equal("This **is** Markdown", schema.DescriptionMarkdown())
}

func TestSchema_EntityName(t *testing.T) {
	require := require.New(t)

	// Defaults to empty string.
	schema := Schema{}
	require.Equal("", schema.EntityName())

	// Uses (pascal cased) `.Ref` if present.
	schema = Schema{
		Ref: "#/definitions/foo_bar",
	}
	require.Equal("FooBar", schema.EntityName())

	// Uses (pascal cased) `.Key` if present.
	schema = Schema{
		Ref: "#/definitions/foo_bar",
		Key: "barFoo",
	}
	require.Equal("BarFoo", schema.EntityName())

	// But prefers `.Title`.
	schema = Schema{
		Ref:   "#/definitions/Foo",
		Key:   "barFoo",
		Title: "Bar",
	}
	require.Equal("Bar", schema.EntityName())
}

func TestSchema_EnumMarkdownItems(t *testing.T) {
	require := require.New(t)

	schema := Schema{}
	require.Equal([]string{}, schema.EnumMarkdownItems())

	schema = Schema{
		Const: Any{"constant value"},
		Enum: []Any{
			{"one"},
			{"two"},
		},
	}
	require.Equal([]string{
		"`\"constant value\"`",
	}, schema.EnumMarkdownItems())

	schema = Schema{
		Enum: []Any{
			{"one"},
			{"two"},
		},
	}
	require.Equal([]string{
		"`\"one\"`",
		"`\"two\"`",
	}, schema.EnumMarkdownItems())

	schema = Schema{
		Enum: []Any{
			{"one"},
			{"two"},
		},
		EnumDescriptions: []string{
			"the first number",
			"the second number",
		},
	}
	require.Equal([]string{
		"`\"one\"`: the first number",
		"`\"two\"`: the second number",
	}, schema.EnumMarkdownItems())
}

func TestSchema_YAMLExamples(t *testing.T) {
	require := require.New(t)

	schema := Schema{}
	require.Equal([]string{}, schema.YAMLExamples())

	schema = Schema{
		Examples: []Any{
			{value: "string one"},
			{value: []string{
				"slice one",
				"slice two",
			}},
		},
	}
	require.Equal([]string{
		"string one",
		"- slice one\n- slice two",
	}, schema.YAMLExamples())

	// Render w/ key when present.
	schema = Schema{
		Key: "foo",
		Examples: []Any{
			{value: "string one"},
			{value: []string{
				"slice one",
				"slice two",
			}},
		},
	}
	require.Equal([]string{
		"foo: string one",
		"foo:\n    - slice one\n    - slice two",
	}, schema.YAMLExamples())
}

func TestSchema_GenPath(t *testing.T) {
	require := require.New(t)

	schema := Schema{
		GenPathTpl: *render.MustCompile(`{{ .Title | underscore }}.md`),
		Title:      "MySchema",
	}
	require.Equal("my_schema.md", schema.GenPath())

	schema = Schema{
		GenPathTpl: *render.MustCompile(`{{ fail "boom" }}`),
	}
	require.Panics(func() {
		schema.GenPath()
	})
}

func TestSchema_GenLink(t *testing.T) {
	require := require.New(t)

	// Example scenario where all schemas are being written to the same file.
	schema1 := Schema{
		Title:      "RootSchema",
		Parent:     nil,
		GenPathTpl: *render.MustCompile(`{{ .Root.Title | underscore }}.md`),
	}
	schema2 := Schema{
		Title:  "SubSchema",
		Parent: &schema1,
	}

	require.Equal("root_schema.md#rootschema", schema1.EntityLink()) // cspell: disable-line
	require.Equal("root_schema.md#subschema", schema2.EntityLink())  // cspell: disable-line
}

func TestSchema_Merge(t *testing.T) {
	require := require.New(t)

	// Empty schemas should be a noop
	schema1 := &Schema{}
	schema2 := &Schema{}
	schema1.Merge(schema2)
	require.Equal(schema1, schema2)

	// Should only overwrite attributes that have been set in the
	// other schema (as determined by their presence in the doc).
	schema1 = &Schema{
		Title:       "schema1 title",
		Description: "schema1 description",
	}
	schema2 = &Schema{
		Description: "schema2 description",
		Document: map[string]any{
			"description": "schema2 description",
		},
		Parent: &Schema{
			Title: "schema2 parent",
		},
	}
	schema1.Merge(schema2)
	require.Equal("schema1 title", schema1.Title)
	require.Equal("schema2 description", schema1.Description)
	require.Equal("schema2 parent", schema1.Parent.Title)

	// Trigger all attributes to be set so coverage doesn't take a hit.
	schema1 = &Schema{}
	schema2 = &Schema{
		Document: map[string]any{
			"unknown attribute":     true,
			"additionalItems":       true,
			"additionalProperties":  true,
			"allOf":                 true,
			"anyOf":                 true,
			"$comment":              true,
			"const":                 true,
			"contains":              true,
			"contentEncoding":       true,
			"contentMediaType":      true,
			"default":               true,
			"definitions":           true,
			"deprecated":            true,
			"description":           true,
			"else":                  true,
			"enum":                  true,
			"enumDescriptions":      true,
			"examples":              true,
			"exclusiveMaximum":      true,
			"exclusiveMinimum":      true,
			"format":                true,
			"$id":                   true,
			"if":                    true,
			"items":                 true,
			"key":                   true,
			"markdownDescription":   true,
			"maxContains":           true,
			"maximum":               true,
			"maxItems":              true,
			"maxLength":             true,
			"maxProperties":         true,
			"minContains":           true,
			"minimum":               true,
			"minItems":              true,
			"minLength":             true,
			"minProperties":         true,
			"multipleOf":            true,
			"not":                   true,
			"oneOf":                 true,
			"pattern":               true,
			"patternProperties":     true,
			"prefixItems":           true,
			"properties":            true,
			"propertyNames":         true,
			"readOnly":              true,
			"$ref":                  true,
			"required":              true,
			"$schema":               true,
			"then":                  true,
			"title":                 true,
			"type":                  true,
			"unevaluatedProperties": true,
			"uniqueItems":           true,
			"writeOnly":             true,
		},
	}
	schema1.Merge(schema2)
}

func TestSchema_RequiredKey(t *testing.T) {
	require := require.New(t)

	schema := Schema{
		Required: []string{
			"firstName",
			"lastName",
		},
	}

	require.Equal(true, schema.RequiredKey("firstName"))
	require.Equal(true, schema.RequiredKey("lastName"))
	require.Equal(false, schema.RequiredKey("age"))
}

func TestSchema_Root(t *testing.T) {
	require := require.New(t)

	schema1 := Schema{
		Title:  "schema1",
		Parent: nil,
	}
	schema2 := Schema{
		Title:  "schema2",
		Parent: &schema1,
	}
	schema3 := Schema{
		Title:  "schema3",
		Parent: &schema2,
	}

	require.Equal("schema1", schema1.Root().Title)
	require.Equal("schema1", schema2.Root().Title)
	require.Equal("schema1", schema3.Root().Title)
}

func TestSchema_TypeInfo(t *testing.T) {
	require := require.New(t)

	schema := Schema{
		Title: "RootSchema",
		Types: Types{
			TypeObject,
			TypeNull,
		},
	}

	require.Equal(
		"[RootSchema](#rootschema) &#124; null", // cspell: disable-line
		schema.TypeInfoMarkdown(),
	)

	schema = Schema{
		Title: "RootSchema",
		Items: &Schema{
			Title: "SubSchema",
			Types: Types{
				TypeObject,
			},
		},
		Types: Types{
			TypeArray,
			TypeObject,
			TypeNull,
		},
	}

	require.Equal(
		"[SubSchema](#subschema)[] &#124; [RootSchema](#rootschema) &#124; null", // cspell: disable-line
		schema.TypeInfoMarkdown(),
	)
}

func TestTypes_Contains(t *testing.T) {
	require := require.New(t)

	types := Types{
		TypeArray,
		TypeNull,
	}
	require.Equal(true, types.Contains(TypeArray))
	require.Equal(true, types.Contains(TypeNull))
	require.Equal(false, types.Contains(TypeBoolean))
}

func TestTypes_UnmarshalJSON(t *testing.T) {
	require := require.New(t)

	types := Types{}

	err := types.UnmarshalJSON([]byte(`"object"`))
	require.NoError(err)
	require.Equal(Types{TypeObject}, types)

	err = types.UnmarshalJSON([]byte(`["array", "null"]`))
	require.NoError(err)
	require.Equal(Types{TypeArray, TypeNull}, types)

	err = types.UnmarshalJSON([]byte(`[`))
	require.Error(err)

	err = types.UnmarshalJSON([]byte(`123`))
	require.Error(err)

	err = types.UnmarshalJSON([]byte(`[123, 456]`))
	require.Error(err)
}

func TestTypeInfo_Markdown(t *testing.T) {
	require := require.New(t)

	ti := TypeInfo{
		Type:   TypeArray,
		Schema: nil,
	}
	require.Equal("array", ti.Markdown())

	ti = TypeInfo{
		Type: TypeArray,
		Schema: &Schema{
			Types: []Type{
				TypeObject,
			},
			Title: "Schema1",
		},
	}
	require.Equal("[Schema1](#schema1)[]", ti.Markdown())

	ti = TypeInfo{
		Type: TypeObject,
		Schema: &Schema{
			Title: "Schema2",
		},
	}
	require.Equal("[Schema2](#schema2)", ti.Markdown())
}

func TestAny_String(t *testing.T) {
	require := require.New(t)

	require.Equal("", (&Any{}).String())
	require.Equal("foo", (&Any{"foo"}).String())
}

func TestAny_JSONString(t *testing.T) {
	require := require.New(t)

	require.Equal("", (&Any{}).JSONString())
	require.Equal("\"foo\"", (&Any{"foo"}).JSONString())

	value := NewErrEncoder(errors.New("boom"))
	require.Contains((&Any{value}).JSONString(), "boom")
}

func TestAny_YAMLString(t *testing.T) {
	require := require.New(t)

	require.Equal("", (&Any{}).YAMLString())
	require.Equal("foo", (&Any{"foo"}).YAMLString())

	value := NewErrEncoder(errors.New("boom"))
	require.Contains((&Any{value}).YAMLString(), "boom")
}

func NewErrEncoder(err error) *ErrEncoder {
	return &ErrEncoder{err}
}

type ErrEncoder struct {
	err error
}

func (e *ErrEncoder) MarshalText() ([]byte, error) {
	return nil, e.err
}

func (e *ErrEncoder) UnmarshalText(_ []byte) error {
	return e.err
}
