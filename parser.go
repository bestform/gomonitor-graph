package main

import (
	"errors"
	"io"
	"io/ioutil"
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
