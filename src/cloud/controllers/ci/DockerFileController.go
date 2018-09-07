package ci

import (
	"github.com/astaxie/beego"
	"cloud/models/ci"
	"cloud/sql"
	"cloud/util"
	"strings"
	"cloud/controllers/base/quota"
	"cloud/userperm"
	"github.com/astaxie/beego/logs"
)

// 2018-01-24 21:32
// 持续集成
type DockerFileController struct {
	beego.Controller
}

// dockerfile管理入口页面
// @router /ci/dockerfile/list [get]
func (this *DockerFileController) DockerFileList() {
	this.TplName = "ci/dockerfile/list.html"
}

// 2018-01-25 10:05
// dockerfile 详情入口页面
// @router /ci/dockerfile/detail/:hi(.*) [get]
func (this *DockerFileController) DockerFileDetail() {
	data := ci.CloudCiDockerfile{}
	searchMap := sql.GetSearchMap("Name", *this.Ctx)
	q := sql.SearchSql(data, ci.SelectCloudCiDockerfile, searchMap)
	sql.Raw(q).QueryRow(&data)
	this.Data["data"] = data
	logs.Info(util.ObjToString(data))
	this.Data["content"] = len(strings.Split(data.Content,"\n"))
	this.TplName = "ci/dockerfile/detail.html"
}


// 生成 镜像服务 html
// 2018-01-26 10:41
func GetDockerFileSelect() string {
	html := make([]string, 0)
	data := GetDockerfileData("")
	for _,v := range data{
		html = append(html, util.GetSelectOptionName(v.Name))
	}
	return strings.Join(html, "")
}


// dockerfile管理添加页面
// @router /ci/dockerfile/add [get]
func (this *DockerFileController) DockerFileAdd() {
	id := this.GetString("FileId")
	copy := this.GetString("Copy")
	update := ci.CloudCiDockerfile{}

	// 更新操作
	if id != "0" {
		searchMap := sql.GetSearchMap("FileId", *this.Ctx)
		q := sql.SearchSql(
			ci.CloudCiDockerfile{},
			ci.SelectCloudCiDockerfile,
			searchMap)
		sql.Raw(q).QueryRow(&update)
	}
	if copy == "1" {
		 update.FileId = 0
	}

	this.Data["data"] = update
	this.TplName = "ci/dockerfile/add.html"
}

// 获取docker数据
// 2018-01-26 11:17
func GetDockerfileData(name string)[]ci.CloudCiDockerfile {
	searchMap := sql.SearchMap{}
	if name != "" {
		searchMap.Put("Name", name)
	}
	// dockerfile数据
	data := make([]ci.CloudCiDockerfile, 0)
	q := sql.SearchSql(ci.CloudCiDockerfile{},
		ci.SelectCloudCiDockerfile,
		searchMap)
	sql.Raw(q).QueryRows(&data)
	return data
}

// 获取dockerfile数据
// 2018-01-24 21:33
// router /api/ci/dockerfile [get]
func (this *DockerFileController) DockerFileData()  {
	setDockerfileJson(this, GetDockerfileData(""))
}


// string
// dockerfile保存
// @router /api/ci/dockerfile [post]
func (this *DockerFileController) DockerFileSave() {
	d := ci.CloudCiDockerfile{}
	err := this.ParseForm(&d)
	if err != nil {
		this.Ctx.WriteString("参数错误" + err.Error())
		return
	}
	util.SetPublicData(d, getDockerfileUser(this), &d)
	
	q := sql.InsertSql(d, ci.InsertCloudCiDockerfile)
	if d.FileId > 0 {
		searchMap := sql.SearchMap{}
		searchMap.Put("FileId", d.FileId)
		q = sql.UpdateSql(
			d,
			ci.UpdateCloudCiDockerfile,
			searchMap,
			ci.UpdateDockerfileExclude)
		DeleteJobCache(d.Name)
	}else{
		status, msg := checkDockerfileQuota(getDockerfileUser(this))
		if !status {
			data := util.ApiResponse(false, msg)
			setDockerfileJson(this, data)
			return
		}
	}
	sql.Raw(q).Exec()
	data,msg := util.SaveResponse(err, "名称已经被使用")
	util.SaveOperLog(this.GetSession("username"),
		*this.Ctx,
		"保存dockerfile配置 "+msg,
		d.Name)
	setDockerfileJson(this, data)
}

