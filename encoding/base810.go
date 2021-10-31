package encoding

import (
	"errors"
	"math/big"
)

var (
	ErrInvalidBase810 = errors.New("Invalid Base810 string")

	base8Table = map[byte]rune{
		byte(0): '2',
		byte(1): '3',
		byte(2): '4',
		byte(3): '5',
		byte(4): '6',
		byte(5): '7',
		byte(6): '8',
		byte(7): '9',
	}

	base8TableReverse = make(map[rune]byte, len(base8Table))

	Base810 = &Base810Encoding{}
)

func init() {
	for key, value := range base8Table {
		base8TableReverse[value] = key
	}
}

// Base810Encoding provides a human-friendly encoding that looks like a phone number but isn't usually a dialable
// phone number, because it doesn't start with 0 or 1 (as is required in most countries). This
// encoding treats a byte array as a big-endian number. The first (most significant) 3 bits of data
// are encoded using a shifted octal encoding (digits 2-9 instead of the standard 0-7) and the
// remaining data is encoded in base10 and left-padded with '0's to meet the desired length.
type Base810Encoding struct{}

// EncodeToString encodes the given bytes using Base810 encoding. The resulting string will be of target length
// using '0's after the first digit in order to pad up to the targetLength.
func (e *Base810Encoding) EncodeToString(b []byte, targetLength int) string {
	_b := make([]byte, len(b))
	copy(_b, b)
	head := base8Table[_b[0]>>5]
	_b[0] = _b[0] << 3
	tail := big.NewInt(0).SetBytes(_b).String()
	result := make([]rune, 79)
	result[0] = head
	padding := targetLength - 1 - len(tail)
	if padding < 0 {
		padding = 0
	}
	// left-pad tail with zeroes to make sure we reach our desired length
	for i := 0; i < padding; i++ {
		result[i+1] = '0'
	}
	for i, r := range tail {
		result[i+1+padding] = r
	}
	return string(result)
}

// DecodeString decodes the given Base810 string into a byte[] of targetSize. If the string doesn't contain
// enough data to fill targetSize, the byte[] will contain leading zeros.
func (e *Base810Encoding) DecodeString(s string, targetSize int) ([]byte, error) {
	head := base8TableReverse[rune(s[0])]
	_tail, ok := big.NewInt(0).SetString(s[1:], 10)
	if !ok {
		return nil, ErrInvalidBase810
	}
	tail := make([]byte, 32)
	_tail.FillBytes(tail)
	tail[0] = head<<5 | tail[0]>>3
	return tail, nil
}
