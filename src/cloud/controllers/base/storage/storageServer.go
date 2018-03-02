package storage

import (
	"github.com/astaxie/beego"
	"cloud/models/storage"
	"cloud/sql"
	"cloud/util"
	"cloud/controllers/ent"
	"cloud/k8s"
	"strings"
)

type StorageServerController struct {
	beego.Controller
}

// 存储管理入口页面
// @router /base/storage/list [get]
func (this *StorageServerController) StorageServerList() {
	this.TplName = "base/storage/server/list.html"
}

// 存储服务添加页面
// 2018-02-08 09:03
// @router /base/storage/add [get]
func (this *StorageServerController) StorageServerAdd() {
	id := this.GetString("ServerId")
	update := storage.CloudStorageServer{}
	entData := ent.GetEntnameSelect()
	var entHtml string
	if id != "0" {
		searchMap := sql.GetSearchMap("ServerId", *this.Ctx)
		q := sql.SearchSql(storage.CloudStorageServer{}, storage.SelectCloudStorageServer, searchMap)
		sql.Raw(q).QueryRow(&update)
		entHtml = util.GetSelectOptionName(update.Entname)
		this.Data["cluster"] = util.GetSelectOptionName(update.ClusterName)
	}

	switch update.StorageType {
	case "Glusterfs":
		this.Data["StorageType2"] = "checked"
		break
	case "Ceph":
		this.Data["StorageType3"] = "checked"
		break
	default:
		this.Data["StorageType1"] = "checked"
		break
	}

	this.Data["data"] = update
	this.Data["entname"] = entHtml + entData

	this.TplName = "base/storage/server/add.html"
}

// json
// @router /api/storage [post]
func (this *StorageServerController) StorageServerSave() {
	d := storage.CloudStorageServer{}
	err := this.ParseForm(&d)
	if err != nil {
		this.Ctx.WriteString("参数错误" + err.Error())
		return
	}
	var q string
	searchMap := sql.SearchMap{}
	if d.ServerId == 0 {
		searchMap.Put("ClusterName", d.ClusterName)
		searchMap.Put("StorageType", d.StorageType)
		masterData := []storage.CloudStorageServer{}
		q = sql.SearchSql(d, storage.SelectCloudStorageServer, searchMap)
		sql.Raw(q).QueryRows(&masterData)
		if len(masterData) > 0 {
			r := util.ApiResponse(false, "报错错误,主节点已经存在，不能重复添加")
			setServerJson(this, r)
			return
		}
		util.SetPublicData(d, util.GetUser(this.GetSession("username")), &d)
		q = sql.InsertSql(d, storage.InsertCloudStorageServer)
		master,port := k8s.GetMasterIp(d.ClusterName)
		param := k8s.StorageParam{
			ClusterName:d.ClusterName,
			Master: master,
			Port: port,
			HostPath:d.HostPath,
			Namespace: util.Namespace("nfs", "nfs"),
		}
		switch d.StorageType {
		case "Nfs":
			k8s.CreateNfsStorageServer(param)
			break
		case "Glusterfs":
			k8s.CreateGlusterfs(param)
			break
		case "Ceph":
			break
		}

	}else{
		searchMap.Put("ServerId", d.ServerId)
		q = sql.UpdateSql(d, storage.UpdateCloudStorageServer, searchMap, storage.UpdateStorageServerWhere)
	}
	_, err = sql.Raw(q).Exec()
	data, msg := util.SaveResponse(err, "名称已经被使用")
	util.SaveOperLog(this.GetSession("username"), *this.Ctx, "操作存储配置 "+msg, d.ServerAddress)
	setServerJson(this, data)
}

// 设置json数据
func setServerJson(this *StorageServerController, data interface{}) {
	this.Data["json"] = data
	this.ServeJSON(false)
}

// 2018-02-07 20:54
// 存储数据获取
// @router /api/storage/server [get]
func (this *StorageServerController) StorageServerData() {
	data := []storage.CloudStorageServer{}
	searchMap := sql.SearchMap{}
	search := this.GetString("key")
	searchSql := sql.SearchSql(
		storage.CloudStorageServer{},
		storage.SelectCloudStorageServer,
		searchMap)

	if len(search) > 0 {
		q := ` where 1=1 and (name like "%?%") `
		searchSql += strings.Replace(q, "?", sql.Replace(search), -1)
	}

	num, err := sql.Raw(searchSql).QueryRows(&data)
	r := util.GetResponseResult(err,
		this.GetString("draw"), data, sql.Count("cloud_storage_server", int(num), search))
	setServerJson(this, r)
}

// 获取单个存储的信息
// 2018-01-31 11:132
func getStorageServerData(id interface{}, clusterName string) storage.CloudStorageServer {
	searchMap := sql.SearchMap{}
	if id != "" {
		searchMap.Put("ServerId", id)
	}else{
		searchMap.Put("ClusterName", clusterName)
	}

	cloudCluster := storage.CloudStorageServer{}
	q := sql.SearchSql(cloudCluster, storage.SelectCloudStorageServer, searchMap)
	sql.Raw(q).QueryRow(&cloudCluster)
	return cloudCluster
}

// 删除存储服务
// @router /api/storage/server/delete [delete]
func (this *StorageServerController) StorageServerDelete() {
	searchMap := sql.SearchMap{}
	id := this.Ctx.Input.Param(":id")
	searchMap.Put("ServerId", id)
	q := sql.DeleteSql(storage.DeleteCloudStorageServer, searchMap)
	sql.Raw(q).Exec()
	setServerJson(this, util.ApiResponse(true, "删除成功"))
}
