package main

import "gitlab.com/thefrol/notty/internal/entity"

// Клиенты которые будут добавлены
var custs = []entity.Customer{
	{
		Id:       "vasya-top",
		Phone:    "+79162325566",
		Tag:      "cool",
		Operator: "МТС",
		Name:     "Василий Иванович",
	},
	{
		Id:       "liza-top",
		Phone:    "+71162325566",
		Tag:      "love",
		Operator: "МТС",
		Name:     "Лизавета Н.",
	},
	{
		Id:       "layla666",
		Phone:    "+71262325566",
		Tag:      "hate",
		Operator: "yota",
		Name:     "Попова М.",
	},
	{
		Id:       "random",
		Phone:    "+71262325566",
		Tag:      "someone",
		Operator: "beeline",
		Name:     "Иванов И.И.",
	},
}
