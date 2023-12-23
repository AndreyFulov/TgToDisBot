package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	tokenDS  string
	tokenTG  string
	tgChanId string
}

func NewBot(ds, tg, chanId string) *Bot {
	return &Bot{tokenDS: ds, tokenTG: tg, tgChanId: chanId}
}

func (bot *Bot) Bot() {
	discord, err := discordgo.New("Bot " + bot.tokenDS)
	if err != nil {
		fmt.Errorf("ERROR! %d", err)
	}
	discord.AddHandler(bot.handleDsBot)
	discord.Identify.Intents = discordgo.IntentGuildMessages
	err = discord.Open()
	discord.Close()
	if err != nil {
		log.Fatalf("Error opening Discord session: %s", err)
		return
	}
	telegramBot, err := tgbotapi.NewBotAPI(bot.tokenTG)
	if err != nil {
		log.Panic(err)
		return
	}
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60
	updates := telegramBot.GetUpdatesChan(updateConfig)
	for update := range updates {
		if update.Message != nil {
			// Проверка наличия текста в сообщении
			if update.Message.Text != "" {
				// Получение текста сообщения из Telegram
				message := update.Message.Text

				// Отправка сообщения из Telegram в Discord
				_, err := discord.ChannelMessageSend(bot.tgChanId, message)
				if err != nil {
					log.Printf("Error sending message to Discord: %s", err)
				}
			}
		}
	}
	// Обработка завершения работы по сигналу Ctrl+C
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Закрытие сессии Discord
	err = discord.Close()
	if err != nil {
		log.Printf("Error closing Discord session: %s", err)
	}
}

func (bot *Bot) handleDsBot(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	if m.Content == "ping" {
		s.ChannelMessageSend(m.ChannelID, "Pong!")
	}
}
