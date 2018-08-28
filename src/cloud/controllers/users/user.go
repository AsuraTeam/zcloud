package users

import (
	"cloud/sql"
	"cloud/util"
	"github.com/astaxie/beego"
	"cloud/models/index"
	"strings"
)

type UserController struct {
	beego.Controller
}

// 用户管理入口页面
// @router /users/user/list [get]
func (this *UserController) UserList() {
	this.TplName = "users/user/list.html"
}

// 用户管理添加页面
// @router /users/user/add [get]
func (this *UserController) UserAdd() {
	id := this.GetString("UserId")
	update := index.DockerCloudAuthorityUser{}
	update.IsDel = 0

	// 更新操作
	if id != "" {
		searchMap := sql.GetSearchMap("UserId", *this.Ctx)
		update.Pwd = "******"
		q := sql.SearchSql(
			index.DockerCloudAuthorityUser{},
			index.SelectDockerCloudAuthorityUser,
			searchMap)
		sql.Raw(q).QueryRow(&update)
	}
	if update.IsDel == 1 {
		this.Data["userchecked"] = "checked"
	}
	this.Data["data"] = update
	this.TplName = "users/user/add.html"
}

// 获取用户数据
// 2018-02-11 11;02
func GetUserSelect() string {
	html := make([]string, 0)
	html = append(html, "<option>--请选择--</option>")
	data := getUserData()
	for _,v := range data{
		html = append(html, util.GetSelectOptionName(v.UserName))
	}
	return strings.Join(html, "\n")
}

// 2018-02-11 11;03
func getUserData() []index.CloudAuthorityUser {
	// 用户数据
	data := make([]index.CloudAuthorityUser, 0)
	q := sql.SearchSql(
		index.CloudAuthorityUser{},
		index.SelectDockerCloudAuthorityUser,
		sql.SearchMap{})
	sql.Raw(q).QueryRows(&data)
	return data
}

// 获取用户数据
// 2018-01-20 12:56
// router /api/users [get]
func (this *UserController) UserData() {
	data := getUserData()
	setUserJson(this, data)
}



// string
// 用户保存
// @router /api/user [post]
func (this *UserController) UserSave() {
	d := index.DockerCloudAuthorityUser{}
	err := this.ParseForm(&d)
	if err != nil {
		setUserJson(this, util.ApiResponse(false,err.Error()))
		return
	}

	u := util.GetUser(this.GetSession("username"))
	if u != "admin" {
		setUserJson(this, util.ApiResponse(false, "无权限"))
		return
	}
	util.SetPublicData(d, util.GetUser(this.GetSession("username")), &d)

	if d.Pwd != "******" {
		d.Pwd = util.Md5String(d.Pwd)
	}
	q := sql.InsertSql(d, index.InsertDockerCloudAuthorityUser)
	if d.UserId > 0 {
		searchMap := sql.SearchMap{}
		searchMap.Put("UserId", d.UserId)
		q = sql.UpdateSql(
			d,
			index.UpdateDockerCloudAuthorityUser,
			searchMap, "CreateTime,CreateUser")
	}
	_, err = sql.Raw(q).Exec()

	data, msg := util.SaveResponse(err, "名称已经被使用")
	util.SaveOperLog(this.GetSession("username"), *this.Ctx, "保存用户配置 "+msg, d.UserName)
	setUserJson(this, data)
}

// 获取用户数据
// 2018-01-20 17:45
// router /api/users/name [get]
func (this *UserController) UserDataName() {
	// 用户数据
	data := make([]index.CloudAuthorityUser, 0)
	q := sql.SearchSql(index.CloudAuthorityUser{}, index.SelectDockerCloudAuthorityUser, sql.SearchMap{})
	sql.Raw(q).QueryRows(&data)
	setUserJson(this, data)
}

// 用户数据
// @router /api/user [get]
func (this *UserController) UserDatas() {
	data := make([]index.DockerCloudAuthorityUser, 0)
	searchMap := sql.SearchMap{}
	id := this.Ctx.Input.Param(":id")
	key := this.GetString("search")
	if id != "" {
		searchMap.Put("UserId", id)
	}

	u := util.GetUser(this.GetSession("username"))
	if u != "admin" {
		searchMap.Put("UserName", u)
	}
	searchSql := sql.SearchSql(
		index.DockerCloudAuthorityUser{},
		index.SelectDockerCloudAuthorityUser,
		searchMap)

	if key != "" && id == "" {
		key = sql.Replace(key)
		searchSql += strings.Replace(index.SelectUserWhere, "?", key, -1)
	}

	num, _ := sql.OrderByPagingSql(
		searchSql, "user_id",
		*this.Ctx.Request,
		&data,
		index.DockerCloudAuthorityUser{})

	result := make([]index.DockerCloudAuthorityUser, 0)
	for _, v := range data {
		v.Pwd = "******"
		result = append(result, v)
	}

	r := util.ResponseMap(
		result,
		sql.Count("cloud_authority_user", int(num), key),
		this.GetString("draw"))
	setUserJson(this, r)

}

// 获取用户token
// 2018-08-28 09:12
// @router /system/users/user/token/:id:int [get]
func (this *UserController) UserToken() {
	searchMap := sql.GetSearchMap("UserId", *this.Ctx)
	searchMap.Put("UserName", util.GetUser(this.GetSession("username")))
	userData := index.DockerCloudAuthorityUser{}
	q := sql.SearchSql(userData, index.SelectDockerCloudAuthorityUser, searchMap)
	sql.Raw(q).QueryRow(&userData)
	if len(userData.Token) == 0 {
		userData.Token = util.Md5String(userData.UserName + util.GetDate())
		q = sql.UpdateSql(userData, index.UpdateDockerCloudAuthorityUser, searchMap, "")
		sql.Exec(q)
	}
	this.Data["data"] = userData
	this.TplName = "users/user/token.html"
}

// json
// 删除用户
// 2018-01-20 17:46
// @router /api/user/:id:int [delete]
func (this *UserController) UserDelete() {
	searchMap := sql.GetSearchMap("UserId", *this.Ctx)
	userData := index.DockerCloudAuthorityUser{}

	q := sql.SearchSql(userData, index.SelectDockerCloudAuthorityUser, searchMap)
	sql.Raw(q).QueryRow(&userData)

	q = sql.DeleteSql(index.DeleteDockerCloudAuthorityUser, searchMap)
	r, err := sql.Raw(q).Exec()

	data := util.DeleteResponse(
		err,
		*this.Ctx,
		"删除用户"+userData.UserName,
		this.GetSession("username"),
		userData.CreateUser, r)
	setUserJson(this, data)
}

func setUserJson(this *UserController, data interface{}) {
	this.Data["json"] = data
	this.ServeJSON(false)
}
