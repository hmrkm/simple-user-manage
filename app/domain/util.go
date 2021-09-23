package domain

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

func CreateULID(now time.Time) string {
	entropy := ulid.Monotonic(rand.New(rand.NewSource(now.UnixNano())), 0)
	id := ulid.MustNew(ulid.Timestamp(now), entropy)
	return id.String()
}

// ページャーの現在位置と最終位置を解決
// count: 全件
// page: 現在ページ
// limit: 表示件数
func Pager(count int, page int, limit int) (now int, last int) {
	if limit <= 0 {
		return 1, 1
	}

	last = count / limit

	if fraction := count % limit; fraction != 0 {
		last++
	}

	if page <= 0 {
		now = 1
	} else if count > (page-1)*limit {
		now = page
	} else {
		now = last
	}

	return now, last
}
