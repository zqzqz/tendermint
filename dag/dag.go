package dag

import (
	"github.com/tendermint/tendermint/types"
)

type DAGNode struct {
	tx     types.Tx
	hash   []byte
	ref    [][]byte
	nounce uint32
}

type DAGGraph struct {
	nodes     map[string]DAGNode // []byte cannot be a key?
	confirmed uint64
}

func NewDAGGraph() *DAGGraph {
	return &DAGGraph{}
}

func (graph *DAGGraph) AddTx(tx types.Tx) DAGNode {
	// Build a new DAGNode of incomming transaction
	// Select optimal reference nodes using SelectTips
	// Add the node to the graph
	
	// two references per node:
	// One is the tip with highest priority
	// another is a random tip (if any)
	return DAGNode{}
}

func (graph *DAGGraph) SelectTips() []DAGNode {
	// return an array of DAGNodes with priority
	// called when add new transactions and create consensus proposals
	
	// The more dag nodes are (indirectly) referred by the tip, the higher priority the tip is.
	return []DAGNode{}
}

func (graph *DAGNode) Commit(hash []byte) {
	// Accept the hash of confirmed DAGNode from consensus;
	// update DAG; update confirmed number for calculation of throughput
}

func (graph *DAGNode) IsValid(node DAGNode) bool {
	// check avaliability of parents: if parents of this node are not learned?
	// ignore other sanity checks
	return true
}
