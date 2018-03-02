package k8s

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"github.com/astaxie/beego/logs"
)

// 创建daemonset的服务
func CreateDeamonSet(param ServiceParam) {
	c1, err := GetYamlClient(param.ClusterName, "extensions", "v1beta1", "/apis")
	if err != nil {
		logs.Error("创建DaemonSet失败", err)
		return
	}
	volumes, volumeMounts := getVolumes(param.StorageData, param.ConfigureData, param)
	conf := map[string]interface{}{
		"apiVersion": "extensions/v1beta1",
		"kind":       "DaemonSet",
		"metadata": map[string]interface{}{
			"labels": map[string]interface{}{
				"k8s-app": param.Name,
			},
			"name":      param.Name,
			"namespace": param.Namespace,
		}, "spec": map[string]interface{}{
			"selector": map[string]interface{}{
				"matchLabels": map[string]interface{}{
					"name": param.Name,
				},
			},
			"template": map[string]interface{}{
				"metadata": map[string]interface{}{
					"labels": map[string]interface{}{
						"name": param.Name,
					},
				},
				"spec": map[string]interface{}{
					"containers": []map[string]interface{}{
						map[string]interface{}{
							"image": param.Image,
							"name":  param.Name,
							"ports": getPorts(param.Port, param.HostPort),
							"resources": map[string]interface{}{
								"limits": map[string]interface{}{
									"memory": param.Memory + "Mi",
									"cpu":    param.Cpu,
								},
								"requests": map[string]interface{}{
									"memory": param.Memory + "Mi",
									"cpu":    param.Cpu,
								},
							},
							"volumeMounts": volumeMounts,
						},
					},
					"volumes": volumes,
				},
			},
		},
	}

	conf = setPrivileged(param, conf)
	conf = setContainerCommand(param, conf)
	resource := &metav1.APIResource{Name: "DaemonSets", Namespaced: true}
	obj := unstructured.Unstructured{Object: conf}
	var d *unstructured.Unstructured
	_, err = c1.Resource(resource, param.Namespace).Get(param.Name, metav1.GetOptions{})

	if err != nil {
		d, err = c1.Resource(resource, param.Namespace).Create(&obj)
		logs.Info("创建DaemonSet", d, err)
	} else {
		d, err = c1.Resource(resource, param.Namespace).Update(&obj)
		logs.Info("更新DaemonSet", d, err)
	}

	if err != nil {
		logs.Error("创建DaemonSets失败", d, err)
	}
}
