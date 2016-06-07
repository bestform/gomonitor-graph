package main

import (
	"flag"
	"fmt"
	"image/color"
	"log"
	"os"
	"strconv"

	"github.com/vdobler/chart"
)

func main() {
	f := flag.String("file", "", "the file to parse")
	n := flag.String("name", "", "the data to parse")
	flag.Parse()
	if *f == "" {
		fmt.Println("Please provide a filename")
		return
	}
	if *n == "" {
		fmt.Println("Please provide a data set name")
		return
	}

	var parser SimpleParser

	file, err := os.Open(*f)
	defer file.Close()

	if err != nil {
		log.Fatal(err)
	}

	data, err := parser.parse(*n, file)
	if err != nil {
		log.Fatal(err)
	}
	plotData(data)
}

func plotData(data MonitorData) {
	dumper := NewDumper(data.Title, 2, 2, 400, 300)
	defer dumper.Close()

	c := chart.ScatterChart{}
	c.Key.Hide = true
	c.XRange.TicSetting.Hide = true
	c.AddData("", convertData(data.Values), chart.PlotStyleLines, chart.Style{Symbol: 'o', SymbolColor: color.NRGBA{0x00, 0xaa, 0x00, 0xff}})
	c.Title = data.Title
	c.Key.Pos = "icr"
	dumper.Plot(&c)

}

func convertData(data []string) []chart.EPoint {
	points := make([]chart.EPoint, len(data))
	for x, d := range data {
		v, err := strconv.ParseFloat(d, 64)
		if err != nil {
			log.Fatal(err)
		}
		points[x].X = float64(x)
		points[x].Y = v
	}

	return points
}
