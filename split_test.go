package chamux

import (
	"bufio"
	"bytes"
	"testing"
)

const FIRST = "checking that my split function works indubitably."
const SECOND = "This is not THE + END. Still not the +EN D. Let's see..."
const THIRD = "this is THE E+ND for real."

var arr = []string{FIRST, SECOND, THIRD}

func TestSplit(t *testing.T) {
	var buf bytes.Buffer
	for _, token := range arr {
		buf.WriteString(token + SPLIT_MARKER)
	}
	s := bufio.NewScanner(&buf)
	s.Split(splitPlusEnd)
	tokens := make([]string, 0)
	for s.Scan() {
		tokens = append(tokens, s.Text())
	}
	if len(tokens) != len(arr) {
		t.FailNow()
	}
	for i, token := range tokens {
		if arr[i] != token {
			t.FailNow()
		}
	}
}
