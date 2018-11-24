package k8s

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"github.com/astaxie/beego/logs"
	"strings"
	"strconv"
	"fmt"
	"encoding/json"
	"cloud/util"
	"cloud/sql"
)

// 创建配置文件
// @param name
// 2018-01-17 16:18
// data := app.GetConfgData("adf","asdfasdfdasf")
// k8s.CreateConfigmap("10.16.55.6","8080",data, "adf","test-app--dfsad")
func CreateConfigmap(param ServiceParam) {
	if len(param.ConfigureData) == 0 {
		return
	}

	namespace := param.Namespace
	for _, v := range param.ConfigureData {
		conf := map[string]interface{}{
			"apiVersion": "v1",
			"kind": "ConfigMap",
			"data": v.ConfigDbData,
			"metadata": map[string]interface{}{
				"name": v.DataName,
			},
		}
		resource := &metav1.APIResource{Name: "ConfigMaps", Namespaced: true}
		obj := unstructured.Unstructured{Object: conf}

		var update bool
		r, err := param.Cl2.Resource(resource, namespace).Get(v.DataName, metav1.GetOptions{})
		if r.GetName() != "" {
			update = true
		}
		if param.NoUpdateConfig == true && update {
			logs.Error("配置不允许更新配置", v.DataName)
			continue
		}
		//var d *unstructured.Unstructured
		if update {
			_, err = param.Cl2.Resource(resource, namespace).Update(&obj)
			logs.Info("更新 configmap", namespace, err, v.DataName, v)
		} else {
			_, err = param.Cl2.Resource(resource, namespace).Create(&obj)
			logs.Info("创建 configmap", namespace,  err, v.DataName)
		}
	}
}


// 生成容器环境变量数据
// 2018-01-11 21:01
func getEnv(envs string) []map[string]interface{} {
	data := make([]map[string]interface{}, 0)
	if len(envs) < 3 {
		return data
	}
	envData := strings.Split(envs, "\n")
	for _, v := range envData {
		vs := strings.Split(v, "=")
		if len(vs) < 2 {
			continue
		}
		if string(vs[0][0]) == "#" ||  string(vs[0][0]) == " "  {
			logs.Info("环境变量获取到注释", vs[0])
			continue
		}
		temp := map[string]interface{}{
			"name":  vs[0],
			"value": strings.Join(vs[1:], "="),
		}
		data = append(data, temp)
	}
	return data
}


// 2018-01-17 17:48
// 获取configmap 挂载某个key的配置
func getConfigKey(keys string) []map[string]interface{} {
	mapdata := make([]map[string]interface{}, 0)
	allkey := strings.Split(keys, ",")
	for _, v := range allkey {
		temp := map[string]interface{}{
			"key":  v,
			"path": v,
		}
		mapdata = append(mapdata, temp)
	}
	return mapdata
}

// 加工磁盘卷数据
// 2018-09-30 20:57 -6
func getFilebeatVolumes(storageData string) ([]map[string]interface{}, []map[string]interface{}) {
	if storageData == ""{
		storageData = `[]`
	}
	storages := make([]map[string]interface{}, 0)
	voluments := make([]map[string]interface{}, 0)
	data := make([]StorageData, 0)
	err := json.Unmarshal([]byte(storageData), &data)
	if err != nil {
		logs.Error("处理Volumes失败", err)
		return storages, voluments
	}

	for k, p := range data {
		data := make(map[string]interface{}, 0)
		is := false

		v := p.HostPath
		if v[len(v)-1:len(v)] != "/" {
			continue
		}
		// 使用物理机的存储
		if len(p.HostPath) > 0 && strings.Contains(p.HostPath, "/") {
			data = map[string]interface{}{
				"name": "volume-filebeat-" + strconv.Itoa(k),
				"emptyDir": map[string]interface{}{
				},
			}
			is = true
		}

		if ! is {
			continue
		}
		volumeMountsData := map[string]interface{}{
			"name":      "volume-filebeat-" + strconv.Itoa(k),
			"mountPath": p.ContainerPath,
		}
		storages = append(storages, data)
		voluments = append(voluments, volumeMountsData)
	}

	fmt.Println("filebeat-", storages, voluments)
	return storages, voluments
}

