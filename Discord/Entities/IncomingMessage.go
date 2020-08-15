package Entities

import "time"

type IncomingMessage struct {
	Id string
	Channel_id string
	Author User
	Content string
	Timestamp time.Time
	Edited_timestamp time.Time
	Type int
}