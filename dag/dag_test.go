package dag

import (
	hex "encoding/hex"
	"fmt"
	"testing"

	"github.com/tendermint/tendermint/types"
)

func Test_dag(t *testing.T) {
	dagGraph := NewDAGGraph()
	for i := 0; i < 1000; i++ {
		tx := types.Tx([]byte(string(i)))
		newNode := dagGraph.AddTx(tx)
		fmt.Println(hex.EncodeToString(NodeSerialize(newNode)))
		dagGraph.AddNode(newNode)
		if i%100 == 0 {
			proposal := dagGraph.SelectProposal()
			dagGraph.Commit(proposal)
		}
	}
}
