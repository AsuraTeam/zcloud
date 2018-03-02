package operlog

import (
	"cloud/sql"
	"cloud/util"
	"github.com/astaxie/beego"
	"cloud/models/cloudLog"
)

type LogController struct {
	beego.Controller
}

// 日志管理入口页面
// @router /users/user/list [get]
func (this *LogController) OperlogList() {
	this.TplName = "operlog/list.html"
}


// 用户数据
// @router /api/operlog [get]
func (this *LogController) OperlogDatas() {
	data := []cloudLog.CloudOperLog{}
	user := util.GetUser(this.GetSession("username"))
	searchMap := sql.GetSearchMapV("User", user)
	key := this.GetString("search")
	searchSql := sql.SearchSql(cloudLog.CloudOperLog{}, cloudLog.SelectCloudOperLog, searchMap)
	if len(key) > 4 {
		key = sql.Replace(key)
		searchSql += "  and (user like \"%" + key + "%\" or messages like \"%" + key + "%\" or ip like \"%" + key + "%\")"
	}
	num,_ := sql.OrderByPagingSql(searchSql, "time", *this.Ctx.Request, &data, cloudLog.CloudOperLog{})
	r := util.ResponseMap(data, sql.CountSearchMap("cloud_oper_log", sql.GetSearchMapV("User", user), int(num), key), this.GetString("draw"))
	setUserJson(this, r)

}

func setUserJson(this *LogController, data interface{}) {
	this.Data["json"] = data
	this.ServeJSON(false)
}
