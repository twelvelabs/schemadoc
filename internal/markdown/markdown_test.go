package markdown

import (
	"os"
	"testing"

	"github.com/sebdah/goldie/v2"
	"github.com/stretchr/testify/require"
)

func TestToHTMLBytes(t *testing.T) {
	require := require.New(t)

	filename := "simple"
	filepath := "testdata/" + filename + ".md"

	content, err := os.ReadFile(filepath)
	require.NoError(err)

	htmlBytes, err := ToHTMLBytes(content)
	require.NoError(err)

	g := goldie.New(t)
	g.Assert(t, filename, htmlBytes)

	htmlString, err := ToHTMLString(string(content))
	require.NoError(err)

	g.Assert(t, filename, []byte(htmlString))
}
