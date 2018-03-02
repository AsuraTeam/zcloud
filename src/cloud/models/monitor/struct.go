package monitor

import (
    "cloud/util"
)

// 2018-02-20 11:03
// 监控服务器访问地址缓存
type PrometheusServer struct {
    Host string
    Port string
}

//2018-02-19 18:13:04.8039696 +0800 CST
type CloudAutoScale struct {
    //开始时间
    Start string
    //查询参数
    Query string
    //命名空间
    Namespace string
    //服务版本号
    ServiceVersion string
    //扩容或缩容间隔
    ActionInterval int
    //集群名称
    ClusterName string
    //最近几次超过阈值
    LastCount int64
    //查询监控时间步长
    Step string
    //扩容或缩容进行时,发送通知组
    MsgGroup string
    //阈值大于多少
    Gt int64
    //扩容步长
    IncreaseStep int64
    //
    ScaleId int64
    //创建时间
    CreateTime string
    //阈值小于多少
    Lt int64
    //服务的名称
    ServiceName string
    //创建用户
    CreateUser string
    //应用名称
    AppName string
    //资源名称
    ResourceName string
    //缩容步长
    ReduceStep int64
    //最近修改时间
    LastModifyTime string
    //最近修改用户
    LastModifyUser string
    //结束时间
    End string
    // 环境名称
    Entname string
    // 描述信息
    Description string
    // 指标类型
    MetricType string
    // 指标名称，系统自带
    MetricName string
    // 数据源,prometheus,es
    DataSource string
    // es主机地址
    Es string
}

// 2018-02-19 21:01
// 获取自动扩容配置默认配置
func GetQueryParamDefault() CloudAutoScale  {
    return CloudAutoScale{
        LastCount:5,
        ActionInterval:120,
        IncreaseStep:1,
        ReduceStep:1,
        Step:"5",
        Gt:80,
        Lt:50,
        MetricName:util.GetSelectOptionName("cpu"),
    }
}
