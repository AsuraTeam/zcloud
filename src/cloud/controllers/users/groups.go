package users

import (
	"cloud/sql"
	"cloud/util"
	"github.com/astaxie/beego"
	"cloud/models/groups"
	"strconv"
	"strings"
)


type GroupsController struct {
	beego.Controller
}

// 部门团队管理入口页面
// @router /users/groups/list [get]
func (this *GroupsController) GroupsList() {
	this.TplName = "users/groups/list.html"
}

// 部门团队管理添加页面
// @router /users/groups/add [get]
func (this *GroupsController) GroupsAdd() {
	id := this.GetString("GroupsId")
	// 更新操作
	if id != "" {
		searchMap := sql.GetSearchMap("GroupsId", *this.Ctx)
		update := groups.CloudUserGroups{}
		q := sql.SearchSql(
			groups.CloudUserGroups{},
			groups.SelectCloudUserGroups,
			searchMap)
		sql.Raw(q).QueryRow(&update)
		this.Data["data"] = update
	}
	this.TplName = "users/groups/add.html"
}

// string
// 部门团队保存
// @router /api/groups [post]
func (this *GroupsController) GroupsSave() {
	d := groups.CloudUserGroups{}
	err := this.ParseForm(&d)
	if err != nil {
		this.Ctx.WriteString("参数错误" + err.Error())
		return
	}
	util.SetPublicData(d, util.GetUser(this.GetSession("username")), &d)

	var q = sql.InsertSql(d, groups.InsertCloudUserGroups)
	if d.GroupsId > 0 {
		searchMap := sql.SearchMap{}
		searchMap.Put("GroupsId", d.GroupsId)
		q = sql.UpdateSql(
			d,
			groups.UpdateCloudUserGroups,
			searchMap,
			groups.UpdateCloudUserGroupsExclude)
	}
	_, err = sql.Raw(q).Exec()

	data,msg := util.SaveResponse(err, "名称已经被使用")
	util.SaveOperLog(
		this.GetSession("username"),
		*this.Ctx, "保存部门团队配置 "+msg,
		d.GroupsName)
	setGroupsJson(this, data)
}

// 2018-01-22 11:04
// 获取组的map数据,关系id和组名
func GetGroupsMap() util.Lock {
	data := getGroupsData()
	gmap := util.Lock{}
	for _,v := range data{
		gmap.Put(strconv.FormatInt(v.GroupsId, 10), v.GroupsName)
	}
	return gmap
}

// 获取组数据
// 2018-02-11 11;11
func GetGroupsSelect() string {
	html := make([]string, 0)
	html = append(html, "<option>--请选择--</option>")
	data := getGroupsData()
	for _,v := range data{
		html = append(html, util.GetSelectOptionName(v.GroupsName))
	}
	return strings.Join(html, "\n")
}


// 获取组的名称和id数据
// 2018-01-22 11:02
func getGroupsData()  []groups.CloudUserGroupsName {
	data := make([]groups.CloudUserGroupsName, 0)
	searchSql := sql.SearchSql(
		groups.CloudUserGroups{},
		groups.SelectCloudUserGroups,
		sql.SearchMap{})
	sql.Raw(searchSql).QueryRows(&data)
	return data
}

// 部门团队名称数据
// @router /api/groups/name [get]
func (this *GroupsController) GetGroupsName() {
	setGroupsJson(this, getGroupsData())
}

// 部门团队名称数据
// @router /api/groups/map [get]
func (this *GroupsController) GetGroupsMap() {
	v := GetGroupsMap()
	setGroupsJson(this, v)
}


// 部门团队数据
// @router /api/groups [get]
func (this *GroupsController) GroupsData() {
	data := make([]groups.CloudUserGroups, 0)
	searchMap := sql.SearchMap{}
	id := this.Ctx.Input.Param(":id")
	key := this.GetString("search")
	if id != "" {
		searchMap.Put("GroupsId", id)
	}

	searchSql := sql.SearchSql(
		groups.CloudUserGroups{},
		groups.SelectCloudUserGroups,
		searchMap)

	if key != "" && id == "" {
		q := groups.SelectCloudUserGroupsWhere
		searchSql += strings.Replace(q, "?", sql.Replace(key), -1)
	}

	num, err := sql.Raw(searchSql).QueryRows(&data)
	var r = util.ResponseMap(data, num, 1)
	if err != nil {
		r = util.ResponseMapError(err.Error())
	}
	setGroupsJson(this, r)
}

// json
// 删除部门团队
// 2018-01-20 10:46
// @router /api/groups/:id:int [delete]
func (this *GroupsController) GroupsDelete() {
	searchMap := sql.SearchMap{}
	id := this.Ctx.Input.Param(":id")
	searchMap.Put("GroupsId", id)
	groupsData := groups.CloudUserGroups{}

	q := sql.SearchSql(
		groupsData,
		groups.SelectCloudUserGroups,
		searchMap)
	sql.Raw(q).QueryRow(&groupsData)

	q = sql.DeleteSql(groups.DeleteCloudUserGroups, searchMap)
	r, err := sql.Raw(q).Exec()
	data := util.DeleteResponse(
		err,
		*this.Ctx, "删除部门团队"+groupsData.GroupsName,
		this.GetSession("username"),
		groupsData.CreateUser,
		r)
	setGroupsJson(this, data)
}

func setGroupsJson(this *GroupsController, data interface{})  {
	this.Data["json"] = data
	this.ServeJSON(false)
}

// 2018-02-11 14:59
// 获取用户部门
func GetUserDept(username string) []string {
	username = sql.Replace(username)
	data := make([]groups.CloudUserGroups, 0)
	q := groups.SelectCloudUserDept
	q = strings.Replace(q, "?", username, -1)
	sql.Raw(q).QueryRows(&data)
	result := make([]string, 0)
	for _, v := range data{
		result = append(result,  v.GroupsName)
	}
	return result
}

// 2018--2011
// 获取组的用户
func GetGroupUsers(groupname []string) []string {
	result := make([]string, 0)
	data := make([]groups.CloudUserGroups, 0)
	q := strings.Replace(groups.SelectGroupUsers, "?", strings.Join(groupname, ","), -1)
	sql.Raw(q).QueryRows(&data)
	for _,v := range data{
		temp := strings.Split(v.Users, ",")
		usersData := make([]string, 0)
		for _,u := range temp {
			usersData = append(usersData, `"` + u + `"`)
		}
		result = append(result, usersData...)
	}
	return result
}