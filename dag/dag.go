package dag

import (
	"crypto/sha256"
	"encoding/hex"
	"math/rand"
	"sort"

	"github.com/tendermint/tendermint/types"
)

type DAGNode struct {
	tx   types.Tx
	hash string
	ref  []string // the ref for geneisus block is empty
	//nounce uint32
	thrpt uint32
}

type DAGNodeList []DAGNode

func (a DAGNodeList) Len() int           { return len(a) }
func (a DAGNodeList) Less(i, j int) bool { return a[i].thrpt > a[j].thrpt }
func (a DAGNodeList) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

type DAGGraph struct {
	nodes     map[string]DAGNode // []byte cannot be a key?
	confirmed map[string]bool    // check whether every node is confirmed
}

func NewDAGGraph() *DAGGraph {
	return &DAGGraph{}
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
	queue := make([]string, 0)
	counter := make(map[string]uint32)
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
	// Select optimal reference nodes using SelectTips
	// Add the node to the graph
	DNL := graph.SelectTips()
	Ref1 := DNL[0]
	newNode := DAGNode{}
	newNode.tx = tx
	if len(DNL) <= 2 {
		newNode.ref = []string{Ref1.hash}
	} else {
		idx := rand.Intn(len(DNL))
		for idx == 0 {
			idx = rand.Intn(len(DNL))
		}
		Ref2 := DNL[idx]
		newNode.ref = []string{Ref1.hash, Ref2.hash}
	}

	newNode.thrpt = graph.calThrpt(newNode)
	newNode.hash = calHash(newNode)
	// two references per node:
	// One is the tip with highest priority
	// another is a random tip (if any)

	return newNode
}

func (graph *DAGGraph) AddNode(newNode DAGNode) {
	graph.nodes[newNode.hash] = newNode
}

func (graph *DAGGraph) SelectTips() []DAGNode { //Sort current nodes according to their thrpt
	// return an array of DAGNodes with priority
	// called when add new transactions and create consensus proposals
	v := make([]DAGNode, len(graph.nodes))
	for _, value := range graph.nodes {
		v = append(v, value)
	}
	sort.Sort(DAGNodeList(v))
	// The more dag nodes are (indirectly) referred by the tip, the higher priority the tip is.
	return v
}

func (graph *DAGGraph) Commit(hash string) {
	// Accept the hash of confirmed DAGNode from consensus;
	// update DAG; update confirmed number for calculation of throughput
	graph.confirmed[hash] = true
}

func (graph *DAGGraph) IsValid(Node DAGNode) bool {
	// check avaliability of parents: if parents of this node are not learned?
	// ignore other sanity checks
	queue := make([]string, 0)
	queue = append(queue, Node.ref...)
	for {
		if !graph.confirmed[queue[0]] {
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
