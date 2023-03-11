# Botty

Botty is Telegram Bot API library for Golang

## Getting Started

Get the library:

```shell
$ go get -u github.com/synthao/botty
```

## Examples

```go
package main

import "github.com/synthao/botty"

func main() {
    errHandler := func(err error) {
        println(err.Error())
    }
    
    client := botty.NewClient("your-bot-token", botty.WithErrorHandler(errHandler))
    
    // handle single command
    client.OnCommand("/start", func(u botty.Update) error {
        return client.Reply(u, "Welcome!")
    })
    
    // handle few commands
    client.OnCommands([]string{"/start", "/help"}, func(u botty.Update) error {
        return client.Reply(u, "Hello!")
    })
	
    // handle message
    client.OnMessage("hey", func(u botty.Update) error {
        return client.Reply(u, "How are you ?")
    })

    // reply to any message
    client.OnMessage("*", func(u botty.Update) error {
        return client.Reply(u, "Some message")
    })

    // send message using the SendMessage() method and add inline keyboard
    _, err := client.SendMessage(&botty.MessageData{
        ParseMode: botty.ParseModeHTML, // add formatting mode
        ChatID:    "-123456789",
        Text:      "Do you like botty ?",
        ReplyMarkup: botty.NewInlineKeyboardMarkup(
            botty.WithRow(
                botty.InlineKeyboardButton{
                    Text: "🚫 No",      // button label
                    CallbackData: "no", // custom data
                    Unique: "btn_no",   // a unique button key to handle it using the OnQuery("btn_no") method
                },
                botty.InlineKeyboardButton{
                    Text: "👍 Yes", 
                    CallbackData: "yes", 
                    Unique: "btn_yes",
                },
            ),
        ),
    })
    if err != nil {
        println(err.Error())
    }

    //Parse mode constants: ParseModeMarkdownV2|ParseModeMarkdown|ParseModeHTML
	
    // handle any query
    client.OnQuery("*", func(u botty.Update) error {
        kb := botty.NewInlineKeyboardMarkup(
            botty.WithRow(
                botty.InlineKeyboardButton{
                    Text: "Google",
                    URL:  "https://google.com",
                },
            ),
            botty.WithRow(
                botty.InlineKeyboardButton{
                    Text: "Apple",
                    URL:  "https://apple.com",
                },
            ),
        )

        _, err := botty.UpdateMessage(&botty.UpdateMessageData{
            ChatID: u.CallbackQuery.Message.Chat.ID,
            MessageID: u.CallbackQuery.Message.MessageID,
            Text: "Some message",
            ReplyMarkup: kb,
        })
    
        return err
    })

    // handle unique query
    client.OnQuery("btn_no", func(u botty.Update) error {
        return client.Reply(u, "Why ?")
    })

    client.OnQuery("btn_yes", func(u botty.Update) error {
        return client.Reply(u, "Great :)")
    })

    // reply using the SendMessage() method 
    client.OnMessage("question", func(u botty.Update) error {
        _, err := client.SendMessage(&botty.MessageData{
            ParseMode: botty.ParseModeHTML, // add formatting mode
            ChatID:    -123456789,
            Text:      "Do you like botty ?",
            ReplyMarkup: botty.NewInlineKeyboardMarkup(
                botty.WithRow(
                    botty.InlineKeyboardButton{
                        Text:         "🚫 No",  // button label
                        CallbackData: "no",     // custom data
                        Unique:       "btn_no", // a unique button key to handle it using the OnQuery("btn_no") method
                    },
                    botty.InlineKeyboardButton{
                        Text:         "👍 Yes",
                        CallbackData: "yes",
                        Unique:       "btn_yes",
                    },
                ),
            ),
        })
        return err
    })

    // update message
    client.OnMessage("message to replace", func(u botty.Update) error {
        _, err := client.UpdateMessage(&botty.UpdateMessageData{
            ChatID:    u.CallbackQuery.Message.Chat.ID,
            MessageID: u.CallbackQuery.Message.MessageID,
            Text:      "Replaced text",
        })
        return err
    })

    // add formatting
    entities := []botty.MessageEntity{
        {Type: "italic"},
    }

    _, err := client.SendMessage(&botty.MessageData{
        ChatID: -123456789,
        Text: "Hello",
        Entities: entities,
    })
    if err != nil {
        println(err.Error())
    }

}


```
