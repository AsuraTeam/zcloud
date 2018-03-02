package ci

import (
	"cloud/util"
	"cloud/sql"
	"cloud/models/ci"
	"cloud/controllers/base/lb"
	"cloud/models/app"
	"cloud/k8s"
	"strings"
	lb2 "cloud/models/lb"
	app2 "cloud/controllers/docker/application/app"
	"strconv"
	"k8s.io/client-go/kubernetes"
	"github.com/astaxie/beego/logs"
)

// 2018-02-10 18:29
func setServiceJson(this *ServiceController, data interface{}) {
	this.Data["json"] = data
	this.ServeJSON(false)
}

// 2018-02-14 07:08
func getServiceData(searchMap sql.SearchMap) ci.CloudCiService {
	serviceData := ci.CloudCiService{}
	q := sql.SearchSql(serviceData, ci.SelectCloudCiService, searchMap)
	sql.Raw(q).QueryRow(&serviceData)
	return serviceData
}

// 2018-02-18 20;26
// 获取域名选择卡
func getServiceDomainSelect(domain string) string {
	domains := strings.Split(domain, ",")
	var html string
	data := []ci.CloudCiService{}
	sql.Raw(ci.SelectCloudCiService).QueryRows(&data)
	for _,v := range data {
		if ! util.ListExistsString(domains, v.Domain) {
			html += util.GetSelectOptionName(v.Domain)
		}
	}
	return html
}

// 获取用户
func getServiceUser(this *ServiceController) string {
	return util.GetUser(this.GetSession("username"))
}

// 2018-02-17 16:47
// 获取负载均衡版本数据
func getLbDomain(domain string) util.Lock {
	domainMap := lb.GetLbServiceMap([]string{`"` + domain + `"`})
	return domainMap
}

// 2018-02-14 21:15
// 检查参数正确性
func checkDeleteServiceParam(serviceCiData ci.CloudCiService, services []app.CloudAppService, service app.CloudAppService) string {
	var msg string
	if len(services) < 2 {
		msg = "没有获取到Service数据,程序退出"
		return msg
	}

	domainMap := getLbDomain(serviceCiData.Domain)
	if _, ok := domainMap.Get(serviceCiData.Domain); ! ok {
		msg = "负载均衡没有该数据,程序退出"
	}

	if domainMap.GetVString(serviceCiData.Domain) == service.ServiceVersion {
		msg = "该版本正在提供服务,不能操作"
	}

	if serviceCiData.Percent > 0 {
		msg = "所有金丝雀加入的服务都下线后,可以删除服务"
	}

	cl, err := k8s.GetClient(service.ClusterName)
	if err != nil {
		msg = "获取k8s连接失败" + err.Error()
	} else {
		svcName := util.Namespace(service.ServiceName, service.ServiceVersion)
		namespace := util.Namespace(service.AppName, service.ResourceName)
		svc := k8s.GetAppService(cl, namespace, svcName)
		if _, ok := svc.Spec.Selector["name"]; ok {
			if svcName != svc.Spec.Selector["name"] {
				msg = "版本交叉出现,不能删除,需要统一到蓝版后再进行删除"
			}
		}
	}

	return msg
}

//// 2018-02-15 06:26
//// 所有服务一直指向蓝版本
//// 新域名指向绿版本
//func updateCiReleaseSvc(serviceData app.CloudAppService, version string) error {
//	namespace := util.Namespace(serviceData.AppName, serviceData.ResourceName)
//	cl, _ := k8s.GetClient(k8s.GetMasterIp(serviceData.ClusterName))
//	name := util.Namespace(serviceData.ServiceName, serviceData.ServiceVersion)
//	svc := k8s.GetAppService(cl, namespace, name)
//
//	if svc.Spec.Selector["name"] == name {
//		msg := "版本无需回滚,蓝版没有更新"
//		logs.Error(msg)
//		return errors.InvalidArgumentError(msg)
//	}
//
//	logs.Info("获取到svc", util.ObjToString(svc))
//	selector := map[string]string{
//		"name": util.Namespace(serviceData.ServiceName, version),
//	}
//	svc.Spec.Selector = selector
//	d, err := cl.CoreV1().Services(namespace).Update(svc)
//	if err != nil {
//		logs.Error("更新服务失败", d, err)
//	} else {
//		logs.Info("更新服务", d)
//	}
//	return err
//}

