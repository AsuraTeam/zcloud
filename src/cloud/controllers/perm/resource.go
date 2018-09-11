package perm

import (
	"cloud/sql"
	"cloud/util"
	"github.com/astaxie/beego"
	"cloud/models/perm"
	"strings"
)

type ResourceController struct {
	beego.Controller
}

// api资源管理入口页面
// @router /perm/resource/list [get]
func (this *ResourceController) ResourceList() {
	this.TplName = "perm/resource/list.html"
}

// api资源管理添加页面
// @router /perm/resource/add [get]
func (this *ResourceController) ResourceAdd() {
	id := this.GetString("ResourceId")
	update := perm.CloudApiResource{}
	// 更新操作
	if id != "" {
		searchMap := sql.GetSearchMap("ResourceId", *this.Ctx)
		sql.Raw(sql.SearchSql(perm.CloudApiResource{}, perm.SelectCloudApiResource, searchMap)).QueryRow(&update)
	}
	this.Data["data"] = update
	this.TplName = "perm/resource/add.html"
}

// 获取api资源数据
// 2018-02-06 8:56
// router /api/perm/resource [get]
func (this *ResourceController) ResourceData() {
	// api资源数据
	data := make([]perm.CloudApiResource, 0)
	q := sql.SearchSql(perm.CloudApiResource{}, perm.SelectCloudApiResource, sql.SearchMap{})
	sql.Raw(q).QueryRows(&data)
	setResourceJson(this, data)
}

// string
// api资源保存
// @router /api/perm/resource [post]
func (this *ResourceController) ResourceSave() {
	d := perm.CloudApiResource{}
	err := this.ParseForm(&d)
	if err != nil {
		this.Ctx.WriteString("参数错误" + err.Error())
		return
	}
	util.SetPublicData(d, util.GetUser(this.GetSession("username")), &d)

	q := sql.InsertSql(d, perm.InsertCloudApiResource)
	if d.ResourceId > 0 {
		searchMap := sql.SearchMap{}
		searchMap.Put("ResourceId", d.ResourceId)
		q = sql.UpdateSql(d, perm.UpdateCloudApiResource, searchMap, "CreateTime,CreateResource")
	}
	_, err = sql.Raw(q).Exec()

	data, msg := util.SaveResponse(err, "名称已经被使用")
	util.SaveOperLog(this.GetSession("username"), *this.Ctx, "保存api资源配置 "+msg, d.ApiUrl)
	setResourceJson(this, data)
}

// 获取api资源数据
// 2018-02-06 08:36
// router /api/perm/resource/name [get]
func (this *ResourceController) ResourceDataName() {
	// api资源数据
	data := make([]perm.CloudApiResource, 0)
	q := sql.SearchSql(perm.CloudApiResource{},
		perm.SelectCloudApiResource,
		sql.SearchMap{})
	sql.Raw(q).QueryRows(&data)
	setResourceJson(this, data)
}

// 2018-08-29 14:31
// 更新资源
func UpdateResource() {
	lock := util.Lock{}
	data := make([]perm.CloudApiResource, 0)
	sql.Raw(perm.SelectCloudApiResource).QueryRows(&data)
	for _, v := range data {
		lock.Put(v.Name+v.ApiUrl+v.Method+v.ApiType+v.Parent, "1")
	}
	//sql.Raw(perm.DeleteCloudApiResource).Exec()
	apis := beego.APIS
	for _, v := range apis {
		vs := strings.Split(v, "|")
		if len(vs) > 4 {
			if len(vs[2]) > 0 {

				p := perm.CloudApiResource{}
				p.CreateUser = "system"
				p.CreateTime = util.GetDate()
				p.ApiUrl = vs[0]
				p.ApiUrl = strings.Replace(p.ApiUrl, "//", "/", -1)
				p.Method = strings.Replace(strings.Split(vs[1], ":")[0], "[", "", -1)
				p.Method = strings.Replace(p.Method, "map", "", -1)
				p.Name = vs[2]
				p.ApiType = vs[3]
				p.Parent = vs[4]
				if _, ok := lock.Get(p.Name + p.ApiUrl + p.Method + p.ApiType + p.Parent); ok {
					continue
				}
				sql.Exec(sql.InsertSql(p, perm.InsertCloudApiResource))
			}
		}
	}
}

