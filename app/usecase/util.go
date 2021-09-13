package usecase

import (
	"crypto/sha256"
	"encoding/hex"
	"math/rand"
	"time"

	"github.com/oklog/ulid"
)

func CreateHash(src string) string {
	sha256ByteArr := sha256.Sum256([]byte(src))
	return hex.EncodeToString(sha256ByteArr[:])
}

func CreateULID() string {
	t := time.Now()
	entropy := ulid.Monotonic(rand.New(rand.NewSource(t.UnixNano())), 0)
	id := ulid.MustNew(ulid.Timestamp(t), entropy)
	return id.String()
}