func getCiService(this *ServiceController) ci.CloudCiService {
	searchMap := sql.GetSearchMap("ServiceId", *this.Ctx)
	serviceCiData := getServiceData(searchMap)
	return serviceCiData
}

// 2018-02-15 06:36
// 获取公共数据,服务信息和ci发布信息
func getCiServiceData(this *ServiceController) ([]app.CloudAppService, ci.CloudCiService) {
	serviceCiData := getCiService(this)
	services := app2.GetServices(serviceCiData, "")
	return services, serviceCiData
}

// 2018-02-15 06:28
// 在lb服务器将版本号更新
func updateLbVersion(version string, serviceCiData ci.CloudCiService, serviceData app.CloudAppService, updateLb bool) error{
	// 更新数据库版本号
	if updateLb {
		q := lb2.UpdateLbServiceServiceVersion
		q = strings.Replace(q, "{0}", version, -1)
		q = strings.Replace(q, "{1}", serviceCiData.Domain, -1)
		sql.Raw(q).Exec()
	}

	param := k8s.UpdateLbNginxUpstream{}
	param.Domain = serviceCiData.Domain
	param.ServiceName = util.Namespace(serviceCiData.ServiceName, version)
	v := lb.GetLbDomainData(serviceCiData.Domain)
	v.AppName = serviceData.AppName
	v.ResourceName = serviceData.ResourceName
	param.V = v
	param.ClusterName = serviceCiData.ClusterName
	param.Namespace = util.Namespace(serviceData.AppName, serviceData.ResourceName)
	err := k8s.UpdateNginxLbUpstream(param)
	return err
}

// 2018-02-15 09:03
// 记录日志
func saveOperLog(this *ServiceController, info string, ciData ci.CloudCiService, serviceData app.CloudAppService) {
	util.SaveOperLog(getServiceUser(this), *this.Ctx, info+ciData.Domain+serviceData.ServiceVersion, serviceData.ClusterName)
	setServiceJson(this, util.ApiResponse(true, info+"成功"))
}

// 2018-02-15 09:29
// 保存发布历史
func saveHistory(this *ServiceController, d ci.CloudCiReleaseHistory) {
	d.CreateTime = util.GetDate()
	d.CreateUser = getServiceUser(this)
	q := sql.InsertSql(d, ci.InsertCloudCiReleaseHistory)
	sql.Raw(q).Exec()
}

// 2018-02-15 10:39
// 获取服务信息
func getImageServiceInfo(services []app.CloudAppService, version string) app.CloudAppService {
	var serviceInfo app.CloudAppService
	if version != "" {
		if services[0].ServiceVersion == version {
			return services[0]
		}
		if len(services) == 2 && services[1].ServiceVersion == version {
			return services[1]
		}
	}
	if services[0].ServiceVersion == "1" {
		serviceInfo = services[0]
	} else {
		if len(services) > 1 {
			serviceInfo.ImageTag = services[1].ImageTag
		}
	}
	if len(services) > 1 && services[1].ServiceVersion == "1" {
		serviceInfo = services[1]
	} else {
		serviceInfo.ImageTag = services[0].ImageTag
	}
	return serviceInfo
}

// 2018-02-18 15:47
// 获取历史数据
func getHistory(searchMap sql.SearchMap) ci.CloudCiReleaseHistory {
	data := ci.CloudCiReleaseHistory{}
	q := ci.SelectCloudCiReleaseHistory
	q = sql.SearchSql(ci.CloudCiReleaseHistory{}, q, searchMap)
	q = sql.SearchOrder(q, "history_id")
	q = q + " limit 1"
	sql.Raw(q).QueryRow(&data)
	return data
}

