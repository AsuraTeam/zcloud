package k8s

import (
	"cloud/util"
	"time"
	"strings"
)

// nginx配置文件路径
const (
	NginxUpstreamPath = "/usr/local/nginx/conf/vhosts/upstream"
	NginxConfigPath   = "/usr/local/nginx/conf/vhosts/conf"
	NginxSslPath      = "/usr/local/nginx/conf/vhosts/ssl"
	LbNginxConfig     = "lb-nginx-config"
	LbNginxUpstream   = "lb-nginx-upstream"
	LbNginxSsl        = "lb-nginx-ssl"
	LbNginxStartPath        = "/start/"
	LbNginxDaemonPath        = "/daemon/"
)

// 获取默认的配置
// 2018-02-02 16:13
func getNgxinDefaulgConfig(containerPath string, dataname string, configData map[string]interface{}, confType string) ConfigureData {
	if configData == nil {
		configData = map[string]interface{}{"default": ""}
	}
	return ConfigureData{
		ContainerPath: containerPath,
		DataName:      dataname + confType,
		ConfigDbData:  configData,
	}
}

// 2018-02-03 06:30
// 获取nginx计划配置信息
func nginxTestJobParam(master string, port string) JobParam {
	param := JobParam{
		Master:            master,
		Port:              port,
		Timeout:           50,
		Memory:            40,
		Jobname:           "job-" + util.Md5Uuid(),
		Cpu:               1,
		Namespace:         util.Namespace("lb", "nginx"),
		Images:            "nginx:v1",
		ConfigureData:     getNginxDefaultConf("-test"),
		NoUpdateConfigMap: true,
		Command:           []string{"sh", "/opt/check.sh"},
	}
	return param
}

// 2018-02-03 07:18
// 分析日志,获取执行结果
func getNginxJobLog(logstr string) string {
	logStr := strings.Split(logstr, "\n")
	r := make([]string, 0)
	for _, v := range logStr {
		if strings.Contains(v, "nginx:") {
			r = append(r, v)
		}
	}
	return strings.Join(r, "\n")
}


// 创建测试任务,检查nginx配置文件
// 2018-02-02 20:30
func MakeTestJob(master string, port string, clusterName string) (string, int64) {
	start := time.Now().Unix()
	param := nginxTestJobParam(master, port)
	param.ClusterName = clusterName
	r := CreateJob(param)
	param.Jobname = r
	logStr := getJobResult(param, "nginx:", 8, "nginx")
	times := time.Now().Unix() - start
	return logStr, times
}

// 2018-02-02 21:48
// 获取默认需要挂载的配置
func getNginxDefaultConf(confType string) []ConfigureData {
	nginxConfigMap := make([]ConfigureData, 0)
	nginxConfigMap = append(nginxConfigMap, getNgxinDefaulgConfig(NginxSslPath, LbNginxSsl+confType, nil, ""))
	nginxConfigMap = append(nginxConfigMap, getNgxinDefaulgConfig(NginxConfigPath, LbNginxConfig+confType, nil, ""))
	nginxConfigMap = append(nginxConfigMap, getNgxinDefaulgConfig(NginxUpstreamPath, LbNginxUpstream+confType, nil, ""))

	// 检查命令
	configData := `/usr/local/nginx/sbin/nginx -t`
	conf :=  ConfigureData{
		ContainerPath: "/opt/",
		DataName:      "check.sh",
		ConfigDbData:   map[string]interface{}{"check.sh": configData},
	}
	nginxConfigMap = append(nginxConfigMap, conf)
	cm := map[string]interface{}{"reload.sh": reloadNginx}
	nginxConfigMap = append(nginxConfigMap, getNgxinDefaulgConfig(LbNginxStartPath, "reload.sh", cm, ""))
	daemon := map[string]interface{}{"daemon.sh": "sh -c 'sh /start/reload.sh'"}
	nginxConfigMap = append(nginxConfigMap, getNgxinDefaulgConfig(LbNginxDaemonPath, "daemon.sh", daemon, ""))
	return nginxConfigMap
}

// 创建nginx集群容器
func CreateNginxLb(param ServiceParam) {
	param.Port = "80,443"
	param.Image = "nginx:v1"
	param.Namespace = "lb--nginx"
	param.Name = "nginx-lb"
	param.Cpu = 1
	param.Memory = "4096"
	param.HostPort = "80,443"
	param.NoUpdateConfig = true
	param.Command = `["sh","/daemon/daemon.sh"]`

	clientSet, _ := GetClient(param.ClusterName)
	cl2, _ := GetYamlClient(param.ClusterName, "", "v1", "api")
	param.Cl2 = cl2
	param.Cl3 = clientSet
	param.ConfigureData = getNginxDefaultConf("")
	YamlCreateNamespace(param.ClusterName, param.Namespace)
	CreateConfigmap(param)
	CreateDeamonSet(param)
}
