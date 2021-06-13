package lib

import (
	"fmt"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api"
	"strings"
)

// Command describes a Telegram commands with all its information
type Command struct {
	// Name of the command as shown in Telegram UI
	Name string
	// A complete description of the command to show in the help message
	Description string
	// The list of flags that can be passed by the user
	Flags map[string]Flag
	// The command actions
	Execute func(bot *telegram.BotAPI, update *telegram.Update, args []string, flags map[string]interface{}) error
}

func (c Command) Help() (content string) {
	content = fmt.Sprintf("*―――――― Aide de la commande %s ――――――*\n%s", c.Name, c.Description)

	if len(c.Flags) > 0 {
		content += "\n\nListe des paramètres disponibles sur cette commande :\n"
		for name, flag := range c.Flags {
			content += fmt.Sprintf("• `%s` : %s _(par défaut : %v)_\n", name, flag.Description, flag.Value)
			if flag.Enum != nil && len(*flag.Enum) > 0 {
				content += fmt.Sprintf("_Valeurs possibles : %s_\n", strings.Join(*flag.Enum, ", "))
			}
		}
		content += "Les paramètres sont à ajouter en début de commande sous la forme `nom=valeur`. Veillez à respecter le type de chacun (nombre entier, réel...)"
	}

	return
}

// Flag holds the data of a command flag, that's to say an optional parameter passed at the beginning
type Flag struct {
	// A complete description of the flag to show in the help message
	Description string
	// Default value of the flag. Currently supported types: string, integer, float
	Value interface{}
	// A list of accepted strings for this flag
	Enum *[]string
}
