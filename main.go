package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/fatih/color"
	telegram "github.com/go-telegram-bot-api/telegram-bot-api"

	"github.com/theovidal/parcolar/info"
	"github.com/theovidal/parcolar/lib"
	"github.com/theovidal/parcolar/pronote"
)

func main() {
	color.Green(" ________  ________  ________  ________  ________  ___       ________  ________     \n|\\   __  \\|\\   __  \\|\\   __  \\|\\   ____\\|\\   __  \\|\\  \\     |\\   __  \\|\\   __  \\    \n\\ \\  \\|\\  \\ \\  \\|\\  \\ \\  \\|\\  \\ \\  \\___|\\ \\  \\|\\  \\ \\  \\    \\ \\  \\|\\  \\ \\  \\|\\  \\   \n \\ \\   ____\\ \\   __  \\ \\   _  _\\ \\  \\    \\ \\  \\\\\\  \\ \\  \\    \\ \\   __  \\ \\   _  _\\  \n  \\ \\  \\___|\\ \\  \\ \\  \\ \\  \\\\  \\\\ \\  \\____\\ \\  \\\\\\  \\ \\  \\____\\ \\  \\ \\  \\ \\  \\\\  \\| \n   \\ \\__\\    \\ \\__\\ \\__\\ \\__\\\\ _\\\\ \\_______\\ \\_______\\ \\_______\\ \\__\\ \\__\\ \\__\\\\ _\\ \n    \\|__|     \\|__|\\|__|\\|__|\\|__|\\|_______|\\|_______|\\|_______|\\|__|\\|__|\\|__|\\|__|")
	lib.LoadEnv(".env")
	lib.OpenCache()

	lib.OpenDirs()
	defer os.RemoveAll(lib.TempDir)

	bot, err := telegram.NewBotAPI(os.Getenv("TELEGRAM_TOKEN"))
	if err != nil {
		log.Panic(err)
	}

	if os.Getenv("DEBUG") == "true" {
		bot.Debug = true
		lib.StandardLogger.Debug = true
	}

	lib.LogSuccess("Authorized on account %s", bot.Self.UserName)

	go pronote.TimetableLoop(bot)

	updateChannel := telegram.NewUpdate(0)
	updateChannel.Timeout = 60

	updates, err := bot.GetUpdatesChan(updateChannel)

	go func() {
		for update := range updates {
			if update.InlineQuery != nil {
				err = info.ParcoursupCommand(bot, &update)
				if err != nil {
					lib.LogError("Error handling an inline request: %s", err)
				}
			} else if update.CallbackQuery != nil {
				err = HandleCommand(bot, update, true)
				if err != nil {
					lib.LogError("Error handling a callback: %s", err)
				}
			} else if update.Message.IsCommand() {
				if update.Message.From.UserName != os.Getenv("TELEGRAM_USER") {
					continue
				}
				err = HandleCommand(bot, update, false)
				if err != nil {
					lib.LogError("Error handling a command: %s", err)
				}
			}
		}
	}()

	// Wait here until CTRL-C or other term signal is received.
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	lib.LogInfo("Gracefully shutting down bot 💤")
}
