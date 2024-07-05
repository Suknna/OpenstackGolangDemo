package vm

import (
	"encoding/json"
	"fmt"
	"log"
	outputfile "openstackClient/mode/outputFile"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/blockstorage/v3/volumes"
	"github.com/gophercloud/gophercloud/openstack/compute/v2/servers"
	"github.com/gophercloud/gophercloud/openstack/db/v1/flavors"
	"github.com/gophercloud/gophercloud/openstack/imageservice/v2/images"
)

/*
获取云服务器的全部信息

*/

// 获取云服务器的状态
func getServer(computeClient *gophercloud.ServiceClient, cinderClient *gophercloud.ServiceClient, glanceClient *gophercloud.ServiceClient, fileName string) error {
	// 存储获取的全部云主机信息
	var datas string
	// 获取云主机的大致信息
	page, err := servers.List(computeClient, servers.ListOpts{
		AllTenants: true,
	}).AllPages()
	if err != nil {
		return fmt.Errorf("get server list failed %s", err.Error())
	}
	allServers, err := servers.ExtractServers(page)
	if err != nil {
		return fmt.Errorf("servers Parsing failed %s", err.Error())
	}
	// 获取镜像列表
	imgs, err := imageList(glanceClient)
	// fmt.Println(imgs)
	if err != nil {
		return fmt.Errorf("run imageList func failed %s", err.Error())
	}
	// 写入表头
	header := "uuid,vmName,instanceName,availabilityZone,hostName,numDisk,imageName,vcps,memory,privateNet,publicNet,status,runTime"
	datas += fmt.Sprintln(header)
	// 获取单个云主机的信息
	for _, s := range allServers {
		// 获取实例的详细信息
		instance, err := vmInfo(computeClient, s.ID)
		if err != nil {
			return fmt.Errorf("get vm info failed %s", err.Error())
		}
		// 获取实例类型的详细信息
		flavor, err := flavorInfo(computeClient, s.Flavor["id"].(string))
		if err != nil {
			return fmt.Errorf("get flavor info failed %s", err.Error())
		}
		// 获取实例所使用的块存储信息
		v, err := volumeInfo(cinderClient, s.AttachedVolumes)
		if err != nil {
			return fmt.Errorf("get volume info failed %s", err.Error())
		}
		// 获取网络
		i, err := networkInfo(s.Addresses)
		if err != nil {
			return fmt.Errorf("get network info failed %s", err.Error())
		}
		data := formatData(*instance, *flavor, v, i, imgs, s.Image)
		datas += data
	}
	err = outputfile.OutputFile(datas, fileName)
	if err != nil {
		return fmt.Errorf("out put file failed %s", err.Error())
	}
	return nil
}

// 获取网络信息
func networkInfo(netdata interface{}) (*IPS, error) {
	// 由于键的名称不固定，这里采用反射来取键名
	// 通过reflect包来进行反射
	var i IPS
	ifaceValue := reflect.ValueOf(netdata)
	if ifaceValue.Kind() != reflect.Map {
		return nil, fmt.Errorf("networkInfo: get key name failed")
	}
	for _, key := range ifaceValue.MapKeys() {
		keyStr := key.String()
		ips := netdata.(map[string]interface{})[keyStr].([]interface{})
		// 循环读取ips中的内容
		for _, v := range ips {
			neType := v.(map[string]interface{})["OS-EXT-IPS:type"].(string)
			if neType == "fixed" {
				i.fixed = v.(map[string]interface{})["addr"].(string)
				// addrs = v.(map[string]interface{})["addr"].(string)
			} else if neType == "floating" {
				i.floating = v.(map[string]interface{})["addr"].(string)
				//addrs = addrs + "," + v.(map[string]interface{})["addr"].(string)
			}
		}
	}
	return &i, nil
}

// 获取云服务器的详细信息
func vmInfo(serviceClient *gophercloud.ServiceClient, id string) (*Instance, error) {
	// 创建一个云服务器对象
	var vm Instance
	// 获取云主机的详细信息
	vmBody := servers.Get(serviceClient, id).Body.(map[string]interface{})["server"]
	// 转换为json
	jsonData, err := json.Marshal(vmBody)
	if err != nil {
		return nil, fmt.Errorf("server get response body marshall json failed %s", err.Error())
	}
	//json转换为struct
	err = json.Unmarshal(jsonData, &vm)
	if err != nil {
		return nil, fmt.Errorf("server get response body unmarshal struct failed %s", err.Error())
	}
	return &vm, nil
}

// 获取云服务器挂载的块存储信息
func volumeInfo(serviceClient *gophercloud.ServiceClient, volumeidSlice []servers.AttachedVolume) (volume []Volume, err error) {
	for _, id := range volumeidSlice {
		var v Volume
		volumeBody := volumes.Get(serviceClient, id.ID).Body.(map[string]interface{})["volume"]
		// json编码
		jsonData, err := json.Marshal(volumeBody)
		if err != nil {
			return nil, fmt.Errorf("volume get response body marshall json failed %s", err.Error())
		}
		// json解码
		err = json.Unmarshal(jsonData, &v)
		if err != nil {
			return nil, fmt.Errorf("volume get response body unmarshal struct failed  %s", err.Error())
		}
		volume = append(volume, v)
	}
	return volume, nil
}

