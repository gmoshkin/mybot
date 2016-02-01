package main

import (
    "github.com/tucnak/telebot"
    "github.com/fatih/color"
    "strconv"
    "time"
    "fmt"
    "os"
)

var (
    ConfigDir        = os.Getenv("HOME") + "/.telebot"
    OwnerIdFileName  = ConfigDir + "/ownerid"
    ApiTokenFileName = ConfigDir + "/apitoken"
)

func ReadConfig() (ownerId int, apiToken string) {
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
    ownerId, apiToken := ReadConfig()
    bot, err := telebot.NewBot(apiToken)
    if err != nil {
        return
    }
    messages := make(chan telebot.Message)
    bot.Listen(messages, 1 * time.Second)
    for message := range messages {
        if message.Sender.ID != ownerId {
            bot.SendMessage(message.Chat,
                            "Nothing to see here, move along!",
                            nil)
            continue
        }
        if message.Text == "/id" {
            bot.SendMessage(message.Chat,
                            "your id is " + strconv.Itoa(message.Sender.ID),
                            nil)
            continue
        }
        bot.SendMessage(message.Chat, message.Text, nil)
    }
}
