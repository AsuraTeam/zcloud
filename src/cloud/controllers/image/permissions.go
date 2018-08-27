package registry

import (
	"github.com/astaxie/beego"
	"cloud/sql"
	"cloud/util"
	"cloud/models/registry"
	"cloud/controllers/users"
	"strings"
)

type RegistryPermController struct {
	beego.Controller
}

// 权限管理入口页面
// @router /image/registry/list [get]
func (this *RegistryPermController) RegistryPermList() {
	this.TplName = "image/permissions/list.html"
}

// @router /image/registry/add [get]
func (this *RegistryPermController) RegistryPermAdd() {
	update := registry.CloudRegistryPermissions{}
	id := this.GetString("PermissionsId")
	update.Project = this.GetString("Project")
	update.ServiceName = this.GetString("ServerDomain")
	// 更新操作
	if id != "0" {
		searchMap := sql.GetSearchMap("PermissionsId", *this.Ctx)
		q := sql.SearchSql(
			registry.CloudRegistryPermissions{},
			registry.SelectCloudRegistryPermissions,
			searchMap)

		sql.Raw(q).QueryRow(&update)
		this.Data["readonly"] = "readonly"
		if strings.Contains(update.Action, "pull") {
			this.Data["pull"] = "checked"
		}
		if strings.Contains(update.Action, "push") {
			this.Data["push"] = "checked"
		}
	}
	if update.Project != "" {
		this.Data["project"] = "readonly"
	}
	this.Data["data"] = update
	this.TplName = "image/permissions/add.html"
}

// json
// @router /api/registry [post]
func (this *RegistryPermController) RegistryPermSave() {
	d := registry.CloudRegistryPermissions{}
	err := this.ParseForm(&d)
	if err != nil {
		this.Ctx.WriteString("参数错误" + err.Error())
		return
	}
	searchMap := sql.SearchMap{}
	searchMap.Put("PermissionsId", d.PermissionsId)
	masterData := make([]registry.CloudRegistryPermissions, 0)

	q := sql.SearchSql(d, registry.SelectCloudRegistryPermissions, searchMap)
	sql.Raw(q).QueryRows(&masterData)
	util.SetPublicData(d, util.GetUser(this.GetSession("username")), &d)


	q = sql.InsertSql(d, registry.InsertCloudRegistryPermissions)
	if d.PermissionsId > 0 {

		q = sql.UpdateSql(d,
			registry.UpdateCloudRegistryPermissions,
			searchMap,
			registry.UpdateCloudRedisPermExclude)
	}
	sql.Raw(q).Exec()
	sql.Exec(registry.UpdateClusterName)
	data, msg := util.SaveResponse(err, "名称已经被使用")
	util.SaveOperLog(
		this.GetSession("username"),
		*this.Ctx, "操作仓库权限 "+msg,
		d.ServiceName)

	setPermissonsJson(this, data)
}

// 仓库服务器数据
// @router /api/registry [get]
func (this *RegistryPermController) RegistryPerm() {
	data := make([]registry.CloudRegistryPermissions, 0)
	searchMap := sql.SearchMap{}
	project := this.GetString("project")
	key := this.GetString("search")
	if project != "" {
		searchMap.Put("Project", project)
	}
	searchSql := sql.SearchSql(registry.CloudRegistryPermissions{},
		registry.SelectCloudRegistryPermissions,
		searchMap)

	if project == "" {
		searchSql += " where 1=1 "
	}
	if key != "" {
		key = sql.Replace(key)
		q := strings.Replace(registry.SelectCloudRegistryPermWhere, "?", key , -1)
		searchSql += q
	}
	groupsMap := users.GetGroupsMap()
	num, err := sql.OrderByPagingSql(searchSql,
		"create_time",
		*this.Ctx.Request,
		&data,
		registry.CloudRegistryPermissions{})

	result := make([]registry.CloudRegistryPermissions, 0)
	for _,v := range data{
		gdata := make([]string,0)
		for _, g := range strings.Split(v.GroupsName, ","){
			gdata = append(gdata, util.ObjToString(groupsMap.GetV(g)))
		}
		v.GroupsName = strings.Join(gdata,",")
		result = append(result, v)
	}


	r := util.ResponseMap(result,
		sql.Count("cloud_registry_permissions", int(num), key),
		this.GetString("draw"))

	if err != nil {
		r = util.ResponseMapError(err.Error())
	}
	setPermissonsJson(this, r)
}

// @router /api/registry/delete [*]
func (this *RegistryPermController) RegistryPermDelete() {
	searchMap := sql.GetSearchMap("PermissionsId", *this.Ctx)
	registrData := registry.CloudRegistryPermissions{}

	q := sql.SearchSql(
		registrData,
		registry.SelectCloudRegistryPermissions,
		searchMap)
	sql.Raw(q).QueryRow(&registrData)

	q = sql.DeleteSql(registry.DeleteCloudRegistryPermissions, searchMap)
	r, err := sql.Raw(q).Exec()
	data := util.DeleteResponse(err,
		*this.Ctx, "删除仓库服务,名称:"+registrData.ServiceName,
		this.GetSession("username"),
		registrData.ServiceName, r)
	setPermissonsJson(this, data)
}

func setPermissonsJson(this *RegistryPermController, data interface{})  {
	this.Data["json"] = data
	this.ServeJSON(false)
}