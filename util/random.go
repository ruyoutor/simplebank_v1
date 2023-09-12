package util

import (
	"math/rand"
	"strings"
	"time"
)

const alphabet = "qwertyuiopasdfghjklzxcvbnm"

func init() {
	rand.Seed(time.Now().UnixNano())
}

// RandonInit generates a random integer between min and max
func RandonInit(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

// RandonString generates a random string of length n
func RandonString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

//RandomOwner generates a random owner name
func RandonOwn() string {
	return RandonString(6)
}

//RandomMoney generates a random amout of money
func RandomMoney() int64 {
	return RandonInit(0, 1000)
}

// RandomCurrency generates a random currency code
func RandomCurrency() string {
	currencies := []string{USD, EUR, BRL, CAD}
	l := len(currencies)
	return currencies[rand.Intn(l)]
}
