package ci

//2018-01-24 16:54:53.7454432 +0800 CST
type CloudCodeRepostitory struct {
    //
    RepostitoryId int64
    //最近修改时间
    LastModifyTime string
    //最近修改用户
    LastModifyUser string
    //密码,base64存储
    Password string
    //gitlab token
    GitlabToken string
    //代码类型, public private 共有和私有
    Type string
    //创建时间
    CreateTime string
    //创建用户
    CreateUser string
    //代码来源, gitlab,github,svn 
    CodeSource string
    //用户名
    Username string
    // 代码路径
    CodeUrl string
    // 描述信息
    Description string
}

//2018-01-24 21:30:36.1428558 +0800 CST
type CloudCiDockerfile struct {
    //文件名称
    Name string
    //最近修改时间
    LastModifyTime string
    //
    FileId int64
    //dockerfile内容
    Content string
    // 编译脚本
    Script string
    //创建时间
    CreateTime string
    //创建用户
    CreateUser string
    //最近修改用户
    LastModifyUser string
    //描述信息
    Description string
    //是否删除,1删除
    IsDel int64
}


//2018-01-25 17:50:32.6516446 +0800 CST
type CloudBuildJob struct {
    //项目名称
    ItemName string
    //镜像tag
    ImageTag string
    //集群名称
    ClusterName string
    //参考code代码仓库
    JobCode int64
    //dockerfile参考
    DockerFile string
    //创建用户
    CreateUser string
    //最近修改用户
    LastModifyUser string
    //任务计划名称
    JobName string
    //创建时间
    CreateTime string
    //build时k8s.job名称
    BuildId string
    //构建状态
    BuildStatus string
    //
    JobId int64
    //镜像仓库配置
    RegistryServer string
    //最近修改时间
    LastModifyTime string
    //最近构建时间
    LastBuildTime string
    // 描述信息
    Description string
    // 自定义dockerfile
    Content string
    // 编译脚本
    Script string
    // 超时时间
    TimeOut int
    // 最近tag
    LastTag string
    // 基础镜像
    BaseImage string
}

//2018-01-26 15:22:01.3732277 +0800 CST
type CloudBuildJobHistory struct {
    //注册服务器
    RegistryServer string
    //
    HistoryId int64
    //任务计划名称
    JobName string
    //创建时间
    CreateTime string
    //构建时间
    BuildTime int64
    //项目名称
    ItemName string
    //镜像tag
    ImageTag string
    //集群名称
    ClusterName string
    //创建用户
    CreateUser string
    //构建状态
    BuildStatus string
    //构建日志
    BuildLogs string
    //参考job表ID
    JobId int64
    //
    DockerFile string
    // 编译脚本
    Script string
    // 仓库组
    RegistryGroup string
    // 基础镜像
    BaseImage string
}

//2018-02-10 18:22:04.0393696 +0800 CST
type CloudCiService struct {
    //应用名称
    AppName string
    //最近修改用户
    LastModifyUser string
    //最近修改时间
    LastModifyTime string
    //服务描述信息
    Description string
    //镜像名称
    ImageName string
    //组名称
    GroupName string
    //创建时间
    CreateTime string
    //创建用户
    CreateUser string
    //发布类型,金丝雀,蓝绿,滚动
    ReleaseType string
    //
    ServiceId int64
    //服务名称
    ServiceName string
    //环境名称
    Entname string
    //域名
    Domain string
    //集群名称
    ClusterName string
    //发布状态
    Status string
    // 当前版本
    CurrentVersion string
    // 负载均衡服务使用的版本
    LbVersion string
    // 镜像信息绿版本信息
    ImageInfoGreen string
    // 镜像信息蓝版本信息
    ImageInfoBlue string
    // 蓝版访问方式
    BlueAccess string
    // 绿版访问方式
    GreenAccess string
    // 蓝版指向的服务
    BluePod string
    //
    LbService string
    // 流量切入百分比
    Percent int
    // 最新版本,按创建时间获取
    NewVersion string
}

//2018-02-14 09:30:07.7331635 +0800 CST
type CloudCiReleaseHistory struct {
    //集群名称
    ClusterName string
    //发布类型
    ReleaseOnlineType string
    //禅道Bug链接
    ReleaseBugPmLink string
    //镜像名称
    ImageName string
    //发布状态
    Status string
    //
    HistoryId int64
    //创建时间
    CreateTime string
    //服务名称
    ServiceName string
    //发布类型,金丝雀,蓝绿,滚动
    ReleaseType string
    //Bug修复功能描述
    ReleaseBugDescription string
    //需求名称
    ReleaseDemandDescription string
    //创建用户
    CreateUser string
    //域名
    Domain string
    //应用名称
    AppName string
    //禅道任务链接
    ReleaseJobPmLink string
    //测试人员
    ReleaseTestUser string
    //环境名称
    Entname string
    //
    Description string
    //项目名称
    ReleaseItemDescription string
    //
    ServiceId int64
    // 是否自动切换
    AutoSwitch string
    // 负载均衡版本
    LbVersion string
    // 发布版本
    ReleaseVersion string
    // 发布或回滚,update,rollback
    Action string
    // 流量切入百分比
    Percent int
    // 旧版本的镜像名称,做滚动更新后回滚使用
    OldImages string
}

//2018-02-17 11:11:53.4396702 +0800 CST
type CloudCiReleaseLog struct {
    //
    LogId int64
    //域名
    Domain string
    //应用名称
    AppName string
    //环境名称
    Entname string
    //
    Messages string
    //服务名称
    ServiceName string
    //集群名称
    ClusterName string
    //发布或回滚,update,rollback
    Action string
    //创建时间
    CreateTime string
    //创建用户
    CreateUser string
    // 操作IP
    Ip string
}

//2018-02-18 12:02:14.5230327 +0800 CST
type CloudCiPerm struct {
    //拥有权限
    Datas string
    //创建时间
    CreateTime string
    //最近修改时间
    LastModifyTime string
    //最近修改用户
    LastModifyUser string
    //
    PermId int64
    //权限用户
    Username string
    //组用户
    GroupsName string
    //创建用户
    CreateUser string
}
