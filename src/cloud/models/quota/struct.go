package quota

type CloudQuotaName struct {
	//配额名称
	QuotaName string
	//
	QuotaId int64
}

type QuotaAppUsed struct {
	//配额名称
	ResourceName string
	// 数量
	Cnt int64
}

//2018-02-11 09:31:19.5279645 +0800 CST
type CloudQuota struct {
	//配额描述信息
	Description string
	//cpu配额多少核心
	QuotaCpu int64
	//资源名称
	ResourceName string
	//最近修改用户
	LastModifyUser string
	//限制pod数量
	PodNumber int64
	//限制应用数量
	ServiceNumber int64
	//配额名称
	QuotaName string
	//创建用户
	CreateUser string
	//最近修改时间
	LastModifyTime string
	//
	QuotaId int64
	//创建时间
	CreateTime string
	//受限业务线名称
	GroupName string
	//内存配额多少MB
	QuotaMemory int64
	//
	Status string
	//限制应用数量
	AppNumber int64
	//限制负载均衡数量
	LbNumber int64
	//限制发布任务数量
	JobNumber int64
	//限制流水线数量
	PipelineNumber int64
	//受限人名称
	UserName string
	// 镜像仓库组
	RegistryGroupNumber int64
	// dockerfield数量
	DockerFileNumber int64
}

// 2018-02-11 16:21
// 配额使用情况
type QuotaUsed struct {
	CloudQuota
	LbUsed int64
	CpuUsed int64
	MemoryUsed int64
	PipelineUsed int64
	JobUsed int64
	AppUsed int64
	ServiceUsed int64
	RegistryGroupUsed int64
	PodUsed int64
	LbPercent int64
	CpuPercent int64
	MemoryPercent int64
	PipelinePercent int64
	JobPercent int64
	AppPercent int64
	ServicePercent int64
	RegistryGroupPercent int64
	PodPercent int64
	LbFree int64
	CpuFree int64
	MemoryFree int64
	PipelineFree int64
	DockerFileNumber int64
	DockerFileUsed int64
	DockerFileFree int64
	DockerFilePercent int64
	JobFree int64
	AppFree int64
	ServiceFree int64
	PodFree int64
	RegistryGroupFree int64
}

// 2018-02-11 10:02
// 获取配置默认数据
func GetDefaultQuota() CloudQuota {
	cloudQuota := CloudQuota{
		AppNumber:      1,
		JobNumber:      1,
		PipelineNumber: 1,
		LbNumber:       1,
		ServiceNumber:  1,
		PodNumber:      1,
		QuotaCpu:       1,
		RegistryGroupNumber:1,
		QuotaMemory:    512,
		DockerFileNumber: 1,
	}
	return cloudQuota
}
