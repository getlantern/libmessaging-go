package encoding

import (
	"errors"
	"math/big"
	"strings"
)

var (
	ErrInvalidPhoneNumber = errors.New("Invalid phone number string")

	base4Table = make(map[byte]rune, 4)

	base4TableReverse = make(map[rune]byte, 4)

	PhoneNumber = &PhoneNumberEncoding{}
)

func addBase4Mapping(b byte, c rune) {
	base4Table[b] = c
	base4TableReverse[c] = b
}

func init() {
	addBase4Mapping(0, '2')
	addBase4Mapping(1, '3')
	addBase4Mapping(2, '4')
	addBase4Mapping(3, '6')
}

// PhoneNumberEncoding provides a human-friendly encoding that looks like a phone number but isn't usually a dialable
// phone number, because it doesn't start with 0 or 1 (as is required in most countries). This
// encoding treats a byte array as a big-endian number. The first (most significant) 2 bits of data
// are encoded using a modified base4 encoding (digits 2, 3, 4, 6 instead of the standard 0-4) and the
// remaining data is encoded in base9 (omitting digit 5) and left-padded with '0's to meet the desired length.
//
// This encoding permits the inclusion of arbitrary '5's anywhere in the encoded string, which it simply ignores. This
// can be used to visually differentiate the beginning of two otherwise very similar numbers, for example, given:
//
// 2222222222222222222222222222222222222222222222222222222222222222222222222222222
// 2222222222222222222222222222222222222222222222222222222222222222222222222222223
//
// We can change the 2nd number to the following equivalent value
//
// 522222222222252222222222222222222222222222222222222222222222222222222222222222223
//
type PhoneNumberEncoding struct{}

// EncodeToString encodes the given bytes using PhoneNumber encoding. The resulting string will be of targetLength
// padding up to the targetLength with '0's following the 1st digit.
func (e *PhoneNumberEncoding) EncodeToString(b []byte, targetLength int) string {
	_b := make([]byte, len(b))
	copy(_b, b)
	head := base4Table[_b[0]>>6]
	_b[0] = _b[0] << 2
	tail := shiftBase9(big.NewInt(0).SetBytes(_b).Text(9))
	if targetLength < len(tail)+1 {
		targetLength = len(tail) + 1
	}
	result := make([]rune, targetLength)
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

// DecodeString decodes the given PhoneNumber string into a byte[] of targetSize. If the string doesn't contain
// enough data to fill targetSize, the byte[] will contain leading zeros.
//
// This function ignores any leading characters other than 2, 3, 4 or 6, and any subsequent 5s.
func (e *PhoneNumberEncoding) DecodeString(s string, targetSize int) ([]byte, error) {
	s = strings.TrimLeft(s, "015789")
	head := base4TableReverse[rune(s[0])]
	_tail, ok := big.NewInt(0).SetString(unshiftBase9(s[1:]), 9)
	if !ok {
		return nil, ErrInvalidPhoneNumber
	}
	tail := make([]byte, targetSize)
	_tail.FillBytes(tail)
	tail[0] = head<<6 | tail[0]>>2
	return tail, nil
}

func shiftBase9(s string) string {
	result := make([]rune, 0, len(s))
	for _, c := range s {
		if c < '5' {
			result = append(result, c)
		} else {
			result = append(result, c+1)
		}
	}
	return string(result)
}

func unshiftBase9(s string) string {
	result := make([]rune, 0, len(s))
	for _, c := range s {
		if c < '5' {
			result = append(result, c)
		} else if c == '5' {
			// ignore 5s
		} else {
			result = append(result, c-1)
		}
	}
	return string(result)
}
