package utils

import (
	"testing"
)

func TestKRand(t *testing.T) {
	t.Log(string(KRand(6, Numbers)))
	t.Log(string(KRand(6, lowerCaseLetters)))
	t.Log(string(KRand(6, UppercaseLetter)))
	t.Log(string(KRand(6, NumberPlusCase)))

}

func BenchmarkGetRandomString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.Log(string(KRand(9, NumberPlusCase)))
	}
}

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
	t.Log(i)
}

func TestToMd5(t *testing.T) {
	t.Log(ToMd5("123"))
}

func TestDeepCopy(t *testing.T) {
	c1 := map[string]int{
		"tom": 11,
	}
	c2 := c1
	c2["tom"] = 1
	t.Log(c1, c2)
	c3 := make(map[string]int)
	if err := DeepCopy(c1, &c3); err != nil {
		t.Error(err)
	}
	c3["tom"] = 2
	t.Log(c1, c2, c3)

}
