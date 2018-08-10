package ci

import (
	"github.com/astaxie/beego"
	"cloud/sql"
	"cloud/models/ci"
	"cloud/util"
	"strings"
	"cloud/controllers/ent"
	"cloud/controllers/docker/application/app"
	"cloud/controllers/base/lb"
	"github.com/astaxie/beego/logs"
	"cloud/controllers/image"
	app2 "cloud/models/app"
	"cloud/k8s"
	lb2 "cloud/models/lb"
	"time"
)

// 2018-02-10 18:14
// 持续集成
type ServiceController struct {
	beego.Controller
}

// 服务发布管理入口页面
// @router /ci/service/list [get]
func (this *ServiceController) ServiceList() {
	this.TplName = "ci/service/list.html"
}

// 服务信息详情
// @router /ci/service/top/:id:int [get]
func (this *ServiceController) ServiceTop() {
	services, SvcCi := getCiServiceData(this)
	SvcCi = getImageInfo(SvcCi)
	SvcCi = getServiceAccess(services, SvcCi)
	this.Data["data"] = SvcCi
	this.TplName = "ci/service/top.html"
}

// 2018-02-16 11:00
// 切入流量页面
// @router /ci/service/flow/:id:int [get]
func (this *ServiceController) StartFlow() {
	version := this.GetString("version")
	_, SvcCi := getCiServiceData(this)
	this.Data["data"] = SvcCi
	this.Data["version"] = version
	this.TplName = "ci/service/start-flow.html"
}

// 服务发布历史入口页面
// @router /ci/service/release/history [get]
func (this *ServiceController) HistoryList() {
	this.Data["data"] = getCiService(this)
	this.TplName = "ci/service/history.html"
}

// 2018-02-17 11:14
// 服务发布日志入口页面
// @router /ci/service/release/logs [get]
func (this *ServiceController) ServiceLog() {
	this.Data["data"] = getCiService(this)
	this.TplName = "ci/service/logs.html"
}

// 2018-02-16 17:32
// 滚动更新页面
// @router /ci/service/rolling/:id:int [get]
func (this *ServiceController) RollingUpdate() {
	svcData := getCiService(this)
	history := getHistoryData(svcData)
	history.ServiceId = svcData.ServiceId
	if history.HistoryId == 0 {
		this.Ctx.WriteString("资源不可用,请发布后操作")
		return
	}
	this.Data["data"] = history
	this.TplName = "ci/service/rolling.html"
}

// 服务发布弹出页面
// @router /ci/service/release [get]
func (this *ServiceController) ServiceRelease() {

	// 历史页面修改发布信息
	historyId,histErr := this.GetInt("historyId")
	if histErr == nil && historyId != 0 {
		update := getHistory(sql.GetSearchMapV("HistoryId", this.GetString("historyId")))
		this.Data["data"] = update
		this.TplName = "ci/service/modify.html"
		return
	}

	searchMap := sql.GetSearchMapV("ServiceId", this.GetString("ServiceId"))
	serviceData := getServiceData(searchMap)
	b := ci.CloudCiReleaseHistory{}
	util.MergerStruct(serviceData, &b)
	services := app.GetServices(serviceData, "")
	if len(services) > 0 {
		imgs := strings.Split(services[0].ImageTag, ":")
		img := imgs[0:len(imgs)-1]
		b.ImageName = strings.Join(img, ":")
		images := registry.GetImageTag(services[0].ImageTag)
		logs.Info(images)
		this.Data["images"] = images
	}
	this.Data["data"] = b
	this.TplName = "ci/service/release.html"
}

