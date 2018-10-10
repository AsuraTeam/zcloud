package log

import (
	"github.com/astaxie/beego"
	"strings"
	"strconv"
	"fmt"
	"gopkg.in/square/go-jose.v1/json"
	"cloud/sql"
	"cloud/util"
	"cloud/models/log"
	"cloud/es"
	"cloud/controllers/ent"
	app2 "cloud/controllers/docker/application/app"
)

type ControllerLog struct {
	beego.Controller
}

// 2018-03-01 13:20
// 日志过滤器页面
// @router /log/filter [get]
func (this *ControllerLog) FilterList() {
	this.TplName = "log/filter.html"
}

func (this *ControllerLog) Index() {
	o := log.LogShowFilter{}
	q := sql.SearchSql(o, log.SelectLastSearch, sql.SearchMap{})
	q = strings.Replace(q, "{0}", util.GetUser(this.GetSession("username")), -1)
	sql.Raw(q).QueryRow(&o)
	this.Data["data"] = o
	this.Data["ent"] = ent.GetEntnameSelectData(true)
	env := this.GetString("env")
	app := this.GetString("app")
	cluster := this.GetString("cluster")
	if len(env) > 0 && len(app) > 0  && len(cluster) > 0 {
		this.Data["selectEnt"] = app2.GetEntDescription(env)
		this.Data["selectApp"] = app
		this.Data["selectCluster"] = cluster
	}
	this.TplName = "log/log.html"
}

// 2018-03-01 13:22
// 日志过滤器页面
// @router /log/query/filter/:id:int [get]
func (this *ControllerLog) QueryFilter() {
	tp := this.GetString("type")
	searchMap := sql.GetSearchMap("Id", *this.Ctx)
	o := log.LogShowFilter{}
	var q string
	if tp == "history" {
		q = sql.SearchSql(o, log.SelectLogShowHistory, searchMap)
	} else {
		q = sql.SearchSql(o, log.SelectLogShowFilter, searchMap)
	}
	sql.Raw(q).QueryRow(&o)
	this.Data["data"] = o
	this.TplName = "index/index.html"
}

// 日志历史页面
// @router /log/history [get]
func (this *ControllerLog) HistoryList() {
	this.TplName = "log/history.html"
}


// 2018-09-14 14:23
// 日志条件搜索
// @router /api/log/filter [post]
func (this *ControllerLog) SaveFilter() {
	tp := this.GetString("type")
	env := this.GetString("env")
	app := this.GetString("appname")
	ip := this.GetString("ip")
	hostname := this.GetString("hostname")
	query := this.GetString("query")
	user := util.GetUser(this.GetSession("username"))
	if len(env) == 0 && len(app) == 0 && len(ip) == 0 && len(hostname) == 0 {
		setLogShowFilterJson(this, util.ApiResponse(false, ""))
		return
	}
	obj := log.LogShowFilter{
		Env:        env,
		Ip:         ip,
		Hostname:   hostname,
		Query:      query,
		Appname:    app,
		CreateTime: util.GetDate(),
		CreateUser: user,
	}
	if tp == "history" {
		r, _ := json.Marshal(obj)
		o := log.LogShowHistory{}
		json.Unmarshal(r, &o)
		sql.Exec(sql.InsertSql(o, log.InsertLogShowHistory))
		q := fmt.Sprintf(`update log_show_filter set click=click+1 where create_user='%v' and appname='%v' and hostname='%v' and ip='%v' and env='%v' and query='%v'`, user, obj.Appname, obj.Hostname, obj.Ip, obj.Env, obj.Query)
		sql.Exec(q)
	} else {
		sql.Exec(sql.InsertSql(obj, log.InsertLogShowFilter))
	}
	setLogShowFilterJson(this, util.ApiResponse(true, "保存成功"))
}

