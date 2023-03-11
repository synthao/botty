package telegram

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

type UpdateMessageData struct {
	ChatID                int
	MessageID             int
	InlineMessageID       string
	Text                  string
	ParseMode             string
	Entities              []MessageEntity
	DisableWebPagePreview bool
	ReplyMarkup           ReplyMarkup
}

func (t *UpdateMessageData) validate() error {
	if t.InlineMessageID == "" && t.ChatID == 0 {
		return fmt.Errorf("chat_id or inline_message_id is required")
	}

	if t.MessageID == 0 && t.InlineMessageID == "" {
		return fmt.Errorf("message_id or inline_message_id is required")
	}

	if t.InlineMessageID == "" && t.ChatID == 0 && t.MessageID == 0 {
		return fmt.Errorf("message_id or inline_message_id is required")
	}

	if strings.TrimSpace(t.Text) == "" {
		return fmt.Errorf("text is required")
	}

	return nil
}

func (c *Client) UpdateMessage(data *UpdateMessageData) (m *Message, err error) {
	defer func() { err = wrapIfErr("can't update message", err) }()

	if err := data.validate(); err != nil {
		return nil, err
	}

	v := url.Values{}
	text := strings.TrimSpace(data.Text)

	v.Add("chat_id", strconv.Itoa(data.ChatID))
	v.Add("message_id", strconv.Itoa(data.MessageID))
	v.Add("text", text)
	v.Add("inline_message_id", data.InlineMessageID)
	v.Add("parse_mode", data.ParseMode)
	v.Add("disable_web_page_preview", strconv.FormatBool(data.DisableWebPagePreview))

	if err := addEntitiesToRequest(v, "entities", text, data.Entities); err != nil {
		return nil, err
	}

	if data.ReplyMarkup != nil {
		v.Add("reply_markup", data.ReplyMarkup.GetText())
	}

	res, err := c.doRequest(methodEditMessageText, v)
	if err != nil {
		return nil, err
	}

	return c.processResponse(res)
}
