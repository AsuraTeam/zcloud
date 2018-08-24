package quota

import (
	"cloud/sql"
	"cloud/util"
	"github.com/astaxie/beego"
	"cloud/models/quota"
	"strings"
	"cloud/controllers/users"
	"github.com/astaxie/beego/logs"
	"encoding/json"
	"strconv"
)

type ControllerQuota struct {
	beego.Controller
}

// 配额管理入口页面
// @router /base/quota/list [get]
func (this *ControllerQuota) QuotaList() {
	this.TplName = "base/quota/list.html"
}

// 集群配额详情页面
// @router /base/quota/detail/:id:int [get]
func (this *ControllerQuota) QuotaDetailPage() {
	var searchMap sql.SearchMap
	id := this.Ctx.Input.Param(":id")

	if id != "" {
		searchMap.Put("QuotaId", id)
	}

	name := this.Ctx.Input.Param(":hi")
	if name != "" {
		searchMap.Put("QuotaName", name)
	}

	data := quota.CloudQuota{}
	q := sql.SearchSql(
		data,
		quota.SelectCloudQuota,
		searchMap)
	sql.Raw(q).QueryRow(&data)
	logs.Info(util.ObjToString(data))

	result := quota.QuotaUsed{}
	t, _ := json.Marshal(data)
	json.Unmarshal(t, &result)
	if data.UserName != "" {
		result = setQuotaUserUsed(data.UserName, result)
	}
	if data.GroupName != "" {
		result = setQuotaGroupUsed(data.GroupName, result)
	}

	this.Data["data"] = result
	this.TplName = "base/quota/detail.html"
}

// 2018-02-11 10:04
// 配额管理添加页面
// @router /base/quota/add [get]
func (this *ControllerQuota) QuotaAdd() {
	id := this.GetString("QuotaId")
	update := quota.GetDefaultQuota()
	var groupsHtml string
	var userHtml string
	// 更新操作
	if id != "0" {
		searchMap := sql.SearchMap{}
		searchMap.Put("QuotaId", id)

		q := sql.SearchSql(
			quota.CloudQuota{},
			quota.SelectCloudQuota,
			searchMap)
		sql.Raw(q).QueryRow(&update)
		if update.GroupName != "" {
			groupsHtml = util.GetSelectOptionName(update.GroupName)
		}
		if update.UserName != "" {
			userHtml = util.GetSelectOptionName(update.UserName)
		}
		this.Data["readonly"] = "readonly"
	}
	this.Data["data"] = update
	this.Data["users"] = userHtml + users.GetUserSelect()
	this.Data["groups"] = groupsHtml + users.GetGroupsSelect()
	this.TplName = "base/quota/add.html"
}

// string
// 配额保存
// @router /api/quota [post]
func (this *ControllerQuota) QuotaSave() {
	d := quota.CloudQuota{}
	d.QuotaName = strings.Replace(d.QuotaName,"--", "-", -1)
	err := this.ParseForm(&d)
	if err != nil {
		this.Ctx.WriteString("参数错误" + err.Error())
		return
	}
	util.SetPublicData(d, getQuotaUser(this), &d)
	q := sql.InsertSql(d, quota.InsertCloudQuota)
	if d.QuotaId > 0 {
		searchMap := sql.SearchMap{}
		searchMap.Put("QuotaId", d.QuotaId)

		q = sql.UpdateSql(
			d,
			quota.UpdateCloudQuota,
			searchMap,
			quota.UpdateCloudQuotaExclude)
		_, err = sql.Raw(q).Exec()
	}
	_, err = sql.Raw(q).Exec()
	data, msg := util.SaveResponse(err, "名称已经被使用")
	util.SaveOperLog(
		getQuotaUser(this),
		*this.Ctx, "保存配额配置 "+msg,
		d.QuotaName)

	setQuotaJson(this, data)
}

