package blockchain

import (
	"log"
	"os"
	"testing"
)

const (
	testDbDir  = ".testing_db"
	difficulty = 10 // Easy difficulty for testing purposes.
)

// MUST defer bc.Close() during the test or test teardown will panic.
func testSetup(t *testing.T) func(t *testing.T) {
	removeTestDB()

	return func(t *testing.T) {
		removeTestDB()
	}
}

func removeTestDB() {
	if _, err := os.Stat(testDbDir); !os.IsNotExist(err) {
		err = os.RemoveAll(testDbDir)
		if err != nil {
			log.Panic(err)
		}
	}
}

func TestBlockchainHappyPath(t *testing.T) {
	testTearDown := testSetup(t)
	defer testTearDown(t)

	log.Println("Test start.")

	bc, err := NewBlockChain(difficulty, testDbDir, nil)
	if err != nil {
		log.Println(err)
	}
	defer bc.Close()

	log.Println("Blockchain created.")

	msg1 := "John has 2 more PrestigeCoin than Jane"
	msg2 := "Jane has 10 more PrestigeCoin than David"

	err = bc.AddBlock([]byte(msg1))
	if err != nil {
		log.Println(err)
	}

	err = bc.AddBlock([]byte(msg2))
	if err != nil {
		log.Println(err)
	}

	log.Println("Blocks added.")

	bci := bc.Iterator()

	currBlock, _ := bci.Next()
	currMsg := string(currBlock.GetData())
	if currMsg != msg2 {
		t.Errorf("Block held incorrect data. Expected: %s but got %s", msg2, currMsg)
	}

	currBlock, _ = bci.Next()
	currMsg = string(currBlock.GetData())
	if string(currBlock.GetData()) != msg1 {
		t.Errorf("Block held incorrect data. Expected: %s but got %s", msg1, currMsg)
	}

	currBlock, _ = bci.Next()
	currMsg = string(currBlock.GetData())
	if string(currBlock.GetData()) != "Genesis Block" {
		t.Errorf("Block held incorrect data. Expected: %s but got %s", "Genesis Block", currMsg)
	}

	currBlock, _ = bci.Next()
	if currBlock != nil {
		t.Errorf("Blockchain should be of length 3 but was not.")
	}
}
