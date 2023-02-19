package testing

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestRand(t *testing.T) {
	rand.Seed(time.Now().UnixNano())
	fmt.Println(rand.Intn(3))
	fmt.Println(rand.Perm(3))
}
