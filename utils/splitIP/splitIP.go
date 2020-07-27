package splitIP

import (
	"github.com/3th1nk/cidr"
	"math"
	"strconv"
	"strings"
)

func SplitIP(ipSeg string) ([]string, error) {
	var ips []string
	maskLen, _ := strconv.Atoi(strings.Split(ipSeg, "/")[1])
	subNetNum := math.Pow(2, float64(24 - maskLen))
	c, _ := cidr.ParseCIDR(ipSeg)
	cs1, _ := c.SubNetting(cidr.SUBNETTING_METHOD_SUBNET_NUM, int(subNetNum))
	for _, ip := range cs1 {
		ips = append(ips, ip.CIDR())
	}
	return ips, nil
}
