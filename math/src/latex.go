package src

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/theovidal/parcolar/lib"
)

// LatexHeader contains the default header for all LaTeX outputs
const LatexHeader = `\documentclass[preview, border=5pt, 12pt]{standalone} \usepackage{pgfplots} \pgfplotsset{compat = newest} \usepackage{amsmath} \usepackage{amssymb}\usepackage[utf8x]{inputenc} \usepackage{xcolor} \everymath{\displaystyle}`

// GenerateLatex renders a LaTeX expression into a PNG file
func GenerateLatex(name string, heading string, expression string) (pngPath string, file *os.File, err error) {
	filename := "latex-" + name
	pdfPath := fmt.Sprintf("%s/%s.pdf", lib.TempDir, filename)
	pngPath = fmt.Sprintf("%s/%s.png", lib.TempDir, filename)
	defer os.Remove(pdfPath)

	pdflatex := exec.Command(
		"pdflatex",
		"-jobname="+filename,
		fmt.Sprintf(
			"%s %s \\begin{document} %s \\end{document}",
			LatexHeader,
			heading,
			expression,
		),
	)
	var commandResult bytes.Buffer
	pdflatex.Stdout = &commandResult
	err = pdflatex.Run()
	os.Remove(filename + ".aux")
	os.Remove(filename + ".log")
	if err != nil {
		lib.LogDebug(commandResult.String())
		pdflatexErrorRegex := regexp.MustCompile("!(.*)\\n(.*)\\n.")
		message := pdflatexErrorRegex.Find(commandResult.Bytes())
		err = fmt.Errorf("Erreur dans l'expression : `%s`", strings.TrimPrefix(string(message), "! "))
		return
	}

	convert := exec.Command(
		"convert",
		"-density", "500",
		"-quality", "90",
		pdfPath, pngPath,
	)
	convert.Stdout = &commandResult
	err = convert.Run()
	if err != nil {
		lib.LogDebug(commandResult.String())
		lib.LogError("Error executing the convert command from ImageMagick: %v", err)
		err = errors.New("Une erreur est survenue lors de l'exécution de la commande.")
		return
	}

	file, err = os.Open(pngPath)
	if err != nil {
		lib.LogError("Error reading generated file at path %s: %v", pngPath, err)
		err = errors.New("Une erreur est survenue lors de l'exécution de la commande.")
	}
	return
}
