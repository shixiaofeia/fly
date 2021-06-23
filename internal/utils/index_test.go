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
