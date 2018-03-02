package registry

import (
	"cloud/sql"
	"cloud/util"
	"github.com/astaxie/beego"
	"cloud/k8s"
	"cloud/models/registry"
	"cloud/controllers/ent"
	"strings"
	"cloud/controllers/base/cluster"
)

type SyncController struct {
	beego.Controller
}

// 镜像同步日志
// @router /image/sync/history [get]
func (this *SyncController) HistoryList() {
	this.TplName = "image/sync/history.html"
}

// 2018-02-05 14:02
// 镜像同步页面
// @router /image/sync/list [get]
func (this *SyncController) SyncList() {
	this.TplName = "image/sync/list.html"
}

// 2018-02-06 17:08
// 审批镜像同步
// @router /api/image/sync/approved/:id:int [post]
func (this *SyncController) ApprovedSave() {
	id := this.Ctx.Input.Param(":id")
	user := getUsesr(this)
	q := `update cloud_image_sync set approved_by="` + user + `", approved_time="` + util.GetDate() + `" where sync_id=` + id
	sql.Raw(q).Exec()
	data, msg := util.SaveResponse(nil, "保存失败")
	util.SaveOperLog(getUsesr(this), *this.Ctx, "同意镜像申请 "+msg, "")
	setSyncJson(this, data)
}

// string
// 镜像同步申请保存
// @router /api/image/sync [post]
func (this *SyncController) SyncSave() {
	d := registry.CloudImageSync{}
	err := this.ParseForm(&d)
	if err != nil {
		this.Ctx.WriteString("参数错误" + err.Error())
		return
	}
	util.SetPublicData(d, util.GetUser(this.GetSession("username")), &d)
	q := sql.InsertSql(d, registry.InsertCloudImageSync)
	if d.SyncId > 0 {
		searchMap := sql.SearchMap{}
		searchMap.Put("SyncId", d.SyncId)
		q = sql.UpdateSql(d,
			registry.UpdateCloudImageSync,
			searchMap,
			"CreateTime,CreateUser")
	}

	_, err = sql.Raw(q).Exec()
	data, msg := util.SaveResponse(err, "保存失败")
	util.SaveOperLog(
		this.GetSession("username"),
		*this.Ctx,
		"保存镜像同步配置 "+msg,
		d.Description)
	setSyncJson(this, data)
}

// 2018-02-07 08:03
// 获取镜像同步页面编辑渲染数据
func getEditData(this *SyncController, data registry.CloudImageSync) {
	copy := this.GetString("copy")
	srcCluster := util.GetSelectOptionName(data.ClusterName)
	targetCluster := util.GetSelectOptionName(data.TargetCluster)
	clusterData := cluster.GetClusterSelect()
	registryData := GetRegistrySelect()
	registryGroupData := GetRegistryGroupSelect(getUsesr(this))

	this.Data["clusterSrc"] = srcCluster + clusterData
	this.Data["clusterTarget"] = targetCluster + clusterData

	imgQ := sql.GetSearchMapV("ClusterName", data.ClusterName, "RegistryGroup", data.RegistryGroup)
	imageData := GetImageSelect(imgQ)
	imgInfo := GetImageDatas(imgQ)
	if copy == "1"{
		data.SyncId = 0
	}
	entData := ent.GetEntnameSelect()
	this.Data["data"] = data
	this.Data["entname"] = entData
	this.Data["entnameSrc"] = util.GetSelectOptionName(data.Entname) + entData
	this.Data["entnameTarget"] = util.GetSelectOptionName(data.TargetEntname) + entData
	this.Data["registrySrc"] = util.GetSelectOptionName(data.Registry) + registryData
	this.Data["registryTarget"] = util.GetSelectOptionName(data.TargetRegistry) + registryData
	this.Data["registryGroup"] = util.GetSelectOptionName(data.RegistryGroup) + registryGroupData
	this.Data["imageData"] = util.GetSelectOptionName(data.ImageName) + imageData
	if len(imgInfo) > 0 {
		this.Data["version"] = util.GetSelectOptionName(data.Version) + GetImageTagSelect(imgInfo[0].Tags)
	}
}

// 2018-02-05 14:02
// 镜像同步页面
// @router /image/sync/add [get]
func (this *SyncController) SyncAdd() {
	data := getSyncData(this,"")
	getEditData(this, data)
	this.TplName = "image/sync/add.html"
}

// 2018-02-06 17:29
// 获取镜像数据
// @router /api/image/sync [get]
func (this *SyncController) SyncDatas() {
	data := []registry.CloudImageSync{}
	user := getUsesr(this)
	searchMap := sql.GetSearchMapV("User", user)
	key := this.GetString("search")
	searchSql := sql.SearchSql(registry.CloudImageSync{}, registry.SelectCloudImageSync, searchMap)
	if len(key) > 4 {
		key = sql.Replace(key)
		searchSql += strings.Replace(registry.SelectCloudImageSyncWhere, "?", key, -1)
	}

	num, _ := sql.OrderByPagingSql(searchSql, "sync_id",
		*this.Ctx.Request,
		&data,
		registry.CloudImageSync{})

	registryData := GetRegistryServerMap()
	result := []registry.CloudImageSync{}
	for _, v := range data {
		temp, ok := registryData.Get(v.ClusterName + v.Registry)
		if ok {
			rv := temp.(registry.CloudRegistryServer)
			server := strings.Split(rv.ServerAddress, ":")
			if len(server) < 2 {
				continue
			}
			v.Registry = rv.ServerDomain + ":" + server[1]
		}
		result = append(result, v)
	}

	r := util.ResponseMap(result,
		sql.CountSearchMap(
			"cloud_image_sync",
			sql.GetSearchMapV("CreateUser", user),
			int(num), key),
		this.GetString("draw"))

	setSyncJson(this, r)
}

