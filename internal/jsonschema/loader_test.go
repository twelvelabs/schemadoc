package jsonschema

import (
	"errors"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/twelvelabs/termite/api"
)

func TestLoad(t *testing.T) {
	require := require.New(t)

	loader := &mockLoader{
		LoadFunc: func(uri *url.URL) (io.ReadCloser, error) {
			return io.NopCloser(strings.NewReader("{}")), nil
		},
	}

	require.Nil(RegisteredLoader("foo"))
	RegisterLoader(loader, "foo")
	t.Cleanup(func() {
		UnregisterLoader("foo")
	})

	// Should error if location can not be parsed.
	_, err := Load("nope://some\nthing")
	require.ErrorContains(err, "parse location:")

	// Should error if unknown scheme.
	_, err = Load("nope://something")
	require.ErrorContains(err, "unknown scheme: nope")

	// Should delegate to the registered loader.
	reader, err := Load("foo://something")
	require.NoError(err)
	actual, err := io.ReadAll(reader)
	defer reader.Close()
	require.NoError(err)
	require.Equal([]byte("{}"), actual)
}

func TestLoaderRegistration(t *testing.T) {
	require := require.New(t)

	require.Nil(RegisteredLoader("foo"))
	require.Nil(RegisteredLoader("bar"))

	loader := &mockLoader{}
	RegisterLoader(loader, "foo", "bar")
	require.Equal(loader, RegisteredLoader("foo"))
	require.Equal(loader, RegisteredLoader("bar"))

	UnregisterLoader("foo")
	UnregisterLoader("bar")
	require.Nil(RegisteredLoader("foo"))
	require.Nil(RegisteredLoader("bar"))
}

func TestFileLoader(t *testing.T) {
	require := require.New(t)

	path := filepath.Join("testdata", "basic.schema.json")

	uri, err := url.Parse(path)
	require.NoError(err)
	require.Equal("", uri.Scheme)
	require.Equal(path, uri.Path)

	reader, err := NewFileLoader().Load(uri)
	require.NoError(err)

	expected, err := os.ReadFile(path)
	require.NoError(err)

	actual, err := io.ReadAll(reader)
	defer reader.Close()
	require.NoError(err)
	require.Equal(expected, actual)
}

func TestHTTPLoader(t *testing.T) {
	require := require.New(t)

	uri, err := url.Parse("https://example.com/foo.schema.json")
	require.NoError(err)

	// Stub out the HTTP client.
	transport := api.NewStubbedTransport()
	transport.RegisterStub(
		api.MatchGet("/foo.schema.json"),
		api.StringResponse(`{}`),
	)
	defer transport.VerifyStubs(t)
	client := &http.Client{}
	client.Transport = transport

	// Load.
	reader, err := NewHTTPLoader(client).Load(uri)
	require.NoError(err)
	actual, err := io.ReadAll(reader)
	defer reader.Close()
	require.NoError(err)
	require.Equal([]byte(`{}`), actual)
}

func TestHTTPLoader_When404(t *testing.T) {
	require := require.New(t)

	uri, err := url.Parse("https://example.com/foo.schema.json")
	require.NoError(err)

	// Stub out the HTTP client.
	transport := api.NewStubbedTransport()
	transport.RegisterStub(
		api.MatchGet("/foo.schema.json"),
		api.WithStatus(404,
			api.StringResponse(`{}`)),
	)
	defer transport.VerifyStubs(t)
	client := &http.Client{}
	client.Transport = transport

	// Load.
	reader, err := NewHTTPLoader(client).Load(uri)
	require.ErrorContains(err, "returned status code 404")
	require.Nil(reader)
}

