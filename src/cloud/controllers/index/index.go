package index

import (
	"github.com/astaxie/beego"
	"strings"
	"time"
	"cloud/sql"
	"cloud/models/index"
	"cloud/util"
	"cloud/controllers/base/cluster"
	"strconv"
	"cloud/models/app"
	_ "github.com/astaxie/beego/cache/redis"
	"github.com/garyburd/redigo/redis"
	"cloud/models/registry"
	"github.com/astaxie/beego/logs"
	"errors"
	"cloud/cache"
)

type IndexController struct {
	beego.Controller
}

// @Title create
// @Description create object
// @Param   body        body    models.Object   true        "The object content"
// @Success 302 {string} models.Object.Id
// @Failure 403 body is empty
// @router / [get]
func (this *IndexController) Get() {
	this.TplName = "login/login.tpl"
}

func (this *IndexController) LoginPage() {
	u := this.GetSession("username")
	if u != nil && u != "" {
		this.Redirect("/index", 302)
		return
	}
	this.TplName = "login/login.html"
}

// web 终端
// 2018-01-15 10:00
func (this *IndexController) WebTty() {
	searchMap := sql.GetSearchMap("ContainerId", *this.Ctx)
	data := app.CloudContainer{}
	q := sql.SearchSql(app.CloudContainer{}, app.SelectCloudContainer, searchMap)
	sql.Raw(q).QueryRow(&data)
	this.Data["username"] = this.GetSession("username")
	this.Data["namespace"] = util.Namespace(data.AppName, data.ResourceName)
	this.Data["pod"] = data.ContainerName
	this.Data["container"] = data.ServiceName
	this.Data["timestamp"] = strconv.FormatInt(time.Now().Unix(), 10)
	this.Data["cluster"] = data.ClusterName
	this.Data["time"] = util.GetDate()

	d := make([]string, 0)
	d = append(d, this.Data["username"].(string))
	d = append(d, this.Data["namespace"].(string))
	d = append(d, this.Data["pod"].(string))
	d = append(d, this.Data["container"].(string))
	d = append(d, this.Data["timestamp"].(string))
	d = append(d, this.Data["cluster"].(string))

	pass := beego.AppConfig.String("ttysecurity")
	token := util.Md5String(strings.Join(d, pass))
	this.Data["token"] = token

	this.TplName = "webtty/tty.html"
}

// @router /api/user/ [get]
func (this *IndexController) GetUser() {
	u := this.GetSession("username")
	if u != nil {
		this.Ctx.WriteString(u.(string))
	} else {
		this.Ctx.WriteString("")
	}
}

// 登录从数据库验证账号
// 2018-01-19 17;08
func DbAuth(user string, password string) bool {
	searchMap := sql.SearchMap{}
	searchMap.Put("UserName", user)
	searchMap.Put("Pwd", util.Md5String(password))
	data := index.DockerCloudAuthorityUser{}
	q := sql.SearchSql(data, index.SelectDockerCloudAuthorityUser, searchMap)
	sql.Raw(q).QueryRow(&data)
	if data.IsDel > 0 {
		return false
	}
	return true
}

var LockUserAuth util.Lock
// 2018-01-16 08:52
// 避免频繁更新数据库,加锁60秒后可操作
func writeLock(username string) bool {
	key := username
	if len(LockUserAuth.GetData()) > 0 {
		v, err := LockUserAuth.Get(key)
		if err {
			last := v.(int64)
			if time.Now().Unix()-last < 60 {
				return false
			}
		}
	}
	LockUserAuth.Put(key, time.Now().Unix())
	return true
}

// 2018-01-21 18:44
// 验证登陆仓库的管理员用户名和密码
func VerifyUser(user string, pass string, service string) bool {
	// 查询对象和用户的权限
	services := strings.Split(service, ".")
	if len(services) < 2 {
		services = append(services, "")
	}
	cacheStr := util.Md5String(pass + beego.AppConfig.String("ttyttysecurity"))
	key := user + "_admin" + "_" + service
	// 区别用户是管理员
	r := cache.RedisUserCache.Get(key)
	redisr, _ := redis.String(r, nil)
	logs.Error("获取到cache数据", redisr, cacheStr, redisr == cacheStr)
	if redisr == cacheStr {
		return true
	}
	d := make([]registry.CloudRegistryServer, 0)
	pass = util.Base64Encoding(pass)
	searchMap := sql.GetSearchMapV("Admin", user, "Password", pass, "Name", services[0], "ClusterName", services[1])
	q := sql.SearchSql(registry.CloudRegistryServer{}, registry.SelectCloudRegistryServer, searchMap)
	sql.Raw(q).QueryRows(&d)
	if len(d) < 1 {
		cache.RedisUserCache.Delete(key)
		return false
	}
	cache.RedisUserCache.Put(key, cacheStr, time.Minute*20)
	return true
}

