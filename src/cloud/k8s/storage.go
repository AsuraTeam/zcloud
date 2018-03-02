package k8s

import (
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	v12 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"github.com/astaxie/beego/logs"
	"fmt"
	"cloud/util"
	"cloud/sql"
)

// 2018-02-22 10:37
// 创建pvc
func createPvc(clustername string, pvc map[string]interface{}, namespace string) error {
	pvcClient, err := GetYamlClient(clustername, "", "v1", "/api")
	resource := &v12.APIResource{Name: "PersistentVolumeClaims", Namespaced: true}
	pvcobj := unstructured.Unstructured{Object: pvc}
	d, err := pvcClient.Resource(resource, namespace).Create(&pvcobj)
	if err != nil {
		logs.Error("创建PVC失败", d, err)
	}
	return err
}

// 2018-01-30 17:37
// 获取serviceaccount信息
func GetServiceAccount(name string, param StorageParam) string {
	cl, err := GetClient(param.ClusterName)
	if err != nil {
		logs.Error("获取ServiceAccount失败", err)
		return ""
	}
	account := cl.CoreV1().ServiceAccounts(param.Namespace)
	info, err := account.Get(name, v12.GetOptions{})
	secrets := info.Secrets
	if len(secrets) > 0 {
		return secrets[0].Name
	}
	fmt.Println(util.ObjToString(info), util.ObjToString(err))
	return ""
}

// 创建ServiceAccount
// 2018-01-30 17:28
func createServiceAccount(param StorageParam) {
	conf := map[string]interface{}{
		"apiVersion": "v1",
		"kind":       "ServiceAccount",
		"metadata": map[string]interface{}{
			"name":      "default",
			"namespace": param.Namespace,
		},
	}
	pvcClient, err := GetYamlClient(param.ClusterName, "", "v1", "/api")
	resource := &v12.APIResource{Name: "ServiceAccounts", Namespaced: true}
	accountobj := unstructured.Unstructured{Object: conf}
	d, err := pvcClient.Resource(resource, param.Namespace).Create(&accountobj)
	if err != nil {
		logs.Error("创建ServiceAccounts失败", d, err)
	}
}


// 2018-02-22 09:59
// 创建StorageClass
func createStorageClass(clustername string , class map[string]interface{}) error {
	classClient, err := GetYamlClient(clustername, "storage.k8s.io", "v1beta1", "/apis")
	resource := &v12.APIResource{Name: "StorageClasses", Namespaced: false}
	classObj := unstructured.Unstructured{Object: class}
	d, err := classClient.Resource(resource, "").Create(&classObj)
	if err != nil {
		logs.Error("创建StorageClass失败", d, err)
	}
	return err
}


// 2018-01-30 21:07
// 创建pvc
func CreateServicePvc(param ServiceParam, pvcName string, containerPath string) error {
	orm := sql.GetOrm()
	cs := CloudStorage{}
	q := "select storage_type,storage_size from cloud_storage where cluster_name=? and name= ?"
	orm.Raw(q, param.ClusterName, pvcName).QueryRow(&cs)
	if param.AccessMode == "" {
		param.AccessMode = "ReadWriteMany"
	}

	storageParam := StorageParam{
		PvcName:     pvcName,
		Master:      param.Master,
		Port:        param.MasterPort,
		Size:        cs.StorageSize,
		Namespace:   param.Namespace,
		AccessMode:  param.AccessMode,
		StorageType: cs.StorageType,
		ClusterName:param.ClusterName,
	}

	err := pvcCreate(storageParam)
	if err != nil {
		return err
	}

	mount := CloudStorageMountInfo{
		ClusterName:  param.ClusterName,
		StorageType:  cs.StorageType,
		ServiceName:  param.Name,
		StorageName:  pvcName,
		ResourceName: param.ResourceName,
		AppName:      param.AppName,
		CreateTime:   util.GetDate(),
		CreateUser:   param.CreateUser,
		MountPath:    containerPath,
		Status:       "1",
	}
	q = sql.InsertSql(mount, "insert into cloud_storage_mount_info")
	orm.Raw(q).Exec()
	return err
}

// 2018-02-22 13:32
// 创建pvc
func pvcCreate(param StorageParam) error {
	logs.Info("创建pvc", util.ObjToString(param))
	var err error
	switch param.StorageType {
	case "Nfs":
		createNfsStorageClass(param)
		err = createNfsPvc(param)
		break
	case "Glusterfs":
		createGlusterfsStorageClass(param.ClusterName)
		createGlusterfsPvc(param)
		break
	case "Ceph":
		break
	default:
		break
	}

	if err != nil {
		logs.Error("创建pvc失败", err)
		return err
	}
	return nil
}

// 2018-01-31 10:35
// 删除存储卷
func DeletePvc(param StorageParam) error {
	pvcClient, err := GetYamlClient(param.ClusterName, "", "v1", "/api")
	resource := &v12.APIResource{Name: "PersistentVolumeClaims", Namespaced: true}
	err = pvcClient.Resource(resource, param.Namespace).Delete(param.PvcName, &v12.DeleteOptions{})
	if err != nil {
		logs.Error("删除pvc失败", err)
	}
	return err
}
