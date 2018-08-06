package app

import (
	"cloud/sql"
	"cloud/util"
	"cloud/models/app"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"golang.org/x/crypto/openpgp/errors"
	"database/sql/driver"
	"strings"
	"cloud/controllers/ent"
)

type ConfigureController struct {
	beego.Controller
}

// 配置文件管理入口页面
// @router /application/configure/list [get]
func (this *ConfigureController) ConfigureList() {
	this.TplName = "application/configure/list.html"
}

// 获取单个应用的数据
func getUpdateData(searchMap sql.SearchMap) app.CloudAppConfigure {
	update := app.CloudAppConfigure{}
	q := sql.SearchSql(app.CloudAppConfigure{}, app.SelectCloudAppConfigure, searchMap)
	sql.Raw(q).QueryRow(&update)
	return update
}

// 配置文件详情入口
// @router /application/configure/detail/:hi(.*) [get]
func (this *ConfigureController) DetailPage() {
	name := this.GetString(":hi")
	searchMap := sql.SearchMap{}
	searchMap.Put("ConfigureName", name)
	this.Data["data"] = getUpdateData(searchMap)
	this.TplName = "application/configure/detail.html"
}

// 配置文件管理添加页面
// @router /application/configure/add [get]
func (this *ConfigureController) ConfigureAdd() {
	var entHtml string
	entData := ent.GetEntnameSelect()
	update := app.CloudAppConfigure{}
	id, _ := this.GetInt("ConfigureId", 0)
	// 更新操作
	if id > 0 {
		searchMap := sql.GetSearchMap("ConfigureId", *this.Ctx)
		this.Data["readonly"] = "readonly"
		update = getUpdateData(searchMap)
		entHtml = util.GetSelectOptionName(update.Entname)
		this.Data["cluster"] = util.GetSelectOptionName(update.ClusterName)
	}
	this.Data["entname"] = entHtml + entData
	this.Data["data"] = update
	this.TplName = "application/configure/add.html"
}

// string
// 配置文件保存
// @router /api/configure [post]
func (this *ConfigureController) ConfigureSave() {
	d := app.CloudAppConfigure{}
	err := this.ParseForm(&d)
	if err != nil {
		this.Ctx.WriteString("参数错误" + err.Error())
		return
	}
	util.SetPublicData(d, getConfigUser(this), &d)

	var q = sql.InsertSql(d, app.InsertCloudAppConfigure)
	if d.ConfigureId > 0 {
		searchMap := sql.SearchMap{}
		searchMap.Put("ConfigureId", d.ConfigureId)

		q = sql.UpdateSql(d,
			app.UpdateCloudAppConfigure,
			searchMap,
			app.UpdateCloudAppConfigureExclude)
	}
	_, err = sql.Raw(q).Exec()
	data, msg := util.SaveResponse(err, "名称已经被使用")
	util.SaveOperLog(this.GetSession("username"),
		*this.Ctx,
		"保存配置文件配置 "+msg,
		d.ConfigureName)
	setConfigJson(this, data)
}

// 配置文件名称数据
// @router /api/configure/name [get]
func (this *ConfigureController) GetConfigureName() {
	data := make([]app.CloudAppConfigureName, 0)
	searchMap := sql.SearchMap{}
	keys := "ClusterName,Entname"
	searchMap = sql.GetSearchMapValue(strings.Split(keys, ","), *this.Ctx, searchMap)
	searchMap.Put("CreateUser", getConfigUser(this))
	searchSql := sql.SearchSql(app.CloudAppConfigure{},
		app.SelectCloudAppConfigure,
		searchMap)
	sql.Raw(searchSql).QueryRows(&data)
	setConfigJson(this, data)
}

// 配置文件数据
// @router /api/configure [get]
func (this *ConfigureController) ConfigureData() {
	data := make([]app.CloudAppConfigure, 0)
	searchMap := sql.SearchMap{}
	id := this.Ctx.Input.Param(":id")
	key := this.GetString("key")
	if id != "" {
		searchMap.Put("ConfigureId", id)
	}

	user := getConfigUser(this)
	searchMap.Put("CreateUser", user)

	searchSql := sql.SearchSql(app.CloudAppConfigure{}, app.SelectCloudAppConfigure, searchMap)
	if key != "" && id == "" {
		key := sql.Replace(key)
		q := strings.Replace(app.SelectCloudAppConfigSearch, "?", key, -1)
		searchSql += q
	}

	sql.OrderByPagingSql(searchSql,
		"configure_id",
		*this.Ctx.Request,
		&data,
		app.CloudAppConfigure{})

	r := util.ResponseMap(data,
		sql.CountSearchMap("cloud_app_configure",
			sql.GetSearchMapV("CreateUser", user),
			len(data), key),
		this.GetString("draw"))
	setConfigJson(this, r)
}

// json
// 删除配置文件
// @router /api/configure/:id:int [delete]
func (this *ConfigureController) ConfigureDelete() {
	searchMap := sql.SearchMap{}
	id := this.Ctx.Input.Param(":id")
	searchMap.Put("ConfigureId", id)
	configure := getUpdateData(searchMap)
	mountData := getMountData(configure.ConfigureName, configure.ClusterName, "")

	if len(mountData) > 0 {
		logs.Info("该项目被挂载不能删除", configure.ConfigureName)
		data := util.DeleteResponse(errors.InvalidArgumentError("已经被挂载,不能删除"),
			*this.Ctx, "删除配置文件数据"+configure.ConfigureName,
			this.GetSession("username"),
			configure.CreateUser,
			driver.ResultNoRows)
		this.Data["json"] = data
		return
	}

	q := sql.DeleteSql(app.DeleteCloudAppConfigure, searchMap)
	r, err := sql.Raw(q).Exec()
	data := util.DeleteResponse(err,
		*this.Ctx, "删除配置文件"+configure.ConfigureName,
		this.GetSession("username"),
		configure.CreateUser,
		r)
	setConfigJson(this, data)
}

// 2018-02-05 21:45
func setConfigJson(this *ConfigureController, data interface{}) {
	this.Data["json"] = data
	this.ServeJSON(false)
}

func getConfigUser(this *ConfigureController) string {
	return util.GetUser(this.GetSession("username"))
}