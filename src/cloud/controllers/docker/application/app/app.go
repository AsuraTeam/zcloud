package app

import (
	"github.com/astaxie/beego"
	"cloud/k8s"
	"cloud/sql"
	"cloud/models/app"
	"cloud/util"
	"database/sql/driver"
	"github.com/astaxie/beego/logs"
	"strings"
	"cloud/models/ci"
	"cloud/models/registry"
	"cloud/controllers/ent"
	"cloud/controllers/base/quota"
	"strconv"
)

type AppController struct {
	beego.Controller
}

var (
	LockContainerUpdate util.Lock
)

// 容器应用入口页面
// @router /application/app/index [get]
func (this *AppController) AppList() {
	this.TplName = "application/app/list.html"
}

// 响应错误数据
// 2018-01-16 21:13
func responseAppData(err error, this *AppController, appName string, info string) {
	data, msg := util.SaveResponse(err, info)
	util.SaveOperLog(this.GetSession("username"), *this.Ctx, info+": "+msg, appName)
	this.Data["json"] = data
	this.ServeJSON(false)
}

// 2018-01-16 20:33
// 应用扩缩容接口
// @router /api/app/scale/:id:int
func (this *AppController) AppScale() {
	searchMap := sql.GetSearchMap("AppId", *this.Ctx)
	d := app.CloudApp{}
	q := sql.SearchSql(d, app.SelectCloudApp, searchMap)
	sql.Raw(q).QueryRow(&d)
	start := this.GetString("start")
	replicas, _ := this.GetInt("replicas")

	serviceData := make([]app.CloudAppService, 0)
	serverMap := sql.GetSearchMapV("AppName", d.AppName)

	q = sql.SearchSql(app.CloudAppService{}, app.SelectCloudAppService, serverMap)
	sql.Raw(q).QueryRows(&serviceData)

	// 启动服务
	if start == "1" {
		for _, v := range serviceData {
			err := k8s.ScalePod(d.ClusterName, util.Namespace(d.AppName, d.ResourceName), v.ServiceName, int32(v.Replicas))
			if err != nil {
				logs.Error("启动服务", err)
			} else {
				logs.Info("启动服务", d.AppName, d.ResourceName, v.ServiceName)
			}
		}
		responseAppData(nil, this, d.AppName, "操作成功")
		return
	}

	// 停止所有服务
	if replicas == 0 {
		for _, v := range serviceData {
			err := k8s.ScalePod(d.ClusterName, util.Namespace(d.AppName, d.ResourceName), v.ServiceName, 0)
			if err != nil {
				logs.Error("停止服务", err)
			} else {
				logs.Info("停止服务", d.AppName, d.ResourceName, v.ServiceName)
			}
		}
		responseAppData(nil, this, d.AppName, "操作成功")
		return
	}
	responseAppData(nil, this, d.AppName, "操作成操作失败")
}

