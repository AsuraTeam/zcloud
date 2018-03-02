package k8s

import (
	"k8s.io/client-go/kubernetes"
	"encoding/json"
	"k8s.io/api/core/v1"
	"github.com/astaxie/beego/logs"
	"strings"
	"cloud/util"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const dockerPullJson = `{
        "auths": {
                "REGISTRY": {
                        "auth": "AUTH"
                }
        },
        "HttpHeaders": {
                "User-Agent": "Docker-Client/18.01.0-ce (linux)"
        }
}`

// 2018-02-04 21:06
func GetDockerImagePullName(name string) string {
	name = strings.Replace(name, ".", "-", -1)
	name = strings.Replace(name, ":", "", -1)
	return name
}

// 2018-02-04 20:51
// 为私有镜像仓库拉取镜像添加拉取权限
func CreateImagePullSecret(param ServiceParam) {
	logs.Info("CreateImagePullSecret", util.ObjToString(param))
	name := GetDockerImagePullName(param.Registry)
	cfg := strings.Replace(dockerPullJson, "REGISTRY", param.Registry, -1)
	cfg = strings.Replace(cfg, "AUTH", util.Base64Encoding(param.RegistryAuth), -1)
	obj := map[string]interface{}{
		"apiVersion": "v1",
		"data": map[string]interface{}{
			".dockerconfigjson": util.Base64Encoding(cfg),
		},
		"kind": "Secret",
		"metadata": map[string]interface{}{
			"name": name,
		},
		"type": "kubernetes.io/dockerconfigjson",
	}
	temp, _ := json.Marshal(obj)
	secret := v1.Secret{}
	json.Unmarshal(temp, &secret)
	isExists := SecretIsExists(param.Cl3, param.Namespace, name)
	var err error
	d := &v1.Secret{}
	if isExists{
		d, err = param.Cl3.CoreV1().Secrets(param.Namespace).Update(&secret)
	}else {
		d, err = param.Cl3.CoreV1().Secrets(param.Namespace).Create(&secret)
	}
	if err != nil {
		logs.Error("创建ImagePull Secret错误", err)
	}
	logs.Info("创建ImagePullSecret", d, err)
}


// 2018-02-04 21:32
// 获取安全密码是否存在
func SecretIsExists(client kubernetes.Clientset, namespace string, name string) bool {
	_, err := client.CoreV1().Secrets(namespace).Get(name, v12.GetOptions{})
	if err == nil {
		return true
	}
	return false
}

// 2018-02-09 15:22
// 删除secret
func DeleteSecret(client kubernetes.Clientset, namespace string)  {
	secrets,err := client.CoreV1().Secrets(namespace).List(v12.ListOptions{})
	if err !=  nil {
		logs.Error("删除Secret错误", err, namespace)
		return
	}
	item := secrets.Items
	for _, v := range item{
		logs.Info("删除Secret", v.Namespace, v.Name)
		client.CoreV1().Secrets(namespace).Delete(v.Name, &v12.DeleteOptions{})
	}
}