package quota

import (
	"cloud/sql"
	"cloud/controllers/users"
	"strings"
	"cloud/models/app"
	"cloud/models/ci"
	"cloud/models/pipeline"
	"cloud/models/lb"
	"github.com/astaxie/beego/orm"
	"strconv"
	"cloud/models/quota"
	"cloud/models/registry"
	"github.com/astaxie/beego/logs"
)

/* 资源限制查询
   包括pod数量
   服务数量
   应用数据
   cpu配额大小
   内存配额大小
   job数量
   lb数量
   对部门
   对用户
   实现 某个用户或部门 整体的使用资源控制
 */

// 2018-02-11 16:40
func getQuotaUsed(userData string, used quota.QuotaUsed) quota.QuotaUsed {

	used.MemoryUsed = getMemory(userData, used.QuotaName)
	used.CpuUsed = getCpu(userData, used.QuotaName)
	used.LbUsed = getLb(userData)
	used.PipelineUsed = getPipeline(userData)
	used.JobUsed = getJob(userData)
	used.AppUsed = getApps(userData, used.QuotaName)
	used.ServiceUsed = getService(userData, used.QuotaName)
	used.PodUsed = getPods(userData, used.QuotaName)
	used.RegistryGroupUsed = getRegistryGroup(userData)
	used.DockerFileUsed = getDockerFile(userData)

	used.PodFree = used.PodNumber - used.PodUsed
	used.ServiceFree = used.ServiceNumber - used.ServiceUsed
	used.AppFree = used.AppNumber - used.AppUsed
	used.LbFree = used.LbNumber - used.LbUsed
	used.JobFree = used.JobNumber - used.JobUsed
	used.PipelineFree = used.PipelineNumber - used.PipelineUsed
	used.CpuFree = used.QuotaCpu - used.CpuUsed
	used.RegistryGroupFree = used.RegistryGroupNumber - used.RegistryGroupUsed
	used.MemoryFree = (used.QuotaMemory - used.MemoryUsed) / 1024
	used.DockerFileFree = used.DockerFileNumber - used.DockerFileUsed

    used.RegistryGroupPercent = getPercent(used.RegistryGroupUsed, used.RegistryGroupNumber)
	used.PodPercent = getPercent(used.PodUsed, used.PodNumber)
	used.ServicePercent = getPercent(used.ServiceUsed, used.ServiceNumber)
	used.AppPercent = getPercent(used.AppUsed, used.AppNumber)
	used.LbPercent = getPercent(used.LbUsed, used.AppNumber)
	used.JobPercent = getPercent(used.JobUsed, used.JobNumber)
	used.PipelinePercent = getPercent(used.PipelineUsed, used.PipelineNumber)
	used.CpuPercent = getPercent(used.CpuUsed, used.QuotaCpu)
	used.MemoryPercent = getPercent(used.MemoryUsed, used.QuotaMemory)
	used.DockerFilePercent = getPercent(used.DockerFileUsed, used.DockerFileNumber)
	used.QuotaMemory = used.QuotaMemory / 1024
	used.MemoryUsed = used.MemoryUsed / 1024
	return used
}

// 2018-02-11 16:49
// 计算百分比
func getPercent(a int64, b int64) int64 {
	if b == 0 {
		b = 1
	}
	r := (float64(a) / float64(b)) * 100
	return int64(r)
}

// 2018-02-11 16:47
// 获取用户的配额使用情况
func setQuotaUserUsed(username string, used quota.QuotaUsed) quota.QuotaUsed {
	userData := `"` + username + `"`
	used = getQuotaUsed(userData, used)
	return used
}

// 2018-02-11 16:45
// 获取组的配额使用情况
func setQuotaGroupUsed(groupname string, used quota.QuotaUsed) quota.QuotaUsed {
	usersData := GetGroupUsers([]string{groupname})
	used = getQuotaUsed(usersData, used)
	return used
}

// 2018-02-11 15:24
// 获取同一部门的用户
// 没有授权的用户可以查询
func getUsers(username string) string {
	depts := users.GetUserDept(username)
	usersData := GetGroupUsers(depts)
	return usersData
}

// 2018-02-11 16:55
// 获取多个组里的所有用户
func GetGroupUsers(depts []string) string {
	userDepts := make([]string, 0)
	for _,v := range depts {
		userDepts = append(userDepts, `"` + v +`"`)
	}
	usersData := users.GetGroupUsers(userDepts)
	return strings.Join(usersData, ",")
}

