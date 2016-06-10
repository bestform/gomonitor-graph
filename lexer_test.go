package main

import (
	"fmt"
	"testing"
)

func TestLexer(t *testing.T) {
	l, c := NewLexer("foo:bar\nbar:baz")
	go l.Lex()

	for item := range c {
		fmt.Printf("%+v\n", item)
	}

	fmt.Println("Channel closed")
}
