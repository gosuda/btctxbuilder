package utils

import (
	"encoding/hex"
	"strings"
)

func Encode(bytes []byte) string {
	return hex.EncodeToString(bytes)
}

func EncodeWith0x(bytes []byte) string {
	return AddPrefix(Encode(bytes))
}

func MustDecode(hexString string) []byte {
	data, err := Decode(hexString)
	if err != nil {
		panic(err)
	}
	return data
}

func Decode(hexString string) ([]byte, error) {
	hexString = TrimPrefix(hexString)

	if len(hexString)%2 != 0 {
		hexString = "0" + hexString
	}

	return hex.DecodeString(hexString)
}

func IsHex(str string) bool {
	str = TrimPrefix(str)

	if len(str)%2 != 0 {
		return false
	}

	for i := 0; i < len(str); i++ {
		if !isHexCharacter(str[i]) {
			return false
		}
	}
	return true
}

func HasPrefix(str string) bool {
	return strings.HasPrefix(str, "0x") || strings.HasPrefix(str, "0X")
}

func TrimPrefix(s string) string {
	if HasPrefix(s) {
		return s[2:]
	}
	return s
}

func AddPrefix(s string) string {
	if !HasPrefix(s) {
		return "0x" + s
	}
	return s
}

func isHexCharacter(c byte) bool {
	return ('0' <= c && c <= '9') || ('a' <= c && c <= 'f') || ('A' <= c && c <= 'F')
}
