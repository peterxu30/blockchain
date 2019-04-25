package blockchain

import "fmt"

type collisionError struct {
	hash []byte
}

func (e *collisionError) Error() string {
	return fmt.Sprintf("A block with this hash already exists. Hash: %s", e.hash)
}