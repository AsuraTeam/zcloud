package groups

//2018-01-20 06:33:32.8534847 +0800 CST
type CloudUserGroups struct {
    //创建用户
    CreateUser string
    //最近修改时间
    LastModifyTime string
    //最近修改用户
    LastModifyUser string
    //创建时间
    CreateTime string
    //组名称
    GroupsName string
    //组成员,用逗号分隔
    Users string
    //
    GroupsId int64

    // 描述信息
    Description string
}

type CloudUserGroupsName struct {
    //组名称
    GroupsName string
    //
    GroupsId int64
}