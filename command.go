package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api"

	"github.com/theovidal/bacbot/info"
	"github.com/theovidal/bacbot/lib"
	"github.com/theovidal/bacbot/math"
	"github.com/theovidal/bacbot/pronote"
)

// commandsList stores the commands available on the Telegram bot
var commandsList = map[string]lib.Command{
	// ―――――― Information ――――――
	"definition":    info.DefinitionCommand(),
	"translate":     info.TranslateCommand(),
	"wordreference": info.WordReferenceCommand(),

	// ―――――― Mathematics ――――――
	"calc":  math.CalcCommand(),
	"plot":  math.PlotCommand(),
	"latex": math.LatexCommand(),

	// ―――――― Pronote ――――――
	"contents":       pronote.ContentsCommand(),
	"homework":       pronote.HomeworkCommand(),
	"timetable":      pronote.TimetableCommand(),
	"timetablechart": pronote.TimetableChartCommand(),
}

// HandleCommand parses an incoming request to execute a bot command
func HandleCommand(bot *telegram.BotAPI, update telegram.Update, isCallback bool) error {
	var commandName string
	var args []string

	if isCallback {
		parts := strings.Split(strings.TrimPrefix(update.CallbackQuery.Data, "/"), " ")
		commandName = parts[0]
		args = parts[1:]
	} else {
		commandName = update.Message.Command()
		if update.Message.CommandArguments() != "" {
			args = strings.Split(update.Message.CommandArguments(), " ")
		}
	}

	command, exists := commandsList[commandName]
	if !exists {
		_, err := bot.Send(telegram.NewMessage(update.Message.Chat.ID, "❓ Oups, il semble que cette commande soit inconnue!"))
		return err
	}

	var flags map[string]interface{}
	var err error
	args, flags, err = ParseFlags(args, command.Flags)
	if err != nil {
		return lib.Error(bot, &update, err.Error())
	}

	return command.Execute(bot, &update, args, flags)
}

// ParseFlags extracts flags at the beginning of the command, holding customizable parameters
func ParseFlags(args []string, commandFlags map[string]lib.Flag) ([]string, map[string]interface{}, error) {
	flags := make(map[string]interface{})

	if commandFlags == nil || len(commandFlags) == 0 {
		for name, flag := range commandFlags {
			flags[name] = flag.Value
		}
		return args, flags, nil
	}

	for index, arg := range args {
		if !strings.Contains(arg, "=") {
			args = args[index:]
			break
		}

		parts := strings.Split(arg, "=")
		name := parts[0]
		flag, found := commandFlags[name]
		if !found {
			return nil, nil, errors.New(fmt.Sprintf("Le paramètre `%s` est inexistant. Vérifiez son orthographe ou consultez la liste des paramètres possibles avec `/help plot`.", name))
		}

		var value interface{}
		switch flag.Value.(type) {
		case float64:
			var err error
			value, err = strconv.ParseFloat(parts[1], 64)
			if err != nil {
				return nil, nil, errors.New(fmt.Sprintf("Le paramètre `%s` attend un nombre réel comme valeur.", name))
			}
		case int:
			var err error
			value, err = strconv.Atoi(parts[1])
			if err != nil {
				return nil, nil, errors.New(fmt.Sprintf("Le paramètre `%s` attend un nombre entier comme valeur.", name))
			}
		case string:
			value = parts[1]
			if exists := lib.Contains(*flag.Enum, value.(string)); !exists && flag.Enum != nil && len(*flag.Enum) > 0 {
				return nil, nil, errors.New(fmt.Sprintf("Les valeurs acceptées pour le paramètre `%s` sont : %s.", name, strings.Join(*flag.Enum, ", ")))
			}
		default:
			panic("Unhandled type for flag " + name)
		}

		flags[name] = value
	}

	for flag, defaultFlag := range commandFlags {
		if _, set := flags[flag]; !set {
			flags[flag] = defaultFlag.Value
		}
	}

	return args, flags, nil
}
