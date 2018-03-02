package app

import (
	"cloud/sql"
	"cloud/util"
	"cloud/models/app"
	"github.com/astaxie/beego"
	"cloud/k8s"
	"github.com/astaxie/beego/logs"
	"golang.org/x/crypto/openpgp/errors"
	"database/sql/driver"
	"strings"
)

type DataController struct {
	beego.Controller
}

// 查询更新数据
func getDataInfo(searchMap sql.SearchMap) app.CloudConfigData {
	update := app.CloudConfigData{}
	q := sql.SearchSql(app.CloudConfigData{}, app.SelectCloudConfigData, searchMap)
	sql.Raw(q).QueryRow(&update)
	return update
}

// 配置文件数据管理添加页面
// @router /application/configure/data/add [get]
func (this *DataController) ConfigDataAdd() {
	this.Data["ConfigureName"] = this.GetString("ConfigureName")
	this.Data["ConfigureId"] = this.GetString("ConfigureId")
	// 更新操作
	if this.GetString("DataId") != "0" {
		searchMap := sql.GetSearchMap("DataId", *this.Ctx)
		this.Data["data"] = getDataInfo(searchMap)
		this.Data["readonly"] = "readonly"
	}
	this.TplName = "application/configure/data/add.html"
}

// string
// 配置文件数据保存
// @router /api/configure/data [post]
func (this *DataController) ConfigDataSave() {
	d := app.CloudConfigData{}
	err := this.ParseForm(&d)
	if err != nil {
		this.Ctx.WriteString("参数错误" + err.Error())
		return
	}

	util.SetPublicData(d, util.GetUser(this.GetSession("username")), &d)

	q := sql.InsertSql(d, app.InsertCloudConfigData)
	if d.DataId > 0 {
		searchMap := sql.SearchMap{}
		searchMap.Put("DataId", d.DataId)
		q = sql.UpdateSql(d, app.UpdateCloudConfigData, searchMap, app.UpdateConfigDataExclude)
	}
	_, err = sql.Raw(q).Exec()

	data, msg := util.SaveResponse(err, "名称已经被使用")
	util.SaveOperLog(this.GetSession("username"), *this.Ctx, "保存配置文件数据配置 "+msg, d.DataName)
	setDataJson(this, data)
	updateK8sConfigMap(d)
}

// 配置文件数据名称数据
// @router /api/configure/data/name [get]
func (this *DataController) GetConfigDataName() {
	data := []app.ConfigDataName{}
	searchSql := sql.SearchSql(app.CloudConfigData{}, app.SelectCloudConfigData, sql.SearchMap{})
	sql.Raw(searchSql).QueryRows(&data)
	this.Data["json"] = data
	this.ServeJSON(false)
}

// 配置文件数据数据
// @router /api/configure/data [get]
func (this *DataController) ConfigData() {
	data := []app.CloudConfigData{}
	searchMap := sql.GetSearchMap("DataId", *this.Ctx)
	configureId := this.GetString("ConfigureId")
	if configureId != "" {
		searchMap.Put("ConfigureId", configureId)
	}
	key := this.GetString("key")
	searchSql := sql.SearchSql(app.CloudConfigData{}, app.SelectCloudConfigData, searchMap)
	if key != "" {
		searchSql += strings.Replace(app.SelectConfigDataWhere, "?", sql.Replace(key), -1)
	}
	num, err := sql.Raw(searchSql).QueryRows(&data)
	r := util.GetResponseResult(err, this.GetString("draw"), data, sql.Count("cloud_config_data", int(num), key))
	setDataJson(this, r)
}

