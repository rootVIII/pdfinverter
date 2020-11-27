package inverter

// rootVIII 2020

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/therecipe/qt/core"
	"github.com/therecipe/qt/gui"
	"github.com/therecipe/qt/widgets"
)

// GUI embeds App Type and controls GUI application startup & processing.
type GUI struct {
	App
	window        *widgets.QMainWindow
	inputTextBox  *widgets.QLineEdit
	outputTextBox *widgets.QLineEdit
	statusLabel   *widgets.QLabel
	userInfo      *user.User
	workingTitle  string
	workingCount  uint8
	statusCount   uint8
	haveStatus    bool
	runningJob    bool
}

// openPDFInput opens the PDF that needs to be inverted.
func (g *GUI) openPDFInput() {
	if !g.runningJob {
		g.PDFIn = widgets.QFileDialog_GetOpenFileName(
			g.window, "Open PDF", g.userInfo.HomeDir,
			"(*.pdf)", "", widgets.QFileDialog__DontUseNativeDialog)
		g.inputTextBox.SetText(g.PDFIn)
	} else {
		g.statusLabel.SetText("Job is currently running... please wait")
	}
}

// openPDFOutput sets the path for the output PDF..
func (g *GUI) openPDFOutput() {
	if !g.runningJob {
		g.PDFOut = widgets.QFileDialog_GetSaveFileName(
			g.window, "Save", g.userInfo.HomeDir, "", "",
			widgets.QFileDialog__DontUseNativeDialog)
		g.outputTextBox.SetText(g.PDFOut)
	} else {
		g.statusLabel.SetText("Job is currently running... please wait")
	}
}

// invert signals to the background go-routine that a new job is ready to be processed.
func (g *GUI) invert() {
	if !g.runningJob {
		g.PDFIn = g.inputTextBox.Text()
		g.PDFOut = g.outputTextBox.Text()
		err := g.shouldExecute()
		if err != nil {
			g.statusLabel.SetText(err.Error())
		} else {
			g.runningJob = true
		}
	} else {
		g.statusLabel.SetText("Job is currently running... please wait")
	}

}

// resetGUI sets the GUI fields and attributes back to default if a job is not running.
// Otherwise the user is warned that a job is currently being processed.
func (g *GUI) resetGUI() {
	if !g.runningJob {
		g.reset()
	} else {
		g.statusLabel.SetText("Job is currently running... please wait")
	}
}

// reset the gui and variables to inital/empty values.
func (g *GUI) reset() {
	g.PDFIn = ""
	g.PDFOut = ""
	g.inputTextBox.SetText(g.PDFIn)
	g.outputTextBox.SetText(g.PDFOut)
	g.imgCount = 0
	g.runningJob = false
	g.workingCount = 0
	g.window.SetWindowTitle("")
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

// clearStatus is a QTimer function that periodically clears any status message after 4 seconds.
func (g *GUI) clearStatus() {
	if g.statusCount > 4 {
		g.statusLabel.SetText("")
		g.haveStatus = false
		g.statusCount = 0
	}
	if len(g.statusLabel.Text()) > 0 {
		g.haveStatus = true
		g.statusCount++
	}
}

// displayStatusRunning is a QTimer() method that
// animates the title bar during processing.
func (g *GUI) displayStatusRunning() {
	if g.runningJob {
		_, err := os.Stat(g.PDFOut)
		if err != nil {
			if g.workingCount > 10 {
				g.workingCount = 0
			} else {
				g.window.SetWindowTitle(g.workingTitle[:g.workingCount])
				g.workingCount++
			}
		}
	}
}

// removePNGs cleans up old PNGs after converting a PDF from within goroutine.
func (g GUI) removePNGs() {
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
			if !g.runningJob {
				time.Sleep(time.Millisecond * 500)
			} else {
				g.extractImage()
				files, _ := ioutil.ReadDir(g.TmpDir)
				for _, batch := range chunk(files) {
					var wg sync.WaitGroup
					for _, fileName := range batch {
						if !strings.Contains(fileName, "out-") {
							continue
						}
						wg.Add(1)
						go g.imageRoutine(fileName, &wg)
						g.imgCount++
					}
					wg.Wait()
				}

				g.writePDF()
				g.statusLabel.SetText(fmt.Sprintf("%s created", filepath.Base(g.PDFOut)))
				g.reset()
				g.removePNGs()
			}
		}
	}()

	userInfo, err := user.Current()
	if err != nil {
		panic("Unable to extract username of current user")
	}
	g.userInfo = userInfo
	g.workingTitle = "working..."

	ui := widgets.NewQApplication(len(os.Args), os.Args)

	g.window = widgets.NewQMainWindow(nil, 0)
	g.window.SetMinimumSize2(600, 250)
	g.window.SetMaximumSize2(600, 250)
	g.window.SetWindowTitle("")

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

	timer1 := core.NewQTimer(g.window)
	timer1.ConnectTimeout(func() { g.clearStatus() })
	timer1.Start(1000)
	timer2 := core.NewQTimer(g.window)
	timer2.ConnectTimeout(func() { g.displayStatusRunning() })
	timer2.Start(500)

	inputLabel := widgets.NewQLabel(nil, 0)
	inputLabel.SetText("Input PDF:")
	g.inputTextBox = widgets.NewQLineEdit(nil)
	g.inputTextBox.SetPlaceholderText("None Selected")
	g.inputTextBox.SetFixedWidth(400)
	g.inputTextBox.SetStyleSheet("color: #00FFFF")
	inputButton := widgets.NewQPushButton2("Browse", nil)
	inputButton.ConnectClicked(func(bool) { g.openPDFInput() })
	outputLabel := widgets.NewQLabel(nil, 0)
	outputLabel.SetText("Output PDF:")
	g.outputTextBox = widgets.NewQLineEdit(nil)
	g.outputTextBox.SetPlaceholderText("None Selected")
	g.outputTextBox.SetFixedWidth(400)
	g.outputTextBox.SetStyleSheet("color: #00FFFF")
	outputButton := widgets.NewQPushButton2("Browse", nil)
	outputButton.ConnectClicked(func(bool) { g.openPDFOutput() })
	resetButton := widgets.NewQPushButton2("Reset", nil)
	resetButton.ConnectClicked(func(bool) { g.resetGUI() })
	invertButton := widgets.NewQPushButton2("Invert", nil)
	invertButton.ConnectClicked(func(bool) { g.invert() })
	g.statusLabel = widgets.NewQLabel(nil, 0)
	g.statusLabel.SetText(fmt.Sprintf("Greetings %s", g.userInfo.Username))

	h1.Layout().AddWidget(view)
	h2.Layout().AddWidget(inputLabel)
	h2.Layout().AddWidget(g.inputTextBox)
	h2.Layout().AddWidget(inputButton)
	h3.Layout().AddWidget(outputLabel)
	h3.Layout().AddWidget(g.outputTextBox)
	h3.Layout().AddWidget(outputButton)
	h4.Layout().AddWidget(resetButton)
	h4.Layout().AddWidget(invertButton)
	h5.Layout().AddWidget(g.statusLabel)

	for _, layout := range []*widgets.QHBoxLayout{h1, h2, h3, h4, h5} {
		v.AddLayout(layout, 0)
	}

	widget := widgets.NewQWidget(nil, 0)
	widget.SetLayout(v)
	g.window.SetCentralWidget(widget)
	g.window.Show()
	ui.Exec()
}
