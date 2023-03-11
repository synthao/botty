package telegram

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"strings"
)

const (
	baseUrl        = "api.telegram.org"
	basePathPrefix = "bot"
)

const (
	methodGetUpdates          = "getUpdates"
	methodSendMessage         = "sendMessage"
	methodEditMessageText     = "editMessageText"
	methodAnswerCallbackQuery = "answerCallbackQuery"
	methodSendPhoto           = "sendPhoto"
)

const (
	ParseModeMarkdownV2 = "MarkdownV2"
	ParseModeMarkdown   = "Markdown"
	ParseModeHTML       = "HTML"
)

type ReplyMarkup interface {
	GetText() string
}

type Client struct {
	client       http.Client
	commands     map[string]func(u Update) error
	messages     map[string]func(u Update) error
	queries      map[string]func(u Update) error
	token        string
	host         string
	basePath     string
	offset       int
	errorHandler func(error)
}

type ClientOption func(*Client)

func WithErrorHandler(handler func(err error)) ClientOption {
	return func(c *Client) {
		c.errorHandler = handler
	}
}

func NewClient(token string, options ...ClientOption) *Client {
	c := &Client{
		client:   http.Client{},
		commands: make(map[string]func(Update) error),
		messages: make(map[string]func(Update) error),
		queries:  make(map[string]func(Update) error),
		token:    token,
		host:     baseUrl,
		basePath: basePathPrefix + token,
	}

	for _, o := range options {
		o(c)
	}

	return c
}

func (c *Client) Run() error {
	for {
		updates, err := c.getUpdates()
		if err != nil {
			return err
		}
		if len(updates) == 0 {
			continue
		}

		for _, u := range updates {
			if err := c.processUpdate(u); err != nil {
				if c.errorHandler != nil {
					c.errorHandler(err)
					continue
				}

				return err
			}
		}
	}
}

func (c *Client) OnCommands(commands []string, f func(Update) error) {
	for _, cmd := range commands {
		c.commands[cmd] = f
	}
}

func (c *Client) OnCommand(cmd string, f func(Update) error) {
	c.commands[cmd] = f
}

func (c *Client) OnMessages(messages []string, f func(Update) error) {
	for _, msg := range messages {
		c.messages[msg] = f
	}
}

func (c *Client) OnMessage(msg string, f func(Update) error) {
	c.messages[msg] = f
}

func (c *Client) OnQuery(query string, f func(Update) error) {
	c.queries[query] = f
}

func (c *Client) processCommand(u Update) (bool, error) {
	if !u.hasMessageText() {
		return false, nil
	}

	cmd := c.parseCommand(u.Message.Text)
	if cmd == "" {
		return false, nil
	}

	if fn, ok := c.commands[cmd]; ok {
		if err := fn(u); err != nil {
			return false, fmt.Errorf("can't process command, %w", err)
		}
	}

	return true, nil
}

func (c *Client) processMessage(u Update) (bool, error) {
	if !u.hasMessageText() {
		return false, nil
	}

	if fn, ok := c.messages["*"]; ok {
		if err := fn(u); err != nil {
			return false, err
		}
	}

	msg := c.parseMessage(u.Message.Text)

	if fn, ok := c.messages[msg]; ok {
		if err := fn(u); err != nil {
			return false, err
		}
	}

	return true, nil
}

func (c *Client) processQuery(u Update) (bool, error) {
	if fn, ok := c.queries["*"]; ok {
		if err := fn(u); err != nil {
			return false, err
		}
	}

	if u.CallbackQuery.Data != "" {
		items := strings.Split(u.CallbackQuery.Data, "|")

		if len(items) > 0 {
			u.CallbackQuery.Data = strings.Trim(u.CallbackQuery.Data, items[0]+"|")
		}

		if fn, ok := c.queries[items[0]]; ok {
			if err := fn(u); err != nil {
				return false, err
			}
		}
	}

	return true, nil
}

func (c *Client) processUpdate(u Update) error {
	processed, err := c.processCommand(u)
	if err != nil {
		return err
	}
	if processed {
		return nil
	}

	processed, err = c.processMessage(u)
	if err != nil {
		return err
	}
	if processed {
		return nil
	}

	processed, err = c.processQuery(u)
	if err != nil {
		return err
	}
	if processed {
		return c.replyToQuery(u)
	}

	return nil
}

func (c *Client) getUpdates() ([]Update, error) {
	updates, err := c.updates(c.offset, 100)
	if err != nil {
		return nil, fmt.Errorf("an error occurred while receiving updates, %w", err)
	}
	if !updates.OK {
		return nil, fmt.Errorf("can't get updates, %s", updates.Description)
	}

	if len(updates.Result) == 0 {
		return make([]Update, 0), nil
	}

	c.offset = updates.Result[len(updates.Result)-1].UpdateID + 1

	return updates.Result, nil
}

func (c *Client) parseCommand(cmd string) string {
	cmd = strings.TrimSpace(cmd)

	cmd = strings.ToLower(cmd)
	if cmd == "" {
		return ""
	}

	if cmd[0:1] != "/" {
		return ""
	}

	return cmd
}

func (c *Client) parseMessage(msg string) string {
	return strings.TrimSpace(msg)
}

func (c *Client) updates(offset, limit int) (UpdateResponse, error) {
	q := url.Values{}

	q.Add("offset", strconv.Itoa(offset))
	q.Add("limit", strconv.Itoa(limit))

	data, err := c.doRequest(methodGetUpdates, q)
	if err != nil {
		return UpdateResponse{}, fmt.Errorf("can't get updates, %w", err)
	}

	var res UpdateResponse

	if err := json.Unmarshal(data, &res); err != nil {
		return UpdateResponse{}, fmt.Errorf("can't get updates, %w", err)
	}

	return res, nil
}

func (c *Client) doRequest(method string, query url.Values) ([]byte, error) {
	u := url.URL{
		Scheme: "https",
		Host:   c.host,
		Path:   path.Join(c.basePath, method),
	}

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("doRequest() - can't create request, %w", err)
	}

	req.URL.RawQuery = query.Encode()

	res, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("doRequest() - can't exec request, %w", err)
	}
	defer func() { _ = res.Body.Close() }()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("doRequest() - can't read response, %w", err)
	}

	return body, nil
}

func (c *Client) doMultipartFormRequest(method string, form MultipartForm) ([]byte, error) {
	u := url.URL{
		Scheme: "https",
		Host:   c.host,
		Path:   path.Join(c.basePath, method),
	}

	req, err := http.NewRequest(http.MethodPost, u.String(), form.Form())
	if err != nil {
		return nil, fmt.Errorf("doMultipartFormRequest() - can't create request, %w", err)
	}

	req.Header.Add("Content-Type", form.FormDataContentType())

	res, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("doMultipartFormRequest() - can't exec request, %w", err)
	}
	defer func() { _ = res.Body.Close() }()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("doMultipartFormRequest() - can't read response, %w", err)
	}

	return body, nil
}

func (c *Client) processResponse(res []byte) (*Message, error) {
	decodedRes := new(Response)

	if err := json.Unmarshal(res, decodedRes); err != nil {
		return nil, fmt.Errorf("can't decode response, %w", err)
	}

	if !decodedRes.OK {
		return nil, fmt.Errorf("code: %d, description: %s", decodedRes.ErrorCode, decodedRes.Description)
	}

	return decodedRes.Result, nil
}