// json
// 删除配置文件数据
// @router /api/configure/data/:id:int [delete]
func (this *DataController) ConfigDataDelete() {
	searchMap := sql.GetSearchMap("DataId", *this.Ctx)
	configure := getDataInfo(searchMap)

	configureData := getMasterConfigure(configure.ConfigureId)
	mountData := getMountData(configure.DataName, configureData.ClusterName, configure.DataName)
	if len(mountData) > 0 {
		logs.Info("该项目被挂载不能删除", configure.ConfigureName)
		data := util.DeleteResponse(errors.InvalidArgumentError("已经被挂载,不能删除"),
			*this.Ctx, "删除配置文件数据"+configure.DataName,
			this.GetSession("username"),
			configure.CreateUser, driver.ResultNoRows)
		setDataJson(this, data)
		return
	}
	q := sql.DeleteSql(app.DeleteCloudConfigData, searchMap)
	r, err := sql.Raw(q).Exec()
	data := util.DeleteResponse(err,
		*this.Ctx,
		"删除配置文件数据"+configure.DataName,
		this.GetSession("username"),
		configure.CreateUser, r)
	setDataJson(this, data)
}

// 从数据库查询组名称
// 2018-01-17 16:24
func GetConfgData(configureName string, cluster string) (map[string]interface{}) {
	configData := app.CloudAppConfigure{}
	searchMap := sql.GetSearchMapV("ConfigureName", configureName, "ClusterName", cluster)
	searchSql := sql.SearchSql(app.CloudAppConfigure{}, app.SelectCloudAppConfigure, searchMap)
	sql.Raw(searchSql).QueryRow(&configData)

	data := []app.CloudConfigData{}
	searchMap = sql.SearchMap{}
	searchMap.Put("ConfigureId", configData.ConfigureId)
	searchSql = sql.SearchSql(app.CloudConfigData{}, app.SelectCloudConfigData, searchMap)
	sql.Raw(searchSql).QueryRows(&data)

	mapData := make(map[string]interface{})
	for _, v := range data {
		mapData[v.DataName] = v.Data
	}
	return mapData
}

// 查询配置挂载的数据
// 2018-01-18 10:50
func getMountData(configname string, cluster string, dataName string) []k8s.CloudConfigureMount {
	// 首先查询已经mount的数据
	mountData := []k8s.CloudConfigureMount{}
	searchMap := sql.GetSearchMapV("ClusterName", cluster, "ConfigureName", configname)
	if dataName != "" {
		searchMap.Put("DataName", dataName)
	}
	q := sql.SearchSql(k8s.CloudConfigureMount{}, k8s.SelectCloudConfigureMount, searchMap)
	sql.Raw(q).QueryRows(&mountData)
	return mountData
}

// 获取主配置文件
// 2018-01-18 11:12
func getMasterConfigure(configId int64) app.CloudAppConfigure {
	searchMap := sql.SearchMap{}
	searchMap.Put("ConfigureId", configId)
	configureData := getUpdateData(searchMap)
	return configureData
}

// 2018-01-18 10:39
// 保存后更新 k8s configmap
func updateK8sConfigMap(d app.CloudConfigData) {
	configureData := getMasterConfigure(d.ConfigureId)
	logs.Info(configureData)
	mountData := getMountData(d.ConfigureName, configureData.ClusterName, d.DataName)
	if len(mountData) < 1 {
		logs.Info("该配置没有被挂载", d.ConfigureName)
		return
	}
	logs.Info("mountdata", util.ObjToString(mountData))
	for _, v := range mountData {

		configData := k8s.ConfigureData{}
		configData.DataName = d.ConfigureName
		configData.ConfigDbData = GetConfgData(d.ConfigureName, configureData.ClusterName)

		param := k8s.ServiceParam{}
		param.ConfigureData = []k8s.ConfigureData{configData}
		param.Namespace = v.Namespace

		cl2, _ := k8s.GetYamlClient(v.ClusterName, "", "v1", "api")
		param.Cl2 = cl2

		k8s.CreateConfigmap(param)
	}
}

func setDataJson(this *DataController, data interface{})  {
	this.Data["json"] = data
	this.ServeJSON(false)
}