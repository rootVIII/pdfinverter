package main

// rootVIII 2020

import (
	"io/ioutil"
	"strings"
)

// CLI inherits all types and controls CLI application startup & processing.
type CLI struct {
	PDFInverter
}

// RunApp inverts a pdf based on cmd-line arguments.
func (cli *CLI) RunApp() {
	cli.ExtractImage()

	ch := make(chan struct{})

	files, _ := ioutil.ReadDir(cli.TmpDir)
	for _, file := range files {
		if !strings.Contains(file.Name(), "out-") {
			continue
		}
		cli.ImgCount++
		go cli.ImageRoutine(file.Name(), ch)
	}

	for i := 0; i < cli.ImgCount; i++ {
		<-ch
	}
	cli.WritePDF()
}
