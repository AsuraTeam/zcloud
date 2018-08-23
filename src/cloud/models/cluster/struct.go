package cluster

import "cloud/k8s"

// 查询集群名称时使用
type CloudClusterName struct {
    ClusterId int64
    //集群显示名称
    ClusterAlias string
    //集群名称,必须英文
    ClusterName string
}

type CloudCluster struct {
    //docker安装路径
    DockerInstallDir string
    //内网网卡名称
    NetworkCart string
    //
    ClusterId int64
    //集群显示名称
    ClusterAlias string
    //最近修改时间
    LastModifyTime string
    //docker版本
    DockerVersion string
    //集群类型
    ClusterType string
    //集群名称,必须英文
    ClusterName string
    //创建时间
    CreateTime string
    //创建用户
    CreateUser string
    // ca证书公钥文件
    CaData string
    // node证书公钥内容
    CertData string
    // node证书私钥内容
    KeyData string
    // 主节点地址
    ApiAddress string
}

// 集群页面管理使用数据
type CloudClusterDetail struct {
    CloudCluster
    k8s.ClusterResources
    ClusterMem int64
    ClusterCpu int64
    ClusterNode int64
    ClusterService int64
    ClusterPods int
    Services int
    Couters int
    Health string
    OsVersion string
}