package k8s

import (
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"cloud/util"
	"cloud/models/hosts"
)

// 配置服务存储数据
type StorageData struct {
	// 宿主机路径
	HostPath string
	// 容器挂载路径
	ContainerPath string
	// 分布式卷名称
	Volume string
	// 读写权限
	Model int
	// 只读挂载
	ReadOnly bool
}

// 亲和性配置结构
// 必须在服务器进行标签配置
// 2018-01-11
type Affinity struct {
	// 类型 disk service server
	Type string
	// 值
	Value string
}

// 选择节点配置
// 2018-01-11 15;57
// 只能选择一个机器
// kubernetes.io/hostname=
type NodeSelector struct {
	// 标签名称,主要是主机名称
	Lables string
	// 标签值
	Value string
}

// 获取服务和pod的映射端口数据
// 2018-01-12 15:36
type PortData struct {
	// 标签名称,主要是主机名称
	NodePort string
	// 容器端口
	ContainerPort string
	// 集群内端口
	Port string
	// 访问类型,集群内, 集群外,集群内外
	Type string
}


// 创建服务参数文件
// 2018-01-11 21:02
type ServiceParam struct {
	// 客户端
	C1 *dynamic.Client
	// 命名空间
	Namespace string
	// 服务名称
	Name string
	// deploy名称
	ServiceName string
	// 网络模式
	NetworkMode string
	// cpu
	Cpu interface{}
	// 内存
	Memory string
	// 容器暴露端口
	Port string
	// 存储数据
	StorageData string
	// 镜像名称
	Image string
	// 标签选择器
	Selector string
	// 亲和性数据
	// 环境变量数据
	Envs string
	// 副本数量
	Replicas int64
	// 健康检查数据
	HealthData string
	// 客户端
	Cl2 *dynamic.Client
	// 端口数据
	PortData string
	// 获取数据客户端
	Cl3 kubernetes.Clientset
	// 更新类型
	Update bool
	// 滚动升级时候,会优先启动的pod数量
	MaxSurge int
	// 滚动升级时候,最大的unavailable数量
	MaxUnavailable int
	// 指定没有任何容器crash的Pod并被认为是可用状态的最小秒数
	MinReady int
	// 更新类型
	UpdateType string
	// 已经存在的端口数据
	OldPort util.Lock
	// 配置文件key
	ConfigureData []ConfigureData
	// 集群名称
	ClusterName string
	// 容器启动命令[]string
	Command string
	// 镜像仓库的
	// 集群主机IP
	Master string
	//
	MasterPort string
	// 资源空间
	ResourceName string
	// 应用名称
	AppName string
	// 创建用户
	CreateUser string
	// 配置文件是否更新参数
	NoUpdateConfig bool
	// 主机端口
	HostPort string
	// 安全设置
	Privileged bool
	// 镜像仓库地址
	Registry string
	// 镜像仓库用户名密码 admin:admin
	RegistryAuth string
	// 最大扩容量
	ReplicasMax int64
	// 标签
	Labels map[string]interface{}
	// 访问状态
	AccessMode string
	// 老的yaml信息,主要是获取端口
	PortYaml string
	// 重建标志
	IsRedeploy bool
	// pod关闭时间
	TerminationSeconds int
    // session 亲和性
	SessionAffinity  string
	// kafka 地址
	Kafka string
	// 日志路径,文件或目录,目录以/结尾
	LogPath string
	// 环境名称
	Ent string
	// es地址
	ElasticSearch string
	// 日志挂载路径
	LogDir string
}

// 配置服务健康检查使用的
type HealthData struct {
	// 检查类型
	HealthType string
	// 检查端口
	HealthPort string
	// http访问路径
	HealthPath string
	// 服务启动预计时间
	HealthInitialDelay string
	// 检查间隔
	HealthInterval string
	// 失败阈值
	HealthFailureThreshold string
	// 检查超时
	HealthTimeout string
	// 通过命令检查
	HealthCmd string
}

