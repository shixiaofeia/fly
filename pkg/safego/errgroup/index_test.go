package errgroup

import (
	"context"
	"fmt"
	"testing"
)

func TestWithContext(t *testing.T) {
	g, _ := WithContext(context.Background())

	g.Go(func() error {
		fmt.Println("111")
		return nil
	})
	g.Go(func() error {
		return fmt.Errorf("go err")
	})
	if err := g.Wait(); err != nil {
		t.Error(err)
	}
	t.Log("end")
}
