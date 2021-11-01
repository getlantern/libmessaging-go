package identity

import (
	"crypto"
	"crypto/ed25519"
	"crypto/rand"
	"testing"

	"github.com/getlantern/libmessaging-go/encoding"
	"github.com/stretchr/testify/require"
)

func TestRoundTrip(t *testing.T) {
	publicKey, err := PublicKeyFromString("rfu2495fqazzpq1e3xkj1skmr9785hwbxggpr17ut1htj4h9nhyy")
	require.NoError(t, err)
	require.Equal(t, "2277029271600308397119018701998194490680040839333862997699030902896411310611021743", publicKey.Number())

	_privateKey, err := encoding.HumanFriendlyBase32Encoding.DecodeString("jkrbbfgym19yz79saxym4mfqxbhzxtndf9r98m76upcxkgyr83cs54x5asgry4x6czscwkakgw476q7mudzgsug1kqrd83t466n1w4e")
	require.NoError(t, err)
	privateKey := ed25519.PrivateKey(_privateKey)

	data := []byte("hello world")
	signature, err := privateKey.Sign(rand.Reader, data, crypto.Hash(0))
	require.NoError(t, err)

	require.True(t, publicKey.Verify(data, signature))
}

func TestNumber(t *testing.T) {
	b := make([]byte, 32)
	rand.Read(b)
	pk := PublicKey(b)
	rt, err := PublicKeyFromNumber(pk.Number())
	require.NoError(t, err)
	require.Equal(t, pk.Number(), rt.Number())
}
