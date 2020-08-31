package Core

import (
	"errors"
	"fmt"
	"log"
	"main/Core/Discord"
	"main/Core/Discord/Entities"
	"main/Core/Discord/Entities/Gateway"
	"main/Core/LightChat"
	"os"
	"strings"
)

const chatLoadInterval = 600.0

type Bot struct {
	discordApi Discord.ApiClient
	discordGateway Discord.GatewayClient
	lightChatContainer LightChat.HistoryContainer

	discordStateDetails *Gateway.ReadyEventPayload
}

// Control

func (bot *Bot) Start() error {
	bot.setupClients()
	bot.lightChatContainer.StartChatHistoryObserving(chatLoadInterval)

	testMsgs := []LightChat.Message{}
	testMsgs = append(testMsgs, LightChat.Message{SenderName: "yuryol", Text: "я встал,у меня пол первого,работать начал =З"})

	bot.OnNewMessagesLoaded(testMsgs)

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

	bot.lightChatContainer = LightChat.HistoryContainer{SiteUrl: os.Getenv("LIGHT_CHAT_URL"), LoadHandler: bot }
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

		return true
	default:
		return false
	}
}

// MessageLoadHandler

func (bot *Bot) OnNewMessagesLoaded(messages []LightChat.Message) {
	var newMessages = "```swift\n"

	for _, msg := range messages {
		newMessages += fmt.Sprintf("%s: %s", strings.Title(msg.SenderName), msg.Text)
	}

	newMessages += "\n```"

	bot.discordApi.SendMessage(newMessages, "743173612843958332")
}