package cache

import "cloud/util"

var (
	// api 访问次数缓存计数
	ApiLimitCache, _ = util.RedisCacheClient("cloud_api_limit_")
	// 缓存监控自动扩容时的服务数据
	ServiceDataCache, _ = util.RedisCacheClient("cloud_cache_service_data_")
	// 服务信息缓存
	ServiceInfoCache, _ = util.RedisCacheClient("cloud_cache_service_info_")
	// 2018-02-18 21:06
	// 设置权限缓存,所有权限查询都从缓存中读取
	PermCache, _ = util.RedisCacheClient("cloud_ci_perm_")
	// 2018-02-08 21:54
	AppCache, AppCacheErr = util.RedisCacheClient("cloud_k8s_app_")
	// 集群数据缓存
	ClusterCache, ClusterCacheErr = util.RedisCacheClient("cloud_cluster_cache_")
	// 集群组件监控状态数据
	ClusterComponentStatusesCache, _ = util.RedisCacheClient("cloud_component_statuses_cache_")
	// 自动扩容缓存
    AutoScaleCache, AutoScaleCacheErr = util.RedisCacheClient("cloud_auto_scale_")
	// 仓库认证缓存数据
    RedisUserCache, _ = util.RedisCacheClient("cloud_user_")
	// 仓库权限缓存
	RegistryPermCache, _ = util.RedisCacheClient("cloud_p_")
	// 仓库更新镜像缓存镜像ID
	RegistryLogCache , _ = util.RedisCacheClient("cloud_image_log_")
	// 集群主机信息缓存
    HostCache, HostCacheErr = util.RedisCacheClient("cluster_node_status")
	// 容器信息缓存
    ContainerCache, ContainerCacheErr = util.RedisCacheClient("cloud_container_")
	// 应用服务数据缓存
	ServiceCache, ServiceCacheErr = util.RedisCacheClient("cloud_service_")
	// 执行构建任务的缓存
    JobDataCache, JobDataErr = util.RedisCacheClient("cloud_job_data_")
	// 任务计划执行缓存
    JobCache, JobCacheErr = util.RedisCacheClient("cloud_job_history_")
	// 缓存集群master数据
    MasterCache, _ = util.RedisCacheClient("cloud_master_")
	// 任务计划,pod数据缓存
	PodCache, podErr = util.RedisCacheClient("job_pod_")
)
