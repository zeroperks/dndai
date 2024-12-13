package dice

import (
	"crypto/rand"
	"encoding/binary"
	"log"
	"time"

	"golang.org/x/crypto/chacha20"
)

type SeedGenerator struct{}

// GenerateSeed creates a cryptographic seed using ChaCha20
func (sg *SeedGenerator) GenerateSeed(userID int) uint64 {
	// Generate a 256-bit key for ChaCha20
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		log.Fatalf("Failed to generate key: %v", err)
	}

	// Create a nonce (12 bytes) using user ID and current time
	nonce := make([]byte, 12)
	binary.LittleEndian.PutUint32(nonce[:4], uint32(userID))
	binary.LittleEndian.PutUint64(nonce[4:], uint64(time.Now().UnixNano()))

	// Initialize ChaCha20 cipher
	chacha, err := chacha20.NewUnauthenticatedCipher(key, nonce)
	if err != nil {
		log.Fatalf("Failed to initialize ChaCha20: %v", err)
	}

	// Generate 8 bytes of random data for the seed
	seedBytes := make([]byte, 8)
	chacha.XORKeyStream(seedBytes, seedBytes)

	// Convert the random bytes to a uint64 seed
	return binary.LittleEndian.Uint64(seedBytes)
}
