package main

import (
	"log"
	"main/Core"
)

func main() {
	bot := Core.Bot{}
	err := bot.Start()

	if err != nil {
		log.Fatal(err)
	}
}


