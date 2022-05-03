package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/google/uuid"
	"pdfinverter/inverter"
)

// runCLI is the entry point to the cmd-line version.
func runCLI(tmpDir string) {
	inputFile := flag.String("i", "", "input file path")
	outputFile := flag.String("o", "", "output file path")
	flag.Parse()
	if len(*inputFile) < 1 || len(*outputFile) < 1 {
		fmt.Println("-i <intput PDF path> -o <output PDF path> are required")
	} else if _, err := os.Stat(*inputFile); err != nil {
		fmt.Println("invalid file path provided for -i <input>")
	} else {
		var cliInit inverter.PDFInverter
		cliInit = &inverter.CLI{
			App: inverter.App{
				TmpDir: tmpDir,
				PDFIn:  *inputFile,
				PDFOut: *outputFile,
			},
		}
		cliInit.RunApp()
	}
}

// runGUI runs the program with a QT front-end..
func runGUI(tmpDir string) {
	var guiInit inverter.PDFInverter
	guiInit = &inverter.GUI{
		App: inverter.App{
			TmpDir: tmpDir,
		},
	}
	guiInit.RunApp()
}

func main() {
	randPrefix, err := uuid.NewRandom()
	if err != nil {
		log.Fatal(err)
	}
	tmpdir, err := ioutil.TempDir("", randPrefix.String())
	if err != nil {
		log.Fatal(err)
	}

	defer os.RemoveAll(tmpdir)

	tmpdir += "/"
	if len(os.Args) > 1 {
		runCLI(tmpdir)
	} else {
		runGUI(tmpdir)
	}
}
