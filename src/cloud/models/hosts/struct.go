package hosts

type CloudClusterHosts struct {
	//主机IP
	HostIp string
	//主机标签
	HostLabel string
	//状态
	Status string
	//容器数量
	ContainerNum int64
	//cpu剩余量
	CpuFree string
	//pod数量
	PodNum int
	//内存大小
	MemSize string
	//创建用户
	CreateUser string
	//最近修改时间
	LastModifyUser string
	//内存剩余量
	MemFree string
	//主机类型
	HostType string
	//cpu使用百分比
	CpuPercent string
	//是否有效
	IsValid int64
	//内使用百分比
	MemPercent string
	//
	HostId int64
	//创建方法
	CreateMethod string
	//cpu数量
	CpuNum int64
	//创建时间
	CreateTime string
	//最近修改时间
	LastModifyTime string
	// 集群名称
	ClusterName string
	// k8sAPi端口,只需要master有就行了
	ApiPort string
	// 镜像数量
	ImageNum int
}

type CloudClusterHostsDetail struct {
	CloudClusterHosts
	K8sVersion string
}
