package lb

import (
	"cloud/sql"
	"cloud/util"
	"github.com/astaxie/beego"
	"cloud/models/lb"
	"github.com/astaxie/beego/logs"
	"strings"
	"cloud/controllers/base/quota"
	"cloud/controllers/ent"
	"cloud/k8s"
)

type LbController struct {
	beego.Controller
}

// 负载均衡管理入口页面
// @router /base/network/lb/list [get]
func (this *LbController) LbList() {
	this.TplName = "base/network/lb/list.html"
}

// 集群负载均衡详情页面
// @router /base/network/lb/detail/:id:int [get]
func (this *LbController) LbDetailPage() {
	id := this.Ctx.Input.Param(":id")
	searchMap := sql.SearchMap{}
	if id != "" {
		searchMap.Put("LbId", id)
	}
	data := lb.CloudLb{}
	q := sql.SearchSql(data, lb.SelectCloudLb, searchMap)
	sql.Raw(q).QueryRow(&data)
	this.Data["data"] = data
	logs.Info("lbdata", data)
	this.TplName = "base/network/lb/detail.html"
}

// 负载均衡管理添加页面
// @router /base/network/lb/add [get]
func (this *LbController) LbAdd() {
	id := this.GetString("LbId")
	entData := ent.GetEntnameSelect()
	update := lb.CloudLb{}
	var entHtml string

	// 更新操作
	if id != "0" {
		searchMap := sql.SearchMap{}
		searchMap.Put("LbId", id)
		q := sql.SearchSql(lb.CloudLb{}, lb.SelectCloudLb, searchMap)
		sql.Raw(q).QueryRow(&update)
		this.Data["readonly"] = "readonly"
		entHtml = util.GetSelectOptionName(update.Entname)
	}
	this.Data["cluster"] = util.GetSelectOptionName(update.ClusterName)
	this.Data["data"] = update
	this.Data["entname"] = entHtml + entData
	this.TplName = "base/network/lb/add.html"
}

// string
// 负载均衡保存
// @router /api/lb [post]
func (this *LbController) LbSave() {
	d := lb.CloudLb{}
	err := this.ParseForm(&d)
	if err != nil {
		this.Ctx.WriteString("参数错误" + err.Error())
		return
	}
	param := k8s.ServiceParam{}
	util.SetPublicData(d, getLbUser(this), &d)
	q := sql.InsertSql(d, lb.InsertCloudLb)
	if d.LbId > 0 {
		param.Update = true
		searchMap := sql.SearchMap{}
		searchMap.Put("LbId", d.LbId)
		q = sql.UpdateSql(d, lb.UpdateCloudLb, searchMap,lb.UpdateLbExclude )
	}else{
		status, msg := checkLbQuota(getLbUser(this))
		if !status {
			data := util.ApiResponse(false, msg)
			setLbJson(this, data)
			return
		}
	}

	param.ClusterName = d.ClusterName
	param.Master, param.MasterPort = k8s.GetMasterIp(d.ClusterName)
	if d.Cpu == "" {
		d.Cpu = "2"
	}
	if d.Memory == "" {
		d.Memory = "4096"
	}
	param.Cpu = d.Cpu
	param.Memory = d.Memory
	if ! strings.Contains(d.HostLogPath, ":") {
		d.HostLogPath = d.HostLogPath + ":" + "/usr/local/nginx/logs/"
	}
	path := strings.Split(d.HostLogPath, ":")
	param.StorageData = `[{"ContainerPath":"`+path[1]+`","HostPath":"`+path[0]+`"}]`
	k8s.CreateNginxLb(param)
	_, err = sql.Raw(q).Exec()

	data, msg := util.SaveResponse(err, "名称已经被使用")
	util.SaveOperLog(
		getLbUser(this),
		*this.Ctx,
		"保存负载均衡配置 "+msg,
		d.LbName)
	setLbJson(this, data)
}

// 2018-02-12 09:45
// 检查镜像仓库配额
// 检查资源配额是否够用
func checkLbQuota(username string) (bool,string) {
	quotaDatas := quota.GetUserQuotaData(username, "")
	for _, v := range quotaDatas {
		if v.LbUsed + 1 > v.LbNumber {
			return false, "负载均衡数量超过配额限制"
		}
	}
	return true, ""
}

// 负载均衡数据
// @router /base/lb [get]
func (this *LbController) LbData() {
	data := make([]lb.CloudLb, 0)
	searchMap := sql.SearchMap{}
	id := this.Ctx.Input.Param(":id")
	key := this.GetString("key")
	if id != "" {
		searchMap.Put("LbId", id)
	}
	searchSql := sql.SearchSql(lb.CloudLb{}, lb.SelectCloudLb, searchMap)
	if key != "" && id == "" {
		searchSql += strings.Replace(lb.SelectCloudLbWhere, "?", sql.Replace(key), -1)
	}

	sql.Raw(searchSql).QueryRows(&data)
	r := util.ResponseMap(data,
		sql.Count("cloud_lb", len(data), key),
		this.GetString("draw"))
	setLbJson(this, r)
}

// 2018-02-01 17:27
func GetLbData(id interface{}) k8s.CloudLb {
	searchMap := sql.SearchMap{}
	searchMap.Put("LbId", id)
	data :=  k8s.GetLbDataSearchMap(searchMap)
	if data != nil {
		dataInterface := data.(interface{})
		return dataInterface.(k8s.CloudLb)
	}
	return k8s.CloudLb{}
}



// json
// 删除负载均衡
// @router /api/network/lb/:id:int [delete]
func (this *LbController) LbDelete() {
	searchMap := sql.SearchMap{}
	id := this.Ctx.Input.Param(":id")
	searchMap.Put("LbId", id)
	template := GetLbData(id)
	q := sql.DeleteSql(lb.DeleteCloudLb, searchMap)
	r, err := sql.Raw(q).Exec()
	data := util.DeleteResponse(err,
		*this.Ctx,
		"删除负载均衡"+template.LbName,
		this.GetSession("username"),
		template.CreateUser, r)
	setLbJson(this, data)
}

func setLbJson(this *LbController, data interface{}) {
	this.Data["json"] = data
	this.ServeJSON(false)
}

func getLbUser(this *LbController) string {
	return util.GetUser(this.GetSession("username"))
}