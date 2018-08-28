package index

type DockerCloudCiEnv struct {
    //部门名称
    Dept string
    //域名
    Domain string
    //组名称
    GroupsName string
    //自动编译时间
    AutoCompileTime string
    //容器最大扩容到数量
    ContainerMaxNumber int64
    //编译脚本
    BuildScript string
    //容器数量
    ContainerNumber int64
    //监控检查脚本
    HealthScript string
    //审批用户
    ApproveUser string
    //创建人
    CreateUser string
    //服务名称和服务管理的名称对应
    ServiceName string
    //是否自动编译
    AutoCompile string
    //
    BuildStatus string
    //
    EnvId int64
    //最近修改时间
    LastModifyTime string
    //代码分支
    CodeBranch string
    //镜像来源,手动上传,脚本编译,公共镜像
    ImagesTp string
    //cpu资源限制
    Cpu int64
    //内存资源限制,单位MB
    Memory int64
    //备注信息
    Description string
    //项目类型
    ItemTp string
    //参考主机组
    GroupsId int64
    //自动编译类型
    AutoCompileTp string
    //编译状态
    AutoCompileStatus string
    //镜像仓库
    Images string
    //审批时间
    ApproveTime string
    //其他未知信息
    GsonData string
    //代码路径
    CodePath string
    //环境组信息
    GroupsIds string
    //最近修改人
    LastModifyUser string
    //创建时间
    CreateTime string
    //环境名称
    Entname string
    //审批状态
    Approve string
}

// 2018-01-20 12:53
// 存储用户
type CloudAuthorityUser struct {
    UserName string
}

type DockerCloudAuthorityUser struct {
    //真实姓名
    ThirdTrueName string
    //用户邮箱
    UserEmail string
    //是否启用（0无效，1有效）
    IsValid int64
    //是否删除（0未删除，1删除）
    IsDel int64
    //用户名称
    UserName string
    //第三方id
    ThirdId string
    //用户电话
    UserMobile string
    //用户id
    UserId int64
    //用户头像
    UserPic string
    // 密码
    Pwd string
    // 描述信息
    Description string
    // 真实姓名
    RealName string
    // 所属部门
    Dept string
    //创建用户
    CreateUser string
    //最近修改时间
    LastModifyTime string
    //最近修改用户
    LastModifyUser string
    //创建时间
    CreateTime string
    // token
    Token string
}

type CloudLoginRecord struct {
    //登录状态
    LoginStatus int64
    //
    RecordId int64
    //登录时间
    LoginTime string
    //登录IP
    LoginIp string
    //登录用户名
    LoginUser string
}

