package app

// 2018-01-22 14:41
// 获取应用名称
type CloudAppName struct {
	//应用名称
	AppName string
	// 创建用户
	CreateUser string
	// 集群名称
	ClusterName string
	// 环境名称
	Entname string
	//
	AppId string
}

type CloudApp struct {
	//
	AppId int64
	//应用名称
	AppName string
	//运行状态
	Status string
	//最近修改时间
	LastModifyTime string
	//应用标签
	AppLabels string
	//其他非固定数据存储
	JsonData string
	//应用类型
	AppType string
	//资源空间
	ResourceName string
	//创建时间
	CreateTime string
	//创建用户
	CreateUser string
	//最近修改用户
	LastModifyUser string
	//最近更新时间
	LastUpdateTime string
	// 编排文件
	Yaml string
	// 集群名称
	ClusterName string
	// 是否创建service
	IsService string
	// uuid
	Uuid string
	// 服务的yaml文件
	ServiceYaml string
	//存放服务的端口,json格式 {name:xxxx,port:8080}
	NodePort string
	// 环境名称
	Entname string
}

type CloudAppTemplateName struct {
	//模板名称
	TemplateName string
	//
	TemplateId int64
}

type CloudAppServiceName struct {
	//service名称
	ServiceName string
	//
	ServiceId int64
	// app
	AppName string
	// 环境名称
	Entname string
	// 集群名称
	ClusterName string
	// 资源名称
	ResourceName string
	// 创建用户
	CreateUser string

}

type CloudAppTemplate struct {
	//模板名称
	TemplateName string
	//创建时间
	CreateTime string
	//创建用户
	CreateUser string
	//最近修改时间
	LastModifyTime string
	//最近修改用户
	LastModifyUser string
	//资源名称
	ResourceName string
	//描述信息
	Description string
	//yaml编排文件
	Yaml string
	//
	TemplateId int64
	// 集群
	Cluster string
	// 环境
	Ent string
	// 服务名称
	ServiceName string
	// 应用名称
	AppName string
	// 域名
	Domain string
}





type CloudAppConfigureName struct {
	//
	ConfigureId int64
	//模板名称
	ConfigureName string
}

type CloudAppConfigure struct {
    //创建用户
    CreateUser string
    //最近修改时间
    LastModifyTime string
    //集群名称
    ClusterName string
    //描述信息
    Description string
    //
    ConfigureId int64
    //模板名称
    ConfigureName string
    //创建时间
    CreateTime string
    //最近修改用户
    LastModifyUser string
	// 环境名称
	Entname string
}

type ConfigDataName struct {
	// 配置文件key名称
	DataName string
	//
	DataId int64
}


type CloudConfigData struct {
    //创建用户
    CreateUser string
    //最近修改时间
    LastModifyTime string
    //参考config的id
    ConfigureId int64
    //配置名称,参考配置名称
    ConfigureName string
    //配置文件数据
    Data string
    //
    DataId int64
    //创建时间
    CreateTime string
    //最近修改用户
    LastModifyUser string
	// 配置文件key名称
	DataName string
}

// 响应命名空间数据
type CloudAppServiceInfo struct {
	Namespace string
	ServiceName string
}
//2018-01-11 11:40:10.8610181 +0800 CST
type CloudAppService struct {
    //部署模式, deployment daemonset statefulset
    DeployType string
    //容器副本数量
    Replicas int64
    //service名称
    ServiceName string
    //最近修改时间
    LastModifyTime string
    //内存数
    Memory int64
    //服务器标签数据
    ServiceLablesData string
    //集群名称
    ClusterName string
    //网络模式 flannel host
    NetworkMode string
    //手动配置文件的内容
    Config string
    //容器最小数量
    ReplicasMin int64
    //健康检查数据
    HealthData string
    //
    ServiceId int64
    //运行状态
    Status string
    //镜像版本号
    ImageTag string
    //容器最多数量
    ReplicasMax int64
    //创建用户
    CreateUser string
    //服务标签,用map标识
    ServiceLabels string
    //负载均衡数据
    LbData string
    //service类型,有状态和无状态
    ServiceType string
    //其他非固定数据存储
    JsonData string
    //容器端口,多个逗号分隔
    ContainerPort string
    //参考环境配置的信息
    EnvFile string
    //存储配置数据
    StorageData string
    //资源空间
    ResourceName string
    //cpu核数
    Cpu float32
    //负载均衡名称
    LbName string
    //
    AppLabels string
    //编排文件
    Yaml string
    //挂载配置文件数据
    ConfigureData string
    //创建时间
    CreateTime string
    //最近修改用户
    LastModifyUser string
    //最近更新时间
    LastUpdateTime string
	// 环境变量值
	Envs string
	// 应用名称
	AppName string
	// 滚动升级时候,会优先启动的pod数量
	MaxSurge int
	// 滚动升级时候,最大的unavailable数量
	MaxUnavailable int
	// 指定没有任何容器crash的Pod并被认为是可用状态的最小秒数
	MinReady int
	// 镜像仓库地址
	ImageRegistry string
	// 环境名称
	Entname string
	// 服务版本做蓝绿,灰度部署标签,
	// 有1和2,如果1存在那么就部署一个2,
	// 如果2存在就部署一个1,
	// 当确认发布完成 ,删除一个未使用的部署
	ServiceVersion string
	// 域名
	Domain string
	// pod关闭时间
	TerminationSeconds int
	// 日志路径
	LogPath string
}

// 存储容器名称数据,在更新数据时做判断使用,不用频繁查库
type CloudContainerName struct {
	//容器名称
	ContainerName string
	//应用名称
	AppName string
	//服务名称
	ServiceName string
	//集群名称
	ClusterName string
}


//2018-01-15 15:20:35.5687711 +0800 CST
type CloudContainer struct {
    //
    ContainerId int64
    //容器名称
    ContainerName string
    //宿主机地址
    ServerAddress string
    //容器ip
    ContainerIp string
    //应用名称
    AppName string
    //创建时间
    CreateTime string
    //服务名称
    ServiceName string
    //集群名称
    ClusterName string
    //镜像名称
    Image string
    //运行状态
    Status string
	// 资源名称
	ResourceName string
	// cpu
	Cpu int64
	// 内存
	Memory int64
	// Env
	Env string
	// 运行的程序
	Process string
	// 存储数据
	StorageData string
	// 容器等待原因信息
	WaitingMessages string
	// 容器等待原因
	WaitingReason string
	// 容器停止原因信息
	TerminatedMessages string
	// 容器停止原因
	TerminatedReason string
	// 创建用户
	CreateUser string
	// 环境名称
	Entname string
	// 事件信息
	Events string
	// 重启次数
	Restart int32
	// 服务信息
	Service string
	// 数据更新时间
	LastUpdateTime int64
}
//2018-08-16 16:04:25.8692888 +0800 CST
type CloudTemplateDeployHistory struct {
    //创建用户
    CreateUser string
    //环境名称
    Entname string
    //创建时间
    CreateTime string
    //service名称
    ServiceName string
    //应用名称
    AppName string
    //环境名称
    ResourceName string
    //集群名称
    ClusterName string
    //
    TemplateName string
    //
    HistoryId int64
    // 域名
    Domain string
}
