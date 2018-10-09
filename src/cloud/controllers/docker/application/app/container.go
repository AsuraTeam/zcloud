package app

import (
	"cloud/k8s"
	"cloud/sql"
	"cloud/models/app"
	"cloud/util"
	"database/sql/driver"
	"strings"
	"cloud/cache"
	"k8s.io/client-go/kubernetes"
	"cloud/models/registry"
	registry2 "cloud/controllers/image"
	"github.com/astaxie/beego/logs"
	"cloud/userperm"
	"time"
)

// 获取详情数据
func getContainerData(this *AppController) app.CloudContainer {
	data := app.CloudContainer{}
	name := this.Ctx.Input.Param(":hi")
	searchMap := sql.GetSearchMapV("ContainerName", name)
	searchSql := sql.SearchSql(app.CloudContainer{}, app.SelectCloudContainer, searchMap)
	sql.Raw(searchSql).QueryRow(&data)
	return data
}

// 容器详情页面
// 2018-01-16 08:34
// @router /application/container/detail/:hi:string [get]
func (this *AppController) ContainerDetail() {
	data := getContainerData(this)
	v := getRedisContainer(data, "", "")
	if v.ContainerId != 0 {
		this.Data["data"] = v
	} else {
		this.Data["data"] = data
	}
	this.TplName = "application/container/detail.html"
}

// 容器镜像提交页面
// 2018-08-21 13:34
// @router /application/container/image [get]
func (this *AppController) ContainerImage() {
	data := app.CloudContainer{}
	id, err := this.GetInt("id")
	if err != nil {
		this.Ctx.WriteString("参数错误")
		return
	}
	searchMap := sql.SearchMap{}
	searchMap.Put("ContainerId", id)
	searchSql := sql.SearchSql(app.CloudContainer{}, app.SelectCloudContainer, searchMap)
	sql.Raw(searchSql).QueryRow(&data)
	this.Data["data"] = data
	if len(strings.Split(data.Image, "/")) < 2 {
		SetAppDataJson(this, "容器没有镜像,不能提交")
		return
	}
	this.Data["Group"] = strings.Split(data.Image, "/")[1]
	this.Data["ItemName"] = strings.Split(strings.Join(strings.Split(data.Image, "/")[2:], "/"), ":")[0]
	this.Data["baseImage"] = registry2.GetBaseImageSelect()
	this.TplName = "application/container/image.html"
}

// 2018-08-21 14:33
// 容器提交镜像
// @router /api/container/commit/:id:int [delete]
func (this *AppController) ContainerCommit() {
	data := app.CloudContainer{}
	id := this.Ctx.Input.Param(":id")
	searchMap := sql.GetSearchMapV("ContainerId", id)
	searchSql := sql.SearchSql(app.CloudContainer{}, app.SelectCloudContainer, searchMap)
	sql.Raw(searchSql).QueryRow(&data)
	sync := registry.CloudImageSync{}
	sync.ClusterName = data.ClusterName
	sync.ImageName = data.Image
	name := strings.Split(sync.ImageName, "/")
	if len(name) > 2 {
		registryData := registry2.GetRegistryServerCluster(strings.Split(name[0], ":")[0], data.ClusterName)
		sync.Registry = registryData.Name

		param := getImageCommitParam(sync, getUser(this))
		param.ContainerId = data.ContainerName
		param.ServerAddress = data.ServerAddress
		param.Version = this.GetString("Version")
		param.ItemName = strings.Split(data.Image, "/")[1] + "/" + this.GetString("ItemName")
		logs.Info("仓库数据信息", util.ObjToString(sync), util.ObjToString(registryData), util.ObjToString(param))
		k8s.ImageCommit(data.ClusterName, param, this.GetString("BaseImage"))
	}
	SetAppDataJson(this, util.ApiResponse(true, "保存成功,正在处理中"))
}

// 2018-01-16 12:20
// 删除容器
// @router /api/container/:id:int [delete]
func (this *AppController) ContainerDelete() {

	data := app.CloudContainer{}
	searchMap := sql.GetSearchMap("ContainerId", *this.Ctx)
	searchSql := sql.SearchSql(app.CloudContainer{}, app.SelectCloudContainer, searchMap)
	sql.Raw(searchSql).QueryRow(&data)

	cl, err := k8s.GetClient(data.ClusterName)
	namespace := util.Namespace(data.AppName, data.ResourceName)

	err = k8s.DeletePod(namespace, data.ContainerName, cl)
	if err == nil {
		r, err := sql.Raw(sql.DeleteSql(app.DeleteCloudContainer, searchMap)).Exec()
		data := util.DeleteResponse(err, *this.Ctx, "删除容器 "+data.ContainerName, this.GetSession("username"), data.ClusterName, r)
		SetAppDataJson(this, data)
		return
	}

	r := driver.ResultNoRows
	json := util.DeleteResponse(err, *this.Ctx, "删除容器失败 "+data.ContainerName, this.GetSession("username"), data.ContainerName, r)
	SetAppDataJson(this, json)
	go MakeContainerData(namespace)
}

