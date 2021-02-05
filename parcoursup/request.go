package parcoursup

import (
	"fmt"
	"log"

	telegram "github.com/go-telegram-bot-api/telegram-bot-api"
)

func HandleRequest(bot *telegram.BotAPI, update *telegram.Update) {
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
