package perm

import (
	"cloud/sql"
	"cloud/util"
	"github.com/astaxie/beego"
	"cloud/models/perm"
)

type PermController struct {
	beego.Controller
}

// 权限管理入口页面
// @router /perm/list [get]
func (this *PermController) PermList() {
	this.TplName = "perm/list.html"
}

// 权限管理添加页面
// @router /perm/add [get]
func (this *PermController) PermAdd() {
	id := this.GetString("PermId")
	update := perm.CloudPerm{}
	// 更新操作
	if id != "" {
		searchMap := sql.GetSearchMap("PermId", *this.Ctx)
		sql.Raw(sql.SearchSql(perm.CloudPerm{}, perm.SelectCloudPerm, searchMap)).QueryRow(&update)
	}
	this.Data["data"] = update
	this.TplName = "perm/add.html"
}

// 获取权限数据
// 2018-02-06 8:56
// router /api/perms [get]
func (this *PermController) PermData() {
	// 权限数据
	data := make([]perm.CloudPerm, 0)
	q := sql.SearchSql(perm.CloudPerm{}, perm.SelectCloudPerm, sql.SearchMap{})
	sql.Raw(q).QueryRows(&data)
	setPermJson(this, data)
}

// string
// 权限保存
// @router /api/perm [post]
func (this *PermController) PermSave() {
	d := perm.CloudPerm{}
	err := this.ParseForm(&d)
	if err != nil {
		this.Ctx.WriteString("参数错误" + err.Error())
		return
	}
	util.SetPublicData(d, util.GetUser(this.GetSession("username")), &d)
	
	q := sql.InsertSql(d, perm.InsertCloudPerm)
	if d.PermId > 0 {
		searchMap := sql.SearchMap{}
		searchMap.Put("PermId", d.PermId)
		q = sql.UpdateSql(d, perm.UpdateCloudPerm, searchMap, "CreateTime,CreatePerm")
	}
	_, err = sql.Raw(q).Exec()

	data, msg := util.SaveResponse(err, "名称已经被使用")
	util.SaveOperLog(this.GetSession("username"), *this.Ctx, "保存权限配置 "+msg, d.Description)
	setPermJson(this, data)
}


// 获取权限数据
// 2018-02-06 08:30
// router /api/perms/name [get]
func (this *PermController) PermDataName() {
	// 权限数据
	data := make([]perm.CloudPerm, 0)
	q := sql.SearchSql(perm.CloudPerm{}, perm.SelectCloudPerm, sql.SearchMap{})
	sql.Raw(q).QueryRows(&data)
	setPermJson(this, data)
}

// 权限数据
// @router /api/perm [get]
func (this *PermController) PermDatas() {
	data := make([]perm.CloudPerm, 0)
	searchMap := sql.SearchMap{}
	id := this.Ctx.Input.Param(":id")
	key := this.GetString("search")
	if id != "" {
		searchMap.Put("PermId", id)
	}
	searchSql := sql.SearchSql(perm.CloudPerm{}, perm.SelectCloudPerm, searchMap)
	if key != "" && id == "" {
		key = sql.Replace(key)
		searchSql += " where 1=1 and (description like \"%" + key + "%\")"
	}
	num, _ := sql.OrderByPagingSql(searchSql, "user_id", *this.Ctx.Request, &data, perm.CloudPerm{})
    r := util.ResponseMap(data, sql.Count("cloud_perm", int(num), key), this.GetString("draw"))
	setPermJson(this, r)
}

// json
// 删除权限
// 2018-02-06 08:29
// @router /api/perm/:id:int [delete]
func (this *PermController) PermDelete() {
	searchMap := sql.GetSearchMap("PermId", *this.Ctx)
	permData := perm.CloudPerm{}
	sql.Raw(sql.SearchSql(permData, perm.SelectCloudPerm, searchMap)).QueryRow(&permData)
	r, err := sql.Raw(sql.DeleteSql(perm.DeleteCloudPerm, searchMap)).Exec()
	data := util.DeleteResponse(err, *this.Ctx, "删除权限"+permData.Description, this.GetSession("username"), permData.CreateUser, r)
	setPermJson(this, data)
}

func setPermJson(this *PermController, data interface{}) {
	this.Data["json"] = data
	this.ServeJSON(false)
}