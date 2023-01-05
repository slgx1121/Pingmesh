package main

import (
	"fmt"
	"log"
	"net"
	"net/rpc/jsonrpc"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

type GetIpRequest struct { //得到IP请求结构体

}

type GetIpResponse struct { //得到IP返回结构体
	Hostip []string //返回所有IP地址
}

type Pingstruct struct {
	Tss                           int64
	Src, Dst, Loss, Min, Avg, Max string
}

type UpIpArrayRequet struct { //返回ping结果组结构体
	UpIparrayrequet []Pingstruct
}

type UpIpResponse struct { //得到IP返回结构体

}

var pingStructArray []Pingstruct

func GetLocalIp() string { //获取本地ip
	var IP string
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				IP = ipnet.IP.String()
			}
		}
	}
	return IP
}

func fPing(ipadd []string) { //获取目标ip,丢包率，ping平均延迟
	runCommand("fping", ipadd...)
	// fmt.Println("ipadd:", ipadd)
}

func runCommand(name string, arg ...string) { //运行fping命令
	var pingstruct Pingstruct
	var Pingstructlist []Pingstruct
	arg = append(arg, "-q")
	arg = append(arg, "-p")
	arg = append(arg, "1000")
	arg = append(arg, "-c")
	arg = append(arg, "2")
	cmd := exec.Command(name, arg...)
	res, _ := cmd.CombinedOutput()

	re := regexp.MustCompile(`(.*) +: xmt/rcv/%loss = (.*), min/avg/max = (.*)`) //正则过滤下丢包率为100%的
	submatchall := re.FindAllStringSubmatch(string(res), -1)
	// fmt.Println(submatchall)
	for _, element := range submatchall {
		pingstruct.Src = GetLocalIp()
		pingstruct.Tss = time.Now().Unix()
		pingstruct.Dst = element[1]
		pingstruct.Loss = strings.Split(element[2], "/")[2]
		pingstruct.Min = strings.Split(element[3], "/")[0]
		pingstruct.Avg = strings.Split(element[3], "/")[1]
		pingstruct.Max = strings.Split(element[3], "/")[2]
		Pingstructlist = append(Pingstructlist, pingstruct)
	}
	pingStructArray = Pingstructlist
	fmt.Println(pingStructArray)

}

func pingHost() []string { //得到所有host组的ip
	conn, err := jsonrpc.Dial("tcp", "127.0.0.1:58098") //服务器的ip
	if err != nil {
		log.Fatalln("dailing error:", err)
	}

	getiprequest := GetIpRequest{}
	var getipresponse GetIpResponse
	err = conn.Call("Ip.GetIp", getiprequest, &getipresponse)
	if err != nil {
		log.Fatalln("getip error: ", err)
	}
	// fmt.Println(getipresponse)
	conn.Close()

	// fmt.Println(getipresponse.Hostip)
	return getipresponse.Hostip
}

func UpIp() { //上传fping的结果
	conn, err := jsonrpc.Dial("tcp", "127.0.0.1:58099")
	if err != nil {
		log.Fatalln("dailing error:", err)
	}

	upip := UpIpArrayRequet{pingStructArray}
	var rippr UpIpResponse
	err = conn.Call("Ip.UpIp", upip, &rippr)
	if err != nil {
		log.Fatalln("return error: ", err)
	}
	conn.Close()
}

func main() {
	//	t_start := time.Now()

	// ticker := time.NewTicker(time.Second * 60)			//每一分钟执行执行一次
	// for {
	// 	select {
	// 	case  <-ticker.C:
	fPing(pingHost()) //得到所有主机ip
	// UpIp()            //提交fping的结果
	// 	}
	// }

}
