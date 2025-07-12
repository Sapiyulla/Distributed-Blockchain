package node

import (
	"distributed-blockchain-net/block"
)

type Node struct {
	ID    int
	Alive bool
	Chain []block.Block
}

func NewNode(ID int) *Node {
	return &Node{
		ID:    ID,
		Alive: true,
		Chain: block.NewBlockchain(""),
	}
}

type NodeContract interface {
	AddBlock(data string)
	LastBlock() block.Block
}

func (c *Node) AddBlock(data string) {
	nextBlock := block.GenerateNextBlock(c.Chain[len(c.Chain)-1], data)
	c.Chain = append(c.Chain, nextBlock)
}
