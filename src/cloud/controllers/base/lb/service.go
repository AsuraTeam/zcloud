package lb

import (
	"cloud/sql"
	"cloud/util"
	"github.com/astaxie/beego"
	"cloud/models/lb"
	"cloud/k8s"
	"cloud/controllers/docker/application/app"
	"strconv"
	"database/sql/driver"
	"golang.org/x/crypto/openpgp/errors"
	"strings"
	"github.com/astaxie/beego/logs"
)

type ServiceController struct {
	beego.Controller
}

// 负载均衡管理添加页面
// @router /base/network/lb/service/add [get]
func (this *ServiceController) ServiceAdd() {
	id := this.GetString("ServiceId")
	lbData := GetLbData(this.GetString("LbId"))
	if id == "" || lbData.LbId == 0{
		this.Ctx.WriteString("参数错误")
		return
	}
	update := k8s.CloudLbService{}
	update.ServiceVersion = "1"
	this.Data["LbMethod1"] = "checked"

	certs := GetCertSelect()

	// 更新操作
	var serviceData string
	var appData string
	certData := util.GetSelectOption("", "", "")
	util.MergerStruct(lbData, &update)
	if id != "0" {
		searchMap := sql.SearchMap{}
		searchMap.Put("ServiceId", id)
		q := sql.SearchSql(k8s.CloudLbService{}, lb.SelectCloudLbService, searchMap)
		sql.Raw(q).QueryRow(&update)
		
		if update.DefaultDomain == "on" {
			this.Data["DefaultDomain"] = "checked"
		}
		if update.LbMethod == "pod" {
			this.Data["LbMethod2"] = "checked"
			this.Data["LbMethod1"] = ""
		}
		
		this.Data["readonly"] = "readonly"
		lbData = GetLbData(update.LbId)
		serviceData = util.GetSelectOption(update.ServiceName, update.LbServiceId, update.ServiceName)
		certData = util.GetSelectOptionName(update.CertFile)
		appData = util.GetSelectOptionName(update.AppName)
	}

	if lbData.LbId == 0 {
		this.TplName = "base/network/lb/list.html"
		return
	}

	user := getServiceUser(this)
	data := app.GetUserLbService(user, update.ClusterName, "")

	for _, v := range data {
		serviceData += util.GetSelectOption(v.ServiceName,
			strconv.FormatInt(v.ServiceId, 10),
			v.ServiceName)
	}

	searchMap := sql.GetSearchMapV("CreateUser", getServiceUser(this),
	"ClusterName", lbData.ClusterName,
	"Entname", update.Entname)
	this.Data["apps"] = appData + app.GetAppSelect(searchMap)
	update.ClusterName = lbData.ClusterName
	this.Data["cert"] = certData + certs
	this.Data["data"] = update
	this.Data["ServiceData"] = serviceData
	this.TplName = "base/network/lb/service.html"
}

// string
// 负载均衡保存
// @router /api/lb/service [post]
func (this *ServiceController) ServiceSave() {
	d := k8s.CloudLbService{}
	err := this.ParseForm(&d)
	if err != nil {
		this.Ctx.WriteString("参数错误" + err.Error())
		return
	}

	user := getServiceUser(this)
	util.SetPublicData(d, getServiceUser(this), &d)

	serviceData := app.GetUserLbService(user, d.ClusterName, d.LbServiceId)
	if len(serviceData) > 0 {
		service := serviceData[0]
		d.ServiceName,d.AppName, d.ClusterName, d.ResourceName = service.ServiceName, service.AppName,service.ClusterName,service.ResourceName
	}

	q := sql.InsertSql(d, lb.InsertCloudLbService)
	if d.ServiceId > 0 {
		searchMap := sql.SearchMap{}
		searchMap.Put("ServiceId", d.ServiceId)
		searchMap.Put("CreateUser", user)
		q = sql.UpdateSql(d, lb.UpdateCloudLbService, searchMap, lb.UpdateLbServiceExclude)
	}
	_, err = sql.Raw(q).Exec()
	data, msg := util.SaveResponse(nil, "名称已经被使用")
	util.SaveOperLog(
		this.GetSession("username"), 
		*this.Ctx, "保存负载均衡服务配置 "+msg, 
		d.LbName)
	setServiceJson(this, data)
	app.UpdateServiceDomain()
	go k8s.CreateNginxConf("")
	go k8s.CreateNginxConf("-test")
}