// 配置文件信息
// 2018-01-17 21:34
type ConfigureData struct {
	ContainerPath string
	DataName string
	DataId string
	ConfigDbData map[string]interface{}
}


const SelectCloudConfigureMount = "select service_name,mount_path,data_name,create_time,mount_id,configure_name,namespace,cluster_name,last_update_time from cloud_configure_mount"
const UpdateCloudConfigureMount = "update cloud_configure_mount"
const DeleteCloudConfigureMount = "delete from cloud_configure_mount"
const InsertCloudConfigureMount = "insert into cloud_configure_mount"
const SelectCloudLb = "select service_number,entname,cpu,memory,host_log_path,lb_ip,lb_type,lb_id,description,create_time,cluster_name,resource_name,last_modify_time,lb_name,lb_domain_prefix,lb_domain_suffix,create_user,last_modify_user,status from cloud_lb"




//2018-01-18 10:45:25.5832512 +0800 CST
type CloudConfigureMount struct {
	//首次挂载时间
	CreateTime string
	//
	MountId int64
	//配置文件名称
	ConfigureName string
	//命名空间
	Namespace string
	//集群名称
	ClusterName string
	//最近更新时间
	LastUpdateTime string
	// 数据名称
	DataName string
	// 挂载路径
	MountPath string
	//
	ServiceName string
}

// 仓库创建使用参数
// 2018-01-21 15:55
type RegistryParam struct {
	Name        string
	ClusterName string
	Master      string
	Port        string
	AuthServer  string
	HostPath    string
	Replicas    int64
}

// 镜像信息获取
// 2018-01-28 13:23

//2018-01-27 15:08:36.2055048 +0800 CST
type CloudImage struct {
	//创建用户
	CreateUser string
	//镜像名称
	Name string
	//所属仓库
	Repositories string
	//镜像类型,分为共有和私有
	ImageType string
	//
	ImageId int64
	//创建时间
	CreateTime string
	//镜像大小
	Size int64
	//镜像仓库组
	RepositoriesGroup string
	// tag数量
	TagNumber int
	// 访问方式
	Access string
	// 镜像层数
	LayersNumber int
	// 版本数据
	Tags string
	// 下载次数
	Download int64
}

//2018-01-18 16:23:05.0430682 +0800 CST
type CloudStorage struct {
	//
	StorageId int64
	//最近修改时间
	LastModifyTime string
	//创建时间
	CreateTime string
	//创建用户
	CreateUser string
	//描述信息
	Description string
	//存储大小,单位GB
	StorageSize string
	//存储格式
	StorageFormat string
	//最近修改用户
	LastModifyUser string
	//集群名称
	ClusterName string
	//glusterfs, nfs, host
	StorageType string
	// 名称
	Name string
	// 服务器地址
	StorageServer string
	// 使用装态
	Status string
	// 环境名称
	Entname string
	// 是否是共享存储 1 独有, 0 共享
	SharedType string
}

//2018-01-31 10:10:22.9723601 +0800 CST
type CloudStorageMountInfo struct {
	//服务名称
	ServiceName string
	//应用名称
	AppName string
	//创建用户
	CreateUser string
	//存储服务器
	StorageServer string
	//读写权限
	Model string
	//存储类型
	StorageType string
	//
	MountId int64
	//创建时间
	CreateTime string
	//集群名称
	ClusterName string
	//容器挂载路径
	MountPath string
	// 资源空间
	ResourceName string
	// 挂载状态
	Status string
	// 存储卷名称
	StorageName string
}

