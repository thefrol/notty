package valid

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_HasNoExplicitLanguage(t *testing.T) {
	tests := []struct {
		s       string
		wantErr error
	}{
		{"хуй-пизд", ErrorExplicitLanguage},
		{"маленький зайчик", nil},
	}
	for _, tt := range tests {
		t.Run(tt.s, func(t *testing.T) {
			err := HasNoExplicitLanguage(tt.s)
			assert.ErrorIs(t, err, tt.wantErr)

		})
	}
}
