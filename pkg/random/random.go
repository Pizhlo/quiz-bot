package random

import "golang.org/x/exp/rand"

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ01234567890123456789"

func Int(min, max int) int {
	return min + rand.Intn(max-min+1)
}

func String(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = charset[rand.Intn(len(charset))]
	}
	return string(b)
}

func Strings(n int) []string {
	res := []string{}

	for i := 0; i < n; i++ {
		res = append(res, String(5))
	}

	return res
}
