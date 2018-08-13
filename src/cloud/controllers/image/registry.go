package registry

import (
	"github.com/astaxie/beego"
	"cloud/sql"
	"cloud/util"
	"cloud/models/registry"
	"strings"
	"cloud/k8s"
	"cloud/controllers/base/hosts"
	"github.com/astaxie/beego/logs"
	"cloud/controllers/base/cluster"
	"strconv"
	"cloud/controllers/ent"
)

type ImageController struct {
	beego.Controller
}

// 主机管理入口页面
// @router /image/registry/list [get]
func (this *ImageController) RegistryServerList() {
	this.TplName = "image/registry/list.html"
}

// @router /image/registry/add [get]
func (this *ImageController) RegistryServerAdd() {
	var entHtml string
	var clusterHtml string
	entData := ent.GetEntnameSelect()

	update := registry.CloudRegistryServer{}

	id := this.GetString("ServerId")
	// 更新操作
	if id != "0" {

		searchMap := sql.GetSearchMap("UserId", *this.Ctx)
		q := sql.SearchSql(registry.CloudRegistryServer{}, registry.SelectCloudRegistryServer, searchMap)
		sql.Raw(q).QueryRow(&update)

		this.Data["readonly"] = "readonly"
		entHtml = util.GetSelectOptionName(update.Entname)
		clusterHtml = util.GetSelectOptionName(update.ClusterName)
	}

	this.Data["entname"] = entHtml + entData
	this.Data["data"] = update
	this.Data["cluster"] = clusterHtml + cluster.GetClusterSelect()
	host := strings.Split(this.Ctx.Request.Host, ":")[0]
	this.Data["AuthServer"] = "https://" + host + ":5001/auth"
	this.TplName = "image/registry/add.html"
}

// 2018-03-02 10:29
// 部署仓库服务
func deployRegistry(d registry.CloudRegistryServer) error {
	master, port := hosts.GetMaster(d.ClusterName)
	param := k8s.RegistryParam{
		Master:      master,
		Port:        port,
		ClusterName: d.ClusterName,
		AuthServer:  d.AuthServer,
		HostPath:    d.HostPath,
		Name:        d.Name}
	err := k8s.CreateRegistry(param)
	return err
}

// 2018-03-02 10:21
// 重建仓库服务
// @param ServerId
// @router /api/registry/recreate [post]
func (this *ImageController) RecreateRegistry() {
	data, _ := getRegistryData(this)
	err := deployRegistry(data)
	var status bool
	if err == nil {
		status = true
	}
	setRegistryServerJson(this, util.ApiResponse(status, err))
}


// json
// @router /api/registry [post]
func (this *ImageController) RegistryServerSave() {
	d := registry.CloudRegistryServer{}
	err := this.ParseForm(&d)
	if err != nil {
		this.Ctx.WriteString("参数错误" + err.Error())
		return
	}

	d.Password = util.Base64Encoding(d.Password)
	searchMap := sql.SearchMap{}
	searchMap.Put("ServerId", d.ServerId)
	masterData := make([]registry.CloudRegistryServer, 0)

	q := sql.SearchSql(d, registry.SelectCloudRegistryServer, searchMap)
	sql.Raw(q).QueryRows(&masterData)

	util.SetPublicData(d, util.GetUser(this.GetSession("username")), &d)

	if d.ServerId > 0 {
		q = sql.UpdateSql(d, registry.UpdateCloudRegistryServer,
			searchMap, registry.UpdateRegistryServerExclude)
		_, err = sql.Raw(q).Exec()
	} else {
		serverData := make([]registry.CloudRegistryServer, 0)
		search := sql.GetSearchMapV("Name", d.Name, "ClusterName", d.ClusterName)
		q := sql.SearchSql(
			registry.CloudRegistryServer{},
			registry.SelectCloudRegistryServer,
			search)

		sql.Raw(q).QueryRows(&serverData)
		if len(serverData) > 0 {
			this.Data["json"], _ = util.SaveResponse(err, "服务名已经被占用")
			return
		}
		err = deployRegistry(d)
		if err == nil {
			q = sql.InsertSql(d, registry.InsertCloudRegistryServer)
			_, err = sql.Raw(q).Exec()
		} else {
			this.Data["json"], _ = util.SaveResponse(err, "操作失败")
			logs.Error("创建 registry server 失败", err)
			return
		}
	}
	data, msg := util.SaveResponse(err, "名称已经被使用")
	util.SaveOperLog(
		this.GetSession("username"),
		*this.Ctx,
		"操作仓库服务 "+msg,
		d.ServerAddress)
	setRegistryServerJson(this, data)
}

// 获取仓库服务器信息
// 2018-01-27 21:17
func GetRegistryServerCluster(serverDomain string, clustername string) registry.CloudRegistryServer {
	data := registry.CloudRegistryServer{}
	searchMap := sql.GetSearchMapV("ServerDomain", serverDomain, "ClusterName", clustername)
	searchSql := sql.SearchSql(
		registry.CloudRegistryServer{},
		registry.SelectCloudRegistryServer,
		searchMap)
	sql.Raw(searchSql).QueryRow(&data)
	return data
}