// 配额数据
// @router /base/quota [get]
func (this *ControllerQuota) QuotaData() {
	data := make([]quota.CloudQuota, 0)
	searchMap := sql.SearchMap{}
	id := this.Ctx.Input.Param(":id")
	key := this.GetString("key")
	if id != "" {
		searchMap.Put("QuotaId", id)
	}

	usedData := queryAppQuotaUsed("")

	searchSql := sql.SearchSql(
		quota.CloudQuota{},
		quota.SelectCloudQuota,
		searchMap)

	if key != "" && id == "" {
		key = sql.Replace(key)
		searchSql += strings.Replace(quota.SelectCloudQuotaWhere, "?", key, -1)
	}

	num, err := sql.Raw(searchSql).QueryRows(&data)

	result := make([]quota.CloudQuota, 0)
	for _, v := range data{
		v.Status = "<span class='RunningNoTop'>未使用</span>"
		if _, ok := usedData.Get(v.QuotaName) ; ok {
			v.Status = "<span class='FailNoTop'>已使用</span>&nbsp;/&nbsp;" + usedData.GetVString(v.QuotaName)
		}
		result = append(result, v)
	}

	var r = util.ResponseMap(result, num, 1)
	if err != nil {
		util.ResponseMapError(err.Error())
	}
	setQuotaJson(this, r)

}

// @router /api/quota/name [get]
func (this *ControllerQuota) GetQuotaName() {
	data := make([]quota.CloudQuotaName, 0)

	searchSql := sql.SearchSql(
		quota.CloudQuota{},
		quota.SelectCloudQuota,
		sql.SearchMap{})

	sql.Raw(searchSql).QueryRows(&data)
	this.Data["json"] = data
	this.ServeJSON(false)
}

// 2018-02-11 18:30
func getQuotaData(searchMap sql.SearchMap) []quota.CloudQuota {
	template := make([]quota.CloudQuota, 0)
	q := sql.SearchSql(quota.CloudQuota{}, quota.SelectCloudQuota, searchMap)
	sql.Raw(q).QueryRows(&template)
	return template
}

// json
// 删除配额
// @router /api/quota/:id:int [delete]
func (this *ControllerQuota) QuotaDelete() {
	searchMap := sql.SearchMap{}
	id := this.Ctx.Input.Param(":id")
	searchMap.Put("QuotaId", id)
	dataQ := getQuotaData(searchMap)
	if len(dataQ) == 0 {
		data := util.ApiResponse(false, "删除失败,没有找到对应的资源")
		setQuotaJson(this, data)
		return
	}
	name := dataQ[0].QuotaName
	usedData := queryAppQuotaUsed(name)
	if _, ok := usedData.Get(name); ok {
		data := util.ApiResponse(false, "删除失败,该配额正在被使用")
		setQuotaJson(this, data)
		return
	}

	qData := dataQ[0]
	q := sql.DeleteSql(quota.DeleteCloudQuota, searchMap)
	r, err := sql.Raw(q).Exec()

	data := util.DeleteResponse(err,
		*this.Ctx, "删除配额"+qData.QuotaName,
		getQuotaUser(this),
		qData.CreateUser,
		r)

	setQuotaJson(this, data)
}

// 2018-02-11 21:53
// 获取用户的配额数据
func GetUserQuotaData(username string, quotaName string) []quota.QuotaUsed {
	freeQuotas := make([]quota.QuotaUsed, 0)

	dataQ := getQuotaData(
		sql.GetSearchMapV(
			"UserName", username,
			"QuotaName", quotaName,
		),
	)

	if len(dataQ) > 0 {
		for _, v := range dataQ {
			result := quota.QuotaUsed{}
			t, _ := json.Marshal(v)
			json.Unmarshal(t, &result)
			result = setQuotaUserUsed(username, result)
			freeQuotas = append(freeQuotas, result)
		}
	}
	logs.Info("freeQuotas", freeQuotas)
	// 如果用户没有可用配额,再查看是否有组的配额
	if len(freeQuotas) == 0 {
		detps := users.GetUserDept(username)
		for _, dept := range detps {
			deptQuotaData := getQuotaData(
				sql.GetSearchMapV("GroupName", dept,
					"QuotaName", quotaName,
				),
			)
			if len(deptQuotaData) == 0 {
				continue
			}
			result := quota.QuotaUsed{}
			t, _ := json.Marshal(deptQuotaData[0])
			json.Unmarshal(t, &result)
			logs.Info("result", util.ObjToString(result))
			result = setQuotaGroupUsed(dept, result)
			freeQuotas = append(freeQuotas, result)
		}
	}
	return freeQuotas
}

