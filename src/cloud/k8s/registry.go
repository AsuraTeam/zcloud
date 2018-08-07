package k8s

import (
	"encoding/json"
	"github.com/astaxie/beego/logs"
	"path/filepath"
	"cloud/util"
	"os"
	"strings"
	"golang.org/x/crypto/openpgp/errors"
	"net/http"
	"crypto/tls"
	"github.com/heroku/docker-registry-client/registry"
	"fmt"
)

var (
	StartCmd = `
	 mkdir /usr/share/zoneinfo/Asia/ -p
	 cp /etc/localtime /usr/share/zoneinfo/Asia/Shanghai
	 registry serve /etc/docker/registry/config-yml
`

	RegistryTemplate = `version: 0.1
log:
  fields:
    service: registry
storage:
    cache:
        blobdescriptor: inmemory
    filesystem:
        rootdirectory: /var/lib/registry
http:
    addr: :5000
    headers:
        X-Content-Type-Options: [nosniff]
    tls:
      certificate: /certs/registry-crt
      key: /certs/registry-key
health:
  storagedriver:
    enabled: true
    interval: 10s
    threshold: 3
auth:
  token:
    realm: "AUTH-SERVER"
    service: "SERVICE"
    issuer: "Acme auth server"
    rootcertbundle: /certs/auth-crt
    `
)

// 在k8s集群中创建仓库使用的配置信息
// 2018-01-20 21;01
// 量REGISTRY_STORAGE_DELETE_ENABLED=true
func getParam(registryParam RegistryParam) ServiceParam {
	param := ServiceParam{}
	param.Name = registryParam.Name
	param.ServiceName = registryParam.Name
	param.Cpu = 1
	param.ClusterName = registryParam.ClusterName
	param.PortData = "5000"
	param.Replicas = 1
	param.Namespace = util.Namespace("registryv2", "registryv2")
	param.Memory = "2048"
	param.Port = "5000"
	param.Image = "registry:2"
	param.MinReady = 1
	param.HealthData = ""
	//param.StorageData = `[{"ContainerPath":"/etc/localtime","Volume":"","HostPath":"/etc/localtime"}]`
	param.Command = `["sh","/start/start-cmd"]`
	// deployment
	c1, _ := GetYamlClient(registryParam.ClusterName, "apps", "v1beta1", "/apis")
	// service
	cl2, _ := GetYamlClient(registryParam.ClusterName, "", "v1", "api")
	cl3, _ := GetClient(registryParam.ClusterName)
	param.Cl3 = cl3
	param.Cl2 = cl2
	param.C1 = c1
	param.Envs = "TZ=Asia/Shanghai"

	config := `[{"ContainerPath":"/certs","DataName":"registry-auth","DataId":"auth-crt,registry-crt,registry-key"},{"ContainerPath":"/start","DataName":"registry-auth","DataId":"start-cmd"},{"ContainerPath":"/etc/docker/registry","DataName":"registry-config","DataId":"config-yml"}]`
	// 生产configmap信息
	configdata := make([]ConfigureData, 0)
	err := json.Unmarshal([]byte(config), &configdata)
	logs.Info(err)
	configureData := make([]ConfigureData, 0)

	// 读取配置文件里的证书信息
	pwd, _ := os.Getwd()
	keyFile := filepath.Join(pwd, "conf", "key", "server.key")
	pemFile := filepath.Join(pwd, "conf", "key", "server.pem")
	registryKey := util.ReadFile(keyFile)
	registryCrt := util.ReadFile(pemFile)

	registryConf := strings.Replace(RegistryTemplate, "AUTH-SERVER", registryParam.AuthServer, -1)
	registryConf = strings.Replace(registryConf, "SERVICE", registryParam.Name+"."+registryParam.ClusterName, -1)

	for _, v := range configdata {
		ConfigDbData := map[string]interface{}{
			"config-yml":   registryConf, // registry 配置文件
			"registry-crt": registryCrt,
			"registry-key": registryKey, // registry 使用
			"auth-crt":     registryCrt, // 验证token使用的
			"start-cmd":    StartCmd,    // 启动命令
		}
		v.ConfigDbData = ConfigDbData
		configureData = append(configureData, v)
	}
	param.ConfigureData = configureData
	return param
}

