package permdata


//2018-01-15 15:20:35.5687711 +0800 CST
type CloudContainer struct {
	//容器名称
	ContainerName string
	//应用名称
	AppName string
	//服务名称
	ServiceName string
	//集群名称
	ClusterName string
	// 资源名称
	ResourceName string
	// 创建用户
	CreateUser string
	// 环境名称
	Entname string
}

//2018-11-28 10:00:06.9680777 +0800 CST
type CloudUserPermDetail struct {
	//
	DetailId int64
	//
	Username string
	//
	ResourceType string
	//
	GroupName string
	//
	Name string
	//
	ClusterName string
	//
	Ent string
}