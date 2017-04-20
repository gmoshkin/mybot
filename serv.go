package main

import (
    "github.com/tucnak/telebot"
    "github.com/fatih/color"
    "strconv"
    "time"
    "fmt"
    "log"
    "os"
)

var (
    ConfigDir        = os.Getenv("HOME") + "/.telebot"
    OwnerIdFileName  = ConfigDir + "/ownerid"
    ApiTokenFileName = ConfigDir + "/apitoken"
    bot *telebot.Bot
    ownerId int
    apiToken string

)

func ReadConfig() (ownerId int, apiToken string) {
    log.Printf("Reading configs from %s... ", ConfigDir)
    ownerIdFile, err := os.Open(OwnerIdFileName)
    if err != nil {
        panic(fmt.Sprintf("Couldn't open file %s", OwnerIdFileName))
    }
    defer ownerIdFile.Close()
    n, err := fmt.Fscanf(ownerIdFile, "%d\n", &ownerId)
    if n < 1 || err != nil {
        panic(fmt.Sprintf("Failed to read file %s:\n%s", OwnerIdFileName, err))
    }
    apiTokenFile, err := os.Open(ApiTokenFileName)
    if err != nil {
        panic(fmt.Sprintf("Couldn't open file %s", ApiTokenFileName))
    }
    defer apiTokenFile.Close()
    n, err = fmt.Fscanln(apiTokenFile, &apiToken)
    if n < 1 || err != nil {
        panic(fmt.Sprintf("Failed to read file %s:\n%s", ApiTokenFileName, err))
    }
    log.Printf("done\n")
    return
}

func main() {
    defer func() {
        msg := recover()
        if msg != nil {
            color.Set(color.FgRed)
            fmt.Fprintln(os.Stderr, msg)
        }
    } ()

    ownerId, apiToken = ReadConfig()

    var err error
    log.Printf("Setting up connection... ")
    bot, err = telebot.NewBot(apiToken)
    if err != nil {
        return
    }
    log.Printf("done\n")

    bot.Messages = make(chan telebot.Message)
    bot.Queries = make(chan telebot.Query)

    go processMessages()
    go processQueries()

    log.Println("Starting bot")
    bot.Start(1 * time.Second)
}

func processMessages() {
    for message := range bot.Messages {
        if message.Sender.ID != ownerId {
            log.Printf("Recieved a message from %s with text:\n%s\n",
                       message.Sender.Username, message.Text)
            bot.SendMessage(message.Chat,
                            "Nothing to see here, move along!",
                            nil)
            continue
        }
        log.Printf("Recieved a message from owner with text:\n\"%s\"\n",
                   message.Text)
        if message.Text == "/id" {
            bot.SendMessage(message.Chat,
                            "your id is " + strconv.Itoa(message.Sender.ID),
                            nil)
            continue
        }
        bot.SendMessage(message.Chat, message.Text, nil)
    }
}

func processQueries() {
    for query := range bot.Queries {
        log.Println("--- new query ---")
        log.Println("from:", query.From.Username)
        log.Println("text:", query.Text)

        // Create an article (a link) object to show in our results.
        article := &telebot.InlineQueryResultArticle{
            Title: "Telegram bot framework written in Go",
            URL:   "https://github.com/tucnak/telebot",
            InputMessageContent: &telebot.InputTextMessageContent{
                Text:           "Telebot is a convenient wrapper to Telegram Bots API, written in Golang.",
                DisablePreview: false,
            },
        }

        article2 := &telebot.InlineQueryResultArticle{
            Title: getTime(),
            URL: "https://time.is",
            InputMessageContent: &telebot.InputTextMessageContent{
                Text: getTime(),
                DisablePreview: false,
            },
        }

        article3 := &telebot.InlineQueryResultArticle{
            Title: query.Text,
            InputMessageContent: &telebot.InputTextMessageContent{
                Text: query.Text,
                DisablePreview: false,
            },
        }

        // Build the list of results. In this instance, just our 1 article from above.
        results := []telebot.InlineQueryResult{article, article2, article3}

        // Build a response object to answer the query.
        response := telebot.QueryResponse{
            Results:    results,
            IsPersonal: true,
        }

        // And finally send the response.
        if err := bot.AnswerInlineQuery(&query, &response); err != nil {
            log.Println("Failed to respond to query:", err)
        }
    }
}

func getTime() (t string) {
    t = time.Now().Format("Mon Jan 2 2006 15:04:05.000")
    return
}
