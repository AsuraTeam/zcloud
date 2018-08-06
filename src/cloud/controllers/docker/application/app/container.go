package app

import (
	"time"
	"cloud/k8s"
	"cloud/sql"
	"cloud/models/app"
	"cloud/util"
	"database/sql/driver"
	"strings"
	"cloud/cache"
	"k8s.io/client-go/kubernetes"
)

// 容器详情页面
// 2018-01-16 08:34
// @router /application/container/detail/:hi:string [get]
func (this *AppController) ContainerDetail() {
	data := app.CloudContainer{}
	name := this.Ctx.Input.Param(":hi")
	searchMap := sql.GetSearchMapV("ContainerName", name)
	searchSql := sql.SearchSql(app.CloudContainer{}, app.SelectCloudContainer, searchMap)
	sql.Raw(searchSql).QueryRow(&data)
	v := getRedisContainer(data, "", "")
	if v.ContainerId != 0 {
		this.Data["data"] = v
	} else {
		this.Data["data"] = data
	}
	this.TplName = "application/container/detail.html"
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

	key := sql.MKeyV("Entname", "ServiceName", "AppName")

	searchMap := sql.GetSearchMapValue(key,
		*this.Ctx,
		sql.SearchMap{})

	searchMap.Put("CreateUser", getUser(this))

	searchSql := sql.SearchSql(app.CloudContainer{},
		app.SelectCloudContainer,
		searchMap)

	if search != "" {
		q := ` and container_name like "%?%"`
		searchSql += strings.Replace(q, "?", sql.Replace(search), -1)
	}

	sql.OrderByPagingSql(searchSql, "container_id",
		*this.Ctx.Request,
		&data,
		app.CloudContainer{})

	datas := make([]interface{}, 0)
	for _, cv := range data {
		key := cv.AppName + cv.ContainerName
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
			sql.GetSearchMapV("CreateUser", getUser(this)),
			len(datas),
			search),
		this.GetString("draw"))

	SetAppDataJson(this, r)
	go MakeContainerData("")
}

// 2018-01-16 8:48
// 更新或写入到数据库
func writeToDb(appDatas util.Lock, appDatasDb util.Lock) {
	for _, d := range appDatas.GetData() {
		o := d.(app.CloudContainer)
		sName := o.AppName + o.ContainerName
		if _, ok := appDatasDb.Get(sName); !ok {
			o.Events = ""
			q := sql.InsertSql(o, app.InsertCloudContainer)
			sql.Raw(q).Exec()
		}
		if cache.ContainerCacheErr == nil {
			cache.ContainerCache.Put(sName, util.ObjToString(o), time.Second*3600)
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

// 2018-01-15 15:25
// 通过任务计划方式获取数据
func MakeContainerData(namespace string) {
	if !util.WriteLock("last_update", &LockContainerUpdate, 10) {
		return
	}

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
		sname := namespace + d.ServiceName

		if _, ok := lockData.Get(sname); !ok {
			c, _ := k8s.GetClient(d.ClusterName)
			appData := k8s.GetContainerStatus(namespace, c)
			for _, all := range appData {
				all = setAppData(all, d, c)
				appDataLock.Put(all.AppName+all.ContainerName, all)
				containerDatas.Put(all.AppName+all.ContainerName, "1")
				lockData.Put(sname, "1")
			}
		}
	}

	// 要删除的数据
	deleteData := util.Lock{}
	appDatasDb := util.Lock{}
	datas := make([]app.CloudContainerName, 0)
	containerSql := sql.SearchSql(app.CloudContainer{}, app.SelectCloudContainer, sql.SearchMap{})
	sql.Raw(containerSql).QueryRows(&datas)
	for _, d := range datas {
		sname := d.AppName + d.ContainerName
		appDatasDb.Put(sname, "1") // 将这个名称写成真
		// 如果k8s里的容器没有的话就删除掉
		if _, ok := containerDatas.Get(sname); !ok {
			deleteData.Put(sname, d)
		}
	}

	// 更新或插入数据
	go writeToDb(appDataLock, appDatasDb)
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

func setAppData(all app.CloudContainer, d app.CloudAppService, c kubernetes.Clientset) app.CloudContainer {
	namespace := util.Namespace(d.AppName, d.ResourceName)

	all.ResourceName = d.ResourceName
	all.ClusterName = d.ClusterName
	all.AppName = d.AppName
	all.CreateUser = d.CreateUser
	all.Entname = d.Entname
	if time.Now().Unix()-util.TimeToStamp(d.CreateTime) < 1800 {
		events := k8s.GetEvents(namespace, all.ContainerName, c)
		all.Events = util.ObjToString(events)
	} else {
		v := getRedisContainer(app.CloudContainer{}, d.AppName, all.ContainerName)
		if v.Events != "" {
			all.Events = v.Events
		}
	}
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
		v.Process = k8s.Exec(v.ClusterName, v.ContainerName, util.Namespace(data.AppName,data.ResourceName), data.ServiceName, cmd)
	}
	return v
}
