package markdown

import (
	"bytes"
	"fmt"
	"strings"
	"unicode"

	"github.com/tdewolff/minify/v2"
	mhtml "github.com/tdewolff/minify/v2/html"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

func newConverter() goldmark.Markdown {
	return goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,
			extension.Footnote,
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			// html.WithHardWraps(),
			html.WithXHTML(),
		),
	)
}

func newMinifier() *minify.M {
	// cspell: words mhtml
	m := minify.New()
	m.AddFunc("text/html", mhtml.Minify)
	m.Add("text/html", &mhtml.Minifier{
		KeepQuotes: true,
	})
	return m
}

// ToHTMLBytes converts the given byte array to HTML.
func ToHTMLBytes(markdown []byte) ([]byte, error) {
	converted := &bytes.Buffer{}
	if err := newConverter().Convert(markdown, converted); err != nil {
		return nil, fmt.Errorf("convert error: %w", err)
	}
	minified := &bytes.Buffer{}
	if err := newMinifier().Minify("text/html", minified, converted); err != nil {
		return nil, fmt.Errorf("minify error: %w", err)
	}
	b := minified.Bytes()
	b = bytes.ReplaceAll(b, []byte("\n"), []byte(" "))
	b = bytes.ReplaceAll(b, []byte("|"), []byte("&#124;"))
	return b, nil
}

func ToHTMLString(markdown string) (string, error) {
	buf, err := ToHTMLBytes([]byte(markdown))
	return string(buf), err
}

func WrapCode(s string) string {
	if s == "" {
		return ""
	}
	return "`" + strings.TrimSpace(s) + "`"
}

func FirstSentence(s string) string {
	if s == "" {
		return ""
	}
	sentence := strings.TrimSpace(strings.Split(s, ".")[0])

	// Only append punctuation if needed.
	runes := []rune(sentence)
	lastRune := runes[len(runes)-1]
	if !unicode.IsPunct(lastRune) {
		sentence += "."
	}

	return sentence
}
