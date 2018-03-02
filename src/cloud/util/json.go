package util

import (
	"strconv"
	"fmt"
	"strings"
	"encoding/json"
	"github.com/ghodss/yaml"
	goyaml "gopkg.in/yaml.v2"
)

func transformData(pIn *interface{}) (err error) {
	switch in := (*pIn).(type) {
	case map[interface{}]interface{}:
		m := make(map[string]interface{}, len(in))
		for k, v := range in {
			if err = transformData(&v); err != nil {
				return err
			}
			var sk string
			switch k.(type) {
			case string:
				sk = k.(string)
			case int:
				sk = strconv.Itoa(k.(int))
			default:
				return fmt.Errorf("type mismatch: expect map key string or int; got: %T", k)
			}
			m[sk] = v
		}
		*pIn = m
	case []interface{}:
		for i := len(in) - 1; i >= 0; i-- {
			if err = transformData(&in[i]); err != nil {
				return err
			}
		}
	}
	return nil
}

// 将yaml配置转换成json
func Yaml2Json(input []byte) (map[string]interface{}, error) {
	temp := make(map[string]interface{})
	var data interface{}
	err := goyaml.Unmarshal(input, &data)
	if err != nil {
		return temp, err
	}
	input = nil
	err = transformData(&data)
	if err != nil {
		return temp, err
	}
	output, err := json.Marshal(data)
	fmt.Println(string(output))
	if err != nil {
		return temp, err
	}
	data = nil
	json.Unmarshal(output, &temp)
	return temp, err
}

// 将json转换成yaml格式
// 2018-01-11 22:10
func Json2Yaml(jsondata string) string {
	data, err := yaml.JSONToYAML([]byte(jsondata))
	if err == nil {
		yaml := string(data)
		yaml = strings.Replace(yaml, "- apiVersion:", "---\napiVersion: ", -1)
		return strings.Join(strings.Split(yaml, "\n")[1:], "\n")
	}
	return err.Error()
}

// 2018-02-04
// 将一些公共数据写好
func SetPublicData(src interface{}, username string, obj interface{})  {
	data := Lock{}
	srcData := make(map[string]interface{})
	temp,_ := json.Marshal(src)
	json.Unmarshal(temp, &data)
	for k, v := range srcData {
		data.Put(k, v)
	}
	data.Put("CreateTime", GetDate())
	data.Put("CreateUser",  username)
	data.Put("LastModifyUser",  username)
	data.Put("LastModifyTime", data.GetV("CreateTime"))
	temp, _ = json.Marshal(data.GetData())
	json.Unmarshal(temp, &obj)
}

// 2018-02-14 10:22
// 将2个struct数据合并
func MergerStruct(a interface{}, b interface{}) {
	data := Lock{}
	temp,_ := json.Marshal(a)
	srcData := make(map[string]interface{})
	json.Unmarshal(temp, &srcData)
	for k, v := range srcData {
		data.Put(k, v)
	}
	destData := make(map[string]interface{})
	temp,_ = json.Marshal(b)
	json.Unmarshal(temp, &destData)
	for k, v := range destData {
		if _, ok := data.Get(k) ; !ok {
			data.Put(k, v)
		}
	}
	temp, _ = json.Marshal(data.GetData())
	json.Unmarshal(temp, &b)
}
