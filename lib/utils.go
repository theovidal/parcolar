package lib

import (
	"regexp"

	"github.com/fatih/color"
)

func ParseTelegramMessage(input string) (output string) {
	escape := regexp.MustCompile("[_\\*\\[\\]\\(\\)~>#\\+-=\\|{}\\.!]+")
	for _, char := range []rune(input) {
		if escape.Match([]byte(string(char))) {
			output += "\\"
		}
		output += string(char)
	}
	return
}

// Red is a tool to display Red color into the term
var Red = color.New(color.FgRed)

// Green is a tool to display Green color into the term
var Green = color.New(color.FgGreen)
