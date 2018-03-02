package k8s

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"github.com/astaxie/beego/logs"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"cloud/controllers/base/quota"
)

// 设置默认资源使用
func SetLimts(clustername string, namespace string, mem string ,cpu interface{}, action string) error{
	cl, err := GetYamlClient(clustername, "", "v1", "/api")
	conf := map[string]interface{}{
		"apiVersion": "v1",
		"kind":"LimitRange",
		"metadata": map[string]interface{}{
			"name": "limit-range-global",
		},
		"spec": map[string]interface{}{
			"limits": []map[string]interface{}{
				map[string]interface{}{
					"max": map[string]interface{}{
						"cpu":    cpu,
						"memory": mem,
					},
					"min": map[string]interface{}{
						"cpu":    cpu,
						"memory": mem,
					},
					"type": "Pod",
				},
				map[string]interface{}{
					"default": map[string]interface{}{
						"cpu":    cpu,
						"memory": mem,
					},
					"max": map[string]interface{}{
						"cpu":    cpu,
						"memory": mem,
					}, "min": map[string]interface{}{
						"cpu":    cpu,
						"memory": mem,
					}, "type": "Container",
				},
			},
		},
	}
	resource := &v12.APIResource{Name: "limitranges", Namespaced: true}
	podobj := unstructured.Unstructured{Object: conf}
	var d interface{}
	if action == "update"{
		d, err = cl.Resource(resource, namespace).Update(&podobj)
	}else{
		d, err = cl.Resource(resource, namespace).Create(&podobj)
	}
	if err != nil {
		logs.Error("创建配额失败", err)
		return err
	}
	logs.Info("创建配额成功", d)
	return err
}

// 2018-02-11 21:59
// 检查配置服务时配额是否够用
// 检查资源配额是否够用
func CheckQuota(username string, podNumber int64, cpu int64, memory int64, resourceName string) (bool, string) {
	quotaDatas := quota.GetUserQuotaData(username, resourceName)
	for _, v := range quotaDatas {
		logs.Info(v.MemoryUsed*1024, v.QuotaMemory*1024, memory)
		if v.MemoryUsed*1024+memory > v.QuotaMemory*1024 {
			return false, "内存超过配额限制"
		}
		if v.CpuUsed+cpu > v.QuotaCpu {
			return false, "cpu超过配额限制"
		}
		if v.PodUsed+podNumber > v.PodNumber {
			return false, "容器超过配额限制"
		}
	}
	return true, ""
}