// 负载均衡数据
// @router /api/lb/service/:hi(.*) [get]
func (this *ServiceController) ServiceData() {
	data := make([]k8s.CloudLbService, 0)
	searchMap := sql.SearchMap{}
	lbName := this.GetString("LbName")
	if lbName == "" {
		lbName = this.Ctx.Input.Param(":hi")
	}
	key := this.GetString("key")
	domain := this.GetString("domain")
	if lbName == "" {
		setServiceJson(this, util.ResponseMapError("缺少lb名称"))
		return
	}

	//user := getServiceUser(this)
	//searchMap.Put("CreateUser", user)
	if domain != "" {
		searchMap.Put("Domain", domain)
	}

	searchMap.Put("LbName", lbName)
	searchSql := sql.SearchSql(k8s.CloudLbService{},
		lb.SelectCloudLbService,
		searchMap)

	if key != "" {
		key = sql.Replace(key)
		searchSql += strings.Replace(lb.SelectCloudLbServiceWhere, "?", key, -1)
	}

	sql.Raw(searchSql).QueryRows(&data)
	r := util.ResponseMap(data,
		sql.CountSearchMap("cloud_lb_service", sql.SearchMap{}, len(data), key),
		this.GetString("draw"))
	setServiceJson(this, r)
}

// json
// 删除负载均衡
// @router /api/network/lb/service/:id:int [delete]
func (this *ServiceController) ServiceDelete() {
	searchMap := sql.GetSearchMap("ServiceId", *this.Ctx)
	template := k8s.CloudLbService{}
	q := sql.SearchSql(template, lb.SelectCloudLbService, searchMap)
	sql.Raw(q).QueryRow(&template)

	q = sql.DeleteSql(lb.DeleteCloudLbService, searchMap)
	r, err := sql.Raw(q).Exec()
	data := util.DeleteResponse(err,
		*this.Ctx,
		"删除负载均衡"+template.LbName,
		this.GetSession("username"),
		template.CreateUser,
		r)

	q = sql.DeleteSql(lb.DeleteCloudLbNginxConf, searchMap)
	sql.Raw(q).Exec()
	setServiceJson(this, data)

	go k8s.CreateNginxConf("")
	go k8s.CreateNginxConf("-test")

}

func setServiceJson(this *ServiceController, data interface{}) {
	this.Data["json"] = data
	this.ServeJSON(false)
}

// 获取nginx配置文件信息
// 2018-02-03 07:31
func getNginxConf(this *ServiceController) k8s.CloudLbNginxConf {
	searchMap := sql.GetSearchMap("ServiceId", *this.Ctx)
	searchMap.Put("CreateUser", getServiceUser(this))
	conf := k8s.CloudLbNginxConf{}
	q := sql.SearchSql(conf, k8s.SelectCloudLbNginxConf, searchMap)
	sql.Raw(q).QueryRow(&conf)
	return conf
}


// 2018-02-01 21:39
// 获取nginx配置文件信息
// @router /api/network/lb/nginx/:id:int [get]
func (this *ServiceController) GetNginxConf() {
	conf := getNginxConf(this)
	this.Ctx.WriteString(conf.Vhost)
}

// 获取负载均衡的域名,通过环境区分
// @param entname
// @router /api/network/lb/domain [get]
func (this *ServiceController) GetLbDomain() {
	data := GetDomainSelect(this.GetString("entname"))
	setServiceJson(this, data)
}

