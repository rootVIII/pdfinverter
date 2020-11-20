package main

// rootVIII 2020

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"time"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

// GUI inherits all types and controls GUI application startup & processing.
type GUI struct {
	PDFInverter
	Window        *widgets.QMainWindow
	InputTextBox  *widgets.QLineEdit
	OutputTextBox *widgets.QLineEdit
	StatusLabel   *widgets.QLabel
	UserInfo      *user.User
	WorkingTitle  string
	WorkingCount  uint8
	StatusCount   uint8
	HaveStatus    bool
	RunningJob    bool
}

// OpenPDFInput opens the PDF that needs to be inverted.
func (g *GUI) OpenPDFInput() {
	if !g.RunningJob {
		g.PDFIn = widgets.QFileDialog_GetOpenFileName(
			g.Window, "Open PDF", g.UserInfo.HomeDir,
			"(*.pdf)", "", widgets.QFileDialog__DontUseNativeDialog)
		g.InputTextBox.SetText(g.PDFIn)
	} else {
		g.StatusLabel.SetText("Job is currently running... please wait")
	}
}

// OpenPDFOutput sets the path for the output PDF..
func (g *GUI) OpenPDFOutput() {
	if !g.RunningJob {
		g.PDFOut = widgets.QFileDialog_GetSaveFileName(
			g.Window, "Save", g.UserInfo.HomeDir, "", "",
			widgets.QFileDialog__DontUseNativeDialog)
		g.OutputTextBox.SetText(g.PDFOut)
	} else {
		g.StatusLabel.SetText("Job is currently running... please wait")
	}
}

// Invert invokes the inherited methods to invert the PDF's colors.
func (g *GUI) Invert() {
	if !g.RunningJob {
		g.PDFIn = g.InputTextBox.Text()
		g.PDFOut = g.OutputTextBox.Text()
		err := g.shouldExecute()
		if err != nil {
			g.StatusLabel.SetText(err.Error())
		} else {
			g.RunningJob = true
		}
	} else {
		g.StatusLabel.SetText("Job is currently running... please wait")
	}

}

// ResetGUI sets the GUI fields and attributes back to default if a job is not running.
func (g *GUI) ResetGUI() {
	if !g.RunningJob {
		g.Reset()
	} else {
		g.StatusLabel.SetText("Job is currently running... please wait")
	}
}

// Reset the gui and variables to inital/empty values.
func (g *GUI) Reset() {
	g.PDFIn = ""
	g.PDFOut = ""
	g.Window.SetWindowTitle("")
	g.InputTextBox.SetText(g.PDFIn)
	g.OutputTextBox.SetText(g.PDFOut)
	g.ImgCount = 0
	g.RunningJob = false
	g.WorkingCount = 0
}

func (g GUI) shouldExecute() error {
	if len(g.PDFIn) < 5 || strings.ToLower(g.PDFIn[len(g.PDFIn)-4:]) != ".pdf" {
		return fmt.Errorf("input file must have .pdf extension and MIME type")
	}
	if len(g.PDFOut) < 1 {
		return fmt.Errorf("invalid output file path provided")
	}
	fileStat, err := os.Stat(g.PDFIn)
	if err != nil || fileStat.IsDir() {
		return fmt.Errorf("invalid file provided: %v", err)
	}
	fileStat, err = os.Stat(g.PDFOut)
	if err == nil && fileStat.IsDir() {
		return fmt.Errorf("output path must be a .pdf file")
	}
	return nil
}

// ClearStatus is a QTimer function that periodically
func (g *GUI) ClearStatus() {
	if g.StatusCount > 4 {
		g.StatusLabel.SetText("")
		g.HaveStatus = false
		g.StatusCount = 0
	}
	if len(g.StatusLabel.Text()) > 0 {
		g.HaveStatus = true
		g.StatusCount++
	}
}

// DisplayStatusRunning is a QTimer() method that
// animates the title bar during processing.
func (g *GUI) DisplayStatusRunning() {
	if g.RunningJob {
		_, err := os.Stat(g.PDFOut)
		if err != nil {
			if g.WorkingCount > 10 {
				g.WorkingCount = 0
			} else {
				g.Window.SetWindowTitle(g.WorkingTitle[:g.WorkingCount])
				g.WorkingCount++
			}
		}
	}
}

// RemovePNGs cleans up old PNGs after converting a PDF from within goroutine.
func (g GUI) RemovePNGs() {
	contents, _ := ioutil.ReadDir(g.TmpDir)
	for _, png := range contents {
		if strings.Contains(png.Name(), "out") {
			err := os.Remove(fmt.Sprintf("%s%s", g.TmpDir, png.Name()))
			if err != nil {
				fmt.Printf("%v\n", err)
			}
		}
	}
}

