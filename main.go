package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api"

	"github.com/theovidal/parcolar/info"
	"github.com/theovidal/parcolar/lib"
	"github.com/theovidal/parcolar/pronote"
)

func main() {
	lib.Green.Println(" ________  ________  ________  ________  ________  ___       ________  ________     \n|\\   __  \\|\\   __  \\|\\   __  \\|\\   ____\\|\\   __  \\|\\  \\     |\\   __  \\|\\   __  \\    \n\\ \\  \\|\\  \\ \\  \\|\\  \\ \\  \\|\\  \\ \\  \\___|\\ \\  \\|\\  \\ \\  \\    \\ \\  \\|\\  \\ \\  \\|\\  \\   \n \\ \\   ____\\ \\   __  \\ \\   _  _\\ \\  \\    \\ \\  \\\\\\  \\ \\  \\    \\ \\   __  \\ \\   _  _\\  \n  \\ \\  \\___|\\ \\  \\ \\  \\ \\  \\\\  \\\\ \\  \\____\\ \\  \\\\\\  \\ \\  \\____\\ \\  \\ \\  \\ \\  \\\\  \\| \n   \\ \\__\\    \\ \\__\\ \\__\\ \\__\\\\ _\\\\ \\_______\\ \\_______\\ \\_______\\ \\__\\ \\__\\ \\__\\\\ _\\ \n    \\|__|     \\|__|\\|__|\\|__|\\|__|\\|_______|\\|_______|\\|_______|\\|__|\\|__|\\|__|\\|__|")
	lib.LoadEnv(".env")
	lib.OpenCache()

	lib.OpenDirs()
	defer os.RemoveAll(lib.TempDir)

	commandsList["help"] = HelpCommand()

	bot, err := telegram.NewBotAPI(os.Getenv("TELEGRAM_TOKEN"))
	if err != nil {
		log.Panic(err)
	}

	if os.Getenv("ENV") == "dev" {
		bot.Debug = true
	}

	log.Println(lib.Green.Sprintf("✅ Authorized on account %s", bot.Self.UserName))

	go pronote.TimetableLoop(bot)

	updateChannel := telegram.NewUpdate(0)
	updateChannel.Timeout = 60

	updates, err := bot.GetUpdatesChan(updateChannel)

	go func() {
		for update := range updates {
			if update.InlineQuery != nil {
				info.ParcoursupCommand(bot, &update)
			} else if update.Message.IsCommand() {
				if update.Message.From.UserName != os.Getenv("TELEGRAM_USER") {
					continue
				}
				err := HandleCommand(bot, update, false)
				if err != nil && bot.Debug {
					log.Println(lib.Red.Sprintf("‼ Error handling a command: %s", err))
				}
			}
		}
	}()

	// Wait here until CTRL-C or other term signal is received.
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	log.Println(lib.Cyan.Sprintf("💤 Closing down bot"))
}
