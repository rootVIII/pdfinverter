package main

import (
	"bytes"
	"image"
	"image/png"
	"io/ioutil"
	"os"
)

// WritePNG writes an inverted PNG to disk.
func WritePNG(path string, newIMG image.Image) {
	buf := &bytes.Buffer{}
	err := png.Encode(buf, newIMG)
	if err != nil {
		panic(err)
	} else {
		err = ioutil.WriteFile(path, buf.Bytes(), 0600)
		if err != nil {
			panic(err)
		}
	}
}

// ReadPNG reads the image to be inverted.
func ReadPNG(path string) image.Image {
	imgRaw, err := os.Open(path)
	defer imgRaw.Close()
	if err != nil {
		panic(err)
	}
	imgDecoded, err := png.Decode(imgRaw)
	if err != nil {
		panic(err)
	}
	return imgDecoded
}

// WriteText creates a script that is used to convert PNGs into a PDF.
func WriteText(writePath string, text []byte) {
	err := ioutil.WriteFile(writePath, text, 0700)
	if err != nil {
		panic(err)
	}
}

// Chunk breaks a slice of file names into evenly sized slices. The
// final slice will contain the remaining filenames.
func Chunk(fileNames []os.FileInfo) [][]string {
	chunked := [][]string{}
	index, chunkSize := 0, 100

	for i := 0; i < len(fileNames)/chunkSize+1; i++ {
		section := make([]string, chunkSize)
		for j := 0; j < chunkSize && index < len(fileNames); j++ {
			section[j] = fileNames[index].Name()
			index++
		}
		chunked = append(chunked, section)
	}
	return chunked
}

// GetPDFConv returns Python code used as a utility shell script.
func GetPDFConv() []byte {
	script := []byte(`
import Quartz as Quartz
from CoreFoundation import NSImage
from os.path import realpath, basename
from sys import argv


def png_to_pdf(args):
    image = NSImage.alloc().initWithContentsOfFile_(args[0])
    page_init = Quartz.PDFPage.alloc().initWithImage_(image)
    pdf = Quartz.PDFDocument.alloc().initWithData_(page_init.dataRepresentation())

    for index, file_path in enumerate(args[1:]):
        image = NSImage.alloc().initWithContentsOfFile_(file_path)
        page_init = Quartz.PDFPage.alloc().initWithImage_(image)
        pdf.insertPage_atIndex_(page_init, index + 1)

    pdf.writeToFile_(realpath(__file__)[:-len(basename(__file__))] + 'aggr.pdf')


if __name__ == '__main__':
	png_to_pdf(argv[1:])
`)
	return bytes.ReplaceAll(script, []byte{0x09}, []byte{0x20, 0x20, 0x20, 0x20})
}
