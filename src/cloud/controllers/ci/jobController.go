package ci

import (
	"github.com/astaxie/beego"

	"cloud/models/ci"
	"cloud/sql"
	"cloud/util"
	"cloud/controllers/image"
	"cloud/controllers/base/cluster"
	"strings"
	"cloud/k8s"
	"cloud/controllers/base/hosts"
	"time"
	"strconv"
	"github.com/astaxie/beego/logs"
	"k8s.io/client-go/kubernetes"
	"github.com/garyburd/redigo/redis"
	"k8s.io/apimachinery/pkg/util/json"
	registry2 "cloud/models/registry"
	"cloud/controllers/docker/application/app"
	"cloud/controllers/base/quota"
	"cloud/cache"
)

// 2018-01-25 17:54
// 构建任务配置
type JobController struct {
	beego.Controller
}

// 构建任务管理入口页面
// @router /ci/job/list [get]
func (this *JobController) JobList() {
	this.TplName = "ci/job/list.html"
}

// 2018-01-26 20:34
// 构建任务历史页面
// @router /ci/job/history/list [get]
func (this *JobController) JobHistoryList() {
	updateTimeOutJob()
	this.TplName = "ci/job/history.html"
}

// 构建任务详情入口
// @router /ci/job/detail/:hi(.*) [get]
func (this *JobController) JobDetail() {
	data := ci.CloudBuildJob{}
	searchMap := sql.GetSearchMap("JobId", *this.Ctx)
	q := sql.SearchSql(data, ci.SelectCloudBuildJob, searchMap)
	sql.Raw(q).QueryRow(&data)
	if data.BuildStatus == "" {
		data.BuildStatus = "未构建"
	}
	this.Data["data"] = data
	this.Data["content"] = len(strings.Split(data.Content, "\n"))
	this.TplName = "ci/job/detail.html"
}

// 2018-01-29 08;13
// 获取构建dockerfile
// @router /api/job/dockerfile/:id:int [get]
func (this *JobController) JobDockerfile() {
	data := ci.CloudBuildJobHistory{}
	historyId := this.Ctx.Input.Param(":id")
	jobName := this.Ctx.Input.Param(":hi")
	searchMap := sql.SearchMap{}
	if historyId != "" {
		searchMap.Put("HistoryId", historyId)
	}
	if jobName != "" {
		searchMap.Put("JobName", jobName)
	}
	q := sql.SearchSql(data, ci.SelectCloudBuildJobHistoryDockerfile, searchMap)
	sql.Raw(q).QueryRow(&data)
	this.Ctx.WriteString(data.DockerFile)
}

// 构建任务管理添加页面
// @router /ci/job/add [get]
func (this *JobController) JobAdd() {
	id := this.GetString("JobId")
	update := ci.CloudBuildJob{}
	update.ImageTag = "000"
	var dockerfile string
	var clusterHtml string
	var baseImageHtml string
	dockerfileData := GetDockerFileSelect()
	clusterData := cluster.GetClusterSelect()
	// 更新操作
	if id != "0" {
		searchMap := sql.GetSearchMap("JobId", *this.Ctx)
		q := sql.SearchSql(ci.CloudBuildJob{}, ci.SelectCloudBuildJob, searchMap)
		sql.Raw(q).QueryRow(&update)
		dockerfile = util.GetSelectOptionName(update.DockerFile)
		clusterHtml = util.GetSelectOptionName(update.ClusterName)
		baseImageHtml = util.GetSelectOptionName(update.BaseImage)
		this.Data["registryGroup"] = util.GetSelectOptionName(update.RegistryServer)
	}
	this.Data["ImageTag1"] = ""
	this.Data["ImageTag2"] = ""
	if update.ImageTag == "000" {
		this.Data["ImageTag1"] = "checked"
	} else {
		this.Data["ImageTag2"] = "checked"
	}
	dockerfile += dockerfileData
	clusterHtml += clusterData
	this.Data["cluster"] = clusterHtml
	this.Data["dockerfile"] = dockerfile
	this.Data["baseImage"] = baseImageHtml + registry.GetBaseImageSelect()
	this.Data["data"] = update
	this.TplName = "ci/job/add.html"
}

// 2018-02-03 21:37
// 获取选项
func GetSelectHtml(username string, clustername string) string {
	data := GetJobName(username, clustername, "")
	var html string
	for _, v := range data {
		html += util.GetSelectOption(v.ItemName, v.ItemName, v.ItemName)
	}
	return html
}

