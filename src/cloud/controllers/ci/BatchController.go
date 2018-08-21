package ci

import (
	"github.com/astaxie/beego"
	"cloud/models/ci"
	"cloud/sql"
	"cloud/util"
	"cloud/controllers/base/cluster"
)

// 2018-01-24 21:32
// 持续集成
type BatchController struct {
	beego.Controller
}

// batch管理入口页面
// @router /ci/batch/list [get]
func (this *BatchController) BatchList() {
	this.TplName = "ci/batch/list.html"
}

// batch管理添加页面
// @router /ci/batch/add [get]
func (this *BatchController) BatchAdd() {
	id := this.GetString("BatchId")
	update := ci.CloudCiBatchJob{}
	this.Data["cluster"] = cluster.GetClusterSelect()
	// 更新操作
	if id != "0" {
		searchMap := sql.GetSearchMap("BatchId", *this.Ctx)
		q := sql.SearchSql(
			ci.CloudCiBatchJob{},
			ci.SelectCloudCiBatchJob,
			searchMap)
		sql.Raw(q).QueryRow(&update)
		this.Data["cluster"] = util.GetSelectOptionName(update.Cluster)
	}

	this.Data["data"] = update
	this.TplName = "ci/batch/add.html"
}


// string
// batch保存
// @router /api/ci/batch [post]
func (this *BatchController) BatchSave() {
	d := ci.CloudCiBatchJob{}
	err := this.ParseForm(&d)
	if err != nil {
		this.Ctx.WriteString("参数错误" + err.Error())
		return
	}
	util.SetPublicData(d, getBatchUser(this), &d)
	
	q := sql.InsertSql(d, ci.InsertCloudCiBatchJob)
	if d.BatchId > 0 {
		searchMap := sql.SearchMap{}
		searchMap.Put("BatchId", d.BatchId)
		q = sql.UpdateSql(
			d,
			ci.UpdateCloudCiBatchJob,
			searchMap,
			"")
	}
	sql.Raw(q).Exec()
	data,msg := util.SaveResponse(err, "名称已经被使用")
	util.SaveOperLog(this.GetSession("username"),
		*this.Ctx,
		"保存batch配置 "+msg,
		d.BatchName)
	setBatchJson(this, data)
}



// batch数据
// @router /api/ci/batch  [get]
func (this *BatchController) BatchData() {
	data := make([]ci.CloudCiBatchJob, 0)
	searchMap := sql.SearchMap{}
	id := this.Ctx.Input.Param(":id")
	key := this.GetString("search")
	if id != "" {
		searchMap.Put("BatchId", id)
	}
	user := getBatchUser(this)

	searchMap.Put("CreateUser", user)
	searchSql := sql.SearchSql(
		ci.CloudCiBatchJob{},
		ci.SelectCloudCiBatchJob,
		searchMap)

	num, err := sql.OrderByPagingSql(searchSql,
		"create_time",
		*this.Ctx.Request,
		&data,
		ci.CloudCiBatchJob{})

	r := util.GetResponseResult(err,
		this.GetString("draw"),
		data,
		sql.Count("cloud_ci_batch_job", int(num), key))

	setBatchJson(this, r)
}

// json
// 删除batch
// 2018-01-24 21:46
// @router /api/ci/batch/:id:int [delete]
func (this *BatchController) BatchDelete() {
	searchMap := sql.SearchMap{}
	id := this.Ctx.Input.Param(":id")
	searchMap.Put("BatchId", id)
	codeData := ci.CloudCiBatchJob{}
	q := sql.SearchSql(codeData, ci.SelectCloudCiBatchJob, searchMap)
	sql.Raw(q).QueryRow(&codeData)

	q = sql.DeleteSql(ci.DeleteCloudCiBatchJob, searchMap)
	r, err := sql.Raw(q).Exec()
	data := util.DeleteResponse(err,
		*this.Ctx,
		"删除批量任务"+codeData.BatchName,
		this.GetSession("username"),
		codeData.CreateUser,
		r)
	setBatchJson(this, data)
}

// 2018-02-12 9:40
// 获取登录用户
func getBatchUser(this *BatchController) string {
	return util.GetUser(this.GetSession("username"))
}

// 设置json数据
func setBatchJson(this *BatchController, data interface{}) {
	this.Data["json"] = data
	this.ServeJSON(false)
}