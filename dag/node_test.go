package dag

import (
	"fmt"
	"testing"

	"github.com/tendermint/tendermint/types"
)

func TestNodeSerialize(t *testing.T) {
	dagGraph := NewDAGGraph()
	tx := types.Tx([]byte("abc"))
	newNode := dagGraph.AddTx(tx)
	n1 := NodeSerialize(newNode)
	n2 := NodeDeserialize(n1)
	fmt.Println(n1)
	fmt.Println(NodeSerialize(n2))
}
