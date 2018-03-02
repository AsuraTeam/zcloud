package monitor

import (
	"cloud/sql"
	"cloud/util"
	"github.com/astaxie/beego"
	"strings"
	"cloud/models/monitor"
	"cloud/controllers/ent"
	"cloud/controllers/docker/application/app"
	"cloud/k8s"
	"strconv"
	"time"
	"cloud/cache"
	app2 "cloud/models/app"
	"github.com/astaxie/beego/logs"
	"math/rand"
)

type AutoScaleController struct {
	beego.Controller
}

// 自动伸缩管理入口页面
// @router /monitor/scale/list [get]
func (this *AutoScaleController) AutoScaleList() {
	this.TplName = "monitor/scale/list.html"
}

// 自动伸缩管理入口页面
// @router /monitor/scale/logs [get]
func (this *AutoScaleController) AutoScaleLogs() {
	this.TplName = "monitor/scale/logs.html"
}

// 2018-02-19 18:19
// 自动伸缩管理添加页面
// @router /monitor/scale/add [get]
func (this *AutoScaleController) AutoScaleAdd() {
	id := this.GetString("ScaleId")
	update := monitor.GetQueryParamDefault()
	entData := ent.GetEntnameSelect()
	appData := app.GetAppSelect(sql.GetSearchMapV("CreateUser", getUser(this)))

	var entHtml string
	var serviceHtml string
	var appHtml string

	this.Data["system"] = "checked"
	// 更新操作
	if id != "0" {
		searchMap := sql.GetSearchMap("ScaleId", *this.Ctx)
		q := sql.SearchSql(monitor.CloudAutoScale{}, monitor.SelectCloudAutoScale, searchMap)
		sql.Raw(q).QueryRow(&update)

		entHtml = util.GetSelectOptionName(update.Entname)
		appHtml = util.GetSelectOptionName(update.AppName)

		searchMap = sql.GetSearchMapV(
			"Entname", update.Entname,
			"CreateUser", getUser(this),
			"AppName", update.AppName,
			"ClusterName", update.ClusterName)
		serviceHtml = util.GetSelectOptionName(update.ServiceName) +
			app.GetServiceHtml(searchMap)

		this.Data["cluster"] = util.GetSelectOptionName(update.ClusterName)
		this.Data["apps"] = appHtml + appData
		this.Data["service"] = serviceHtml
	}

	this.Data["dataSource"] = getDataSourceSelect(update.DataSource)
	this.Data["metric"] = getSystemMetricSelect(update.MetricName)
	this.Data["entname"] = entHtml + entData
	this.Data["data"] = update
	this.TplName = "monitor/scale/add.html"
}

// 获取自动伸缩数据
// 2018-01-20 12:56
// router /api/monitor/scale [get]
func (this *AutoScaleController) AutoScaleData() {
	// 自动伸缩数据
	data := []monitor.CloudAutoScale{}
	q := sql.SearchSql(monitor.CloudAutoScale{}, monitor.SelectCloudAutoScale, sql.SearchMap{})
	sql.Raw(q).QueryRows(&data)
	setAutoScaleJson(this, data)
}

// string
// 自动伸缩保存
// @router /api/monitor/scale [post]
func (this *AutoScaleController) AutoScaleSave() {
	d := monitor.CloudAutoScale{}
	err := this.ParseForm(&d)
	if err != nil {
		this.Ctx.WriteString("参数错误" + err.Error())
		return
	}
	util.SetPublicData(d, getUser(this), &d)
	q := sql.InsertSql(d, monitor.InsertCloudAutoScale)
	if d.ScaleId > 0 {
		searchMap := sql.SearchMap{}
		searchMap.Put("ScaleId", d.ScaleId)
		q = sql.UpdateSql(d, monitor.UpdateCloudAutoScale, searchMap, monitor.UpdateAutoScaleExclude)
	}
	_, err = sql.Raw(q).Exec()

	data, msg := util.SaveResponse(err, "名称已经被使用")
	util.SaveOperLog(this.GetSession("username"), *this.Ctx, "保存自动伸缩配置 "+msg, d.ClusterName)
	setAutoScaleJson(this, data)
}

