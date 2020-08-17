package Gateway

import "main/Core/Discord/Entities"

type ReadyEventPayload struct {
	GatewayVersion int           `json:"v"`
	UserDetails    Entities.User `json:"user"`
	SessionId      string        `json:"session_id"`
	Shard          *[]int        `json:"shard"`
}