func TestHTTPLoader_WhenHTTPError(t *testing.T) {
	require := require.New(t)

	uri, err := url.Parse("https://example.com/foo.schema.json")
	require.NoError(err)

	// Stub out the HTTP client.
	transport := api.NewStubbedTransport()
	transport.RegisterStub(
		api.MatchGet("/foo.schema.json"),
		api.ErrorResponse(errors.New("boom")),
	)
	defer transport.VerifyStubs(t)
	client := &http.Client{}
	client.Transport = transport

	// Load.
	reader, err := NewHTTPLoader(client).Load(uri)
	require.ErrorContains(err, "boom")
	require.Nil(reader)
}

func TestCachedLoader(t *testing.T) {
	require := require.New(t)

	uri, err := url.Parse("foo.schema.json")
	require.NoError(err)

	mReadCloser := &mockReadCloser{
		ReadFunc: func(p []byte) (int, error) {
			copy(p, []byte("{}"))
			return 2, io.EOF
		},
	}
	mLoader := &mockLoader{
		LoadFunc: func(uri *url.URL) (io.ReadCloser, error) {
			return mReadCloser, nil
		},
	}
	cachedLoader := NewCachedLoader(mLoader)

	// Load.
	reader, err := cachedLoader.Load(uri)
	require.NoError(err)
	actual, err := io.ReadAll(reader)
	defer reader.Close()
	require.NoError(err)
	require.Equal([]byte(`{}`), actual)

	// Reload - should return cached content.
	reader, err = cachedLoader.Load(uri)
	require.NoError(err)
	actual, err = io.ReadAll(reader)
	defer reader.Close()
	require.NoError(err)
	require.Equal([]byte(`{}`), actual)

	// Should have only called the wrapped loader once.
	require.Equal(1, mLoader.LoadCalls)
	require.Equal(1, mReadCloser.ReadCalls)
}

func TestCachedLoader_WhenLoadError(t *testing.T) {
	require := require.New(t)

	uri, err := url.Parse("foo.schema.json")
	require.NoError(err)

	mLoader := &mockLoader{
		LoadFunc: func(uri *url.URL) (io.ReadCloser, error) {
			return nil, errors.New("load error")
		},
	}
	cachedLoader := NewCachedLoader(mLoader)

	// Load.
	reader, err := cachedLoader.Load(uri)
	require.ErrorContains(err, "load error")
	require.Nil(reader)

	// Reload - should return cached error.
	reader, err = cachedLoader.Load(uri)
	require.ErrorContains(err, "load error")
	require.Nil(reader)

	// Should have only called the wrapped loader once.
	require.Equal(1, mLoader.LoadCalls)
}

func TestCachedLoader_WhenReadError(t *testing.T) {
	require := require.New(t)

	uri, err := url.Parse("foo.schema.json")
	require.NoError(err)

	mReadCloser := &mockReadCloser{
		ReadFunc: func(p []byte) (int, error) {
			return 0, errors.New("read error")
		},
	}
	mLoader := &mockLoader{
		LoadFunc: func(uri *url.URL) (io.ReadCloser, error) {
			return mReadCloser, nil
		},
	}
	cachedLoader := NewCachedLoader(mLoader)

	// Load.
	reader, err := cachedLoader.Load(uri)
	require.ErrorContains(err, "read error")
	require.Nil(reader)

	// Reload - should return cached error.
	reader, err = cachedLoader.Load(uri)
	require.ErrorContains(err, "read error")
	require.Nil(reader)

	// Should have only called the wrapped loader once.
	require.Equal(1, mLoader.LoadCalls)
	require.Equal(1, mReadCloser.ReadCalls)
}

type mockLoader struct {
	LoadCalls int
	LoadFunc  func(uri *url.URL) (io.ReadCloser, error)
}

func (l *mockLoader) Load(uri *url.URL) (io.ReadCloser, error) {
	l.LoadCalls++
	return l.LoadFunc(uri)
}

type mockReadCloser struct {
	ReadCalls int
	ReadFunc  func(p []byte) (n int, err error)
}

func (r *mockReadCloser) Read(p []byte) (int, error) {
	r.ReadCalls++
	return r.ReadFunc(p)
}
func (r *mockReadCloser) Close() error {
	return nil
}