// 服务发布管理添加页面
// @router /ci/service/add [get]
func (this *ServiceController) ServiceAdd() {
	id := this.GetString("ServiceId")
	update := ci.CloudCiService{}
	entData := ent.GetEntnameSelect()
	appData := app.GetAppSelect(sql.GetSearchMapV("CreateUser", getServiceUser(this)))
	var entHtml string
	var serviceHtml string
	var appHtml string
	// 更新操作
	if id != "0" {
		searchMap := sql.GetSearchMap("ServiceId", *this.Ctx)
		q := sql.SearchSql(ci.CloudCiService{}, ci.SelectCloudCiService, searchMap)
		sql.Raw(q).QueryRow(&update)

		entHtml = util.GetSelectOptionName(update.Entname)
		appHtml = util.GetSelectOptionName(update.AppName)

		searchMap = sql.GetSearchMapV(
			"Entname", update.Entname,
			"CreateUser", getServiceUser(this),
			"AppName", update.AppName,
			"ClusterName", update.ClusterName)
		serviceHtml = util.GetSelectOptionName(update.ServiceName) +
			app.GetServiceHtml(searchMap)

		this.Data["cluster"] = util.GetSelectOptionName(update.ClusterName)
		this.Data["apps"] = appHtml + appData
		this.Data["service"] = serviceHtml
	}

	this.Data["entname"] = entHtml + entData
	this.Data["data"] = update
	this.TplName = "ci/service/add.html"
}

// 获取服务发布数据
// 2018-02-17 11:20
// router /api/ci/service/logs [get]
func (this *ServiceController) ServiceLogs() {
	data := make([]ci.CloudCiReleaseLog, 0)
	domain := this.GetString("domain")
	searchMap := sql.SearchMap{}
	if domain != "" {
		searchMap.Put("Domain", domain)
	}
	key := this.GetString("search")
	searchSql := sql.SearchSql(
		ci.CloudCiReleaseLog{},
		ci.SelectCloudCiReleaseLog,
		searchMap)
	if domain == "" {
		searchSql += " where 1=1 "
	}

	if key != "" {
		key = sql.Replace(key)
		q := ci.SelectCiReleaseLogWhere
		searchSql += strings.Replace(q, "?", key, -1)
	}

	sql.OrderByPagingSql(searchSql, "log_id",
		*this.Ctx.Request, &data,
		ci.CloudCiReleaseLog{})

	var r = util.ResponseMap(data, sql.Count("cloud_ci_release_log", len(data), key),
		this.GetString("draw"))
	setServiceJson(this, r)
}

// string
// 服务发布保存
// @router /api/ci/service [post]
func (this *ServiceController) ServiceSave() {
	d := ci.CloudCiService{}
	err := this.ParseForm(&d)
	if err != nil {
		this.Ctx.WriteString("参数错误" + err.Error())
		return
	}
	util.SetPublicData(d, getServiceUser(this), &d)

	q := sql.InsertSql(d, ci.InsertCloudCiService)
	if d.ServiceId > 0 {
		searchMap := sql.SearchMap{}
		searchMap.Put("ServiceId", d.ServiceId)
		q = sql.UpdateSql(d,
			ci.UpdateCloudCiService,
			searchMap,
			"CreateTime,CreateUser")
	}
	_, err = sql.Raw(q).Exec()
	data, msg := util.SaveResponse(err, "名称已经被使用")
	util.SaveOperLog(getServiceUser(this),
		*this.Ctx,
		"保存服务发布配置 "+msg,
		d.Entname+d.AppName+d.ServiceName)
	setServiceJson(this, data)
}

// 获取服务发布数据
// 2018-02-10 19:45
// router /api/ci/service/name [get]
func (this *ServiceController) ServiceDataName() {
	// 服务发布数据
	data := make([]ci.CloudCiService, 0)
	q := sql.SearchSql(
		ci.CloudCiService{},
		ci.SelectCloudCiService,
		sql.SearchMap{})
	sql.Raw(q).QueryRows(&data)
	setServiceJson(this, data)
}

