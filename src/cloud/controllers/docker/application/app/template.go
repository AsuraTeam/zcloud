package app

import (
	"cloud/sql"
	"cloud/util"
	"cloud/models/app"
	"cloud/controllers/base/cluster"
	"strings"
	"encoding/json"
	"cloud/controllers/ent"
	"cloud/controllers/base/quota"
	"github.com/astaxie/beego/logs"
	"fmt"
	"cloud/k8s"
	"strconv"
	"time"
)

// 模板管理入口页面
// @router /application/template/list [get]
func (this *AppController) TemplateList() {
	this.TplName = "application/template/list.html"
}

// 编排历史页面
// @router /application/template/deploy/history [get]
func (this *AppController) HistoryList() {
	this.TplName = "application/template/history.html"
}

// 模板管理添加页面
// @router /application/template/add [get]
func (this *AppController) TemplateAdd() {
	id := this.GetString("TemplateId")
	update := app.CloudAppTemplate{}
	this.Data["cluster"] = cluster.GetClusterSelect()
	// 更新操作
	if id != "0" {
		searchMap := sql.GetSearchMap("TemplateId", *this.Ctx)
		sql.Raw(sql.SearchSql(app.CloudAppTemplate{}, app.SelectCloudAppTemplate, searchMap)).QueryRow(&update)
		this.Data["data"] = update
		this.Data["cluster"] = util.GetSelectOptionName(update.Cluster)
	}
	this.TplName = "application/template/add.html"
}

// 2018-08-16 09:45
// 模板yaml更新添加页面
// @router /application/template/update/add [get]
func (this *AppController) TemplateUpdateAdd() {
	update := app.CloudAppTemplate{}
	searchMap := sql.GetSearchMap("TemplateId", *this.Ctx)
	sql.Raw(sql.SearchSql(app.CloudAppTemplate{}, app.SelectCloudAppTemplate, searchMap)).QueryRow(&update)
	update.Yaml = util.Base64Decoding(update.Yaml)
	this.Data["data"] = update
	this.TplName = "application/template/update.html"
}

// 2018-08-16 11:06
// 模板管理应用拉起添加页面
// @router /application/template/deploy/add [get]
func (this *AppController) TemplateDeployAdd() {
	update := app.CloudAppTemplate{}
	entData := ent.GetEntnameSelect()
	this.Data["cluster"] = cluster.GetClusterSelect()
	searchMap := sql.GetSearchMap("TemplateId", *this.Ctx)
	sql.Raw(sql.SearchSql(app.CloudAppTemplate{}, app.SelectCloudAppTemplate, searchMap)).QueryRow(&update)
	this.Data["data"] = update
	this.Data["entname"] = entData
	quotas := quota.GetUserQuota(getUser(this), "app")
	this.Data["quotas"] = quotas
	this.TplName = "application/template/deploy.html"
}

// 2018-08-16 09:16
// yaml 获取
// 将当期的服务信息存储,用来创建新的服务
func getServiceYaml(serviceName string, clusterName string) string {
	names := strings.Split(serviceName, ",")
	yaml := make([]app.CloudAppService, 0)
	for _, v := range names {
		searchMap := sql.SearchMap{}
		vs := strings.Split(v, ";")
		if len(vs) < 2 {
			logs.Error("参数异常", util.ObjToString(vs))
			continue
		}
		searchMap.Put("AppName", vs[0])
		searchMap.Put("ResourceName", vs[1])
		searchMap.Put("ServiceName", vs[2])
		searchMap.Put("ClusterName", clusterName)
		searchMap.Put("ServiceVersion", 1)
		q := sql.SearchSql(app.CloudAppService{}, app.SelectCloudAppService, searchMap)
		d := make([]app.CloudAppService, 0)
		sql.Raw(q).QueryRows(&d)

		for _, v := range d {
			v.CreateUser = ""
			v.CreateTime = ""
			v.LastModifyUser = ""
			v.LastModifyTime = ""
			v.ServiceId = 0
			yaml = append(yaml, v)
		}

	}
	content, _ := json.MarshalIndent(yaml, "", "  ")
	return util.Base64Encoding(util.StringsToJSON(string(content)))
}

