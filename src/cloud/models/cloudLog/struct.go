package cloudLog

type CloudOperLog struct {
    //操作信息
    Messages string
    //在哪个集群操作的
    Cluster string
    //操作IP地址
    Ip string
    //
    LogId int64
    //操作时间
    Time string
    //操作用户
    User string
}
