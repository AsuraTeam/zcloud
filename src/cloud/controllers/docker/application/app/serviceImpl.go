package app

import (
	"github.com/astaxie/beego/logs"
	"cloud/k8s"
	"cloud/models/app"
	"cloud/util"
	"golang.org/x/crypto/openpgp/errors"
	"cloud/sql"
	"strconv"
	"encoding/json"
	"strings"
	"cloud/controllers/image"
	"time"
	"cloud/models/ci"
	"cloud/cache"
	"cloud/userperm"
	"cloud/models/log"
	"cloud/models/ent"
)

// 2018-02-13 16:36
func getServiceData(searchMap sql.SearchMap, q string) []app.CloudAppService {
	if q == "" {
		q = app.SelectCloudAppService
	}
	data := make([]app.CloudAppService, 0)
	searchSql := sql.SearchSql(
		app.CloudAppService{}, q, searchMap)

	searchSql = sql.SearchOrder(searchSql, "create_time")
	sql.Raw(searchSql).QueryRows(&data)
	return data
}

// 2018-02-13 15:53 getServiceData
// service选择查询
func GetServiceHtml(searchMap sql.SearchMap) string {
	data := getServiceData(searchMap, app.SelectServiceName)
	var opt string
	for _, v := range data {
		opt += util.GetSelectOptionName(v.ServiceName)
	}
	return opt
}

// 2018-02-03 21:30
// 获取服务名称数据
func GetSelectHtml(username string, cluster string) string {
	data := GetUserLbService(username, cluster, "")
	var html string
	for _, v := range data {
		html += util.GetSelectOptionName(v.ServiceName)
	}
	return html
}

// 2018-02-04 16:08
// 更新服务
func ExecUpdate(service app.CloudAppService, updateType string, username string) error {
	param := getParam(service, username)
	param.Update = true
	param.UpdateType = updateType
	var err error
	if updateType == "image" {
		logs.Info("ExecUpdate", param.Image)
		rollparam := k8s.RollingParam{
			Client:          param.Cl3,
			Images:          param.Image,
			Namespace:       param.Namespace,
			Name:            param.ServiceName,
			MaxSurge:        1,
			MinReadySeconds: 80,
		}
		if param.Replicas == 1 {
			rollparam.MaxUnavailable = 0
		}
		logs.Info("更新镜像", util.ObjToString(rollparam))
		_, err = k8s.UpdateDeploymentImage(rollparam)
	} else {
		_, err = k8s.CreateServicePod(param)
	}
	if err == nil {
		updateServiceData(service, username)
		go updateContainerData(service)
	} else {
		logs.Error("ExecUpdate 失败 ", err.Error())
	}
	return err
}

// 2018-02-15 15:39
// 更新完成后,更新容器数据
func updateContainerData(service app.CloudAppService) {
	for i := 0; i <= 5; i++ {
		go MakeContainerData(util.Namespace(service.AppName, service.ResourceName))
		time.Sleep(time.Second * 5)
	}
}

// 更新服务数据
// 2018-01-14 07:54
func updateServiceData(service app.CloudAppService, username string) {
	// 操作成功后更新数据库数
	searchMap := sql.SearchMap{}
	searchMap.Put("ServiceId", service.ServiceId)
	service.LastModifyTime = util.GetDate()
	service.LastModifyUser = username
	updateSql := sql.UpdateSql(service,
		app.UpdateCloudAppService,
		searchMap,
		"CreateTime,CreateUser")
	sql.Raw(updateSql).Exec()
}

// 检查服务参数是否正确
// 2018-01-14 07:53
func checkParam(service app.CloudAppService, cpuerr error, memerr error, this *ServiceController) bool {

	// 检查内存和cpu的数据准确性
	if cpuerr != nil || service.Cpu > 128 || service.Cpu < 0.1 {
		if cpuerr == nil {
			cpuerr = errors.InvalidArgumentError("")
		}
		responseData(cpuerr, this, service.ServiceName, "操作失败: cpu参数不正确,正常范围:0.1-128 ")
		return false
	}

	if memerr != nil || service.Memory > 262144 || service.Memory < 100 {
		if memerr == nil {
			memerr = errors.InvalidArgumentError("")
		}
		responseData(memerr, this, service.ServiceName, "操作失败: mem参数不正确,正常范围 256-262144 ")
		return false
	}
	return true
}

