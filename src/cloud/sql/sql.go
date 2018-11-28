package sql

import (
	_ "github.com/go-sql-driver/mysql"
	"encoding/json"
	"strconv"
	"strings"
	"sync"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/context"
	"net/http"
	"github.com/astaxie/beego/logs"
	"reflect"
	"database/sql"
)

func StringToUpper(str string) string {
	strs := strings.Split(str, "_")
	temp := ""
	for _, s := range strs {
		ss := strings.Split(s, "")
		temp += strings.ToUpper(ss[0])
		temp += s[1:]
		continue
	}
	return temp
}

// 判断某个元素是否在splice
func IsExists(strs []string, key string) bool {
	for _, k := range strs {
		if k == key {
			return true
		}
	}
	return false
}

func StringToLower(str string) string {
	strs := strings.Split(str, "")
	count := 0
	var temp = ""
	for _, s := range strs {
		if count == 0 {
			s = strings.ToLower(s)
			temp += s
			count += 1
			continue
		}
		if s == strings.ToUpper(s) {
			s = "_" + strings.ToLower(s)
		}
		temp += s
		count += 1
	}
	return temp
}

func replace(sql string) string {
	sql = strings.Replace(sql, "\"", "\\\"", -1)
	return sql
}

func Replace(sql string) string {
	return replace(sql)
}

type SearchMap struct {
	Lock sync.RWMutex
	Data map[string]interface{}
}

func (m *SearchMap) GetData() map[string]interface{} {
	m.Lock.RLock()
	defer m.Lock.RUnlock()
	return m.Data
}

// 获取数据
func (m *SearchMap) Get(key string) (interface{}) {
	m.Lock.RLock()
	defer m.Lock.RUnlock()
	data := m.GetData()
	if _, ok := data[key]; ok {
		return data[key]
	}
	return nil
}

func (m *SearchMap) Put(k string, v interface{}) {
	if len(m.Data) < 1 {
		m.Data = make(map[string]interface{})
	}
	m.Lock.Lock()
	m.Data[k] = v
	defer m.Lock.Unlock()
}

func getMap(obj interface{}) map[string]interface{} {
	var maps = make(map[string]interface{})
	d, _ := json.Marshal(obj)
	err := json.Unmarshal(d, &maps)
	if err != nil {
		logs.Info("getMap", err)
	}
	return maps
}

func getReturnSql(tempSql string) string {
	tempSql = strings.TrimSpace(tempSql)
	temps := strings.Split(tempSql, " ")
	v := temps[0:len(temps)-1]
	return strings.Join(v, " ")
}

func getValue(v interface{}) string {
	switch v.(type) {
	case int:
		return strconv.Itoa(v.(int))
	case int64:
		return strconv.FormatInt(v.(int64), 10)
	case float64:
		return strconv.FormatFloat(v.(float64), 'g', -1, 64)
	case string:
		return "\"" + replace(v.(string)) + "\" "
	case bool:
		return "\"" + replace(strconv.FormatBool(v.(bool))) + "\" "
	default:
		return "\"" + replace(v.(string)) + "\" "
	}
}

func getSql(searchMap map[string]interface{}, k string, tempSql string, connect string) string {
	v := searchMap[k]
	if v == nil {
		return tempSql
	}
	k = StringToLower(k)
	switch v.(type) {
	case int:
		tempSql += k + "=" + strconv.Itoa(v.(int)) + " " + connect + " "
		break
	case int64:
		tempSql += k + "=" + strconv.FormatInt(v.(int64), 10) + " " + connect + " "
		break
	case float64:
		tempSql += k + "=" + strconv.FormatFloat(v.(float64), 'g', -1, 64) + " " + connect + " "
		break
	case string:
		tempSql += k + "=\"" + replace(v.(string)) + "\" " + connect + " "
		break
	default:
		tempSql += k + "=\"" + replace(v.(string)) + "\" " + connect + " "
		break
	}
	return tempSql
}

