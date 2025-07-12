package block

import (
	"crypto/sha256"
	"encoding/hex"
	"strconv"
	"time"
)

type Block struct {
	Index        int
	Timestamp    string
	Data         string
	Hash         string
	PreviousHash string
}

func NewBlockchain(data string) []Block {
	genesisBlock := Block{
		Index:     0,
		Timestamp: time.Now().Format(time.RFC3339),
		Data:      data,
	}
	genesisBlock.Hash = CalculateHash(genesisBlock)
	return []Block{genesisBlock}
}

func CalculateHash(b Block) string {
	data := strconv.Itoa(b.Index) +
		b.Timestamp + b.Data + b.PreviousHash
	hash := sha256.New()
	hash.Write([]byte(data))
	hashed := hash.Sum(nil)
	return hex.EncodeToString(hashed)
}

func GenerateNextBlock(prev Block, data string) Block {
	nextBlock := Block{
		Index:        prev.Index + 1,
		Timestamp:    time.Now().Format(time.RFC3339),
		Data:         data,
		PreviousHash: prev.Hash,
	}
	nextBlock.Hash = CalculateHash(nextBlock)
	return nextBlock
}
