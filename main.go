package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
	"sync"
)

type namegroup struct {
	name string
	ips  []string
}

var ipss []namegroup

func dns(domain string, group *sync.WaitGroup) {

	ips, err := net.LookupHost(domain)
	if err != nil {
		fmt.Printf("failed to lookup %s: %v\n", domain, err)
	}
	if len(ips) > 0 {
		ng := &namegroup{
			name: domain,
			ips:  []string{},
		}
		for _, ip := range ips {
			fmt.Printf("%s resolves to %s\n", domain, ip)
			ng.ips = append(ng.ips, ip)
		}
		ipss = append(ipss, *ng)
	}

	group.Done()
}

func main() {

	var inputfile string
	var outputfile string
	flag.StringVar(&inputfile, "f", "name.txt", "输入文件名")
	flag.StringVar(&outputfile, "o", "out.txt", "输出文件名")
	flag.Parse()
	//无参数时打印帮助信息
	if len(os.Args) < 2 {
		flag.Usage()
		os.Exit(0)
	}

	var wg sync.WaitGroup
	file, err := os.Open(inputfile)
	if err != nil {
		fmt.Println("打开文件时出错:", err)
		return
	}
	defer file.Close()
	rd := bufio.NewScanner(file)
	for rd.Scan() {
		//println()
		wg.Add(1)
		go dns(rd.Text(), &wg)
	}
	wg.Wait()
	if ipss != nil {

		outfile, e := os.OpenFile(outputfile, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
		if e != nil {
			println("无法创建输出文件")
			println(e.Error())
			os.Exit(0)
		}
		for _, s := range ipss {
			for _, ip := range s.ips {
				ss := fmt.Sprintf("%s %s\n", s.name, ip)
				outfile.WriteString(ss)
			}

		}
		outfile.Close()

	}
}
