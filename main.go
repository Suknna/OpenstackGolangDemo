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
		//AuthEndpoint: "http://192.168.200.10:5000/v3",
		AuthEndpoint: "http://132.91.181.250:30500/v3",
		UserName:     "admin",
		//Password:     "297b0c495cf84464",
		Password:   "e29995bc97aaTi0*",
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
