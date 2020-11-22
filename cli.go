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
	files, _ := ioutil.ReadDir(cli.TmpDir)
	for _, batch := range Chunk(files) {
		ch := make(chan struct{})
		routines := 0
		for _, fileName := range batch {
			if !strings.Contains(fileName, "out-") {
				continue
			}
			go cli.ImageRoutine(fileName, ch)
			cli.ImgCount++
			routines++
		}
		for i := 0; i < routines; i++ {
			<-ch
		}
		routines = 0
	}

	cli.WritePDF()

}