// 2018-02-12 21:57
// 通过任务计划写入服务缓存
func CronServiceCache() {
	logs.Info("写入服务缓存")
	data := make([]app.CloudAppService, 0)
	q := sql.SearchSql(app.CloudAppService{}, app.SelectCloudAppService, sql.SearchMap{})
	sql.Raw(q).QueryRows(&data)
	go GoServerThread(data)
}

// 2018-02-04
// 从redis里获取应用服务运行状态数据
func GetServiceRunData(data []app.CloudAppService, user string) []k8s.CloudApp {
	//result := make([]interface{}, 0)
	perm := userperm.GetResourceName("服务", user)
	permApp := userperm.GetResourceName("应用", user)

	result := make([]k8s.CloudApp, 0)
	for _, d := range data {
		// 不是自己创建的才检查
		if d.CreateUser != user && user != "admin"{
			pName := d.AppName+";"+d.ResourceName+";"+d.ServiceName
			if ! userperm.CheckPerm(pName, d.ClusterName, d.Entname, perm) && len(user) > 0 {
				if ! userperm.CheckPerm(d.AppName, d.ClusterName, d.Entname, permApp) {
					continue
				}
			}
		}

		namespace := util.Namespace(d.AppName, d.ResourceName) +
			d.ServiceName  + d.Entname + d.ClusterName
		if len(d.ServiceVersion) > 0 {
			namespace = util.Namespace(namespace, d.ServiceVersion)
		}
		namespace += strconv.FormatInt(d.ServiceId, 10)
		r := cache.ServiceCache.Get(namespace)
		var v = k8s.CloudApp{}
		status := util.RedisObj2Obj(r, &v)

		if status {
			result = append(result, v)
		} else {
			t := k8s.CloudApp{}
			dbyte, err := json.Marshal(d)
			if err == nil {
				json.Unmarshal(dbyte, &t)
				result = append(result, t)
			}
		}
	}
	return result
}

// 修改数据时公共数据
// 2018-01-14 13:35
func setChangeData(this *ServiceController) app.CloudAppService {
	service := getService(this)
	this.Data["data"] = service
	return service
}

// 2018-10-03 14:57
// 获取日志驱动
func getLogDriver(ent string, cluster string) log.LogDataSource {
	searchMap :=  sql.SearchMap{}
	searchMap.Put("Ent", ent)
	searchMap.Put("ClusterName", cluster)
	searchMap.Put("DataType", "driver")
	q := sql.SearchSql(log.LogDataSource{}, log.SelectLogDataSource,searchMap)
	data := log.LogDataSource{}
	sql.Raw(q).QueryRow(&data)
	return data
}

// 获取环境英文名称
func GetEntDescription(entname string)  string{
	searchMap := sql.SearchMap{}
	searchMap.Put("Entname", entname)
	q := sql.SearchSql(ent.CloudEnt{}, ent.SelectCloudEnt, searchMap)
	e := ent.CloudEnt{}
	sql.Raw(q).QueryRow(&e)
	return e.Description
}

// 设置filebeat需要的参数
func setFilebeatParam(param k8s.ServiceParam, d app.CloudAppService)   k8s.ServiceParam{
	if len(d.LogPath) > 0 {
		dataDriver := getLogDriver(d.Entname, param.ClusterName)
		param.LogPath = d.LogPath
		param.Kafka = dataDriver.Address
		if dataDriver.DriverType == "elasticsearch"{
			param.ElasticSearch = dataDriver.Address
		}
		param.Ent = GetEntDescription(d.Entname)
	}
	return param
}

