package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/rootVIII/pdfinverter/inverter"
)

// runCLI is the entry point to the cmd-line version.
func runCLI(tmpDir string, pngtopdf string) {
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
				TmpDir:     tmpDir,
				PDFIn:      *inputFile,
				PDFOut:     *outputFile,
				PyPNGToPDF: pngtopdf,
			},
		}
		cliInit.RunApp()
	}
}

// runGUI runs the program with a QT front-end..
func runGUI(tmpDir string, pngtopdf string) {

	var guiInit inverter.PDFInverter
	guiInit = &inverter.GUI{
		App: inverter.App{
			TmpDir:     tmpDir,
			PyPNGToPDF: pngtopdf,
		},
	}
	guiInit.RunApp()
}

// getPDFConv returns Python code used as a utility shell script to write the final PDF.
func getPDFConv() []byte {
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

func main() {
	// Use system python2.7 until Apple includes NSImage/Quartz with Python3.
	if _, err := exec.LookPath("python"); err != nil {
		log.Fatal(fmt.Errorf("Failed to find system Python in path: %v", err))
	}

	randPrefix, err := uuid.NewRandom()
	if err != nil {
		log.Fatal(err)
	}
	randFileName, err := uuid.NewRandom()
	if err != nil {
		log.Fatal(err)
	}
	tmpdir, err := ioutil.TempDir("", randPrefix.String())
	if err != nil {
		log.Fatal(err)
	}

	defer os.RemoveAll(tmpdir)

	pngtopdfTMP := filepath.Join(tmpdir, randFileName.String())
	err = ioutil.WriteFile(pngtopdfTMP, getPDFConv(), 0700)
	if err != nil {
		panic(err)
	}

	tmpdir += "/"
	if len(os.Args) > 1 {
		runCLI(tmpdir, pngtopdfTMP)
	} else {
		runGUI(tmpdir, pngtopdfTMP)
	}
}
