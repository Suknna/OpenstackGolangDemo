package node

// 每台节点的信息
type hypervisorData struct {
	CPUInfo            string `json:"cpu_info"`
	CurrentWorkload    int    `json:"current_workload"`
	DiskAvailableLeast int    `json:"disk_available_least"`
	FreeDiskGB         int    `json:"free_disk_gb"`
	FreeRamMB          int    `json:"free_ram_mb"`
	HostIP             string `json:"host_ip"`
	HypervisorHostname string `json:"hypervisor_hostname"`
	HypervisorType     string `json:"hypervisor_type"`
	HypervisorVersion  int    `json:"hypervisor_version"`
	ID                 int    `json:"id"`
	LocalGB            int    `json:"local_gb"`
	LocalGBUsed        int    `json:"local_gb_used"`
	MemoryMB           int    `json:"memory_mb"`
	MemoryMBUsed       int    `json:"memory_mb_used"`
	RunningVMs         int    `json:"running_vms"`
	State              string `json:"state"`
	Status             string `json:"status"`
	VCPUs              int    `json:"vcpus"`
	VCPUsUsed          int    `json:"vcpus_used"`
}
