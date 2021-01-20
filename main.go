package main

import (
	"fmt"
	"log"
	"os"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("💾 No .env file at the root - Ignoring")
	}

	bot, err := telegram.NewBotAPI(os.Getenv("TELEGRAM_TOKEN"))
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := telegram.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.InlineQuery == nil {
			continue
		}

		records := SearchRecords(update.InlineQuery.Query).Records

		var results []interface{}
		for _, record := range records {
			messageContent := fmt.Sprintf("[%s](%s)\n%s - %s", record.Fields.Name, record.Fields.Link, record.Fields.Course, record.Fields.Specialization)
			if record.Fields.CourseDetail != "" {
				messageContent += " - " + record.Fields.CourseDetail
			}
			messageContent += fmt.Sprintf("\n_%s, %s_", record.Fields.Department, record.Fields.Region)

			results = append(results, telegram.InlineQueryResultArticle{
				Type:  "article",
				ID:    record.ID,
				Title: fmt.Sprintf("%s - %s - %s", record.Fields.Name, record.Fields.Course, record.Fields.Specialization),
				URL:   record.Fields.Link,
				InputMessageContent: telegram.InputTextMessageContent{
					Text:                  messageContent,
					ParseMode:             "Markdown",
					DisableWebPagePreview: false,
				},
			})
		}

		_, err := bot.AnswerInlineQuery(telegram.InlineConfig{
			InlineQueryID: update.InlineQuery.ID,
			Results:       results,
		})
		if err != nil {
			log.Println(err)
		}
	}
}
