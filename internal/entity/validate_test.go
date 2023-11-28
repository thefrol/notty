package entity

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"gitlab.com/thefrol/notty/internal/entity/valid"
)

// Этот тест показывает, что если мы соберем массив нил ошибок
// и заджоиним их, то все будет хорошо, массив нилов не даст ошибку
// в конце
//
// От инефрфейсов бывают странные артефакты
// лучше это сразу иметь в виду
func Test_AppendNil(t *testing.T) {
	var errs []error

	errs = append(errs, nil, nil)
	assert.Equal(t, 2, len(errs))

	errs = append(errs, valid.Phone("+71112223344"))

	joined := errors.Join(errs...)
	assert.NoError(t, joined)
}
