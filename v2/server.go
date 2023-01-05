package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"
	"strings"

	// "net"
	"os"
)

type Ip struct { //IP结构体

}

type GetIpRequest struct { //得到IP请求结构体

}

type GetIpResponse struct { //得到IP返回结构体
	Hostip []string //返回所有IP地址
}

//上传IP方法UpIp(客户端角度)
func (this *Ip) GetIp(getiprequest GetIpRequest, getipresponse *GetIpResponse) error {
	getipresponse.Hostip = PingList() //返回值从数据库函数getHostip获取
	fmt.Println("getipresponse", getipresponse)
	return nil
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
		hostIp = strings.Replace(hostIp, "\n", "", -1)
		ipArray = append(ipArray, hostIp) //IP存入数组
	}
	fmt.Println(ipArray)
	return ipArray //返回IP数组
}

//添加IP
func InsertHostip(localIp string) {
	file, err := os.OpenFile("/Users/slgx/study/GoLang/Pingmesh/v2/pinglist.txt", os.O_WRONLY|os.O_APPEND, 0666) //文件里追加数据
	CheckError(err)
	defer file.Close()
	write := bufio.NewWriter(file)
	write.WriteString(localIp + "\n")
	write.Flush()
}

//多线程监听
func ListenIp() {
	rpc.Register(new(Ip)) //rpc注册IP方法

	listen, err := net.Listen("tcp", "127.0.0.1:58098") //监听端口
	CheckError(err)

	defer listen.Close()

	for {
		conn, err := listen.Accept()
		if err != nil {
			continue
		}

		go func(conn net.Conn) {
			fmt.Fprintf(os.Stdout, "%s", "new client is comming\n")
			fmt.Println(conn.RemoteAddr())
			jsonrpc.ServeConn(conn)
			ipaddress := conn.RemoteAddr().String()
			ipadd := strings.Split(ipaddress, ":")[0]
			fmt.Printf("%s\n", ipadd)

			InsertHostip(ipadd) //先调用插入ip函数

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
	ListenIp()
}
