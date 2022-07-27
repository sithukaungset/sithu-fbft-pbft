package main

import (
	"fmt"
	"os"

	"github.com/harmony-one/harmony/consensus"
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
			"Apple":  "localhost:1111",
			"MS":     "localhost:1112",
			"Google": "localhost:1113",
			"IBM":    "localhost:1114",
		},
		View: &View{
			ID:      viewID,
			Primary: "Apple", // The primary node is Apple

		},
		// Consensus-related struct
		CurrentState:  nil,
		CommittedMsgs: make([]*consensus.RequestMsg, 0), // Submitted information
		MsgBuffer: &MsgBuffer{
			ReqMsgs:        make([]*consensus.RequestMsg, 0),
			PrePrepareMsgs: make([]*consensus.PrePrepareMsg, 0),
			PrePrepareMsgs: make([]*consensus.VoteMsg, 0),
			CommitMsgs:     make([]*consensus.VoteMsg, 0),
		},
		// Channels
		MsgEntrance: make(chan interface{}), // Unbuffered information receiving channel
		MsgDelivery: make(chan interface{}), // Unbuffered message transmission channel
		Alarm:       make(chan bool),
	}

	// Start message scheduler
	go node.dispatchMsg()

	// Start alarm trigger
	go node.alarmToDispatcher()

	// Start information voting
	go node.resolveMsg()

	return node

}

// The NewNode here actually creates a new node and defines some attrbute information in the node structure.
// The node structure is as follows:

// node
type Node struct {
	NodeID        string
	NodeTable     map[string]string // key=nodeID, value=url
	View          *View
	CurrentState  *consensus.State
	CommittedMsgs []*consensus.RequestMsg // kinda block.
	MsgBuffer     *MsgBuffer
	MsgEntrance   chan interface{}
	MsgDelivery   chan interface{}
	Alarm         chan bool
}

// Analyzing NewNode
// Initializing view number const viewid = 100000000
// Defines some attribute information in the node structure

// Each node has three goroutines open
// dispatchMsg
// alarmToDispatcher
// resolveMsg

// Lets focus on two cooperative processes goroutine

//collaboration 1: dispatchMsg

func (node *Node) dispatchMsg() {
	for {
		select {
		case msg := <-node.MsgEntrance: // If a message is sent from the MsgEntrance channel, get msg
			err := node.routeMsg(msg) // routeMsg
			if err != nil {
				fmt.Println(err)
				// TODO: send err to ErrorChannel

			}

		case <-node.Alarm:
			err := node.routeMsgWhenAlarmed()
			if err != nil {
				fmt.Println(err)
				// ToDo: send err to ErrorChannel
			}
		}
	}

}

// Here we need to know something about for selecting multiplexing

// We can see from the code of dispatchMsg that as long as there is a value in the MsgEnrance channel, it will be passed to an intermediate variable
// message and then the message will be routed and forwarded to routeMsg.

