package botty

import (
	"encoding/json"
	"fmt"
	"net/url"
)

type AnswerCallbackQueryData struct {
	CallbackQueryID string
	Text            string
}

type ReplyToQueryResponse struct {
	OK          bool   `json:"ok"`
	Result      bool   `json:"result"`
	Description string `json:"description"`
}

func (c *Client) replyToQuery(u Update) error {
	v := url.Values{}
	m := &AnswerCallbackQueryData{
		CallbackQueryID: u.CallbackQuery.ID,
	}

	v.Add("callback_query_id", m.CallbackQueryID)

	res, err := c.doRequest(methodAnswerCallbackQuery, v)
	if err != nil {
		return fmt.Errorf("can't send query callback response, %w", err)
	}

	decodedRes := new(ReplyToQueryResponse)
	if err := json.Unmarshal(res, decodedRes); err != nil {
		return fmt.Errorf("can't decode response after replying to query")
	}
	if !decodedRes.OK {
		if decodedRes.Description != "" {
			return fmt.Errorf("can't decode response after replying to query, %s", decodedRes.Description)
		}
		return fmt.Errorf("can't decode response after replying to query")
	}

	return nil
}
