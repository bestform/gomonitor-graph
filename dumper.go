package main

import (
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"

	"github.com/vdobler/chart"
	"github.com/vdobler/chart/imgg"
)

// dumper code taken and modified from the original example included in the library

type dumper struct {
	N, M, W, H, Cnt int
	I               *image.RGBA
	imgFile         *os.File
}

func NewDumper(name string, n, m, w, h int) *dumper {
	var err error
	dumper := dumper{N: n, M: m, W: w, H: h}

	dumper.imgFile, err = os.Create(name + ".png")
	if err != nil {
		panic(err)
	}
	dumper.I = image.NewRGBA(image.Rect(0, 0, n*w, m*h))
	bg := image.NewUniform(color.RGBA{0xff, 0xff, 0xff, 0xff})
	draw.Draw(dumper.I, dumper.I.Bounds(), bg, image.ZP, draw.Src)

	return &dumper
}
func (d *dumper) Close() {
	png.Encode(d.imgFile, d.I)
	d.imgFile.Close()
}

func (d *dumper) Plot(c chart.Chart) {
	row, col := d.Cnt/d.N, d.Cnt%d.N

	igr := imgg.AddTo(d.I, col*d.W, row*d.H, d.W, d.H, color.RGBA{0xff, 0xff, 0xff, 0xff}, nil, nil)
	c.Plot(igr)

	d.Cnt++

}
