package registry

type CloudRegistryServerName struct {
    //集群名称
    ClusterName string
    // 服务器地址
    ServerAddress string
    // 名字
    Name string
    // 域名
    ServerDomain string
}

//2018-01-18 16:23:17.4387771 +0800 CST
type CloudDeployImage struct {
    //
    ServerAddress string
    //集群名称
    ClusterName string
    //
    ServerDomain string
    // 镜像名称
    Name string
    // 镜像tag
    Tags string
}

//2018-01-18 16:23:17.4387771 +0800 CST
type CloudRegistryServer struct {
    //
    ServerAddress string
    //集群名称
    ClusterName string
    //
    ServerId int64
    //最近修改用户
    LastModifyUser string
    //创建时间
    CreateTime string
    //创建用户
    CreateUser string
    //镜像数量
    ImagesNumber int64
    //镜像用户名名
    Username string
    //
    ServerDomain string
    //描述信息
    Description string
    //镜像前缀 如 online test develop ${user}
    Prefix string
    //镜像密码
    Password string
    //最近修改时间
    LastModifyTime string
    //组名称,不同的服务属于不同的组
    Groups string
    //
    GroupsId int64
    // 认证服务器
    AuthServer string
    // 管理员
    Admin string
    // 名字
    Name string
    // 访问信息
    Access string
    // 环境名称
    Entname string
    // 主机挂载路径
    HostPath string
    // 运行状态
    Status string
    // 副本数量
    Replicas int64

}

//2018-01-19 22:15:41.9929294 +0800 CST
type CloudRegistryPermissions struct {
    //注册中心名称
    ServiceName string
    //用户名称
    UserName string
    //项目名称
    Project string
    //镜像名称
    ImageName string
    //创建时间
    CreateTime string
    //最近修改用户
    LastModifyUser string
    //
    PermissionsId int64
    //用户类型,用户或组
    UserType string
    //注册中心地址
    RegistryServer string
    //用户组名称
    GroupsName string
    //创建用户
    CreateUser string
    //最近修改时间
    LastModifyTime string
    // 操作权限
    Action string
    // 集群名称
    ClusterName string
    // 描述信息
    Description string
}

//2018-01-27 15:08:25.7829086 +0800 CST
type CloudRegistryGroup struct {
    // tag 数量
    TagNumber int64
    //镜像数量
    ImageNumber int64
    //最近修改时间
    LastModifyTime string
    //创建用户
    CreateUser string
    //镜像类型,分为共有和私有
    GroupType string
    //最近修改用户
    LastModifyUser string
    //镜像总大小
    SizeTotle int64
    //
    GroupId int64
    //创建时间
    CreateTime string
    // 组名称
    GroupName string
    // 镜像服务器域名
    ServerDomain string
    // 集群名称
    ClusterName string
}


//2018-01-28 14:35:48.4221703 +0800 CST
type CloudImageLog struct {
    //操作时间
    CreateTime string
    //操作人
    CreateUser string
    //镜像仓库组
    RepositoriesGroup string
    //镜像获取类型,pull,create,push
    OperType string
    //标签名称
    Label string
    //
    LogId int64
    //镜像名称
    Name string
    //所属仓库
    Repositories string
    //所属集群
    ClusterName string
    // ip
    Ip string
}

//2018-02-06 16:46:02.7429371 +0800 CST
type CloudImageSync struct {
    //
    SyncId int64
    //集群名称
    ClusterName string
    //目标集群名称
    TargetCluster string
    //备注信息
    Description string
    //仓库组
    RegistryGroup string
    //镜像名称
    ImageName string
    //镜像源环境名称
    Entname string
    //创建时间
    CreateTime string
    //创建用户
    CreateUser string
    //目标仓库服务器
    TargetRegistry string
    //版本号
    Version string
    //目标集群
    TargetEntname string
    //最近修改时间
    LastModifyTime string
    //审批时间
    ApprovedTime string
    //最近修改用户
    LastModifyUser string
    //审批人
    ApprovedBy string
    // 同步状态
    Status string
    // 来源仓库地址
    Registry string
}

//2018-02-09 16:18:14.7085342 +0800 CST
type CloudImageBase struct {
    //
    BaseId int64
    //镜像仓库地址
    RegistryServer string
    //创建时间
    CreateTime string
    //创建用户
    CreateUser string
    //镜像名称
    ImageName string
    //最近修改时间
    LastModifyTime string
    //最近修改用户
    LastModifyUser string
    //镜像图标
    Icon string
    //镜像描述信息
    Description string
    // 镜像类型
    ImageType string
}