//searchMap := util.SearchMap{}
//searchMap.Put("CreateUser", "zhaozq14")
//searchMap.Put("AppId", 89)
//data := index.DockerCloudApp{}
//data.AppName = "1111"
//data.Status = "1"
//data.AppId = 10
//data.AppDescription = "测试update"
// extCloumnt 排除的列,多个逗号隔开
//updateSql := util.UpdateSql(data, index.UpdateDockerCloudApp, searchMap)
func UpdateSql(obj interface{}, sql string, data SearchMap, extCloumnt string) string {
	exts := strings.Split(extCloumnt, ",")
	searchMap := data.GetData()
	var tempSql = " "
	updateSql := " "
	whereSql := " where "
	var maps = getMap(obj)
	for k := range maps {
		if IsExists(exts, StringToUpper(k)) || IsExists(exts, StringToLower(k)) {
			continue
		}
		if maps[k] != nil {
			updateSql += getSql(maps, k, "", ",")
		}
		if _, ok := searchMap[k]; ok {
			whereSql = getSql(searchMap, k, whereSql, "and")
		}
	}
	updateSql = strings.TrimSpace(updateSql)
	updateSql = updateSql[0:len(updateSql)-1]
	tempSql = sql + " set " + updateSql + " " + whereSql
	return getReturnSql(tempSql)
}

func GetWhere(searchSql string, searchMap SearchMap) string  {
	if len(searchMap.GetData()) == 0 {
		searchSql += " where 1=1 "
	}
	return  searchSql
}

func getSearchSql(obj interface{}, sql string, data SearchMap) string {
	searchMap := data.GetData()
	sql += " where 1=1 and "
	var tempSql = sql
	var maps = getMap(obj)
	for k := range maps {
		if _, ok := searchMap[k]; ok {
			tempSql = getSql(searchMap, k, tempSql, "and")
		}
	}
	tempSql = strings.Replace(tempSql, "1=1 and", "", -1)
	return getReturnSql(tempSql)
}

//searchMap := util.SearchMap{}
//searchMap.Put("CreateUser", "zhaozq14")
//sql := util.SearchSql(index.DockerCloudApp{},index.SelectDockerCloudApp, searchMap)
func SearchSql(obj interface{}, sql string, data SearchMap) string {
	return getSearchSql(obj, sql, data)
}

// 2018-02-01 09:36
func getCount(table string, searchMap SearchMap) int {
	maps := []orm.Params{}
	countSql := "select count(1) as cnt from " + table
	var searchQ = make([]string, 0)
	for k, v := range searchMap.GetData() {
		searchQ = append(searchQ, StringToLower(k)+"=\""+replace(v.(string))+"\"")
	}
	if len(searchQ) > 0 {
		countSql += " where 1=1 and " + strings.Join(searchQ, " and ")
	}
	GetOrm().Raw(countSql).Values(&maps)
	if len(maps) > 0 {
		v, err := strconv.Atoi(maps[0]["cnt"].(string))
		if err == nil {
			return v
		}
	}
	return 0
}

type Total struct {
	Total int64
}

/**
 2018-11-27 09:09
 获取查询sql的总行数
 */
func CountSqlTotal(query string) int64  {
	qs := strings.Split(query, " from ")
	total :=  Total{}
	if len(qs) > 1 {
		q := strings.Split(query, " limit ")[0]
		query := "select count(*) as total from (" + q + ") as temp"
		GetOrm().Raw(query).QueryRow(&total)
	}
	return total.Total
}

// 2018-01-17 12:40
// 获取表的行数
func Count(table string, count int, search string) int {
	if search != "" {
		return count
	}
	return getCount(table, SearchMap{})
}

// 获取有条件的计算
// 2018-02-01 9:38
func CountSearchMap(table string, searchMap SearchMap, num int, search string) int {
	if search != "" {
		return num
	}
	return getCount(table, searchMap)
}

