package helpers

import (
	"math/rand"
	"time"

	"github.com/oklog/ulid/v2"
)

func NewULID() string {
	entropy := rand.New(rand.NewSource(time.Now().UnixNano()))
	ms := ulid.Timestamp(time.Now())
	result, _ := ulid.New(ms, entropy)
	return result.String()
}