// RunApp runs the GUI version of the app. Check for new jobs and process
// found jobs in a single background goroutine to prevent the hanging GUI
// and possible spinning beach-ball for larger-sized PDFs.
func (g *GUI) RunApp() {
	go func() {
		for {
			if !g.RunningJob {
				time.Sleep(time.Millisecond * 500)
			} else {
				g.ExtractImage()
				ch := make(chan struct{})
				files, _ := ioutil.ReadDir(g.TmpDir)
				for _, file := range files {
					if !strings.Contains(file.Name(), "out-") {
						continue
					}
					g.ImgCount++
					go g.ImageRoutine(file.Name(), ch)
				}

				for i := 0; i < g.ImgCount; i++ {
					<-ch
				}

				g.WritePDF()
				g.ImgCount = 0
				g.StatusLabel.SetText(fmt.Sprintf("%s created", filepath.Base(g.PDFOut)))
				g.Reset()
				g.RemovePNGs()
			}
		}
	}()

	userInfo, err := user.Current()
	if err != nil {
		ExitErr(fmt.Errorf("Unable to extract username of current user"))
	}
	g.UserInfo = userInfo
	g.WorkingTitle = "working..."

	ui := widgets.NewQApplication(len(os.Args), os.Args)

	g.Window = widgets.NewQMainWindow(nil, 0)
	g.Window.SetMinimumSize2(600, 250)
	g.Window.SetMaximumSize2(600, 250)
	g.Window.SetWindowTitle("")

	h1 := widgets.NewQHBoxLayout()
	h2 := widgets.NewQHBoxLayout()
	h3 := widgets.NewQHBoxLayout()
	h4 := widgets.NewQHBoxLayout()
	h5 := widgets.NewQHBoxLayout()
	v := widgets.NewQVBoxLayout()

	title := widgets.NewQGraphicsScene(nil)
	title.AddText("P D F   I N V E R T E R", gui.NewQFont2("Menlo", 20, 1, false))
	view := widgets.NewQGraphicsView(nil)
	view.SetScene(title)

	timer1 := core.NewQTimer(g.Window)
	timer1.ConnectTimeout(func() { g.ClearStatus() })
	timer1.Start(1000)
	timer2 := core.NewQTimer(g.Window)
	timer2.ConnectTimeout(func() { g.DisplayStatusRunning() })
	timer2.Start(500)

	inputLabel := widgets.NewQLabel(nil, 0)
	inputLabel.SetText("Input PDF:")
	g.InputTextBox = widgets.NewQLineEdit(nil)
	g.InputTextBox.SetPlaceholderText("None Selected")
	g.InputTextBox.SetFixedWidth(400)
	g.InputTextBox.SetStyleSheet("color: #00FFFF")
	inputButton := widgets.NewQPushButton2("Browse", nil)
	inputButton.ConnectClicked(func(bool) { g.OpenPDFInput() })
	outputLabel := widgets.NewQLabel(nil, 0)
	outputLabel.SetText("Output PDF:")
	g.OutputTextBox = widgets.NewQLineEdit(nil)
	g.OutputTextBox.SetPlaceholderText("None Selected")
	g.OutputTextBox.SetFixedWidth(400)
	g.OutputTextBox.SetStyleSheet("color: #00FFFF")
	outputButton := widgets.NewQPushButton2("Browse", nil)
	outputButton.ConnectClicked(func(bool) { g.OpenPDFOutput() })
	resetButton := widgets.NewQPushButton2("Reset", nil)
	resetButton.ConnectClicked(func(bool) { g.ResetGUI() })
	invertButton := widgets.NewQPushButton2("Invert", nil)
	invertButton.ConnectClicked(func(bool) { g.Invert() })
	g.StatusLabel = widgets.NewQLabel(nil, 0)
	g.StatusLabel.SetText(fmt.Sprintf("Greetings %s", g.UserInfo.Username))

	h1.Layout().AddWidget(view)
	h2.Layout().AddWidget(inputLabel)
	h2.Layout().AddWidget(g.InputTextBox)
	h2.Layout().AddWidget(inputButton)
	h3.Layout().AddWidget(outputLabel)
	h3.Layout().AddWidget(g.OutputTextBox)
	h3.Layout().AddWidget(outputButton)
	h4.Layout().AddWidget(resetButton)
	h4.Layout().AddWidget(invertButton)
	h5.Layout().AddWidget(g.StatusLabel)

	for _, layout := range []*widgets.QHBoxLayout{h1, h2, h3, h4, h5} {
		v.AddLayout(layout, 0)
	}

	widget := widgets.NewQWidget(nil, 0)
	widget.SetLayout(v)
	g.Window.SetCentralWidget(widget)
	g.Window.Show()
	ui.Exec()
}
