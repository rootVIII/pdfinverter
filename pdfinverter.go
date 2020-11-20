package main

// rootVIII 2020

import (
	"fmt"
	"image"
	"image/color"
	"os"
	"strconv"
	"strings"

	"github.com/gen2brain/go-fitz"
)

// PDFInverter inherits all types and controls CLI application startup & processing.
type PDFInverter struct {
	TmpDir, PDFIn, PDFOut string
	ImgCount              int
	Executor
}

// ImageRoutine inverts the image within a goroutine.
func (p *PDFInverter) ImageRoutine(imgName string, fin chan<- struct{}) {
	p.IterImage(imgName)
	fin <- struct{}{}
}

// ExtractImage extracts PNG images from the input PDF.
func (p PDFInverter) ExtractImage() {
	doc, err := fitz.New(p.PDFIn)
	if err != nil {
		ExitErr(err)
	}

	defer doc.Close()

	for pageCount := 0; pageCount < doc.NumPage(); pageCount++ {
		currentImg, err := doc.Image(pageCount)
		if err != nil {
			ExitErr(err)
		}
		WritePNG(fmt.Sprintf("%sout-%06d.png", p.TmpDir, pageCount), currentImg)
	}
}

// IterImage examines each row of pixels in a PNG while creating a new image.
func (p PDFInverter) IterImage(imgName string) {
	pathPNG := fmt.Sprintf("%s%s", p.TmpDir, imgName)
	currentPNG := ReadPNG(pathPNG)
	perimeter := currentPNG.Bounds()
	width, height := perimeter.Max.X, perimeter.Max.Y
	revised := image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{width, height}})
	var currentPixel color.Color
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			red, green, blue, _ := currentPNG.At(x, y).RGBA()
			r, g, b := uint8(red), uint8(green), uint8(blue)
			if r == 0x7F && g == 0x7F && b == 0x7F {
				currentPixel = color.RGBA{0x1E, 0x1B, 0x24, 0xFF}
			} else {
				currentPixel = color.RGBA{0xFF - r, 0xFF - g, 0xFF - b, 0xFF}
			}
			revised.Set(x, y, currentPixel)
		}
	}
	WritePNG(pathPNG, revised)
}

// WritePDF uses gofpdf to write images to the PDF file.
func (p *PDFInverter) WritePDF() {
	var paths []string
	for index := 0; index < p.ImgCount; index++ {
		var indexString string = strconv.Itoa(index)
		var leadingZeroes string
		for i := 0; i < (6 - len(indexString)); i++ {
			leadingZeroes += "0"
		}
		inputPath := fmt.Sprintf("%sout-%s%s.png", p.TmpDir, leadingZeroes, indexString)
		paths = append(paths, inputPath)
	}
	cmd := fmt.Sprintf("/usr/bin/python %spngtopdf.py %s", p.TmpDir, strings.Join(paths, " "))
	p.SetCommand(cmd)
	p.RunCommand()
	err := os.Rename(p.TmpDir+"aggr.pdf", p.PDFOut)
	if err != nil {
		ExitErr(err)
	}
}
