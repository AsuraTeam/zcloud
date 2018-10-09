package es

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"strings"
	"cloud/util"
	"cloud/sql"
	"cloud/models/log"
)

// 获取es数据
func getEsBuckets(data map[string]interface{}) []string {
	result := make([]string, 0)
	if data != nil {
		hits := data["aggregations"]
		if hits != nil {
			hits1 := hits.(map[string]interface{})
			if hits1 != nil {
				r := hits1["1"]
				if r != nil {
					hitsData := r.(map[string]interface{})
					if hitsData != nil {
						logs.Info(hitsData)
						for k, v := range hitsData {
							if k != "buckets" || v == nil {
								continue
							}
							vs := v.([]interface{})
							for _, bv := range vs {
								if bv == nil {
									continue
								}
								kvd := bv.(map[string]interface{})
								if kvd != nil {
									result = append(result, kvd["key"].(string))
								}
							}
						}
					}
					return result
				}
			}
		}
	}
	return result
}

func getMess(r1 string) string {
	r1 = strings.Replace(r1, "@kibana-highlighted-field@", "<span style='color:red'>", -1)
	r1 = strings.Replace(r1, "@/kibana-highlighted-field@", "</span>", -1)
	return r1
}

func getIp(source map[string]interface{}) string {
	if source != nil {
		fields := source["fields"]
		if fields != nil {
			ip := fields.(map[string]interface{})
			if ip != nil {
				ipadd := ip["ip_address"]
				if ipadd != nil {
					return ipadd.(string)
				}
			}
		}
	}
	return ""
}

// 获取es数据
func getEsValues(data map[string]interface{}, query string) (int64, []string) {
	result := make([]string, 0)
	total := float64(0)
	if data != nil {
		hits := data["hits"]
		if hits != nil {
			totalMap := hits.(map[string]interface{})["total"]
			if nil != totalMap {
				total = totalMap.(float64)
			}
			hits1 := hits.(map[string]interface{})
			if hits1 != nil {
				r := hits1["hits"]
				if r != nil {
					hitsData := r.([]interface{})
					if hitsData != nil {
						for _, v := range hitsData {
							var ok = false
							vs := v.(map[string]interface{})
							if query != "*" {
								if vs["highlight"] != nil {
									source := vs["highlight"].(map[string]interface{})
									if source != nil {
										if source["message"] == nil {
											source := vs["_source"].(map[string]interface{})
											if source != nil {
												mess := source["message"]
												if mess != nil {
													result = append(result, getIp(source)+" "+getMess(mess.(string)))
												}
											}
											continue
										}
										mess := source["message"].([]interface{})
										if len(mess) > 0 {
											ok = true
											result = append(result, getIp(vs["_source"].(map[string]interface{}))+" "+getMess(mess[0].(string)))
										}
									}
								}
							} else {
								if (vs["_source"] != nil && !ok) || query == "*" {
									source := vs["_source"].(map[string]interface{})
									if source != nil {
										mess := source["message"]
										if mess != nil {
											result = append(result, getIp(source)+" "+getMess(mess.(string)))
										}
									}
								}
							}

						}
					}
					return int64(total), result
				}
			}
		}
	}
	return int64(total), result
}

// 获取数据源地址
func getDataSource(env string, cluster string)  string {
	searchMap := sql.SearchMap{}
	searchMap.Put("Ent", env)
	searchMap.Put("ClusterName", cluster)
	data := log.LogDataSource{}
	sql.GetOrm().Raw(log.SelectDataSource, env, env, cluster).QueryRow(&data)
	return data.Address
}

// 2018-03-13 19:26
// 获取监控数据es
func RequestEs(query, index, q, env string, cluster string) (int64, []string) {
	//data := map[string]interface{}{}
	//json.Unmarshal([]byte(query), &data)
	server := getDataSource(env, cluster)
	server = fmt.Sprintf("%s/%s/_search", server, index)
	logs.Info(server, query)
	result := util.HttpGetJson(query, server)
	r := make([]string, 0)
	total := int64(0)
	if q == "buckets" {
		r = getEsBuckets(result)
	} else {
		total, r = getEsValues(result, q)
	}
	return total, r
}
