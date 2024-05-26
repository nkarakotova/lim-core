package servicesErrors

import "errors"

var (
	SubscriptionStartDateAfterEndDate = errors.New("Service error! Дата начала действия абонимента позже, чем дата конца действия!")
)
