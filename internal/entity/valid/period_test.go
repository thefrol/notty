package valid

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPeriod(t *testing.T) {
	tests := []struct {
		name       string
		start, end string // время в формате ГГГГ.ММ.ДД
		wantErr    error
	}{
		{
			name:    "нормально",
			start:   "2013.03.01",
			end:     "2013.03.02",
			wantErr: nil,
		},
		{
			name:    "начало раньше конца",
			start:   "2013.03.02",
			end:     "2013.03.01",
			wantErr: ErrorInvalidPeriod,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s, err := time.Parse("2006.01.02", tt.start)
			require.NoError(t, err)

			e, err := time.Parse("2006.01.02", tt.end)
			require.NoError(t, err)

			err = Period(s, e)
			assert.ErrorIs(t, tt.wantErr, err)
		})
	}
}