// 2018-02-12 08:40
// 检查镜像仓库配额
// 检查资源配额是否够用
func checkDockerfileQuota(username string) (bool,string) {
	quotaDatas := quota.GetUserQuotaData(username, "")
	for _, v := range quotaDatas {
		if v.DockerFileUsed + 1 > v.DockerFileNumber {
			return false, "Dockerfile数量超过配额限制"
		}
	}
	return true, ""
}


// 获取dockerfile数据
// 2018-01-24 21:45
// router /api/ci/dockerfile/name [get]
func (this *DockerFileController) DockerFileDataName()  {
	// dockerfile数据
	data := make([]ci.CloudCiDockerfile, 0)
	q := sql.SearchSql(
		ci.CloudCiDockerfile{},
		ci.SelectCloudCiDockerfile,
		sql.SearchMap{})
	sql.Raw(q).QueryRows(&data)
	setDockerfileJson(this, data)
}

// dockerfile数据
// @router /api/ci/dockerfile  [get]
func (this *DockerFileController) DockerFileDatas() {
	data := make([]ci.CloudCiDockerfile, 0)
	searchMap := sql.SearchMap{}
	id := this.Ctx.Input.Param(":id")
	key := this.GetString("search")
	if id != "" {
		searchMap.Put("FileId", id)
	}
	user := getDockerfileUser(this)

	searchSql := sql.SearchSql(
		ci.CloudCiDockerfile{},
		ci.SelectCloudCiDockerfile,
		searchMap)

	if key != "" && id == "" {
		key = sql.Replace(key)
		searchSql += strings.Replace(ci.SelectDockerfileWhere, "?", key, -1)
	}

	num, err := sql.OrderByPagingSql(searchSql,
		"create_time",
		*this.Ctx.Request,
		&data,
		ci.CloudCiDockerfile{})

	perm := userperm.GetResourceName("DockerFile", user)
	result := make([]ci.CloudCiDockerfile, 0)
	for _, v := range data{
		// 不是自己创建的才检查
		if v.CreateUser != user {
			if ! userperm.CheckPerm(v.Name, "", "", perm)  {
				continue
			}
		}
		result = append(result, v)
	}

	r := util.GetResponseResult(err,
		this.GetString("draw"),
		data,
		sql.Count("cloud_ci_dockerfile", int(num), key))


	setDockerfileJson(this, r)
}

// json
// 删除dockerfile
// 2018-01-24 21:46
// @router /api/ci/dockerfile/:id:int [delete]
func (this *DockerFileController) DockerFileDelete() {
	searchMap := sql.SearchMap{}
	id := this.Ctx.Input.Param(":id")
	searchMap.Put("FileId", id)
	codeData := ci.CloudCiDockerfile{}
	q := sql.SearchSql(codeData, ci.SelectCloudCiDockerfile, searchMap)
	sql.Raw(q).QueryRow(&codeData)

	q = sql.DeleteSql(ci.DeleteCloudCiDockerfile, searchMap)
	r, err := sql.Raw(q).Exec()
	data := util.DeleteResponse(err,
		*this.Ctx,
		"删除dockerfile"+codeData.Name,
		this.GetSession("username"),
		codeData.CreateUser,
		r)
	setDockerfileJson(this, data)
}

// 2018-02-12 9:40
// 获取登录用户
func getDockerfileUser(this *DockerFileController) string {
	return util.GetUser(this.GetSession("username"))
}

// 设置json数据
func setDockerfileJson(this *DockerFileController, data interface{}) {
	this.Data["json"] = data
	this.ServeJSON(false)
}