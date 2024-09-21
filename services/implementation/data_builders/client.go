package servicesDataBuilder

import (
	"github.com/nkarakotova/lim-core/models"
)

type ClientBuilder struct {
	client *models.Client
}

func NewClientBuilder() *ClientBuilder {
	return &ClientBuilder{
		client: &models.Client{
			ID: 1,
			Name: "Name",
			Telephone: "1234567890",
			Mail: "mail@mail.ru",
			Password: "123",
		},
	}
}

func (b *ClientBuilder) WithID(id uint64) *ClientBuilder {
	b.client.ID = id
	return b
}

func (b *ClientBuilder) WithName(name string) *ClientBuilder {
	b.client.Name = name
	return b
}

func (b *ClientBuilder) WithTelephone(telephone string) *ClientBuilder {
	b.client.Telephone = telephone
	return b
}

func (b *ClientBuilder) WithMail(mail string) *ClientBuilder {
	b.client.Mail = mail
	return b
}

func (b *ClientBuilder) WithPassword(password string) *ClientBuilder {
	b.client.Password = password
	return b
}

func (b *ClientBuilder) Build() *models.Client {
	return b.client
}
