package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"time"
)

// AppInitializer provides an interface to the CLI and GUI types.
type AppInitializer interface {
	ImageRoutine(imgName string, fin chan<- struct{})
	ExtractImage()
	IterImage(imgName string)
	WritePDF()
	RunApp()
}

// CLI is the entry point to the cmd-line version.
func runCLI(tmpDir string) {
	inputFile := flag.String("i", "", "input file path")
	outputFile := flag.String("o", "", "output file path")
	flag.Parse()
	if len(*inputFile) < 1 || len(*outputFile) < 1 {
		ExitErr(fmt.Errorf("-i <intput PDF path> -o <output PDF path> are required"))
		if _, err := os.Stat(*inputFile); err != nil {
			ExitErr(fmt.Errorf("invalid file path provided for -i <input>"))
		}
	}

	var cliInit AppInitializer
	cliInit = &CLI{
		PDFInverter: PDFInverter{
			TmpDir:   tmpDir,
			PDFIn:    *inputFile,
			PDFOut:   *outputFile,
			ImgCount: 0,
		},
	}

	defer os.RemoveAll(tmpDir)
	cliInit.RunApp()
}

// GUI runs the program with a QT front-end..
func runGUI(tmpDir string) {

	var guiInit AppInitializer
	guiInit = &GUI{
		PDFInverter: PDFInverter{
			TmpDir:   tmpDir,
			ImgCount: 0,
		},
	}
	guiInit.RunApp()
	os.RemoveAll(tmpDir)
}

func run() {
	CleanDirs()
	tmp := fmt.Sprintf("/var/tmp/invertpdf--%s/", time.Now().Format("20060102150405"))
	os.Mkdir(tmp, 0700)
	WriteText(tmp+"pngtopdf.py", GetPDFConv())
	if len(os.Args) > 1 {
		runCLI(tmp)
	} else {
		runGUI(tmp)
	}
}

func main() {
	// Use system python2.7 until Apple includes Python3.
	if _, err := exec.LookPath("python"); err != nil {
		panic("System python not found")
	}
	run()
}
