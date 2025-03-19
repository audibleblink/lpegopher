package collectors

import (
	"encoding/hex"
	"fmt"
	"hash"
	"sync"

	"github.com/audibleblink/logerr"
	"github.com/minio/highwayhash"

	"github.com/audibleblink/lpegopher/util"
)

// Hasher provides an interface for hashing strings
type Hasher interface {
	HashString(data string, normalize bool) (string, error)
}

// HighwayHasher implements Hasher using HighwayHash algorithm
type HighwayHasher struct {
	key       []byte
	hashPool  *sync.Pool
	normalize bool
}

// NewHighwayHasher creates a new HighwayHasher with the given key and normalization preference
func NewHighwayHasher(key []byte, normalize bool) *HighwayHasher {
	return &HighwayHasher{
		key:       key,
		normalize: normalize,
		hashPool: &sync.Pool{
			New: func() any {
				h, err := highwayhash.New(key)
				if err != nil {
					// This is a rare case that would only happen with an invalid key
					log := logerr.Add("hash")
					log.Errorf("Failed to create HighwayHash instance: %v", err)
					return nil
				}
				return h
			},
		},
	}
}

// HashString hashes a string with optional normalization
func (h *HighwayHasher) HashString(data string, normalize bool) (string, error) {
	// Normalize the data if requested
	if normalize {
		data = util.PathFix(data)
	}

	// Get a hash instance from the pool
	hashInstance := h.hashPool.Get()
	if hashInstance == nil {
		return "", fmt.Errorf("could not obtain hash instance from pool")
	}
	hash := hashInstance.(hash.Hash)
	defer func() {
		hash.Reset()
		h.hashPool.Put(hash)
	}()

	// Write data directly to hash
	_, err := hash.Write([]byte(data))
	if err != nil {
		return "", fmt.Errorf("hash write failed: %w", err)
	}

	// The HighwayHash output size is 32 bytes
	checksum := hash.Sum(nil)
	return hex.EncodeToString(checksum), nil
}

// Variables for thread-safe initialization of the default hasher
var (
	defaultHasherOnce sync.Once
	defaultHasher     *HighwayHasher
)

// initDefaultHasher initializes the default hasher with proper synchronization
func initDefaultHasher() {
	defaultHasherOnce.Do(func() {
		defaultHasher = NewHighwayHasher(key, true)
	})
}

// hashFor generates a hash string for the given data using the default hasher
// For backward compatibility, this function maintains the original signature
// but uses the more advanced implementation internally
func hashFor(data string) string {
	// Initialize hasher if needed
	if defaultHasher == nil {
		initDefaultHasher()
	}

	result, err := defaultHasher.HashString(data, true)
	if err != nil {
		log := logerr.Add("hash")
		log.Errorf("Hashing failed: %v", err)
		return ""
	}
	return result
}

// HashWithOptions provides a more flexible API for hashing with options
// It returns both the hash and any error that occurred
func HashWithOptions(data string, normalize bool) (string, error) {
	// Initialize hasher if needed
	if defaultHasher == nil {
		initDefaultHasher()
	}

	return defaultHasher.HashString(data, normalize)
}