// 带分页的sql语句
// 2018-01-15
func SearchSqlPages(sql string, request http.Request) string {
	length := getParam(&request, "length")
	start := getParam(&request, "start")
	if start == "" {
		start = "1"
	}
	if length == "" {
		length = "10"
	}
	pstart, serr := strconv.ParseInt(start, 10, 64)
	plength, lerr := strconv.ParseInt(length, 10, 64)
	if lerr != nil || serr != nil {
		pstart = 0
		plength = 10
	}

	l := strconv.FormatInt(plength, 10)
	st := (pstart - 1)
	if st < 0 {
		st = 0
	}
	s := strconv.FormatInt(st, 10)
	sql += " limit " + s + "," + l
	return sql
}

// 添加排序
// 2018-01-13 07:08
func SearchOrder(sql string, columnt ...string) string {
	sql += " order by " + strings.Join(columnt, ",") + " desc"
	return sql
}

// 2018-01-20 9:30
// 将对象转成字符串
func ObjToString(v interface{}) string {
	t, _ := json.Marshal(v)
	return string(t)
}

// 获取某个对象的类型
// 2018-02-05 10;38
func getObjStructMap(structObj interface{}) map[string]string {
	var result = make(map[string]string)
	var v = reflect.ValueOf(structObj)
	t := v.Type()
	if v.Kind() == reflect.Struct { //反射结构体成员信息
		for i := 0; i < t.NumField(); i++ {
			f := t.Field(i)
			result[f.Name] = f.Type.String()
		}
	}
	return result
}

// 2018-02-05 09:09
// 公共查询方法，包含orderby和分页sql
func OrderByPagingSql(searchSql string, column string, request http.Request, obj interface{}, structObj interface{}) (int, error) {
	searchSql = SearchOrder(searchSql, column)
	searchSql = SearchSqlPages(searchSql, request)
	c := []orm.Params{}
	_, err := GetOrm().Raw(searchSql).Values(&c)
	structMap := getObjStructMap(structObj)

	var objData = make([]map[string]interface{}, 0)
	for _, v := range c {
		result := make(map[string]interface{}, 0)
		for mk, v := range v {
			key := StringToUpper(mk)
			if structMap[key] == "string" {
				result[StringToUpper(mk)] = v
			}
			if structMap[key] == "int64" {
				if v == nil {
					result[StringToUpper(mk)] = 0
					continue
				}
				vint64, err := strconv.ParseInt(v.(string), 10, 64)
				if err == nil {
					result[StringToUpper(mk)] = vint64
				} else {
					result[StringToUpper(mk)] = 0
				}
			}
			if structMap[key] == "float64" {
				if v == nil {
					result[StringToUpper(mk)] = 0.0
					continue
				}
				vFloat64, err := strconv.ParseFloat(v.(string),  64)
				if err == nil {
					result[StringToUpper(mk)] = vFloat64
				} else {
					result[StringToUpper(mk)] = 0.0
				}
			}
			if structMap[key] == "int" || structMap[key] == "int32" {
				if v == nil {
					result[StringToUpper(mk)] = 0
					continue
				}
				vint, err := strconv.Atoi(v.(string))
				if err == nil {
					result[StringToUpper(mk)] = vint
				} else {
					result[StringToUpper(mk)] = 0
				}
			}
		}
		objData = append(objData, result)
	}
	temp, _ := json.Marshal(objData)
	err = json.Unmarshal(temp, &obj)
	if err != nil {
		logs.Info(err)
	}
	return len(objData), err
}

//insertSql := util.InsertSql(data,"insert into docker_application")
func InsertSql(obj interface{}, sql string) string {
	var tempSql = sql
	var insertSql = "("
	var valueSql = ""
	var maps = getMap(obj)
	for k := range maps {
		if maps[k] != "" {
			insertSql += StringToLower(k) + ","
			valueSql += getValue(maps[k]) + ","
		}
	}
	insertSql = insertSql[0:len(insertSql)-1]
	valueSql = valueSql[0:len(valueSql)-1]
	tempSql += insertSql + ") values(" + valueSql + ")"
	return tempSql
}

