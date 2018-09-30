package cluster

import (
	"time"
	"strconv"
	"fmt"
	"cloud/models/cluster"
	"strings"
	"cloud/k8s"
	"cloud/util"
	"cloud/sql"
	"encoding/json"
	"cloud/cache"
	"github.com/astaxie/beego/logs"
	"github.com/garyburd/redigo/redis"
)

// 2018-02-19 09:37
// 获取集群数据
func GetClusterDetailData(name string) cluster.CloudClusterDetail {
	detail := cluster.CloudClusterDetail{}
	status := util.RedisObj2Obj(cache.ClusterCache.Get("detail"+name), &detail)
	if ! status {
		go CacheClusterDetailData(name)
	}
	r := cache.ClusterComponentStatusesCache.Get(name)
	health, _ := redis.String(r, nil)
	detail.Health = health
	return detail
}

// 2018-02-19 09:13
// 缓存集群详情到redis里
// 获取集群资源使用详细情况
func CacheClusterDetailData(name string) {
	data := cluster.CloudCluster{}
	searchMap := sql.GetSearchMapV("ClusterName", name)
	searchSql := sql.SearchSql(
		cluster.CloudCluster{},
		cluster.SelectCloudCluster,
		searchMap)

	sql.Raw(searchSql).QueryRow(&data)
	detail := cluster.CloudClusterDetail{}
	if data.ClusterId == 0 {
		return
	}
	temp, _ := json.Marshal(data)
	json.Unmarshal(temp, &detail)
	c, err := k8s.GetClient(name)
	detail.ClusterId = data.ClusterId
	if err == nil {
		clusterStatus := k8s.GetNodeFromCluster(c)
		detail.ClusterCpu = clusterStatus.CpuNum
		detail.ClusterMem = clusterStatus.MemSize
		detail.ClusterNode = clusterStatus.Nodes
		detail.ClusterPods = clusterStatus.PodNum
		detail.OsVersion = clusterStatus.OsVersion
		if detail.ClusterCpu > 0 && detail.ClusterMem > 0 {
			used := k8s.GetClusterUsed(c)
			detail.UsedMem = used.UsedMem
			detail.UsedCpu = used.UsedCpu
			floatCpu := (float64(detail.UsedCpu) / float64(detail.ClusterCpu)) * 100
			floatMem := (float64(detail.UsedMem) / float64(detail.ClusterMem)) * 100
			cp, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", floatCpu), 64)
			mp, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", floatMem), 64)
			detail.CpuUsePercent = cp
			detail.MemUsePercent = mp
			detail.MemFree = detail.ClusterMem - detail.UsedMem
			detail.CpuFree = detail.ClusterCpu - detail.UsedCpu
			detail.Services = used.Services
		}
	}
	if cache.ClusterCacheErr == nil {
		// 30小时
		cache.ClusterCache.Put("detail"+name, util.ObjToString(detail), time.Hour*30)
	}
}

// 获取cluster数据
// 2018-01-26 13:38
func GetClusterSelect() string {
	html := make([]string, 0)
	data := GetClusterName()
	for _, v := range data {
		html = append(html, util.GetSelectOptionName(v.ClusterName))
	}
	return strings.Join(html, "\n")
}

// 设置公共数据
// 2018-01-18 20:52
func setClusterData(cStatus k8s.ClusterStatus, k cluster.CloudClusterDetail) k8s.ClusterStatus {
	cStatus.ClusterId = k.ClusterId
	cStatus.ClusterName = k.ClusterName
	cStatus.ClusterAlias = k.ClusterAlias
	cStatus.ClusterType = k.ClusterType
	return cStatus
}

// 使用线程获取cluster数据信息
func GetClusterInfo(k cluster.CloudClusterDetail, cData *[]k8s.ClusterStatus) {
	client, err := k8s.GetClient(k.ClusterName)
	cStatus := k8s.ClusterStatus{}
	if err != nil {
		cStatus = setClusterData(cStatus, k)
		*cData = append(*cData, cStatus)
		return
	}
	if err == nil {
		cStatus = k8s.GetNodeFromCluster(client)
	}
	cStatus = setClusterData(cStatus, k)
	*cData = append(*cData, cStatus)
}

// 获取集群名称和别名的对应关系
// 2018-01-21 17:46
func GetClusterMap() util.Lock {
	data := GetClusterName()
	r := util.Lock{}
	for _, v := range data {
		r.Put(v.ClusterName, v.ClusterAlias)
	}
	return r
}

