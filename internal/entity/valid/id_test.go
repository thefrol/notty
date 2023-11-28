package valid

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Id(t *testing.T) {
	tests := []struct {
		s       string
		wantErr error
	}{
		{"ok-id", nil},
		{"ok", nil},
		{"a977f79b-e923-4bdd-86a7-7be0d8c6e06a", nil},
		{"bad_id", ErrorIdValidation},
		{"айди", ErrorIdValidation},
	}
	for _, tt := range tests {
		t.Run(tt.s, func(t *testing.T) {
			err := Id(tt.s)
			assert.ErrorIs(t, err, tt.wantErr)
		})
	}
}
