package cluster

import (
	"github.com/astaxie/beego"
	"cloud/models/cluster"
	"cloud/sql"
	"cloud/util"
	"cloud/k8s"
	"cloud/controllers/base/hosts"
	"strings"
	hosts2 "cloud/models/hosts"
	"cloud/cache"
)

type ClusterController struct {
	beego.Controller
}

// 集群管理入口页面
// @router /base/cluster/index [get]
func (this *ClusterController) List() {
	this.TplName = "base/cluster/index.html"
}

// 集群管理入口页面
// @router /base/cluster/image/:id:int [get]
func (this *ClusterController) Images() {
	this.Data["hostId"] = this.Ctx.Input.Param(":id")
	this.TplName = "base/cluster/img.html"
}

// 节点报表入口页面
// @router /base/cluster/report/:id:int [get]
func (this *ClusterController) Report() {
	this.Data["hostId"] = this.Ctx.Input.Param(":id")
	this.TplName = "base/cluster/report.html"
}


// 添加集群
// @router /base/cluster/add [get]
func (this *ClusterController) Add() {
	id,_ := this.GetInt("ClusterId")
	update := cluster.CloudCluster{}
	update.NetworkCart = "6443"
	if id != 0 {
		q := sql.SearchSql(update, cluster.SelectCloudCluster, sql.GetSearchMap("ClusterId", *this.Ctx))
		sql.Raw(q).QueryRow(&update)
		h, p := k8s.GetMasterIp(update.ClusterName)
		update.ApiAddress= h
		update.NetworkCart = p
	}

	this.Data["data"] = update
	this.TplName = "base/cluster/add.html"
}


// @router /base/cluster/detail [get]
func (this *ClusterController) DetailPage() {
	name := this.Ctx.Input.Param(":hi")
	if len(name) < 1 {
		this.Redirect("/base/cluster/list", 302)
		return
	}

	detail := GetClusterDetailData(name)
	if detail.ClusterId < 1 {
		this.Redirect("/base/cluster/list", 302)
		return
	}
	this.Data["data"] = detail
	this.TplName = "base/cluster/detail.html"
}

// 保存集群，初始化集群节点IP
// json
// @router /api/cluster [post]
func (this *ClusterController) Save() {
	d := cluster.CloudCluster{}
	err := this.ParseForm(&d)
	if err != nil {
		this.Ctx.WriteString("参数错误" + err.Error())
		return
	}
	util.SetPublicData(d, getUsername(this), &d)
	q := sql.InsertSql(d, cluster.InsertCloudCluster)
	if d.ClusterId > 0 {
		searchMap := sql.SearchMap{}
		searchMap.Put("ClusterId", d.ClusterId)
		q = sql.UpdateSql(d, cluster.UpdateCloudCluster, searchMap, "CreateTime,CreateUser")
	}

	_, err = sql.Raw(q).Exec()

	// 插入集群节点数据
	h := hosts2.CloudClusterHosts{}
	h.HostType = "master"
	h.HostIp = d.ApiAddress
	h.ApiPort = d.NetworkCart
	h.CreateTime = util.GetDate()
	h.CreateUser = getUsername(this)
	h.LastModifyTime = h.CreateTime
	h.LastModifyUser = h.CreateUser
	h.ClusterName = d.ClusterName
	i := sql.InsertSql(h, hosts2.InsertCloudClusterHosts)
	sql.Raw(i).Exec()
	CacheClusterData()
	data, msg := util.SaveResponse(err, "名称已经被使用")
	util.SaveOperLog(this.GetSession("username"), *this.Ctx, "保存集群操作 "+msg, d.ClusterName)
	setClusterJson(this, data)
}


// json响应
// 集群数据,直返回,集群名称和id的数据
// @router /api/cluster/name [get]
func (this *ClusterController) ClusterName() {
	setClusterJson(this, GetClusterName())
}

// json 响应
// 集群数据获取
// @router /api/cluster [get]
func (this *ClusterController) ClusterData() {
	searchMap := sql.SearchMap{}
	id := this.Ctx.Input.Param(":id")
	key := this.GetString("key")
	if id != "" {
		searchMap.Put("ClusterId", id)
	}

	searchSql := sql.SearchSql(
		cluster.CloudCluster{},
		cluster.SelectCloudCluster,
		searchMap)

	if key != "" && id == "" {
		pkey := sql.Replace(key)
		searchSql += strings.Replace(cluster.SelectCloudClusterWhere, "?", pkey, -1)
	}
	data := make([]k8s.ClusterStatus, 0)
	sql.Raw(searchSql).QueryRows(&data)
	result := make([]k8s.ClusterStatus, 0)
	for _, v := range data{
		r := cache.ClusterCache.Get("data" + v.ClusterName)
		v1 := k8s.ClusterStatus{}
		status := util.RedisObj2Obj(r, &v1)
		if status {
			result = append(result, v1)
		}else{
			result = append(result, v)
			CacheClusterData()
		}
	}
	var r = util.ResponseMap(result, len(result), 1)
	setClusterJson(this, r)
}



// @router /api/cluster/nodes [get]
func (this *ClusterController) NodesData() {
	clusterName := this.GetString("clusterName")
	var check bool = true
	c, err := k8s.GetClient(clusterName)
	if err != nil {
		check = false
	}
	if !check {
		setClusterJson(this, k8s.NodeIp{})
		return
	}
	this.Data["json"] = k8s.GetNodesIp(c)
	this.ServeJSON(false)
}

// @router /api/cluster/delete [*]
func (this *ClusterController) Delete() {
	searchMap := sql.SearchMap{}
	id := this.Ctx.Input.Param(":id")
	searchMap.Put("ClusterId", id)
	cloudCluster := cluster.CloudCluster{}

	q := sql.SearchSql(cloudCluster, cluster.SelectCloudCluster, searchMap)
	sql.Raw(q).QueryRow(&cloudCluster)

	size := len(hosts.GetClusterHosts(cloudCluster.ClusterName))
	if size > 0 {
		msg := "删除失败: 该集群还有节点没有清理,不能删除"
		r := util.ApiResponse(false, msg)
		util.SaveOperLog(
			this.GetSession("username"),
			*this.Ctx, "删除集群 "+msg,
			cloudCluster.ClusterName)

		setClusterJson(this, r)
		return
	}

	q = sql.DeleteSql(cluster.DeleteCloudCluster, searchMap)
	r, err := sql.Raw(q).Exec()
	data := util.DeleteResponse(
		err,
		*this.Ctx,
		"删除集群"+cloudCluster.ClusterName,
		this.GetSession("username"),
		cloudCluster.ClusterName,
		r)
	setClusterJson(this, data)
}