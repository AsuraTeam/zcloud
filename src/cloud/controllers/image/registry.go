package registry

import (
	"github.com/astaxie/beego"
	"cloud/sql"
	"cloud/util"
	"cloud/models/registry"
	"github.com/cesanta/docker_auth/auth_server/server"
	"strings"
	"cloud/k8s"
	"cloud/controllers/base/hosts"
	"github.com/astaxie/beego/logs"
	"cloud/controllers/base/cluster"
	"strconv"
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
	this.Data["cluster"] = cluster.GetClusterSelect()
	id := this.GetString("ServerId")
	// 更新操作
	if id != "0" {
		searchMap := sql.GetSearchMap("UserId", *this.Ctx)
		update := registry.CloudRegistryServer{}
		q := sql.SearchSql(registry.CloudRegistryServer{}, registry.SelectCloudRegistryServer, searchMap)
		sql.Raw(q).QueryRow(&update)
		this.Data["data"] = update
		this.Data["readonly"] = "readonly"
		this.Data["cluster"] = util.GetSelectOptionName(update.ClusterName)
	}
	cf := util.AuthServerConfigFile()
	c, err := server.LoadConfig(cf)
	if err == nil {
		host := strings.Split(this.Ctx.Request.Host, ":")[0]
		this.Data["AuthServer"] = "https://" + host + c.Server.ListenAddress + "/auth"
	}
	this.TplName = "image/registry/add.html"
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
	masterData := []registry.CloudRegistryServer{}

	q := sql.SearchSql(d, registry.SelectCloudRegistryServer, searchMap)
	sql.Raw(q).QueryRows(&masterData)

	util.SetPublicData(d, util.GetUser(this.GetSession("username")), &d)

	this.ServeJSON(false)
	if d.ServerId > 0 {
		q = sql.UpdateSql(d,
			registry.UpdateCloudRegistryServer,
			searchMap,
			registry.UpdateRegistryServerExclude)
		_, err = sql.Raw(q).Exec()
	} else {
		registryData := []registry.CloudRegistryServer{}
		search := sql.GetSearchMapV("Name", d.Name, "ClusterName", d.ClusterName)
		q := sql.SearchSql(
			registry.CloudRegistryServer{},
			registry.SelectCloudRegistryServer,
			search)

		sql.Raw(q).QueryRows(&registryData)
		if len(registryData) > 0 {
			this.Data["json"], _ = util.SaveResponse(err, "服务名已经被占用")
			return
		}
		master, port := hosts.GetMaster(d.ClusterName)
		param := k8s.RegistryParam{
			Master:      master,
			Port:        port,
			ClusterName: d.ClusterName,
			AuthServer:  d.AuthServer,
			Name:        d.Name}
		err = k8s.CreateRegistry(param)
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
	data := []registry.CloudRegistryServer{}
	if name != "" && name != "1" {
		searchMap.Put("ServerDomain", name)
	}

	searchSql := sql.SearchSql(
		registry.CloudRegistryServer{},
		registry.SelectCloudRegistryServer,
		searchMap)
	sql.Raw(searchSql).QueryRows(&data)
	result := []registry.CloudRegistryServer{}
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
			html = append(html, util.GetSelectOption(v.ServerDomain, v.ServerDomain, v.ServerAddress))
		}
	}
	return strings.Join(html, "")
}

// 选择仓库组

// 仓库服务器数据
// @router /api/registry [get]
func (this *ImageController) RegistryServer() {
	data := []registry.CloudRegistryServer{}
	searchMap := sql.SearchMap{}
	id := this.Ctx.Input.Param(":id")
	key := this.GetString("search")
	if id != "" {
		searchMap.Put("UserId", id)
	}
	clusterName := this.GetString("ClusterName")
	if clusterName != "" {
		searchMap.Put("ClusterName", clusterName)
	}

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

	result := []registry.CloudRegistryServer{}
	namespace := util.Namespace("registryv2", "registryv2")
	for _, v := range data {
		if len(v.Access) > 10 {
			v.Password = "******"
			v.ClusterName = util.ObjToString(clusterMap.GetV(v.ClusterName))
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
				h := hostdata[0].HostIp + ":" + strings.Replace(port, "<br>", "", -1)
				v.Access += strings.Replace(registry.SelectRegistryAccess, "?", h, -1)
				v.ServerAddress = hostdata[0].HostIp + ":" + port
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

// @router /api/registry/delete [*]
func (this *ImageController) RegistryServerDelete() {
	searchMap := sql.GetSearchMap("ServerId", *this.Ctx)

	registrData := registry.CloudRegistryServer{}
	q := sql.SearchSql(registrData, registry.SelectCloudRegistryServer, searchMap)
	sql.Raw(q).QueryRow(&registrData)

	err := k8s.DeletelDeployment(
		util.Namespace(
			"registryv2",
			"registryv2"),
		true,
		registrData.Name, registrData.ClusterName)
	q = sql.DeleteSql(registry.DeleteCloudRegistryServer, searchMap)
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
