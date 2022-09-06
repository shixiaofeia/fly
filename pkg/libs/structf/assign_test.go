package structf

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

type (
	User struct {
		Name string
		Age  int
	}
	Info struct {
		Name string
		Age  int
	}
)

func TestAssign(t *testing.T) {

	userM := &User{Name: "aa", Age: 18}
	infoM := &Info{}
	Assign(userM, infoM, "Age")

	assert.Equal(t, infoM.Name, userM.Name)
	assert.Equal(t, infoM.Age, 0)
}

func BenchmarkAssign(b *testing.B) {
	userM := &User{Name: "aa", Age: 18}
	infoM := &Info{}
	for i := 0; i < b.N; i++ {
		Assign(userM, infoM, "Age")
	}
}

func TestAssignUnmarshal(t *testing.T) {
	userRecords := []*User{
		{Name: "a", Age: 11},
		{Name: "b", Age: 12},
	}
	infoRecords := make([]*Info, 0)

	AssignUnmarshal(&userRecords, &infoRecords)

	for k, v := range infoRecords {
		assert.Equal(t, userRecords[k].Name, v.Name)
		assert.Equal(t, userRecords[k].Age, v.Age)
	}
}

func BenchmarkAssignUnmarshal(b *testing.B) {
	userRecords := []*User{
		{Name: "a", Age: 11},
		{Name: "b", Age: 12},
	}
	infoRecords := make([]*Info, 0)

	for i := 0; i < b.N; i++ {
		AssignUnmarshal(&userRecords, &infoRecords)
	}
}
