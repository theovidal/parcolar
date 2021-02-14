package main

import (
	"errors"
	"fmt"
	"github.com/theovidal/bacbot/lib"
	"github.com/theovidal/bacbot/math"
	"strconv"
	"strings"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api"

	"github.com/theovidal/bacbot/pronote"
)

var commandsList = map[string]lib.Command{
	// ------- Mathematics -------
	"plot": math.PlotCommand(),

	// ------- Pronote -------
	"contents":  pronote.ContentsCommand(),
	"homework":  pronote.HomeworkCommand(),
	"timetable": pronote.TimetableCommand(),
}

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

func ParseFlags(args []string, commandFlags map[string]interface{}) ([]string, map[string]interface{}, error) {
	if commandFlags == nil || len(commandFlags) == 0 {
		return args, commandFlags, nil
	}

	flags := make(map[string]interface{})

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
		switch flag.(type) {
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
		default:
			panic("Unhandled type for flag " + name)
		}

		flags[name] = value
	}

	for flag, defaultValue := range commandFlags {
		if _, set := flags[flag]; !set {
			flags[flag] = defaultValue
		}
	}

	return args, flags, nil
}
