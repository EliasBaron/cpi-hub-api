package helpers

import (
	"math/rand"

	"github.com/oklog/ulid/v2"
)

func NewULID() string {
	now := NowBuenosAires()
	entropy := rand.New(rand.NewSource(now.UnixNano()))
	ms := ulid.Timestamp(now)
	result, _ := ulid.New(ms, entropy)
	return result.String()
}
