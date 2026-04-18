package engine

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"math/rand"
)

type Direction int

const (
	Left Direction = iota
	Right
)

type ProvablyFairEngine struct {
	ServerSeed string
	ClientSeed string
	Nonce      int
}

func NewEngine(s, c string, n int) *ProvablyFairEngine {
	return &ProvablyFairEngine{s, c, n}
}

func (e *ProvablyFairEngine) GenerateHash() string {
	combined := fmt.Sprintf("%s:%s:%d", e.ServerSeed, e.ClientSeed, e.Nonce)
	return fmt.Sprintf("%x", sha256.Sum256([]byte(combined)))
}

// CalculatePath derives a deterministic path from the SHA-256 hash of
// serverSeed:clientSeed:nonce. Pure p=0.5, fully verifiable externally.
func (e *ProvablyFairEngine) CalculatePath(rows int) []Direction {
	hash := sha256.Sum256([]byte(fmt.Sprintf("%s:%s:%d", e.ServerSeed, e.ClientSeed, e.Nonce)))
	seed := int64(binary.BigEndian.Uint64(hash[:8]))
	r := rand.New(rand.NewSource(seed))

	path := make([]Direction, rows)
	for i := 0; i < rows; i++ {
		if r.Float64() < 0.5 {
			path[i] = Right
		} else {
			path[i] = Left
		}
	}
	return path
}
