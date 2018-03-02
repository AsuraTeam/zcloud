package k8s

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/api/core/v1"
	"github.com/astaxie/beego/logs"
)

// 2018-02-11
// 创建默认serviceAccount
func CreateServiceAccount(client kubernetes.Clientset, namespace string, name string) {
	account := v1.ServiceAccount{}
	account.Name = name
	d, err := client.CoreV1().ServiceAccounts(namespace).Create(&account)
	if err != nil {
		logs.Error("创建serviceAccount失败", err, d)
	}
}
