package main

import (
	"os"

	"github.com/harmony-one/harmony/staking/network"
)

func main() {
	nodeID := os.Args[1] // Here is the name of the company
	server := network.NewServer(nodeID)
	server.Start()

}

func NewServer(nodeID string) *Server {
	// A new node is created according to the passed NodeID. The default view of the node is 1000000, and the node starts three processes

}