// 2018-02-03 19:03
// 应用 名称数据
// @param ClusterName
// @router /api/app/name [get]
func (this *AppController) GetAppName() {
	data := make([]app.CloudAppName, 0)
	searchMap := sql.SearchMap{}
	q := strings.Split("Entname,ClusterName", ",")
	searchMap = sql.GetSearchMapValue(q, *this.Ctx, searchMap)
	searchMap.Put("CreateUser", getUser(this))
	searchSql := sql.SearchSql(app.CloudAppService{}, app.SelectCloudApp, searchMap)
	sql.Raw(searchSql).QueryRows(&data)
	SetAppDataJson(this, data)
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

// 容器列表入口
// 2018-01-15 14:57
// @router /application/container/list
func (this *AppController) ContainerList() {
	this.Data["Entname"] = ent.GetEntnameSelect()
	this.Data["AppData"] = GetAppSelect(sql.GetSearchMapV("CreateUser", getUser(this)))
	this.TplName = "application/container/list.html"
}

// 获取应用名称信息
func getAppDataQ(searchMap sql.SearchMap) []app.CloudAppName {
	data := make([]app.CloudAppName, 0)
	searchSql := sql.SearchSql(app.CloudAppName{}, app.GetAppName, searchMap)
	logs.Info("searchSql", searchSql, searchMap)
	sql.Raw(searchSql).QueryRows(&data)
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
func selectAppData(searchMap sql.SearchMap)[]app.CloudApp  {
	data := make([]app.CloudApp, 0)
	searchSql := sql.SearchSql(app.CloudAppService{}, app.SelectCloudApp, searchMap)
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

// 应用详情页面
// @router /application/app/detail/:id:int [get]
func (this *AppController) AppDetail() {
	id := this.Ctx.Input.Param(":id")
	data := app.CloudApp{}
	searchMap := sql.SearchMap{}
	searchMap.Put("AppId", id)
	searchMap.Put("CreateUser", getUser(this))
	datas := selectAppData(searchMap)

	if len(datas) > 0 {
		data = datas[0]
	}
	if data.AppId == 0 {
		this.TplName = "application/app/list.html"
		return
	}
	yamlShow := this.GetString("yaml")
	this.Data["detault"] = "active"
	this.Data["yamlActive"] = ""
	if yamlShow == "1" {
		this.Data["yamlActive"] = "active"
		this.Data["detault"] = ""
	}
	this.Data["namespace"] = util.Namespace(data.AppName, data.ResourceName)
	this.Data["data"] = data
	yaml := util.Json2Yaml(data.Yaml)
	this.Data["Yaml"] = yaml
	this.TplName = "application/app/detail.html"
}

// 添加应用页面
// @router /application/app/add [get]
func (this *AppController) AppAdd() {
	clusterName := this.GetString("ClusterName")
	data := app.CloudAppService{}
	data.ClusterName = clusterName
	imageId := this.GetString("imageId")
	historyId := this.GetString("historyId")
	entData := ent.GetEntnameSelect()
	if imageId != "" {
		image := k8s.CloudImage{}
		q := sql.SearchSql(image, registry.SelectImageTgs, sql.GetSearchMapV("ImageId", imageId))
		sql.Raw(q).QueryRow(&image)
		if image.ImageId > 0 {
			data.ImageRegistry = image.Access + "/" + image.Name
			data.ImageTag = makeImageTags(image.Tags)
		}
	}
	if historyId != "" {
		history := ci.CloudBuildJobHistory{}
		q := sql.SearchSql(history, ci.SelectBuildJobToApp, sql.GetSearchMapV("HistoryId", historyId))
		sql.Raw(q).QueryRow(&history)
		if history.HistoryId > 0 {
			data.ServiceName = history.ItemName + "-service"
			data.ImageRegistry = history.RegistryServer + "/" + history.RegistryGroup + "/" + history.ItemName
			data.ImageTag = util.GetSelectOption(history.ImageTag, history.ImageTag, history.ImageTag)
		}
	}

	var ent string
	d, _ := getApp(this)
	if d.AppId > 0 {
		ent = util.GetSelectOptionName(d.Entname)
		this.Data["cluster"] = util.GetSelectOptionName(d.ClusterName)
		data.AppName = d.AppName
	}
	this.Data["entname"] = ent + entData
	this.Data["data"] = data

	quotas := quota.GetUserQuota(getUser(this), "app")
	this.Data["quotas"] = quotas
	logs.Info("quotas", quotas)
	this.TplName = "application/app/add.html"
}

// 2018-02-26 09:24
// 重新部署应用
// @router /api/app/redeploy [post]
func (this *AppController) RedeployApp() {
	ids := this.GetString("apps")
	user := getUser(this)
	for _, v:= range strings.Split(ids, ","){

		if _, err := strconv.Atoi(v); err != nil {
			continue
		}

		services, status := getRedeployService(v, user)
		if status {
			for _, service := range services {
				ExecDeploy(service, true)
			}
		}
	}
	SetAppDataJson(this, util.ApiResponse(true, "成功,重建中..."))
}

// 2018-02-26 09:32
// 获取重建的应用信息
func getRedeployApp(v string, user string) ([]app.CloudAppName,bool) {
	searchMap := sql.SearchMap{}
	searchMap.Put("AppId", v)
	searchMap.Put("CreateUser", user)
	r := getAppDataQ(searchMap)
	if len(r) == 0 {
		return []app.CloudAppName{},false
	}
	return r, true
}

// 2018-02-26 09:54
// 获取重建应用的服务信息
func getRedeployService(v string, user string) ([]app.CloudAppService, bool) {
	logs.Info("开始重建服务", util.ObjToString(v))
	data,status := getRedeployApp(v, user)
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

// 删除应用
// @router /api/app/:id:int [delete]
func (this *AppController) AppDelete() {
	force := this.GetString("force")
	d, searchMap := getApp(this)

	// 先去服务器删除,成功后再删除数据库
	namespace := util.Namespace(d.AppName, d.ResourceName)

	client, err := k8s.GetClient(d.ClusterName)
	k8s.DeleteSecret(client, namespace)
	err = k8s.DeletelDeployment(namespace, true, "", d.ClusterName)

	if err != nil && force == "" {
		data := util.DeleteResponse(err, *this.Ctx,
			"删除应用"+d.AppName,
			this.GetSession("username"),
			d.CreateUser,
			driver.ResultNoRows)

		this.Data["json"] = data
		this.ServeJSON(false)
		return
	}

	q := sql.DeleteSql(app.DeleteCloudApp, searchMap)
	r, err := sql.Raw(q).Exec()

	searchMap = sql.GetSearchMapV(
		"AppName", d.AppName,
		"Entname", d.Entname,
		"ClusterName", d.ClusterName)
	q = sql.DeleteSql(app.DeleteCloudAppService, searchMap)
	sql.Raw(q).Exec()

	data := util.DeleteResponse(err,
		*this.Ctx,
		"删除应用"+d.AppName,
		getUser(this),
		d.CreateUser, r)

	this.Data["json"] = data
	this.ServeJSON(false)
}

// @router /api/app [get]
func (this *AppController) AppData() {
	data := make([]app.CloudApp, 0)
	searchMap := sql.SearchMap{}
	ip := this.GetString("ip")
	searchMap = sql.GetSearchMapValue(
		sql.MKeyV("AppName"),
		*this.Ctx, searchMap)

	searchMap.Put("CreateUser", getUser(this))
	searchSql := sql.SearchSql(app.CloudApp{}, app.SelectCloudApp, searchMap)
	if ip != "" {
		q := ` and (app_name like "%?%")`
		searchSql += strings.Replace(q, "?", sql.Replace(ip), -1)
	}

	sql.OrderByPagingSql(searchSql, "app_id",
		*this.Ctx.Request, &data,
		app.CloudApp{})

	cloudApps := getCacheAppData(data)
	r := util.ResponseMap(cloudApps, len(data), this.GetString("draw"))

	this.Data["json"] = r
	this.ServeJSON(false)
	go getK8sAppData(data)
	if len(data) > 0 {
		go MakeContainerData(util.Namespace(data[0].AppName, data[0].ResourceName))
	}
}

// 服务管理入口页面
// @router /application/app/service/service [get]
func (this *AppController) ServiceList() {
	this.TplName = "application/app/service/service.html"
}

// 模板管理入口页面
// @router /application/app/service/envfile [get]
func (this *AppController) EnvfileList() {
	//this.TplName = "application/app/service/envfile.html"
	this.Ctx.WriteString("建设中")

}

// 模板管理入口页面
// @router /application/app/service/configure [get]
func (this *AppController) ConfigureList() {
	//this.TplName = "application/app/service/configure.html"
	this.Ctx.WriteString("建设中")
}
