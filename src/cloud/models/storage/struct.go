package storage

//2018-01-18 16:22:58.0616688 +0800 CST
type CloudStorageServer struct {
    //
    ServerAddress string
    //最近修改用户
    LastModifyUser string
    //创建时间
    CreateTime string
    //创建用户
    CreateUser string
    //glusterfs, nfs, host
    StorageType string
    //独享型，共享型
    UsedType string
    //
    ServerId int64
    //最近修改时间
    LastModifyTime string
    //描述信息
    Description string
    //集群名称
    ClusterName string
    // 环境名称
    Entname string
    // 磁盘路径或目录路径
    HostPath string
}



