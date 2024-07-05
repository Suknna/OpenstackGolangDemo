package node

import (
	"encoding/json"
	"fmt"
	"log"
	outputfile "openstackClient/mode/outputFile"
	"strconv"
	"strings"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/extensions/hypervisors"
)

// 传入nova client查看hypervisor 列表
func getHypervisorsList(serviceClient *gophercloud.ServiceClient, fileName string) error {
	// 存储集群全部hypervisor的信息变量
	var datas string
	page, err := hypervisors.List(serviceClient, hypervisors.ListOpts{}).AllPages()
	if err != nil {
		return fmt.Errorf("get hypervisor list failed %s", err.Error())
	}
	// 获取body
	bodyData := page.GetBody().(map[string]interface{})["hypervisors"]
	// 将body转换为json
	jsB, err := json.Marshal(bodyData)
	if err != nil {
		return fmt.Errorf("json marshal failed, %s", err.Error())
	}
	var data []hypervisorData
	err = json.Unmarshal(jsB, &data)
	if err != nil {
		return fmt.Errorf("json unmarshal failed %s", err.Error())
	}
	header := "clusterId,serverIP,serverName,vcpus,vcpusUsed,memoryGB,memoryUsedGB,runningVms,status,state,uptime"
	datas += fmt.Sprintln(header)
	// 定义一个temp变量，用于临时存放hypervisorData的值，用于下面判断是否重复添加
	seen := make(map[hypervisorData]bool)
	for i := range data {
		// 判断data是否被重复添加，如果被重复添加则结束此次循环
		if !seen[data[i]] {
			r := hypervisors.GetUptime(serviceClient, strconv.Itoa(data[i].ID))
			// 处理字符串，原始内容为：15:56:17 up  7:09,  0 users,  load average: 0.58, 0.54, 0.60
			upTimeStr := strings.Split(r.Body.(map[string]interface{})["hypervisor"].(map[string]interface{})["uptime"].(string), ",")[0]
			upTimeStr = strings.TrimSpace(strings.Split(upTimeStr, "up")[1])
			seen[data[i]] = true
			data := fmt.Sprintf("%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v\n", data[i].ID, data[i].HostIP, data[i].HypervisorHostname, data[i].VCPUs, data[i].VCPUsUsed, data[i].MemoryMB, data[i].MemoryMBUsed, data[i].RunningVMs, data[i].Status, data[i].State, upTimeStr)
			// 向存储集群全部hypervisor的信息变量中写入单一hypervisor的信息
			datas += data
		} else {
			continue
		}
	}
	// 写入文件
	err = outputfile.OutputFile(datas, fileName)
	if err != nil {
		return fmt.Errorf("out put file failed, err: %s", err.Error())
	}
	return nil
}

// 传入认证后的client，获取hypervisor的列表
func RunHypervisorList(client *gophercloud.ProviderClient, fileName string) {
	// 创建service client
	computeClient, err := openstack.NewComputeV2(client, gophercloud.EndpointOpts{
		Region: "RegionOne",
		Type:   "compute",
	})
	if err != nil {
		log.Fatalf("Have some error. %s\n", err.Error())
	}
	// 获取hypervisor信息
	err = getHypervisorsList(computeClient, fileName)
	if err != nil {
		log.Fatalf("Have some error. %s\n", err.Error())
	}
}