// 2018-02-03 21:55
// 获取构建任务信息
func GetJobName(username string, clustername string, itemname string) []ci.CloudBuildJob {
	searchMap := sql.GetSearchMapV("CreateUser", username)
	if clustername != "" {
		searchMap.Put("ClusterName", clustername)
	}
	if itemname != "" {
		searchMap.Put("ItemName", itemname)
	}
	// 构建任务数据
	data := make([]ci.CloudBuildJob, 0)
	q := sql.SearchSql(ci.CloudBuildJob{}, ci.SelectCloudBuildJob, searchMap)
	sql.Raw(q).QueryRows(&data)
	return data
}

// 获取构建任务数据
// 2018-01-25 17:57
// router /api/ci/job [get]
func (this *JobController) JobData() {
	setJson(this, GetJobName(util.GetUser(this.GetSession("username")), "", ""))
}

// string
// 构建任务保存
// @router /api/ci/job [post]
func (this *JobController) JobSave() {
	d := ci.CloudBuildJob{}
	err := this.ParseForm(&d)
	if err != nil {
		this.Ctx.WriteString("参数错误" + err.Error())
		return
	}

	util.SetPublicData(d, getUser(this), &d)
	d = updateJobContent(d)

	var q = sql.InsertSql(d, ci.InsertCloudBuildJob)
	if d.JobId > 0 {
		searchMap := sql.SearchMap{}
		searchMap.Put("JobId", d.JobId)
		searchMap.Put("CreateUser", getUser(this))
		q = sql.UpdateSql(d, ci.UpdateCloudBuildJob, searchMap, ci.UpdateCloudBuildJobExclude2)
		cache.JobDataCache.Delete(strconv.FormatInt(d.JobId, 10))
	}else{
		status, msg := checkQuota(getUser(this))
		if !status {
			data := util.ApiResponse(false, msg)
			setJson(this, data)
			return
		}
	}

	_, err = sql.Raw(q).Exec()
	data, msg := util.SaveResponse(err, "名称已经被使用")
	util.SaveOperLog(
		this.GetSession("username"),
		*this.Ctx, "保存构建任务配置 "+msg,
		d.JobName)
	setJson(this, data)
}

// 2018-02-12 09:30
// 检查镜像仓库配额
// 检查资源配额是否够用
func checkQuota(username string) (bool,string) {
	quotaDatas := quota.GetUserQuotaData(username, "")
	for _, v := range quotaDatas {
		if v.JobUsed + 1 > v.JobNumber {
			return false, "构建任务数量超过配额限制"
		}
	}
	return true, ""
}

// 更新任务脚本和dockerfile
func updateJobContent(jobData ci.CloudBuildJob) ci.CloudBuildJob  {
	if jobData.DockerFile != "0" {
		file := GetDockerfileData(jobData.DockerFile)
		if len(file) > 0 {
			jobData.Content = file[0].Content
			jobData.Script = file[0].Script
		}
	}
	return jobData
}


// 获取构建任务数据
// 2018-01-25 17:45
// router /api/ci/job/name [get]
func (this *JobController) JobDataName() {
	clustername := this.GetString("ClusterName")
	searchMap := sql.SearchMap{}
	searchMap.Put("CreateUser", getUser(this))
	if clustername != "" {
		searchMap.Put("ClusterName", clustername)
	}
	// 构建任务数据
	data := make([]ci.CloudBuildJob, 0)
	q := sql.SearchSql(ci.CloudBuildJob{}, ci.SelectCloudBuildJob, searchMap)
	sql.Raw(q).QueryRows(&data)
	setJson(this, data)
}

// 设置json数据
func setJson(this *JobController, data interface{}) {
	this.Data["json"] = data
	this.ServeJSON(false)
}

// 获取用户名
func getUser(this *JobController) string {
	return util.GetUser(this.GetSession("username"))
}

// 构建任务数据
// @router /api/ci/job/history [get]
func (this *JobController) JobHistoryDatas() {
	data := make([]ci.CloudBuildJobHistory, 0)
	searchMap := sql.SearchMap{}
	id := this.Ctx.Input.Param(":id")
	key := this.GetString("search")
	if id != "" {
		searchMap.Put("HistoryId", id)
	}
	searchMap.Put("CreateUser", getUser(this))

	searchSql := sql.SearchSql(ci.CloudBuildJobHistory{}, ci.SelectCloudBuildJobHistory, searchMap)
	if key != "" && id == "" {
		key = sql.Replace(key)
		searchSql += strings.Replace(ci.SelectBuildHistoryWhere, "?", key, -1)
	}

	num,err := sql.OrderByPagingSql(searchSql,
		"create_time",
		*this.Ctx.Request,
		&data,
		ci.CloudBuildJobHistory{})

	r := util.GetResponseResult(err,
		this.GetString("draw"),
		data,
		sql.Count("cloud_build_job_history", int(num), key))
	setJson(this, r)
}

