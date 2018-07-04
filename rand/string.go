package rand

import (
	"bytes"
	"math/rand"
)

const validChars = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

func AlphanumericString(length int) string {
	buf := bytes.NewBuffer(nil)

	for i := 0; i < length; i++ {
		ch := validChars[rand.Intn(len(validChars))]
		buf.WriteByte(ch)
	}

	return buf.String()
}