// 2018-02-19 18:22
// 获取自动伸缩信息
func GetAutoScaledata() []monitor.CloudAutoScale {
	// 自动伸缩数据
	data := []monitor.CloudAutoScale{}
	q := sql.SearchSql(monitor.CloudAutoScale{}, monitor.SelectCloudAutoScale, sql.SearchMap{})
	sql.Raw(q).QueryRows(&data)
	return data
}

// 2018-02-20 17:42
// 自动伸缩日志
// @router /api/monitor/scale/logs [get]
func (this *AutoScaleController) AutoScaleLogsData() {
	data := []k8s.CloudAutoScaleLog{}
	searchMap := sql.SearchMap{}
	key := this.GetString("search")
	searchSql := sql.SearchSql(
		k8s.CloudAutoScaleLog{},
		monitor.SelectCloudAutoScaleLog, searchMap)

	if key != "" {
		key = sql.Replace(key)
		searchSql += strings.Replace(monitor.SelectAutoScaleLogWhere, "?", key, -1)
	}

	num, _ := sql.OrderByPagingSql(searchSql,
		"log_id",
		*this.Ctx.Request,
		&data,
		k8s.CloudAutoScaleLog{})

	r := util.ResponseMap(data,
		sql.Count("cloud_auto_scale_log", int(num), key),
		this.GetString("draw"))
	setAutoScaleJson(this, r)
}

// 自动伸缩数据
// @router /api/monitor/scale [get]
func (this *AutoScaleController) AutoScaleDatas() {
	data := []monitor.CloudAutoScale{}
	searchMap := sql.SearchMap{}
	key := this.GetString("search")
	searchSql := sql.SearchSql(
		monitor.CloudAutoScale{},
		monitor.SelectCloudAutoScale, searchMap)

	if key != "" {
		key = sql.Replace(key)
		searchSql += strings.Replace(monitor.SelectCloudAutoScaleWhere, "?", key, -1)
	}

	num, _ := sql.OrderByPagingSql(searchSql,
		"scale_id",
		*this.Ctx.Request,
		&data,
		monitor.CloudAutoScale{})

	r := util.ResponseMap(data,
		sql.Count("cloud_auto_scale", int(num), key),
		this.GetString("draw"))
	setAutoScaleJson(this, r)
}

// json
// 删除自动伸缩
// 2018-02-05 18:05
// @router /api/monitor/scale/:id:int [delete]
func (this *AutoScaleController) AutoScaleDelete() {
	searchMap := sql.SearchMap{}
	id := this.Ctx.Input.Param(":id")
	searchMap.Put("ScaleId", id)
	scaleData := monitor.CloudAutoScale{}

	q := sql.SearchSql(scaleData, monitor.SelectCloudAutoScale, searchMap)
	sql.Raw(q).QueryRow(&scaleData)

	q = sql.DeleteSql(monitor.DeleteCloudAutoScale, searchMap)
	r, err := sql.Raw(q).Exec()
	data := util.DeleteResponse(
		err,
		*this.Ctx,
		"删除自动伸缩"+scaleData.ServiceName,
		this.GetSession("username"),
		scaleData.ClusterName,
		r)
	setAutoScaleJson(this, data)
}

func setAutoScaleJson(this *AutoScaleController, data interface{}) {
	this.Data["json"] = data
	this.ServeJSON(false)
}

func getUser(this *AutoScaleController) string {
	return util.GetUser(this.GetSession("username"))
}

// 2018-02-20 08:27
func getHtml(lock util.Lock, name string) string {
	var data string
	if _, ok := lock.Get(name); ok {
		data = util.GetSelectOption(lock.GetVString(name), name, "")
	}
	for k, v := range lock.GetData() {
		data += util.GetSelectOption(v.(string), k, "")
	}
	return data
}

