package Gateway

type IdentityPayload struct {
	Token string 							`json:"token"`
	Properties IdentityPayloadProperties	`json:"properties"`
}

type IdentityPayloadProperties struct {
	Os string								`json:"$os"`
	Browser string							`json:$browser`
	Device string							`json:$device`
}