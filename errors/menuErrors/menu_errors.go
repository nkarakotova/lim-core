package menuErrors

import "errors"

var (
	ErrorCase = errors.New("Menu Error! Нет пункта меню с таким номером!")
	ErrorMenu = errors.New("Menu Error")
)
