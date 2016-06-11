package main

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"strings"
)

// MonitorData contains all values for a given, named data set
type MonitorData struct {
	// Title contains the name of the data set
	Title string
	// Values contains the corresponding values
	Values []string
}

// Parser will define a parser for gomonitor logs
type Parser interface {
	Parse(string, io.Reader) (MonitorData, error)
}

// SimpleParser will parse the simple standard logs written by gomonitor
type SimpleParser struct{}

func (s SimpleParser) parse(name string, r io.Reader) (MonitorData, error) {
	var parsedData MonitorData

	rawdata, err := ioutil.ReadAll(r)

	if err != nil {
		return parsedData, err
	}

	if len(rawdata) == 0 {
		return parsedData, errors.New("Empty input")
	}

	// do the parsey parse
	parsedData.Title = name
	lines := strings.Split(string(rawdata), "\n")

	for _, line := range lines {
		if !strings.HasPrefix(line, name) {
			continue
		}
		lineData := strings.TrimPrefix(line, name+": ")
		parsedData.Values = append(parsedData.Values, lineData)
	}

	return parsedData, nil
}

type state int

const (
	stateName state = iota
	stateValue
	stateDivider
	stateIgnoreLine
)

type LexParser struct{}

func (p LexParser) parse(name string, r io.Reader) (MonitorData, error) {
	var parsedData MonitorData

	rawdata, err := ioutil.ReadAll(r)

	if err != nil {
		return parsedData, err
	}

	l, c := NewLexer(string(rawdata))

	go l.Lex()

	// do the parsey parse
	parsedData.Title = name
	f := parseName
	for i := range c {
		f, err = f(i, &parsedData)
		if err != nil {
			log.Fatal("error while parsing input")
		}
	}

	return parsedData, nil
}

type parseFn func(item, *MonitorData) (parseFn, error)

func parseName(i item, d *MonitorData) (parseFn, error) {
	if i.typ == itemNewLine {
		return parseName, nil
	}
	if i.typ == itemString && i.content == d.Title { // todo: this is ugly. We should carry the name separatly
		return parseDivider, nil
	}
	return parseIgnoreLine, nil
}

func parseDivider(i item, d *MonitorData) (parseFn, error) {
	if i.typ != itemDivider {
		return nil, fmt.Errorf("Expected Divider")
	}
	return parseValue, nil
}

func parseValue(i item, d *MonitorData) (parseFn, error) {
	d.Values = append(d.Values, strings.TrimSpace(i.content))
	return parseIgnoreLine, nil
}

func parseIgnoreLine(i item, d *MonitorData) (parseFn, error) {
	if i.typ == itemNewLine {
		return parseName, nil
	}
	return parseIgnoreLine, nil
}
