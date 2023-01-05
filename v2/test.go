package main

import (
	"fmt"
	"os/exec"
	"regexp"
)

func main() {
	var args []string
	args = append(args, "114.114.114.114")
	args = append(args, "-p")
	args = append(args, "1000")
	args = append(args, "-c")
	args = append(args, "2")
	cmd := exec.Command("fping", args...)
	fmt.Println(cmd)
	res, _ := cmd.CombinedOutput()
	fmt.Println(string(res))
	re := regexp.MustCompile(`(.*) +: xmt/rcv/%loss = (.*), min/avg/max = (.*)`) //正则过滤下丢包率为100%的
	submatchall := re.FindAllStringSubmatch(string(res), -1)
	fmt.Println(submatchall)
}