// 2018-02-14 16:40
// 服务发布历史数据
// @router /api/ci/service/history [get]
func (this *ServiceController) ReleaseHistory() {
	data := make([]ci.CloudCiReleaseHistory, 0)
	domain := this.GetString("domain")
	searchMap := sql.SearchMap{}
	if domain != "" {
		searchMap.Put("Domain", domain)
	}
	key := this.GetString("search")
	searchSql := sql.SearchSql(
		ci.CloudCiReleaseHistory{},
		ci.SelectCloudCiReleaseHistory,
		searchMap)
	if domain == "" {
		searchSql += " where 1=1 "
	}
	if key != "" {
		key = sql.Replace(key)
		q := ci.SelectReleaseHistoryWhere
		searchSql += strings.Replace(q, "?", key, -1)
	}

	sql.OrderByPagingSql(searchSql, "history_id",
		*this.Ctx.Request, &data,
		ci.CloudCiReleaseHistory{})

	var r = util.ResponseMap(data, sql.Count("cloud_ci_release_history", len(data), key),
		this.GetString("draw"))
	setServiceJson(this, r)
}

// 服务发布数据
// @router /api/ci/service [get]
func (this *ServiceController) ServiceDatas() {
	data := []ci.CloudCiService{}
	searchMap := sql.SearchMap{}
	key := this.GetString("search")
	//searchMap.Put("CreateUser", getServiceUser(this))
	searchSql := sql.SearchSql(
		ci.CloudCiService{},
		ci.SelectCloudCiService,
		searchMap)
	if key != "" {
		key = sql.Replace(key)
		q := ci.SelectCloudCiServiceWhere
		searchSql += strings.Replace(q, "?", key, -1)
	}

	sql.OrderByPagingSql(searchSql, "service_id",
		*this.Ctx.Request, &data,
		ci.CloudCiService{})

	domain := make([]string, 0)
	for _, v := range data {
		domain = append(domain, `"`+v.Domain+`"`)
	}

	domainMap := lb.GetLbServiceMap(domain)
	result := []ci.CloudCiService{}
	for _, v := range data {
		if ! CheckUserPerms(getServiceUser(this), v.Domain){
			continue
		}
		v = getImageInfo(v)
		v.LbVersion = domainMap.GetVString(v.Domain)
		result = append(result, v)
	}

	var r = util.ResponseMap(result, sql.Count("cloud_ci_service", len(data), key),
		this.GetString("draw"))
	setServiceJson(this, r)
}

// json
// 删除服务发布
// 2018-02-10 18:27
// @router /api/ci/service/:id:int [delete]
func (this *ServiceController) ServiceDelete() {
	searchMap := sql.GetSearchMap("ServiceId", *this.Ctx)
	serviceData := getServiceData(searchMap)

	q := sql.DeleteSql(ci.DeleteCloudCiService, searchMap)
	r, err := sql.Raw(q).Exec()
	data := util.DeleteResponse(err,
		*this.Ctx,
		"删除服务发布"+serviceData.ServiceName+serviceData.AppName+serviceData.Entname,
		getServiceUser(this),
		serviceData.CreateUser, r)
	setServiceJson(this, data)
	saveServiceLog(this, serviceData, "删除了服务发布")
}

// 2018-02-16 18:57
// 滚动镜像
// @param version
// @router /api/ci/service/rolling/:id:int [post]
func (this *ServiceController) RollingUpdateExec() {
	d := k8s.RollingParam{}
	err := this.ParseForm(&d)
	version := this.GetString("version")
	serviceData, history, cl, serviceCiData, msg := getRollingData(err, this)

	if ! CheckUserPerms(getServiceUser(this), serviceCiData.Domain){
		setServiceJson(this, util.ApiResponse(false, "用户无权限操作"))
		return
	}

	if msg != "" {
		setServiceJson(this, util.ApiResponse(false, msg))
		return
	}

	if serviceData.ImageTag == history.ImageName {
		msg = "滚动更新失败,新旧版本一致"
		logs.Error(msg)
		setServiceJson(this, util.ApiResponse(false, msg))
		return
	}

	logs.Info("更新镜像", history.ImageName)
	d.Client = cl
	d.Images = history.ImageName
	d.Namespace = util.Namespace(serviceData.AppName, serviceData.ResourceName)
	d.Name = util.Namespace(serviceCiData.ServiceName, version)

	err = execCiUpdate(serviceData, d, getServiceUser(this))
	if err != nil {
		setServiceJson(this, util.ApiResponse(false, "更新失败"+err.Error()))
		return
	}
	saveOperLog(this, "滚动更新服务 ", serviceCiData, serviceData)
	saveServiceLog(this, serviceCiData, "滚动更新服务"+d.Images)
}

