package servicesErrors

import "errors"

var (
	ClientDoesNotExists = errors.New("Service error! Такого клиента не существует!")
	ClientAlreadyExists = errors.New("Service error! Клиент уже существует в базе!")
	ErrorGetClientByLogin = errors.New("Service error! Ошибка при получении клиента по логину!")
	InvalidPassword = errors.New("Service error! Неверный пароль!")
	AssignmentOnThisTimeAlreadyExists = errors.New("Service error! Клиент уже записан на тренировку в это время!")
	AgeNotCorrespondToAcceptableAge = errors.New("Service error! Возраст клиента не соответствует допустимому возрасту!")
	GenderNotCorrespondToAcceptableGender = errors.New("Service error! Пол клиента не соответствует допустимому полу!")
	ClientHasntGotSubscription = errors.New("Service error! Клиент не имеет абонемент!")
	ClientSubscriptionIsOver = errors.New("Service error! Абонемент клиента израсходован!")
	ClientTelephoneIncorrect = errors.New("Service error! Некорректный телефон клиента!")
	ClientMailIncorrect = errors.New("Service error! Некорректная почта клиента!")
)
