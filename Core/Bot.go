package Core

import (
	"errors"
	"log"
	"main/Core/Discord"
	"main/Core/Discord/Entities"
	"main/Core/Discord/Entities/Gateway"
	"os"
)

type Bot struct {
	discordApi Discord.ApiClient
	discordGateway Discord.GatewayClient

	discordStateDetails *Gateway.ReadyEventPayload
}

// Control

func (bot *Bot) Start() error {
	bot.setupClients()

	return bot.discordGateway.Connect()
}

// Setup

func (bot *Bot) setupClients() {
	bot.discordApi = Discord.ApiClient{
		BaseUrl: os.Getenv("DISCORD_API_URL"),
		DiscordToken: os.Getenv("DISCORD_TOKEN"),
	}

	bot.discordGateway = Discord.GatewayClient{
		GatewayUrl: os.Getenv("DISCORD_GATEWAY_URL"),
		DiscordToken: os.Getenv("DISCORD_TOKEN"),
		ContentMessageHandler: bot,
	}
}

// Actions

func (bot *Bot) sendMessageToDiscord(message string, channelId string) (*Entities.IncomingMessage, error) {
	if !bot.discordGateway.IsIdentified { return nil, errors.New("Not identified by Gateway") }
	return bot.discordApi.SendMessage(message, channelId)
}

// Message handling

func (bot *Bot) HandleMessage(message *Gateway.Message) bool {
	switch *message.T {
	case "READY":
		payload := new(Gateway.ReadyEventPayload)
		err := remapJson(message.D, payload)

		if err != nil {
			log.Fatal("Invalid 'READY' event payload. Error: ", err)
			return false
		}

		bot.discordStateDetails = payload
		bot.discordGateway.IsIdentified = true

		msg, err := bot.sendMessageToDiscord("Hello!", "743736138258448474")

		if err != nil {
			log.Fatal("Error: ", err)
		} else {
			log.Println("Message: ", msg)
		}

		return true
	default:
		return false
	}
}