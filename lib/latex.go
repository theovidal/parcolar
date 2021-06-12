package lib

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
)

const LatexHeader = `\documentclass[preview, border=5pt, 12pt]{standalone} \usepackage{pgfplots} \pgfplotsset{compat = newest} \usepackage{amsmath} \usepackage{amssymb}\usepackage[utf8x]{inputenc} \usepackage{xcolor} \pagecolor{white} \everymath{\displaystyle}`

func GenerateLatex(name string, expression string) (pngPath string, file *os.File, err error) {
	filename := "latex-" + name
	pdfPath := fmt.Sprintf("%s/%s.pdf", TempDir, filename)
	pngPath = fmt.Sprintf("%s/%s.png", TempDir, filename)
	defer os.Remove(pdfPath)

	pdflatex := exec.Command(
		"pdflatex",
		"-jobname="+filename,
		fmt.Sprintf(
			"%s \\begin{document} %s \\end{document}",
			LatexHeader,
			expression,
		),
	)
	var out bytes.Buffer
	pdflatex.Stdout = &out
	err = pdflatex.Run()
	if os.Getenv("ENV") == "dev" {
		log.Println(out.String())
	}
	os.Remove(filename + ".aux")
	os.Remove(filename + ".log")
	if err != nil {
		return
	}

	convert := exec.Command(
		"convert",
		"-density", "500",
		"-quality", "90",
		pdfPath, pngPath,
	)
	out = bytes.Buffer{}
	convert.Stdout = &out
	err = convert.Run()
	if os.Getenv("ENV") == "dev" {
		log.Println(out.String())
	}
	if err != nil {
		log.Panicln("Error executing the convert imagick command: " + err.Error())
	}

	file, err = os.Open(pngPath)
	if err != nil {
		log.Panicln("Error reading generated file: " + err.Error())
	}
	return
}
