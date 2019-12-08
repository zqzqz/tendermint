package dag

import (
	"bytes"
	"github.com/tendermint/tendermint/types"
)

type DAGNode struct {
	tx   types.Tx
	hash string
	ref  []string // the ref for geneisus block is empty
	//nounce uint32
	thrpt int
}


func NodeSerialize(node DAGNode) []byte {
	split := byte('#')
	res := []byte{}
	res = append(res, node.tx...)
	res = append(res, split)
	res = append(res, []byte(node.hash)...)
	for _, r := range node.ref {
		res = append(res, split)
		res = append(res, []byte(r)...)
	}
	return res
}

func NodeDeserialize(txBytes []byte) DAGNode {
	node := DAGNode{}
	tokens := bytes.Split(txBytes, []byte("#"))
	for index, token := range tokens {
		if index == 0 {
			node.tx = token
		} else if index == 1 {
			node.hash = string(token)
		} else {
			node.ref = append(node.ref, string(token))
		}
	}
	return node
}
