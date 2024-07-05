package vm

// 实例对象
type Instance struct {
	OSEXTAZAvailabilityZone        string `json:"OS-EXT-AZ:availability_zone"`
	OSEXTSRVATTRHost               string `json:"OS-EXT-SRV-ATTR:host"`
	OSEXTSRVATTRHypervisorHostname string `json:"OS-EXT-SRV-ATTR:hypervisor_hostname"`
	OSEXTSRVATTRInstanceName       string `json:"OS-EXT-SRV-ATTR:instance_name"`
	OSEXTSTSPowerState             int    `json:"OS-EXT-STS:power_state"`
	OSEXTSTSVMState                string `json:"OS-EXT-STS:vm_state"`
	OSSRVUSGLaunchedAt             string `json:"OS-SRV-USG:launched_at"`
	HostId                         string `json:"hostId"`
	ID                             string `json:"id"`
	Name                           string `json:"name"`
	Status                         string `json:"status"`
	Updated                        string `json:"updated"`
}

// 云主机类型对象
type Flavor struct {
	ID    string `json:"id"`
	Disk  int    `json:"disk"`
	Name  string `json:"name"`
	RAM   int    `json:"ram"`
	VCPUs int    `json:"vcpus"`
}

//挂载点信息
type Attachment struct {
	Device   string `json:"device"`
	HostName string `json:"host_name"`
	ID       string `json:"id"`
	ServerID string `json:"server_id"`
	VolumeID string `json:"volume_id"`
}

// volume拷贝的镜像信息
type VolumeImageMetadata struct {
	ImageID   string `json:"image_id"`
	ImageName string `json:"image_name"`
	Size      string `json:"size"`
}

// 块磁盘信息
type Volume struct {
	Attachments         []Attachment        `json:"attachments"`
	AvailabilityZone    string              `json:"availability_zone"`
	ID                  string              `json:"id"`
	Size                int                 `json:"size"`
	Status              string              `json:"status"`
	VolumeImageMetadata VolumeImageMetadata `json:"volume_image_metadata"`
	VolumeType          string              `json:"volume_type"`
}

// 镜像对象
type Image struct {
	File   string `json:"file"`
	ID     string `json:"id"`
	Name   string `json:"name"`
	Schema string `json:"schema"`
	Self   string `json:"self"`
	// 单位M 值除以1024
	Size   int    `json:"size"`
	Status string `json:"status"`
}

// IP对象

type IP struct {
	OSEXTIPSMAC string `json:"OS-EXT-IPS-MAC:mac_addr"`
	OSEXTIPS    string `json:"OS-EXT-IPS:type"`
	Addr        string `json:"addr"`
	Version     int    `json:"version"`
}

type IPS struct {
	fixed    string
	floating string
}
