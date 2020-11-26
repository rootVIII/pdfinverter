package pdfinverter

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

// PDFInverter provides an interface to the CLI and GUI types.
type PDFInverter interface {
	imageRoutine(imgName string, fin chan<- struct{})
	extractImage()
	iterImage(imgName string)
	writePDF()
	RunApp()
}

// App controls CLI and GUI application startup & processing.
type App struct {
	TmpDir, PyPNGToPDF, PDFIn, PDFOut string
	imgCount                          int
	Executor
}

// ImageRoutine inverts the image within a goroutine.
func (app *App) imageRoutine(imgName string, fin chan<- struct{}) {
	app.iterImage(imgName)
	fin <- struct{}{}
}

// ExtractImage extracts PNG images from the input PDF.
func (app App) extractImage() {
	doc, err := fitz.New(app.PDFIn)
	if err != nil {
		panic(err)
	}

	defer doc.Close()

	for pageCount := 0; pageCount < doc.NumPage(); pageCount++ {
		currentImg, err := doc.Image(pageCount)
		if err != nil {
			panic(err)
		}
		writePNG(fmt.Sprintf("%sout-%06d.png", app.TmpDir, pageCount), currentImg)
	}
}

// IterImage examines each row of pixels in a PNG while creating a new image.
func (app App) iterImage(imgName string) {
	pathPNG := fmt.Sprintf("%s%s", app.TmpDir, imgName)
	currentPNG := readPNG(pathPNG)
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
	writePNG(pathPNG, revised)
}

// WritePDF uses gofpdf to write images to the PDF file.
func (app *App) writePDF() {
	var paths []string
	for index := 0; index < app.imgCount; index++ {
		var indexString string = strconv.Itoa(index)
		var leadingZeroes string
		for i := 0; i < (6 - len(indexString)); i++ {
			leadingZeroes += "0"
		}
		inputPath := fmt.Sprintf("%sout-%s%s.png", app.TmpDir, leadingZeroes, indexString)
		paths = append(paths, inputPath)
	}
	cmd := fmt.Sprintf("/usr/bin/python %s %s", app.PyPNGToPDF, strings.Join(paths, " "))
	app.SetCommand(cmd)
	app.RunCommand()
	err := os.Rename(app.TmpDir+"aggr.pdf", app.PDFOut)
	if err != nil {
		panic(err)
	}
}