// 加工磁盘卷数据
// 2018-01-11 13::57
func getVolumes(storagesData string, configData []ConfigureData, param ServiceParam) ([]map[string]interface{}, []map[string]interface{}) {
	if storagesData == ""{
		storagesData = `[]`
	}
	storages := make([]map[string]interface{}, 0)
	voluments := make([]map[string]interface{}, 0)
	data := make([]StorageData, 0)
	err := json.Unmarshal([]byte(storagesData), &data)
	if err != nil {
		logs.Error("处理Volumes失败", err)
		return storages, voluments
	}

	data = append(data,  StorageData{ContainerPath:"/etc/localtime",ReadOnly:true, Volume:"", HostPath:"/etc/localtime"})
	data = append(data,  StorageData{ContainerPath:"/dev/random",ReadOnly:true, Volume:"", HostPath:"/dev/urandom"})

	for k, p := range data {
		data := make(map[string]interface{}, 0)
		is := false
		// 使用物理机的存储
		if len(p.HostPath) > 0 && strings.Contains(p.HostPath, "/") {
			data = map[string]interface{}{
				"name": "volume-" + strconv.Itoa(k),
				"hostPath": map[string]interface{}{
					"path": p.HostPath,
				},
			}
			if p.Model != 0 {
				data["hostPath"].(map[string]interface{})["defaultMode"] = p.Model
			}
			is = true
		}

		// 使用k8s pv的存
		if len(p.Volume) > 3 {
			data = map[string]interface{}{
				"name": "volume-" + strconv.Itoa(k),
				"persistentVolumeClaim": map[string]interface{}{
					"claimName": p.Volume,
				},
			}
			CreateServicePvc(param, p.Volume, p.ContainerPath)
			is = true
		}

		fmt.Println("ddddddddd", data)
		if ! is {
			continue
		}
		volumeMountsData := map[string]interface{}{
			"name":      "volume-" + strconv.Itoa(k),
			"mountPath": p.ContainerPath,
		}
		if p.ReadOnly {
			volumeMountsData["readOnly"] = true
		}
		storages = append(storages, data)
		voluments = append(voluments, volumeMountsData)
	}

	if len(configData) > 0 {
		for id, conf := range configData {
			mapdata := make(map[string]interface{})
			if conf.DataId != "" {
				// 挂载单个key的
				paths := getConfigKey(conf.DataId)
				mapdata = map[string]interface{}{
					"configMap": map[string]interface{}{
						"defaultMode": 420,
						"name":        conf.DataName,
						"items":       paths,
					},
					"name": "configmap-volume-"+ strconv.Itoa(id),
				}

				for _, v := range paths {
					mountPath := conf.ContainerPath + "/" + v["path"].(string)
					mountData := map[string]interface{}{
						"name":      "configmap-volume-" + strconv.Itoa(id),
						"mountPath": mountPath,
						"subPath":   v["path"],
					}
					WriteMountDataToDb(conf.DataName, v["path"].(string), param.ClusterName, param.Namespace, mountPath, param.Name)
					voluments = append(voluments, mountData)
				}

			} else {
				// 挂载整个组
				mapdata = map[string]interface{}{
					"configMap": map[string]interface{}{
						"defaultMode": 420,
						"name":        conf.DataName,
					},
					"name": "configmap-volume-"+strconv.Itoa(id),
				}
				mountData := map[string]interface{}{
					"name":      "configmap-volume-"+strconv.Itoa(id),
					"mountPath": conf.ContainerPath,
				}
				if conf.ContainerPath == "/etc/filebeat/" {
					mountData["name"] = "filebeat-config-" + param.ServiceName
				}
				WriteMountDataToDb(conf.DataName, "", param.ClusterName, param.Namespace, conf.ContainerPath, param.Name)
				voluments = append(voluments, mountData)
			}
			storages = append(storages, mapdata)
		}
	}
	fmt.Println(storages, voluments)
	return storages, voluments
}

// 2018-01-18 11:31
// 将挂载数据写入到数据
func WriteMountDataToDb(configname string, dataName string, cluster string, namespace string, mountpath string, serviceName string)  {
	if cluster == "" || serviceName == "" {
		return
	}
	// 首先查询已经mount的数据
	mountData := CloudConfigureMount{}
	mountData.CreateTime = util.GetDate()
	mountData.LastUpdateTime = util.GetDate()
	mountData.DataName = dataName
	mountData.ClusterName = cluster
	mountData.Namespace = namespace
	mountData.MountPath = mountpath
	mountData.ConfigureName = configname
	mountData.ServiceName = serviceName
	searchMap := sql.GetSearchMapV("ClusterName", cluster, "ConfigureName", configname,"ServiceName", serviceName, "Namespace", namespace)
	if dataName != "" {
		searchMap.Put("DataName", dataName)
	}
	mounts := make([]CloudConfigureMount, 0)
	sql.Raw(sql.SearchSql(mountData, SelectCloudConfigureMount, searchMap)).QueryRows(&mounts)
	var action string
	if len(mounts) > 0 {
		action = sql.UpdateSql(mountData, UpdateCloudConfigureMount, searchMap, "CreateTime,MountId")
	}else{
		action = sql.InsertSql(mountData, InsertCloudConfigureMount)
	}
	sql.Raw(action).Exec()
}