// 2018-02-15 10:49
// 将篮板更新镜像
// @router /api/ci/service/blue/:id:int [post]
func (this *ServiceController) UpdateBlueService() {
	// 第一次发布为2，就是绿的版本
	// 第二次发布为1，就是蓝的版本,老版本
	services, serviceCiData := getCiServiceData(this)
	history := getHistoryData(serviceCiData)

	if ! CheckUserPerms(getServiceUser(this), serviceCiData.Domain){
		setServiceJson(this, util.ApiResponse(false, "用户无权限操作"))
		return
	}

	var serviceInfo app2.CloudAppService
	if len(services) < 2 {
		setServiceJson(this, util.ApiResponse(false, "更新失败,只有一个版本,不能升级"))
		return
	}
	serviceInfo = getImageServiceInfo(services, "1")
	serviceInfo.ServiceVersion = "1"
	serviceInfo.ImageTag = history.ImageName
	err := app.ExecUpdate(serviceInfo, "image", getServiceUser(this))
	if err != nil {
		setServiceJson(this, util.ApiResponse(false, "更新失败"+err.Error()))
		return
	}
	saveOperLog(this, "更新蓝版 ", serviceCiData, serviceInfo)
	saveServiceLog(this, serviceCiData, "蓝绿发布,更新蓝版服务")
}

// 服务发布入口
// 2018-02-10 18:27
// @param ServiceVersion  [ 1 | 2 ]
// @router /api/ci/service/release/:id:int [post]
func (this *ServiceController) ServiceReleaseExec() {
	services, serviceCiData := getCiServiceData(this)
	d := ci.CloudCiReleaseHistory{}

	if ! CheckUserPerms(getServiceUser(this), serviceCiData.Domain){
		setServiceJson(this, util.ApiResponse(false, "用户无权限操作"))
		return
	}

	err := this.ParseForm(&d)
	if err != nil || d.ImageName == "" || d.ReleaseTestUser == "" || d.Description == "" || len(services) == 0 {
		setServiceJson(this, util.ApiResponse(false, "参数错误"+err.Error()))
		return
	}

	if len(services) > 1 {
		setServiceJson(this, util.ApiResponse(false, "不能重复发布"))
		return
	}

	if services[0].ImageTag == d.ImageName {
		setServiceJson(this, util.ApiResponse(false, "蓝绿版本一致,请注意镜像选择!"))
		return
	}

	// 保存旧版本数据
	d.OldImages = services[0].ImageTag

	serviceCiData.ImageName = d.ImageName
	var user = getServiceUser(this)

	msg, status, serviceInfo := app.CreateGreenService(serviceCiData, user)
	if ! status {
		logs.Error("创建绿版错误", msg)
		setServiceJson(this, util.ApiResponse(false, msg))
		return
	}

	util.MergerStruct(serviceCiData, &d)
	setServiceJson(this, util.ApiResponse(status, "创建绿版成功"))

	if d.AutoSwitch == "on" && status {
		updateLbVersion(serviceInfo.ServiceVersion, serviceCiData, serviceInfo, true)
	}

	// 记录历史
	if status {
		d.ReleaseVersion = "2"
		d.Action = "update"
		saveHistory(this, d)
		updateLbPercent(0, serviceCiData, ci.UpdateCiServicePercent, serviceInfo.ServiceVersion)
		updateLbPercent(0, serviceCiData, lb2.UpdateLbServicePercent, serviceInfo.ServiceVersion)
		saveServiceLog(this, serviceCiData, "发布服务,创建了绿色版本")
	}
}

