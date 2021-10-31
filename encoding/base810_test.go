package encoding

import (
	"crypto/rand"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBase810(t *testing.T) {
	for i := 0; i < 100000; i++ {
		b := make([]byte, 32)
		rand.Read(b)
		rt, err := Base810.DecodeString(Base810.EncodeToString(b, 79), 32)
		require.NoError(t, err)
		require.Equal(t, Base810.EncodeToString(b, 79), Base810.EncodeToString(rt, 79))
	}
}
