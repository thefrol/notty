package entity_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/thefrol/notty/internal/entity"
	"gitlab.com/thefrol/notty/internal/entity/valid"
)

// TestCustomerRequest проверяет, что у такого человека
// есть конкретная ошибка валидации, или вовсе их нет
func TestCustomerRequest(t *testing.T) {

	tests := []struct {
		name string
		c    entity.Customer
		has  error
	}{
		//
		// проверки номера телефона
		//
		{
			name: "Валидный",
			c:    VasyaValid(),
			has:  nil,
		},
		{
			name: "короткий номер",
			c:    VasyaShortNumber(),
			has:  valid.ErrorPhoneValidation,
		},
		{
			name: "номер без + в начале",
			c:    NotPlusNumber(),
			has:  valid.ErrorPhoneValidation,
		},

		// //
		// // проверки Id
		// //

		{
			name: "UUID подходит",
			c:    VasyaUUID(),
			has:  nil,
		},
		{
			name: "кириллица не подходит",
			c:    VasyaCyrillic(),
			has:  valid.ErrorIdValidation,
		},
		{
			name: "нижнее подчеркивание не подходит",
			c:    VasyaUnderscore(),
			has:  valid.ErrorIdValidation,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.c.Validate()
			if tt.has == nil {
				assert.NoError(t, err, "не должно быть ошибки")
				return
			}
			assert.ErrorIs(t, err, tt.has)
		})
	}
}

func VasyaValid() entity.Customer {
	return entity.Customer{
		Id:       "vasya-top",
		Phone:    "+79162325566",
		Tag:      "cool",
		Operator: "МТС",
		Name:     "Василий Иванович",
	}
}

func VasyaShortNumber() entity.Customer {
	v := VasyaValid()
	v.Phone = "+123123"
	return v
}

func NotPlusNumber() entity.Customer {
	v := VasyaValid()
	v.Phone = "123123"
	return v
}

func VasyaUUID() entity.Customer {
	v := VasyaValid()
	v.Id = "19b70ac1-0ff4-4522-8c03-a0aad1aece65"
	return v
}

func VasyaCyrillic() entity.Customer {
	v := VasyaValid()
	v.Id = "васясупер"
	return v
}

func VasyaUnderscore() entity.Customer {
	v := VasyaValid()
	v.Id = "my_id"
	return v
}
