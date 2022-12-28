package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

type Pingstruct struct {
	//源IP地址，目的IP地址，丢失率，平均延迟
	Src, Dst, Loss, Avg string
}

//获取本地IP
func Getip() string {
	var Localip string
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				Localip = ipnet.IP.String()
			}
		}
	}
	return Localip
}

func main() {
	var Pingstructlist []Pingstruct
	var pingstruct Pingstruct
	t_start := time.Now()
	// fmt.Println(Getip())
	cmd := exec.Command("fping", "-q", "-p100", "-c5", "10.176.40.181:80", "10.176.40.182")
	res, _ := cmd.CombinedOutput()
	lists := strings.Split(string(res), "\n")
	for _, list := range lists {
		fmt.Println(list)
	}
	re := regexp.MustCompile(`(.*) +: xmt/rcv/%loss = (.*), min/avg/max = (.*)`)
	submatchall := re.FindAllStringSubmatch(string(res), -1)
	for _, element := range submatchall {
		pingstruct.Src = Getip()
		pingstruct.Dst = element[1]
		pingstruct.Loss = strings.Split(element[2], "/")[2]
		pingstruct.Avg = strings.Split(element[3], "/")[2]
		Pingstructlist = append(Pingstructlist, pingstruct)
	}
	//将结果生成json格式
	data, _ := json.Marshal(Pingstructlist)
	fmt.Println(string(data))

	t_end := time.Now()
	fmt.Println(t_end.Sub(t_start))
}
