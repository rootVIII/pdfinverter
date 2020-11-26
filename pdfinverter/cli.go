package pdfinverter

// rootVIII 2020

import (
	"io/ioutil"
	"strings"
)

// CLI inherits all types and controls CLI application startup & processing.
type CLI struct {
	App
}

// RunApp inverts a pdf based on cmd-line arguments.
func (cli *CLI) RunApp() {
	cli.extractImage()
	files, _ := ioutil.ReadDir(cli.TmpDir)
	for _, batch := range chunk(files) {
		ch := make(chan struct{})
		routines := 0
		for _, fileName := range batch {
			if !strings.Contains(fileName, "out-") {
				continue
			}
			go cli.imageRoutine(fileName, ch)
			cli.imgCount++
			routines++
		}
		for i := 0; i < routines; i++ {
			<-ch
		}
		routines = 0
	}

	cli.writePDF()

}
