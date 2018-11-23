package k8s

import (
	"strings"
	"cloud/util"
	"github.com/astaxie/beego/logs"
	"path/filepath"
	"github.com/astaxie/beego"
)

// 2018-03-30 13:36
// 配置kafka
func getFilebeatOutput(param ServiceParam) string {
	template := `
output.kafka:
  enable: True
  hosts: HOSTS
  topic: 'k8s-%{[fields][appname]}'
  partition.round_robin:
    reachable_only: false
  #required_acks: 1
  #compression: gzip
  #max_message_bytes: 1000000`
	topic := beego.AppConfig.String("kafka.topic")
	if topic != "" {
		template = strings.Replace(template, "k8s-%{[fields][appname]}", topic, -1)
	}
	if len(param.Kafka) > 0 && len(param.LogPath) > 0 || (len(param.ElasticSearch) > 0 && len(param.LogPath) > 0) {
		addrs := strings.Split(param.Kafka, ",")
		// ["kafka-node-01:9092"]
		host := util.ObjToString(addrs)
		host = strings.Replace(host, "\\\"", "\"", -1)
		logs.Info(host)
		template = strings.Replace(template, "HOSTS", host, -1)
		if len(param.ElasticSearch) > 0 {
			template = strings.Replace(template, "output.kafka", "output.elasticsearch", -1)
		}
		return template
	}
	return ""
}

func filebeatCmd() {
	cmd := `#!/bin/bash
cd /etc/filebeat-1
export PATH=$PATH:/usr/local/share/filebeat/bin/
parm=$(echo $* | sed 's#filebeat/filebeat.yml#filebeat-1/temp.yml#g')
cp /etc/filebeat/filebeat.yml /etc/filebeat-1/temp.yml
ip=$(ifconfig |awk '$0 ~ /broadcast/ {print $2}')
sed -i "s/IP_ADDRESS/$ip/g" /etc/filebeat-1/temp.yml
cat /etc/filebeat-1/temp.yml
echo $parm
filebeat-1 -e -c  /etc/filebeat-1/temp.yml`
	logs.Info(cmd)
}

// 检查日志路径是否包含文件
// 如果不包含文件那么只有目录
func getFilebeatPaths(param ServiceParam) map[string][]string {
	paths := strings.Split(param.LogPath, "\n")
	dir := map[string][]string{}
	//file := map[string]string{}
	for _, v := range paths {
		logs.Info(v[len(v)-1:len(v)] == "/", v, v[len(v)-1:len(v)])
		if v[len(v)-1:len(v)] == "/" {
			dir[v] = []string{}
		} else {
			p := strings.Split(v, "/")
			r := strings.Join(p[0:len(p)-1], "/") + "/"
			dir[r] = []string{}
		}
	}
	for _, v := range paths {
		if v[len(v)-1:len(v)] == "/" {
			continue
		}
		dirname := strings.Replace(filepath.Dir(v)+"/", "\\", "/", -1)
		if _, ok := dir[dirname]; ok {
			dir[dirname] = append(dir[dirname], v)
		}
	}
	return dir
}

// 2018-10-04 15:24
// 获取所有需要挂载的目录
func resetFilebeatPath(param ServiceParam) []string {
	path := make([]string, 0)
	paths := strings.Split(param.LogPath, "\n")
	for _, v := range paths {
		//if v[len(v)-1:len(v)] == "/" {
		//	continue
		//}
		p := strings.Split(v, "/")
		r := strings.Join(p[0:len(p)-1], "/") + "/"
		if !util.ListExistsString(path, r) {
			path = append(path, r)
		}
	}
	return path
}

