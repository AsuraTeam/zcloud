package app

import (
	"cloud/sql"
	"cloud/util"
	"cloud/models/app"
	"github.com/astaxie/beego"
	"cloud/k8s"
	"strconv"
	"github.com/astaxie/beego/logs"
	"encoding/json"
	"time"
	"golang.org/x/crypto/openpgp/errors"
	"strings"
	"cloud/controllers/base/storage"
	"cloud/controllers/image"
	"cloud/userperm"
)

type ServiceController struct {
	beego.Controller
}



// Service 管理入口页面
// @router /application/service/list [get]
func (this *ServiceController) ServiceList() {
	go registry.UpdateGroupImageInfo()
	this.Data["ServiceName"] = this.GetString("name")
	this.TplName = "application/service/list.html"
}

// Service 创建服务添加配置文件页面
// @router /application/service/configure/add [get]
func (this *ServiceController) ConfigureAdd() {
	this.TplName = "application/service/add_config.html"
}

// 2018-01-13 17:01
// Service 服务管理扩容页面
// @param
// @router /application/service/scale/add/:id:int [get]
func (this *ServiceController) ScaleAdd() {
	service := setChangeData(this)

	v, err := k8s.GetAutoScale(service.ClusterName, util.Namespace(service.AppName, service.ResourceName), service.ServiceName)
	this.Data["MaxReplicas"] = service.ReplicasMax
	if err == nil {
		this.Data["MaxReplicas"] = v.Spec.MaxReplicas
	}
	this.TplName = "application/service/scale.html"
}

// 2018-01-13 18:43
// Service 服务管理修改配置
// @param
// @router /application/service/config/add/:id:int [get]
func (this *ServiceController) ConfigAdd() {
	this.Data["data"] = getService(this)
	this.TplName = "application/service/change_cpu_mem.html"
}

// 2018-01-14 09:31
// 应用Service 修改滚动升级页面
// @router /application/service/image/add/:id:int [get]
func (this *ServiceController) ImageAdd() {
	service := setChangeData(this)
	tag := strings.Split(service.ImageTag, ":")
	if len(tag) > 1 {
		this.Data["tag"] = tag[1]
	}
	images := registry.GetImageTag(service.ImageTag)
	this.Data["images"] = images
	this.TplName = "application/service/image.html"
}

// 2018-01-14 11:13
// 应用Service 修改环境变量
// @router /application/service/env/add/:id:int [get]
func (this *ServiceController) EnvAdd() {
	setChangeData(this)
	this.TplName = "application/service/env.html"
}

// 2018-01-14 13:31
// 应用Service 修改端口数据
// @router /application/service/port/add/:id:int [get]
func (this *ServiceController) PortChange() {
	setChangeData(this)
	this.TplName = "application/service/port.html"
}

// Service 创建服务添加存储页面
// @router /application/service/storage/add [get]
func (this *ServiceController) StorageAdd() {
	storage := storage.GetStorageName(getServiceUser(this), this.GetString("ClusterName"))
	var html = "<option value=''>--请选择--</option>"
	for _, v := range storage {
		html += util.GetSelectOptionName(v.Name)
	}
	this.Data["storage"] = html
	this.TplName = "application/service/add_storage.html"
}

// Service 创建服务添加健康检查页面
// @router /application/service/health/add [get]
func (this *ServiceController) HealthAdd() {
	this.TplName = "application/service/add_health.html"
}

// Service 创建服务添加健康检查页面
// @router /application/service/health/add/:id:int [get]
func (this *ServiceController) HealthChange() {
	service := getService(this)
	conf := k8s.HealthData{}
	conf.HealthType = "TCP"
	conf.HealthPort = "8080"
	conf.HealthFailureThreshold = "0"
	conf.HealthInterval = "60"
	conf.HealthPath = "/"
	conf.HealthTimeout = "20"
	conf.HealthCmd = "ls /tmp"
	conf.HealthInitialDelay = "30"
	if len(service.HealthData) > 10 {
		err := json.Unmarshal([]byte(service.HealthData), &conf)
		if err != nil{
			logs.Error("检查检查转换错误", err.Error())
		}
	}
	logs.Info(util.ObjToString(conf), service.HealthData)
	this.Data["config"] = conf
	this.Data["data"] = service
	this.TplName = "application/service/change_health.html"
}

