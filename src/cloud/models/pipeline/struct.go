package pipeline

//2018-02-03 17:08:23.9574221 +0800 CST
type CloudPipeline struct {
    //创建时间
    CreateTime string
    //
    PipelineId int64
    //应用名称
    AppName string
    //集群名称
    ClusterName string
    //服务名称
    ServiceName string
    //最近修改时间
    LastModifyTime string
    //资源空间
    ResourceName string
    //流水线名称
    PipelineName string
    // 最近修改人
    LastModifyUser string
    //构建任务内容,关联构建任务
    JobName string
    //创建用户
    CreateUser string
    //执行时间
    ExecTime string
    // 描述信息
    Description string
    // 过程失败执行动作
    FailAction string
    // 构建服务Id
    JobId int64
    // 状态
    Status string
}

//2018-02-05 12:53:22.6536907 +0800 CST
type CloudPipelineLog struct {
    //应用名称
    AppName string
    //集群名称
    ClusterName string
    //创建用户
    CreateUser string
    //结束构建任务时间时间
    EndJobTime string
    //更新服务结束时间
    UpdateServiceEndTime string
    //服务名称
    ServiceName string
    //资源空间
    ResourceName string
    //启动执行时间
    StartTime string
    //构建任务是否成功
    JobStatus string
    //提交镜像时间
    PushImageStartTime string
    //更新服务启动时间
    UpdateServiceStartTime string
    //创建时间
    CreateTime string
    //执行时间
    ExecTime string
    //构建任务ID
    JobId int64
    //运行时间
    RunTime int64
    //启动构建任务时间时间
    StartJobTime string
    //提交镜像完成时间
    PushImageEndTime string
    //流程结束时间
    EndTime string
    //
    LogId int64
    //流水线名称
    PipelineName string
    //构建任务内容,关联构建任务
    JobName string
    //执行状态,成功或失败
    Status string
    //执行日志
    Messages string
    // 更新服务状态
    UpdateServiceStatus string
    // 更新服务错误信息
    UpdateServiceErrormsg string
    // 构建状态
    BuildJobStatus string
    // 构建错误信息
    BuildJobErrormsg string
}