// 2018-03-30 14:01
// 获取日志路径
func getLogPaths(param ServiceParam) string {
	if len(param.LogPath) == 0 {
		return ""
	}
	paths := getFilebeatPaths(param)
	path := make([]string, 0)
	counter := 0
	logs.Info("获取到日志文件", util.ObjToString(paths))
	for dir, file := range paths {
		logs.Info(dir, file)
		if len(file) > 0 {
			for _, v := range file {
				if counter == 0 {
					path = append(path, "- "+v)
				} else {
					path = append(path, "    - "+v)
				}
			}
		} else {
			if counter == 0 {
				path = append(path, "- "+dir+"*")
			} else {
				path = append(path, "    - "+dir+"*")
			}
		}
		counter += 1

	}
	return strings.Join(path, "\n")
}

// 生产
func CreateFilebeatConfig(param ServiceParam) map[string]interface{} {
	if len(param.Kafka) == 0 || len(param.LogPath) == 0 {
		return map[string]interface{}{}
	}

	temp := `
filebeat.prospectors:
- input_type: log
  paths:
    PATHS 
  document_type: "k8s-APP_NAME"
  fields:
    runtime_env: RUN_TIME_ENV 
    appname: APP_NAME
    ip_address: IP_ADDRESS

  include_lines: ^(\[201|201)
  multiline.pattern: ^(\[201|201)
  multiline.negate: true 
  multiline.match: after

registry_file: /dev/shm/registry
#output.logstash:
#  hosts: ["logstash-02:5514"]
KAFKA
`

	temp = strings.Replace(temp, "KAFKA", getFilebeatOutput(param), -1)
	item := strings.Replace(param.ServiceName, "--1", "", -1)
	temp = strings.Replace(temp, "PATHS", getLogPaths(param), -1)
	temp = strings.Replace(temp, "$item", item, -1)
	temp = strings.Replace(temp, "APP_NAME", param.ServiceName, -1)
	temp = strings.Replace(temp, "RUN_TIME_ENV", param.Ent, -1)
	conf := map[string]interface{}{
		"filebeat.yml": strings.Replace(temp, "\\\"", "\"", -1),
	}
	confData := ConfigureData{}
	confData.DataName = "filebeat-config-" + param.ServiceName
	confData.ConfigDbData = conf
	confData.ContainerPath = "/etc/filebeat/"
	if len(param.ConfigureData) == 0 {
		param.ConfigureData = make([]ConfigureData, 0)
	}
	param.ConfigureData = append(param.ConfigureData, confData)
	CreateConfigmap(param)
	return getFilebeatContainer(param)
}

// filebeat
func getFilebeatConfig(param ServiceParam) map[string]interface{} {
	v := map[string]interface{}{
		"name": "filebeat-config-" + param.ServiceName,
		"configMap": map[string]interface{}{
			"name": "filebeat-config-" + param.ServiceName,
		},
	}
	return v
}

// 2018-09-30 20:37 -6
// filebeat容器启动
func getFilebeatContainer(param ServiceParam) map[string]interface{} {
	image := beego.AppConfig.String("filebeat.image")
	if len(image) == 0 {
		image = "prima/filebeat"
	}
	_, volumeMounts := getVolumes(param.StorageData, param.ConfigureData, param)
	data := map[string]interface{}{
		"image":           "" + image + "",
		"name":            "filebeat",
		"imagePullPolicy": "Always",
		"volumeMounts":    volumeMounts,
		"command":         strings.Split("filebeat,-e,-c,/etc/filebeat/filebeat.yml", ","),
	}
	return data
}

// 2018-09-30 20:56 -6
// 获取filebeat挂载的路径信息
func getFilebeatStorage(param ServiceParam) ([]map[string]interface{}, []map[string]interface{}) {
	data := make([]StorageData, 0)
	//paths := strings.Split(param.LogPath, "\n")
	paths := resetFilebeatPath(param)
	for _, v := range paths {
		data = append(data, StorageData{ContainerPath: v, HostPath: v})
	}

	logs.Info("获取到需要挂载的路径", util.ObjToString(data))
	return getFilebeatVolumes(util.ObjToString(data))
}