// 2018-08-30 15:24
// 获取资源树
// @router /api/perm/resource/tree [get]
func (this *ResourceController) GetResourceTree() {
	data3 := make([]perm.CloudApiResource, 0)
	data4 := make([]perm.CloudApiResource, 0)
	data5 := make([]perm.CloudApiResource, 0)
	sql.Raw(perm.SelectPerm3).QueryRows(&data3)
	sql.Raw(perm.SelectPerm4).QueryRows(&data4)
	sql.Raw(perm.SelectPerm5).QueryRows(&data5)
	d := map[string]interface{}{
	}
	// 获取一级菜单
	for _, v := range data3 {
		if _, ok := d[v.Parent ]; !ok {
			if v.Parent == "网络管理" || v.Parent == "集群管理" || v.Parent == "" || v.Parent == "用户中心"{
				continue
			}
			d[v.Parent ] = map[string]interface{}{
				v.ApiType: map[string]interface{}{
				},
			}
		}
		if _, ok := d[v.Parent ].(map[string]interface{})[v.ApiType ]; !ok {
			d[v.Parent ].(map[string]interface{})[v.ApiType ] = map[string]interface{}{}

		}
		if _, ok := d[v.Parent ].(map[string]interface{})[v.ApiType ].(map[string]interface{})[v.Name ]; !ok {
			d[v.Parent ].(map[string]interface{})[v.ApiType ].(map[string]interface{})[v.Name ] = map[string]interface{}{
			}
		}
		for _, v2 := range data4 {
			if v2.ApiType == v.Name {
				d[v.Parent].(map[string]interface{})[v.ApiType ].(map[string]interface{})[v.Name ].(map[string]interface{})[v2.Name ] = map[string]interface{}{
				}
			}
			for _, v3 := range data5 {
				if v3.ApiType == v2.Name {
					if _, ok := d[v.Parent].(map[string]interface{})[v.ApiType].(map[string]interface{})[v.Name].(map[string]interface{})[v2.Name]; ok {
						d[v.Parent].(map[string]interface{})[v.ApiType].(map[string]interface{})[v.Name].(map[string]interface{})[v2.Name].(map[string]interface{})[v3.Name ] = map[string]interface{}{
						}
					}
				}
			}
		}
	}
	setResourceJson(this, d)
}

// api资源数据
// @router /api/perm/resource [get]
func (this *ResourceController) ResourceDatas() {
	data := make([]perm.CloudApiResource, 0)
	searchMap := sql.SearchMap{}
	id := this.Ctx.Input.Param(":id")
	key := this.GetString("search")
	if id != "" {
		searchMap.Put("ResourceId", id)
	}
	searchSql := sql.SearchSql(perm.CloudApiResource{}, perm.SelectCloudApiResource, searchMap)
	if key != "" && id == "" {
		key = sql.Replace(key)
		searchSql += " where 1=1 and (api_url like \"%" + key + "%\" or description like \"%" + key + "%\")"
	}

	num, _ := sql.OrderByPagingSql(searchSql, "api_url",
		*this.Ctx.Request,
		&data,
		perm.CloudApiResource{})

	r := util.ResponseMap(data, sql.Count("cloud_api_resource", int(num), key), this.GetString("draw"))
	setResourceJson(this, r)
}

// json
// 删除api资源
// 2018-02-06 08:36
// @router /api/perm/resource/:id:int [delete]
func (this *ResourceController) ResourceDelete() {
	searchMap := sql.GetSearchMap("ResourceId", *this.Ctx)
	permData := perm.CloudApiResource{}
	q := sql.SearchSql(permData, perm.SelectCloudApiResource, searchMap)
	sql.Raw(q).QueryRow(&permData)

	q = sql.DeleteSql(perm.DeleteCloudApiResource, searchMap)
	r, err := sql.Raw(q).Exec()
	data := util.DeleteResponse(err,
		*this.Ctx, "删除api资源"+permData.ApiUrl,
		this.GetSession("username"),
		permData.CreateUser, r)
	setResourceJson(this, data)
}

func setResourceJson(this *ResourceController, data interface{}) {
	this.Data["json"] = data
	this.ServeJSON(false)
}
