package k8s

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"cloud/util"
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
	health := []ClusterHealth{}
	cl, _ := GetClient(clusterName)
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
	return util.ObjToString(health)
}