//  2018-08-16 10:31
// 拉起应用模板环境
func startDeploy(yaml string, appName string, ent string, clusterName string, resourceName string, user string, templateName string, domain string, envs string) {
	services := make([]app.CloudAppService, 0)
	err := json.Unmarshal([]byte(util.Base64Decoding(yaml)), &services)
	if err != nil {
		logs.Error("部署服务转换失败", err)
		return
	}
	logs.Info("正在拉起环境", appName, ent, clusterName, resourceName)
	appData := app.CloudApp{}
	appData.AppName = appName
	appData.ClusterName = clusterName
	appData.Entname = ent
	appData.ResourceName = resourceName
	appData.CreateTime = util.GetDate()
	appData.CreateUser = user
	_, err = sql.Exec(sql.InsertSql(appData, app.InsertCloudApp))
	if err != nil {
		logs.Error("添加应用失败 ", err.Error())
		return
	}

	for _, service := range services {
		service.AppName = appName
		service.Entname = ent
		service.ClusterName = clusterName
		service.ResourceName = resourceName
		service.CreateUser = user
		if len(envs) > 0 {
			service.Envs = envs
		}
		service.CreateTime = util.GetDate()
		history := app.CloudTemplateDeployHistory{AppName: appName,
			Entname: ent,
			ClusterName: clusterName,
			ResourceName: resourceName,
			CreateUser: user,
			CreateTime: util.GetDate(),
			TemplateName: templateName,
			ServiceName: service.ServiceName,
			Domain: domain,
		}
		if len(domain) > 0 {
			history.Domain = fmt.Sprintf("%s.%s.%s", appName, service.ServiceName, domain)
		}
		sql.Exec(sql.InsertSql(history, app.InsertCloudTemplateDeployHistory))
		ExecDeploy(service, false)

	}

	if len(domain) > 0 {
		logs.Info("获取到拉起服务的域名后缀", domain)
		for _, service := range services {
			createLbConfig(service, clusterName, ent, appName, domain, user, resourceName)
		}
		go k8s.CreateNginxConf("")
	}
}

// 2018-08-29 07:51
// 创建负载均衡配置
func createLbConfig(service app.CloudAppService, clusterName string, ent string, appName string, domain string, user string, resourceName string)  {
	if len(service.ContainerPort) < 1 {
		logs.Error("获取服务端口错误", service.ServiceName)
		return
	}

	searchMap := sql.SearchMap{}
	searchMap.Put("ClusterName", clusterName)
	searchMap.Put("Entname", ent)
	var lbData k8s.CloudLb
	lb := k8s.GetLbDataSearchMap(searchMap)
	if lb!= nil{
		lbData = lb.(k8s.CloudLb)
	}
	if lbData.LbId == 0 {
		logs.Error("服务获取负载均衡失败，该集群环境没有配置负载均衡", ent, clusterName)
		return
	}
	searchMap.Put("AppName", appName)
	searchMap.Put("ResourceName", resourceName)
	searchMap.Put("ServiceName", service.ServiceName)
	services := getServiceData(searchMap, "")
	if len(services) > 0 {
		service.ServiceId = services[0].ServiceId
	}
	conf := k8s.CloudLbService{
		ServiceName:    service.ServiceName,
		AppName:        appName,
		Domain:         fmt.Sprintf("%s.%s", service.ServiceName, domain),
		LbType:         "Nginx",
		ClusterName:    clusterName,
		ResourceName:   resourceName,
		Protocol:       "HTTP",
		Entname:        ent,
		ServiceVersion: "1",
		LbMethod:       "service",
		LbId:           lbData.LbId,
		LbName:         lbData.LbName,
		LbServiceId:    strconv.FormatInt(service.ServiceId, 10),
		CreateTime:     util.GetDate(),
		CreateUser:     user,
		LastModifyTime: util.GetDate(),
		LastModifyUser: user,
		ContainerPort:  strings.Split(service.ContainerPort, ".")[0],
	}
	sql.Exec(sql.InsertSql(conf, "insert into cloud_lb_service" ))

	q := fmt.Sprintf(app.UpdateServiceDomain, domain, appName, clusterName, service.ServiceName, resourceName)
	sql.Exec(q)
	time.Sleep(time.Second * 6)
	UpdateServiceDomain()
}