// 构建任务数据
// @router /api/ci/job [get]
func (this *JobController) JobDatas() {

	data := make([]ci.CloudBuildJob, 0)
	searchMap := sql.SearchMap{}
	id := this.Ctx.Input.Param(":id")
	key := this.GetString("search")
	if id != "" {
		searchMap.Put("JobId", id)
	}
	searchMap.Put("CreateUser", getUser(this))
	searchSql := sql.SearchSql(ci.CloudBuildJob{}, ci.SelectCloudBuildJob, searchMap)

	if key != "" && id == "" {
		key = sql.Replace(key)
		searchSql += strings.Replace(ci.SelectCloudBuildJobWhere, "?", key, -1)
	}
	num,err := sql.OrderByPagingSql(searchSql,
		"create_time",
		*this.Ctx.Request,
		&data,
		ci.CloudBuildJob{})

	clusterMap := cluster.GetClusterMap()
	result := make([]ci.CloudBuildJob, 0)
	for _, v := range data {
		v.ClusterName = clusterMap.GetVString(v.ClusterName)
		result = append(result, v)
	}

	r := util.GetResponseResult(err,
		this.GetString("draw"),
		data,
		sql.Count("cloud_build_job", int(num), key))
	setJson(this, r)
}

// 获取数据
// 2018-01-26 15:02
func getJobData(this *JobController) ci.CloudBuildJob {
	jobData := ci.CloudBuildJob{}
	id := this.Ctx.Input.Param(":id")
	if cache.JobDataErr == nil {
		r := cache.JobDataCache.Get(id)
		s := util.RedisObj2Obj(r, &jobData)
		if s {
			logs.Info("从redis获取job数据", id, util.ObjToString(jobData))
			return jobData
		}
	}
	searchMap := sql.SearchMap{}
	searchMap.Put("JobId", id)
	searchMap.Put("CreateUser", getUser(this))
	sql.Raw(sql.SearchSql(jobData, ci.SelectCloudBuildJob, searchMap)).QueryRow(&jobData)
	if jobData.JobId > 0 {
		if cache.JobDataCache != nil {
			cache.JobDataCache.Put(id, util.ObjToString(jobData), time.Minute*10)
		}
	}
	logs.Info("从数据库获取job数据", util.ObjToString(jobData))
	return jobData
}

// json
// 删除构建任务
// 2018-01-25 17:46
// @router /api/ci/job/:id:int [delete]
func (this *JobController) JobDelete() {
	searchMap := sql.SearchMap{}
	id := this.Ctx.Input.Param(":id")
	searchMap.Put("JobId", id)
	jobData := getJobData(this)
	r, err := sql.Raw(sql.DeleteSql(ci.DeleteCloudBuildJob, searchMap)).Exec()
	data := util.DeleteResponse(err, *this.Ctx, "删除构建任务"+jobData.JobName, getUser(this), jobData.CreateUser, r)
	setJson(this, data)
}

// 2018-01-26 15:56
// 查看构建日志
// @router /ci/job/logs/:id:int [get]
func (this *JobController) JobLogsPage() {

	historyId := this.GetString("history")
	jobName := this.GetString("jobName")
	history := ci.CloudBuildJobHistory{}
	if historyId != "" || jobName != "" {
		searchMap := sql.SearchMap{}
		if historyId != "" {
			searchMap = sql.GetSearchMapV("HistoryId", historyId)
		}
		if jobName != "" {
			searchMap = sql.GetSearchMapV("JobName", jobName)
		}
		q := sql.SearchSql(history, ci.SelectCloudBuildJobHistory, searchMap)
		sql.Raw(q).QueryRow(&history)
	} else {
		jobData := getJobData(this)
		history = getLastLog(jobData.JobId)
	}

	if history.HistoryId == 0 {
		history.BuildLogs = "没有找到日志<span style='display:none;'>构建完成</span>"
	}
	this.Data["data"] = history
	this.TplName = "ci/job/log.html"
}