// 获取创建服务的配置参数
// 2018-01-12 8:56
func getParam(d app.CloudAppService, user string) k8s.ServiceParam {
	param := k8s.ServiceParam{}
	param.Name = d.ServiceName
	if d.ServiceVersion == "" {
		d.ServiceVersion = "1"
	}
	param.ServiceName = util.Namespace(d.ServiceName, d.ServiceVersion)
	param.Cpu = d.Cpu
	param.ClusterName = d.ClusterName
	param.PortData = d.ContainerPort
	param.Replicas = d.Replicas
	param.Namespace = util.Namespace(d.AppName, d.ResourceName)
	param.Memory = strconv.FormatInt(d.Memory, 10)
	param.Port = d.ContainerPort
	param.Image = d.ImageTag
	param.MinReady = d.MinReady
	param.HealthData = d.HealthData
	param.ResourceName = d.ResourceName
	param.StorageData = d.StorageData
	param.PortYaml = d.Yaml
	param.NetworkMode = d.NetworkMode

	param = setFilebeatParam(param, d)

	// 关闭容器时间
	if param.TerminationSeconds == 0 {
		param.TerminationSeconds = 50
	}

	if d.ReplicasMax > param.Replicas {
		param.ReplicasMax = d.ReplicasMax
	} else {
		param.ReplicasMax = d.Replicas
	}
	master, port := k8s.GetMasterIp(d.ClusterName)
	// deployment
	c1, _ := k8s.GetYamlClient(d.ClusterName, "apps", "v1beta1", "/apis")
	// service
	cl2, _ := k8s.GetYamlClient(d.ClusterName, "", "v1", "api")
	cl3, _ := k8s.GetClient(d.ClusterName)
	param.Cl3 = cl3
	param.Cl2 = cl2
	param.C1 = c1
	param.Envs = d.Envs
	param.Master = master
	param.MasterPort = port
	param.CreateUser = user
	config := d.ConfigureData
	if config != "" {
		configureData := make([]k8s.ConfigureData, 0)
		configData := make([]k8s.ConfigureData, 0)
		json.Unmarshal([]byte(config), &configData)
		for _, v := range configData {
			v.ConfigDbData = GetConfgData(v.DataName, d.ClusterName)
			configureData = append(configureData, v)
		}
		param.ConfigureData = configureData
	}
	param = CreateSecretFile(param)
	return param
}

// 响应错误数据
func responseData(err error, this *ServiceController, serviceName string, info string) {
	data, msg := util.SaveResponse(err, info)
	util.SaveOperLog(getServiceUser(this), *this.Ctx, info+": "+msg, serviceName)
	setServiceJson(this, data)
}

// 2018-02-09 21:32
// 创建secret文件
func CreateSecretFile(param k8s.ServiceParam) k8s.ServiceParam {
	// 创建私有仓库镜像获取私密文件
	param.Registry = strings.Split(param.Image, "/")[0]
	logs.Info("获取到仓库地址", param.Registry)
	servers := registry.GetRegistryServer(strings.Split(param.Registry, ":")[0])
	if len(servers) > 0 {
		param.RegistryAuth = servers[0].Admin + ":" + util.Base64Decoding(servers[0].Password)
		k8s.CreateImagePullSecret(param)
		logs.Info("创建私密文件完成")
	}
	return param
}

// 查询某个服务的数据
func GetServiceData(name string, cluster string, appname string) app.CloudAppService {
	data := app.CloudAppService{}
	searchMap := sql.GetSearchMapV(
		"ServiceName",
		name, "ClusterName",
		cluster, "AppName",
		appname)
	searchSql := sql.SearchSql(app.CloudAppService{}, app.SelectCloudAppService, searchMap)

	sql.Raw(searchSql).QueryRow(&data)
	return data
}

// 2018-02-01 15:15
// 获取某个用户的所有服务
func GetUserLbService(user string, clusterName string, id string) []app.CloudAppService {
	data := make([]app.CloudAppService, 0)
	searchMap := sql.GetSearchMapV(
		"CreateUser",
		user,
		"ClusterName",
		clusterName)

	if id != "" {
		searchMap.Put("ServiceId", id)
	}

	searchSql := sql.SearchSql(
		app.CloudAppService{},
		app.SelectCloudAppService,
		searchMap)
	sql.Raw(searchSql).QueryRows(&data)
	return data
}

