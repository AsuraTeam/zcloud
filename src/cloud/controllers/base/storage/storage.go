package storage

import (
	"github.com/astaxie/beego"
	"cloud/models/storage"
	"cloud/sql"
	"cloud/util"
	"cloud/k8s"
	"golang.org/x/crypto/openpgp/errors"
	"cloud/controllers/base/hosts"
	"cloud/controllers/ent"
)

type StorageController struct {
	beego.Controller
}

// 存储管理入口页面
// @router /base/storage/index [get]
func (this *StorageController) StorageList() {
	this.TplName = "base/storage/list.html"
}

// 存储卷添加页面
// 2018-01-31 08:28
// @router /base/storage/add [get]
func (this *StorageController) StorageAdd() {
	id := this.GetString("StorageId")
	update := k8s.CloudStorage{}
	update.StorageSize = "512"
	entData := ent.GetEntnameSelect()
	var entHtml string

	this.Data["selectStorageType1"] = ""
	this.Data["selectStorageType2"] = "checked"
	this.Data["StorageType1"] = "checked"
	update.SharedType = "0"

	if id != "0" {
		searchMap := sql.GetSearchMap("StorageId", *this.Ctx)
		q := sql.SearchSql(k8s.CloudStorage{}, storage.SelectCloudStorage, searchMap)
		sql.Raw(q).QueryRow(&update)
		if update.SharedType == "1" {
			this.Data["selectStorageType1"] = "checked"
			this.Data["selectStorageType2"] = ""
		}
		switch update.StorageType {
		case "Glusterfs":
			this.Data["StorageType2"] = "checked"
			break
		case "Ceph":
			this.Data["StorageType3"] = "checked"
			break
		}
		entHtml = util.GetSelectOptionName(update.Entname)
		this.Data["cluster"]  = util.GetSelectOptionName(update.ClusterName)
	}

	this.Data["entname"] = entHtml + entData
	this.Data["data"] = update
	this.TplName = "base/storage/add.html"
}

// json
// @router /api/storage [post]
func (this *StorageController) StorageSave() {
	d := k8s.CloudStorage{}
	err := this.ParseForm(&d)
	if err != nil {
		this.Ctx.WriteString("参数错误" + err.Error())
		return
	}
	searchMap := sql.SearchMap{}
	searchMap.Put("StorageId", d.StorageId)
	util.SetPublicData(d, util.GetUser(this.GetSession("username")), &d)

	var q string
	if d.StorageId > 0 {
		q = sql.UpdateSql(d, storage.UpdateCloudStorage, searchMap, storage.UpdateStorageExclude)
	} else {
		q = sql.InsertSql(d, storage.InsertCloudStorage)
	}
	_, err = sql.Raw(q).Exec()
	data, msg := util.SaveResponse(err, "名称已经被使用")
	util.SaveOperLog(this.GetSession("username"), *this.Ctx, "操作存储配置 "+msg, d.Name)
	setJson(this, data)
}

// 设置json数据
func setJson(this *StorageController, data interface{}) {
	this.Data["json"] = data
	this.ServeJSON(false)
}

// 存储数据获取
// @router /api/storage/data [get]
func (this *StorageController) StorageData() {
	data := []k8s.CloudStorage{}
	searchMap := sql.SearchMap{}
	id := this.Ctx.Input.Param(":id")
	search := this.GetString("key")
	if len(id) > 0 {
		searchMap.Put("StorageId", id)
	}

	searchSql := sql.SearchSql(k8s.CloudStorage{}, storage.SelectCloudStorage, searchMap)
	if len(search) > 0 {
		searchSql += ` and (name like "%` + sql.Replace(search) + `%")`
	}

	sql.OrderByPagingSql(searchSql, "storage_id",
		*this.Ctx.Request, &data,
		k8s.CloudStorage{})

	result := []k8s.CloudStorage{}

	for _, v := range data {
		mount := k8s.CloudStorageMountInfo{}
		searchMap = sql.GetSearchMapV("ClusterName", v.ClusterName, "StorageName", v.Name)
		q := sql.SearchSql(mount, storage.SelectCloudStorageMountInfo, searchMap)
		sql.Raw(q).QueryRow(&mount)
		v.Status = mount.Status
		result = append(result, v)
	}

	this.Data["json"] = util.GetResponseResult(nil, this.GetString("draw"),
		result, sql.Count("cloud_storage", len(data), search))
	this.ServeJSON(false)
}

// 获取单个存储的信息
// 2018-01-29 21:11
func getStorageData(id interface{}) k8s.CloudStorage {
	searchMap := sql.SearchMap{}
	searchMap.Put("StorageId", id)
	storageData := k8s.CloudStorage{}
	sql.Raw(sql.SearchSql(storageData, storage.SelectCloudStorage, searchMap)).QueryRow(&storageData)
	return storageData
}

// 2018-01-31 14:47
// 获取存储所有数据
func GetStorageName(username string, clustername string) []k8s.CloudStorage {
	data := []k8s.CloudStorage{}
	searchMap := sql.GetSearchMapV("username", username, "ClusterName", clustername)
	q := sql.SearchSql(k8s.CloudStorage{}, "select name from cloud_storage", searchMap)
	sql.Raw(q).QueryRows(&data)
	return data
}

// 删除存储数据
// @router /api/storage/delete/:id:int [delete]
func (this *StorageController) StorageDelete() {
	force := this.GetString("force")
	searchMap := sql.GetSearchMap("StorageId", *this.Ctx)
	storageData := getStorageData(searchMap.Get("StorageId"))
	if storageData.StorageId == 0 && force == "" {
		data := util.ApiResponse(false, "删除存储失败,数据不存在")
		setJson(this, data)
		return
	}

	server := getStorageServerData("", storageData.ClusterName)
	if server.ClusterName == "" && force == "" {
		data := util.ApiResponse(false, errors.UnsupportedError("没有找到集群信息"))
		setJson(this, data)
		return
	}

	master, port := hosts.GetMaster(storageData.ClusterName)

	mount := k8s.CloudStorageMountInfo{}
	qSearchMap := sql.GetSearchMapV("StorageName", storageData.Name, "ClusterName", storageData.ClusterName)
	q := sql.SearchSql(mount, storage.SelectCloudStorageMountInfo, qSearchMap)
	sql.Raw(q).QueryRow(&mount)

	if mount.Status == "1" && force == "" {
		data := util.ApiResponse(false, errors.UnsupportedError("该卷被挂载,不能删除"))
		setJson(this, data)
		return
	}

	namespace := util.Namespace(mount.AppName, mount.ResourceName)

	param := k8s.StorageParam{
		Namespace: namespace,
		Master:    master,
		Port:      port,
		PvcName:   mount.ServiceName,
		ClusterName:storageData.ClusterName,
	}

	err := k8s.DeletePvc(param)
	if err != nil && force == "" {
		data := util.ApiResponse(false, "删除存储失败"+err.Error())
		setJson(this, data)
		return
	}
	orm := sql.GetOrm()
	q = sql.DeleteSql(storage.DeleteCloudStorageMountInfo, searchMap)
	orm.Raw(q).Exec()
	q = sql.DeleteSql(storage.DeleteCloudStorage, searchMap)
	r, err := orm.Raw(q).Exec()
	data := util.DeleteResponse(err, *this.Ctx, "删除存储,存储IP"+storageData.StorageType, this.GetSession("username"), storageData.ClusterName, r)

	setJson(this, data)
}
