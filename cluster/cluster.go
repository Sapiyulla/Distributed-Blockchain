package cluster

import (
	"distributed-blockchain-net/block"
	"distributed-blockchain-net/mode"
	"distributed-blockchain-net/node"
	"errors"
	"log"
	"slices"
	"strconv"
	"sync"
)

type Cluster struct {
	Mode  mode.Mode
	Nodes []*node.Node
}

func NewCluster() Cluster {
	return Cluster{
		Mode:  mode.ModeAP,
		Nodes: []*node.Node{node.NewNode(0)},
	}
}

type ClusterContract interface {
	AddNode()
	SetMode(mode mode.Mode)
	SetNodeAlive(id int, alive bool)
	WriteTransaction(data string) error
	ReadLastBlock() (block.Block, error)
	PrintChains()
}

func (c *Cluster) AddNode() {
	c.Nodes = append(c.Nodes, node.NewNode(len(c.Nodes)))
}

func (c *Cluster) SetMode(_mode mode.Mode) {
	if _mode != mode.ModeAP && _mode != mode.ModeCA && _mode != mode.ModeCP {
		log.Println("invalid mode: not setted")
	}
	c.Mode = _mode
}

func (c *Cluster) SetNodeAlive(id int, alive bool) {
	if len(c.Nodes) >= id {
		log.Println("cluster: node with id=" + strconv.Itoa(id) + " not found: not setted")
	}
	c.Nodes[id].Alive = alive
}

func (c *Cluster) WriteTransaction(data string) error {
	mu := sync.Mutex{}
	switch c.Mode {
	case mode.ModeAP, mode.ModeCP:
		aliveNodesCount := 0
		mu.Lock()
		defer mu.Unlock()
		for i := range c.Nodes {
			if c.Nodes[i].Alive {
				aliveNodesCount++
				c.Nodes[i].Chain = append(c.Nodes[i].Chain, block.GenerateNextBlock(c.Nodes[i].Chain[len(c.Nodes[i].Chain)-1], data))
			}
		}
		if aliveNodesCount == 0 {
			return errors.New("cluster: not found alive node(s)")
		}
	default:
		for i := range c.Nodes {
			if !c.Nodes[i].Alive {
				return errors.New("cluster: not aviable")
			}
		}
		for i := range c.Nodes {
			c.Nodes[i].Chain = append(c.Nodes[i].Chain, block.GenerateNextBlock(c.Nodes[i].Chain[len(c.Nodes[i].Chain)-1], data))
		}
	}
	return nil
}

func (c *Cluster) ReadLastBlock() (*block.Block, error) {
	mu := sync.Mutex{}
	switch c.Mode {
	case mode.ModeCP:
		mu.Lock()
		defer mu.Unlock()
		aliveNodes := []*node.Node{}
		for i := range c.Nodes {
			if c.Nodes[i].Alive {
				aliveNodes = append(aliveNodes, c.Nodes[i])
			}
		}
		if len(aliveNodes) < (len(c.Nodes)/2)+1 {
			return nil, errors.New("cluster: not aviable")
		}
		trueBlockHashCounts := map[string]int{}
		for i := range aliveNodes {
			trueBlockHashCounts[aliveNodes[i].Chain[len(aliveNodes[i].Chain)-1].Hash]++
		}
		maxKey := ""
		maxVal := 0
		for k, v := range trueBlockHashCounts {
			if v > maxVal {
				maxKey = k
				maxVal = v
			}
		}
		for i := range aliveNodes {
			if aliveNodes[i].Chain[len(aliveNodes[i].Chain)-1].Hash == maxKey && maxVal >= (len(c.Nodes)/2)+1 {
				return &aliveNodes[i].Chain[len(aliveNodes[i].Chain)-1], nil
			}
		}
		return nil, errors.New("cluster: inconsistent quorum")
	case mode.ModeAP:
		mu.Lock()
		defer mu.Unlock()
		IdNodeWithHotChainBlock := 0
		aliveNodesCount := 0
		for i := range c.Nodes {
			if c.Nodes[i].Alive {
				aliveNodesCount++
				if len(c.Nodes[i].Chain) > IdNodeWithHotChainBlock {
					IdNodeWithHotChainBlock = i
				}
			}
		}
		if aliveNodesCount < 1 {
			return nil, errors.New("cluster: not found alive node(s)")
		}
		return &c.Nodes[IdNodeWithHotChainBlock].Chain[len(c.Nodes[IdNodeWithHotChainBlock].Chain)-1], nil
	case mode.ModeCA:
		mu.Lock()
		defer mu.Unlock()
		lengthMaxNodeBlockchain := []int{}
		for i := range c.Nodes {
			if !c.Nodes[i].Alive {
				return nil, errors.New("cluster: not aviable")
			}
			lengthMaxNodeBlockchain = append(lengthMaxNodeBlockchain, len(c.Nodes[i].Chain))
		}
		for i := range c.Nodes {
			if slices.Max(lengthMaxNodeBlockchain) == len(c.Nodes[i].Chain) {
				return &c.Nodes[i].Chain[len(c.Nodes[i].Chain)-1], nil
			}
		}
		return nil, errors.New("cluster: please, retry later")
	default:
		return nil, errors.New("cluster: unavaible mode")
	}
}

func (c *Cluster) PrintChains() {
	for i := range c.Nodes {
		if c.Nodes[i].Alive {
			log.Printf("🟢 Node %d (alive)", c.Nodes[i].ID)
			for _, block := range c.Nodes[i].Chain {
				log.Printf("Block %d: Hash=%s PrevHash=%s Data=%s", block.Index, block.Hash, block.PreviousHash, block.Data)
			}
			log.Println("")
		} else {
			log.Printf("🔴 Node %d (not aviable)", c.Nodes[i].ID)
			log.Println("")
		}
	}
}
