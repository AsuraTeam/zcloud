package k8s

import (
	"strings"
	"cloud/util"
	"time"
	"github.com/astaxie/beego/logs"
	"cloud/sql"
)

// 镜像推送使用的参数
// 2018-02-06 08:13
type ImagePushParam struct {
	// 仓库IP
	Registry1Ip string
	Registry2Ip string
	// 仓库域名
	Registry1Domain string
	Registry2Domain string
	// 仓库端口
	Registry1Port string
	Registry2Port string
	// user:pass b64
	Registry2Auth string
	Registry1Auth string
	// 仓库组
	RegistryGroup string
	// 项目名
	ItemName string
	// 版本号
	Version string
	// 操作用户
	User string
	// 启动时间
	CreateTime string
}

// 2018-02-06 09:21
// 替换命令
func replacePushCmd(cmd string, param ImagePushParam) string {
	cmd = strings.Replace(cmd, "REGISTRY-1", param.Registry1Domain+":"+param.Registry1Port, -1)
	cmd = strings.Replace(cmd, "REGISTRY-2", param.Registry2Domain+":"+param.Registry1Port, -1)
	cmd = strings.Replace(cmd, "REGISTRYDOMAIN-1", param.Registry1Domain, -1)
	cmd = strings.Replace(cmd, "REGISTRYDOMAIN-2", param.Registry2Domain, -1)
	cmd = strings.Replace(cmd, "REGISTRYIP-1", param.Registry1Ip, -1)
	cmd = strings.Replace(cmd, "REGISTRYIP-2", param.Registry2Ip, -1)
	cmd = strings.Replace(cmd, "AUTH-1", param.Registry1Auth, -1)
	cmd = strings.Replace(cmd, "AUTH-2", param.Registry2Auth, -1)
	cmd = strings.Replace(cmd, "REGISTRYGROUP", param.RegistryGroup, -1)
	cmd = strings.Replace(cmd, "ITEMNAME", param.ItemName, -1)
	cmd = strings.Replace(cmd, "VERSION", param.Version, -1)
	return cmd
}

// 2018-02-06 09:06
// 在docker中执行镜像提交
func getPushCmd(param ImagePushParam) string {
	cmd := `mkdir /root/.docker -p
d=$(date +"%F %T")
echo "开始提交镜像...$d"
ping REGISTRYDOMAIN-1 -c 1
ping REGISTRYDOMAIN-2 -c 1
/usr/local/bin/dockerd --ip-forward=false --iptables=false --insecure-registry REGISTRY-1 --insecure-registry REGISTRY-2 &>/dev/null &
sleep 6
echo REGISTRYIP-1 REGISTRYDOMAIN-1 >> /etc/hosts
echo REGISTRYIP-2 REGISTRYDOMAIN-2 >> /etc/hosts
cat > /root/.docker/config.json <<EOF
{
	"auths": {
	      "REGISTRY-1": {
	         "auth": "AUTH-1"
	      },
	      "REGISTRY-2": {
	         "auth": "AUTH-2"
	      }
	},
	"HttpHeaders": {
	      "User-Agent": "Docker-Client/18.01.0-ce (linux)"
	}
}
EOF
chmod 700 /root/.docker -R
echo docker pull REGISTRY-1/ITEMNAME:VERSION 2>&1
echo
docker pull REGISTRY-1/ITEMNAME:VERSION 2>&1
if [ $? -gt 0 ] ; then
	   echo "pull镜像失败"
	   d=$(date +"%F %T")
       echo "完成提交... $d"
	   exit
fi
docker images|grep ITEMNAME:VERSION
echo docker push REGISTRY-2/ITEMNAME:VERSION 2>&1
echo
docker push REGISTRY-2/ITEMNAME:VERSION 2>&1
if [ $? -gt 0 ] ; then
	   echo "push镜像失败"
	   d=$(date +"%F %T")
       echo "完成提交... $d"
	   exit
fi
echo
d=$(date +"%F %T")
echo "push镜像成功..."
echo "完成提交... $d"
	`
	cmd = replacePushCmd(cmd, param)
	return cmd
}

// 2018-02-06 09;37
// 将操作命令挂载到容器中
func getImagePushConfig(param ImagePushParam) []ConfigureData {
	configureData := []ConfigureData{}
	config := ConfigureData{}
	config.ContainerPath = "/build"
	config.ConfigDbData = map[string]interface{}{"push.sh": getPushCmd(param)}
	name := strings.Replace(param.Registry1Domain+param.Registry1Port+param.Registry2Domain+param.Registry2Port, ".", "-", -1)
	config.DataName = "image-push" + name
	config.DataId = "push.sh"
	configureData = append(configureData, config)
	return configureData
}

// 2018-02-06 08:57
// 获取镜像推送的job参数
func imagePushJobParam(clusterName string, pushParam ImagePushParam) JobParam {
	param := JobParam{
		ClusterName:clusterName,
		//Master:    master,
		//Port:      port,
		Timeout:   600,
		Memory:    40,
		Cpu:       1,
		Namespace: util.Namespace("registryv2", "registryv2"),
		//Images:            "nginx:v1", 使用默认docker镜像
		ConfigureData: getImagePushConfig(pushParam),
		Command:       []string{"sh", "/build/push.sh"},
	}
	return param
}

// 2018-02-06 10:45
// 镜像提交完成后写入日志
const InsertCloudImageSyncLog = "insert into cloud_image_sync_log"
const UpdateCloudImageSyncLog = "update cloud_image_sync_log"

func writeImagePushToHistory(messages string, runtime int64, param ImagePushParam) {
	status := "同步中"
	if strings.Contains(messages, "push镜像成功...") {
		status = "成功"
	}
	if strings.Contains(messages, "pull镜像失败"){
		status = "失败"
	}
	synclog := CloudImageSyncLog{
		CreateTime:      param.CreateTime,
		CreateUser:      param.User,
		Messages:        messages,
		Runtime:         runtime,
		RegistryGroup:   param.RegistryGroup,
		RegistryServer1: param.Registry1Domain + ":" + param.Registry1Port,
		RegistryServer2: param.Registry2Domain + ":" + param.Registry2Port,
		ItemName:        param.ItemName,
		Version:         param.Version,
		Status:          status,
	}
	var q string
	if runtime == 0 {
		q = sql.InsertSql(synclog, InsertCloudImageSyncLog)
	}else{
		searchMap := sql.SearchMap{}
		searchMap.Put("Runtime", 0)
		searchMap.Put("CreateTime", param.CreateTime)
		q = sql.UpdateSql(synclog, UpdateCloudImageSyncLog, searchMap, "LogId")
	}
	sql.Raw(q).Exec()
}

// 将A的镜像推送到B集群去
// 镜像推送服务
func ImagePush(clusterName string, imagePushParam ImagePushParam) {
	start := time.Now().Unix()
	jobParam := imagePushJobParam(clusterName, imagePushParam)
	jobParam.Jobname = "job-" + util.Md5Uuid()
	jobName := CreateJob(jobParam)
	writeImagePushToHistory("", 0, imagePushParam)
	jobParam.Jobname = jobName
	logstr := getJobResult(jobParam, "完成提交", 300, "")
	times := time.Now().Unix() - start
	writeImagePushToHistory(logstr, times, imagePushParam)
	logs.Info(logstr, times)
}
