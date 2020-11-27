package inverter

// rootVIII 2020

import (
	"fmt"
	"image"
	"image/color"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/gen2brain/go-fitz"
)

// PDFInverter provides an interface to the CLI and GUI types.
type PDFInverter interface {
	imageRoutine(imgName string, wg *sync.WaitGroup)
	extractImage()
	iterImage(imgName string)
	writePDF()
	RunApp()
}

// App controls CLI and GUI application startup & processing.
type App struct {
	TmpDir, PyPNGToPDF, PDFIn, PDFOut string
	imgCount                          int
	executor
}

// imageRoutine inverts the image within a goroutine.
func (app *App) imageRoutine(imgName string, wg *sync.WaitGroup) {
	defer wg.Done()
	app.iterImage(imgName)
}

// extractImage extracts PNG images from the input PDF using fitz.
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

// iterImage examines each row of pixels in a PNG while creating a new image.
func (app App) iterImage(imgName string) {
	pathPNG := fmt.Sprintf("%s%s", app.TmpDir, imgName)
	currentPNG := readPNG(pathPNG)
	perimeter := currentPNG.Bounds()
	width, height := perimeter.Max.X, perimeter.Max.Y
	revised := image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{width, height}})
	var currentPixel color.Color
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			red, green, blue, alpha := currentPNG.At(x, y).RGBA()
			r, g, b, a := uint8(red), uint8(green), uint8(blue), uint8(alpha)
			if r == 0x7F && g == 0x7F && b == 0x7F {
				currentPixel = color.RGBA{0x1E, 0x1B, 0x24, a}
			} else {
				currentPixel = color.RGBA{0xFF - r, 0xFF - g, 0xFF - b, a}
			}
			revised.Set(x, y, currentPixel)
		}
	}
	writePNG(pathPNG, revised)
}

// writePDF uses write images to the PDF file.
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
	app.setCommand(cmd)
	app.runCommand()
	err := os.Rename(app.TmpDir+"aggr.pdf", app.PDFOut)
	if err != nil {
		panic(err)
	}
}
