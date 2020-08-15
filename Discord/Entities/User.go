package Entities

type User struct {
	Id string
	Username string
	Discriminator string
	Avatar *string
	Bot bool
	System bool
	Locale string
	Email *string
}