// 获取容器运行情况
// 2018-01-15 15:11
// @router /api/container [get]
func (this *AppController) ContainerData() {
	data := make([]app.CloudContainer, 0)
	search := this.GetString("search")

	key := sql.MKeyV("Entname", "Service", "AppName")

	searchMap := sql.GetSearchMapValue(key,
		*this.Ctx,
		sql.SearchMap{})

	searchSql := sql.SearchSql(app.CloudContainer{},
		app.SelectCloudContainer,
		searchMap)

	searchSql = sql.GetWhere(searchSql, searchMap)
	if search != "" {
		q := ` and container_name like "%?%"`
		searchSql += strings.Replace(q, "?", sql.Replace(search), -1)
	}

	sql.OrderByPagingSql(searchSql, "create_time",
		*this.Ctx.Request,
		&data,
		app.CloudContainer{})

	user := getUser(this)
	perm := userperm.GetResourceName("服务", user)
	permApp := userperm.GetResourceName("应用", user)
	datas := make([]interface{}, 0)
	for _, cv := range data {
		key := cv.AppName + cv.ContainerName
		d := cv
		// 不是自己创建的才检查
		if d.CreateUser != user {
			service := strings.Replace(d.ServiceName, "--1", "", -1)
			service = strings.Replace(service, "--2", "", -1)

			if ! userperm.CheckPerm(d.AppName+";"+d.ResourceName+";"+service, d.ClusterName, d.Entname, perm) && len(user) > 0 {
				if ! userperm.CheckPerm(d.AppName, d.ClusterName, d.Entname, permApp) {
					continue
				}
			}
		}
		r := cache.ContainerCache.Get(key)
		var v interface{}
		status := util.RedisObj2Obj(r, &v)
		if status {
			v.(map[string]interface{})["ContainerId"] = cv.ContainerId
			v.(map[string]interface{})["CreateTime"] = util.GetMinTime(cv.CreateTime)
			datas = append(datas, v)
		} else {
			datas = append(datas, cv)
		}
	}

	r := util.ResponseMap(datas,
		sql.CountSearchMap("cloud_container",
			sql.SearchMap{},
			len(datas),
			search),
		this.GetString("draw"))

	SetAppDataJson(this, r)
	go MakeContainerData("")
}

// 2018-01-16 8:48
// 更新或写入到数据库
func writeToDb(appData util.Lock, appDatasDb util.Lock) {
	for _, d := range appData.GetData() {
		o := d.(app.CloudContainer)
		sName := o.AppName + o.ContainerName
		if _, ok := appDatasDb.Get(sName); !ok {
			o.Events = ""
			q := sql.InsertSql(o, app.InsertCloudContainer)
			sql.Raw(q).Exec()
		}
	}
}

// 2018-01-16 08:51
// 删除数据库多余的数据
func deleteDbContainer(deleteData util.Lock) {
	// 删除数据库中的内容
	for _, d := range deleteData.GetData() {
		v := d.(app.CloudContainerName)

		deleteSql := sql.DeleteSql(
			app.DeleteCloudContainer,

			sql.GetSearchMapV("ContainerName",
				v.ContainerName,
				"AppName",
				v.AppName))

		go sql.Raw(deleteSql).Exec()
	}
}

// 2018-09-04 18:19
// 获取容器日志
// @router /api/container/logs/:hi(.*) [get]
func (this *AppController) GetDockerLogs() {
	data := getContainerData(this)
	cl, err := k8s.GetClient(data.ClusterName)
	if err != nil {
		this.Ctx.WriteString(err.Error())
		return
	}
	line, err := this.GetInt64("tailLine")
	if err != nil {
		line = 5000
	}
	log := k8s.GetJobLogs(cl, data.ContainerName, util.Namespace(data.AppName, data.ResourceName), line)
	logs.Info(log, data.ClusterName, data.AppName, data.ResourceName)
	this.Ctx.WriteString(log)
}

