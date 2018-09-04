package hosts

import (
	"github.com/astaxie/beego"
	"cloud/models/hosts"
	"cloud/sql"
	"cloud/util"
	"cloud/k8s"
	"strings"
)

type HostsController struct {
	beego.Controller
}

// 主机管理入口页面
// @router /base/hosts/index [get]
func (this *HostsController) List() {
	this.TplName = "base/cluster/hosts/index.html"
}

// 添加主机标签
// 2018-01-18 21:15
func (this *HostsController) LabelAdd() {
	cloudCluster := getHostData(this.GetString("hostId"))
	this.Data["data"] = cloudCluster
	this.TplName = "base/cluster/hosts/add_label.html"
}

// @router /base/hosts/add [get]
func (this *HostsController) Add() {
	clusterName := this.GetString("ClusterName")
	if len(clusterName) < 1 {
		this.TplName = "base/cluster/index.html"
		return
	}
	this.Data["ClusterName"] = clusterName
	this.TplName = "base/cluster/hosts/add.html"
}

// json
// @router /api/hosts [post]
func (this *HostsController) Save() {
	d := hosts.CloudClusterHosts{}
	err := this.ParseForm(&d)
	if err != nil {
		this.Ctx.WriteString("参数错误" + err.Error())
		return
	}
	searchMap := sql.SearchMap{}
	searchMap.Put("ClusterName", d.ClusterName)
	searchMap.Put("HostType", "master")
	masterData := make([]hosts.CloudClusterHosts, 0)

	q := sql.SearchSql(d, hosts.SelectCloudClusterHosts, searchMap)
	sql.Raw(q).QueryRows(&masterData)
	if len(masterData) > 0 && d.HostType == "master" {
		data := util.ApiResponse(false, "报错错误,主节点已经存在，不能重复添加")
		setHostJson(this, data)
		return
	}

	util.SetPublicData(d, getHostUser(this), &d)

	q = sql.InsertSql(d, hosts.InsertCloudClusterHosts)
	_, err = sql.Raw(q).Exec()
	data, msg := util.SaveResponse(err, "名称已经被使用")
	util.SaveOperLog(getHostUser(this), *this.Ctx, "操作主机配置 "+msg, d.HostIp)
	setHostJson(this, data)
}

// 2018-09-04 10:17
// 获取主机报表数据
// @router /api/cluster/hosts/report/:id:int [get]
func (this *HostsController) GetHostReport() {
	cloudHost := getHostData(this.Ctx.Input.Param(":id"))
	cl, err := k8s.GetClient(cloudHost.ClusterName)
	if err != nil {
		setHostJson(this, util.ApiResponse(false, "获取数据失败: " + err.Error()))
		return
	}
	data := k8s.DescribeNodeResource(cl,cloudHost.HostIp)
	r := util.ResponseMap(data, len(data), this.GetString("draw"))
	setHostJson(this, r)
}

// 获取主机镜像
// @router /api/cluster/hosts/images/:id:int [get]
func (this *HostsController) GetHostImages() {
	cloudHost := getHostData(this.Ctx.Input.Param(":id"))
	data := k8s.GetNodeImage(cloudHost.ClusterName, cloudHost.HostIp)
	r := util.ResponseMap(data, len(data), this.GetString("draw"))
	setHostJson(this, r)
}

// 主机数据获取
// @router /api/hosts/data [get]
func (this *HostsController) HostsData() {
	data := make([]hosts.CloudClusterHosts, 0)
	searchMap := sql.SearchMap{}
	id := this.Ctx.Input.Param(":id")
	ip := this.GetString("ip")
	cluster := this.GetString("cluster")
	if len(cluster) < 1 {
		setHostJson(this, util.ResponseMapError("缺少集群名称"))
		return
	}
	if len(id) > 0 {
		searchMap.Put("HostId", id)
	}
	searchMap.Put("ClusterName", cluster)
	searchSql := sql.SearchSql(
		hosts.CloudClusterHosts{},
		hosts.SelectCloudClusterHosts,
		searchMap)

	if len(ip) > 0 {
		q := ` and (host_ip like "%?%")`
		searchSql += strings.Replace(q, "?", sql.Replace(ip), -1)
	}
	searchSql = sql.SearchSqlPages(searchSql, *this.Ctx.Request)

	num, err := sql.Raw(searchSql).QueryRows(&data)
	c, err := k8s.GetClient(cluster)
	var r map[string]interface{}

	if err == nil {
		returnData := getRedisNodeData(data)
		r = util.ResponseMap(returnData, num, this.GetString("draw"))
	}else{
		r = util.ResponseMapError(err.Error())
	}
	setHostJson(this, r)
	if err == nil {
		go CacheNodeStatus(data, c)
	}
}


// @router /api/hosts/delete [*]
func (this *HostsController) Delete() {
	searchMap := sql.SearchMap{}
	id := this.Ctx.Input.Param(":id")
	searchMap.Put("HostId", id)
	cloudCluster := getHostData(id)
	q := sql.DeleteSql(hosts.DeleteCloudClusterHosts, searchMap)
	r, err := sql.Raw(q).Exec()
	data := util.DeleteResponse(err,
		*this.Ctx,
		"删除主机,主机IP"+cloudCluster.HostIp,
		this.GetSession("username"),
		cloudCluster.ClusterName, r)
	setHostJson(this, data)
}

// 2018-02-12 19:04
// node节点调度操作
// @router /api/cluster/hosts/:id:int [post]
func (this *HostsController) Schedulable() {
	cordonStr := this.GetString("cordon")
	id := this.Ctx.Input.Param(":id")
	var cordon bool
	if cordonStr == "0" {
		cordon = true
	}
	cloudHost := getHostData(id)
	cl, err := k8s.GetClient(cloudHost.ClusterName)
	if err != nil {
		data := util.ApiResponse(false, "获取集群连接失败"+err.Error())
		setHostJson(this, data)
		return
	}
	err = k8s.UpdateNodeStatus(cl, cloudHost.HostIp, cordon)
	if err != nil {
		data := util.ApiResponse(false, "设置失败"+err.Error())
		setHostJson(this, data)
		return
	}
	data := util.ApiResponse(true, "操作成功")
	setHostJson(this, data)
}

// 保存标签
// 2018-01-18 21:25
// @router /api/cluster/label [post]
func (this *HostsController) LabelSave() {
	d := hosts.CloudClusterHosts{}
	this.ParseForm(&d)
	hostData := getHostData(d.HostId)
	err := k8s.UpdateNodeLabels(hostData.ClusterName, hostData.HostIp, d.HostLabel)
	if err != nil {
		setHostJson(this, util.ResponseMapError("操作失败:" + err.Error()))
		return
	}
	searchMap := sql.SearchMap{}
	searchMap.Put("HostId", d.HostId)
	util.SetPublicData(d, getHostUser(this), &d)
	update := sql.UpdateSql(
		d,
		hosts.UpdateCloudClusterHosts,
		searchMap,
		hosts.UpdateExclude)
	_, err = sql.Raw(update).Exec()
	data, msg := util.SaveResponse(err, "数据库操作失败")
	util.SaveOperLog(getHostUser(this), *this.Ctx, "保存主机标签 "+msg, d.HostIp)
	setHostJson(this, data)
}

