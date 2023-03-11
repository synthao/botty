package botty

type ReplyOption func(s *MessageData)

func WithParseMode(mode string) ReplyOption {
	return func(s *MessageData) {
		s.ParseMode = mode
	}
}

func WithMessageThreadID(messageThreadID int) ReplyOption {
	return func(s *MessageData) {
		s.MessageThreadID = messageThreadID
	}
}

func WithEntities(entities []MessageEntity) ReplyOption {
	return func(s *MessageData) {
		s.Entities = entities
	}
}

func WithDisableWebPagePreview(disableWebPagePreview bool) ReplyOption {
	return func(s *MessageData) {
		s.DisableWebPagePreview = disableWebPagePreview
	}
}

func WithDisableNotification(disableNotification bool) ReplyOption {
	return func(s *MessageData) {
		s.DisableNotification = disableNotification
	}
}

func WithProtectContent(protectContent bool) ReplyOption {
	return func(s *MessageData) {
		s.ProtectContent = protectContent
	}
}

func WithReplyToMessageID(replyToMessageID int) ReplyOption {
	return func(s *MessageData) {
		s.ReplyToMessageID = replyToMessageID
	}
}

func WithAllowSendingWithoutReply(allowSendingWithoutReply bool) ReplyOption {
	return func(s *MessageData) {
		s.AllowSendingWithoutReply = allowSendingWithoutReply
	}
}

func WithReplyMarkup(markup ReplyMarkup) ReplyOption {
	return func(s *MessageData) {
		s.ReplyMarkup = markup
	}
}

func (c *Client) Reply(u Update, text string, options ...ReplyOption) error {
	m := &MessageData{
		ChatID: u.Message.Chat.ID,
		Text:   text,
	}

	for _, o := range options {
		o(m)
	}

	_, err := c.SendMessage(m)
	return err
}