// Service 管理添加页面
// @router /application/service/add [get]
func (this *ServiceController) ServiceAdd() {
	id := this.GetString("ServiceId")
	// 更新操作

	if id != "" {
		searchMap := sql.GetSearchMap("ServiceId", *this.Ctx)
		update := app.CloudAppService{}
		q := sql.SearchSql(
			app.CloudAppService{},
			app.SelectCloudAppService,
			searchMap)
		sql.Raw(q).QueryRow(&update)
		this.Data["data"] = update
	}

	this.TplName = "application/service/add.html"
}

// string
// Service 保存
// @router /api/service [post]
func (this *ServiceController) ServiceSave() {
	d := app.CloudAppService{}
	err := this.ParseForm(&d)
	if err != nil {
		this.Ctx.WriteString("参数错误" + err.Error())
		return
	}
	util.SetPublicData(d, getServiceUser(this), &d)

	serviceData := GetServiceData(d.ServiceName, d.ClusterName, d.AppName)

	if serviceData.ServiceId > 0 {
		logs.Error("创建服务失败", "该服务已经存在")
		responseData(err, this, d.ServiceName, "该服务已经存在")
		return
	}

	status, msg := k8s.CheckQuota(
		getServiceUser(this), d.Replicas,
		int64(d.Cpu), d.Memory,
		d.ResourceName)

	if !status {
		logs.Error("用户超过配额", msg)
		responseData(errors.InvalidArgumentError(msg), this, d.ServiceName, msg)
		return
	}

	d, err  = ExecDeploy(d, false)
	if err != nil {
		logs.Error("创建服务失败", "k8s执行错误", err.Error())
		responseData(err, this, d.ServiceName, "创建服务时失败")
		return
	}

	if len(d.Domain) > 0 {
		createLbConfig(d, d.ClusterName, d.Entname, d.AppName, d.Domain, getServiceUser(this), d.ResourceName)
		go k8s.CreateNginxConf("")
	}

	data, msg := util.SaveResponse(nil, "保存成功")
	util.SaveOperLog(getServiceUser(this), *this.Ctx,
		"保存Service 配置 "+msg, d.ServiceName)
	setServiceJson(this, data)
	saveAppData(d)

}

// 创建服务公用
func ExecDeploy(d app.CloudAppService, isRedeploy bool) (app.CloudAppService, error) {
	serviceParam := getParam(d, d.CreateUser)
	serviceParam.IsRedeploy = isRedeploy
	yaml, err := k8s.CreateServicePod(serviceParam)
	if err != nil {
		logs.Error("创建服务失败", "k8s执行错误", err.Error())
		return d, err
	}

	d.Yaml = yaml
	saveServiceDeploy(d)
	return d, nil
}

// Service 名称数据
// @param AppName
// @param ServiceName
// @param ClusterName
// @router /api/service/name [get]
func (this *ServiceController) GetServiceName() {
	data := make([]app.CloudAppServiceName, 0)
	key := strings.Split(app.ServiceSearchKey, ",")

	searchMap := sql.GetSearchMapValue(key, *this.Ctx, sql.SearchMap{})
	searchSql := sql.SearchSql(
		app.CloudAppService{},
		app.SelectCloudAppService,
		searchMap)

	sql.Raw(searchSql).QueryRows(&data)

	user := getServiceUser(this)
	perm := userperm.GetResourceName("服务", user)
	permApp := userperm.GetResourceName("应用", user)
	result := make([]app.CloudAppServiceName, 0)
	for _, d := range data {
		// 不是自己创建的才检查
		if d.CreateUser != user && user != "admin" {
			if ! userperm.CheckPerm(d.AppName+";"+d.ResourceName+";"+d.ServiceName, d.ClusterName, d.Entname, perm) {
				if ! userperm.CheckPerm(d.AppName, d.ClusterName, d.Entname, permApp) {
					continue
				}
			}
		}
		result = append(result, d)
	}

	setServiceJson(this, result)
}