// 历史数据
// @router /api/image/sync/history [get]
func (this *SyncController) HistorDatas() {
	data := []k8s.CloudImageSyncLog{}
	user := getUsesr(this)
	searchMap := sql.GetSearchMapV("User", user)
	key := this.GetString("search")
	id := this.GetString("id")
	if id != "" {
		searchMap.Put("LogId", id)
	}
	searchSql := sql.SearchSql(k8s.CloudImageSyncLog{}, registry.SelectCloudImageSyncLog, searchMap)
	if len(key) > 4 {
		key = sql.Replace(key)
		searchSql += strings.Replace(registry.SelectImageSyncLogWhere, "?", key, -1)
	}

	num, _ := sql.OrderByPagingSql(searchSql,
		"log_id",
		*this.Ctx.Request,
		&data,
		k8s.CloudImageSyncLog{})

	r := util.ResponseMap(data,
		sql.CountSearchMap("cloud_image_sync_log", sql.GetSearchMapV("CreateUser", user), int(num), key),
		this.GetString("draw"))
	setSyncJson(this, r)
}

func setSyncJson(this *SyncController, data interface{}) {
	this.Data["json"] = data
	this.ServeJSON(false)
}

// 2018-02-06 20:39
// 获取用户
func getUsesr(this *SyncController) string {
	return util.GetUser(this.GetSession("username"))
}

// 2018-02-06 20:38
// 删除镜像同步
// @router /api/image/sync/:id:int [delete]
func (this *SyncController) SyncDelete() {
	searchMap := sql.GetSearchMap("SyncId", *this.Ctx)
	searchMap.Put("CreateUser", getUsesr(this))
	q := sql.DeleteSql(registry.DeleteCloudImageSync, searchMap) + " and approved_by is null"
	r, _ := sql.Raw(q).Exec()
	data := util.DeleteResponse(nil,
		*this.Ctx,
		"删除镜像同步",
		this.GetSession("username"),
		"",
		r)
	setSyncJson(this, data)
}

// 2018-02-06 22:01
// 获取镜像同步参数
func getImagePushParam(this *SyncController, d registry.CloudImageSync) k8s.ImagePushParam {
	registryData := GetRegistryServerMap()

	imagePushParam := k8s.ImagePushParam{
		RegistryGroup: "",
		ItemName:      d.ImageName,
		Version:       d.Version,
		CreateTime:    util.GetDate(),
		User:          getUsesr(this),
	}
	reg1, ok := registryData.Get(d.ClusterName + d.Registry)
	if ok {
		reg1data := reg1.(registry.CloudRegistryServer)
		servers := strings.Split(reg1data.ServerAddress, ":")
		imagePushParam.Registry1Domain = reg1data.ServerDomain
		imagePushParam.Registry1Auth = util.Base64Encoding(reg1data.Admin + ":" + util.Base64Decoding(reg1data.Password))
		if len(servers) == 2 {
			imagePushParam.Registry1Ip = servers[0]
			imagePushParam.Registry1Port = servers[1]
		}
	}
	reg2, ok := registryData.Get(d.TargetCluster + d.TargetRegistry)
	if ok {
		reg1data := reg2.(registry.CloudRegistryServer)
		servers := strings.Split(reg1data.ServerAddress, ":")
		imagePushParam.Registry2Domain = reg1data.ServerDomain
		imagePushParam.Registry2Auth = util.Base64Encoding(reg1data.Admin + ":" + util.Base64Decoding(reg1data.Password))
		if len(servers) == 2 {
			imagePushParam.Registry2Ip = servers[0]
			imagePushParam.Registry2Port = servers[1]
		}
	}

	return imagePushParam
}

// 2018-02-07 08:04
// 获取数据
func getSyncData(this *SyncController, approved string) registry.CloudImageSync {
	data := registry.CloudImageSync{}
	searchMap := sql.GetSearchMap("SyncId", *this.Ctx)
	searchMap.Put("CreateUser", getUsesr(this))
	q := sql.SearchSql(data, registry.SelectCloudImageSync, searchMap)
	if approved != ""{
		q += `and approved_by != ""`
	}
	sql.Raw(q).QueryRow(&data)
	return data
}

// 2018-02-06 21:50
// 镜像同步启动
// @router /api/image/sync/:id:int [get]
func (this *SyncController) SyncExec() {
	data := getSyncData(this,"1")
	if len(data.ApprovedBy) > 1 {
		param := getImagePushParam(this, data)
		go k8s.ImagePush(data.ClusterName, param)
	}
	this.Ctx.WriteString("后台正在同步中,请稍后查看日志!")
}