type CloudLbService struct {
	//转到容器的端口
	ContainerPort string
	//要负载的服务的名称
	ServiceName string
	//负载均衡名称
	LbName string
	//证书文件
	CertFile string
	//服务描述信息
	Description string
	//监听端口
	ListenPort string
	//负载均衡类型,tcp,http,https
	LbType string
	//集群名称
	ClusterName string
	//最近修改用户
	LastModifyUser string
	//
	ServiceId int64
	//最近修改时间
	LastModifyTime string
	//创建时间
	CreateTime string
	//创建用户
	CreateUser string
	// 访问地址
	Domain string
	// 应用名称
	AppName string
	// 资源空间
	ResourceName string
	// 应用服务对应的ID
	LbServiceId string
	// 负载均衡的ID
	LbId int64
	// 是否配置默认域名
	DefaultDomain string
	// 负载方式 pod node
	LbMethod string
	// 负载协议
	Protocol string
	// 服务版本
	ServiceVersion string
	// 环境名称
	Entname string
	// 流量切入百分比
	Percent int
	// 流量切入名称
	FlowServiceName string
}

type CloudClusterHosts struct {
	//主机IP
	HostIp string
	//主机标签
	HostLabel string
	//主机类型
	HostType string
	//是否有效
	IsValid int64
	//
	HostId int64
	// 集群名称
	ClusterName string
	// k8sAPi端口,只需要master有就行了
	ApiPort string
}


//2018-02-01 13:32:07.5158035 +0800 CST
type CloudLbNginxConf struct {
	//
	ConfId int64
	//创建用户
	CreateUser string
	//参考lb服务id
	LbServiceId string
	//资源空间
	ResourceName string
	//应用名称
	AppName string
	//集群名称
	ClusterName string
	//最近修改时间
	LastModifyTime string
	// 最近修改人
	LastModifyUser string
	//域名
	Domain string
	//vhost数据
	Vhost string
	//创建时间
	CreateTime string
	//服务名称
	ServiceName string
	// 负载服务ID
	ServiceId int64
	// 使用证书名称
	CertFile string
}

type CloudLb struct {
	//负载均衡名称
	LbName string
	//域名前缀
	LbDomainPrefix string
	//域名后缀
	LbDomainSuffix string
	//集群名称
	ClusterName string
	//资源空间
	ResourceName string
	//最近修改时间
	LastModifyTime string
	//创建用户
	CreateUser string
	//最近修改用户
	LastModifyUser string
	//
	Status string
	//IP地址
	LbIp string
	//负载均衡类型,nginx,haproxy
	LbType string
	//
	LbId int64
	//配额描述信息
	Description string
	//创建时间
	CreateTime string
	//
	ServiceNumber int64
	// 环境名称
	Entname string
	// cpu
	Cpu string
	// Memory
	Memory string
	// 日志挂载路径
	HostLogPath string
}

//2018-02-02 10:01:17.2337629 +0800 CST
type CloudLbCert struct {
	//描述信息
	Description string
	//
	CertId int64
	//创建时间
	CreateTime string
	//最近修改时间
	LastModifyTime string
	//最近修改用户
	LastModifyUser string
	//证书名称
	CertKey string
	//证书内容
	CertValue string
	//创建用户
	CreateUser string
	// 证书公钥文件
	PemValue string
}


//2018-02-06 10:40:26.9362807 +0800 CST
type CloudImageSyncLog struct {
	//镜像仓库组
	RegistryGroup string
	//仓库服务器2
	RegistryServer2 string
	//项目名称
	ItemName string
	//
	LogId int64
	//执行内容
	Messages string
	//创建用户
	CreateUser string
	//创建时间
	CreateTime string
	//程序运行时间
	Runtime int64
	//仓库服务器1
	RegistryServer1 string
	// 版本号
	Version string
	// 同步状态
	Status string
}

// 2018-02-13 09:46
// 镜像数据,
type HostImages struct {
	Id int
	Name string
	Tag string
	Size string
}

// 2018-02-16 18:36
// 服务滚动更新参数
type RollingParam struct {
	MinReadySeconds int32
	TerminationGracePeriodSeconds int64
	MaxUnavailable int32
	MaxSurge int32
	Namespace string
	Name string
	Client kubernetes.Clientset
	Images string
}

