package util

import (
	"github.com/astaxie/beego/context"
	"cloud/models/cloudLog"
	"strings"
	"cloud/sql"
	gosql "database/sql"
	"strconv"
	"io/ioutil"
	"path/filepath"
	"os"
	"fmt"
	"net/http"
	"bytes"
	"github.com/astaxie/beego/logs"
	"gopkg.in/square/go-jose.v1/json"
)



// 获取用户名
func GetUser(user interface{}) string {
	if user != nil {
		return user.(string)
	}
	return ""
}

// 记录操作日志
func SaveOperLog(user interface{}, ctx context.Context, info string, clusterName string) {
	username := GetUser(user)
	data := cloudLog.CloudOperLog{}
	data.Ip = GetClientIp(ctx.Request)
	data.Time = GetDate()
	data.User = username
	data.Messages = info
	data.Cluster = clusterName
	o := sql.GetOrm()
	o.Raw(sql.InsertSql(data, cloudLog.InsertCloudOperLog)).Exec()
}

// 获取来源页面
func GetUri(ctx context.Context) string {
	uri := ctx.Request.RequestURI
	uris := strings.Split(uri, "/")
	rd := strings.Join(uris, "/")
	return rd
}

// 获取来源页面
func GetReferer(ctx context.Context) string {
	referer := ctx.Request.Referer()
	referers := strings.Split(referer, "?referer=")
	if len(referers) > 1 {
		return referers[1]
	}
	return "/index"
}


// 删除时返回的响应数据
func DeleteResponse(err error, ctx context.Context, info string, username interface{}, cluster string, r gosql.Result) map[string]interface{} {
	var data map[string]interface{}
	if err == nil {
		d, _ := r.RowsAffected()
		if d == 0 {
			data = ApiResponse(false, "删除失败: 数据不存在")
		} else {
			data = ApiResponse(true, "删除成功: 成功删除了"+strconv.FormatInt(d, 10)+"条"+info)
		}
	} else {
		data = ApiResponse(false, "删除失败:"+err.Error())
	}
	SaveOperLog(username, ctx, info, cluster)
	return data
}

// 保存成功提示信息
func SaveResponse(err error, errmsg string) (map[string]interface{}, string) {
	var msg string
	var status bool = false
	if err != nil {
		if strings.Contains(err.Error(), "Duplicate") {
			msg = "保存失败: " + errmsg
		} else {
			msg = "保存失败: " + err.Error() + " " + errmsg
		}
	} else {
		status = true
		msg = "保存成功"
	}
	return ApiResponse(status, msg), msg
}

// 将字符串转成int类型
// 2018-01-12 10:19
func StringToInt(v string) interface{} {
	r, err := strconv.Atoi(v)
	if err != nil {
		return 0
	}
	return r
}

// 获取namespace
// 2018-01-13 07:02
func Namespace(appname string, resource string) string {
	return appname + "--" + resource
}

// 检查某个slice有某个数据
// 2018-01-14 14:26
func ListExistsInt(arr []int, value int) bool {
	for _, v := range arr {
		if v == value {
			return true
		}
	}
	return false
}

// 检查某个slice有某个数据
// 2018-02-05 21:06
func ListExistsString(arr []string, value string) bool {
	for _, v := range arr {
		if v == value {
			return true
		}
	}
	return false
}


// 2018-01-21 07:39
// 读取文件内容
func ReadFile(filename string) string {
	r, err := ioutil.ReadFile(filename)
	if err != nil {
		return ""
	}
	return string(r)
}

// 获取认证服务器的配置文件路径
// 2018-01-21 10:29
func AuthServerConfigFile() string {
	pwd,_ := os.Getwd()
	cf := filepath.Join(pwd,"conf", "simple.yaml")
	return cf
}



// 生成htmlselect的内容
// 2018-01-26 10:57
func GetSelectOption(name string, value string, title string) string {
	return"<option title='"+title+"' value='"+value+"'>"+name+"</option>"
}

// 生成htmlselect的内容
// 2018-02-07
func GetSelectOptionName(name string) string {
	return"<option title='"+name+"' value='"+name+"'>"+name+"</option>"
}

// 2018-01-31 09:30
// 获取响应数据
func GetResponseResult(err error,draw interface{}, returnData interface{}, totle interface{}) map[string]interface{} {
	var r map[string]interface{}
	if err == nil {
		r = ResponseMap(returnData, totle, draw)
	} else {
		r = ResponseMapError(err.Error())
	}
	return r
}


// 2018-03-24 17:03
// 发送json请求
func HttpGetJson(data string, url string) map[string]interface{} {
	fmt.Println("URL:>", url)
	//var jsonStr, _ = json.Marshal(data)
	v := map[string]interface{}{}
	req, err := http.NewRequest("GET", url, bytes.NewBuffer([]byte(data)))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		logs.Error("GET请求异常", err)
		return v
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	json.Unmarshal(body, &v)
	return v
}
