package Entities

type OutcomingMessage struct {
	Content string					`json:"content"`
	Tts 	bool					`json:"tts"`
	Embed	*OutcomingMessageEmbed	`json:"embed"`
}

type OutcomingMessageEmbed struct {
	Title string		`json:"title"`
	Description string	`json:"description"`
}