package main

import (
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"fmt"
	"strings"
	"path/filepath"
	"os"
	"time"
)

func init() {
	orm.Debug = true
	orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.RegisterDataBase("default", "mysql", "user:pass@tcp(1.1.1.1:3306)/cloud?charset=utf8")
}

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

func getType(tp string) string {
	switch tp {
	case "varchr":
		return "string"
	case "int":
		return "int64"
	case "double":
		return "float64"
	case "text":
		return "string"
	case "INTEGER":
		return "int64"
	case "BIGINT":
		return "int64"
	case "date":
		return "string"
	case "TIMESTAMP":
		return "int64"
	default:
		return "string"
	}
}

func main() {
	var table = "cloud_ci_service"
	var packageName = "ci"
	var maps []orm.Params
	var keyMaps []orm.Params
	keySql := " desc " + table
	sql := "select * from " + table + " limit 1"
	commentSql := "select column_comment,DATA_TYPE from information_schema.columns where table_name='"+table+"' and column_name='COLUMN_NAME'"
	var column = ""
	var structDatas = ""
	//var primaryKey = ""
	o := orm.NewOrm()
	o.Raw(sql).Values(&maps)
	for _, e := range maps {
		for k := range e {
			column += k + ","
		}
	}
	for _, e1 := range maps {
		for k1 := range e1 {
			var maps2 []orm.Params
			o.Raw(strings.Replace(commentSql, "COLUMN_NAME",k1,-1)).Values(&maps2)
			for _, e := range maps2 {
				structDatas += "    //" + e["column_comment"].(string) + "\n"
				t := getType(e["DATA_TYPE"].(string))
				structDatas += "    " + StringToUpper(k1) + " " + t + "\n"
			}
		}
	}
	o.Raw(keySql).Values(&keyMaps)
	//for _,k := range keyMaps{
	//	if k["Key"] == "PRI" {
	//		primaryKey = k["Field"].(string)
	//	}
	//}
	path := "D:\\F\\code\\workspace\\zcloud\\src\\cloud\\models\\" + packageName + "\\"
	structName := StringToUpper(table)
	column = column[0:len(column)-1]
	baseSql   := "\nconst Select"+structName+" = \"select " + column + " from " + table +"\""
	updateSql := "const Update"+structName+" = \"update " + table +"\""
	deleteSql := "const Delete"+structName+" = \"delete from " + table +"\" "
	insertSql := "const Insert"+structName+" = \"insert into " + table +"\" "
	//findById  := "const FindById" +structName + " = Select"+structName+" + \" where {0}={1}\""
	//findById = strings.Replace(findById, "{0}", primaryKey, -1)
	fmt.Println(baseSql)
	fmt.Println(updateSql)
	fmt.Println(deleteSql)
	fmt.Println(insertSql)
	structData := "\n//" +time.Now().Local().String() + "\ntype " + structName + " struct {\n"+structDatas+"}"
	file := filepath.Join(path,"struct.go")
	_,err := os.Stat(file)
	f, _ := os.OpenFile(file, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil{
		f.Write([]byte("package " + packageName+"\n"))
	}
	f.Write([]byte(structData+"\n"))
	f.Close()
	fmt.Println(structData)
	fileMap := filepath.Join(path,"structMap.go")
	_,err = os.Stat(fileMap)
	f1, _ := os.OpenFile(fileMap, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil{
		f1.Write([]byte("package " + packageName+"\n"))
	}
	f1.Write([]byte(baseSql+"\n"))
	//f1.Write([]byte(findById+"\n"))
	f1.Write([]byte(updateSql+"\n"))
	f1.Write([]byte(insertSql+"\n"))
	f1.Write([]byte(deleteSql+"\n"))
	f.Close()
}
