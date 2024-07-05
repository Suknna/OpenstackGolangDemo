package main

import (
	"log"
	"openstackClient/authkeystone"
	"openstackClient/mode/node"
	"openstackClient/mode/vm"
	"os"
	"path/filepath"
)

const (
	hypervisorFileName = "hypervisorList.txt"
	serverFileName     = "serverInfo.txt"
)

func main() {
	// 获取程序工作目录
	path, err := os.Executable()
	if err != nil {
		log.Fatalf("Failed to obtain program working directory, err: %s\n", err.Error())
	}
	path = filepath.Dir(path)
	// 认证
	auth := authkeystone.AuthSetUp{
		AuthEndpoint: "http://192.168.200.10:5000/v3",
		UserName:     "admin",
		Password:     "000000",
		DomainName: "Default",
		TenantName: "admin",
	}
	pc, err := auth.Auth()
	if err != nil {
		log.Fatalf("Have some error.  %s\n", err.Error())
	}
	vm.RunServerList(pc, filepath.Join(path, serverFileName))
	node.RunHypervisorList(pc, filepath.Join(path, hypervisorFileName))
}
