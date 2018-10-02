package app

import (
	"cloud/sql"
	"cloud/models/app"
	"encoding/json"
	"cloud/k8s"
	"strconv"
	"time"
	"cloud/util"
	"strings"
	"cloud/cache"
	"github.com/astaxie/beego/logs"
	"cloud/userperm"
)

// 2018-02-19 10:40
// 获取应用数据
func getApp(this *AppController) (app.CloudApp, sql.SearchMap) {
	searchMap := sql.GetSearchMap("AppId", *this.Ctx)
	d := app.CloudApp{}
	q := sql.SearchSql(d, app.SelectCloudApp, searchMap)
	sql.Raw(q).QueryRow(&d)
	return d, searchMap
}

// 将数据写入到应用表中
// 2018-01-16 17:23
func saveAppData(service app.CloudAppService) {
	v, _ := json.Marshal(service)
	data := app.CloudApp{}
	json.Unmarshal(v, &data)
	insert := sql.InsertSql(data, app.InsertCloudApp)
	sql.Raw(insert).Exec()
}

// 2018-02-11 08:40
// 从缓存中读取应用状态数据
func getCacheAppData(data []app.CloudApp) []k8s.CloudApp {
	var cloudApps []k8s.CloudApp
	for _, d := range data {
		if cache.AppCacheErr == nil {
			redisAppData := cache.AppCache.Get(strconv.FormatInt(d.AppId, 10))
			if redisAppData != nil {
				tempAppData := k8s.CloudApp{}
				status := util.RedisObj2Obj(redisAppData, &tempAppData)
				if status {
					cloudApps = append(cloudApps, tempAppData)
				}
			}
		}
	}
	return cloudApps
}

// 2018-02-22 18:18
// 任务计划app数据缓存
func CacheAppData() {
	data := make([]app.CloudApp, 0)
	sql.Raw(app.SelectCloudApp).QueryRows(&data)
	getK8sAppData(data)
}

// 2018-02-27 11:21
// 将应用数据写入到redis
func putAppDataToRedis(id int64, app interface{}) {
	key := strconv.FormatInt(id, 10)
	cache.AppCache.Put(key, util.ObjToString(app), time.Minute*10)
}

// 2018-09-04 08:06
// 生产应用缓存数据
func getK8sAppData(data []app.CloudApp) {

	for _, d := range data {
		//putAppDataToRedis(d.AppId, d)
		namespace := util.Namespace(d.AppName, d.ResourceName)
		c, err := k8s.GetClient(d.ClusterName)
		if err != nil {
			logs.Error("获取客户端失败", err.Error())
			continue
		}

		cloudAppData := k8s.GetDeploymentApp(c, namespace, "")
		for _, app := range cloudAppData {
			app.ResourceName = d.ResourceName
			app.AppId = d.AppId
			app.Entname = d.Entname
			app.ClusterName = d.ClusterName
			app.CreateUser = d.CreateUser
			putAppDataToRedis(app.AppId, app)
		}

	}
}

// 获取登录用户
func getUser(this *AppController) string {
	return util.GetUser(this.GetSession("username"))
}

// 2018-02-09 15:55
// 获取所有应用和服务的名称,在流水线判断应用服务是否存在
func GetAppServiceDataMap() util.Lock {
	lock := util.Lock{}
	data := make([]app.CloudAppServiceName, 0)

	searchSql := sql.SearchSql(
		app.CloudAppService{},
		app.SelectAppServiceName,
		sql.SearchMap{})

	sql.Raw(searchSql).QueryRows(&data)
	for _, v := range data {
		lock.Put(v.ClusterName+v.AppName+v.ServiceName, 1)
	}
	return lock
}

// 生成镜像tag
// 2018-02-01 08:25
func makeImageTags(tag string) string {
	var html string
	tags := strings.Split(tag, ",")
	for _, v := range tags {
		html += util.GetSelectOptionName(v)
	}
	return html
}

// 2018-02-26 09:32
// 获取重建的应用信息
func getRedeployApp(v string, user string) ([]app.CloudAppName, bool) {
	searchMap := sql.SearchMap{}
	searchMap.Put("AppId", v)
	searchMap.Put("CreateUser", user)
	r := getAppDataQ(searchMap)
	if len(r) == 0 {
		return []app.CloudAppName{}, false
	}
	return r, true
}

// 2018-02-26 09:54
// 获取重建应用的服务信息
func getRedeployService(v string, user string) ([]app.CloudAppService, bool) {
	logs.Info("开始重建服务", util.ObjToString(v))
	data, status := getRedeployApp(v, user)
	if ! status {
		return []app.CloudAppService{}, false
	}
	searchMap := sql.SearchMap{}
	searchMap.Put("ClusterName", data[0].ClusterName)
	searchMap.Put("Entname", data[0].Entname)
	searchMap.Put("AppName", data[0].AppName)
	serviceData := getServiceData(searchMap, app.SelectCloudAppService)
	logs.Info("获取到服务数据", util.ObjToString(serviceData))
	if len(serviceData) == 0 {
		return []app.CloudAppService{}, false
	}
	return serviceData, true
}

// 获取应用名称信息
func getAppDataQ(searchMap sql.SearchMap) []app.CloudAppName {
	userInterface := searchMap.Get("CreateUser")
	searchMap.Put("CreateUser", nil)
	data := make([]app.CloudAppName, 0)
	dataApp := make([]app.CloudAppName, 0)
	searchSql := sql.SearchSql(app.CloudAppName{}, app.GetAppName, searchMap)
	logs.Info("searchSql", searchSql, searchMap)
	sql.Raw(searchSql).QueryRows(&data)

	if userInterface != nil {
		user := userInterface.(string)
		permApp := userperm.GetResourceName("应用", user)
		for _, v := range data {
			// 不是自己创建的才检查
			if v.CreateUser != user && user != "admin" {
				if ! userperm.CheckPerm(v.AppName, v.ClusterName, v.Entname, permApp) {
					continue
				}
			}
			dataApp = append(dataApp, v)
		}
	}
	return data
}

// 2018-02-03 21:44
// 获取选项卡
func GetAppHtml(cluster string, username string) string {
	data := getAppData("", cluster, username)
	var html string
	for _, v := range data {
		html += util.GetSelectOptionName(v.AppName)
	}
	return html
}

// 2018-02-27 11:45
// 加载应用数据
func selectAppData(searchMap sql.SearchMap) []app.CloudApp {
	data := make([]app.CloudApp, 0)
	searchSql := sql.SearchSql(app.CloudApp{}, app.SelectCloudApp, searchMap)
	sql.Raw(searchSql).QueryRows(&data)
	return data
}

// 查询某个服务的数据
func getAppData(name string, cluster string, username string) []app.CloudApp {

	searchMap := sql.GetSearchMapV("ClusterName", cluster, "CreateUser", username)
	if name != "" {
		searchMap.Put("AppName", name)
	}
	return selectAppData(searchMap)
}

// 2018-02-13 15:46
// 获取应用选择项
func GetAppSelect(searchMap sql.SearchMap) string {
	data := getAppDataQ(searchMap)
	var opt = "<option>--请选择--</option>"
	for _, v := range data {
		opt += util.GetSelectOptionName(v.AppName)
	}
	return opt
}
