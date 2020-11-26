package pdfinverter

import (
	"bytes"
	"image"
	"image/png"
	"io/ioutil"
	"os"
)

// WritePNG writes an inverted PNG to disk.
func writePNG(path string, newIMG image.Image) {
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
func readPNG(path string) image.Image {
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

// Chunk breaks a slice of file names into evenly sized slices. The
// final slice will contain the remaining filenames.
func chunk(fileNames []os.FileInfo) [][]string {
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
