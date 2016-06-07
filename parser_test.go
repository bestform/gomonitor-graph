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

func TestCases(t *testing.T) {
	cases := []struct {
		Input, Name, ExpectedTitle string
		ExpectedValues             []string
	}{
		{"foo: bar", "foo", "foo", []string{"bar"}},
	}

	for _, c := range cases {
		p := SimpleParser{}
		parsedData, err := p.parse(c.Name, strings.NewReader(c.Input))
		if err != nil {
			t.Errorf("an error occured when parsing '%s': %s", c.Input, err)
		}
		if parsedData.Title != c.ExpectedTitle {
			t.Errorf("Expected title: %s, but got: %s", c.ExpectedTitle, parsedData.Title)
		}
		if len(parsedData.Values) != len(c.ExpectedValues) {
			t.Errorf("Expected values: %v, but got: %v", c.ExpectedValues, parsedData.Values)
		} else {
			for i := range parsedData.Values {
				if parsedData.Values[i] != c.ExpectedValues[i] {
					t.Errorf("Expected values: %v, but got: %v", c.ExpectedValues, parsedData.Values)
				}
			}
		}

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
