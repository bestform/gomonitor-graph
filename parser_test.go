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
		{"foo: bar\nfoo: bar2", "foo", "foo", []string{"bar", "bar2"}},
		{"foo: bar\nfoo: bar2", "bar", "bar", []string{}},
		{"foo: bar\nbar: bar\nfoo: bar2", "foo", "foo", []string{"bar", "bar2"}},
		{"\n\nfoo: bar\n\n\n\nfoo: bar2", "foo", "foo", []string{"bar", "bar2"}},
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
