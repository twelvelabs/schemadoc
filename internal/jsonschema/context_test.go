package jsonschema

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestContext_Get(t *testing.T) {
	require := require.New(t)

	context := NewContext()

	schema, err := context.Get(filepath.Join("testdata", "unknown.json"))
	require.Error(err)
	require.Nil(schema)

	schema, err = context.Get(filepath.Join("testdata", "invalid.json"))
	require.Error(err)
	require.Nil(schema)

	path := filepath.Join("testdata", "arrays.schema.json")
	schema, err = context.Get(path)
	require.NoError(err)
	require.NotNil(schema)
	require.Equal("https://example.com/arrays.schema.json", schema.ID)
	require.Same(context, schema.Context)
	require.Nil(schema.Parent)

	// Reloading the same schema URI should return the same pointer.
	schema2, err := context.Get(filepath.Join("testdata", "arrays.schema.json"))
	require.NoError(err)
	require.NotNil(schema2)
	require.Same(schema, schema2)
	require.Same(context, schema2.Context)

	// Should be able to load an in-schema ref.
	ref := schema.RefURI("#/definitions/veggie").String()
	require.Equal("https://example.com/arrays.schema.json#/definitions/veggie", ref)
	veggie, err := context.Get(ref)
	require.NoError(err)
	require.NotNil(veggie)
	require.Same(schema, veggie.Parent)
	require.Same(context, veggie.Context)

	// Reloading the same schema URI should return the same pointer.
	veggie2, err := context.Get(ref)
	require.NoError(err)
	require.NotNil(veggie2)
	require.Same(veggie, veggie2)
}

func TestContext_parseSubSchema(t *testing.T) {
	require := require.New(t)

	context := NewContext()

	schema, err := context.Get(filepath.Join("testdata", "arrays.schema.json"))
	require.NoError(err)

	root, err := context.parseSubSchema(schema, "")
	require.NoError(err)
	require.NotNil(root)
	require.Same(schema, root)
	require.Same(context, root.Context)

	veggie, err := context.parseSubSchema(schema, "/definitions/veggie")
	require.NoError(err)
	require.NotNil(veggie)
	require.Equal("Veggie", veggie.EntityName())
	require.Equal(true, veggie.Types.Contains(TypeObject))
	require.Equal(true, veggie.RequiredKey("veggieName"))
	require.Equal(true, veggie.RequiredKey("veggieLike"))
	require.Same(schema, veggie.Parent)
	require.Same(context, veggie.Context)

	unknown, err := context.parseSubSchema(schema, "/definitions/unknown")
	require.Error(err)
	require.Nil(unknown)

	invalid, err := context.parseSubSchema(schema, "invalid path")
	require.Error(err)
	require.Nil(invalid)
}

func TestSchema_resolveRef(t *testing.T) {
	require := require.New(t)

	// Load the schema into the context.
	context := NewContext()
	loaded, err := context.Get(filepath.Join("testdata", "ref.schema.json"))
	require.NoError(err)
	require.NotNil(loaded)

	// Manually recreate the schema that was loaded so we can test
	// resolving refs (the `Get` call above does so automatically).
	schema := &Schema{
		ID:       loaded.ID,
		Title:    loaded.Title,
		Ref:      loaded.Ref,
		Document: loaded.Document,
		Context:  context,
	}
	require.Equal("RefSchema", schema.Title)
	require.Equal("#/definitions/ConcreteSchema", schema.Ref)
	require.Equal(false, schema.Resolved)
	require.Equal(false, schema.Types.Contains(TypeObject))
	require.Equal(0, len(schema.Properties))

	schema, err = resolveRef(schema)
	require.NoError(err)
	require.NotNil(schema)
	require.Equal("RefSchema", schema.Title)
	require.Equal("#/definitions/ConcreteSchema", schema.Ref)
	require.Equal(true, schema.Resolved)
	require.Equal(true, schema.Types.Contains(TypeObject))
	require.Equal(1, len(schema.Properties))
	require.Equal(true, schema.Properties["name"].Types.Contains(TypeString))
	require.Equal("The display name", schema.Properties["name"].Description)

	// Re-resolving should be a noop.
	schema, err = resolveRef(schema)
	require.NoError(err)
	require.NotNil(schema)
	require.Equal(true, schema.Resolved)

	schema3 := &Schema{
		ID:       loaded.ID,
		Title:    loaded.Title,
		Ref:      "#invalid-json-pointer",
		Document: loaded.Document,
		Context:  context,
	}
	schema3, err = resolveRef(schema3)
	require.Error(err)
	require.Nil(schema3)

	schema4 := &Schema{
		ID:       loaded.ID,
		Title:    loaded.Title,
		Ref:      "",
		Document: loaded.Document,
		Context:  context,
	}
	schema4, err = resolveRef(schema4)
	require.NoError(err)
	require.NotNil(schema4)
	require.Equal("RefSchema", schema4.Title)
	require.Equal("", schema4.Ref)
	require.Equal(true, schema4.Resolved)
	require.Equal(false, schema4.Types.Contains(TypeObject))
	require.Equal(0, len(schema4.Properties))
}
