package dag

import (
	"crypto/sha256"
	"encoding/hex"
	"math/rand"
	"sort"

	"github.com/tendermint/tendermint/types"
)

type DAGGraph struct {
	nodes     map[string]DAGNode // []byte cannot be a key?
	confirmed map[string]bool    // check whether every node is confirmed
}

func NewDAGGraph() *DAGGraph {
	graph := DAGGraph{}
	graph.nodes = make(map[string]DAGNode)
	graph.confirmed = make(map[string]bool)
	return &graph
}

func calHash(Node DAGNode) string { //compute the hash of Node, include {tx, {ref}, thrpt}
	record := string(Node.tx.Hash()) + string(Node.thrpt)
	for _, preHash := range Node.ref {
		record += preHash
	}
	h := sha256.New()
	h.Write([]byte(record))
	hashed := h.Sum(nil)
	return hex.EncodeToString(hashed)
}

func (graph *DAGGraph) calThrpt(Node DAGNode) uint32 { //use queue to enumerate one's ancestors
	queue := []string{}
	counter := map[string]int{}
	queue = append(queue, Node.ref...)
	for {
		counter[queue[0]] = 1
		queue = queue[1:]
		if len(queue) > 0 {
			newList := graph.nodes[queue[0]].ref
			if len(newList) > 0 {
				queue = append(queue, newList...)
			}
		} else {
			break
		}
	}

	return uint32(len(counter))

}

func (graph *DAGGraph) AddTx(tx types.Tx) DAGNode {
	// Build a new DAGNode of incomming transaction
	newNode := DAGNode{}
	newNode.tx = tx

	parents := graph.SelectTxParents()
	for _, p := range parents {
		newNode.ref = append(newNode.ref, p.hash)
	}

	newNode.thrpt = graph.calThrpt(newNode)
	newNode.hash = calHash(newNode)

	return newNode
}

func (graph *DAGGraph) AddNode(newNode DAGNode) {
	graph.nodes[newNode.hash] = newNode
}

func (graph *DAGGraph) SelectTips() []DAGNode { //Sort current nodes according to their thrpt
	// return an array of DAGNodes with priority
	// called when add new transactions and create consensus proposals
	v := []DAGNode{}
	for _, value := range graph.nodes {
		v = append(v, value)
	}
	sort.Sort(DAGNodeList(v))
	// The more dag nodes are (indirectly) referred by the tip, the higher priority the tip is.
	return v
}

func (graph *DAGGraph) SelectProposal() DAGNode {
	candidates := graph.SelectTips()
	if len(candidates) > 0 {
		return candidates[0]
	} else {
		return DAGNode{}
	}
}

func (graph *DAGGraph) SelectTxParents() []DAGNode {
	// Select optimal reference nodes using SelectTips
	// Add the node to the graph
	candidates := graph.SelectTips()

	// two references per node:
	// One is the tip with highest priority
	// another is a random tip (if any)
	Ref1 := candidates[0]
	if len(candidates) <= 2 {
		return []DAGNode{Ref1}
	} else {
		idx := rand.Intn(len(candidates))
		for ; idx == 0; idx = rand.Intn(len(candidates)) {
		}
		Ref2 := candidates[idx]
		return []DAGNode{Ref1, Ref2}
	}
}

func (graph *DAGGraph) Commit(node DAGNode) {
	// Accept the hash of confirmed DAGNode from consensus;
	// update DAG; update confirmed number for calculation of throughput
	graph.confirmed[node.hash] = true
}

func (graph *DAGGraph) IsValid(Node DAGNode) bool {
	// check avaliability of parents: if parents of this node are not learned?
	// ignore other sanity checks
	queue := []string{}
	queue = append(queue, Node.ref...)
	for {
		if graph.confirmed[queue[0]] == false {
			return false
		}
		queue = queue[1:]
		if len(queue) > 0 {
			newList := graph.nodes[queue[0]].ref
			if len(newList) > 0 {
				queue = append(queue, newList...)
			}
		} else {
			break
		}
	}
	return true
}