// 2018-02-17 21:09
// 更新upstream，参数
type UpdateLbNginxUpstream struct {
	Master string
	Port string
	Domain string
	Namespace string
	ServiceName string
	V CloudLbService
	ClusterName string
}

// 2018-02-19 14:50
// 获取自动扩容的扩容数量参数
type AutoScaleParam struct {
	ReplicasMax int32
	ReplicasMin int32
	Cpu int64
	Memory int64
	CreateUser string
	ResourceName string
}

//2018-02-20 09:39:59.7024273 +0800 CST
type CloudAutoScaleLog struct {
	//
	LogId int64
	//创建时间
	CreateTime string
	//指标类型
	MetricType string
	//缩容步长
	ReduceStep int64
	//扩展到
	Replicas int64
	//指标名称
	MetricName string
	//es连接地址
	Es string
	//扩容状态,成功失败
	Status string
	//监控值
	MonitorValue float64
	//最小值
	ReplicasMin int32
	//扩容步长
	IncreaseStep int64
	//最大值
	ReplicasMax int32
	//扩容或缩容间隔
	ActionInterval int64
	//查询参数
	Query string
	//应用名称
	AppName string
	//集群名称
	ClusterName string
	//最近几次超过阈值
	LastCount int64
	// 配置的阈值
	Gt int64
	// 步长
	Step string
	// 服务名称
	ServiceName string
	// 环境名称
	Entname string
}

// 2018-01-29 14:28
// 创建nfs存储
type StorageParam struct {
	// 存储名称
	Name string
	// 存储大小
	Size string
	// 访问模式
	AccessMode string
	// 集群MasterIP
	Master string
	// 集群MasterPort
	Port string
	//
	Namespace string
	// pvc name
	PvcName string
	// 宿主机地址
	HostPath string
	// pvc类型
	StorageType string
	// 集群名称
	ClusterName string
}

// 2018-02-27 20:51
type EventData struct {
	// 事件事件
	EventTime string
	// 信息
	Messages string
	// 原因
	Reason string
	// 主机Ip
	Host string
	// 类型
	Type string
}


type NodeStatus struct {
	hosts.CloudClusterHosts
	Lables     []string
	K8sVersion string
	ErrorMsg   string
	MemSize    int64
	OsVersion string
}

type ClusterStatus struct {
	ClusterId    int64
	ClusterType  string
	NodeStatus
	ClusterAlias string
	ClusterName  string
	Nodes        int64
	Services     int
	OsVersion string
}

type ClusterResources struct {
	UsedCpu       int64
	UsedMem       int64
	CpuUsePercent float64
	MemUsePercent float64
	MemFree       int64
	CpuFree       int64
	Cpu           int64
	Mmeory        int64
	Services      int
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
	// 容器数量
	ContainerNumber int
	// 集群名称
	ClusterName string
	// 失败的容器数量
	ContainerFail int
	// 失败的服务数量
	ServiceFail int
	// 服务总量
	ServiceNumber int
	// 域名
	Domain string
	// svc + namespace:8888
	// 访问方式
	Access []string
	// 镜像名称
	Image string
	// 服务名称
	ServiceName string
	// 服务ID
	ServiceId int64
	// 正在运行的数量
	AvailableReplicas int32
	// 环境名称
	Entname string
	// 检查时间
	CheckTime int64
}

// 2018-03-01 14:24
// 获取集群证书文件
type CertData struct {
	// ca证书公钥文件
	CaData string
	// node证书公钥内容
	CertData string
	// node证书私钥内容
	KeyData string
}

// 2018-09-04 09:51
// 集群节点资源使用情况
type NodeReport struct {
	//
	Ip string
	Namespace string
	Name string
	CpuRequests string
	MemoryRequests string
	CpuLimits string
	MemoryLimits string
}