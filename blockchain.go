package blockchain

import (
	"fmt"
	"log"
	"os"

	"github.com/boltdb/bolt"
)

const (
	dbFile       = "main.db"
	blocksBucket = "blocksBucket"
	headBlock    = "head"
)

type Blockchain struct {
	metadata   *BlockchainMetadata
	head       []byte
	difficulty int
	db         *bolt.DB
}

type BlockchainMetadata struct {
	DbPath string
}

type BlockchainIterator struct {
	currentHash []byte
	db          *bolt.DB
}

func NewBlockChain(dbPath string, difficulty int, genesisData []byte) (*Blockchain, error) {
	fullDbPath := fmt.Sprintf("%s/%s", dbPath, dbFile)
	if _, err := os.Stat(fullDbPath); os.IsNotExist(err) {
		err = os.MkdirAll(dbPath, 0700)
		if err != nil {
			return nil, err
		}
	}

	db, err := bolt.Open(fullDbPath, 0600, nil)
	if err != nil {
		return nil, err
	}

	var encodedHead []byte
	err = db.Update(func(tx *bolt.Tx) error {
		blocksBucketKey := []byte(blocksBucket)
		headBlockKey := []byte(headBlock)

		b := tx.Bucket(blocksBucketKey)

		if b == nil {
			b, err = tx.CreateBucket(blocksBucketKey)
			if err != nil {
				return err
			}

			genesisBlock := NewGenesisBlock(genesisData)
			encodedGenesisBlock, err := genesisBlock.Serialize()
			if err != nil {
				return err
			}

			genesisBlockHash := genesisBlock.Header.Hash

			err = b.Put(genesisBlockHash, encodedGenesisBlock)
			if err != nil {
				return err
			}

			err = b.Put(headBlockKey, genesisBlockHash)
			if err != nil {
				return err
			}

			encodedHead = genesisBlockHash
		} else {
			encodedHead = b.Get(headBlockKey)
		}

		return nil
	})

	metadata := &BlockchainMetadata{
		DbPath: dbPath,
	}

	blockchain := &Blockchain{
		metadata:   metadata,
		head:       encodedHead,
		difficulty: difficulty,
		db:         db,
	}

	return blockchain, nil
}

func NewGenesisBlock(data []byte) *Block {
	if data == nil {
		data = []byte("Genesis Block")
	}
	return NewBlock(
		0,
		nil,
		data)
}

func (bc *Blockchain) AddBlock(data []byte) (*Block, error) {
	var finishedNewBlock *Block
	var newBlockHash []byte
	err := bc.db.Batch(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))

		newBlock := NewBlock(bc.difficulty, bc.head, data)

		encodedNewBlock, err := newBlock.Serialize()
		if err != nil {
			return err
		}

		newBlockHash = newBlock.Header.Hash

		// Blocks should never be updated after initial creation.
		if b.Get(newBlockHash) != nil {
			return &collisionError{newBlockHash}
		}

		err = b.Put(newBlockHash, encodedNewBlock)
		if err != nil {
			return err
		}

		// Last operation that can fail. If fails, need to remove the added block but functionality of the blockchain remains unaffected.
		err = b.Put([]byte(headBlock), newBlockHash)
		if err != nil {
			return err
		}

		// Only set head when function cannot fail anymore. Since BoltDB processes batches sequentially,
		// head will be correct for successful additions and not be updated for failed ones.
		bc.head = newBlockHash

		finishedNewBlock = newBlock
		return nil
	})

	if err != nil {
		bc.deleteBadBlock(newBlockHash)
		return nil, err
	}

	return finishedNewBlock, nil
}

func (bc *Blockchain) deleteBadBlock(hash []byte) {
	go func() {
		bc.db.Batch(func(tx *bolt.Tx) error {
			b := tx.Bucket([]byte(blocksBucket))

			// Not a read-only transaction so should always return nil
			return b.Delete(hash)
		})
	}()
}

func (bc *Blockchain) Difficulty() int {
	return bc.difficulty
}

func (bc *Blockchain) Close() error {
	err := bc.db.Close()
	return err
}

func (bc *Blockchain) Iterator() *BlockchainIterator {
	bci := &BlockchainIterator{bc.head, bc.db}
	return bci
}

func DeleteBlockchain(bc *Blockchain) {
	bc.Close()

	if _, err := os.Stat(bc.metadata.DbPath); !os.IsNotExist(err) {
		err = os.RemoveAll(bc.metadata.DbPath)
		if err != nil {
			log.Panic(err)
		}
	}
}

func (bci *BlockchainIterator) Next() (*Block, error) {
	if bci.currentHash == nil {
		return nil, nil
	}

	var encodedBlock []byte

	err := bci.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(blocksBucket))
		encodedBlock = b.Get(bci.currentHash)

		return nil
	})

	if err != nil {
		return nil, err
	}

	if encodedBlock == nil {
		return nil, fmt.Errorf("No block found for the hash %v", bci.currentHash)
	}

	block, err := DeserializeBlock(encodedBlock)
	if err != nil {
		return nil, err
	}

	bci.currentHash = block.Header.PreviousHash

	return block, nil
}
