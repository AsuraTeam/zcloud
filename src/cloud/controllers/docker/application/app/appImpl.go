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
)

// 2018-02-19 10:40
// 获取应用数据
func getApp(this *AppController) (app.CloudApp, sql.SearchMap) {
	searchMap := sql.GetSearchMap("AppId", *this.Ctx)
	searchMap.Put("CreateUser", getUser(this))
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
				if  status {
					cloudApps = append(cloudApps, tempAppData)
				}
			}
		}
	}
	return cloudApps
}

// 2018-02-22 18:18
// 任务计划app数据缓存
func CacheAppData()  {
	data := make([]app.CloudApp, 0)
	sql.Raw(app.SelectCloudApp).QueryRows(&data)
	getK8sAppData(data)
}

// 2018-02-27 11:21
// 将应用数据写入到redis
func putAppDataToRedis(id int64, app interface{})  {
	key := strconv.FormatInt(id, 10)
	cache.AppCache.Put(key, util.ObjToString(app), time.Second * 86400)
}

func getK8sAppData(data []app.CloudApp)  {
	allData := util.Lock{}
	for _, d := range data {
		putAppDataToRedis(d.AppId, d)
		namespace := util.Namespace(d.AppName, d.ResourceName)
		number := 0
		if _, ok := allData.Get(namespace); !ok {
			c, _ := k8s.GetClient(d.ClusterName)
			cloudAppData := k8s.GetDeploymentApp(c, namespace, "")
			for _, app := range cloudAppData {
				app.ResourceName = d.ResourceName
				app.AppId = d.AppId
				app.Entname = d.Entname
				if cache.AppCacheErr == nil {
					putAppDataToRedis(app.AppId, app)
				}
				allData.Put(namespace, "1")
				number += app.ContainerNumber
			}
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