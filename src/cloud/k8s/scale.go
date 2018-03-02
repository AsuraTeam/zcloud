package k8s

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"github.com/astaxie/beego/logs"
	"k8s.io/api/autoscaling/v2beta1"
	"cloud/sql"
	"cloud/util"
	"time"
	"github.com/garyburd/redigo/redis"
	"strconv"
	"cloud/cache"
)

const SelectServiceReplicas = `select replicas_max,resource_name,replicas_min,cpu,memory,create_user from cloud_app_service where  cluster_name=? and service_name=? and app_name=? and service_version=?`

// 在监控数据超过后查询
// 获取扩容实例数量
func getServiceReplicas(param QueryParam) AutoScaleParam {
	data := AutoScaleParam{}
	sql.GetOrm().Raw(SelectServiceReplicas,
		param.ClusterName, param.ServiceName,
		param.AppName, param.ServiceVersion).QueryRow(&data)
	logs.Info(util.ObjToString(data))
	return data
}

// 2018-01-13 10:32
// 将应用扩展或停止
// k8s.ScalePod("10.16.55.6","8080","auto-3--dfsad","auto-3",4)
func ScalePod(clustername string, namespace string, name string, replicas int32) error {
	cl, _ := GetClient(clustername)
	logs.Info("伸缩节点", namespace, name)
	deployments, err := cl.ExtensionsV1beta1().Deployments(namespace).Get(name, metav1.GetOptions{})
	if err != nil {
		logs.Error("获取服务失败", err, namespace, name, deployments)
		return err
	}
	var replicasv *int32
	replicasv = &replicas
	deployments.Spec.Replicas = replicasv
	d, err := cl.ExtensionsV1beta1().Deployments(namespace).Update(deployments)
	logs.Info("更新Deploy节点数量:", d, err)
	return err
}

// 2018-01-13 17:54
// 获取扩展信息
// {"metadata":{"name":"auto-3","namespace":"auto-3--dfsad","selfLink":"/apis/autoscaling/v2beta1/namespaces/auto-3--dfsad/horizontalpodautoscalers/auto-3","uid":"3a4ec944-f83e-11e7-8d1c-0894ef37b2d2","resourceVersion":"4149862","creationTimestamp":"2018-01-13T08:46:20Z"},"spec":{"scaleTargetRef":{"kind":"Deployment","name":"auto-3","apiVersion":"extensions/v1beta1"},"minReplicas":1,"maxReplicas":2,"metrics":[{"type":"Resource","resource":{"name":"cpu","targetAverageUtilization":80}}]},"status":{"lastScaleTime":"2018-01-13T09:45:51Z","currentReplicas":2,"desiredReplicas":2,"currentMetrics":null,"conditions":[{"type":"AbleToScale","status":"True","lastTransitionTime":"2018-01-13T08:46:50Z","reason":"SucceededGetScale","message":"the HPA controller was able to get the target's current scale"},{"type":"ScalingActive","status":"False","lastTransitionTime":"2018-01-13T08:46:50Z","reason":"FailedGetResourceMetric","message":"the HPA was unable to compute the replica count: unable to get metrics for resource cpu: failed to get pod resource metrics: the server could not find the requested resource (get services http:heapster:)"}]}} <nil>
func GetAutoScale(clustername string, namespace string, name string) (v2beta1.HorizontalPodAutoscaler, error) {
	cl, _ := GetClient(clustername)
	v, err := cl.AutoscalingV2beta1().HorizontalPodAutoscalers(namespace).Get(name, metav1.GetOptions{})
	return *v, err
}

// 2018-02-20 18:59
// 添加锁
func setLock(lockKey string, interval time.Duration)  {
	logs.Info(lockKey, interval)
	cache.AutoScaleCache.Put(lockKey, time.Now().Unix(), interval)
}

