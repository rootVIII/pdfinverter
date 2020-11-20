### PDFINVERTER for macOS - darken (or lighten) a PDF

PDFInverter (GUI and CLI) will create a new PDF at the specified
location from a source PDF. All colors will be inverted as shown:

<img src="https://github.com/rootVIII/pdfinverter/blob/master/screenshots/inverted.png" alt="example1">


Unfortunately page links are not preserved.


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


A signed/notarized .pkg installer for macOS may also be downloaded from the <a href="https://github.com/rootVIII/pdfinverter/releases">releases</a> page.


<hr>
This project was developed on macOS Big Sur 11.0.1