// 2018-02-01 22:11
// 保存nginx配置
// @router /api/network/lb/nginx/:id:int [post]
func (this *ServiceController) SaveNginxConf() {
	searchMap := sql.GetSearchMap("ServiceId", *this.Ctx)
	if searchMap.Get("ServiceId") == "" {
		data := util.DeleteResponse(errors.UnsupportedError("配置异常"),
			*this.Ctx,
			"保存nginx配置",
			this.GetSession("username"), "",
			driver.ResultNoRows)
		setServiceJson(this, data)
		return
	}

	user := getServiceUser(this)
	searchMap.Put("CreateUser", user)
	vhost := this.GetString("Vhost")
	if len(vhost) < 50 {
		data := util.DeleteResponse(errors.UnsupportedError("配置异常"),
			*this.Ctx, "保存nginx配置",
			this.GetSession("username"),
			vhost,
			driver.ResultNoRows)
		setServiceJson(this, data)
		return
	}
	conf := getNginxConf(this)

	master, port := k8s.GetMasterIp(conf.ClusterName)
	configData := map[string]interface{}{
		conf.Domain + ".conf": vhost,
	}

	sslDbData := make(map[string]interface{})
	if conf.CertFile != "" {
		sslDbData = k8s.GetCertConfigData(conf.CertFile, sslDbData)
	}

	logs.Info("获取到检查配置", util.ObjToString(configData))
	k8s.MakeTestNginxConfMap(configData, sslDbData, conf.ClusterName)

	logStr, logTime := k8s.MakeTestJob(master, port, conf.ClusterName)
	if ! strings.Contains(logStr, "test is successful") {
		data, _ := util.SaveResponse(errors.InvalidArgumentError("配置检查失败"), logStr)
		setServiceJson(this, data)
		logs.Error("检查nginx配置失败", user, logStr, logTime)
		return
	}

	q := `update cloud_lb_nginx_conf set vhost="` + sql.Replace(vhost) + `" where create_user="` + user + `" and service_id=` + this.Ctx.Input.Param(":id")
	sql.Raw(q).Exec()
	r := util.ApiResponse(true, "保存成功 " + logStr)
	setServiceJson(this, r)
	go k8s.CreateNginxConf("")
}



// 2018-02-14 14:40
// 获取域名选项卡
func GetDomainSelect(entname string) string {
	data := make([]lb.LbServiceVersion, 0)

	q := sql.SearchSql(
		lb.LbServiceVersion{},
		lb.SelectLbDomain,
		sql.GetSearchMapV("Entname", entname))
	sql.Raw(q).QueryRows(&data)

	var html string
	for _,v := range data{
		html += util.GetSelectOptionName(v.Domain)
	}
	return html
}

// 2018-02-17 21:20
// 按域名查找lb服务信息
func GetLbDomainData(domain string) k8s.CloudLbService {
	data := k8s.CloudLbService{}
	q := strings.Replace(lb.SelectLbDomainData, "{0}", domain, -1)
	sql.Raw(q).QueryRow(&data)
	return data
}

// 2018-02-17 21:17
// 获取指定域名的服务
func GetLbServiceData(domains []string) []lb.LbServiceVersion {
	data := make([]lb.LbServiceVersion, 0)
	q := sql.SearchSql(
		lb.LbServiceVersion{},
		lb.SelectLbServiceVersion,
		sql.SearchMap{})
	q = strings.Replace(q, "?", strings.Join(domains, ","), -1)
	sql.Raw(q).QueryRows(&data)
	return data
}

// 2018-02-14 13:29
// 获取指定域名在负载均衡的版本号
func GetLbServiceMap(domains []string) util.Lock {
	lock := util.Lock{}
	data := GetLbServiceData(domains)
	for _, v := range data {
		lock.Put(v.Domain, v.ServiceVersion)
	}
	return lock
}

func getServiceUser(this *ServiceController) string {
	return util.GetUser(this.GetSession("username"))
}