// 2018-02-11 15:30
// 获取服务使用数量
func getService(userData string, quotaName string) int64 {
	services := make([]app.CloudAppService, 0)
	q := replace(
		userData,
		app.SelectUserServices)
	if quotaName != ""{
		q += ` and resource_name="{0}"`
	}
	q = strings.Replace(q, "{0}", quotaName, -1)
	sql.Raw(q).QueryRows(&services)
	return int64(len(services))
}

// 2018-02-11 15:35
// 获取用户服务使用量
func getPods(userData string, quotaName string) int64 {
	pods := make([]app.CloudContainer, 0)
	q := replace(
		userData,
		app.SelectUserContainer)
	if quotaName != ""{
		q += ` and resource_name="{0}"`
	}
	q = strings.Replace(q, "{0}", quotaName, -1)
	sql.Raw(q).QueryRows(&pods)
	return int64(len(pods))
}

// 2018-02-11 15:39
// 获取用户应用数量
func getApps(userData string, quotaName string) int64 {
	apps := make([]app.CloudApp, 0)
	q := replace(
		userData,
		app.SelectUserApp)
	if quotaName != ""{
		q += ` and resource_name="{0}"`
	}
	q = strings.Replace(q, "{0}", quotaName, -1)
	logs.Info("quotaName", quotaName, q)
	sql.Raw(q).QueryRows(&apps)
	return int64(len(apps))
}

// 2018-02-11 16:00
// 获取cpu或内存
func getCpuMemory(q string, key string) int64 {
	maps := make([]orm.Params, 0)

	sql.Raw(q).Values(&maps)
	if len(maps) > 0 {
		if maps[0][key] == nil {
			return 0
		}
		v, err := strconv.ParseInt(maps[0][key].(string), 10, 64)
		if err == nil {
			return v
		}
		return 0
	}
	return 0
}

// 2018-02-11 15:59
// 获取用户内存使用量
func getMemory(userData string, quotaName string) int64 {
	q := replace(userData, app.SelectUsersMemory)
	if quotaName != ""{
		q += ` and resource_name="{0}"`
	}
	logs.Info("quotaName", quotaName, q)
	q = strings.Replace(q, "{0}", quotaName, -1)
	return getCpuMemory(q, "memory")
}

// 2018-02-11 16:01
// 获取用户cpu使用量
func getCpu(userData string, quotaName string) int64 {
	q := replace(userData, app.SelectUsersCpu)
	if quotaName != ""{
		q += ` and resource_name="{0}"`
	}
	q = strings.Replace(q, "{0}", quotaName, -1)
	logs.Info("quotaName", quotaName, q)
	return getCpuMemory(q, "cpu")
}

// 2018-02-11 15:46
// 获取用户负载均衡数量
func getLb(userData string) int64 {
	lbs := make([]lb.CloudLb, 0)
	sql.Raw(
		replace(
			userData,
			lb.SelectUserLbs)).QueryRows(&lbs)
	return int64(len(lbs))
}

// 2018-02-12 09:35
// 获取用户负载均衡数量
func getDockerFile(userData string) int64 {
	dockerFile := make([]ci.CloudCiDockerfile, 0)
	sql.Raw(
		replace(
			userData,
			ci.SelectDockerfiles)).QueryRows(&dockerFile)
	return int64(len(dockerFile))
}

// 2018-02-12 08:05
// 获取用户负载均衡数量
func getRegistryGroup(userData string) int64 {
	groups := make([]registry.CloudRegistryGroup, 0)
	sql.Raw(
		replace(
			userData,
			registry.SelectUserRegistryGroups)).QueryRows(&groups)
	return int64(len(groups))
}

// 获取用户构建任务数量
// 获取job数量
func getJob(userData string) int64 {
	jobs := make([]ci.CloudBuildJob, 0)
	sql.Raw(
		replace(
			userData,
			ci.SelectUserJobs)).QueryRows(&jobs)
	return int64(len(jobs))
}

// 2018-02-11 15:41
// 获取用户流水线数量
func getPipeline(userData string) int64 {
	pipelines := make([]pipeline.CloudPipeline, 0)
	sql.Raw(
		replace(
			userData,
			pipeline.SelectUserPipeline)).QueryRows(&pipelines)
	return int64(len(pipelines))
}

func replace(userData string, sqlData string) string {
	q := strings.Replace(sqlData, "?", userData, -1)
	return q
}
