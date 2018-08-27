package users

import (
	"cloud/sql"
	"cloud/util"
	"github.com/astaxie/beego"
	"cloud/models/perm"
)

type UserPermController struct {
	beego.Controller
}

// 权限管理入口页面
// @router /system/users/perm/list [get]
func (this *UserPermController) PermList() {
	this.TplName = "users/perm/list.html"
}

// 权限管理添加页面
// @router /system/users/perm/add [get]
func (this *UserPermController) PermAdd() {
	id := this.GetString("PermId")
	update := perm.CloudUserPerm{}
	var entHtml string
	// 更新操作
	if id != "0" {
		searchMap := sql.GetSearchMap("PermId", *this.Ctx)
		sql.Raw(sql.SearchSql(perm.CloudUserPerm{}, perm.SelectCloudUserPerm, searchMap)).QueryRow(&update)
		entHtml = util.GetSelectOptionName(update.Ent)
		this.Data["cluster"] = util.GetSelectOptionName(update.ClusterName)
		this.Data["resourceType"] = util.GetSelectOptionName(update.ResourceType)
	}
	this.Data["entname"] = entHtml
	this.Data["data"] = update
	this.TplName = "users/perm/add.html"
}


// string
// 权限保存
// @router /api/users/perm [post]
func (this *UserPermController) PermSave() {
	d := perm.CloudUserPerm{}
	err := this.ParseForm(&d)
	if err != nil {
		this.Ctx.WriteString("参数错误" + err.Error())
		return
	}
	util.SetPublicData(d, util.GetUser(this.GetSession("username")), &d)
	
	q := sql.InsertSql(d, perm.InsertCloudUserPerm)
	if d.PermId > 0 {
		searchMap := sql.SearchMap{}
		searchMap.Put("PermId", d.PermId)
		searchMap.Put("CreateUser", util.GetUser(this.GetSession("username")))
		q = sql.UpdateSql(d, perm.UpdateCloudUserPerm, searchMap, "CreateTime,CreateUser")
	}
	_, err = sql.Raw(q).Exec()

	data, msg := util.SaveResponse(err, "名称已经被使用")
	util.SaveOperLog(this.GetSession("username"), *this.Ctx, "保存用户权限配置 "+msg, d.Description)
	setPermJson(this, data)
}

// 2018-08-23 08:46
// 权限数据
// @router /api/users/perm [get]
func (this *UserPermController) PermData() {
	data := make([]perm.CloudUserPerm, 0)
	searchMap := sql.SearchMap{}
	searchMap.Put("CreateUser", util.GetUser(this.GetSession("username")))
	key := this.GetString("search")
	searchSql := sql.SearchSql(perm.CloudUserPerm{}, perm.SelectCloudUserPerm, searchMap)
	if key != ""  {
		key = sql.Replace(key)
		searchSql += " where 1=1 and (description like \"%" + key + "%\")"
	}
	num, _ := sql.OrderByPagingSql(searchSql, "perm_id", *this.Ctx.Request, &data, perm.CloudUserPerm{})
    r := util.ResponseMap(data, sql.Count("cloud_user_perm", int(num), key), this.GetString("draw"))
	setPermJson(this, r)
}

// json
// 删除权限
// 2018-08-23 09:10
// @router /api/users/perm/:id:int [delete]
func (this *UserPermController) PermDelete() {
	searchMap := sql.GetSearchMap("PermId", *this.Ctx)
	searchMap.Put("CreateUser", util.GetUser(this.GetSession("username")))
	permData := perm.CloudUserPerm{}
	sql.Raw(sql.SearchSql(permData, perm.SelectCloudUserPerm, searchMap)).QueryRow(&permData)
	r, err := sql.Raw(sql.DeleteSql(perm.DeleteCloudUserPerm, searchMap)).Exec()
	data := util.DeleteResponse(err, *this.Ctx, "删除权限"+permData.Description, this.GetSession("username"), permData.CreateUser, r)
	setPermJson(this, data)
}

func setPermJson(this *UserPermController, data interface{}) {
	this.Data["json"] = data
	this.ServeJSON(false)
}