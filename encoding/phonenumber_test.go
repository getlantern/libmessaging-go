package encoding

import (
	"crypto/rand"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestShiftBase9(t *testing.T) {
	s := "012345678"
	require.Equal(t, "012346789", shiftBase9(s))
	require.Equal(t, s, unshiftBase9(shiftBase9(s)))
}

func TestPhoneNumber(t *testing.T) {
	for i := 0; i < 100000; i++ {
		b := make([]byte, 32)
		rand.Read(b)
		str := PhoneNumber.EncodeToString(b, 82)
		// prepend some ignored characters and an interior 5 to make sure they don't interfere with decoding
		str = fmt.Sprintf("015789%s5%s", str[:12], str[12:])
		rt, err := PhoneNumber.DecodeString(str, 32)
		require.NoError(t, err)
		require.Equal(t, PhoneNumber.EncodeToString(b, 82), PhoneNumber.EncodeToString(rt, 82))
	}
}
