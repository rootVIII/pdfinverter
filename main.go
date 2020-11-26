package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/google/uuid"
)

// AppInitializer provides an interface to the CLI and GUI types.
type AppInitializer interface {
	ImageRoutine(imgName string, fin chan<- struct{})
	ExtractImage()
	IterImage(imgName string)
	WritePDF()
	RunApp()
}

// RunCLI is the entry point to the cmd-line version.
func RunCLI(tmpDir string, pngtopdf string) {
	inputFile := flag.String("i", "", "input file path")
	outputFile := flag.String("o", "", "output file path")
	flag.Parse()
	if len(*inputFile) < 1 || len(*outputFile) < 1 {
		fmt.Println("-i <intput PDF path> -o <output PDF path> are required")
	} else if _, err := os.Stat(*inputFile); err != nil {
		fmt.Println("invalid file path provided for -i <input>")
	} else {
		var cliInit AppInitializer
		cliInit = &CLI{
			PDFInverter: PDFInverter{
				TmpDir:     tmpDir,
				PDFIn:      *inputFile,
				PDFOut:     *outputFile,
				PyPNGToPDF: pngtopdf,
				ImgCount:   0,
			},
		}
		cliInit.RunApp()
	}
}

// RunGUI runs the program with a QT front-end..
func RunGUI(tmpDir string, pngtopdf string) {

	var guiInit AppInitializer
	guiInit = &GUI{
		PDFInverter: PDFInverter{
			TmpDir:     tmpDir,
			PyPNGToPDF: pngtopdf,
			ImgCount:   0,
		},
	}
	guiInit.RunApp()
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

	tmpdir, err := ioutil.TempDir("", randPrefix.String())
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(tmpdir)

	randFileName, err := uuid.NewRandom()
	if err != nil {
		log.Fatal(err)
	}

	pngtopdfTMP := filepath.Join(tmpdir, randFileName.String())
	tmpdir += "/"
	WriteText(pngtopdfTMP, GetPDFConv())

	if len(os.Args) > 1 {
		RunCLI(tmpdir, pngtopdfTMP)
	} else {
		RunGUI(tmpdir, pngtopdfTMP)
	}
}
