package jsonschema

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sync"
)

func init() {
	// Register the default loaders.
	// Package consumers are free to unregister or replace
	// if they need something different.
	RegisterLoader(
		NewCachedLoader(NewFileLoader()),
		"", // missing schemes are assumed to be `file`.
		"file",
	)
	RegisterLoader(
		NewCachedLoader(NewHTTPLoader(http.DefaultClient)),
		"http",
		"https",
	)
}

var (
	loaders   = map[string]Loader{}
	loadersMu = sync.RWMutex{}
)

// Load loads the content from location using the appropriate Loader.
func Load(location string) (io.ReadCloser, error) {
	u, err := url.Parse(location)
	if err != nil {
		return nil, fmt.Errorf("parse location: %w", err)
	}

	loader := RegisteredLoader(u.Scheme)
	if loader == nil {
		return nil, fmt.Errorf("unknown scheme: %s", u.Scheme)
	}
	return loader.Load(u)
}

// RegisterLoader registers a loader for the given schemes.
func RegisterLoader(loader Loader, schemes ...string) {
	loadersMu.Lock()
	defer loadersMu.Unlock()
	for _, scheme := range schemes {
		loaders[scheme] = loader
	}
}

// RegisteredLoader returns the loader registered for a scheme.
func RegisteredLoader(scheme string) Loader {
	loadersMu.Lock()
	defer loadersMu.Unlock()
	return loaders[scheme]
}

// UnregisterLoader unregisters the loader for the given scheme.
func UnregisterLoader(scheme string) {
	loadersMu.Lock()
	defer loadersMu.Unlock()
	delete(loaders, scheme)
}

/****************************************
* Loader
****************************************/

// Loader is a type that can load schemas from a URI.
type Loader interface {
	// Load returns an [io.ReadCloser] for the content at the given URI.
	Load(uri *url.URL) (io.ReadCloser, error)
}

/****************************************
* FileLoader
****************************************/

// NewFileLoader returns a new FileLoader.
func NewFileLoader() *FileLoader {
	return &FileLoader{}
}

// FileLoader loads files from the filesystem.
type FileLoader struct {
}

// Load opens the given URI.
func (l *FileLoader) Load(uri *url.URL) (io.ReadCloser, error) {
	return os.Open(filepath.FromSlash(uri.Path))
}

/****************************************
* HTTPLoader
****************************************/

// NewHTTPLoader returns a new HTTPLoader.
func NewHTTPLoader(client *http.Client) *HTTPLoader {
	return &HTTPLoader{
		client: client,
	}
}

// HTTPLoader loads files via [net/http].
type HTTPLoader struct {
	client *http.Client
}

// Load requests the given URI via HTTP GET.
func (l *HTTPLoader) Load(uri *url.URL) (io.ReadCloser, error) {
	url := uri.String()

	// No way to trigger an error given these params.
	request, _ := http.NewRequestWithContext(
		context.Background(), http.MethodGet, url, nil,
	)
	response, err := l.client.Do(request)
	if err != nil {
		return nil, err
	}
	if response.StatusCode != http.StatusOK {
		_ = response.Body.Close()
		return nil, fmt.Errorf("%s returned status code %d", url, response.StatusCode)
	}

	return response.Body, nil
}

/****************************************
* CachedLoader
****************************************/

// NewCachedLoader wraps the given loader with one that
// caches results in-memory.
func NewCachedLoader(loader Loader) *CachedLoader {
	return &CachedLoader{
		loader:  loader,
		results: map[string]cachedLoadResult{},
	}
}

// CachedLoader caches the results of another Loader in-memory.
type CachedLoader struct {
	loader  Loader
	results map[string]cachedLoadResult
}

// Load delegates to the wrapped loader and caches the result.
func (l *CachedLoader) Load(uri *url.URL) (io.ReadCloser, error) {
	key := uri.String()

	result, ok := l.results[key]
	if !ok {
		result = l.newResultFor(uri)
		l.results[key] = result
	}

	if result.buf == nil {
		return nil, result.err
	}
	return io.NopCloser(bytes.NewReader(result.buf)), result.err
}

// newResultFor loads the URI and returns a new cachedLoadResult.
func (l *CachedLoader) newResultFor(uri *url.URL) cachedLoadResult {
	reader, err := l.loader.Load(uri)
	if err != nil {
		return cachedLoadResult{buf: nil, err: err}
	}
	defer reader.Close()

	buf, err := io.ReadAll(reader)
	if err != nil {
		return cachedLoadResult{buf: nil, err: err}
	}

	return cachedLoadResult{buf: buf, err: nil}
}

type cachedLoadResult struct {
	buf []byte
	err error
}