// 2018-02-14 21:10
// 上线服务,将新版本替换到线上
// 蓝绿发布切换
// @param ServerVersion
// @router /api/ci/service/online/:id:int [post]
func (this *ServiceController) ServiceOnline() {
	version := this.GetString("ServiceVersion")
	services, serviceCiData := getCiServiceData(this)

	if ! CheckUserPerms(getServiceUser(this), serviceCiData.Domain){
		setServiceJson(this, util.ApiResponse(false, "用户无权限操作"))
		return
	}

	var msg string
	if len(services) < 2 {
		msg = "没有获取到Service数据,程序退出"
		setServiceJson(this, util.ApiResponse(false, msg))
		return
	}

	if version == "" {
		msg = "缺少版本号参数"
	}

	if serviceCiData.Percent > 0 {
		msg = "所有蓝版剔除负载均衡,才可以做蓝绿切换"
	}

	domainMap := getLbDomain(serviceCiData.Domain)
	if _, ok := domainMap.Get(serviceCiData.Domain); ! ok {
		msg = "负载均衡没有该数据,程序退出"
	}

	if version == domainMap.GetVString(serviceCiData.Domain) {
		msg = "当前负载均衡已经使用该版本提供服务,无需再次切换"
	}

	if msg != "" {
		logs.Error(msg, serviceCiData.AppName, serviceCiData.ServiceName, serviceCiData.Domain)
		setServiceJson(this, util.ApiResponse(false, msg))
		return
	}
	// 只做lb更新数据参考
	serviceData := services[1]

	updateLbVersion(version, serviceCiData, serviceData, true)
	saveOperLog(this, "上线服务 ", serviceCiData, serviceData)
	saveServiceLog(this, serviceCiData, "蓝绿版本切换,切入版本为"+version)
}

// 2018-02-14 18:20
// 下线服务,删除不要的服务
// @router /api/ci/service/release/:id:int [delete]
func (this *ServiceController) ServiceOffline() {
	version := this.GetString("version")
	services, serviceCiData := getCiServiceData(this)

	if ! CheckUserPerms(getServiceUser(this), serviceCiData.Domain){
		setServiceJson(this, util.ApiResponse(false, "用户无权限操作"))
		return
	}

	domainMap := getLbDomain(serviceCiData.Domain)
	if _, ok := domainMap.Get(serviceCiData.Domain); ! ok {
		msg := "负载均衡没有该数据,程序退出"
		logs.Error(msg)
		setServiceJson(this, util.ApiResponse(false, msg))
		return
	}

	if domainMap.GetVString(serviceCiData.Domain) == version {
		msg := "该版本正在提供服务,无法删除"
		logs.Error(msg)
		setServiceJson(this, util.ApiResponse(false, msg))
		return
	}

	serviceData := getImageServiceInfo(services, version)
	var msg = checkDeleteServiceParam(serviceCiData, services, serviceData)
	if msg != "" {
		logs.Error(msg, serviceCiData.AppName, serviceCiData.ServiceName, serviceCiData.Domain)
		setServiceJson(this, util.ApiResponse(false, msg))
		return
	}
	err := app.DeleteK8sService(serviceData, "")
	if err != nil {
		setServiceJson(this, util.ApiResponse(false, err))
		return
	}
	saveOperLog(this, "删除服务 ", serviceCiData, serviceData)
	saveServiceLog(this, serviceCiData, "删除服务")
}