// 登录后用户记录到数据库
// 2018-01-19 17:06
func RecordLoginUser(username string, password string) (bool, error) {
	cacheStr := util.Md5String(password + beego.AppConfig.String("ttyttysecurity"))
	r := cache.RedisUserCache.Get(username)
	redisr, _ := redis.String(r, nil)
	if redisr == cacheStr {
		return true, nil
	}
	r1, _ := util.LdapLoginAuth(username, password)
	if ! r1 {
		r1 = DbAuth(username, password)
		logs.Info("通过db验证用户", username, r)
	}
	if r1 {
		cache.RedisUserCache.Put(username, cacheStr, time.Minute*180)
		// 如果是ldap登录的,将数据记录到数据库里面
		v := index.DockerCloudAuthorityUser{}
		v.IsValid = 1
		v.Pwd = util.Md5String(password)
		v.UserName = username
		v.LastModifyTime = util.GetDate()
		v.Token = util.Md5String(username + util.GetDate())
		searchMap := sql.GetSearchMapV("UserName", username)
		user := index.DockerCloudAuthorityUser{}
		q := sql.SearchSql(v, index.SelectDockerCloudAuthorityUser, searchMap)
		sql.Raw(q).QueryRow(&user)
		if user.UserName == "" && v.UserName != ""{
			q = sql.InsertSql(v, index.InsertDockerCloudAuthorityUser)
			sql.Raw(q).Exec()
		} else {
			if ! writeLock(username) {
				q = sql.UpdateSql(v, index.UpdateDockerCloudAuthorityUser, searchMap, "Token")
				sql.Raw(q).Exec()
			}
		}
		return true, nil
	}
	return false, errors.New("ldap和数据库验证失败")
}

// 获取用户是否禁用
// 2018-01-22 09:18
func getUserIsDel(username string) bool {
	data := make([]index.DockerCloudAuthorityUser, 0)
	searchMap := sql.SearchMap{}
	searchMap.Put("IsDel", 1)
	searchMap.Put("UserName", username)
	q := sql.SearchSql(index.DockerCloudAuthorityUser{}, index.SelectDockerCloudAuthorityUser, searchMap)
	sql.Raw(q).QueryRow(&data)
	if len(data) > 0 {
		return true
	}
	return false
}

// @router /api/user/login [post]
func (this *IndexController) Login() {
	username := this.GetString("username")
	password := this.GetString("password")

	if getUserIsDel(util.GetUser(username)) {
		this.Ctx.WriteString("false,用户已经禁用")
		return
	}
	r, err := RecordLoginUser(util.GetUser(username), password)
	r = true
	ip := this.Ctx.Request.RemoteAddr
	o := sql.GetOrm()
	data := index.CloudLoginRecord{
		LoginStatus: 0,
		LoginTime:   util.GetDate(),
		LoginUser:   username,
		LoginIp:     ip}

	if !r {
		o.Raw(sql.InsertSql(data, index.InsertCloudLoginRecord)).Exec()
		this.Ctx.WriteString("false,验证失败" + err.Error())
		return
	}

	data.LoginStatus = 1
	o.Raw(sql.InsertSql(data, index.InsertCloudLoginRecord)).Exec()

	this.SetSession("username", username)
	this.SetSession("logintime", time.Now().Unix())
	this.SetSession("clientIp", ip)
	rd := util.GetReferer(*this.Ctx)
	if len(rd) > 0 && rd != "/login" {
		this.Ctx.WriteString("true," + rd)
		return
	} else {
		this.Ctx.WriteString("true,/index")
		return
	}
}

// 快捷入口页面
// @router /shortcut [get]
func (this *IndexController) Index() {
	// 获取全部集群数据
	this.TplName = "index/shortcut.html"
}

// @router /index [get]
func (this *IndexController) Shortcut() {
	// 获取全部集群数据
	this.TplName = "index/index.html"
}

// 首页获取集群详细数据
// @router /index/detail/hi(.*) [get]
func (this *IndexController) IndexDetail() {
	// 获取指定机器的数据
	this.Data["data"] = cluster.GetClusterData(this.Ctx.Input.Param(":hi"))
	this.TplName = "index/detail_use.html"
}

// @router /logout [get]
func (this *IndexController) OutLogin() {
	this.DelSession("username")
	this.DelSession("logintime")
	this.DelSession("clientIp")
	this.Redirect("/login", 302)
}
