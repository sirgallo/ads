package lfnodepooltests

import "testing"

import "github.com/sirgallo/ads/pkg/node"


func TestNodePool(t *testing.T) {
	poolsize := 10000
	poolRange := make([]int, poolsize)
	np := node.NewLFNodePool[string](poolsize)

	lfNodeArr := make([]*node.LFNode[string], 0, poolsize)
	
	for range poolRange {
		newNode := np.GetLFNode()
		lfNodeArr = append(lfNodeArr, newNode)
	}

	for _, val := range lfNodeArr {
		np.PutLFNode(val)
	}

	t.Logf("actual poolsize: %d, expected poolsize: %d\n", np.PoolSize.GetValue(), poolsize)
	if np.PoolSize.GetValue() != int64(poolsize) {
		t.Errorf("actual poolsize not expected poolsize: actual(%d), expected(%d)", np.PoolSize, poolsize)
	}

	for range poolRange {
		node := np.GetLFNode()
		if node == nil {
			t.Error("unable to get node")
		}
	}

	t.Logf("actual poolsize: %d, expected poolsize: %d\n", np.PoolSize.GetValue(), 0)
	if int(np.PoolSize.GetValue()) != 0 {
		t.Errorf("actual poolsize not expected poolsize: actual(%d), expected(%d)", np.PoolSize.GetValue(), 0)
	}
}