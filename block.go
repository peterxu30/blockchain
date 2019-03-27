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

type BlockHeader struct {
	Timestamp  int64
	Hash       []byte
	PrevHash   []byte
	Nonce      int
	Difficulty int
}

func NewBlock(difficulty int, prevHash []byte, data []byte) *Block {
	header := &BlockHeader{
		Timestamp:  time.Now().Unix(),
		PrevHash:   prevHash,
		Difficulty: difficulty,
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

func (block *Block) GetTimestamp() int64 {
	return block.Header.GetTimestamp()
}

func (block *Block) GetHash() []byte {
	return block.Header.GetHash()
}

func (block *Block) GetPreviousHash() []byte {
	return block.Header.GetPreviousHash()
}

func (block *Block) GetNonce() int {
	return block.Header.GetNonce()
}

func (block *Block) GetDifficulty() int {
	return block.Header.GetDifficulty()
}

func (block *Block) GetData() []byte {
	return block.Data
}

func (block *Block) IsLastBlock() bool {
	return block.Header.IsLastBlock()
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

func (header *BlockHeader) GetTimestamp() int64 {
	return header.Timestamp
}

func (header *BlockHeader) GetHash() []byte {
	return header.Hash
}

func (header *BlockHeader) GetPreviousHash() []byte {
	return header.PrevHash
}

func (header *BlockHeader) GetNonce() int {
	return header.Nonce
}

func (header *BlockHeader) GetDifficulty() int {
	return header.Difficulty
}

func (header *BlockHeader) IsLastBlock() bool {
	return len(header.PrevHash) == 0
}