// 2018-01-20 20:56
// 创建镜像仓库
func CreateRegistry(param RegistryParam) (error) {
	if !strings.Contains(param.AuthServer, "https://") || !strings.Contains(param.AuthServer, "/auth") {
		logs.Error("认证服务失败", param.AuthServer)
		return errors.InvalidArgumentError("认证服务器失败")
	}
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	r, err := client.Get(param.AuthServer)
	defer r.Body.Close()
	logs.Info("验证服务器返回信息", r, err)
	if err != nil  {
		return errors.InvalidArgumentError("认证服务器失败" + err.Error())
	}
	yaml, err := CreateServicePod(getParam(param))
	logs.Info(yaml, err)
	return err
}

// 2018-01-28 13:36
// 获取访问连接
func getHubClient(host string, username string, password string) *registry.Registry {
	url := "https://" + host + "/"
	hub, err := registry.New(url, username, password)
	if err != nil {
		logs.Error("获取仓库连接失败", url, err, username)
		return nil
	}
	return hub
}

// 2018-01-27 21:06
// 获取仓库中不同组的镜像数量和tag数量
func GetRegistryInfo(host string, username string, password string, registryName string) (util.Lock, util.Lock, util.Lock) {
	hub := getHubClient(host, username, password)
	if hub == nil {
		return util.Lock{}, util.Lock{}, util.Lock{}
	}
	repositories, err := hub.Repositories()
	fmt.Println(repositories)
	if err != nil {
		logs.Error("获取仓库数据失败", err, username)
		return util.Lock{}, util.Lock{}, util.Lock{}
	}
	imagesLock := util.Lock{}
	lock := util.Lock{}
	for _, v := range repositories {
		vs := strings.Split(v, "/")
		key := vs[0]
		img := CloudImage{}
		img.RepositoriesGroup = key
		img.Name = v

		imagesLock.Put(v, img)
		if _, ok := lock.Get(key); ok {
			lock.Put(key, lock.GetV(key).(int)+1)
			continue
		}
		lock.Put(key, 1)
	}
	tagLock := util.Lock{}
	for v := range lock.GetData() {
		tag, _ := hub.Tags(v)
		tagLock.Put(v, len(tag))
	}
	for k, v := range imagesLock.GetData() {
		img := v.(CloudImage)
		tag, _ := hub.Tags(k)
		img.TagNumber = len(tag)
		img.Tags = strings.Join(tag, ",")
		if len(tag) > 0 {
			size := int64(0)
			m, err := hub.ManifestV2(k, tag[len(tag)-1])
			if err == nil {
				size += m.Config.Size
				for _, mani := range m.Layers {
					size += mani.Size
				}
			}
			img.Repositories = registryName
			img.LayersNumber = len(m.Layers)
			img.Size = size
		}
		img.Access = host
		imagesLock.Put(k, img)
	}
	return lock, tagLock, imagesLock
}

// 2018-01-29 08:44
// 删除镜像
func deleteImage(hub *registry.Registry, imagename string, tag string) (bool, error) {
	digest, err := hub.ManifestDigest(imagename, tag)
	if err != nil {
		return false, err
	}
	err = hub.DeleteManifest(imagename, digest)
	if err != nil {
		return false, err
	}
	return true, nil
}

// 2018-01-29 8:27
// 删除镜像
func DeleteRegistryImage(host string, username string, password string, imagename string, tag string) (bool, error) {
	hub := getHubClient(host, username, password)
	if hub == nil {
		return false, errors.UnsupportedError("连接registry server失败")
	}
	var r bool
	var err error
	if tag != "" {
		r, err = deleteImage(hub, imagename, tag)
	} else {
		tags, err := hub.Tags(imagename)
		if err != nil {
			return false, err
		}
		for _, tag := range tags {
			r, err = deleteImage(hub, imagename, tag)
			if err != nil {
				return false, err
			}
			logs.Info("删除镜像", imagename, tag)
		}
	}
	return r, err
}

// 2018-02-09 17:01
// 检查镜像是否存在
func CheckImageExists(host string, username string, password string, imagename string, tag string) (bool) {
	hub := getHubClient(host, username, password)
	if hub == nil {
		logs.Error("连接失败")
		return false
	}
	_, err := hub.ManifestDigest(imagename, tag)
	if err != nil {
		return false
	}
	return true
}

