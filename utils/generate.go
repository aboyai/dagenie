package utils

import (
	"encoding/hex"
	"math/rand"
	"sync"
	"time"
)

var counter uint32
var mu sync.Mutex

func GenerateObjectID() string {
	mu.Lock()
	defer mu.Unlock()

	now := uint32(time.Now().Unix())
	randPart := rand.Uint64() & 0xFFFFFFFFFFFF // 6 bytes
	counter++

	// Compose 12-byte ID: 4 bytes time + 6 bytes rand + 2 bytes counter
	id := make([]byte, 12)
	id[0] = byte(now >> 24)
	id[1] = byte(now >> 16)
	id[2] = byte(now >> 8)
	id[3] = byte(now)

	id[4] = byte(randPart >> 40)
	id[5] = byte(randPart >> 32)
	id[6] = byte(randPart >> 24)
	id[7] = byte(randPart >> 16)
	id[8] = byte(randPart >> 8)
	id[9] = byte(randPart)

	id[10] = byte(counter >> 8)
	id[11] = byte(counter)

	return hex.EncodeToString(id)
}
