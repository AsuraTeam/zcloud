package ci

import (
	"github.com/astaxie/beego"
	"cloud/sql"
	"cloud/util"
	"strings"
	"cloud/models/ci"
	"time"
	"github.com/garyburd/redigo/redis"
	"cloud/controllers/users"
	"cloud/cache"
)

type CiPermController struct {
	beego.Controller
}

// 权限管理入口页面
// @router /ci/service/perm/list [get]
func (this *CiPermController) CiPermList() {
	this.TplName = "ci/service/perm/list.html"
}

// 发布权限添加页面
// @router /ci/service/perm/add [get]
func (this *CiPermController) CiPermAdd() {
	update := ci.CloudCiPerm{}
	id := this.GetString("PermId")
	// 更新操作
	if id != "0" {
		searchMap := sql.GetSearchMap("PermId", *this.Ctx)
		q := sql.SearchSql(
			ci.CloudCiPerm{},
			ci.SelectCloudCiPerm,
			searchMap)
		sql.Raw(q).QueryRow(&update)
		var domainSelect string
		for _, v := range strings.Split(update.Datas, ",") {
			domainSelect += util.GetSelectOptionName(v)
		}
		this.Data["domainSelect"] = domainSelect
	}
	domains := getServiceDomainSelect(update.Datas)
	domains = strings.Replace(domains, "</option>", "</option>\n", -1)
	domains = strings.TrimSuffix(domains, "\n")
	this.Data["domains"] = domains
	this.Data["data"] = update
	this.TplName = "ci/service/perm/add.html"
}

// 保存数据
// @router /api/ci/service/perm [post]
func (this *CiPermController) CiPermSave() {
	d := ci.CloudCiPerm{}
	err := this.ParseForm(&d)
	if err != nil {
		this.Ctx.WriteString("参数错误" + err.Error())
		return
	}
	searchMap := sql.SearchMap{}
	searchMap.Put("PermId", d.PermId)
	masterData := make([]ci.CloudCiPerm, 0)

	q := sql.SearchSql(d, ci.SelectCloudCiPerm, searchMap)
	sql.Raw(q).QueryRows(&masterData)
	util.SetPublicData(d, getPermUser(this), &d)

	q = sql.InsertSql(d, ci.InsertCloudCiPerm)
	if d.PermId > 0 {
		q = sql.UpdateSql(d,
			ci.UpdateCloudCiPerm,
			searchMap,
			"CreateUser,CreateTime")
	}
	sql.Raw(q).Exec()

	data, msg := util.SaveResponse(err, "名称已经被使用")
	util.SaveOperLog(
		getPermUser(this),
		*this.Ctx, "保存发布权限 "+msg,
		d.Username+d.GroupsName)

	setPermissonsJson(this, data)
}

// 权限数据
// @router /api/ci/service/perm [get]
func (this *CiPermController) CiPerm() {
	data := make([]ci.CloudCiPerm, 0)
	searchMap := sql.SearchMap{}
	key := this.GetString("search")
	searchSql := sql.SearchSql(ci.CloudCiPerm{},
		ci.SelectCloudCiPerm,
		searchMap)

	if key != "" {
		key = sql.Replace(key)
		q := strings.Replace(ci.SelectCloudCiPermWhere, "?", key, -1)
		searchSql += q
	}
	num, err := sql.OrderByPagingSql(searchSql,
		"create_time",
		*this.Ctx.Request,
		&data,
		ci.CloudCiPerm{})

	r := util.ResponseMap(data,
		sql.Count("cloud_ci_perm", int(num), key),
		this.GetString("draw"))

	if err != nil {
		r = util.ResponseMapError(err.Error())
	}
	setPermissonsJson(this, r)
	go setPermCache()
}

// @router /api/ci/service/perm/delete [*]
func (this *CiPermController) CiPermDelete() {
	searchMap := sql.GetSearchMap("PermId", *this.Ctx)
	permData := ci.CloudCiPerm{}

	q := sql.SearchSql(
		permData,
		ci.SelectCloudCiPerm,
		searchMap)
	sql.Raw(q).QueryRow(&permData)

	q = sql.DeleteSql(ci.DeleteCloudCiPerm, searchMap)
	r, err := sql.Raw(q).Exec()
	data := util.DeleteResponse(err,
		*this.Ctx, "删除权限信息,用户名:"+permData.Username+permData.GroupsName,
		this.GetSession("username"),
		permData.GroupsName, r)
	cache.PermCache.Delete(permData.Username)
	cache.PermCache.Delete(permData.GroupsName)
	setPermissonsJson(this, data)
}

func setPermissonsJson(this *CiPermController, data interface{}) {
	this.Data["json"] = data
	this.ServeJSON(false)
}

func getPermUser(this *CiPermController) string {
	return util.GetUser(this.GetSession("username"))
}



func setPermCache() {
	data := make([]ci.CloudCiPerm, 0)
	sql.Raw(ci.SelectCloudCiPerm).QueryRows(&data)
	for _, v := range data {
		if v.Username != "" {
			cache.PermCache.Put(v.Username, util.ObjToString(v.Datas), time.Second*86400*5)
		}
		if v.GroupsName != "" {
			cache.PermCache.Put(v.GroupsName, util.ObjToString(v.Datas), time.Second*86400*5)
		}
	}
}

// 2018-02-18 21:24
// 检查权限是否存在
func checkPermExists(r interface{}, domain string) bool {
	redisR, err := redis.String(r, nil)
	if err == nil {
		if strings.Contains(redisR, domain) {
			return true
		}
	}
	return false
}

// 2018-02-18
// 缓存和获取用户部门数据
func getUserDept(username string) []string {
	var depts []string
	r := cache.PermCache.Get(username + "_groups")
	redisR, err := redis.String(r, nil)
	if err != nil {
		depts = users.GetUserDept(username)
		cache.PermCache.Put(username+"_groups", strings.Join(depts, ","), time.Second*600)
		return depts
	} else {
		return strings.Split(redisR, ",")
	}
	return depts
}

// 2018-02-18 21:13
// 检查用户权限
// 发布项目权限检查使用
func CheckUserPerms(username string, domain string) bool {

	r := cache.PermCache.Get(username)
	if checkPermExists(r, domain) {
		return true
	}

	depts := getUserDept(username)
	for _, groups := range depts {
		r = cache.PermCache.Get(groups)
		if checkPermExists(r, domain) {
			return true
		}
	}

	return false
}
