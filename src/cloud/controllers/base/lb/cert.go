package lb

import (
	"github.com/astaxie/beego"
	"cloud/sql"
	"cloud/util"
	"strings"
	"cloud/models/lb"
	"cloud/k8s"
)

// 2018-02-02 10:32
// 持续集成
type CertController struct {
	beego.Controller
}

// cert管理入口页面
// @router /ci/cert/list [get]
func (this *CertController) CertList() {
	this.TplName = "base/network/cert/list.html"
}


// 生成 镜像服务 html
// 2018-01-26 10:41
func GetCertSelect() string {
	html := make([]string, 0)
	data := GetCertfileData("")
	for _,v := range data{
		html = append(html, util.GetSelectOption(v.CertKey, v.CertKey, v.CertKey))
	}
	return strings.Join(html, "")
}


// cert管理添加页面
// @router /ci/cert/add [get]
func (this *CertController) CertAdd() {
	id := this.GetString("CertId")
	update := k8s.CloudLbCert{}

	// 更新操作
	if id != "0" {
		searchMap := sql.GetSearchMap("CertId", *this.Ctx)
		sql.Raw(sql.SearchSql(k8s.CloudLbCert{}, k8s.SelectCloudLbCert, searchMap)).QueryRow(&update)
	}

	this.Data["data"] = update
	this.TplName = "base/network/cert/add.html"
}

// 获取docker数据
// 2018-01-26 11:17
func GetCertfileData(name string)[]k8s.CloudLbCert {
	searchMap := sql.SearchMap{}
	if name != "" {
		searchMap.Put("Name", name)
	}
	// cert数据
	data := make([]k8s.CloudLbCert, 0)
	q := sql.SearchSql(k8s.CloudLbCert{}, k8s.SelectCloudLbCert, searchMap)
	sql.Raw(q).QueryRows(&data)
	return data
}



// string
// cert保存
// @router /api/network/cert [post]
func (this *CertController) CertSave() {
	d := k8s.CloudLbCert{}
	err := this.ParseForm(&d)
	if err != nil {
		this.Ctx.WriteString("参数错误" + err.Error())
		return
	}
	util.SetPublicData(d, util.GetUser(this.GetSession("username")), &d)


	q := sql.InsertSql(d, lb.InsertCloudLbCert)
	if d.CertId > 0 {
		searchMap := sql.SearchMap{}
		searchMap.Put("CertId", d.CertId)
		q = sql.UpdateSql(d, lb.UpdateCloudLbCert, searchMap, "CreateTime,CreateUser")
	}
	sql.Raw(q).Exec()
	data,msg := util.SaveResponse(err, "名称已经被使用")
	util.SaveOperLog(this.GetSession("username"), *this.Ctx, "保存证书配置 "+msg, d.CertKey)
	setCertJson(this, data)
	go k8s.CreateNginxConf("")
}




// cert数据
// @router /api/network/cert [get]
func (this *CertController) CertData() {
	data := make([]k8s.CloudLbCert, 0)
	searchMap := sql.SearchMap{}
	id := this.Ctx.Input.Param(":id")
	key := this.GetString("search")
	if id != "" {
		searchMap.Put("CertId", id)
	}
	searchSql := sql.SearchSql(k8s.CloudLbCert{}, k8s.SelectCloudLbCert, searchMap)
	if key != "" && id == "" {
		key = sql.Replace(key)
		searchSql += " where 1=1 and (cert_key like \"%" + key + "%\" or description like \"%" + key + "%\" )"
	}
	num, err := sql.Raw(searchSql).QueryRows(&data)
	r := util.GetResponseResult(err, this.GetString("draw"), data, sql.Count("cloud_lb_cert", int(num), key))
	setCertJson(this, r)
}

// json
// 删除cert
// 2018-02-02 21:46
// @router /api/network/cert/:id:int [delete]
func (this *CertController) CertDelete() {
	searchMap := sql.SearchMap{}
	id := this.Ctx.Input.Param(":id")
	searchMap.Put("CertId", id)
	codeData := k8s.CloudLbCert{}
	sql.Raw(sql.SearchSql(codeData, k8s.SelectCloudLbCert, searchMap)).QueryRow(&codeData)
	r, err := sql.Raw(sql.DeleteSql(lb.DeleteCloudLbCert, searchMap)).Exec()
	data := util.DeleteResponse(err, *this.Ctx, "删除cert"+codeData.CertKey, this.GetSession("username"), codeData.CreateUser, r)
	setCertJson(this, data)
}

func setCertJson(this *CertController ,data interface{})  {
	this.Data["json"] = data
	this.ServeJSON(false)
}