// 2018-02-15 06:40
// 将所有配置都回归到蓝版本, svc和lb
// 服务回滚
// @router /api/ci/service/rollback/:id:int [post]
func (this *ServiceController) ServiceRollback() {
	var history ci.CloudCiReleaseHistory
	historyId := this.GetString("history")
	var services []app2.CloudAppService
	var serviceCiData ci.CloudCiService

	// 在服务发布页面回滚
	if historyId == "" {
		services, serviceCiData = getCiServiceData(this)
		history = getHistoryData(serviceCiData)
	} else {
		// 在历史页面回滚
		history = getHistory(sql.GetSearchMapV("HistoryId", historyId))
		s := util.TimeToStamp(history.CreateTime)
		e := time.Now().Unix()
		// 超过俩周的不会滚
		if e-s > 1209600 || history.HistoryId == 0 {
			setServiceJson(this, util.ApiResponse(false, "该数据已经超过2周,不能回滚"))
			return
		}
		serviceCiData = getServiceData(sql.GetSearchMapV("Domain", history.Domain))
		services = app.GetServices(serviceCiData, "")
	}

	if ! CheckUserPerms(getServiceUser(this), serviceCiData.Domain){
		setServiceJson(this, util.ApiResponse(false, "用户无权限操作"))
		return
	}

	if len(services) < 2 {
		setServiceJson(this, util.ApiResponse(false, "目前只有一个服务提供,无法回滚"))
		return
	}

	lbData := lb.GetLbDomainData(serviceCiData.Domain)
	serviceCiData = getImageInfo(serviceCiData)
	var version = lbData.ServiceVersion

	serviceData := getImageServiceInfo(services, version)


	msg := execUpdateRollbackService(serviceCiData, history, serviceData, this, version)
	if msg != "" {
		logs.Error(msg)
		setServiceJson(this, util.ApiResponse(false, msg))
	}
	setServiceJson(this, util.ApiResponse(false, "正在操作中,请稍后重试验证!"))
}

// 2018-02-16 12:57
// 流量切入,按百分比切入
// 主要是修改lb的upstream配置
// @param Percent
// @router /api/ci/service/flow/:id:int [post]
func (this *ServiceController) StartFlowExec() {
	percent, err := this.GetInt("percent")
	if err != nil {
		setServiceJson(this, util.ApiResponse(false, "percent参数错误"+err.Error()))
		return
	}

	services, serviceCiData := getCiServiceData(this)

	if ! CheckUserPerms(getServiceUser(this), serviceCiData.Domain){
		setServiceJson(this, util.ApiResponse(false, "用户无权限操作"))
		return
	}

	if percent > 0 && percent == serviceCiData.Percent {
		setServiceJson(this, util.ApiResponse(false, "已经切入流量,无需重新切入"))
		return
	}

	lbData := lb.GetLbDomainData(serviceCiData.Domain)
	if lbData.ServiceId == 0 {
		setServiceJson(this, util.ApiResponse(false, "负载服务不存在"))
		return
	}
	version := lbData.ServiceVersion
	updateLbPercent(percent, serviceCiData, lb2.UpdateLbServicePercent, version)

	serviceData := getImageServiceInfo(services, version)
	err = updateLbVersion(version, serviceCiData, serviceData, false)
	if err == nil {
		updateLbPercent(percent, serviceCiData, ci.UpdateCiServicePercent, version)
		setServiceJson(this, util.ApiResponse(true, "操作成功,服务正在后台更新,等待30秒后生效"))
		saveServiceLog(this, serviceCiData, "金丝雀切入流量,切入百分比"+this.GetString("percent"))
	} else {
		setServiceJson(this, util.ApiResponse(true, "操作失败"+err.Error()))
		saveServiceLog(this, serviceCiData, "金丝雀切入流量,切入百分比"+this.GetString("percent")+"操作失败"+err.Error())
	}
}

// 2018-02-18 17:30
// 更新发布历史中的描述信息
// @router /api/ci/service/history/:id:int [post]
func (this *ServiceController) UpdateHistory()  {
	d := ci.CloudCiReleaseHistory{}
	this.ParseForm(&d)
	searchMap := sql.SearchMap{}
	searchMap.Put("CreateUser", getServiceUser(this))
	searchMap.Put("HistoryId", d.HistoryId)
	if d.HistoryId == 0 {
		setServiceJson(this, util.ApiResponse(false, "参数错误"))
		return
	}
	q := ci.UpdateCloudCiReleaseHistory
	q = sql.UpdateSql(d, q, searchMap, ci.UpdateHistoryExclude)
	sql.Raw(q).Exec()
	setServiceJson(this, util.ApiResponse(false, "更新成功"))
}