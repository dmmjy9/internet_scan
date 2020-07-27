package pingScan

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"sync"
)

func PingScan(ipSeg string) ([]string, error) {
	//pingRet := make(map[string]string)
	pingcmd := "fping -aeg " + ipSeg

	cmd := exec.Command("/bin/bash", "-c", pingcmd)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return []string{}, err
	}

	var pingRets []string

	cmd.Start()

	reader := bufio.NewReader(stdout)
	for {
		line, err2 := reader.ReadString('\n')
		if err2 != nil || io.EOF == err2 {
			break
		}
		pingRets = append(pingRets, strings.Split(line, "\n")[0])
	}
	cmd.Wait()
	return pingRets, nil
}

func IPFileScanAlive(filename string) ([]string, error) {
	var workerSlice []string
	var ipSegAlive []string
	worker := make(chan string, 200000)
	goLimit := make(chan struct{}, 500)
	ipAlive := make(chan string, 200000)

	var wg sync.WaitGroup

	fh, err := os.Open("ips_ret/" + filename)
	if err != nil {
		fmt.Println("open file " + filename + " fail")
		return []string{}, err
	}

	contents, err := ioutil.ReadAll(fh)
	if err != nil {
		fmt.Println("read file " + filename + " fail")
		return []string{}, err
	}
	ctBytes := bytes.Split(contents, []byte{'\n'})
	for _, val := range ctBytes {
		workerSlice = append(workerSlice, string(val))
	}


	fh.Close()

	wg.Add(1)
	go func(worker chan string) {
		defer wg.Done()
		for _, val := range workerSlice {
			worker <- val
		}
		close(worker)
	}(worker)

	wg.Add(1)
	fmt.Printf("start scan: %s\n", filename)
	go func(ipAlive chan string) {
		defer wg.Done()
		var wg2 sync.WaitGroup
		for ip := range worker {
			goLimit <- struct{}{}
			wg2.Add(1)
			go func(ipAlive chan string, _ip string) {
				defer func() {
					wg2.Done()
					<-goLimit
				}()
				fmt.Printf("start scan ipseg: %s, from file: %s\n", _ip, filename)
				ipsAlive, err := PingScan(_ip)
				if err != nil {
					fmt.Println("scan error: ", _ip, err, ", skip!")
					return
				}
				if len(ipsAlive) != 0 {
					fmt.Printf("find alive ipseg: %s, from file: %s\n", _ip, filename)
					ipAlive <- _ip
				} else {
					return
				}
			}(ipAlive, ip)
		}
		wg2.Wait()
		close(ipAlive)
	}(ipAlive)

	wg.Wait()

	for {
		if ip, ok := <-ipAlive; ok {
			ipSegAlive = append(ipSegAlive, ip)
		} else {
			break
		}
	}
	return ipSegAlive, nil
}