// Service 数据
// @router /api/service [get]
func (this *ServiceController) ServiceData() {
	data := make([]app.CloudAppService, 0)
	key := this.GetString("key")
	qk := strings.Split(app.ServiceSearchKey, ",")
	user := getServiceUser(this)
	searchMap := sql.GetSearchMapValue(qk, *this.Ctx, sql.SearchMap{})
	//searchMap.Put("CreateUser", user)

	searchSql := sql.SearchSql(
		app.CloudAppService{},
		app.SelectCloudAppService,
		searchMap)
	if len(searchMap.GetData()) == 0 {
		searchSql += " where 1=1 "
	}
	if key != "" {
		q := `and (service_name like "%?%") `
		searchSql += strings.Replace(q, "?", sql.Replace(key), -1)
	}

	num, err := sql.OrderByPagingSql(searchSql,
		"create_time",
		*this.Ctx.Request,
		&data,
		app.CloudAppService{})

	result := GetServiceRunData(data, user)
	setServiceJson(this,
		util.GetResponseResult(err, this.GetString("draw"),
			result,
			sql.CountSearchMap("cloud_app_service",
				sql.SearchMap{},
				int(num), key)))

	go GoServerThread(data)
}

// 2018-02-14 18:06
// 删除服务
func DeleteK8sService(service app.CloudAppService, force string) interface{} {
	namespace := util.Namespace(service.AppName, service.ResourceName)

	serviceName := service.ServiceName
	if service.ServiceVersion != "" {
		serviceName = util.Namespace(serviceName, service.ServiceVersion)
	}
	// 先将pod数量改成0
	k8s.ScalePod(
		service.ClusterName,
		util.Namespace(service.AppName, service.ResourceName),
		serviceName,
		int32(0))
	time.Sleep(time.Second * 5)

	err := k8s.DeletelDeployment(
		namespace,
		true,
		serviceName, service.ClusterName)
	if err != nil && force == "" {
		data := util.ApiResponse(false, err.Error())
		return data
	}
	return nil
}

// json
// 删除Service 
// @router /api/service/:id:int [delete]
func (this *ServiceController) ServiceDelete() {
	searchMap := sql.GetSearchMap("ServiceId", *this.Ctx)
	searchMap.Put("CreateUser", getServiceUser(this))
	service := getService(this)
	force := this.GetString("force")

	err := DeleteK8sService(service, force)
	if err != nil {
		setServiceJson(this, err)
		return
	}

	searchMap = sql.SearchMap{}
	searchMap.Put("ServiceId", service.ServiceId)
	q := sql.DeleteSql(app.DeleteCloudAppService, searchMap)
	r, delErr := sql.Raw(q).Exec()

	data := util.DeleteResponse(delErr,
		*this.Ctx, "删除Service "+service.ServiceName,
		this.GetSession("username"),
		service.CreateUser,
		r)

	setServiceJson(this, data)
	go MakeContainerData(util.Namespace(service.AppName, service.ResourceName))
}

