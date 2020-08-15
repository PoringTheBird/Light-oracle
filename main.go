package main

import (
	"log"
	"main/Discord"
	"os"
)

func main() {
	discordGateway := Discord.GatewayClient{GatewayUrl: os.Getenv("DISCORD_GATEWAY_URL")}
	//discordApi := Discord.ApiClient{os.Getenv("DISCORD_API_URL"), os.Getenv("DISCORD_TOKEN")}
	//
	//message, err := discordApi.SendMessage("Some message","743173612843958332")
	//
	//if err != nil {
	//	log.Fatal(err)
	//} else {
	//	text := fmt.Sprintf("%s: %s", message.Author.Username, message.Content)
	//	log.Println(text)
	//}

	err := discordGateway.Connect()

	if err != nil {
		log.Fatal(err)
	}
}


