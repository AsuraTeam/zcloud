package perm

//2018-02-06 08:22:35.6731915 +0800 CST
type CloudPermRole struct {
    //创建时间
    CreateTime string
    //最近修改时间
    LastModifyTime string
    //最近修改用户
    LastModifyUser string
    //服务描述信息
    Description string
    //
    RoleId int64
    //角色名称
    RoleName string
    //拥有权限
    Permissions string
    //创建用户
    CreateUser string
}

//2018-02-06 08:22:40.5954731 +0800 CST
type CloudPerm struct {
    //最近修改时间
    LastModifyTime string
    //最近修改用户
    LastModifyUser string
    //拥有团队
    Groups string
    //拥有权限角色
    Roles string
    //创建时间
    CreateTime string
    //创建用户
    CreateUser string
    //服务描述信息
    Description string
    //拥有用户
    User string
    //
    PermId int64
}

//2018-02-06 08:22:46.2597971 +0800 CST
type CloudApiResource struct {
    //创建时间
    CreateTime string
    //创建用户
    CreateUser string
    //
    ResourceId int64
    //最近修改用户
    LastModifyUser string
    //服务描述信息
    Description string
    //apiurl地址
    ApiUrl string
    //api名称
    Name string
    //是否是公开的,公开的将不受权限控制
    ApiType string
    //最近修改时间
    LastModifyTime string
    // 请求方法
    Method string
    // 父节点
	Parent string
}


//2018-08-23 08:38:08.9273685 +0800 CST
type CloudUserPerm struct {
    //
    ResourceType string
    //
    LastModifyTime string
    //
    LastModifyUser string
    //
    Ent string
    //
    Description string
    //
    ParentUser string
    //
    Name string
    //
    UserName string
    //
    ResourceName string
    //
    SubUser string
    //
    CreateUser string
    //
    CreateTime string
    //
    ClusterName string
    //
    PermId int64
    //
    GroupName string
}