// 2018-02-06 21:01
func GetRegistryServerMap() util.Lock {
	lock := util.Lock{}
	data := GetRegistryServer("1")
	for _, v := range data {
		lock.Put(v.ClusterName+v.Name, v)
	}
	return lock
}

// 获取镜像服务器信息
// 2018-01-26 10:37
func GetRegistryServer(name string) []registry.CloudRegistryServer {
	searchMap := sql.SearchMap{}
	data := make([]registry.CloudRegistryServer, 0)
	if name != "" && name != "1" {
		searchMap.Put("ServerDomain", name)
	}

	searchSql := sql.SearchSql(
		registry.CloudRegistryServer{},
		registry.SelectCloudRegistryServer,
		searchMap)
	sql.Raw(searchSql).QueryRows(&data)

	result := make([]registry.CloudRegistryServer, 0)
	for _, v := range data {
		if name == "" {
			v.Password = "****"
		}
		result = append(result, v)
	}
	return result
}

// 生成 镜像服务 html
// 2018-01-26 10:41
func GetRegistrySelect() string {
	html := make([]string, 0)
	data := GetRegistryServer("")
	for _, v := range data {
		if v.ServerDomain != "" {
			html = append(html, util.GetSelectOptionName(v.ServerDomain))
		}
	}
	return strings.Join(html, "")
}

// 仓库服务器数据
// @router /api/registry [get]
func (this *ImageController) RegistryServer() {
	data := make([]registry.CloudRegistryServer, 0)
	searchMap := sql.SearchMap{}
	id := this.Ctx.Input.Param(":id")
	key := this.GetString("search")
	if id != "" {
		searchMap.Put("UserId", id)
	}

	searchMap = sql.GetSearchMapValue([]string{"ClusterName"}, *this.Ctx, searchMap)

	searchSql := sql.SearchSql(
		registry.CloudRegistryServer{},
		registry.SelectCloudRegistryServer,
		searchMap)

	if key != "" && id == "" {
		key = sql.Replace(key)
		searchSql += strings.Replace(registry.SelectRegistryServerWhere, "?", key, -1)
	}

	num, _ := sql.Raw(searchSql).QueryRows(&data)
	clusterMap := cluster.GetClusterMap()
	result := make([]registry.CloudRegistryServer, 0)
	namespace := util.Namespace("registryv2", "registryv2")
	for _, v := range data {
		if len(v.Access) > 10 {
			v.Password = "******"
			//v.ClusterName = util.ObjToString(clusterMap.GetV(v.ClusterName))
			result = append(result, v)
			continue
		}
		c, _ := k8s.GetClient(v.ClusterName)
		appData := k8s.GetServicePort(c, namespace, v.Name)

		if len(appData.Spec.Ports) > 0 {
			port := strconv.Itoa(int(appData.Spec.Ports[0].NodePort))
			v.Access = "容器内&nbsp;<br>" + v.Name + "." + namespace + ":" + port + "<br>"
			hostdata := hosts.GetClusterHosts(v.ClusterName)
			if len(hostdata) > 0 {
				h := v.ServerDomain + ":" + strings.Replace(port, "<br>", "", -1)
				logs.Info(h)
				v.Access += strings.Replace(registry.SelectRegistryAccess, "?", h, -1)
				v.ServerAddress = hostdata[0].HostIp + ":" + port
				logs.Info(v.ServerAddress)
				searchMap = sql.GetSearchMapV("Serverid", strconv.FormatInt(v.ServerId, 10))
				u := sql.UpdateSql(v,
					registry.UpdateCloudRegistryServer,
					searchMap,
					registry.UpdateRegistryServerExcludePass)
				sql.Raw(u).Exec()
			}
		}
		v.Password = "******"
		v.ClusterName = util.ObjToString(clusterMap.GetV(v.ClusterName))
		result = append(result, v)
	}
	r := util.ResponseMap(result,
		sql.Count("cloud_registry_server", int(num), key),
		this.GetString("draw"))
	setRegistryServerJson(this, r)
}

// 2018-03-02 10:23
// 获取仓库服务数据
func getRegistryData(this *ImageController) (registry.CloudRegistryServer, sql.SearchMap)  {
	searchMap := sql.GetSearchMap("ServerId", *this.Ctx)
	data := registry.CloudRegistryServer{}
	q := sql.SearchSql(data, registry.SelectCloudRegistryServer, searchMap)
	sql.Raw(q).QueryRow(&data)
	return data,searchMap
}


// @router /api/registry/delete [*]
func (this *ImageController) RegistryServerDelete() {
	registrData,searchMap := getRegistryData(this)
	err := k8s.DeletelDeployment(
		util.Namespace(
			"registryv2",
			"registryv2"),
		true,
		registrData.Name, registrData.ClusterName)
	q := sql.DeleteSql(registry.DeleteCloudRegistryServer, searchMap)
	r, _ := sql.Raw(q).Exec()
	data := util.DeleteResponse(err, *this.Ctx,
		"删除仓库服务,名称:"+registrData.Name,
		this.GetSession("username"),
		registrData.ClusterName, r)
	setRegistryServerJson(this, data)
}

func setRegistryServerJson(this *ImageController, data interface{}) {
	this.Data["json"] = data
	this.ServeJSON(false)
}
