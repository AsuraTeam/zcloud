package k8s

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"cloud/util"
	"github.com/astaxie/beego/logs"
	"cloud/cache"
	"time"
)

type ClusterHealth struct {
	// 插件名称
	Name string
	// 状态
	Status string
	// 信息
	Message string
}

// 2018-02-28 11:16
// 每次读取
// 获取集群组件监控状态
func GetClusterStatus(clusterName string) string {
	health := make([]ClusterHealth, 0)
	cl, err := GetClient(clusterName)
	if err != nil {
		logs.Error("获取集群连接失败", clusterName, err)
		return ""
	}
	data ,err := cl.CoreV1().ComponentStatuses().List(metav1.ListOptions{})
	if err == nil {
		for _,v := range data.Items {
			temp := ClusterHealth{}
			if len(v.Conditions) == 0 {
				continue
			}
			temp.Message = v.Conditions[0].Message
			temp.Name = v.Name
			temp.Status = util.ObjToString(v.Conditions[0].Status)
			health = append(health, temp)
		}
	}
	cache.ClusterComponentStatusesCache.Put(clusterName, util.ObjToString(health), time.Minute * 10)
	return util.ObjToString(health)
}