// 2018-08-10 15:36
// 数据写入redis
func serviceToRedis(namespace string, id int64, sv k8s.CloudApp) {
	cache.ServiceCache.Put(
		namespace+strconv.FormatInt(id, 10),
		util.ObjToString(sv),
		time.Minute * 10)
}

// 2018-01-31 16:04
// 后台执行服务状态,更新到缓存里
func GoServerThread(data []app.CloudAppService) {
	appDatas := util.Lock{}
	for _, d := range data {
		go goServiceData(d, &appDatas)
	}

	result := make([]interface{}, 0)
	counter := 0
	for {
		result = make([]interface{}, 0)
		for _, d := range data {
			namespace := util.Namespace(d.AppName, d.ResourceName) + d.ServiceName + d.Entname + d.ClusterName
			if d.ServiceVersion != "" {
				namespace = util.Namespace(namespace, d.ServiceVersion)
			}

			v, ok := appDatas.Get(namespace)
			if ok {
				sv := v.(k8s.CloudApp)
				sv.ClusterName = d.ClusterName
				sv.CheckTime = time.Now().Unix()
				sv.Domain = d.Domain
				sv.Entname = d.Entname
				serviceToRedis(namespace, d.ServiceId, sv)
				result = append(result, v)
			} else {
				sv := k8s.CloudApp{}
				r := cache.ServiceCache.Get(namespace + strconv.FormatInt(d.ServiceId, 10))
				s := util.RedisObj2Obj(r, &sv)
				now := time.Now().Unix()
				if s && now-sv.CheckTime > 60 * 15 {
					sv.Status = "False"
					sv.AvailableReplicas = 0
					serviceToRedis(namespace, d.ServiceId, sv)
				}
			}
		}
		if len(result) >= len(data) {
			break
		}
		time.Sleep(time.Second * 1)
		counter += 1
		if counter > 5 {
			break
		}
	}
}

// 通过多线程去跑
// 2018-01-13 08-07
func goServiceData(d app.CloudAppService, appDatas *util.Lock) {

	namespace := util.Namespace(d.AppName, d.ResourceName)
	sname := namespace + d.ServiceName + d.Entname + d.ClusterName
	if d.ServiceVersion != "" {
		sname = util.Namespace(sname, d.ServiceVersion)
	}

	if _, ok := appDatas.Get(sname); !ok {
		c, _ := k8s.GetClient(d.ClusterName)
		serviceName := d.ServiceName
		if d.ServiceVersion != "" {
			serviceName = util.Namespace(serviceName, d.ServiceVersion)
		}
		appData := k8s.GetDeploymentApp(c, namespace, serviceName)
		for _, all := range appData {
			all.ResourceName = d.ResourceName
			all.ServiceId = d.ServiceId
			appDatas.Put(sname, all)
		}
	} else {
		d1 := k8s.CloudApp{}
		t, _ := json.Marshal(d)
		json.Unmarshal(t, &d1)
		d1.Status = "获取失败"
		d1.Access = []string{""}
		appDatas.Put(sname, d1)
	}
}

// 获取某个服务的数据
// 2018-01-13 11:26
func getService(this *ServiceController) app.CloudAppService {
	searchMap := sql.SearchMap{}
	id := this.Ctx.Input.Param(":id")
	searchMap.Put("ServiceId", id)
	user := getServiceUser(this)
	perm := userperm.GetResourceName("服务", user)
	permApp := userperm.GetResourceName("应用", user)

	d := app.CloudAppService{}
	q := sql.SearchSql(d, app.SelectCloudAppService, searchMap)
	sql.Raw(q).QueryRow(&d)
	// 不是自己创建的才检查
	if d.CreateUser != user && user != util.ADMIN {
		if ! userperm.CheckPerm(d.AppName+";"+d.ResourceName+";"+d.ServiceName, d.ClusterName, d.Entname, perm) {
			if ! userperm.CheckPerm(d.AppName, d.ClusterName, d.Entname, permApp) {
				return app.CloudAppService{}
			}
		}
	}
	return d
}

// 设置json数据
func setServiceJson(this *ServiceController, data interface{}) {
	this.Data["json"] = data
	this.ServeJSON(false)
}