// 获取云服务器的主机配置
func flavorInfo(serviceClient *gophercloud.ServiceClient, id string) (*Flavor, error) {
	var f Flavor
	flavorBody := flavors.Get(serviceClient, id).Body.(map[string]interface{})["flavor"]
	// 将interface转换为json
	jsonData, err := json.Marshal(flavorBody)
	if err != nil {
		return nil, fmt.Errorf("flavor get response body marshall json failed %s", err.Error())
	}
	// 将json转换为struct
	err = json.Unmarshal(jsonData, &f)
	if err != nil {
		return nil, fmt.Errorf("flavor json unmarshal struct failed %s", err.Error())
	}
	return &f, nil
}

// 获取镜像列表
func imageList(serviceClient *gophercloud.ServiceClient) (map[string]string, error) {
	imgMap := make(map[string]string)
	pagers, err := images.List(serviceClient, images.ListOpts{}).AllPages()
	if err != nil {
		return nil, fmt.Errorf("get image list failed %s", err.Error())
	}
	imgs, err := images.ExtractImages(pagers)
	if err != nil {
		return nil, fmt.Errorf("get image data parse failed %s", err.Error())
	}
	for _, v := range imgs {
		_, ok := imgMap[v.ID]
		if !ok {
			imgMap[v.ID] = v.Name
		} else {
			continue
		}
	}
	return imgMap, nil
}

// 根据imageList函数获取的镜像列表匹配云主机使用的镜像名称
func imageGetName(imgList map[string]string, serverImgInfo interface{}) string {
	id, ok := serverImgInfo.(map[string]interface{})["ID"].(string)
	if !ok {
		return ""
	}
	for k, v := range imgList {
		ok := strings.Contains(id, k)
		if ok {
			return v
		}
	}
	return ""
}

// 格式化输出信息
func formatData(ins Instance, fl Flavor, vo []Volume, addrs *IPS, imgList map[string]string, serverImgInfo interface{}) string {
	var (
		id, name, insName, azName, hostName, numDisk, imageName, vcpus, memory, privateNet, publicNet, stat, runTime string
		diskTotal                                                                                                    int
	)
	id = ins.ID
	name = ins.Name
	insName = ins.OSEXTSRVATTRInstanceName
	azName = ins.OSEXTAZAvailabilityZone
	hostName = ins.OSEXTSRVATTRHypervisorHostname
	numDisk = strconv.Itoa(len(vo))
	vcpus = strconv.Itoa(fl.VCPUs)
	memory = strconv.Itoa(fl.RAM)
	stat = ins.Status
	privateNet = addrs.fixed
	publicNet = addrs.floating
	// 获取镜像名称
	imageName = imageGetName(imgList, serverImgInfo)
	// 计算云主机运行时间
	ATTime, err := time.Parse("2006-01-02T15:04:05", ins.OSSRVUSGLaunchedAt)
	if err != nil {
		log.Fatalf("OS-SRV-USG:launched_at time parse failed %s\n", err.Error())
		runTime = "0"
	} else {
		now := time.Now()
		runTime = fmt.Sprintf("%0.2f", now.Sub(ATTime).Hours()/24)
	}

	// 获取磁盘相关信息
	for _, v := range vo {
		diskTotal = diskTotal + v.Size
		// 如果镜像名为空使用sda克隆的磁盘名称
		if imageName == "" {
			imageName = v.VolumeImageMetadata.ImageName
		}
	}
	// "uuid,vmName,instanceName,availabilityZone,hostNmae,numDisk,total,Attachment,imageName,vcps,memory,privateNet,publicNet,status,runTime"
	data := fmt.Sprintf("%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v,%v\n", id, name, insName, azName, hostName, numDisk, imageName, vcpus, memory, privateNet, publicNet, stat, runTime)
	return data
}

// 传入认证后的client，创建不同部分的client
func RunServerList(client *gophercloud.ProviderClient, fileName string) {
	computeClient, err := openstack.NewComputeV2(client, gophercloud.EndpointOpts{
		Region: "RegionOne",
		Type:   "compute",
	})
	if err != nil {
		log.Fatalf("create compute client failed, err: %s\n", err.Error())
	}
	cinderClient, err := openstack.NewBlockStorageV3(client, gophercloud.EndpointOpts{
		Region: "RegionOne",
		Type:   "volumev3",
	})
	if err != nil {
		log.Fatalf("create block storage client failed, err: %s\n", err.Error())
	}
	glanceClient, err := openstack.NewImageServiceV2(client, gophercloud.EndpointOpts{
		Region: "RegionOne",
		Type:   "image",
	})
	if err != nil {
		log.Fatalf("create glance client failed, err: %s", err.Error())
	}
	err = getServer(computeClient, cinderClient, glanceClient, fileName)
	if err != nil {
		log.Fatalf("Have some error. %s\n", err.Error())
	}

}
