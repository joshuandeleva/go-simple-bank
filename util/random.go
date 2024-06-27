package util

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"


func init(){
	rand.Seed(time.Now().UnixNano())
}

func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

func RandomString(n int) string {
	var sb strings.Builder

	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}
	return sb.String()

}

// random owner

func RandomOwner() string {
	return RandomString(6)
}

// random amount money

func RandomMoney() int64 {
	return RandomInt(0, 1000)
}

// random currency
func RandomCurrency() string {
	currencies := []string{EUR, USD, CAD, AUD, JPY, CNY, KSH, GBP}
	n := len(currencies)
	return currencies[rand.Intn(n)]
}

func RandomEamil() string {
	return fmt.Sprintf("%s@email.com", RandomString(6))
}