// 2018-08-21 11:16
// 更新服务域名
func UpdateServiceDomain()  {
	logs.Info("更新服务域名")
	data := make([]k8s.CloudLbService, 0)
	q := sql.SearchSql(k8s.CloudLbService{}, k8s.SelectCloudLbService, sql.SearchMap{})
	sql.Raw(q).QueryRows(&data)
	for _, v := range data {
		service := app.CloudAppService{}
		searchMap := sql.SearchMap{}
		searchMap.Put("ClusterName", v.ClusterName)
		searchMap.Put("AppName", v.AppName)
		searchMap.Put("ServiceName", v.ServiceName)
		searchMap.Put("ResourceName", v.ResourceName)
		if len(v.ServiceVersion) == 0 {
			v.ServiceVersion = "1"
		}
		searchMap.Put("ServiceVersion", v.ServiceVersion)
		q = sql.SearchSql(app.CloudAppService{}, app.SelectCloudAppService, searchMap)
		sql.Raw(q).QueryRow(&service)
		service.Domain = v.Domain
		if service.ServiceId != 0 {
			sql.Exec(sql.UpdateSql(service, app.UpdateCloudAppService, searchMap, ""))
		}
	}
}

// 2018-08-16 10:35
// 执行环境拉起操作
// @router /api/template/deploy/:id:int [post]
func (this *AppController) StartDeploy() {
	d := app.CloudAppTemplate{}
	envs := this.GetString("Envs")
	this.ParseForm(&d)
	searchMap := sql.SearchMap{}
	id := this.Ctx.Input.Param(":id")
	searchMap.Put("TemplateId", id)
	template := getTemplateData(searchMap)
	go startDeploy(template.Yaml, d.AppName, d.Ent, d.Cluster, d.ResourceName, util.GetUser(this.GetSession("username")), d.TemplateName, d.Domain, envs)
	this.Data["json"] = util.ApiResponse(true, "保存成功,正在拉起环境,请耐心等待")
	this.ServeJSON(false)
}

// 2018-08-16 09:39
// 模板更新yaml数据
// @router /api/template/update [post]
func (this *AppController) TemplateUpdate() {
	d := app.CloudAppTemplate{}
	this.ParseForm(&d)
	util.SetPublicData(d, util.GetUser(this.GetSession("username")), &d)
	searchMap := sql.SearchMap{}
	searchMap.Put("TemplateId", d.TemplateId)
	d.Yaml = util.Base64Encoding(d.Yaml)
	q := sql.UpdateSql(d, app.UpdateCloudAppTemplate, searchMap, "CreateTime,CreateUser,Cluster,ServiceName")
	sql.Exec(q)
	this.Data["json"] = util.ApiResponse(true, "保存成功")
	this.ServeJSON(false)
}

// string
// 模板保存
// @router /api/template [post]
func (this *AppController) TemplateSave() {
	d := app.CloudAppTemplate{}
	err := this.ParseForm(&d)
	if err != nil {
		this.Ctx.WriteString("参数错误" + err.Error())
		return
	}
	util.SetPublicData(d, util.GetUser(this.GetSession("username")), &d)
	var q string
	if d.TemplateId > 0 {
		searchMap := sql.SearchMap{}
		searchMap.Put("TemplateId", d.TemplateId)
		q = sql.UpdateSql(d, app.UpdateCloudAppTemplate, searchMap, "CreateTime,CreateUser,Cluster,Yaml")
	} else {
		d.Yaml = getServiceYaml(d.ServiceName, d.Cluster)
		q = sql.InsertSql(d, app.InsertCloudAppTemplate)
	}
	_, err = sql.Exec(q)
	data, msg := util.SaveResponse(err, "名称已经被使用")
	util.SaveOperLog(this.GetSession("username"), *this.Ctx, "保存模板配置 "+msg, d.TemplateName)
	this.Data["json"] = data
	this.ServeJSON(false)
}

