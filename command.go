package main

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api"

	"github.com/theovidal/parcolar/info"
	"github.com/theovidal/parcolar/lib"
	"github.com/theovidal/parcolar/math"
	"github.com/theovidal/parcolar/pronote"
	"github.com/theovidal/parcolar/wolfram"
)

// commandsList stores the commands available on the Telegram bot
var commandsList = map[string]lib.Command{
	// ―――――― Default ――――――
	"help":  HelpCommand(),
	"start": HelpCommand(),

	// ―――――― Information ――――――
	"definition":    info.DefinitionCommand(),
	"translate":     info.TranslateCommand(),
	"wordreference": info.WordReferenceCommand(),

	"wolfram": wolfram.Command(),

	// ―――――― Mathematics ――――――
	"calc":  math.CalcCommand(),
	"latex": math.LatexCommand(),
	"plot":  math.PlotCommand(),

	// ―――――― PRONOTE ――――――
	"contents":       pronote.ContentsCommand(),
	"homework":       pronote.HomeworkCommand(),
	"timetable":      pronote.TimetableCommand(),
	"timetablechart": pronote.TimetableChartCommand(),
}

// HandleCommand parses an incoming request to execute a bot command
func HandleCommand(bot *telegram.BotAPI, update telegram.Update, isCallback bool) (err error) {
	var commandName string
	var args []string
	var chatID int64

	if isCallback {
		parts := strings.Split(strings.TrimPrefix(update.CallbackQuery.Data, "/"), " ")
		commandName = parts[0]
		args = parts[1:]
		chatID = update.CallbackQuery.Message.Chat.ID
	} else {
		commandName = update.Message.Command()
		if update.Message.CommandArguments() != "" {
			args = strings.Split(update.Message.CommandArguments(), " ")
		}
		chatID = update.Message.Chat.ID
	}

	command, exists := commandsList[commandName]
	if !exists {
		return lib.Error(bot, chatID, "Oups, il semble que cette commande soit inconnue!")
	}

	var flags map[string]interface{}
	args, flags, err = ParseFlags(args, command.Flags)
	if err != nil {
		return lib.Error(bot, chatID, err.Error())
	}

	return command.Execute(bot, &update, chatID, args, flags)
}

// ParseFlags extracts flags at the beginning of the command, holding customizable parameters
func ParseFlags(args []string, flags map[string]lib.Flag) (parsedArgs []string, parsedFlags map[string]interface{}, err error) {
	parsedFlags = make(map[string]interface{})

	if flags == nil || len(flags) == 0 {
		for name, flag := range flags {
			parsedFlags[name] = flag.Value
		}
		return args, parsedFlags, nil
	}

	for index, arg := range args {
		if !strings.Contains(arg, "=") {
			parsedArgs = append(parsedArgs, args[index])
			continue
		}

		parts := strings.Split(arg, "=")
		name := parts[0]
		flag, found := flags[name]
		if !found {
			parsedArgs = append(parsedArgs, args[index])
			continue
		}

		var value interface{}
		switch flag.Value.(type) {
		case float64:
			value, err = strconv.ParseFloat(parts[1], 64)
			if err != nil {
				err = fmt.Errorf("Le paramètre `%s` attend un nombre réel comme valeur.", name)
				return
			}
		case int, bool:
			value, err = strconv.Atoi(parts[1])
			if err != nil {
				err = fmt.Errorf("Le paramètre `%s` attend un nombre entier comme valeur.", name)
				return
			}

			if reflect.TypeOf(flag.Value).Name() == "bool" {
				if value == 0 {
					value = false
				} else if value == 1 {
					value = true
				} else {
					err = fmt.Errorf("Le paramètre `%s` attend un booléen (0 ou 1) comme valeur.", name)
					return
				}
			}
		case string:
			value = parts[1]
			if exists := lib.Contains(*flag.Enum, value.(string)); !exists && flag.Enum != nil && len(*flag.Enum) > 0 {
				err = fmt.Errorf("Les valeurs acceptées pour le paramètre `%s` sont : %s.", name, strings.Join(*flag.Enum, ", "))
				return
			}
		default:
			lib.Fatal("Unhandled type for flag %s", name)
		}

		parsedFlags[name] = value
	}

	for flag, defaultFlag := range flags {
		if _, set := parsedFlags[flag]; !set {
			parsedFlags[flag] = defaultFlag.Value
		}
	}

	return
}
