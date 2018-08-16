package app

import (
	"cloud/sql"
	"cloud/util"
	"cloud/models/app"
	"cloud/controllers/base/cluster"
	"strings"
	"encoding/json"
)

// 模板管理入口页面
// @router /application/template/list [get]
func (this *AppController) TemplateList() {
	this.TplName = "application/template/list.html"
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

	this.Data["data"] = update
	this.TplName = "application/template/update.html"
}

// 2018-08-16 09:16
// yaml 获取
// 将当期的服务信息存储,用来创建新的服务
func getServiceYaml(serviceName string) string {
	q := strings.Replace(app.SelectCloudAppService, "?", serviceName, -1)
	d := make([]app.CloudAppService, 0)
	sql.Raw(q).QueryRows(&d)
	yaml := make([]app.CloudAppService, 0)
	for _, v := range d {
		v.CreateUser = ""
		v.CreateTime = ""
		v.LastModifyUser = ""
		v.LastModifyTime = ""
		yaml = append(yaml, v)
	}
	content, _ := json.MarshalIndent(yaml, "", "  ")
	return strings.Replace(util.StringsToJSON(string(content)), "\\", "", -1)
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
	q := sql.UpdateSql(d, app.UpdateCloudAppTemplate, searchMap, "CreateTime,CreateUser,Cluster,ServiceName")
	sql.Exec(q)
	this.Data["json"] = util.ApiResponse(true, "更新完成")
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
		d.Yaml = getServiceYaml(d.ServiceName)
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
	template := app.CloudAppTemplate{}
	sql.Raw(sql.SearchSql(template, app.SelectCloudAppTemplate, searchMap)).QueryRow(&template)
	r, err := sql.Raw(sql.DeleteSql(app.DeleteCloudAppTemplate, searchMap)).Exec()
	data := util.DeleteResponse(err, *this.Ctx, "删除模板"+template.TemplateName, this.GetSession("username"), template.CreateUser, r)
	this.Data["json"] = data
	this.ServeJSON(false)
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
