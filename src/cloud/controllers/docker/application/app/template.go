package app

import (
	"cloud/sql"
	"cloud/util"
	"cloud/models/app"
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
	// 更新操作
	if id != "" {
		searchMap := sql.GetSearchMap("TemplateId", *this.Ctx)
		update := app.CloudAppTemplate{}
		sql.Raw(sql.SearchSql(app.CloudAppTemplate{}, app.SelectCloudAppTemplate, searchMap)).QueryRow(&update)
		this.Data["data"] = update
	}
	this.TplName = "application/template/add.html"
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
	if d.TemplateId > 0 {
		searchMap := sql.SearchMap{}
		searchMap.Put("TemplateId", d.TemplateId)
		_, err = sql.Raw(sql.UpdateSql(d, app.UpdateCloudAppTemplate, searchMap, "CreateTime,CreateUser")).Exec()
	} else {
		_, err = sql.Raw(sql.InsertSql(d, app.InsertCloudAppTemplate)).Exec()
	}
	data,msg := util.SaveResponse(err, "名称已经被使用")
	util.SaveOperLog(this.GetSession("username"), *this.Ctx, "保存模板配置 "+msg, d.TemplateName)
	this.Data["json"] = data
	this.ServeJSON(false)
}

// 模板名称数据
// @router /api/template/name [get]
func (this *AppController) GetTemplateName() {
	data := []app.CloudAppTemplateName{}
	searchSql := sql.SearchSql(app.CloudAppTemplate{}, app.SelectCloudAppTemplate, sql.SearchMap{})
	sql.Raw(searchSql).QueryRows(&data)
	this.Data["json"] = data
	this.ServeJSON(false)
}

// 模板数据
// @router /api/template [get]
func (this *AppController) TemplateData() {
	data := []app.CloudAppTemplate{}
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
func (this *AppController) YamlCheck()  {
	yaml := this.GetString("yaml")
	_, err := util.Yaml2Json([]byte(yaml))
	if err != nil {
		this.Ctx.WriteString("false")
		return
	}else{
		this.Ctx.WriteString("true")
	}
}


