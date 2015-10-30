.PHONY: all build

PDFLATEX := pdflatex -shell-escape

all: build


build: slides.tex
	/bin/rm -rf slides.{pdf,ps,dvi}
	$(PDFLATEX) slides.tex
	$(PDFLATEX) slides.tex

