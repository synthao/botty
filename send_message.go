package botty

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

type MessageData struct {
	ChatID                   int
	MessageThreadID          int
	ReplyToMessageID         int
	Text                     string
	ParseMode                string
	Entities                 []MessageEntity
	DisableWebPagePreview    bool
	DisableNotification      bool
	ProtectContent           bool
	AllowSendingWithoutReply bool
	ReplyMarkup              ReplyMarkup
}

func (d *MessageData) validate() error {
	if d.ChatID == 0 {
		return fmt.Errorf("chat_id is required")
	}

	if strings.TrimSpace(d.Text) == "" {
		return fmt.Errorf("text is required")
	}

	return nil
}

func (c *Client) SendMessage(data *MessageData) (m *Message, err error) {
	defer func() { err = wrapIfErr("can't send message", err) }()

	if err := data.validate(); err != nil {
		return nil, err
	}

	v := url.Values{}
	text := strings.TrimSpace(data.Text)

	v.Add("chat_id", strconv.Itoa(data.ChatID))
	v.Add("text", text)
	v.Add("message_thread_id", strconv.Itoa(data.MessageThreadID))
	v.Add("parse_mode", data.ParseMode)
	v.Add("disable_web_page_preview", strconv.FormatBool(data.DisableWebPagePreview))
	v.Add("disable_notification", strconv.FormatBool(data.DisableNotification))
	v.Add("protect_content", strconv.FormatBool(data.ProtectContent))
	v.Add("reply_to_message_id", strconv.Itoa(data.ReplyToMessageID))
	v.Add("allow_sending_without_reply", strconv.FormatBool(data.AllowSendingWithoutReply))

	if err := addEntitiesToRequest(v, "entities", text, data.Entities); err != nil {
		return nil, err
	}

	if data.ReplyMarkup != nil {
		v.Add("reply_markup", data.ReplyMarkup.GetText())
	}

	res, err := c.doRequest(methodSendMessage, v)
	if err != nil {
		return nil, fmt.Errorf("can't send message, %w", err)
	}

	return c.processResponse(res)
}
