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
)

// 获取默认的配置
// 2018-02-02 16:13
func getNgxinDefaulgConfig(containerPath string, dataname string, configdata map[string]interface{}, conftype string) ConfigureData {
	if configdata == nil {
		configdata = map[string]interface{}{"default": ""}
	}
	return ConfigureData{
		ContainerPath: containerPath,
		DataName:      dataname + conftype,
		ConfigDbData:  configdata,
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
		Cpu:               1,
		Namespace:         util.Namespace("lb", "nginx"),
		Images:            "nginx:v1",
		ConfigureData:     getNginxDefaultConf("-test"),
		NoUpdateConfigMap: true,
		Command:           []string{"sh", "/check.sh"},
	}
	return param
}

// 2018-02-03 07:18
// 分析日志,获取执行结果
func getNginxJobLog(logstr string) string {
	logstrs := strings.Split(logstr, "\n")
	r := make([]string, 0)
	for _, v := range logstrs {
		if strings.Contains(v, "nginx:") {
			r = append(r, v)
		}
	}
	return strings.Join(r, "\n")
}


// 创建测试任务,检查nginx配置文件
// 2018-02-02 20:30
func MakeTestJob(master string, port string) (string, int64) {
	start := time.Now().Unix()
	param := nginxTestJobParam(master, port)
	r := CreateJob(param)
	param.Jobname = r
	logstr := getJobResult(param, "nginx:", 8, "nginx")
	times := time.Now().Unix() - start
	return logstr, times
}

// 2018-02-02 21:48
// 获取默认需要挂载的配置
func getNginxDefaultConf(conftype string) []ConfigureData {
	nginxConfigMap := []ConfigureData{}
	nginxConfigMap = append(nginxConfigMap, getNgxinDefaulgConfig(NginxSslPath, LbNginxSsl+conftype, nil, ""))
	nginxConfigMap = append(nginxConfigMap, getNgxinDefaulgConfig(NginxConfigPath, LbNginxConfig+conftype, nil, ""))
	nginxConfigMap = append(nginxConfigMap, getNgxinDefaulgConfig(NginxUpstreamPath, LbNginxUpstream+conftype, nil, ""))
	return nginxConfigMap
}

func CreateNginxLb(param ServiceParam) {
	param.Port = "80,443"
	param.MasterPort = "8080"
	param.Master = "10.16.55.114"
	param.Image = "nginx:v1"
	param.Namespace = "lb--nginx"
	param.Name = "nginx-lb"
	param.Cpu = 1
	param.Memory = "4096"
	param.HostPort = "80,443"
	param.NoUpdateConfig = true
	param.ClusterName = ""
	//param.Command = `["/usr/local/nginx/sbin/nginx"]`

	param.StorageData = `[{"ContainerPath":"/usr/local/nginx/logs/","HostPath":"/home/data/nginx/logs/"}]`

	clientset, _ := GetClient(param.ClusterName)
	cl2, _ := GetYamlClient(param.ClusterName, "", "v1", "api")
	param.Cl2 = cl2
	param.Cl3 = clientset
	param.ConfigureData = getNginxDefaultConf("")
	CreateConfigmap(param)
	CreateDeamonSet(param)
}
