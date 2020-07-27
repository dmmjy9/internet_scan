package main

import (
	"bufio"
	"fmt"
	"internet_scan/pingScan"
	"internet_scan/screen_ip"
	"internet_scan/utils/parseDir"
	"math/rand"
	"os"
	"sync"
	"time"
)

func ParseIPFiles() {
	files, err := parseDir.GetFilesName("ips/")
	if err != nil {
		panic(err)
	}
	for _, val := range files {
		func() {
			dstFile, err := os.OpenFile("ips_ret/"+val, os.O_CREATE|os.O_WRONLY, os.ModePerm)
			if err != nil {
				panic(err)
			}

			bufWriter := bufio.NewWriter(dstFile)

			segs, err := parseDir.GetIPSeg("ips/"+val)
			if err != nil {
				panic(err)
			}
			for _, val := range segs {
				bufWriter.WriteString(val+"\n")
			}

			bufWriter.Flush()
			dstFile.Close()
		}()
	}
}

func getAliveIPSeg() {
	start := time.Now()

	filenames, _ := parseDir.GetFilesName("ips_ret/")

	fmt.Println("discovered files: ")
	for _, val := range filenames {
		fmt.Printf("%s\n", val)
	}

	for _, file := range filenames {
		ips, err := pingScan.IPFileScanAlive(file)
		if err != nil {
			fmt.Println("file scan fail: ", file, err)
			continue
		}

		dstFile, err := os.OpenFile("ips_connected/"+file, os.O_CREATE|os.O_WRONLY, os.ModePerm)
		if err != nil {
			fmt.Println("open file to write fail: ", file)
			continue
		}

		bufWriter := bufio.NewWriter(dstFile)
		for _, val := range ips {
			if val != "" {
				bufWriter.WriteString(val+"\n")
			}
		}
		bufWriter.Flush()
		dstFile.Close()
		fmt.Printf("file scan done: %s\n", file)
	}
	cost := time.Since(start)
	fmt.Printf("run time: [%s]\n", cost)
}

func splitIPAlive() {
	//start := time.Now()

	//var ipSlice []string

	// Find files.
	filenames, _ := parseDir.GetFilesName("ips_connected")
	fmt.Println("discovered files:")
	for _, val := range filenames {
		fmt.Printf("%s\n", val)
	}

	for _,file := range filenames {
		// Open file && dived ips into 100 per slice.
		ipDivided, err := screen_ip.GetIPSplit(file)
		if err != nil {
			fmt.Println("opsn file " + file + " fail!")
			return
		}

		var wg sync.WaitGroup
		var alives []string
		rand.Seed(time.Now().UnixNano())

		// start scan.
		fmt.Println("scan from file: ", file)
		for _, ipSlices := range ipDivided {
			wg.Add(1)
			go func(ipSlicesParm []string) {
				defer wg.Done()
				for {
					if len(ipSlicesParm) <= 1 {
						return
					}
					if len(ipDivided) == 1 {
						scanIP1 := ipSlicesParm[rand.Intn(len(ipSlicesParm)-1)]
						scanIP2 := ipSlices[rand.Intn(len(ipSlicesParm)-1)]
						fmt.Println("start scan ip seg: ", scanIP1)
						fmt.Println("start scan ip seg: ", scanIP2)
						ret1, _ := pingScan.PingScan(scanIP1)
						ret2, _ := pingScan.PingScan(scanIP2)
						alives = append(alives, ret1[0])
						alives = append(alives, ret2[0])
						break
					}

					scanIP := ipSlicesParm[rand.Intn(len(ipSlicesParm)-1)]
					fmt.Println("start scan ip seg: ", scanIP)

					ret, _ := pingScan.PingScan(scanIP)
					if len(ret) <= 1 {
						continue
					}
					alives = append(alives, ret[rand.Intn(len(ret)-1)])
					break
				}
			}(ipSlices)
		}
		wg.Wait()
		fmt.Println("scan file: ", file, " finished")


		dstFile, err := os.OpenFile("ips_alive/"+file, os.O_CREATE|os.O_WRONLY, os.ModePerm)
		if err != nil {
			fmt.Println("open file to write fail: ", file)
			continue
		}

		bufWriter := bufio.NewWriter(dstFile)
		for _, val := range alives {
			if val != "" {
				bufWriter.WriteString(val+"\n")
			}
		}
		bufWriter.Flush()
		dstFile.Close()
	}
}

func main() {
	splitIPAlive()
}
