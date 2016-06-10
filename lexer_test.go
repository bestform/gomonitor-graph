package main

import "testing"

func TestLexer(t *testing.T) {

	cases := []struct {
		Input    string
		Expected []item
	}{
		{"foo:bar", []item{{itemString, "foo"}, {itemDivider, ":"}, {itemString, "bar"}, {itemEOF, ""}}},
		{"\nfoo:bar", []item{{itemNewLine, "\n"}, {itemString, "foo"}, {itemDivider, ":"}, {itemString, "bar"}, {itemEOF, ""}}},
		{"foo:bar\n\n", []item{{itemString, "foo"}, {itemDivider, ":"}, {itemString, "bar"}, {itemNewLine, "\n"}, {itemNewLine, "\n"}, {itemEOF, ""}}},
	}

	for _, cs := range cases {
		l, c := NewLexer(cs.Input)
		go l.Lex()

		i := 0
		for item := range c {
			if len(cs.Expected) < i+1 {
				t.Errorf("Lexing of input '%s' failed", cs.Input)
				break
			}
			if item.typ != cs.Expected[i].typ {
				t.Errorf("Lexing of input '%s' failed. Wrong type. Expected: %q, Got: %q", cs.Input, cs.Expected[i], item)
				break
			}
			if item.content != cs.Expected[i].content {
				t.Errorf("Lexing of input '%s' failed. Wrong content", cs.Input)
				break
			}
			i++
		}
	}
}