// 2018-02-19 15:14
// 扩容pod
func increasePod(param QueryParam) {
	lockKey := param.Namespace + param.ClusterName + param.ServiceName + param.ServiceVersion + param.AppName + "increase"
	if !checkAutoScaleLock(lockKey) {
		return
	}

	// 扩容
	logs.Info("开始扩容")
	client, err := GetClient(param.ClusterName)
	if err == nil {
		serviceName := util.Namespace(param.ServiceName, param.ServiceVersion)
		deployment := GetDeployment(param.Namespace, client, serviceName)
		if deployment.Name == "" {
			logs.Info("获取到deployment失败", param.Namespace, serviceName)
			return
		}
		replicas := *deployment.Spec.Replicas
		if replicas > 0 {
			data := getServiceReplicas(param)
			if replicas < data.ReplicasMax {
				replicas += param.IncreaseStep
				if replicas > data.ReplicasMax {
					replicas -= data.ReplicasMax - replicas
					logs.Info("扩容超过实际配额", replicas, data.ReplicasMax)
				}
			} else {
				logs.Info("已经扩容到最大配额", serviceName, param.ClusterName, param.AppName)
				setLock(lockKey, time.Second* time.Duration(param.ActionInterval))
				return
			}
			status, msg := CheckQuota(data.CreateUser,int64(replicas), data.Cpu, data.Memory, data.ResourceName)
			if !status {
				logs.Error(msg)
				// 超过配额后锁定
				setLock(lockKey, time.Second* time.Duration(param.ActionInterval))
				return
			}
			logs.Info("扩容到", replicas, serviceName)
			if cache.AutoScaleCacheErr == nil {
				setLock(lockKey, time.Second* time.Duration(param.ActionInterval))
			}
			writeScaleLog(param, data, replicas,"扩容")
			ScalePod(param.ClusterName, param.Namespace,serviceName, replicas)
		}
	}
}

// 2018-02-19 16:23
// 检查自动扩容是否锁定
func checkAutoScaleLock(key string) bool {
	r := cache.AutoScaleCache.Get(key)
	if r != nil {
		redisR, err := redis.String(r, nil)
		logs.Error("扩容间隔太短", redisR, err)
		return false
	}
	return true
}

// 2018-02-19 15:20
// 缩小pod
func reducePod(param QueryParam) {
	lockKey := param.Namespace + param.ClusterName + param.ServiceName + param.ServiceVersion + param.AppName + "reduce"
	if !checkAutoScaleLock(lockKey) {
		return
	}
	// 缩容
	logs.Info("开始缩容")
	client, err := GetClient(param.ClusterName)
	if err == nil {
		serviceName := util.Namespace(param.ServiceName, param.ServiceVersion)
		deployment := GetDeployment(param.Namespace, client, serviceName)
		if deployment.Name == "" {
			logs.Info("获取到deployment失败", param.Namespace, serviceName)
			return
		}
		replicas := *deployment.Spec.Replicas
		logs.Info("获取到实际pod数量为", replicas)
		if replicas > 1 {
			data := getServiceReplicas(param)

			if replicas > data.ReplicasMin {
				replicas -= param.IncreaseStep
				logs.Info("缩小到实际数量", replicas)
				if replicas < data.ReplicasMin {
					replicas += data.ReplicasMin - replicas
					logs.Info("缩容小于实际配额", replicas, data.ReplicasMax)
				}
			} else {
				logs.Info("已经缩小到最小配额", serviceName, param.ClusterName, param.AppName)
				setLock(lockKey, time.Second* time.Duration(param.ActionInterval))
				return
			}
			logs.Info("缩容到", replicas, serviceName)
			if cache.AutoScaleCacheErr == nil {
				setLock(lockKey, time.Second* time.Duration(param.ActionInterval))
			}
			writeScaleLog(param, data, replicas,"缩容")
			ScalePod(param.ClusterName, param.Namespace, serviceName, replicas)
		}else{
			setLock(lockKey, time.Second * time.Duration(param.ActionInterval))
		}
	}
}

// 2018-02-19 14:36
// 分析监控数据,并做出扩容和缩容操作
func ParseMonitorData(param QueryParam) {
	for i := 1; i < 3; i ++ {
		param.ServiceVersion = strconv.Itoa(i)
		// 如果没有配置,默认2分钟
		param.Interval = time.Second * time.Duration(param.ActionInterval)
		v := GetLastCountData(param)
		param.MonitorValue = v
		logs.Info("获取到监控数据为", v, param.Gt, param.Interval)
		if v > float64(param.Gt) {
			// 扩容
			increasePod(param)
			// 如果扩容和缩容配置不合理,先扩容,扩容后不缩容
			return
		}
		if v > 0 && v < float64(param.Gt) {
			// 缩容
			reducePod(param)
		}
	}
}

// 2018-02-20 17:10
// 记录扩容操作日志
const InsertCloudAutoScaleLog = "insert into cloud_auto_scale_log"
func writeScaleLog(param QueryParam,scaleParam AutoScaleParam, replicas int32, status string)  {
	data := CloudAutoScaleLog{}
	util.MergerStruct(param, &data)
	data.ReplicasMax = scaleParam.ReplicasMax
	data.ReplicasMin = scaleParam.ReplicasMin
	data.Replicas = int64(replicas)
	data.Status = status
	data.CreateTime = util.GetDate()
	q := sql.InsertSql(data, InsertCloudAutoScaleLog)
	sql.Raw(q).Exec()
}