// 扩容或缩容服务容器
// @router /api/service/scale/:id:int [*]
func (this *ServiceController) ServiceScale() {
	service := getService(this)

	start, starte := this.GetInt("start")
	replicas, err := this.GetInt("replicas")
	if err != nil && starte != nil {
		responseData(errors.UnsupportedError(""),
			this, service.ServiceName,
			"扩容数量不对")
		return
	}

	max := int(service.ReplicasMax)
	min := int(service.ReplicasMin)
	if replicas > max || (replicas != 0 && replicas < min) {
		responseData(errors.UnsupportedError(""), this,
			service.ServiceName,
			"超过或比预期值小,最大:"+strconv.Itoa(max)+" 最小:"+strconv.Itoa(min))
		return
	}

	// 如果是启动服务,就恢复到配置好的数量
	if start > 0 {
		replicas = int(service.Replicas)
	}

	mem := service.Memory
	cpu := int64(service.Cpu)
	if service.Replicas > int64(replicas) {
		n := service.Replicas - int64(replicas)
		mem = mem - mem*n
		cpu = cpu - cpu*n
	} else {
		n := int64(replicas) - service.Replicas
		mem = mem * n
		cpu = cpu * n
	}
	status, msg := k8s.CheckQuota(
		getServiceUser(this), int64(replicas),
		int64(service.Cpu), mem,
		service.ResourceName)

	if ! status {
		logs.Error("用户超过配额", msg)
		responseData(errors.InvalidArgumentError(msg), this, service.ServiceName, msg)
		return
	}

	serviceName := service.ServiceName
	if service.ServiceVersion != "" {
		serviceName = util.Namespace(serviceName, service.ServiceVersion)
	}

	err = k8s.ScalePod(service.ClusterName,
		util.Namespace(service.AppName, service.ResourceName),
		serviceName,
		int32(replicas))

	if err != nil {
		responseData(err, this, service.ServiceName, "操作出现异常")
		return
	}

	// 更新数据库副本数量
	if err == nil {
		if replicas > 0 && int64(replicas) != service.Replicas {
			service.Replicas = int64(replicas)
			updateServiceData(service, getServiceUser(this))
		}
	}

	go updateServiceRedisCache(service)
	go MakeContainerData(
		util.Namespace(
			service.AppName,
			service.ResourceName))
	responseData(nil, this, service.ServiceName, "操作成功")
}

// @parame type 更新类型 image config port env health
// 2018-01-13 19:37
// @router /api/service/update/:id:int [post]
func (this *ServiceController) ServiceUpdate() {

	service := getService(this)
	updateType := this.GetString("type")

	// 修改端口数据
	if updateType == "port" {
		port := this.GetString("port")
		if port == "" {
			responseData(
				errors.InvalidArgumentError("数据不能为空 :port 80,8080"),
				this,
				service.ServiceName,
				"操作失败",
			)
			return
		}
		service.ContainerPort = port
	}

	// 健康检查升级
	// 2018-01-14 13:23
	if updateType == "health" {
		healthData := this.GetString("healthData")
		if healthData == "" || len(healthData) < 20 {
			responseData(errors.InvalidArgumentError("数据不能为空"), this, service.ServiceName, "操作失败")
			return
		}
		service.HealthData = healthData
	}

	// 升级环境变量
	if updateType == "env" {
		env := this.GetString("env")
		if env == "" || !strings.Contains(env, "=") || len(env) < 3 {
			responseData(errors.InvalidArgumentError("变量数据异常"), this, service.ServiceName, "操作失败")
			return
		}
		service.Envs = env
	}

	// 升级镜像版本
	if updateType == "image" {
		version := this.GetString("version")
		version = strings.TrimSpace(version)
		tags := strings.Split(service.ImageTag, ":")
		if version != tags[1] {
			service.ImageTag = tags[0] + ":" + version
		}
		//} else {
		//	responseData(errors.InvalidArgumentError("镜像版本一致"), this, service.ServiceName, "操作失败")
		//	return
		//}
		interval, err := this.GetInt("MinReady")
		if err != nil || interval > 60 || interval < 2 {
			responseData(errors.InvalidArgumentError("更新间隔错误,可选范围为:2-60"), this, service.ServiceName, "操作失败")
			return
		}
		service.MinReady = interval
	}

	// 更新内存,cpu配置
	if updateType == "config" {
		v, cpuerr := strconv.ParseFloat(this.GetString("cpu"), 32)
		service.Cpu = float32(v)
		mem, memerr := this.GetInt64("mem")
		service.Memory = mem
		if !checkParam(service, cpuerr, memerr, this) {
			return
		}

		status, msg := k8s.CheckQuota(
			getServiceUser(this), service.Replicas,
			int64(service.Cpu), service.Memory,
			service.ResourceName)

		if ! status {
			logs.Error("用户超过配额", msg)
			responseData(errors.InvalidArgumentError(msg), this, service.ServiceName, msg)
			return
		}
	}

	user := util.GetUser(this.GetSession("username"))
	err := ExecUpdate(service, updateType, user)
	if err == nil {
		responseData(err, this, service.ServiceName, "操作成功")
		return
	}
	responseData(err, this, service.ServiceName, "操作失败")
}
