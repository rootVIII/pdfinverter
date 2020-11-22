### PDFINVERTER for macOS - darken (or lighten) a PDF

PDFInverter (GUI and CLI) will create a new PDF at the specified
location from a source PDF. All colors will be inverted (original shown on left):

<img src="https://github.com/rootVIII/pdfinverter/blob/master/screenshots/inverted.png" alt="example1">


Unfortunately page links are not preserved, but this program will darken PDFs making them suitable for night reading.


A 2-3 page pdf will invert very quickly. However a 400 page PDF may take 3-4 minutes (not recommended to exceed 200 pages).


A small/simple Python script is used (instead of ImageMagick) for creating the output PDF after editing. It uses the system's Python2 interpreter with NSImage/Quartz libs. This may cause the Python launcher to open in the Dock while the program is running.

The GUI is developed with Golang QT bindings:
<img src="https://github.com/rootVIII/pdfinverter/blob/master/screenshots/gui.png" alt="example2">


###### go get and run/build yourself:
<pre>
  <code>
go get github.com/rootVIII/pdfinverter
  </code>
</pre>


###### command-line usage:
<pre>
  <code>
# Required
-i     input PDF file path
-o     output PDF file path

Note:  If no command line arguments are provided, the GUI version will open.
  </code>
</pre>


<hr>
This project was developed on macOS Big Sur 11.0.1