// 获取所有集群名称信息
func GetClusterName() []cluster.CloudClusterName {
	data := make([]cluster.CloudClusterName, 0)

	searchSql := sql.SearchSql(
		cluster.CloudClusterName{},
		cluster.SelectCloudCluster,
		sql.SearchMap{})

	sql.Raw(searchSql).QueryRows(&data)
	return data
}

// 通过多线程添加和获取数据
func GoGetCluseterDetail(clusterName string, details *cluster.CloudClusterDetail) {
	detail := GetClusterDetailData(clusterName)
	details.ClusterCpu = details.ClusterCpu + detail.ClusterCpu
	details.ClusterMem = details.ClusterMem + detail.ClusterMem
	details.ClusterNode = details.ClusterNode + detail.ClusterNode
	details.ClusterPods = details.ClusterPods + detail.ClusterPods
	details.Services = details.Services + detail.Services
	details.UsedMem = details.UsedMem + detail.UsedMem
	details.UsedCpu = details.UsedCpu + detail.UsedCpu
	details.Couters = details.Couters + 1
}

// 获取集群数据,首页使用
func GetClusterData(clusterName string) cluster.CloudClusterDetail {
	details := cluster.CloudClusterDetail{}
	if clusterName == "" {
		data := GetClusterName()
		for _, d := range data {
			go GoGetCluseterDetail(d.ClusterName, &details)
		}
		t := 0
		for {
			if details.Couters >= len(data) {
				break
			}
			if t > 6 {
				break
			}
			t += 1
			time.Sleep(time.Second * 1)
		}
	} else {
		GoGetCluseterDetail(clusterName, &details)
	}
	if details.ClusterCpu > 0 && details.ClusterMem > 0 {
		floatCpu := (float64(details.UsedCpu) / float64(details.ClusterCpu)) * 100
		floatMem := (float64(details.UsedMem) / float64(details.ClusterMem)) * 100
		cp, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", floatCpu), 64)
		mp, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", floatMem), 64)
		details.CpuUsePercent = cp
		details.MemUsePercent = mp
		details.MemFree = details.ClusterMem - details.UsedMem
		details.CpuFree = details.ClusterCpu - details.UsedCpu
	}
	return details
}

func setClusterJson(this *ClusterController, data interface{}) {
	this.Data["json"] = data
	this.ServeJSON(false)
}

func getUsername(this *ClusterController) string {
	return util.GetUser(this.GetSession("username"))
}

// 2018-02-20 08:46
// 缓存集群信息到redis中,不用每次查库
// 任务计划自动更新数据
// 每一分钟一次
func CacheClusterHealthData() {
	logs.Info("生成组件健康数据")
	data := make([]cluster.CloudClusterDetail, 0)
	sql.Raw(cluster.SelectCloudCluster).QueryRows(&data)
	for _, k := range data {
		go k8s.GetClusterStatus(k.ClusterName)
	}
}

// 2018-02-19 08:46
// 缓存集群信息到redis中,不用每次查库
// 任务计划自动更新数据
func CacheClusterData() {
	cData := make([]k8s.ClusterStatus, 0)
	data := make([]cluster.CloudClusterDetail, 0)
	sql.Raw(cluster.SelectCloudCluster).QueryRows(&data)
	for _, k := range data {
		go CacheClusterDetailData(k.ClusterName)
		go GetClusterInfo(k, &cData)

	}
	counter := 0
	for {
		if len(cData) >= len(data) {
			break
		}
		time.Sleep(time.Second * 3)
		if counter > 3 {
			break
		}
		counter += 3
	}
	if cache.ClusterCacheErr == nil {
		for _, v := range cData {
			cache.ClusterCache.Put("data" + v.ClusterName, util.ObjToString(v), time.Hour*80)
		}
		// 80 小时
		//cache.ClusterCache.Put("data", util.ObjToString(cData), time.Hour*80)
	}

}

// 2018-02-19 08:53
// 获取缓存的集群数据数据
func getClusterCacheData() []k8s.ClusterStatus {
	cData := make([]k8s.ClusterStatus, 0)
	r := cache.ClusterCache.Get("data")
	status := util.RedisObj2Obj(r, &cData)
	if ! status {
		go CacheClusterData()
	}
	return cData
}
