package sb3

import (
	"math/rand/v2"
	"strings"
)

const (
	idCharSet = "!#%()*+,-./:;=?@[]^_`{|}~ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
	idLength  = 20
)

func generateId() string {
	var stringBuilder strings.Builder
	for range idLength {
		stringBuilder.WriteByte(idCharSet[rand.IntN(len(idCharSet))])
	}
	return stringBuilder.String()
}
