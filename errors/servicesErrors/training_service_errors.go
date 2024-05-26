package servicesErrors

import "errors"

var (
	NoAvailablePlacesNum = errors.New("Service error! Не осталось свободных мест на тренировке!")
	PlacesNumMoreThenCapacity = errors.New("Service error! На тренировке больше мест, чем вместительность зала!")
	IncorrectTrainingTime = errors.New("Service error! Тренировка начиниется в недопустимое время!")
	BysyDateTime = errors.New("Service error! В данное время занят зал или тренер!")
	StartDateAfterEndDate = errors.New("Service error! Начальная дата позже конечной!")
	NotCoachByDirection = errors.New("Service error! Тренер не ведёт данное направление!")
)
