package jsonschema

import (
	"fmt"
	"strings"
)

func WordSequence(words []string, separator string) string {
	if len(words) <= 2 {
		return strings.Join(words, " "+separator+" ")
	}
	lastWord := words[len(words)-1]
	wordSeq := strings.Join(words[:len(words)-1], ", ")
	return fmt.Sprintf("%s, %s %s", wordSeq, separator, lastWord)
}
