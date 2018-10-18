package log

import (
	"github.com/astaxie/beego"
	"cloud/sql"
	"cloud/util"
	"cloud/models/log"
	"cloud/controllers/docker/application/app"
)

type DataSourceController struct {
	beego.Controller
}

// 数据源管理入口页面
// @router /log/datasrc/list [get]
func (this *DataSourceController) DataSourceList() {
	this.TplName = "log/datasrc/list.html"
}

// 日志驱动管理入口页面
// @router /log/driver/list [get]
func (this *DataSourceController) DriverList() {
	this.TplName = "log/datasrc/list-driver.html"
}

// 数据源管理添加页面
// @router /log/driver/add [get]
func (this *DataSourceController) DriverAdd() {
	id := this.GetString("DataSourceId")
	update := log.LogDataSource{}

	// 更新操作
	if id != "" {
		searchMap := sql.GetSearchMap("DataSourceId", *this.Ctx)
		q := sql.SearchSql(
			log.LogDataSource{},
			log.SelectLogDataSource,
			searchMap)
		sql.Raw(q).QueryRow(&update)
	}
	this.Data["data"] = update
	this.TplName = "log/datasrc/add-driver.html"
}

// 数据源管理添加页面
// @router /log/datasrc/add [get]
func (this *DataSourceController) DataSourceAdd() {
	id := this.GetString("DataSourceId")
	update := log.LogDataSource{}

	// 更新操作
	if id != "" {
		searchMap := sql.GetSearchMap("DataSourceId", *this.Ctx)
		q := sql.SearchSql(
			log.LogDataSource{},
			log.SelectLogDataSource,
			searchMap)
		sql.Raw(q).QueryRow(&update)
	}
	this.Data["data"] = update
	this.TplName = "log/datasrc/add.html"
}


// 数据源保存
// @router /api/log/datasrc [post]
func (this *DataSourceController) DataSourceSave() {
	d := log.LogDataSource{}
	err := this.ParseForm(&d)
	if err != nil {
		this.Ctx.WriteString("参数错误" + err.Error())
		return
	}
	util.SetPublicData(d, util.GetUser(this.GetSession("username")), &d)

	q := sql.InsertSql(d, log.InsertLogDataSource)
	if d.DataSourceId > 0 {
		searchMap := sql.SearchMap{}
		searchMap.Put("DataSourceId", d.DataSourceId)
		q = sql.UpdateSql(
			d, log.UpdateLogDataSource, searchMap, "CreateTime,CreateUser")
	}
	_, err = sql.Exec(q)

	data, msg := util.SaveResponse(err, "名称已经被使用")
	util.SaveOperLog(this.GetSession("username"), *this.Ctx, "保存数据源配置 "+msg, "")
	if d.DataType == "driver" {
		go app.MakeFilebeatConfig(d.Ent, d.ClusterName)
	}
	setDataSourceJson(this, data)
}

// 数据源数据
// @router /api/log/datasrc [get]
func (this *DataSourceController) DataSourceDatas() {
	data := make([]log.LogDataSource, 0)
	searchMap := sql.SearchMap{}
	tp := this.GetString("type")
	key := this.GetString("search")
	searchMap.Put("DataType", tp)
	searchSql := sql.SearchSql(
		log.LogDataSource{},
		log.SelectLogDataSource,
		searchMap)

	num, _ := sql.OrderByPagingSql(
		searchSql, "data_source_id",
		*this.Ctx.Request,
		&data,
		log.LogDataSource{})

	r := util.ResponseMap(
		data,
		sql.Count("log_data_source", int(num), key),
		this.GetString("draw"))
	setDataSourceJson(this, r)

}

// 删除数据源
// @router /api/log/datasrc/:id:int [delete]
func (this *DataSourceController) DataSourceDelete() {
	searchMap := sql.GetSearchMap("DataSourceId", *this.Ctx)
	data := log.LogDataSource{}

	q := sql.SearchSql(data, log.SelectLogDataSource, searchMap)
	sql.Raw(q).QueryRow(&data)

	q = sql.DeleteSql(log.DeleteLogDataSource, searchMap)
	r, err := sql.Exec(q)

	datar := util.DeleteResponse(
		err,
		*this.Ctx,
		"删除数据源"+data.Name,
		this.GetSession("username"),
		data.Name, r)
	setDataSourceJson(this, datar)
}

func setDataSourceJson(this *DataSourceController, data interface{}) {
	this.Data["json"] = data
	this.ServeJSON(false)
}