// 2018-01-15 15:25
// 通过任务计划方式获取数据
func MakeContainerData(namespace string) {
	if !util.WriteLock("last_update", &LockContainerUpdate, 10) {
		return
	}
	logs.Info("生成容器数据")
	searchMap := sql.SearchMap{}
	if namespace != "" {
		searchMap.Put("Namespace", namespace)
	}
	data := make([]app.CloudAppService, 0)
	searchSql := sql.SearchSql(
		app.CloudAppService{},
		app.SelectCloudAppService,
		searchMap)
	sql.Raw(searchSql).QueryRows(&data)

	containerDatas := util.Lock{}
	appDataLock := util.Lock{}
	lockData := util.Lock{}
	for _, d := range data {
		namespace := util.Namespace(d.AppName, d.ResourceName)
		sName := namespace + d.ServiceName + d.ClusterName
		if _, ok := lockData.Get(sName); !ok {
			c, err := k8s.GetClient(d.ClusterName)
			if err != nil {
				logs.Error("获取客户端错误", err.Error())
				continue
			}
			appData := k8s.GetContainerStatus(namespace, c)
			for _, all := range appData {
				all = setAppData(all, d, c, util.Namespace(d.ServiceName, d.ServiceVersion))
				cache.ContainerCache.Put(all.AppName+all.ContainerName, util.ObjToString(all), time.Second*3600)
				appDataLock.Put(all.AppName+all.ContainerName, all)
				containerDatas.Put(all.AppName+all.ContainerName, "1")
				lockData.Put(sName, "1")
			}
		}
	}

	// 要删除的数据
	deleteData := util.Lock{}
	appDataDb := util.Lock{}
	dataS := make([]app.CloudContainerName, 0)
	containerSql := sql.SearchSql(app.CloudContainer{}, app.SelectCloudContainer, sql.SearchMap{})
	sql.Raw(containerSql).QueryRows(&dataS)
	for _, d := range dataS {
		sName := d.AppName + d.ContainerName
		appDataDb.Put(sName, "1") // 将这个名称写成真
		// 如果k8s里的容器没有的话就删除掉
		if _, ok := containerDatas.Get(sName); !ok {
			deleteData.Put(sName, d)
		}
	}

	// 更新或插入数据
	go writeToDb(appDataLock, appDataDb)
	// 删除数据
	go deleteDbContainer(deleteData)
}

func SetAppDataJson(this *AppController, data interface{}) {
	this.Data["json"] = data
	this.ServeJSON(false)
}

// 2018-02-28 09:52
// 填充容器数据
var cmd = []string{"ps", "aux"}

func setAppData(all app.CloudContainer, d app.CloudAppService, c kubernetes.Clientset, serviceName string) app.CloudContainer {
	namespace := util.Namespace(d.AppName, d.ResourceName)
	all.ResourceName = d.ResourceName
	all.ClusterName = d.ClusterName
	all.AppName = d.AppName
	all.CreateUser = d.CreateUser
	all.Entname = d.Entname
	all.ServiceName = serviceName
	events := k8s.GetEvents(namespace, all.ContainerName, c)
	all.Events = util.ObjToString(events)
	return all
}

// 2018-02-28 09:57
// 从redis中获取容器数据
func getRedisContainer(data app.CloudContainer, appName string, containerName string) app.CloudContainer {
	var v app.CloudContainer
	var r interface{}
	if appName == "" {
		r = cache.ContainerCache.Get(data.AppName + data.ContainerName)
	} else {
		r = cache.ContainerCache.Get(appName + containerName)
	}
	status := util.RedisObj2Obj(r, &v)
	if status && data.ContainerId > 0 {
		v.ContainerId = data.ContainerId
		v.CreateTime = util.GetMinTime(data.CreateTime)
		v.Process = k8s.Exec(v.ClusterName, v.ContainerName, util.Namespace(data.AppName, data.ResourceName), data.ServiceName, cmd)
	}
	return v
}

// 2018-08-22 11:17
// 获取容器提交参数
func getImageCommitParam(d registry.CloudImageSync, user string) k8s.ImagePushParam {
	registryData := registry2.GetRegistryServerMap()

	imagePushParam := k8s.ImagePushParam{
		RegistryGroup: "",
		ItemName:      d.ImageName,
		Version:       d.Version,
		CreateTime:    util.GetDate(),
		User:          user,
	}
	reg1, ok := registryData.Get(d.ClusterName + d.Registry)
	if ok {
		reg1data := reg1.(registry.CloudRegistryServer)
		servers := strings.Split(reg1data.ServerAddress, ":")
		imagePushParam.Registry1Domain = reg1data.ServerDomain
		imagePushParam.Registry1Auth = util.Base64Encoding(reg1data.Admin + ":" + util.Base64Decoding(reg1data.Password))
		if len(servers) == 2 {
			imagePushParam.Registry1Ip = servers[0]
			imagePushParam.Registry1Port = servers[1]
		}
	}
	return imagePushParam
}
