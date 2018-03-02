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
	data := []index.CloudAuthorityUser{}
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
		this.Ctx.WriteString("参数错误" + err.Error())
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
	data := []index.CloudAuthorityUser{}
	q := sql.SearchSql(index.CloudAuthorityUser{}, index.SelectDockerCloudAuthorityUser, sql.SearchMap{})
	sql.Raw(q).QueryRows(&data)
	setUserJson(this, data)
}

// 用户数据
// @router /api/user [get]
func (this *UserController) UserDatas() {
	data := []index.DockerCloudAuthorityUser{}
	searchMap := sql.SearchMap{}
	id := this.Ctx.Input.Param(":id")
	key := this.GetString("search")
	if id != "" {
		searchMap.Put("UserId", id)
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

	result := []index.DockerCloudAuthorityUser{}
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