// 2018-01-26 15:59
// 获取最近一次日志信息
func getLastLog(jobId int64) ci.CloudBuildJobHistory {
	history := ci.CloudBuildJobHistory{}
	if cache.JobCacheErr == nil {
		r := cache.JobCache.Get(strconv.FormatInt(jobId, 10))
		if r != nil {
			rData, _ := redis.String(r, nil)
			json.Unmarshal([]byte(rData), &history)
			if history.HistoryId != 0 {
				logs.Info("从redis缓存获取到 cloud_job_history_", history.JobName)
				return history
			}
		}
	}
	q := strings.Replace(ci.SelectCloudBuildJobHistoryLast,
		"{0}",
		strconv.FormatInt(jobId, 10),
		-1)
	sql.Raw(q).QueryRow(&history)
	if cache.JobCacheErr == nil {
		cache.JobCache.Put(strconv.FormatInt(jobId, 10), util.ObjToString(history), time.Second*10)
	}
	return history
}

// 2018-02-04 15;57
// 获取历史数据,在流水线时候使用到
func GetHistoryData(jobName string) ci.CloudBuildJobHistory {
	history := ci.CloudBuildJobHistory{}
	q := sql.SearchSql(history, ci.SelectCloudBuildJobHistory, sql.GetSearchMapV("JobName", jobName))
	sql.Raw(q).QueryRow(&history)
	return history
}

/**
更新日志
2018-01-26 6:29
 */
func updateBuildLog(history ci.CloudBuildJobHistory, logsR string) {
	if len(logsR) < 10 {
		return
	}
	history.BuildLogs = logsR
	success := strings.Split(logsR, "完成构建...")
	if len(success) > 1 {
		t := strings.Split(success[1], "\n")
		finishTime := util.TimeToStamp(t[0])
		startTime := util.TimeToStamp(history.CreateTime)
		history.BuildTime = finishTime - startTime
		logs.Info("构建时间", util.ObjToString(history), t, finishTime, startTime, history.BuildTime)

	}
	searchMap := sql.SearchMap{}
	searchMap.Put("HistoryId", history.HistoryId)
	q := sql.UpdateSql(history, ci.UpdateCloudBuildJobHistory, searchMap, ci.ExcludeUpdateHistoryColumn)
	sql.Raw(q).Exec()
}

// 完成后更新数据
// 2018-01-26 20:27
func updateBuildResult(history ci.CloudBuildJobHistory, logsR string, jobData ci.CloudBuildJob, cl kubernetes.Clientset) {
	updateBuildLog(history, logsR)
	searchMap := sql.SearchMap{}
	searchMap.Put("JobId", jobData.JobId)
	jobData.BuildStatus = history.BuildStatus
	q := sql.UpdateSql(jobData, ci.UpdateCloudBuildJob, searchMap, ci.UpdateCloudBuildJobExclude)
	sql.Raw(q).Exec()
	k8s.DeleteJob(cl, history.JobName, "")
}

// 2018-02-04 08:19
// 获取job日志
// 流水线也用这个日志
func GetJobLogs(jobData ci.CloudBuildJob) string {
	cl, _ := k8s.GetClient(jobData.ClusterName)
	history := getLastLog(jobData.JobId)
	if history.BuildStatus != "构建中" {
		logs.Info("获取到开始构建")
		return history.BuildLogs
	}
	logsR := k8s.GetJobLogs(cl, history.JobName, util.Namespace("job", "job"))
	if logsR == "" {
		logsR = "没有找到日志<span style='display:none;'></span>"
		return logsR
	}

	if strings.Contains(logsR, "构建完成") || strings.Contains(logsR, "构建失败") {
		history.BuildStatus = "构建失败"
		if strings.Contains(logsR, "构建完成") {
			history.BuildStatus = "构建成功"
		}
		updateBuildResult(history, logsR, jobData, cl)
		go registry.UpdateGroupImageInfo()
	}
	updateBuildLog(history, logsR)
	return logsR
}

// 执行任务计划
// 2018-01-26 15;47
// @router /api/ci/job/logs/:id:int [get]
func (this *JobController) JobLogs() {
	jobData := getJobData(this)
	logs := GetJobLogs(jobData)
	this.Ctx.WriteString(logs)
}

// 2018-02-08 12:26
// 获取job参数
func getJobParam(jobData ci.CloudBuildJob,jobName string, registryServer string, groupData registry2.CloudRegistryServer) k8s.JobParam{
	master, port := hosts.GetMaster(jobData.ClusterName)
	param := k8s.JobParam{
		Master:         master,
		Port:           port,
		Jobname:        jobName,
		Itemname:       jobData.ItemName,
		RegistryServer: registryServer,
		Version:        jobData.LastTag,
		Timeout:        jobData.TimeOut,
		Dockerfile:     jobData.Content,
		Script:         jobData.Script,
		RegistryGroup:  jobData.RegistryServer,
		Images: jobData.BaseImage,
	}
	registryServers := registry.GetRegistryServer(groupData.ServerDomain)

	var authUser string
	if len(registryServers) > 0 {
		access := strings.Split(groupData.ServerAddress, ":")
		if len(access) > 1 {
			param.RegistryIp = access[0]
			param.RegistryServer = groupData.ServerDomain + ":" + access[1]
		}
		logs.Info("获取到镜像服务器地址", registryServers)
		authUser = registryServers[0].Admin + ":" + util.Base64Decoding(registryServers[0].Password)
		param.RegistryDomain = registryServers[0].ServerDomain
	}

	param.Auth = util.Base64Encoding(authUser)
	return param
}


