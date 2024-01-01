package jsonschema

import (
	"encoding/json"
	"io"
	"path"
	"strings"

	"github.com/xeipuuv/gojsonpointer"
)

func NewContext() *Context {
	return &Context{
		schemas: map[string]*Schema{},
	}
}

type Context struct {
	schemas map[string]*Schema
}

func (c *Context) Get(ref string) (*Schema, error) {
	// Lookup full ref.
	if schema := c.lookup(ref); schema != nil {
		return schema, nil
	}

	// Try looking up a sub-schema ref.
	base, frag, _ := strings.Cut(ref, "#")
	if schema := c.lookup(base); schema != nil {
		subSchema, err := c.parseSubSchema(schema, frag)
		if err != nil {
			return nil, err
		}
		c.store(ref, subSchema)
		return subSchema, nil
	}

	// Not in the store (nor a sub-schema of something in the store),
	// so we need to load and parse a new root schema.
	doc, err := c.load(ref)
	if err != nil {
		return nil, err
	}
	schema, err := c.parse(doc, nil)
	if err != nil {
		return nil, err
	}
	schema.RetrievalURI = ref

	// Store by both the retrieval ref _and_ the base URI.
	// Allows us to correctly lookup relative and/or in-schema refs
	// (which resolve against the base URI).
	c.store(ref, schema)
	c.store(schema.BaseURI().String(), schema)

	// Note that we must store _before_ resolving, otherwise resolving
	// in-schema refs will trigger an infinite loop.
	return resolveSubSchemas(schema)
}

// lookup returns the stored Schema for ref,
// or nil if the ref has yet to be loaded.
func (c *Context) lookup(ref string) *Schema {
	if schema, ok := c.schemas[ref]; ok {
		return schema
	}
	return nil
}

// store stores the schema for ref.
func (c *Context) store(ref string, schema *Schema) {
	c.schemas[ref] = schema
}

// load fetches the schema document located at ref.
func (c *Context) load(ref string) ([]byte, error) {
	reader, err := Load(ref)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	doc, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	return doc, nil
}

// parse parses and returns the given schema doc.
func (c *Context) parse(doc []byte, parent *Schema) (*Schema, error) {
	schema := Schema{}
	if err := json.Unmarshal(doc, &schema); err != nil {
		return nil, err
	}
	if err := json.Unmarshal(doc, &schema.Document); err != nil {
		return nil, err
	}
	schema.Context = c
	schema.Parent = parent
	return &schema, nil
}

// parseSubSchema returns the sub-schema located at ref.
func (c *Context) parseSubSchema(schema *Schema, ref string) (*Schema, error) {
	if ref == "" {
		return schema, nil
	}

	// Resolve the sub document and convert to JSON byte array.
	pointer, err := gojsonpointer.NewJsonPointer(ref)
	if err != nil {
		return nil, err
	}
	subDoc, _, err := pointer.Get(schema.Document)
	if err != nil {
		return nil, err
	}
	buf, err := json.Marshal(subDoc)
	if err != nil {
		return nil, err
	}

	// Parse the doc back into a `Schema` struct.
	subSchema, err := c.parse(buf, schema)
	if err != nil {
		return nil, err
	}
	// Keep track of the ref base so we can use it in `.EntityName()`.
	subSchema.Key = path.Base(ref)

	return subSchema, nil
}

// resolve recursively resolves all refs in the given schema.
func resolveSubSchemas(schema *Schema) (*Schema, error) {
	if schema.Items != nil {
		schema.Items.Parent = schema
		tmp, err := resolveRef(schema.Items)
		if err != nil {
			return nil, err
		}
		*schema.Items = *tmp
		tmp, err = resolveSubSchemas(schema.Items)
		if err != nil {
			return nil, err
		}
		*schema.Items = *tmp
	}

	for key, subSchema := range schema.Definitions {
		subSchema.Key = key
		subSchema.Parent = schema
		tmp, err := resolveRef(subSchema)
		if err != nil {
			return nil, err
		}
		*subSchema = *tmp
		tmp, err = resolveSubSchemas(subSchema)
		if err != nil {
			return nil, err
		}
		*subSchema = *tmp
	}

	for key, subSchema := range schema.Properties {
		subSchema.Key = key
		subSchema.Parent = schema
		tmp, err := resolveRef(subSchema)
		if err != nil {
			return nil, err
		}
		*subSchema = *tmp
		tmp, err = resolveSubSchemas(subSchema)
		if err != nil {
			return nil, err
		}
		*subSchema = *tmp
	}

	for _, subSchema := range schema.OneOf {
		subSchema.Parent = schema
		tmp, err := resolveRef(subSchema)
		if err != nil {
			return nil, err
		}
		*subSchema = *tmp
		tmp, err = resolveSubSchemas(subSchema)
		if err != nil {
			return nil, err
		}
		*subSchema = *tmp
	}

	return schema, nil
}

func resolveRef(schema *Schema) (*Schema, error) {
	// Already resolved.
	if schema.Resolved {
		return schema, nil
	}

	// Empty ref - nothing to do.
	if schema.Ref == "" {
		schema.Resolved = true
		return schema, nil
	}

	// Get the fully qualified ref URI.
	refURI := schema.Root().RefURI(schema.Ref)
	// Load the schema for that ref.
	loaded, err := schema.Root().Context.Get(refURI.String())
	if err != nil {
		return nil, err
	}
	// Clone the loaded schema so we don't mutate the context.
	resolved, err := loaded.Clone()
	if err != nil {
		return nil, err
	}

	// Merge in the origin schema's properties and return.
	resolved.Merge(schema)
	resolved.Resolved = true
	return resolved, nil
}
