### PDFINVERTER - darken (or lighten) a PDF

PDFInverter (GUI and CLI) will create a new PDF at the specified
location from a source PDF. All colors will be inverted (original shown on left):

<img src="https://user-images.githubusercontent.com/30498791/166346009-2b635dda-3c79-4557-9a7b-20f5bb64f075.png" alt="example1"><br>
<img src="https://user-images.githubusercontent.com/30498791/166346010-9d05b846-c924-4012-9693-928eafbc2a83.png" alt="example2"><br>
<img src="https://user-images.githubusercontent.com/30498791/166346011-c470d255-602c-4379-a8bd-bc0e8a2085ed.png" alt="example3"><br>


Unfortunately page links are not preserved, but this program will darken PDFs making them suitable for night reading.


A 2-3 page PDF will invert very quickly. However a 400 page PDF may take 3-4 minutes.


This project should build on any platform with <a href="https://github.com/gographics/imagick">ImageMagick bindings</a> for Golang. <code>export CGO_CFLAGS_ALLOW='-Xpreprocessor'</code> may need to be executed to run/build.



The GUI is developed with Golang QT bindings:
<img src="https://user-images.githubusercontent.com/30498791/166346008-b40e110c-9fb9-4ca1-9434-0e1f5a330171.png" alt="example2">


###### Get the project and build:
<pre>
  <code>
git clone https://github.com/rootVIII/pdfinverter.git
cd &lt;project root&gt;
go build -o bin/pdfinverter
./bin/pdfinverter 
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
