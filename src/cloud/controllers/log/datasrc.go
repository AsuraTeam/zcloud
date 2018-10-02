package log

import (
	"github.com/astaxie/beego"
	"strings"
	"strconv"
	"cloud/sql"
	"cloud/util"
	"cloud/models/log"
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
			log.SelectAlarmDataSource,
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
			log.SelectAlarmDataSource,
			searchMap)
		sql.Raw(q).QueryRow(&update)
	}
	this.Data["data"] = update
	this.TplName = "log/datasrc/add.html"
}

// 获取数据源数据
func GetDataSourceSelect(searchMap sql.SearchMap) string {
	html := make([]string, 0)
	html = append(html, "<option>--请选择--</option>")
	data := getDataSourceData(searchMap)
	for _, v := range data {
		html = append(html, util.GetSelectOption(v.Name, strconv.FormatInt(v.DataSourceId, 10), v.Name))
	}
	return strings.Join(html, "\n")
}

func getDataSourceData(searchMap sql.SearchMap) []log.LogDataSource {
	// 数据源数据
	data := make([]log.LogDataSource, 0)
	q := sql.SearchSql(
		log.LogDataSource{},
		log.SelectAlarmDataSource,
		searchMap)
	sql.Raw(q).QueryRows(&data)
	return data
}

// 获取数据源名称对应关系
func GetDataSourceMap() util.Lock {
	data := getDataSourceData(sql.SearchMap{})
	r := util.Lock{}
	for _, v := range data {
		r.Put(strconv.FormatInt(v.DataSourceId, 10), v.Name)
	}
	return r
}

// 获取数据源数据
// router /api/log/datasrc [get]
func (this *DataSourceController) DataSourceData() {
	data := getDataSourceData(sql.SearchMap{})
	setDataSourceJson(this, data)
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

	q := sql.InsertSql(d, log.InsertAlarmDataSource)
	if d.DataSourceId > 0 {
		searchMap := sql.SearchMap{}
		searchMap.Put("DataSourceId", d.DataSourceId)
		q = sql.UpdateSql(
			d, log.UpdateAlarmDataSource, searchMap, "CreateTime,CreateUser")
	}
	_, err = sql.Exec(q)

	data, msg := util.SaveResponse(err, "名称已经被使用")
	util.SaveOperLog(this.GetSession("username"), *this.Ctx, "保存数据源配置 "+msg, "")
	setDataSourceJson(this, data)
}

// 数据源数据
// @router /api/log/datasrc [get]
func (this *DataSourceController) DataSourceDatas() {
	data := make([]log.LogDataSource, 0)
	searchMap := sql.SearchMap{}
	id := this.Ctx.Input.Param(":id")
	key := this.GetString("search")
	if id != "" {
		searchMap.Put("DataSourceId", id)
	}
	searchSql := sql.SearchSql(
		log.LogDataSource{},
		log.SelectAlarmDataSource,
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

	q := sql.SearchSql(data, log.SelectAlarmDataSource, searchMap)
	sql.Raw(q).QueryRow(&data)

	q = sql.DeleteSql(log.DeleteAlarmDataSource, searchMap)
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