// 2018-02-08 12:30
// 插入历史数据
func writeJobHistory(jobData ci.CloudBuildJob, jobId string, username string, registryServer string)  {
	history := ci.CloudBuildJobHistory{
		ImageTag:       jobData.LastTag,
		JobName:        jobId,
		JobId:          jobData.JobId,
		ItemName:       jobData.ItemName,
		DockerFile:     jobData.Content,
		Script:          jobData.Script,
		BuildStatus:    "构建中",
		BuildLogs:      "开始构建",
		RegistryGroup:  jobData.RegistryServer,
		CreateUser:     username,
		CreateTime:     util.GetDate(),
		ClusterName:    jobData.ClusterName,
		RegistryServer: registryServer,
	}
	i := sql.InsertSql(history, ci.InsertCloudBuildJobHistory)
	sql.Raw(i).Exec()
}

// 2018-02-03 22:13
// 执行构建任务
func JobExecStart(jobData ci.CloudBuildJob, username string, jobname string, registryAuth string) string {
	jobData = updateJobContent(jobData)
	// 按时间戳自动生成
	jobData.LastTag = jobData.ImageTag
	if jobData.ImageTag == "000" {
		jobData.LastTag = util.MakeImageTag()
	} else {
		jobData.LastTag = jobData.ImageTag
	}
	logs.Info("jobData", util.ObjToString(jobData))

	groupData, nodeIp, authServer, authDomain := registry.GetRegistryGroup(jobData.RegistryServer, jobData.ClusterName)
	if groupData.ServerAddress == "" {
		logs.Error("获取仓库服务失败", groupData)
		return ""
	}
	logs.Info("获取到仓库地址", groupData)

	param := getJobParam(jobData, jobname, groupData.ServerAddress, groupData)
	param.RegistryIp = nodeIp
	param.AuthServerIp = authServer
	param.AuthServerDomain = authDomain
	param.Script = jobData.Script
	param.ClusterName = jobData.ClusterName
	if registryAuth != "" {
		param.RegistryAuth = registryAuth
	}
	jobId := k8s.CreateJob(param)
	writeJobHistory(jobData, jobId, username, param.RegistryServer)
	jobData.BuildStatus = "构建中"

	searchMap := sql.SearchMap{}
	searchMap.Put("JobId", jobData.JobId)
	u := sql.UpdateSql(jobData,
		ci.UpdateCloudBuildJob,
		searchMap, ci.UpdateCloudBuildJobExclude)
	sql.Raw(u).Exec()
	return jobId
}

// 执行任务计划
// 2018-01-26 15;00
// @router /api/ci/job/exec/:id:int [get]
func (this *JobController) JobExec() {

	jobData := getJobData(this)
	logs.Info("获取到job数据", util.ObjToString(jobData))
	// 创建私密文件
	param := k8s.ServiceParam{}
	param.Image = jobData.BaseImage
	param.Namespace = util.Namespace("job", "job")
	cl,_ := k8s.GetClient(jobData.ClusterName)
	param.Cl3 = cl
	param = app.CreateSecretFile(param)
	jobName := "job-" + util.Md5Uuid()
	JobExecStart(jobData, getUser(this), jobName, param.Registry)
	data, msg := util.SaveResponse(nil, "构建中")
	util.SaveOperLog(this.GetSession("username"), *this.Ctx, "保存构建任务配置 "+msg, jobData.ItemName)
	setJson(this, data)
}

// 任务计划超时更新
func updateTimeOutJob()  {
	jobHistory := make([]ci.CloudBuildJobHistory, 0)
	sql.Raw(ci.SelectJobTimeout).QueryRows(&jobHistory)
	for _, v := range jobHistory{
		if util.TimeToStamp(util.GetDate()) - util.TimeToStamp(v.CreateTime) > v.BuildTime {
			q := ci.UpdateCloudBuildJobTimeout + util.ObjToString(v.HistoryId)
			sql.Raw(q).Exec()
		}
	}
}