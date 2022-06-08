package structf

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAssign(t *testing.T) {
	type User struct {
		Name string
		Age  int
	}
	type Info struct {
		Name string
		Age  int
	}
	u := &User{Name: "aa", Age: 18}
	i := &Info{}
	Assign(u, i, "Age")

	assert.Equal(t, i.Name, u.Name)
	assert.Equal(t, i.Age, 0)
}
