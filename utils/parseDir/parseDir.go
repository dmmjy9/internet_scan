package parseDir

import (
	"bufio"
	"internet_scan/utils/splitIP"
	"io"
	"io/ioutil"
	"os"
)

func GetFilesName(dir string) ([]string, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return []string{}, err
	}
	var filename []string
	for _, val := range files {
		filename = append(filename, val.Name())
	}
	return filename, nil
}

func GetIPSeg(filename string) ([]string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return []string{}, err
	}
	line := bufio.NewReader(file)
	var IPSeg []string
	for {
		content, _, err := line.ReadLine()
		if err == io.EOF {
			break
		}
		ipseg, err1 := splitIP.SplitIP(string(content))
		if err1 != nil {
			continue
		}
		IPSeg = append(IPSeg, ipseg...)
	}
	file.Close()
	return IPSeg, nil
}
