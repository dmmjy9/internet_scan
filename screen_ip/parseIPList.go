package screen_ip

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
)

func GetIPSplit(filename string) ([][]string, error) {
	var	ipSplit	[][]string
	var	ipSlice	[]string

	// Open file and read to slice.
	fh, err := os.Open("ips_connected/" + filename)
	if err != nil {
		fmt.Println("open file " + filename + " fail")
		return [][]string{}, err
	}
	contents, err := ioutil.ReadAll(fh)
	if err != nil {
		fmt.Println("read file " + filename + " fail")
		return [][]string{}, nil
	}
	ctBytes := bytes.Split(contents, []byte{'\n'})
	for _, val := range ctBytes {
		if len(string(val)) != 0 {
			ipSlice = append(ipSlice, string(val))
		}
	}

	// Parse if ip count <= 100.
	if len(ipSlice) <= 100 {
		ipSplit = append(ipSplit, ipSlice)
		return ipSplit, nil
	}

	// Parse if ip count > 100.
	var ipTmp []string
	for idx := range ipSlice {
		ipTmp = append(ipTmp, ipSlice[idx])
		if len(ipTmp) == 100 {
			ipSplit = append(ipSplit, ipTmp)
			ipTmp = []string{}
			continue
		}
	}
	ipSplit = append(ipSplit, ipTmp)

	return ipSplit, nil
}
