package ent

import (
	"cloud/sql"
	"cloud/util"
	"github.com/astaxie/beego"
	"cloud/models/ent"
	"cloud/controllers/base/cluster"
	"strings"
)

type EntController struct {
	beego.Controller
}

// 环境管理入口页面
// @router /ent/list [get]
func (this *EntController) EntList() {
	this.TplName = "ent/list.html"
}

// 2018-02-05 21:02
// 环境管理添加页面
// @router /ent/add [get]
func (this *EntController) EntAdd() {
	id := this.GetString("EntId")
	update := ent.CloudEnt{}
	clusterData := cluster.GetClusterSelect()
	var clusterHtml string
	// 更新操作
	if id != "0" {
		searchMap := sql.GetSearchMap("EntId", *this.Ctx)
		q := sql.SearchSql(ent.CloudEnt{}, ent.SelectCloudEnt, searchMap)
		sql.Raw(q).QueryRow(&update)
		clusters := update.Clusters
		for _, v := range strings.Split(clusters, ",") {
			clusterHtml += util.GetSelectOption(v, v, v)
		}
	}
	this.Data["selectCluster"] = clusterHtml
	this.Data["clusters"] = strings.Replace(clusterData, "<option></option>", "", -1)
	this.Data["data"] = update
	this.TplName = "ent/add.html"
}

// 获取环境数据
// 2018-01-20 12:56
// router /api/ent [get]
func (this *EntController) EntData() {
	// 环境数据
	data := make([]ent.CloudEnt, 0)
	q := sql.SearchSql(ent.CloudEnt{}, ent.SelectCloudEnt, sql.SearchMap{})
	sql.Raw(q).QueryRows(&data)
	setEntJson(this, data)
}

// string
// 环境保存
// @router /api/ent [post]
func (this *EntController) EntSave() {
	d := ent.CloudEnt{}
	err := this.ParseForm(&d)
	if err != nil {
		this.Ctx.WriteString("参数错误" + err.Error())
		return
	}
	clustersSet := make([]string, 0)
	clusters := d.Clusters
	for _, v := range strings.Split(clusters, ",") {
		if !util.ListExistsString(clustersSet, v) {
			clustersSet = append(clustersSet, v)
		}
	}
	d.Clusters = strings.Join(clustersSet, ",")
	util.SetPublicData(d, util.GetUser(this.GetSession("username")), &d)

	q := sql.InsertSql(d, ent.InsertCloudEnt)
	if d.EntId > 0 {
		searchMap := sql.SearchMap{}
		searchMap.Put("EntId", d.EntId)
		q = sql.UpdateSql(d, ent.UpdateCloudEnt, searchMap, "CreateTime,CreateUser,Entname")
	}
	_, err = sql.Raw(q).Exec()

	data, msg := util.SaveResponse(err, "名称已经被使用")
	util.SaveOperLog(this.GetSession("username"), *this.Ctx, "保存环境配置 "+msg, d.Entname)
	setEntJson(this, data)
}

func GetEntnameSelectData(isLog bool)  string {
	html := make([]string, 0)
	html = append(html, "<option>--请选择--</option>")
	data := getEntdata()
	for _, v := range data {
		e := util.GetSelectOptionName(v.Entname)
		if isLog {
			e =  util.GetSelectOptionName(v.Description)
		}
		html = append(html, e)
	}
	return strings.Join(html, "\n")
}

// 获取select选项
// 2018-02-06 15:32
func GetEntnameSelect() string {
	return GetEntnameSelectData(false)
}

// 2018-02-06 15:30
// 获取环境信息
func getEntdata() []ent.CloudEnt {
	// 环境数据
	data := make([]ent.CloudEnt, 0)
	q := sql.SearchSql(ent.CloudEnt{}, ent.SelectCloudEnt, sql.SearchMap{})
	sql.Raw(q).QueryRows(&data)
	return data
}

// 获取环境数据
// 2018-01-20 17:45
// router /api/ent/name [get]
func (this *EntController) EntDataName() {
	setEntJson(this, getEntdata())
}

// 环境数据
// @router /api/ent [get]
func (this *EntController) EntDatas() {
	data := make([]ent.CloudEnt, 0)
	searchMap := sql.SearchMap{}
	id := this.Ctx.Input.Param(":id")
	key := this.GetString("search")

	if id != "" {
		searchMap.Put("EntId", id)
	}

	searchSql := sql.SearchSql(ent.CloudEnt{}, ent.SelectCloudEnt, searchMap)
	if key != "" && id == "" {
		key = sql.Replace(key)
		searchSql += strings.Replace(ent.SelectCloudEntWhere, "?", key, -1)
	}

	num, _ := sql.OrderByPagingSql(searchSql,
		"ent_id",
		*this.Ctx.Request,
		&data,
		ent.CloudEnt{})

	r := util.ResponseMap(data,
		sql.Count("cloud_ent", int(num), key),
		this.GetString("draw"))
	setEntJson(this, r)

}

// json
// 删除环境
// 2018-02-05 18:05
// @router /api/ent/:id:int [delete]
func (this *EntController) EntDelete() {
	searchMap := sql.SearchMap{}
	id := this.Ctx.Input.Param(":id")
	searchMap.Put("EntId", id)
	entData := ent.CloudEnt{}
	sql.Raw(sql.SearchSql(entData, ent.SelectCloudEnt, searchMap)).QueryRow(&entData)
	r, err := sql.Raw(sql.DeleteSql(ent.DeleteCloudEnt, searchMap)).Exec()
	data := util.DeleteResponse(
		err,
		*this.Ctx,
		"删除环境"+entData.Entname,
		this.GetSession("username"),
		entData.Entname,
		r)
	setEntJson(this, data)
}

func setEntJson(this *EntController, data interface{}) {
	this.Data["json"] = data
	this.ServeJSON(false)
}