// 2018-02-15 10:45
// 获取历史数据
func getHistoryData(serviceCiData ci.CloudCiService) ci.CloudCiReleaseHistory {
	searchMap := sql.GetSearchMapV("ServiceName", serviceCiData.ServiceName)
	searchMap.Put("Domain", serviceCiData.Domain)
	searchMap.Put("Entname", serviceCiData.Entname)
	return getHistory(searchMap)
}

// 2018-02-15 22:34
// 获取镜像名称
func getImageName(img1 string) string {
	img := strings.Split(img1, ":")
	return strings.Join(img[0:len(img)-1], ":")
}

// 2018-02-15 15:56
// 获取蓝绿版本的镜像信息
func getImageInfo(v ci.CloudCiService) ci.CloudCiService {
	services := app2.GetServices(v, app.SelectCurrentVersion)
	lock := app2.GetCurrentVersion(v, services)
	img1 := lock.GetVString("1")
	v.ImageInfoBlue = "<span style='color: #4489e4a6' class='th-top-8 text-default m-r-10'>蓝:  " + getImageTag(img1) + "</span>"
	if _, ok := lock.Get("2"); ok {
		v.ImageInfoGreen = "<span style='color: #33b867' class='text-default m-r-10'>绿:  " + getImageTag(lock.GetVString("2")) + "</span>"
	}
	v.ImageName = getImageName(img1)
	if img1 == "" {
		v.ImageName = getImageName(lock.GetVString("2"))
	}
	// 获取那个服务是比较新的
	if len(services) == 2 {
		if services[0].ServiceId > services[1].ServiceId {
			v.NewVersion = services[0].ServiceVersion
		}else{
			v.NewVersion = services[1].ServiceVersion
		}
	}
	return v
}

// 2018-02-15 16:05
// 获取镜像tag
func getImageTag(image string) string {
	imgs := strings.Split(image, ":")
	if len(imgs) > 1 {
		return imgs[len(imgs)-1]
	}
	return "版本不存在"
}

// 2018-02-15 22:28
// 获取服务访问信息
func getServiceAccess(services []app.CloudAppService, SvcCi ci.CloudCiService) ci.CloudCiService {
	if len(services) > 0 {
		SvcCi.ImageName = getImageName(services[0].ImageTag)
		c, _ := k8s.GetClient(SvcCi.ClusterName)
		namespace := util.Namespace(SvcCi.AppName, services[0].ResourceName)
		n1 := namespace + SvcCi.ServiceName

		blueNamespace := util.Namespace(SvcCi.ServiceName, "1")
		greenNamespace := util.Namespace(SvcCi.ServiceName, "2")
		blue := k8s.GetDeploymentApp(c, namespace, blueNamespace)
		green := k8s.GetDeploymentApp(c, namespace, greenNamespace)

		// 获取svc指向的pod
		svc := k8s.GetAppService(c, namespace, blueNamespace)
		if _, ok := svc.Spec.Selector["name"]; ok {
			selector := svc.Spec.Selector["name"]
			SvcCi.BluePod = selector
		}

		bnl := util.Namespace(n1, "1")
		if _, ok := blue[bnl]; ok {
			SvcCi.BlueAccess = strings.Join(blue[bnl].Access, "")
		}
		gnl := util.Namespace(n1, "2")
		if _, ok := green[gnl]; ok {
			SvcCi.GreenAccess = strings.Join(green[gnl].Access, "")
		}

		// 获取
		lock := lb.GetLbServiceMap([]string{`"` + SvcCi.Domain + `"`})
		lbVersion := lock.GetVString(SvcCi.Domain)
		SvcCi.LbService = SvcCi.BlueAccess
		if lbVersion == "2" {
			SvcCi.LbService = SvcCi.GreenAccess
		}
	}
	return SvcCi
}

// 2018-02-16 15:00
// 更新负载均衡服务数据
func updateLbPercent(percent int, serviceCiData ci.CloudCiService, q string, version string) {
	q = strings.Replace(q, "{0}", strconv.Itoa(percent), -1)
	q = strings.Replace(q, "{1}", serviceCiData.Domain, -1)
	q = strings.Replace(q, "{2}", util.Namespace(serviceCiData.ServiceName, version), -1)
	sql.Raw(q).Exec()
}

