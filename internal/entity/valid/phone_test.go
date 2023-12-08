package valid

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPhone(t *testing.T) {
	tests := []struct {
		s       string
		wantErr error
	}{
		{"+79163332211", nil},
		{"83332221111", ErrorPhoneValidation}, // todo можно попробовать конвертировать
		{"+73332211", ErrorPhoneValidation},
		{"+733322a2233", ErrorPhoneValidation},
	}
	for _, tt := range tests {
		t.Run(tt.s, func(t *testing.T) {
			err := Phone(tt.s)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}