// 2018-09-14 09:21
// 日志条件搜索
// @router /api/log/query [get]
func (this *ControllerLog) Query() {
	env := this.GetString("env")
	cluster := this.GetString("cluster")
	app := this.GetString("appname")
	query := this.GetString("query")
	start := this.GetString("start")
	end := this.GetString("end")
	size, _ := this.GetInt64("size", 500)

	if len(env) == 0 {
		env = "prod"
	}

	selectQuery := make([]string, 0)
	if len(env) > 0 {
		selectQuery = append(selectQuery, "fields.runtime_env: \\\""+env+`\"`)
	}
	if len(app) > 0 {
		selectQuery = append(selectQuery, "fields.appname: \\\""+app+`\"`)
	}
	if len(query) == 0 {
		query = "*"
	}

	if len(selectQuery) > 0 {
		if !strings.Contains(query, " ") && query != "*"  && len(query) > 0 {
			query = fmt.Sprintf("*%v*", query)
		}
		query = fmt.Sprintf("%v AND (%v) ", query, strings.Join(selectQuery, " AND "))
	}

	q := `{
    "_source": {
        "excludes": []
    },
    "docvalue_fields": [
        "@timestamp"
    ],
    "highlight": {
        "fields": {
            "*": {
                "highlight_query": {
                    "bool": {
                        "must": [
                            {
                                "query_string": {
                                    "all_fields": true,
                                    "analyze_wildcard": true,
                                    "query": "` + query + `"
                                }
                            },
                            {
                                "range": {
                                    "@timestamp": {
                                        "format": "epoch_millis",
                                           "gte": ` + start + `,
                                           "lte": ` + end + `
                                    }
                                }
                            }
                        ],
                        "must_not": []
                    }
                }
            }
        },
        "fragment_size": 2147483647,
        "post_tags": [
            "@/kibana-highlighted-field@"
        ],
        "pre_tags": [
            "@kibana-highlighted-field@"
        ]
    },
    "query": {
        "bool": {
            "must": [
                {
                    "query_string": {
                        "analyze_wildcard": true,
                        "query": "` + query + `"
                    }
                },
                {
                    "range": {
                        "@timestamp": {
                            "format": "epoch_millis",
                              "gte": ` + start + `,
                              "lte": ` + end + `
                        }
                    }
                }
            ],
            "must_not": []
        }
    },
    "script_fields": {},
  "size": ` + strconv.FormatInt(size, 10) + `,
    "sort": [
        {
            "@timestamp": {
                "order": "desc",
                "unmapped_type": "boolean"
            }
        }
    ],
    "stored_fields": [
        "*"
    ],
    "version": true
}`
	total, r := es.RequestEs(q, "*", query, env, cluster)
	r1 := strings.Join(r, "<br>")
	this.Ctx.ResponseWriter.Header().Add("Content-Type", "text/html; charset=utf-8")
	this.Ctx.WriteString(fmt.Sprintf("共匹配到<span class='text-danger'>%d</span>", total) +"行,最多显示500行<br>" +  r1)
}

// 2018-09-14 09:21
// 日志条件搜索
// @router /api/log/search [get]
func (this *ControllerLog) Search() {
	key := this.GetString("key")
	tp := this.GetString("type")
	appname := this.GetString("appname")
	q := ""
	switch tp {
	case "ip":
		q = "select ip as value from log_show_ip where ip like '%key%'"
		if len(appname) > 0 {
			q += fmt.Sprintf(" and app_name='%v'", appname)
		}
		break
	case "appname":
		q = "select appname as value from log_show_appname where appname like '%key%'"
		break
	case "hostname":
		q = "select hostname as value from log_show_hostname where hostname like '%key%'"
		break
	}
	q = strings.Replace(q, "key", sql.Replace(key), -1)
	data := make([]log.Search, 0)
	sql.Raw(q).QueryRows(&data)
	setLogShowFilterJson(this, data)
}

// 日志过滤器数据
// @router /api/log/filter [get]
func (this *ControllerLog) LogShowFilterData() {
	data := make([]log.LogShowFilter, 0)
	searchMap := sql.SearchMap{}
	key := this.GetString("search")
	all := this.GetString("all")
	searchSql := sql.SearchSql(
		log.LogShowFilter{},
		log.SelectLogShowFilter,
		searchMap)

	searchSql += " where 1=1 "

	if len(key) > 0 {
		searchSql += ` and query like '%V%' or ip like '%V%' or appname like '%V%' or hostname like '%V%'`
		searchSql = strings.Replace(searchSql, "V", sql.Replace(key), -1)
	}

	if all != "1" {
		searchSql += ` and create_user='` + util.GetUser(this.GetSession("username")) + "'"
	}

	num, _ := sql.OrderByPagingSql(
		searchSql, "id",
		*this.Ctx.Request,
		&data,
		log.LogShowFilter{})

	r := util.ResponseMap(
		data,
		sql.Count("log_show_filter", int(num), key),
		this.GetString("draw"))
	setLogShowFilterJson(this, r)

}

// 日志历史
// @router /api/log/history [get]
func (this *ControllerLog) LogShowHistoryData() {
	data := make([]log.LogShowHistory, 0)
	searchMap := sql.SearchMap{}
	key := this.GetString("search")
	searchSql := sql.SearchSql(
		log.LogShowHistory{},
		log.SelectLogShowHistory,
		searchMap)
	all := this.GetString("all")
	searchSql += " where 1=1 "
	if len(key) > 0 {
		searchSql += ` and query like '%V%' or ip like '%V%' or appname like '%V%' or hostname like '%V%'`
		searchSql = strings.Replace(searchSql, "V", sql.Replace(key), -1)
	}
	if all != "1" {
		searchSql += ` and create_user='` + util.GetUser(this.GetSession("username")) + "'"
	}
	num, _ := sql.OrderByPagingSql(
		searchSql, "id",
		*this.Ctx.Request,
		&data,
		log.LogShowHistory{})

	r := util.ResponseMap(
		data,
		sql.Count("log_show_history", int(num), key),
		this.GetString("draw"))
	setLogShowFilterJson(this, r)

}

// 删除日志过滤器
// @router /api/log/filter/:id:int [delete]
func (this *ControllerLog) LogShowFilterDelete() {
	searchMap := sql.GetSearchMap("Id", *this.Ctx)
	searchMap.Put("CreateUser", util.GetUser(this.GetSession("username")))
	q := sql.DeleteSql(log.DeleteLogShowFilter, searchMap)
	sql.Exec(q)
	setLogShowFilterJson(this, util.ApiResponse(true, "删除成功"))
}

func setLogShowFilterJson(this *ControllerLog, data interface{}) {
	this.Data["json"] = data
	this.ServeJSON(false)
}