// 模板名称数据
// @router /api/template/name [get]
func (this *AppController) GetTemplateName() {
	data := make([]app.CloudAppTemplateName, 0)
	searchSql := sql.SearchSql(app.CloudAppTemplate{}, app.SelectCloudAppTemplate, sql.SearchMap{})
	sql.Raw(searchSql).QueryRows(&data)
	this.Data["json"] = data
	this.ServeJSON(false)
}

// 模板数据
// @router /api/template/deploy/history [get]
func (this *AppController) HistoryData() {
	data := make([]app.CloudTemplateDeployHistory, 0)
	key := this.GetString("key")
	searchSql := sql.SearchSql(app.CloudAppTemplate{}, app.SelectCloudTemplateDeployHistory, sql.SearchMap{})
	if key != "" {
		searchSql += " where 1=1 and (template_name like \"%" + sql.Replace(key) + "%\" or service_name like \"%" + sql.Replace(key) + "%\")"
	}
	sql.OrderByPagingSql(searchSql, "history_id",
		*this.Ctx.Request,
		&data,
		app.CloudTemplateDeployHistory{})
	r := util.ResponseMap(data,
		sql.CountSearchMap("cloud_template_deploy_history",
			sql.GetSearchMapV("CreateUser", getUser(this)),
			len(data),
			""),
		this.GetString("draw"))
	this.Data["json"] = r
	this.ServeJSON(false)
}

// 模板数据
// @router /api/template [get]
func (this *AppController) TemplateData() {
	data := make([]app.CloudAppTemplate, 0)
	searchMap := sql.SearchMap{}
	id := this.Ctx.Input.Param(":id")
	key := this.GetString("key")
	if id != "" {
		searchMap.Put("TemplateId", id)
	}
	searchSql := sql.SearchSql(app.CloudAppTemplate{}, app.SelectCloudAppTemplate, searchMap)
	if key != "" && id == "" {
		searchSql += " where 1=1 and (template_name like \"%" + sql.Replace(key) + "%\" or yaml like \"%" + sql.Replace(key) + "%\")"
	}
	num, err := sql.Raw(searchSql).QueryRows(&data)
	var r = util.ResponseMap(data, num, this.GetString("draw"))
	if err != nil {
		r = util.ResponseMapError(err.Error())
	}
	this.Data["json"] = r
	this.ServeJSON(false)
}

// json
// 删除模板
// @router /api/template/:id:int [delete]
func (this *AppController) TemplateDelete() {
	searchMap := sql.SearchMap{}
	id := this.Ctx.Input.Param(":id")
	searchMap.Put("TemplateId", id)
	template := getTemplateData(searchMap)
	r, err := sql.Raw(sql.DeleteSql(app.DeleteCloudAppTemplate, searchMap)).Exec()
	data := util.DeleteResponse(err, *this.Ctx, "删除模板"+template.TemplateName, this.GetSession("username"), template.CreateUser, r)
	this.Data["json"] = data
	this.ServeJSON(false)
}

func getTemplateData(searchMap sql.SearchMap) app.CloudAppTemplate {
	template := app.CloudAppTemplate{}
	sql.Raw(sql.SearchSql(template, app.SelectCloudAppTemplate, searchMap)).QueryRow(&template)
	return template
}

// string
// 检查yaml是否可以转换成json格式
// @router /api/template/yaml/check [post]
func (this *AppController) YamlCheck() {
	yaml := this.GetString("yaml")
	_, err := util.Yaml2Json([]byte(yaml))
	if err != nil {
		this.Ctx.WriteString("false")
		return
	} else {
		this.Ctx.WriteString("true")
	}
}
