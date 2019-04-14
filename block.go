// Package blockchain provides a simple blockchain implementation.
package blockchain

import (
	"bytes"
	"encoding/gob"
	"time"
)

type Block struct {
	Header *BlockHeader
	Data   []byte
}

// BlockHeader stores metadata relevant to a Block.
type BlockHeader struct {
	Timestamp    int64
	Hash         []byte
	PreviousHash []byte
	Nonce        int
	Difficulty   int
}

func NewBlock(difficulty int, previousHash []byte, data []byte) *Block {
	header := &BlockHeader{
		Timestamp:    time.Now().Unix(),
		PreviousHash: previousHash,
		Difficulty:   difficulty,
	}

	block := &Block{
		Header: header,
		Data:   data,
	}

	pow := NewProofOfWork(block, difficulty)
	nonce, hash := pow.Run()

	header.Hash = hash
	header.Nonce = nonce

	return block
}

func (block *Block) Serialize() ([]byte, error) {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)

	err := encoder.Encode(block)

	if err != nil {
		return nil, err
	}

	return result.Bytes(), nil
}

func DeserializeBlock(d []byte) (*Block, error) {
	var block Block

	decoder := gob.NewDecoder(bytes.NewReader(d))
	err := decoder.Decode(&block)

	if err != nil {
		return nil, err
	}

	return &block, nil
}

func (header *BlockHeader) IsLastBlock() bool {
	return len(header.PreviousHash) == 0
}
