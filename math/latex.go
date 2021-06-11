package math

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api"

	"github.com/theovidal/bacbot/lib"
)

func LatexCommand() lib.Command {
	return lib.Command{
		Name:        "latex",
		Description: "Rendu LaTeX",
		Execute: func(bot *telegram.BotAPI, update *telegram.Update, args []string, _ map[string]interface{}) (err error) {
			expression := "$" + strings.Join(args, " ") + "$"

			// TODO: customize output, error handling, documentation on installing imagick and pdflatex

			tmp, _ := ioutil.TempDir("", "bacbot-latex")
			defer os.RemoveAll(tmp)

			workdir, err := os.Getwd()
			if err != nil {
				return errors.New("couldn't get working directory")
			}
			os.Chdir(tmp)
			defer os.Chdir(workdir)

			filename := fmt.Sprintf("%d", update.UpdateID)

			cmd := exec.Command(
				"pdflatex",
				"-jobname="+filename,
				fmt.Sprintf(
					"\\documentclass[preview, border=5pt, 12pt]{standalone} \\usepackage{amsmath} \\usepackage{amssymb}\\usepackage[utf8x]{inputenc} \\usepackage{xcolor} \\pagecolor{white} \\everymath{\\displaystyle} \\begin{document} %s \\end{document}",
					expression,
				),
			)
			var out bytes.Buffer
			cmd.Stdout = &out
			err = cmd.Run()
			fmt.Println(out.String())
			if err != nil {
				return
			}

			cmd = exec.Command("convert",
				"-density", "500",
				"-quality", "90",
				filename+".pdf", filename+".png",
			)
			out = bytes.Buffer{}
			cmd.Stdout = &out
			err = cmd.Run()
			fmt.Println(out.String())
			if err != nil {
				return
			}

			file, err := os.Open(filename + ".png")
			if err != nil {
				return
			}
			reader := telegram.FileReader{
				Name:   "expression.png",
				Reader: file,
				Size:   -1,
			}
			photo := telegram.NewPhotoUpload(update.Message.Chat.ID, reader)
			_, err = bot.Send(photo)
			return
		},
	}
}