// 2018-02-20 08:01
// 默认系统指标选择
func getSystemMetricSelect(name string) string {
	lock := util.Lock{}
	lock.Put("cpu", "cpu")
	lock.Put("memory", "内存")
	lock.Put("trafficInput", "网络流量输入")
	lock.Put("trafficOutput", "网络流量输出")
	return getHtml(lock, name)
}

// 2018-02-20 08:21
// 监控数据源
func getDataSourceSelect(name string) string {
	lock := util.Lock{}
	lock.Put("prometheus", "Prometheus")
	lock.Put("es", "ElasticSearch")
	return getHtml(lock, name)
}

// 2018-02-20 10:20
// 缓存服务数据
func cacheServiceData(clusterName string, appName string, serviceName string) app2.CloudAppService {
	data := app2.CloudAppService{}
	key := clusterName + appName + serviceName
	r := cache.ServiceDataCache.Get(key)
	status := util.RedisObj2Obj(r, &data)
	if !status || data.ServiceId == 0 {
		data = app.GetServiceData(serviceName, clusterName, appName)
		cache.ServiceDataCache.Put(key, util.ObjToString(data), time.Minute*30)
	}
	return data
}

// 2018-02-20 11:02
// 缓存prometheus服务信息
func cachePrometheus(data monitor.CloudAutoScale) monitor.PrometheusServer {
	key := data.ClusterName + data.AppName + data.ServiceName + "prometheus"
	r := cache.AutoScaleCache.Get(key)
	server := monitor.PrometheusServer{}
	status := util.RedisObj2Obj(r, &server)
	if ! status {
		master, _ := k8s.GetMasterIp(data.ClusterName)
		client, _ := k8s.GetClient(data.ClusterName)
		serviceData := k8s.GetServerPort(client, "monitoring", "prometheus")
		logs.Info(util.ObjToString(serviceData))
		if _, ok := serviceData.Get("nodePort"); ok {
			server.Port = serviceData.GetVString("nodePort")
			server.Host = master
			logs.Info(serviceData)
			// 2小时缓存
			cache.AutoScaleCache.Put(key, util.ObjToString(server), time.Hour * 2)
			return server
		}
	}
	return server
}

// 2018-02-20 10:02
// 执行监控检查和操作扩容
func execAutoScale(data monitor.CloudAutoScale) {
	serviceData := cacheServiceData(data.ClusterName, data.AppName, data.ServiceName)
	param := k8s.QueryParam{}
	util.MergerStruct(data, &param)
	param.Namespace = util.Namespace(serviceData.AppName, serviceData.ResourceName)
	prometheusServer := cachePrometheus(data)
	if prometheusServer.Port == "" {
		logs.Error("获取监控Prometheus服务器失败", prometheusServer)
		return
	}
	param.Host = prometheusServer.Host
	param.Port = prometheusServer.Port
	step, err := strconv.ParseInt(data.Step, 10, 64)
	if err == nil {
		param.Start = strconv.FormatInt(time.Now().Unix()-data.LastCount*step, 10)
	}
	param.End = strconv.FormatInt(time.Now().Unix(), 10)
	k8s.ParseMonitorData(param)
}

// 2018-02-20 9:52
// 执行扩容监控和操作
// 通过任务计划调用
func CronAutoScale() {
	key := "auto_scale_"
	// 缓存5分钟
	data := []monitor.CloudAutoScale{}
	r := cache.AutoScaleCache.Get(key)
	util.RedisObj2Obj(r, &data)
	if len(data) == 0 {
		data = GetAutoScaledata()
		cache.AutoScaleCache.Put(key, util.ObjToString(data), time.Minute * 5)
	}
	for _, v := range data {
		key := "cloud_auto_scale_lock" + strconv.FormatInt(v.ScaleId, 10)
		r := cache.AutoScaleCache.Get(key)
		if r != nil {
			//continue
		}
		cache.AutoScaleCache.Put(key, 1, time.Second * time.Duration(rand.Int31n(120)))
		go execAutoScale(v)
	}
}
