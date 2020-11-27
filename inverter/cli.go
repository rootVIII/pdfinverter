package inverter

// rootVIII 2020

import (
	"io/ioutil"
	"strings"
	"sync"
)

// CLI embeds App type and controls CLI application startup & processing.
type CLI struct {
	App
}

// RunApp inverts a pdf based on cmd-line arguments.
func (cli *CLI) RunApp() {
	cli.extractImage()
	files, _ := ioutil.ReadDir(cli.TmpDir)
	for _, batch := range chunk(files) {
		var wg sync.WaitGroup
		for _, fileName := range batch {
			if !strings.Contains(fileName, "out-") {
				continue
			}
			wg.Add(1)
			go cli.imageRoutine(fileName, &wg)
			cli.imgCount++
		}
		wg.Wait()
	}

	cli.writePDF()

}
