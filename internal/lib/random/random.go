package random

import (
	"math/rand"
	"strings"
	"time"
)

const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func NewRandomString(aliasLength int) string {
	source := rand.NewSource(time.Now().UnixNano())
	randomGenerator := rand.New(source)

	var alias []string

	for i := 0; i < aliasLength; i++ {
		alias = append(alias, string(chars[randomGenerator.Intn(len(chars))]))
	}

	return strings.Join(alias, "")

}
