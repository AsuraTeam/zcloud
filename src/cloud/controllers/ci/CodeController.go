package ci

import (
	"github.com/astaxie/beego"
	"cloud/sql"
	"cloud/models/ci"
	"cloud/util"
	"strings"
)

// 2018-01-24 16:36
// 持续集成
type CodeController struct {
	beego.Controller
}

// 代码仓库管理入口页面
// @router /ci/code/list [get]
func (this *CodeController) CodeList() {
	this.TplName = "ci/code/list.html"
}

// 代码仓库管理添加页面
// @router /ci/code/add [get]
func (this *CodeController) CodeAdd() {
	id := this.GetString("RepostitoryId")
	update := ci.CloudCodeRepostitory{}

	// 更新操作
	if id != "0" {
		searchMap := sql.GetSearchMap("RepostitoryId", *this.Ctx)
		q := sql.SearchSql(ci.CloudCodeRepostitory{}, ci.SelectCloudCodeRepostitory, searchMap)
		sql.Raw(q).QueryRow(&update)
	}

	this.Data["data"] = update
	this.TplName = "ci/code/add.html"
}

// 获取代码仓库数据
// 2018-01-20 12:56
// router /api/ci/code [get]
func (this *CodeController) CodeData()  {
	// 代码仓库数据
	data := []ci.CloudCodeRepostitory{}
	q := sql.SearchSql(ci.CloudCodeRepostitory{}, ci.SelectCloudCodeRepostitory, sql.SearchMap{})
	sql.Raw(q).QueryRows(&data)
	this.Data["json"] = data
	this.ServeJSON(false)
}


// string
// 代码仓库保存
// @router /api/ci/code [post]
func (this *CodeController) CodeSave() {
	d := ci.CloudCodeRepostitory{}
	err := this.ParseForm(&d)
	if err != nil {
		this.Ctx.WriteString("参数错误" + err.Error())
		return
	}
	util.SetPublicData(d, util.GetUser(this.GetSession("username")), &d)

	q := sql.InsertSql(d, ci.InsertCloudCodeRepostitory)
	if d.RepostitoryId > 0 {
		searchMap := sql.SearchMap{}
		searchMap.Put("RepostitoryId", d.RepostitoryId)
		q = sql.UpdateSql(d, ci.UpdateCloudCodeRepostitory, searchMap, "CreateTime,CreateCode")
	}
	_, err = sql.Raw(q).Exec()
	data,msg := util.SaveResponse(err, "名称已经被使用")
	util.SaveOperLog(this.GetSession("username"), *this.Ctx, "保存代码仓库配置 "+msg, d.CodeUrl)
	this.Data["json"] = data
	this.ServeJSON(false)
}


// 获取代码仓库数据
// 2018-01-20 17:45
// router /api/ci/code/name [get]
func (this *CodeController) CodeDataName()  {
	// 代码仓库数据
	data := []ci.CloudCodeRepostitory{}
	q := sql.SearchSql(ci.CloudCodeRepostitory{}, ci.SelectCloudCodeRepostitory, sql.SearchMap{})
	sql.Raw(q).QueryRows(&data)
	this.Data["json"] = data
	this.ServeJSON(false)
}

// 代码仓库数据
// @router /api/ci/code [get]
func (this *CodeController) CodeDatas() {
	data := []ci.CloudCodeRepostitory{}
	searchMap := sql.SearchMap{}
	id := this.Ctx.Input.Param(":id")
	key := this.GetString("search")
	if id != "" {
		searchMap.Put("RepostitoryId", id)
	}
	searchSql := sql.SearchSql(ci.CloudCodeRepostitory{}, ci.SelectCloudCodeRepostitory, searchMap)
	if key != "" && id == "" {
		key = sql.Replace(key)
		q := ci.SelectCloudCodeRepostitoryWhere
		searchSql += strings.Replace(q, "?", key, -1)
	}
	num, err := sql.Raw(searchSql).QueryRows(&data)
	
	var r = util.ResponseMap(data, num, 1)
	if err != nil {
		r = util.ResponseMapError(err.Error())
	}
	this.Data["json"] = r
	this.ServeJSON(false)
}

// json
// 删除代码仓库
// 2018-01-20 17:46
// @router /api/ci/code/:id:int [delete]
func (this *CodeController) CodeDelete() {
	searchMap := sql.SearchMap{}
	id := this.Ctx.Input.Param(":id")
	searchMap.Put("RepostitoryId", id)
	codeData := ci.CloudCodeRepostitory{}

	q := sql.SearchSql(codeData, ci.SelectCloudCodeRepostitory, searchMap)
	sql.Raw(q).QueryRow(&codeData)

	q = sql.DeleteSql(ci.DeleteCloudCodeRepostitory, searchMap)
	r, err := sql.Raw(q).Exec()
	data := util.DeleteResponse(err,
		*this.Ctx,
		"删除代码仓库"+codeData.CodeUrl,
		this.GetSession("username"),
		codeData.CreateUser, r)
	this.Data["json"] = data
	this.ServeJSON(false)
}