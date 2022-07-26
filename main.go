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
	// A new node is created according to the passed NodeID. The default view of the node is 1000000, and the node starts three processes: dispatchMsg,
	// alarmToDispatcher and resolveMsg
	node := NewNode(nodeID)
	//Start service
	server := &Server{node.NodeTable[nodeID], node}

	//Set route
	server.setRoute()
	return server

	//1. Create a node according to the passed nodeID
	//2. Create a server service
	//3. Set route
}


// NewNode(node ID)
func NewNode(nodeID string) *Node {
	const viewID = 1000000000 // temporary.

	node := &Node{
		// Hard-coded for test
		NodeID: nodeID,
		NodeTable: map[string]string{
			"Apple"

		}
	}


}