package registry

import (
	"github.com/astaxie/beego"
	"cloud/sql"
	"cloud/util"
	"cloud/models/registry"
	"strings"
	"cloud/k8s"
)

type BaseController struct {
	beego.Controller
}

// 主机管理入口页面
// @router /image/registry/list [get]
func (this *BaseController) BaseList() {
	this.TplName = "image/base/list.html"
}

// 生成 基础镜像服务 html
// 2018-02-09 17:54
func GetBaseImageSelect() string {
	data := make([]registry.CloudImageBase, 0)
	searchSql := sql.SearchSql(registry.CloudImageBase{}, registry.SelectCloudImageBase, sql.SearchMap{})
	sql.Raw(searchSql).QueryRows(&data)
	html := make([]string, 0)
	for _, v := range data {
		server := GetRegistryServer(v.RegistryServer)
		if len(server) > 0 {
			servers := strings.Split(server[0].ServerAddress, ":")
			if len(servers) > 1 {
				html = append(html, util.GetSelectOptionName(v.RegistryServer+":"+servers[1]+"/"+v.ImageName))
			}
		}
	}
	return strings.Join(html, "")
}

// @router /image/registry/add [get]
func (this *BaseController) BaseAdd() {
	update := registry.CloudImageBase{}
	id := this.GetString("BaseId")
	regData := GetRegistrySelect()
	var regHtml string
	// 更新操作
	if id != "0" {
		searchMap := sql.GetSearchMap("BaseId", *this.Ctx)
		q := sql.SearchSql(
			registry.CloudImageBase{},
			registry.SelectCloudImageBase,
			searchMap)

		sql.Raw(q).QueryRow(&update)
		regHtml = util.GetSelectOptionName(update.RegistryServer)
	}
	this.Data["registryServer"] = regHtml + regData
	this.Data["data"] = update
	this.TplName = "image/base/add.html"
}

// 2018-02-09 17:10
// 检查镜像是否存在
func checkImageExists(d registry.CloudImageBase) (interface{}, bool) {
	server := GetRegistryServer(d.RegistryServer)
	if len(server) == 0  {
		data := util.ApiResponse(false, "没有找到仓库服务器")
		return data, false
	}

	images := strings.Split(d.ImageName, ":")
	if len(images) == 1 {
		data := util.ApiResponse(false, "镜像格式错误")
		return data, false
	}

	servers := strings.Split(server[0].ServerAddress, ":")
	if len(servers) == 1 {
		data := util.ApiResponse(false, "不能获取到仓库访问地址和端口")
		return data, false
	}

	status := k8s.CheckImageExists(d.RegistryServer+":"+servers[1],
		server[0].Admin,
		util.Base64Decoding(server[0].Password),
		images[0],
		images[1])
	if ! status {
		data := util.ApiResponse(false, "镜像不存在")
		return data, false
	}
	return "", true
}

// json
// @router /api/registry [post]
func (this *BaseController) BaseSave() {
	d := registry.CloudImageBase{}
	err := this.ParseForm(&d)
	if err != nil {
		this.Ctx.WriteString("参数错误" + err.Error())
		return
	}

	// 检查镜像是否存在
	data, status := checkImageExists(d)
	if ! status {
		setBaseJson(this, data)
		return
	}

	searchMap := sql.SearchMap{}
	searchMap.Put("BaseId", d.BaseId)
	masterData := make([]registry.CloudImageBase, 0)

	q := sql.SearchSql(d, registry.SelectCloudImageBase, searchMap)
	sql.Raw(q).QueryRows(&masterData)
	util.SetPublicData(d, util.GetUser(this.GetSession("username")), &d)

	q = sql.InsertSql(d, registry.InsertCloudImageBase)
	if d.BaseId > 0 {

		q = sql.UpdateSql(d,
			registry.UpdateCloudImageBase,
			searchMap,
			registry.UpdateCloudRedisPermExclude)
	}
	sql.Raw(q).Exec()

	data, msg := util.SaveResponse(err, "名称已经被使用")
	util.SaveOperLog(
		this.GetSession("username"),
		*this.Ctx, "操作基础镜像"+msg,
		d.RegistryServer+d.ImageName)

	setBaseJson(this, data)
}

// 基础镜像服务器数据
// @router /api/registry [get]
func (this *BaseController) Base() {
	data := make([]registry.CloudImageBase, 0)
	searchMap := sql.SearchMap{}
	key := this.GetString("search")
	searchSql := sql.SearchSql(registry.CloudImageBase{},
		registry.SelectCloudImageBase,
		searchMap)

	if key != "" {
		key = sql.Replace(key)
		q := strings.Replace(registry.SelectCloudBaseWhere, "?", key, -1)
		searchSql += q
	}

	num, err := sql.OrderByPagingSql(searchSql,
		"create_time",
		*this.Ctx.Request,
		&data,
		registry.CloudImageBase{})

	r := util.ResponseMap(data,
		sql.Count("cloud_image_base", int(num), key),
		this.GetString("draw"))

	if err != nil {
		r = util.ResponseMapError(err.Error())
	}
	setBaseJson(this, r)
}

// @router /api/registry/delete [*]
func (this *BaseController) BaseDelete() {
	searchMap := sql.GetSearchMap("BaseId", *this.Ctx)
	baseData := registry.CloudImageBase{}

	q := sql.SearchSql(
		baseData,
		registry.SelectCloudImageBase,
		searchMap)
	sql.Raw(q).QueryRow(&baseData)

	q = sql.DeleteSql(registry.DeleteCloudImageBase, searchMap)
	r, err := sql.Raw(q).Exec()
	data := util.DeleteResponse(err,
		*this.Ctx, "删除基础镜像服务,名称:"+baseData.RegistryServer+baseData.ImageName,
		this.GetSession("username"),
		baseData.RegistryServer+baseData.ImageName, r)
	setBaseJson(this, data)
}

func setBaseJson(this *BaseController, data interface{}) {
	this.Data["json"] = data
	this.ServeJSON(false)
}
