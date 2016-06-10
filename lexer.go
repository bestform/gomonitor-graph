package main

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

type itemType int
type pos int

const (
	itemString itemType = iota
	itemDivider
	itemNewLine
	itemEOF
	itemError
)

const eof = -1

type item struct {
	typ     itemType
	content string
}

func (i item) String() string {
	var t, c string
	c = i.content
	switch i.typ {
	case itemString:
		t = "STRING"
	case itemDivider:
		t = "DIVIDER"
	case itemNewLine:
		t = "NEWLINE"
		c = "\\n"
	case itemEOF:
		t = "EOF"
	case itemError:
		t = "ERROR"
	}

	return fmt.Sprintf("[%s:'%s']", t, c)
}

type lexer struct {
	input    string
	start    pos
	pos      pos
	width    pos
	emitChan chan (item)
}

func NewLexer(input string) (lexer, chan (item)) {
	l := lexer{
		input: input,
	}
	l.emitChan = make(chan (item))
	return l, l.emitChan
}

type stateFn func(*lexer) stateFn

func (l *lexer) Lex() {
	s := stateFnText
	for s != nil {
		s = s(l)
	}

	close(l.emitChan)
}

func (l *lexer) emit(typ itemType) {
	i := item{
		typ:     typ,
		content: l.input[l.start:l.pos],
	}
	l.emitChan <- i
	l.start = l.pos
}

func (l *lexer) next() rune {
	if int(l.pos) >= len(l.input) {
		l.width = 0
		return eof
	}

	r, w := utf8.DecodeRuneInString(l.input[l.pos:])
	l.width = pos(w)
	l.pos += l.width

	return r
}

func (l *lexer) backup() {
	l.pos -= l.width
}

func (l *lexer) peek() rune {
	r := l.next()
	l.backup()

	return r
}

func (l *lexer) ignore() {
	l.start = l.pos
}

func (l *lexer) accept(valid string) bool {
	if strings.IndexRune(valid, l.next()) >= 0 {
		return true
	}
	l.backup()

	return false
}

func (l *lexer) acceptRun(valid string) {
	for strings.IndexRune(valid, l.next()) >= 0 {
	}
	l.backup()
}

func stateFnText(l *lexer) stateFn {
	for {
		n := l.next()
		if n == eof {
			l.emit(itemString)
			return stateFnEOF
		}
		if r, _ := utf8.DecodeRuneInString(":"); r == n {
			l.backup()
			l.emit(itemString)
			return stateFnDivider
		}
		if r, _ := utf8.DecodeRuneInString("\n"); r == n {
			l.backup()
			l.emit(itemString)
			return stateFnNewLine
		}
	}
}

func stateFnEOF(l *lexer) stateFn {
	l.emit(itemEOF)
	return func(l *lexer) stateFn {
		return nil
	}
}

func stateFnDivider(l *lexer) stateFn {
	l.next()
	l.emit(itemDivider)
	l.next()

	return stateFnText(l)
}

func stateFnNewLine(l *lexer) stateFn {
	l.next()
	l.emit(itemNewLine)
	l.next()

	return stateFnText(l)
}