// 2018-02-11 21:30
// 获取用户可用配额
func GetUserQuota(username string, quotaType string) string {
	freeQuotas := make([]string, 0)
	dataQ := getQuotaData(sql.GetSearchMapV("UserName", username))
	if len(dataQ) > 0 {
		for _, v := range dataQ {
			result := quota.QuotaUsed{}
			t, _ := json.Marshal(v)
			json.Unmarshal(t, &result)
			result = setQuotaUserUsed(username, result)
			if result.CpuFree > 0 && result.MemoryFree * 1024 > 512 {
				freeQuotas = append(
					freeQuotas,
					getFreeQuota(
						quotaType,
						result,
						freeQuotas)...)
			}
		}
	}
	logs.Info("freeQuotas", freeQuotas)
	// 如果用户没有可用配额,再查看是否有组的配额

	if len(freeQuotas) == 0 {
		detps := users.GetUserDept(username)

		for _, dept := range detps {
			searchMap := sql.GetSearchMapV("GroupName", dept)
			deptQuotaData := getQuotaData(searchMap)
			if len(deptQuotaData) > 0 {
				for _, v := range deptQuotaData {
					result := quota.QuotaUsed{}
					t, _ := json.Marshal(v)
					json.Unmarshal(t, &result)
					result = setQuotaGroupUsed(dept, result)
					logs.Info(result.CpuFree, result.MemoryFree, result.MemoryFree*1024)
					if result.CpuFree > 0 && result.MemoryFree*1024 > 512 {
						freeQuotas = getFreeQuota(quotaType, result, freeQuotas)
					}
				}
			}
		}
	}
	var option string
	for _, v := range freeQuotas {
		option += util.GetSelectOptionName(v)
	}
	return option
}

// 2018-02-11 18:57
// 获取是否有可用配额
func getFreeQuota(quotaType string, result quota.QuotaUsed, freeQuotas []string) []string {
	var free bool
	switch quotaType {
	case "app":
		if result.AppFree > 0  && result.ServiceFree > 0{
			free = true
		}
		break
	case "service":
		if result.ServiceFree > 0 && result.PodFree > 0 {
			free = true
		}
		break
	case "pod":
		if result.PodFree > 0 {
			free = true
		}
		break
	case "lb":
		if result.LbFree > 0 {
			free = true
		}
		break
	case "job":
		if result.JobFree > 0 {
			free = true
		}
		break
	case "pipiline":
		if result.PipelineFree > 0 {
			free = true
		}
		break
	default:
		break
	}
	if free {
		freeQuotas = append(freeQuotas, result.QuotaName)
	}
	return freeQuotas
}

// 2018-02-12 06:24
// 查询是否有应用使用资源配额
func queryAppQuotaUsed(quotaName string) util.Lock {
	var qname string
	if quotaName != "" {
		qname = strings.Replace(quota.SelectAppQuotaUsedWhere, "?", quotaName, -1)
	}
	qgroupby := strings.Replace(quota.SelectAppQuotaUsed, "?", qname, -1)

	qdata := make([]quota.QuotaAppUsed, 0)
	sql.Raw(qgroupby).QueryRows(&qdata)
	mapData := util.Lock{}
	for _, v := range qdata{
		mapData.Put(v.ResourceName, strconv.FormatInt(v.Cnt, 10))
	}
	return mapData
}

func getQuotaUser(this *ControllerQuota) string {
	return util.GetUser(this.GetSession("username"))
}

func setQuotaJson(this *ControllerQuota, data interface{})  {
	this.Data["json"] = data
	this.ServeJSON(false)
}