package perm

import (
	"cloud/sql"
	"cloud/util"
	"github.com/astaxie/beego"
	"cloud/models/perm"
	"strings"
)

type PermRoleController struct {
	beego.Controller
}

// 角色管理入口页面
// @router /perm/role/list [get]
func (this *PermRoleController) PermRoleList() {
	this.TplName = "perm/role/list.html"
}

// @router /perm/role/perm/add [get]
func (this *PermRoleController) PermRoleAddList() {
	update := perm.CloudPermRole{}
	id , err := this.GetInt64("RoleId")
	if err != nil {
		setPermRoleJson(this, util.ApiResponse(false, "资源不存在"))
		return
	}
	update.RoleId = id
	update.Permissions =  util.ObjToString(GetRolePermMap(id))
	this.Data["data"] = update
	this.TplName = "perm/role/perm/add.html"
}

// 角色分配用户页面
// 2018-09-11 10:22
// @router /perm/role/user/add [get]
func (this *PermRoleController) PermRoleUserList() {
	update := perm.CloudPermRoleUser{}
	id , err := this.GetInt64("RoleId")
	if err != nil {
		setPermRoleJson(this, util.ApiResponse(false, "资源不存在"))
		return
	}
	update.RoleId = id
	this.Data["data"] = update
	this.TplName = "perm/role/user/add.html"
}

// 2018-09-11 08:44
// 获取权限map
func GetRolePermMap(roleId int64)  map[string]interface{} {
	rolePerm := make([]perm.CloudPermRolePerm, 0)
	searchMap := sql.SearchMap{}
	searchMap.Put("RoleId", roleId)
	q := sql.SearchSql(perm.CloudPermRolePerm{}, perm.SelectCloudPermRolePerm, searchMap)
	sql.Raw(q).QueryRows(&rolePerm)
	roleMap := map[string]interface{}{}
	for _, v := range rolePerm {
		roleMap[v.PermName] = "1"
	}
	return roleMap
}

// 角色管理添加页面
// @router /perm/role/add [get]
func (this *PermRoleController) PermRoleAdd() {
	id := this.GetString("RoleId")
	update := perm.CloudPermRole{}
	// 更新操作
	if id != "" {
		searchMap := sql.GetSearchMap("RoleId", *this.Ctx)
		sql.Raw(sql.SearchSql(perm.CloudPermRole{}, perm.SelectCloudPermRole, searchMap)).QueryRow(&update)
	}
	this.Data["data"] = update
	this.TplName = "perm/role/add.html"
}

// 获取角色数据
// 2018-02-06 8:56
// router /api/perm/role [get]
func (this *PermRoleController) PermRoleData() {
	// 角色数据
	data :=make([]perm.CloudPermRole, 0)
	q := sql.SearchSql(perm.CloudPermRole{}, perm.SelectCloudPermRole, sql.SearchMap{})
	sql.Raw(q).QueryRows(&data)
	setPermRoleJson(this, data)
}

// string
// 角色保存
// @router /api/perm/role [post]
func (this *PermRoleController) PermRoleSave() {
	d := perm.CloudPermRole{}
	err := this.ParseForm(&d)
	if err != nil {
		this.Ctx.WriteString("参数错误" + err.Error())
		return
	}
	util.SetPublicData(d, util.GetUser(this.GetSession("username")), &d)
	
	q := sql.InsertSql(d, perm.InsertCloudPermRole)
	if d.RoleId > 0 {
		searchMap := sql.SearchMap{}
		searchMap.Put("RoleId", d.RoleId)
		q = sql.UpdateSql(d, perm.UpdateCloudPermRole, searchMap, "CreateTime,CreateUser")
	}
	_, err = sql.Raw(q).Exec()

	data, msg := util.SaveResponse(err, "名称已经被使用")
	util.SaveOperLog(this.GetSession("username"), *this.Ctx, "保存角色配置 "+msg, d.RoleName)
	setPermRoleJson(this, data)
}

// 2018-09-11 08:17
// 角色权限保存
// @router /api/perm/role/perm/:id:int [post]
func (this *PermRoleController) PermRoleSavePerm() {
	d := perm.CloudPermRolePerm{}
	err := this.ParseForm(&d)
	if err != nil {
		this.Ctx.WriteString("参数错误" + err.Error())
		return
	}

	user := util.GetUser(this.GetSession("username"))
	util.SetPublicData(d, user, &d)


	if d.RoleId > 0 {
		searchMap := sql.SearchMap{}
		searchMap.Put("RoleId", d.RoleId)
		if user != "admin" {
			searchMap.Put("CreateUser", user)
		}
		q := sql.DeleteSql(perm.DeleteCloudPermRolePerm, searchMap)
		sql.Exec(q)
	}

	role := this.GetString("PermName")
	if len(role) > 0 {
		roles := strings.Split(role, ",")
		for _, v := range roles {
			i := perm.CloudPermRolePerm{}
			i.CreateTime = d.CreateTime
			i.CreateUser = d.CreateUser
			i.RoleId = d.RoleId
			i.PermName = v
			insert := sql.InsertSql(i, perm.InsertCloudPermRolePerm)
			sql.Exec(insert)
		}
	}

	data, msg := util.SaveResponse(err, "名称已经被使用")
	util.SaveOperLog(this.GetSession("username"), *this.Ctx, "保存角色权限 "+msg, "")
	setPermRoleJson(this, data)
}

// 获取角色数据
// 2018-02-06 08:36
// router /api/perm/role/name [get]
func (this *PermRoleController) PermRoleDataName() {
	// 角色数据
	data :=make([]perm.CloudPermRole, 0)
	q := sql.SearchSql(perm.CloudPermRole{}, perm.SelectCloudPermRole, sql.SearchMap{})
	sql.Raw(q).QueryRows(&data)
	setPermRoleJson(this, data)
}

// 角色数据
// @router /api/perm/role [get]
func (this *PermRoleController) PermRoleDatas() {
	data :=make([]perm.CloudPermRole, 0)
	searchMap := sql.SearchMap{}
	id := this.Ctx.Input.Param(":id")
	key := this.GetString("search")
	if id != "" {
		searchMap.Put("RoleId", id)
	}
	searchSql := sql.SearchSql(perm.CloudPermRole{}, perm.SelectCloudPermRole, searchMap)
	if key != "" && id == "" {
		key = sql.Replace(key)
		searchSql += " where 1=1 and (description like \"%" + key + "%\")"
	}
	num, _ := sql.OrderByPagingSql(searchSql, "role_id", *this.Ctx.Request, &data, perm.CloudPermRole{})
    r := util.ResponseMap(data, sql.Count("cloud_perm_role", int(num), key), this.GetString("draw"))
	setPermRoleJson(this, r)
}

// json
// 删除角色
// 2018-02-06 08:36
// @router /api/perm/role/:id:int [delete]
func (this *PermRoleController) PermRoleDelete() {
	searchMap := sql.GetSearchMap("RoleId", *this.Ctx)
	permData := perm.CloudPermRole{}
	sql.Raw(sql.SearchSql(permData, perm.SelectCloudPermRole, searchMap)).QueryRow(&permData)
	r, err := sql.Raw(sql.DeleteSql(perm.DeleteCloudPermRole, searchMap)).Exec()
	data := util.DeleteResponse(err, *this.Ctx, "删除角色"+permData.RoleName, this.GetSession("username"), permData.CreateUser, r)
	setPermRoleJson(this, data)
}

func setPermRoleJson(this *PermRoleController, data interface{}) {
	this.Data["json"] = data
	this.ServeJSON(false)
}