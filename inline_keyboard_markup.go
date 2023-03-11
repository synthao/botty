package telegram

import "encoding/json"

type InlineKeyboardButton struct {
	Text         string `json:"text"`
	URL          string `json:"url,omitempty"`
	CallbackData string `json:"callback_data,omitempty"`
	Unique       string `json:"unique,omitempty"`
}

type InlineKeyboardMarkup struct {
	InlineKeyboardMarkup [][]InlineKeyboardButton `json:"inline_keyboard"`
}

func WithRow(buttons ...InlineKeyboardButton) func(markup *InlineKeyboardMarkup) {
	return func(markup *InlineKeyboardMarkup) {
		row := make([]InlineKeyboardButton, len(buttons))

		for i, button := range buttons {
			if button.Unique != "" {
				button.CallbackData = button.Unique + "|" + button.CallbackData
				button.Unique = ""
			}
			row[i] = button
		}

		markup.InlineKeyboardMarkup = append(markup.InlineKeyboardMarkup, row)
	}
}

func NewInlineKeyboardMarkup(opts ...func(markup *InlineKeyboardMarkup)) *InlineKeyboardMarkup {
	m := new(InlineKeyboardMarkup)

	for _, opt := range opts {
		opt(m)
	}

	return m
}

func (m *InlineKeyboardMarkup) GetText() string {
	// TODO handle error
	data, _ := json.Marshal(m)

	return string(data)
}
