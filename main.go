package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"gopkg.in/telegram-bot-api.v4"
)

// Constants
const (
	Poloniex = 1
	Bitfinex = 2
)

func main() {
	telegramAPIKey := os.Getenv("CoinTronTelegramAPIKey")
	log.Print(telegramAPIKey)
	bot, err := tgbotapi.NewBotAPI(telegramAPIKey)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = false
	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Printf("Unable to get bot updates %s", err)
	}

	for update := range updates {
		if update.Message == nil {
			continue
		}
		messageStr, msgErr := messageHandler(update.Message.Text, bot.Self.UserName)
		if msgErr != nil {
			log.Printf("Message error: %s", msgErr)
		} else {
			log.Printf("Responding to query \"%s\" made by %s", update.Message.Text, update.Message.From)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, messageStr)
			_, msgSendErr := bot.Send(msg)
			if msgSendErr != nil {
				log.Printf("Bot message send error %s", msgSendErr)
			}
		}
	}
}

func enabledBotCommand(s string, botname string) (int, error) {
	str := strings.Split(s, " ")
	if str[0] == fmt.Sprintf("/polo@%s", botname) || str[0] == "/polo" {
		return Poloniex, nil
	} else if str[0] == fmt.Sprintf("/bitfinex@%s", botname) || str[0] == "/bitfinex" {
		return Bitfinex, nil
	}
	return 0, errors.New("Invalid command")
}

func messageHandler(s string, botname string) (string, error) {
	str := strings.Split(s, " ")
	requestExchange, commandError := enabledBotCommand(s, botname)
	if commandError != nil {
		log.Printf("Bot command error %s", commandError)
	}

	switch requestExchange {
	case Poloniex:
		if len(str) == 1 {
			data, err := getCurrentPoloniex("BTC LTC")
			return data, err
		} else if len(str) == 3 {
			data, err := getCurrentPoloniex(fmt.Sprintf("%s %s", strings.ToUpper(str[1]), strings.ToUpper(str[2])))
			return data, err
		}
	case Bitfinex:
		if len(str) == 1 {
			data, err := getCurrentBitfinex("LTC BTC")
			return data, err
		} else if len(str) == 3 {
			data, err := getCurrentBitfinex(fmt.Sprintf("%s %s", strings.ToUpper(str[1]), strings.ToUpper(str[2])))
			return data, err
		}
	}
	return "", errors.New("Invalid command")
}
