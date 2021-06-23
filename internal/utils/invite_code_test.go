package utils

import "testing"

func TestGetInviteCode(t *testing.T) {
	t.Log(NewInviteCode().IdToCode(12345678))
}

func TestDecodeInviteCode(t *testing.T) {
	t.Log(NewInviteCode().CodeToId("LE338BRR"))
	index := "ABC"
	res := ""
	num := 27
	for num != 0 {
		i := num % 2
		num = num / 2
		res += string(index[i])
	}

	print(res)

}