func FindById(sql string, id int) string {
	sql = strings.Replace(sql, "{1}", strconv.Itoa(id), -1)
	return sql
}

//deleteSql := util.DeleteSql("delete from docker_application", searchMap)
func DeleteSql(sql string, data SearchMap) string {
	var tempSql = sql
	searchMap := data.GetData()
	if len(searchMap) > 0 {
		tempSql += " where "
	} else {
		return ""
	}
	for k := range searchMap {
		if _, ok := searchMap[k]; ok {
			tempSql = getSql(searchMap, k, tempSql, "and")
		}
	}
	return getReturnSql(tempSql)
}

//searchMap := sql.SearchMap{}
//searchMap.Put("CreateUser", "zhaozq14")
//searchMap.Put("AppId", 89)
//data := index.DockerCloudApp{}
//data.AppName = "1111"
//data.Status = "1"
//data.AppId = 10
//data.AppDescription = "测试update"
//updateSql := sql.UpdateSql(data, index.UpdateDockerCloudApp, searchMap)
//deleteSql := sql.DeleteSql("delete from docker_application", searchMap)
//insertSql := sql.InsertSql(data,"insert into docker_application")
//fmt.Println(updateSql)
//fmt.Println(deleteSql)
//fmt.Println(insertSql)

func init() {
	orm.Debug = true
	orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.RegisterDataBase("default", "mysql", beego.AppConfig.String("mysql"))
}

func GetOrm() orm.Ormer {
	var o = orm.NewOrm()
	return o
}

// 插入数据
func Insert(obj interface{}) (int64, error) {
	return GetOrm().Insert(obj)
}

// 执行sql语句
func Raw(q string) orm.RawSeter {
	o := GetOrm()
	return o.Raw(q)
}

// 执行sql语句
func Exec(q string) (sql.Result, error) {
	o := GetOrm()
	return o.Raw(q).Exec()
}

func getParam(req *http.Request, key string) string {
	req.ParseForm()
	var id = ""
	if len(req.Form[key]) > 0 {
		id = req.Form[key][0]
	}
	return id
}

// 2018-02-09 06:40
// 将要获取的参数都写到一起
func MKeyV(key ...string) []string {
	return key
}

// 2018-02-09 07:57
// 将参数值写到SearchMap
func GetSearchMapValue(key []string, ctx context.Context, searchMap SearchMap) SearchMap {
	for _, paramKey := range key {
		paramValue := getParam(ctx.Request, paramKey)
		if paramValue != "" {
			searchMap.Put(paramKey, paramValue)
		}
	}
	return searchMap
}

// 创建一个带参数ID的 SearchMap
func GetSearchMap(key string, ctx context.Context) SearchMap {
	searchMap := SearchMap{}
	id := ctx.Input.Param(":id")
	if id == "" {
		id = getParam(ctx.Request, key)
	}
	if id != "" {
		searchMap.Put(key, id)
	}
	str := ctx.Input.Param(":hi")
	if str != "" {
		searchMap.Put(key, str)
	}
	return searchMap
}

// 获取默认带数据的searchmap
// 2018-01-12 14:05
func GetSearchMapV(key ...string) SearchMap {
	searchMap := SearchMap{}
	counter := 0
	var k string
	for _, s := range key {
		if counter == 0 {
			k = s
		}
		if counter == 1 {
			v := s
			searchMap.Put(k, v)
			counter = 0
		} else {
			counter += 1
		}
	}
	return searchMap
}

// 获取参数并写入到map里面
// 2018-01-13 19:36
func GetString(ctx context.Context, key ...string) map[string]string {
	data := make(map[string]string)
	for _, k := range key {
		v := getParam(ctx.Request, k)
		if v != "" {
			data[k] = v
		}
	}
	return data
}
