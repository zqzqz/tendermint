package dag

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/rand"
	"sync"

	"github.com/tendermint/tendermint/types"
)

type DAGGraph struct {
	nodes          map[string]DAGNode // []byte cannot be a key?
	confirmed      map[string]bool    // check whether every node is confirmed
	is_tip         map[string]bool
	cache          map[string]DAGNode
	pendingCommits map[string]DAGNode
	mux            sync.Mutex
}

func NewDAGGraph() *DAGGraph {
	graph := DAGGraph{}
	graph.nodes = make(map[string]DAGNode)
	graph.confirmed = make(map[string]bool)
	graph.is_tip = make(map[string]bool)
	graph.cache = make(map[string]DAGNode)
	graph.pendingCommits = make(map[string]DAGNode)
	genesis := DAGNode{}
	genesis.thrpt = 0
	genesis.hash = calHash(genesis)
	graph.nodes[genesis.hash] = genesis
	graph.is_tip[genesis.hash] = true
	graph.confirmed[genesis.hash] = true
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

func (graph *DAGGraph) calThrpt(Node DAGNode) int { //use queue to enumerate one's ancestors
	queue := []string{Node.hash}
	counter := map[string]int{}

	for {
		if len(queue) > 0 {
			counter[queue[0]] = 1
			graph.mux.Lock()
			newList := graph.nodes[queue[0]].ref
			for _, n := range newList {
				if graph.confirmed[n] == false {
					_, ok := counter[n]
					if ok == false {
						queue = append(queue, n)
					}
				}
			}
			graph.mux.Unlock()
			queue = queue[1:]
		} else {
			break
		}
	}

	return len(counter)

}

func (graph *DAGGraph) AddTx(tx types.Tx) DAGNode {
	// Build a new DAGNode of incomming transaction
	newNode := DAGNode{}
	newNode.tx = tx
	for _, node := range graph.SelectTxParents() {
		newNode.ref = append(newNode.ref, node.hash)
	}
	newNode.thrpt = graph.calThrpt(newNode)
	newNode.hash = calHash(newNode)

	return newNode
}

func (graph *DAGGraph) AddNode(newNode DAGNode) {
	if _, ok := graph.nodes[newNode.hash]; ok {
		return
	}
	valid := graph.IsValid(newNode)
	if !valid {
		graph.mux.Lock()
		graph.cache[newNode.hash] = newNode
		graph.mux.Unlock()
		return
	}

	graph.mux.Lock()
	graph.nodes[newNode.hash] = newNode
	graph.is_tip[newNode.hash] = true
	graph.mux.Unlock()

	if _, pending := graph.pendingCommits[newNode.hash]; pending {
		graph.mux.Lock()
		delete(graph.pendingCommits, newNode.hash)
		graph.mux.Unlock()
		graph.Commit(newNode)
	}

	queue := make([]string, 0)
	for _, ref := range newNode.ref {
		queue = append(queue, ref)
	}
	counter := map[string]int{}

	for len(queue) > 0 {
		counter[queue[0]] = 1
		graph.mux.Lock()
		graph.is_tip[queue[0]] = false
		newList := graph.nodes[queue[0]].ref
		for _, n := range newList {
			if graph.is_tip[n] == true {
				_, ok := counter[n]
				if ok == false {
					queue = append(queue, n)
				}
			}
		}
		graph.mux.Unlock()
		queue = queue[1:]
	}

	graph.mux.Lock()
	for h, Node := range graph.cache {
		flag := graph.IsValid(Node)
		if flag == true {
			delete(graph.cache, h)
			graph.nodes[Node.hash] = Node
			if _, pending := graph.pendingCommits[Node.hash]; pending {
				delete(graph.pendingCommits, Node.hash)
				graph.Commit(Node)
			}
			graph.is_tip[Node.hash] = true
			for _, p := range Node.ref {
				_, ok := graph.is_tip[p]
				if ok == true {
					delete(graph.is_tip, p)
				}
			}
		}
	}
	graph.mux.Unlock()
}

func (graph *DAGGraph) SelectTips() []DAGNode { //Sort current nodes according to their thrpt
	// return an array of DAGNodes with priority
	// called when add new transactions and create consensus proposals
	res := []DAGNode{}
	v := []DAGNode{}
	max_tip := DAGNode{}
	k := -1
	graph.mux.Lock()
	for value, flag := range graph.is_tip {
		if flag {
			v = append(v, graph.nodes[value])
			if graph.nodes[value].thrpt > k {
				max_tip = graph.nodes[value]
				k = graph.nodes[value].thrpt
			}
		}
	}
	graph.mux.Unlock()
	res = append(res, max_tip)
	for _, value := range v {
		if value.hash != max_tip.hash {
			res = append(res, value)
		}
	}
	// The more dag nodes are (indirectly) referred by the tip, the higher priority the tip is.
	return res
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
	if len(candidates) <= 1 {
		return candidates
	} else {
		Ref1 := candidates[0]
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
	if _, ok := graph.nodes[node.hash]; !ok {
		graph.pendingCommits[node.hash] = node
		return
	}

	queue := []string{node.hash}
	counter := map[string]int{}
	for len(queue) > 0 {
		counter[queue[0]] = 1
		graph.mux.Lock()
		graph.confirmed[queue[0]] = true
		newList := graph.nodes[queue[0]].ref
		for _, n := range newList {
			if graph.confirmed[n] == false {
				_, ok := counter[n]
				if ok == false {
					queue = append(queue, n)
				}
			}
		}
		graph.mux.Unlock()
		queue = queue[1:]
	}
	fmt.Printf("Current graph throughput: %d\n", len(graph.confirmed))
}

func (graph *DAGGraph) IsValid(Node DAGNode) bool {
	// check avaliability of parents: if parents of this node are not learned?
	// ignore other sanity checks
	parents := Node.ref
	graph.mux.Lock()
	for _, p := range parents {
		_, ok := graph.nodes[p]
		if ok == false {
			return false
		}
	}
	graph.mux.Unlock()
	return true
}
