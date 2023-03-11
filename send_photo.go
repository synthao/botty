package telegram

import (
	"fmt"
	"net/url"
	"os"
	"strconv"
)

type SendPhotoData struct {
	ChatID                   int
	MessageThreadID          int
	ReplyToMessageID         int
	Photo                    string
	Caption                  string
	ParseMode                string
	DisableNotification      bool
	ProtectContent           bool
	AllowSendingWithoutReply bool
	CaptionEntities          []MessageEntity
	ReplyMarkup              ReplyMarkup
}

func (d *SendPhotoData) validate() (err error) {
	if d.ChatID == 0 {
		return fmt.Errorf("chat_id is required")
	}

	if d.Photo == "" {
		return fmt.Errorf("photo is required")
	}

	return nil
}

func (d *SendPhotoData) isLocalFilePath() bool {
	uri, _ := url.ParseRequestURI(d.Photo)
	if uri.Host == "" {
		return true
	}

	return false
}

func (c *Client) SendPhoto(d *SendPhotoData) (_ *Message, err error) {
	defer func() { err = wrapIfErr("can't send photo", err) }()

	if err := d.validate(); err != nil {
		return nil, err
	}

	if d.isLocalFilePath() {
		return c.sendLocalFile(d)
	}

	return c.sendRemoteFile(d)
}

func (c *Client) sendLocalFile(d *SendPhotoData) (*Message, error) {
	f, err := os.Open(d.Photo)
	if err != nil {
		return nil, err
	}

	form := NewMultipartForm()

	if err = form.AddField("chat_id", strconv.Itoa(d.ChatID)); err != nil {
		return nil, err
	}

	if err = form.AddField("message_thread_id", strconv.Itoa(d.MessageThreadID)); err != nil {
		return nil, err
	}

	if err := form.AddField("caption", d.Caption); err != nil {
		return nil, err
	}

	if err := form.AddField("parse_mode", d.ParseMode); err != nil {
		return nil, err
	}

	if err := addEntitiesToForm(form, "caption_entities", d.Caption, d.CaptionEntities); err != nil {
		return nil, err
	}

	if err = form.AddField("disable_notification", strconv.FormatBool(d.DisableNotification)); err != nil {
		return nil, err
	}

	if err = form.AddField("protect_content", strconv.FormatBool(d.ProtectContent)); err != nil {
		return nil, err
	}

	if err = form.AddField("reply_to_message_id", strconv.Itoa(d.ReplyToMessageID)); err != nil {
		return nil, err
	}

	if err = form.AddField("allow_sending_without_reply", strconv.FormatBool(d.AllowSendingWithoutReply)); err != nil {
		return nil, err
	}

	if d.ReplyMarkup != nil {
		if err = form.AddField("reply_markup", d.ReplyMarkup.GetText()); err != nil {
			return nil, err
		}
	}

	if err = form.AddFile("photo", f); err != nil {
		return nil, err
	}

	res, err := c.doMultipartFormRequest(methodSendPhoto, form)
	if err != nil {
		return nil, err
	}

	return c.processResponse(res)
}

func (c *Client) sendRemoteFile(d *SendPhotoData) (*Message, error) {
	v := url.Values{}

	v.Add("chat_id", strconv.Itoa(d.ChatID))
	v.Add("photo", d.Photo)
	v.Add("message_thread_id", strconv.Itoa(d.MessageThreadID))
	v.Add("caption", d.Caption)
	v.Add("parse_mode", d.ParseMode)
	v.Add("disable_notification", strconv.FormatBool(d.DisableNotification))
	v.Add("protect_content", strconv.FormatBool(d.ProtectContent))
	v.Add("reply_to_message_id", strconv.Itoa(d.ReplyToMessageID))
	v.Add("allow_sending_without_reply", strconv.FormatBool(d.AllowSendingWithoutReply))

	if d.ReplyMarkup != nil {
		v.Add("reply_markup", d.ReplyMarkup.GetText())
	}

	if err := addEntitiesToRequest(v, "caption_entities", d.Caption, d.CaptionEntities); err != nil {
		return nil, err
	}

	res, err := c.doRequest(methodSendPhoto, v)
	if err != nil {
		return nil, err
	}

	return c.processResponse(res)
}
