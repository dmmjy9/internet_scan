package screen_ip

import (
	"fmt"
	"testing"
)

func TestScreenIP(t *testing.T) {
	parseRet, err := GetIPSplit("ipstest2.txt")
	if err != nil {
		panic(err)
	}
	for _, sp := range parseRet {
		fmt.Println(len(sp), sp)
	}
}
