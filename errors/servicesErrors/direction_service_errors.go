package servicesErrors

import "errors"

var (
	DirectionAlreadyExists = errors.New("Service error! Направление уже существует в базе!")
)