// 2018-02-17 07:38
// 执行滚动更新
func execCiUpdate(serviceData app.CloudAppService, d k8s.RollingParam, username string) error {
	serviceData.MaxUnavailable = int(d.MaxUnavailable)
	serviceData.MaxSurge = int(d.MaxSurge)
	serviceData.MinReady = int(d.MinReadySeconds)
	serviceData.LastModifyTime = util.GetDate()
	serviceData.LastModifyUser = username
	serviceData.ImageTag = d.Images
	logs.Info("执行更新", util.ObjToString(serviceData))
	err := app2.ExecUpdate(serviceData, "image", username)
	return err
}

// 2018-02-17 07:52
// 获取滚动更新数据,并坚持错误
func getRollingData(err error, this *ServiceController) (app.CloudAppService, ci.CloudCiReleaseHistory, kubernetes.Clientset, ci.CloudCiService, string) {
	version := this.GetString("version")
	var msg string
	if err != nil {
		msg = "参数错误"
	}

	services, serviceCiData := getCiServiceData(this)
	if len(services) == 0 {
		msg = "获取服务数据失败"
	}

	serviceData := getImageServiceInfo(services, version)
	if serviceData.ServiceId == 0 {
		msg = "获取服务数据失败,当前版本不存在"
	}

	if len(services) == 2 {
		if services[0].ImageTag == services[1].ImageTag {
			msg = "当前蓝绿版本,版本已经一样,不能继续滚动更新了"
		}
	}

	history := getHistoryData(serviceCiData)
	cl, kerr := k8s.GetClient(serviceCiData.ClusterName)
	if kerr != nil {
		msg = "获取k8s连接失败"
	}
	return serviceData, history, cl, serviceCiData, msg
}

// 2018-02-17 11:50
// 记录服务发布操作日志
func saveServiceLog(this *ServiceController, ciData ci.CloudCiService, msg string) {
	log := ci.CloudCiReleaseLog{}
	log.CreateUser = getServiceUser(this)
	log.Ip = this.Ctx.Request.RemoteAddr
	log.CreateTime = util.GetDate()
	log.Messages = msg
	log.AppName = ciData.AppName
	log.ServiceName = ciData.ServiceName
	log.ClusterName = ciData.ClusterName
	log.Entname = ciData.Entname
	log.Domain = ciData.Domain
	q := sql.InsertSql(log, ci.InsertCloudCiReleaseLog)
	sql.Raw(q).Exec()
}

// 2018-02-18 13:29
// 执行回滚操作
func execUpdateRollbackService(serviceCiData ci.CloudCiService, history ci.CloudCiReleaseHistory,
	serviceData app.CloudAppService, this *ServiceController,
	version string) string {

	serviceData.ServiceVersion = version
	d := k8s.RollingParam{}

	cl, kerr := k8s.GetClient(serviceCiData.ClusterName)
	if kerr != nil {
		return "获取k8s连接失败: "
	}

	if serviceData.ServiceId == 0 {
		return "获取服务的数据失败"
	}

	if serviceData.ImageTag == history.OldImages {
		return "新旧版本一致,不能回滚"
	}

	d.Client = cl
	d.Images = history.OldImages
	d.Namespace = util.Namespace(serviceData.AppName, serviceData.ResourceName)
	d.Name = util.Namespace(serviceCiData.ServiceName, version)
	serviceData.ImageTag = history.OldImages
	logs.Info("获取到更新镜像为", serviceData.ImageTag, history.OldImages)
	err := execCiUpdate(serviceData, d, getServiceUser(this))
	if err != nil {
		return "回滚失败" + err.Error()
	}
	updateLbVersion(version, serviceCiData, serviceData, true)
	saveOperLog(this, "回滚服务 ", serviceCiData, serviceData)
	saveServiceLog(this, serviceCiData, "回滚服务")
	return ""
}