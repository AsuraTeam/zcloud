package lb

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

// 2018-02-14 13:26
type LbServiceVersion struct {
    // 域名
    Domain string
    // 服务版本
    ServiceVersion string
    // 环境名称
    Entname string
}


