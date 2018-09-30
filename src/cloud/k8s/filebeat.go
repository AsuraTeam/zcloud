package k8s

import (
	"strings"
	"cloud/util"
	"github.com/astaxie/beego/logs"
)

// 2018-03-30 13:36
// 配置kafka
func getKafka(param ServiceParam) string  {
	kafka := `
output.kafka:
  enable: True
  hosts: HOSTS
  topic: 'access'
  partition.round_robin:
    reachable_only: false
  #required_acks: 1
  #compression: gzip
  #max_message_bytes: 1000000`
  if len(param.Kafka) > 0 && len(param.LogPath) > 0 {
  	addrs := strings.Split(param.Kafka, ",")
  	// ["kafka-node-01:9092"]
  	host := util.ObjToString(addrs)
  	host = strings.Replace(host, "\\\"", "\"", -1)
  	logs.Info(host)
  	kafka = strings.Replace(kafka, "HOSTS", host , -1)
  	return kafka
  }
	return ""
}

// 2018-03-30 14:01
// 获取日志路径
func getLogPaths(param ServiceParam)  string {
	if len(param.LogPath) ==0 {
		return ""
	}
	paths := strings.Split(param.LogPath, ",")
	path := make([]string, 0)
	for _, v := range paths {
		path = append(path, "- " + v)
	}
	return strings.Join(path, "\n")
}

// 生产
func CreateFilebeatConfig(param ServiceParam)  map[string]interface{} {
	if len(param.Kafka) == 0 || len(param.LogPath) == 0 {
		return  map[string]interface{}{}
	}
	temp := `
filebeat.prospectors:
- input_type: log
  paths:
    PATHS 
  document_type: "java_access-APP_NAME"
  fields:
    runtime_env: RUN_TIME_ENV 
    appname: APP_NAME


  include_lines: ^(\[201|201)
  multiline.pattern: ^(\[201|201)
  multiline.negate: true 
  multiline.match: after

registry_file: /dev/shm/registry
#output.logstash:
#  hosts: ["logstash-02:5514"]
KAFKA
`

	temp = strings.Replace(temp, "KAFKA", getKafka(param), -1)
	temp = strings.Replace(temp, "PATHS", getLogPaths(param), -1)
	temp = strings.Replace(temp, "APP_NAME", param.ServiceName, -1)
	temp = strings.Replace(temp, "RUN_TIME_ENV", param.Ent, -1)
	conf := map[string]interface{}{
		"filebeat.yml":strings.Replace(temp, "\\\"", "\"", -1),
	}
	confData := ConfigureData{}
	confData.DataName = "filebeat-config-"+ param.ServiceName
	confData.ConfigDbData = conf
	confData.ContainerPath = "/etc/filebeat/"
	if len(param.ConfigureData) == 0 {
		param.ConfigureData = make([]ConfigureData, 0)
	}
	param.ConfigureData = append(param.ConfigureData, confData)
	CreateConfigmap(param)
	return getFilebeatContainer(param)
}

// 2018-09-30 20:37 -6
// filebeat容器启动
func getFilebeatContainer(param ServiceParam)  map[string]interface{}{
	_, volumeMounts := getVolumes(param.StorageData, param.ConfigureData, param)
	data := map[string]interface{}{
		"image": "prima/filebeat",
		"name":  "filebeat",
		"imagePullPolicy": "Always",
		"volumeMounts": volumeMounts,
	}
	return data
}

// 2018-09-30 20:56 -6
// 获取filebeat挂载的路径信息
func getFilebeatStorage(param ServiceParam) ([]map[string]interface{},[]map[string]interface{} ) {
	data := make([]StorageData, 0)
	paths := strings.Split(param.LogPath, ",")
	for _, v := range paths{
		data = append(data,  StorageData{ContainerPath:v, HostPath:v})
	}
	return getFilebeatVolumes(util.ObjToString(data))
}