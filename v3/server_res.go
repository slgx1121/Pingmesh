package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"

	// "net"
	"os"
)

type Ip struct { //IP结构体

}

type UpIpRequest struct { //返回ping结果结构体
	Tss  int64
	Src  string
	Dst  string
	Loss string
	Avg  string
	Min  string
	Max  string
}

type UpIpArrayRequet struct { //返回ping结果组结构体
	UpIparrayrequet []UpIpRequest
}

type UpIpResponse struct { //得到IP返回结构体

}

var Upiparrayrequet []UpIpRequest

//获取ping结果
func (this *Ip) UpIp(upiparrayrequet UpIpArrayRequet, upiprespone *UpIpResponse) error {
	Upiparrayrequet = upiparrayrequet.UpIparrayrequet

	return nil

}

func insertResult(res []UpIpRequest) { //把客户端的结果存入数据库
	file, err := os.OpenFile("/Users/slgx/study/GoLang/Pingmesh/v3/result.txt", os.O_WRONLY|os.O_APPEND, 0666) //文件里追加数据
	CheckError(err)
	defer file.Close()
	content, err := json.Marshal(res)
	write := bufio.NewWriter(file)
	write.WriteString(string(content) + "\n")
	write.Flush()
}

func PingList() []string { //返回ip地址组
	file, err := os.Open("/Users/slgx/study/GoLang/Pingmesh/v2/pinglist.txt")
	CheckError(err)
	defer file.Close()
	var ipArray []string            //保存ip数组
	reader := bufio.NewReader(file) //读取文件内容
	for {
		hostIp, err := reader.ReadString('\n') //获取IP  #数据得有回车
		if err == io.EOF {
			break
		}
		fmt.Println(hostIp)
		ipArray = append(ipArray, hostIp) //IP存入数组
	}
	return ipArray //返回IP数组
}

func listenResult() { //多线程监听
	rpc.Register(new(Ip))                            //rpc注册IP方法
	lis, err := net.Listen("tcp", "127.0.0.1:58099") //监听端口
	CheckError(err)
	defer lis.Close()
	fmt.Fprint(os.Stdout, "%s", "start connection aprilmadaha 58099")
	for {
		conn, err := lis.Accept()
		if err != nil {
			continue
		}

		go func(conn net.Conn) {
			jsonrpc.ServeConn(conn)
			insertResult(Upiparrayrequet)
		}(conn)
	}
}

func CheckError(err error) { //错误函数
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}
}

func main() {
	// InsertHostip("8.8.8.8")
	// PingList()
	listenResult()
}
