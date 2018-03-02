package k8s

import (
	"k8s.io/client-go/kubernetes"
)

// 获取资源使用情况,cpu，内存
func GetClusterUsed(clientset kubernetes.Clientset) ClusterResources {
	clusterResouces := ClusterResources{}
	resources := GetPods("", clientset)
	var cpu int64
	var memory int64
	for _, item := range resources {
		containers := item.Spec.Containers
		for _, container := range containers {
			cpu += container.Resources.Limits.Cpu().Value()
			memory += container.Resources.Limits.Memory().Value()
		}
	}
	clusterResouces.Services = GetServiceNumber(clientset, "")
	clusterResouces.UsedCpu = cpu
	clusterResouces.UsedMem = memory / 1024 / 1024 / 1024
	return clusterResouces
}