func (node *Node) routeMsg(msg interface{}) []error {
	switch msg.(type) {
	case *consensus.RequestMsg:
		if node.CurrentState == nil {
			// Copy buffered messages first.
			msgs := make([]*consensus.RequestMsg, len(node.MsgBuffer.ReqMsgs))
			copy(msgs, node.MsgBuffer.ReqMsgs)

			// Append a newly arrived message.
			msgs = append(msgs, msg.(*consensus.RequestMsg))

			// Empty the buffer.
			node.MsgBuffer.ReqMsgs = make([]*consensus.RequestMsg, 0)

			// Send messages.
			node.MsgDelivery <- msgs

		} else {
			node.MsgBuffer.ReqMsgs = append(node.MsgBuffer.ReqMsgs, msg.(*consensus.RequestMsg))
		}
		//... Other types of information data
		return nil
	}

// Collaboration 2: resolveMsg

func (node *Node) resolveMsg() {
		for {
			// Get cache information from scheduler
			msgs := <-node.MsgDelivery
			switch msgs.(type) {
			case []*consensus.RequestMsg:
				// Node voting decision information
				errs := node.resolveRequestMsg(msgs.([]*consensus.RequestMsg))
				if len(errs) != 0{
					for _, err := range errs {
						fmt.Println(err)
					}
					// ToDo: send err to ErrorChannel
				}
				//... For other types of information data
			}
		}
// Here we can see clearly MsgDelivery is waiting for the dispatcher dispatchMsg to deliver messages.
// Because the two buffer channels are unbuffered, no messages will be blocked here all the time

// Here is how to execute the reolveRequestMsg function:
	
// Information of node voting request stage
func (node *Node) resolveRequestMsg(msgs []*consensus.RequestMsg) []error {
	errs := make([]error, 0)

	// Voting information
	for _, reqMsg := range msgs {
		err := node.GetReq(reqMsg)
		if err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) != 0{
		return errs
	}
	return nil

}
// In the resolveRequestMsg code block, we can see that the main execution logic is to call the GetReq function

// The master node starts the global consensus
func (node *Node) GetReq(reqMsg *consensus.RequestMsg) error {
	LogMsg(reqMsg)

	// Create a new state for the new consensus.
	err := node.createStateForNewConsensus()
	if err != nil {
		return err
	}
	// Start consensus process
	prePrepareMsg, err := node.CurrentState.StartConsensus(reqMsg)
	if err != nil {
		return err
	}
	LogState(fmt.Sprintf("Consensus Process (ViewID:%d)", node.CurrentState.ViewID), false)
	
	// Getprepare information found
	// Here, the master node starts sending pre preparation messages to other nodes
	if perPrepareMsg != nil {
		node.Broadcast(prePrepareMsg, "/preprepare")
		LogStage("Pre-prepare", true)
	}
	return nil
}
	// 1. Create a new consensus state createStateForNewConsensus
	// 2. Start consensus on the current state of the node startconnections
	// 3. After the consensus is completed, the master node broadcasts the information of the preparemsg phase to other nodes!

	
}

	// 2.2 setRoute function
	func (server *Server) setRoute() {
		http.HandleFunc("/req", server.getReq)
		http.HandleFunc("/preprepare", server.getPrePrepare)
		http.HandleFunc("/prepare", server.getPrepare)
		http.HandleFunc("/commit", server.getCommit)
		http.HandleFunc("/reply", server.getReply)
	}

	// 3. Request consensus process
	// 3.1 create status createStateForNewConsensus

	func (node *Node) createStateForNewConsensus() error {
		// Check if there is an ongoing consensus process.
		if node.CurrentState != nil {
			return errors.New("another consensus is ongoing")


		}

		// Get the last sequence ID

		var lastSequenceID int64
		if len(node.CommittedMsgs) == 0 {
				lastSequenceID = -1
		}
		else {
			lastSequenceID = node.CommittedMsgs[len(node.CommittedMsgs) - 1].lastSequenceID
		}

		// Create a new state for this new consensus process in the Primary
		node.CurrentState = consensus.CreateState(node.View.ID, lastSequenceID)


		LogStage("Create the replica status", true)

		return nil



	}
	// It is preferred to judge whether the status of the current node is nil, that is, whether the current node is in other stages(pre preparation stage or preparation stage, etc.)
	// Judge whether messages have been sent in the current stage. If consensus is reached for the first time, the lastSequenceID of the previous serial number is set to -1, otherwise we take 
	// out the previous serial number.
	// Create a statecreatestate, which is mained returned after initializing the State 

	func CreateState(viewID int64, lastSequenceID int64) *State {

		return &State{
			ViewID: viewID, 
			MsgLogs: &MsgLogs{
				ReqMsg: nil, 
				PrepareMsgs: make(map[string]*VoteMsg),
				commitMsgs: make(map[string]*VoteMsg),
			},
			LastSequenceID: lastSequenceID,
			CurrentStage: Idle,
		}
	}
	// Start Consensus
	
	// Master node start consensus
	func (state *State) StartConsensus(request *RequestMsg) (*PrePrepareMsg, error){
		// The sequence number of the message is sequenceID, that is, the current time
		sequenceID := time.Now().UnixNano()

		// Find the unique and largest number for the sequence ID

		if state.LastSequenceID != -1 {
			for state.LastSequenceID >= sequenceID {
				sequenceID += 1
			}
		}
		// Assign a new sequence ID to the request message object.
		request.SequenceID = sequenceID 
		
		// Save ReqMsgs to its logs.
		state.MsgLogs.ReqMsg = request 
	
		// Get the summary of the client request message requestMsg
		digest, err := digest(request)
		if err != nil {
			fmt.Println(err)

		}

		// Convert the current phase to the prepared phase
		state.CurrentStage = PrePrepared

		// this is actually the format of messages sent by the master node to other nodes: View ID, request information and summary of request information
		return &PrePrepareMsg{
			ViewID: state.ViewID,
			SequenceID: sequenceID,
			Digest: digest,
			RequestMsg: request, 

		}, nil
	}
	// First, use the current timestamp as the serial number
	// If the previous serial number >= the current serial number, the current serial number + 1. The purpose is very simple. Every time the master node starts a consensus, the serial number + 1
	// Update the serial number SequenceID, obtain the summary digest(request) of the request information, and the convert the current status to prepared


	// Node broadcast function
 
	func (node *Node) Broadcast(msg interface{}, path string) map[string]error {
		// Because we dont need to broadcast to yourself, we just skip it 
		if nodeID == node.NodeID {
			continue

		}
		// Encoding msg information into json format
		jsonMsg, err := json.Marshal(msg)
		if err != nil {
			errorMap[nodeID] = err 
			continue
		}
		// Pass json format to other nodes
		send(url + path, jsonMsg) // URL localhost: 1111 path: / prepare, etc
		// send function: http Post("http://"+url, "application/json", buff)

	}
	if len(errorMap) == 0{
		return nil 
	}else {
		return errorMap
	}

}
// Traverse the node list that comes with node initialization NodeTable
// Encode msg information into json format Marshal(msg), and then send it to other nodes(URL + path, jsonmsg)


func send(url string, msg []byte) {
	buff := bytes.NewBuffer(msg)
	_, _ = http.Post("http://"+ url, "applicaiton/json", buff)
}

