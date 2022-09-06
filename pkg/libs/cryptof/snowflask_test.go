package cryptof

import (
	"fmt"
	"testing"
)

func TestNode_Generate(t *testing.T) {
	n, err := NewNode(1)
	if err != nil {
		t.Error(err.Error())
		return
	}

	for i := 0; i < 10000; i++ {
		fmt.Println(n.Generate())
	}
}

func BenchmarkNewNode(b *testing.B) {
	n, _ := NewNode(1)

	for i := 0; i < b.N; i++ {
		n.Generate()
	}
}
