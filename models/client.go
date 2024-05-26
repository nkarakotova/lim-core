package models

type Client struct {
	ID             uint64
	SubscriptionID uint64
	Name           string
	Telephone      string
	Mail           string
	Password       string
	Age            uint16
	Gender         Gender
}