func getServiceUser(this *ServiceController) string {
	return util.GetUser(this.GetSession("username"))
}

// 2018-2-13 17:53
// 保存部署数据
func saveServiceDeploy(d app.CloudAppService) {
	var q = sql.InsertSql(d, app.InsertCloudAppService)
	if d.ServiceId > 0 {
		id := strconv.FormatInt(d.ServiceId, 10)
		searchMap := sql.GetSearchMapV("ServiceId", id)

		q = sql.UpdateSql(d,
			app.UpdateCloudAppService,
			searchMap,
			app.UpdateCloudAppServiceWhere)
	}
	sql.Raw(q).Exec()
}

// 2018-02-13 16:49
// 获取蓝是否存在
func getBlueExists(greenData app.CloudAppService, version string) bool {
	searchMap := sql.GetSearchMapV("ServiceVersion", version,
		"ClusterName", greenData.ClusterName,
		"AppName", greenData.AppName,
		"ServiceName", greenData.ServiceName,
		"Entname", greenData.Entname,
		"ResourceName", greenData.ResourceName)
	blueData := getServiceData(searchMap, "")
	if len(blueData) > 0 {
		return true
	}
	return false
}

// 2018-02-14 07:20
// 获取蓝绿部署的服务信息
func GetServices(ciData ci.CloudCiService, q string) []app.CloudAppService {
	searchMap := sql.GetSearchMapV(
		"ServiceName", ciData.ServiceName,
		"ClusterName", ciData.ClusterName,
		"AppName", ciData.AppName,
	)
	services := getServiceData(searchMap, q)
	return services
}

// 2018-02-14 07;30
// 获取当前版本
func GetCurrentVersion(ciData ci.CloudCiService, services []app.CloudAppService) util.Lock {

	lock := util.Lock{}
	for _, v := range services {
		lock.Put(v.ServiceVersion, v.ImageTag)
	}
	return lock
}

// 蓝绿部署,启动一个绿的
func CreateGreenService(ciData ci.CloudCiService, username string) (string, bool, app.CloudAppService) {
	services := GetServices(ciData, "")
	if len(services) == 0 || len(services) == 2 {
		msg := "没有获取到Service数据或则版本已经发布,程序退出"
		logs.Error(msg)
		return msg, false, app.CloudAppService{}
	}

	serviceData := services[0]
	version := "1"
	if serviceData.ServiceVersion == "1" {
		version = "2"
	}

	if getBlueExists(serviceData, version) {
		msg := "俩个部署版本同时存在不能同时存在,可能是没有结束部署"
		logs.Error(msg)
		return msg, false, app.CloudAppService{}
	}

	serviceData.ImageTag = ciData.ImageName
	serviceData.ServiceVersion = version
	//serviceData.ServiceName = serviceData.ServiceName
	serviceParam := getParam(serviceData, username)

	yaml, err := k8s.CreateServicePod(serviceParam)
	if err != nil {
		msg := "创建服务失败 k8s执行错误 " + err.Error()
		logs.Error(msg)
		return msg, false, app.CloudAppService{}
	}
	serviceData.CreateTime = util.GetDate()
	serviceData.Yaml = yaml
	serviceData.ServiceId = 0
	saveServiceDeploy(serviceData)
	return "", true, serviceData
}

// 2018-02-18 11:04
// 获取需要刷新redis缓存的数据
func getServices(d app.CloudAppService) []app.CloudAppService {
	data := make([]app.CloudAppService, 0)
	searchSql := sql.SearchSql(
		app.CloudAppService{},
		app.SelectCloudAppService,
		sql.GetSearchMapV("ServiceName", d.ServiceName,
			"ClusterName", d.ClusterName,
			"AppName", d.AppName))
	sql.Raw(searchSql).QueryRows(&data)
	return data
}

// 2018-02-18 11:00
// 服务创建后持续更新redis换成
func updateServiceRedisCache(d app.CloudAppService) {
	services := getServices(d)
	for i := 1; i < 5; i ++ {
		go GoServerThread(services)
		go MakeContainerData(util.Namespace(d.AppName, d.ResourceName))
		time.Sleep(time.Second * 5)
	}
}
