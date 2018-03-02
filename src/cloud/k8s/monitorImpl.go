package k8s

import (
	"strings"
	"net/http"
	"github.com/astaxie/beego/logs"
	"io/ioutil"
	"net/url"
	"encoding/json"
	"strconv"
	"time"
	"cloud/util"
)

// sum(rate(container_cpu_usage_seconds_total{id="/",instance=~"^.*$"}[1m])) / sum (machine_cpu_cores{instance=~"^.*$"}) * 100 cpu使用率
// sum (container_memory_working_set_bytes{id="/",instance=~"^.*$"}) / sum (machine_memory_bytes{instance=~"^.*$"}) * 100 内存使用率
// query:sum(rate(container_cpu_usage_seconds_total{image!="",name=~"^k8s_.*",instance=~"^.*$",namespace=~"^kube-system$"}[1m])) by (pod_name)
type QueryParam struct {
	// 集群名称
	ClusterName string
	// 服务名称
	ServiceName string
	// 应用名称
	AppName string
	// namespace
	Namespace string
	// prometheus 主机IP地址
	Host string
	// prometheus api端口
	Port string
	// 开始时间
	Start string
	// 结束时间
	End string
	// 步长
	Step string
	// 查询参数
	Query string
	// 最近几次超过阈值
	LastCount int
	// 阈值大于多少
	Gt int64
	// 阈值小于多少
	Lt int64
	// 操作,是扩容还是缩容,还是不操作, increase | reduce | none
	LtAction string
	// 阈值大于成立
	GtTrue bool
	// 阈值小成立
	LtTrue bool
	// 扩容步长
	IncreaseStep int32
	// 扩容或缩容间隔
	Interval time.Duration
	// 操作间隔
	ActionInterval int64
	// 扩容
	// 服务版本
	ServiceVersion string
	// es地址
	Es string
	// 数据源来源
	DataSource string
	// 指标名称，系统自带
	MetricName string
	// 监控值
	MonitorValue float64
	// 环境名称
	Entname string
}

// 2018-02-19 11:43
// 获取prometheus api接口
func getUri(param QueryParam) string {
	url := "http://" + param.Host + ":" + param.Port + "/api/v1/query_range?"
	return url
}

// 2018-02-19 13:40
// 查询监控接口
func Query(param QueryParam) string {
	api := getUri(param)
	uri := api
	query := make([]string, 0)
	if param.Query == ""{
		param.Query = systemQuery(param.Namespace, param.ServiceName, param.MetricName, param.ServiceVersion)
	}
	query = append(query, "query="+param.Query)
	query = append(query, "start="+param.Start)
	query = append(query, "end="+param.End)
	query = append(query, "step="+param.Step)
	uri += strings.Join(query, "&")
	qurl, err := url.Parse(uri)
	uri = api + qurl.Query().Encode()
	logs.Info(uri)
	resp, err := http.Get(uri)
	if err != nil {
		logs.Error("获取监控接口失败", err)
		return ""
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		logs.Error("获取数据失败", err)
	}
	data := string(body)
	return data
}

// 2018-02-19 13;49
// 获取监控数据值
func getMonitorData(data string) []interface{} {
	obj := make(map[string]interface{})
	err := json.Unmarshal([]byte(data), &obj)
	if err != nil {
		logs.Error("解析监控数据失败", err)
		return make([]interface{}, 0)
	}

	if _, ok := obj["data"]; ok {
		result := obj["data"].(map[string]interface{})["result"].([]interface{})
		if len(result) > 0 {
			values := result[0].(map[string]interface{})["values"].([]interface{})
			return values
		}
	}
	return make([]interface{}, 0)
}

// 2018-02-19 13;50
// 获取最近几次的结果值
func GetLastCountData(param QueryParam) float64 {
	data := Query(param)
	values := getMonitorData(data)
	size := len(values)
	if param.LastCount == 0 {
		param.LastCount = 5
	}

	if len(values) > param.LastCount {
		values = values[size-param.LastCount:size]
	}

	avg := 0.00
	for _, v := range values {
		rv := v.([]interface{})
		if len(rv) > 1 {
			strV := rv[1].(string)
			floatV, err := strconv.ParseFloat(strV, 64)
			if err == nil {
				avg += floatV
			}
		}
	}
	return avg / float64(param.LastCount)
}

var (
	// 求平均值,多个实例的平均使用率
	monitorCpu    = `avg(container_cpu_usage_seconds_total{$PARAM})`
	monitorMem    = `sum (container_memory_working_set_bytes{$PARAM}) / sum (container_spec_memory_limit_bytes{$PARAM})  * 100 `
	trafficInput  = `sum(rate(container_network_receive_bytes_total{$PARAM}[1m]))`
	trafficOutput = `sum(rate(container_network_transmit_bytes_total{$PARAM}[1m]))`
)

func replaceParam(item string, param string) string {
	r :=  strings.Replace(item, "$PARAM", param, -1)
	logs.Info(r)
	return r
}

// 2018-02-20 15:14
// 设置 prometheus 监控查询参数
func systemQuery(namespace string, serviceName string, metricName string, serviceVersion string) string {
	param := `image!="",pod_name=~"^$SERVICE-NAME.*",name=~"^k8s_.*",instance=~"^.*$",namespace=~"^$NAMESPACE$"`
	param = strings.Replace(param, "$SERVICE-NAME", util.Namespace(serviceName, serviceVersion), -1)
	param = strings.Replace(param, "$NAMESPACE", namespace, -1)
	switch metricName {
	case "cpu":
		return replaceParam(monitorCpu, param)
	case "memory":
		return replaceParam(monitorMem, param)
	case "trafficInput":
		return replaceParam(trafficInput, param)
	case "trafficOutput":
		return  replaceParam(trafficOutput, param)
	}
	return ""
}