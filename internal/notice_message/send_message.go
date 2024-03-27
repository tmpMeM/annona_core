package notice_message

func SendNoticeMessage(
	messageText string,
	botToken, chatID string, topicID *int64,
	disableWebPagePreview, disableNotification bool, disableButtons bool,
) error {
	// formattedText := fmt.Sprintf(
	// 	"%s"+"%s"+"\n\n"+`——<a href=%q>%s</a><a href=%q>%s</a>`, //%q	a single-quoted character literal safely escaped with Go syntax.
	// 	entryTitle,
	// 	tgHtmlText,
	// )
	formattedText := messageText

	message := &MessageRequest{
		ChatID:                chatID,
		Text:                  formattedText,
		ParseMode:             HTMLFormatting,
		DisableWebPagePreview: disableWebPagePreview,
		DisableNotification:   disableNotification,
	}

	if topicID != nil {
		if *topicID != 0 {
			message.MessageThreadID = *topicID
		}
	}

	// if !disableButtons {
	// 	var markupRow []*InlineKeyboardButton

	// 	articleURLButton3 := InlineKeyboardButton{Text: "", URL: "https://"}
	// 	markupRow = append(markupRow, &articleURLButton3)

	// 	message.ReplyMarkup = &InlineKeyboard{}
	// 	message.ReplyMarkup.InlineKeyboard = append(message.ReplyMarkup.InlineKeyboard, markupRow)
	// }

	client := NewClient(botToken, chatID)
	_, err := client.SendMessage(message)
	return err
}
