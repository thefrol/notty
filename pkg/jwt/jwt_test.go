package auth

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_JWTExpired(t *testing.T) {
	jwt, err := BuildExpirable(time.Microsecond, "123")
	require.NoError(t, err)

	time.Sleep(10 * time.Microsecond)
	err = Validate(jwt, "123")
	assert.Error(t, err)
}

func Test_JWTUnexpired(t *testing.T) {
	jwt, err := BuildExpirable(time.Second, "123")
	require.NoError(t, err)

	err = Validate(jwt, "123")
	assert.NoError(t, err)
}
