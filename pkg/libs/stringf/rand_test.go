package stringf

import "testing"

func TestRandStr(t *testing.T) {
	t.Log(string(RandStr(6, Numbers)))
	t.Log(string(RandStr(6, LowerCaseLetters)))
	t.Log(string(RandStr(6, UppercaseLetter)))
	t.Log(string(RandStr(6, NumberPlusCase)))
}

func BenchmarkRandStr(b *testing.B) {
	for i := 0; i < b.N; i++ {
		RandStr(6, NumberPlusCase)
	}
}
