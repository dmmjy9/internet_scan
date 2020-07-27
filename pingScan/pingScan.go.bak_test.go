package pingScan

import (
	"fmt"
	"testing"
)

func TestPingScan(t *testing.T) {
	ipseg := "114.114.114.0/24"
	rets, _ := PingScan(ipseg)
	fmt.Println(len(rets))
	for _, val := range rets {
		fmt.Println(val)
	}
}
