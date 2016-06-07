package main

import (
	"strings"
	"testing"
)

func TestEmptyStringParse(t *testing.T) {
	data := ""
	p := SimpleParser{}
	_, err := p.parse("foo", strings.NewReader(data))

	if err == nil {
		t.Fatal("Empty string should result in an error")
	}
}

func TestSimpleEntry(t *testing.T) {
	data := "foo: bar"
	p := SimpleParser{}
	parsedData, err := p.parse("foo", strings.NewReader(data))

	if err != nil {
		t.Fatal("Unexpected error")
	}

	if parsedData.Title != "foo" {
		t.Fatal("wrong title")
	}

	if len(parsedData.Values) == 0 {
		t.Fatal("No data")
	}

	if parsedData.Values[0] != "bar" {
		t.Fatal("Wrong data")
	}
}

func TestSelectiveEntry(t *testing.T) {
	data := "foo: bar\nbar: baz"
	p := SimpleParser{}
	parsedData, _ := p.parse("bar", strings.NewReader(data))

	if parsedData.Title != "bar" {
		t.Fatal("wrong title")
	}

	if len(parsedData.Values) != 1 {
		t.Fatal("Wrong amount of data")
	}

	if parsedData.Values[0] != "baz" {
		t.Fatal("Wrong data")
	}
}

func TestmultipleEntries(t *testing.T) {
	data := "bar: baz1\nfoo: bar\nbar: baz2"
	p := SimpleParser{}
	parsedData, _ := p.parse("bar", strings.NewReader(data))

	if len(parsedData.Values) != 2 {
		t.Fatal("Wrong amount of data")
	}

	if parsedData.Values[0] != "baz1" || parsedData.Values[1] != "baz2" {
		t.Fatal("Wrong data")
	}
}
