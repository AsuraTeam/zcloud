package pipeline

import (
	"cloud/sql"
	"cloud/util"
	"github.com/astaxie/beego"
	"cloud/models/lb"
	"cloud/models/pipeline"
	"cloud/controllers/base/cluster"
	"cloud/controllers/docker/application/app"
	"cloud/controllers/ci"
	"encoding/json"
	"github.com/astaxie/beego/logs"
	"strings"
	"time"
	app2 "cloud/models/app"
	"strconv"
	"cloud/k8s"
	"cloud/controllers/image"
	ci2 "cloud/models/ci"
	"cloud/controllers/base/quota"
	"cloud/userperm"
)

type ControllerPipeline struct {
	beego.Controller
}

// 2018-02-03 17:30
// 流水线入口页面
// @router /pipeline/list [get]
func (this *ControllerPipeline) PipelineList() {
	this.TplName = "pipeline/list.html"
}

// 流水线历史入口页面
// @router /pipeline/history/list [get]
func (this *ControllerPipeline) PipelineHistoryList() {
	this.TplName = "pipeline/history.html"
}

// 2018-09-07 08:58
// 获取任务容器状态数据
// @router /pipeline/container/:id:int [get]
func (this *ControllerPipeline) GetPipelineContainer() {
	searchMap := sql.GetSearchMap("PipelineId", *this.Ctx)
	q := sql.SearchSql(pipeline.CloudPipeline{}, pipeline.SelectCloudPipeline, searchMap)
	data := pipeline.CloudPipeline{}
	sql.Raw(q).QueryRow(&data)
	this.Data["data"] = data
	this.TplName = "pipeline/container.html"
}

// 流水线详情页面
// @router /pipeline/detail/:hi(.*) [get]
func (this *ControllerPipeline) PipelineDetail() {
	data, _ := getPipeData(this)
	jobName := this.Ctx.Input.Param("JobName")
	searchMap := sql.SearchMap{}
	if jobName != "" {
		searchMap.Put("JobName", jobName)
	}
	if data.PipelineName != "" {
		searchMap.Put("PipelineName", data.PipelineName)
	}

	logData := pipeline.CloudPipelineLog{}
	q := sql.SearchSql(data, pipeline.SelectCloudPipelineLog, searchMap)
	q += pipeline.OrderByLogIdLimt1
	sql.Raw(q).QueryRow(&logData)

	this.Data["data"] = logData
	this.TplName = "pipeline/detail.html"
}

// 流水线管理添加页面
// @router /pipeline/add [get]
func (this *ControllerPipeline) PipelineAdd() {
	id := this.GetString("PipelineId")
	update := pipeline.CloudPipeline{}
	clusterData := cluster.GetClusterSelect()
	user := getUser(this)

	// 更新操作
	var PipelineData string
	var clusterHtml string
	var serviceHtml string
	var appHtml string
	var jobHtml string
	if id != "0" {
		searchMap := sql.SearchMap{}
		searchMap.Put("PipelineId", id)
		q := sql.SearchSql(pipeline.CloudPipeline{}, pipeline.SelectCloudPipeline, searchMap)
		sql.Raw(q).QueryRow(&update)
		this.Data["readonly"] = "readonly"

		clusterHtml = util.GetSelectOptionName(update.ClusterName)
		serviceHtml = util.GetSelectOptionName(update.ServiceName)
		appHtml = util.GetSelectOptionName(update.AppName)
		jobHtml = util.GetSelectOptionName(update.JobName)

		jobData := ci.GetSelectHtml(user, update.ClusterName)
		serviceData := app.GetSelectHtml(user, update.ClusterName)
		appData := app.GetAppHtml(user, update.ClusterName)
		this.Data["appData"] = appHtml + appData
		this.Data["serviceData"] = serviceHtml + serviceData
		this.Data["jobData"] = jobHtml + jobData
	}

	this.Data["cluster"] = clusterHtml + clusterData
	this.Data["data"] = update
	this.Data["PipelineData"] = PipelineData
	this.TplName = "pipeline/add.html"
}

// string
// 流水线保存
// @router /api/pipeline [post]
func (this *ControllerPipeline) PipelineSave() {
	d := pipeline.CloudPipeline{}
	err := this.ParseForm(&d)
	if err != nil {
		this.Ctx.WriteString("参数错误" + err.Error())
		return
	}

	user := getUser(this)
	util.SetPublicData(d, user, &d)
	jobData := ci.GetJobName(user, "", d.JobName)

	if len(jobData) == 0 {
		this.Ctx.WriteString("参数错误" + err.Error())
		return
	}

	d.JobId = jobData[0].JobId

	q := sql.InsertSql(d, pipeline.InsertCloudPipeline)
	if d.PipelineId > 0 {
		searchMap := sql.SearchMap{}
		searchMap.Put("PipelineId", d.PipelineId)
		searchMap.Put("CreateUser", user)
		q = sql.UpdateSql(d, pipeline.UpdateCloudPipeline, searchMap, "CreateTime,CreateUser,LbName")
	} else {
		status, msg := checkPipelineQuota(getUser(this))
		if !status {
			data := util.ApiResponse(false, msg)
			setPipelineJson(this, data)
			return
		}
	}
	_, err = sql.Raw(q).Exec()
	data, msg := util.SaveResponse(err, "名称已经被使用")

	util.SaveOperLog(this.GetSession("username"),
		*this.Ctx, "保存流水线服务配置 "+msg,
		d.PipelineName)
	setPipelineJson(this, data)
}

// 2018-02-12 11:22
// 检查流水线配额
// 检查资源配额是否够用
func checkPipelineQuota(username string) (bool, string) {
	quotaData := quota.GetUserQuotaData(username, "")
	for _, v := range quotaData {
		if v.PipelineUsed+1 > v.PipelineNumber {
			return false, "流水线数量超过配额限制"
		}
	}
	return true, ""
}

// 2018-02-04 13:29
// 流水线运行历史
// @router /api/pipeline/history [get]
func (this *ControllerPipeline) PipelineHistoryData() {
	data := make([]pipeline.CloudPipelineLog, 0)
	searchMap := sql.SearchMap{}
	key := this.GetString("key")
	searchSql := sql.SearchSql(pipeline.CloudPipelineLog{}, pipeline.SelectCloudPipelineLog, searchMap)
	searchSql = sql.GetWhere(searchSql, searchMap)
	if key != "" {
		q := `and (pipeline_name like "%?%" or service_name like "%?%")`
		replace := strings.Replace(q, "?", sql.Replace(key), -1)
		searchSql += replace
	}
	sql.OrderByPagingSql(searchSql,
		"log_id",
		*this.Ctx.Request,
		&data,
		pipeline.CloudPipelineLog{})

	r := util.ResponseMap(data,
		sql.CountSearchMap("cloud_pipeline_log", sql.SearchMap{}, len(data), key),
		this.GetString("draw"))
	setPipelineJson(this, r)
}

// 流水线数据
// @router /api/pipeline [get]
func (this *ControllerPipeline) PipelineData() {
	data := make([]pipeline.CloudPipeline, 0)
	searchMap := sql.SearchMap{}
	key := this.GetString("key")
	user := getUser(this)
	//searchMap.Put("CreateUser", user)
	searchSql := sql.SearchSql(pipeline.CloudPipeline{},
		pipeline.SelectCloudPipeline,
		searchMap)
	if key != "" {
		searchSql += strings.Replace(pipeline.SelectCloudPipelineWhere, "?", sql.Replace(key), -1)
	}

	sql.OrderByPagingSql(searchSql,
		"create_time",
		*this.Ctx.Request,
		&data,
		pipeline.CloudPipeline{})

	perm := userperm.GetResourceName("流水线", getUser(this))
	appDataMap := app.GetAppServiceDataMap()
	result := make([]pipeline.CloudPipeline, 0)
	for _, v := range data {
		tk := v.ClusterName + v.AppName + v.ServiceName
		if _, ok := appDataMap.Get(tk); ! ok {
			v.Status = "false"
		}
		// 不是自己创建的才检查
		if v.CreateUser != user && user != "admin" {
			if ! userperm.CheckPerm(v.PipelineName, v.ClusterName, "", perm) && len(user) > 0 {
				continue
			}
		}
		result = append(result, v)
	}

	r := util.ResponseMap(result,
		sql.CountSearchMap("cloud_pipeline", sql.SearchMap{}, len(data), key),
		this.GetString("draw"))
	setPipelineJson(this, r)
}

// 2018-02-03 22:03
// 获取流水线数据
func getPipeData(this *ControllerPipeline) (pipeline.CloudPipeline, sql.SearchMap) {
	searchMap := sql.SearchMap{}
	id := this.Ctx.Input.Param(":id")
	if id != "" {
		searchMap.Put("PipelineId", id)
	}
	name := this.Ctx.Input.Param(":hi")
	if name != "" {
		searchMap.Put("PipelineName", name)
	}
	//searchMap.Put("CreateUser", getUser(this))
	data := pipeline.CloudPipeline{}
	q := sql.SearchSql(data, pipeline.SelectCloudPipeline, searchMap)
	sql.Raw(q).QueryRow(&data)
	user := getUser(this)
	perm := userperm.GetResourceName("流水线", user)
	// 不是自己创建的才检查
	if data.CreateUser != user {
		if ! userperm.CheckPerm(data.PipelineName, data.ClusterName, "", perm) {
			return pipeline.CloudPipeline{}, sql.SearchMap{}
		}
	}
	return data, searchMap
}

// 2018-02-05 12:56
// 获取用户名
func getUser(this *ControllerPipeline) string {
	return util.GetUser(this.GetSession("username"))
}

// json
// 删除流水线
// @router /api/network/lb/service/:id:int [delete]
func (this *ControllerPipeline) PipelineDelete() {
	pipeData, searchMap := getPipeData(this)
	q := sql.DeleteSql(pipeline.DeleteCloudPipeline, searchMap)
	r, err := sql.Raw(q).Exec()
	data := util.DeleteResponse(err,
		*this.Ctx, "删除流水线"+pipeData.PipelineName,
		this.GetSession("username"),
		pipeData.CreateUser,
		r)
	q = sql.DeleteSql(lb.DeleteCloudLbNginxConf, searchMap)
	sql.Raw(q).Exec()
	setPipelineJson(this, data)

}

func setPipelineJson(this *ControllerPipeline, data interface{}) {
	this.Data["json"] = data
	this.ServeJSON(false)
}

// 2018-02-04 21:22
// 创建镜像安全数据
func createImagePullSecret(jobHistoryData ci2.CloudBuildJobHistory, start int64, jobName string, pipelog pipeline.CloudPipelineLog, serviceData app2.CloudAppService) bool {
	registryServer := strings.Split(jobHistoryData.RegistryServer, ":")
	servers := registry.GetRegistryServer(registryServer[0])
	user := jobHistoryData.CreateUser
	if len(servers) == 0 {
		logs.Error("获取镜像服务失败", "", serviceData, user)
		pipelog.Messages = "获取镜像服务失败"
		updatePipelineLog(pipelog, start, jobName, user)
		return false
	}
	namespace := util.Namespace(serviceData.AppName, serviceData.ResourceName)
	cl, err := k8s.GetClient(serviceData.ClusterName)
	if err != nil {
		pipelog.Messages = "获取镜像master失败"
		updatePipelineLog(pipelog, start, jobName, user)
		return false
	}

	if jobHistoryData.RegistryServer == "" {
		pipelog.Messages = "获取仓库服务器错误"
		updatePipelineLog(pipelog, start, jobName, user)
		return false
	}

	serviceParam := k8s.ServiceParam{
		RegistryAuth: servers[0].Admin + ":" + util.Base64Decoding(servers[0].Password),
		Registry:     jobHistoryData.RegistryServer,
		Namespace:    namespace,
		Cl3:          cl,
	}
	k8s.CreateImagePullSecret(serviceParam)
	return true
}

// 2018-02-04 08:32
// 后台执行流水线程序
func startPipeline(user string, pipeData pipeline.CloudPipeline) {
	jobData := ci.GetJobName(user, "", pipeData.JobName)
	if len(jobData) == 0 {
		logs.Error("获取构建程序失败", jobData, pipeData.ClusterName, pipeData.JobName)
		return
	}

	job := jobData[0]
	jobName := "job-" + util.Md5Uuid()
	go ci.JobExecStart(job, user, jobName, "")
	pipelog := pipeline.CloudPipelineLog{}

	temp, _ := json.Marshal(pipeData)
	json.Unmarshal(temp, &pipelog)
	job.JobName = jobName

	pipelog.JobId = job.JobId
	pipelog.CreateUser = user
	pipelog.CreateTime = pipeData.CreateTime
	pipelog.JobName = jobName
	pipelog.StartTime = util.GetDate()

	q := sql.InsertSql(pipelog, pipeline.InsertCloudPipelineLog)
	sql.Raw(q).Exec()

	start := time.Now().Unix()
	var count int
	updatePipeLogTime("start_job_time", jobName)
	for {
		if count > 500 || count > job.TimeOut {
			break
		}
		// 触发日志更新
		ci.GetJobLogs(job)
		jobHistory := ci.GetHistoryData(jobName)
		if jobHistory.BuildStatus == "构建失败" {
			updateBuildFaild(jobName, "构建失败", "执行构建失败")
			return
		}
		if jobHistory.BuildStatus == "构建超时" {
			updateBuildFaild(jobName, "构建超时", "执行构建超时")
			return
		}
		if jobHistory.BuildStatus == "构建成功" {
			updateBuildFaild(jobName, "构建成功", "执行构建成功")
			break
		}
		time.Sleep(time.Second * 5)
		count += 1
		logs.Info("job等待-->", jobName, count)
	}

	updatePipeLogTime("end_job_time", jobName)
	serviceData := app.GetServiceData(pipeData.ServiceName, pipeData.ClusterName, pipeData.AppName)
	serviceStatus := updateServiceErrorStatus(serviceData, jobName)
	if !serviceStatus {
		logs.Error("serviceStatus 异常")
		return
	}

	jobStatus, jobHistoryData := getHistoryDataStatus(jobName)
	if !jobStatus {
		logs.Error("jobStatus 异常")
		return
	}

	logs.Info("getHistoryDataStatus", util.ObjToString(jobHistoryData))
	// 创建拉取镜像时的安全=数据
	createImage := createImagePullSecret(jobHistoryData, start, jobName, pipelog, serviceData)
	if !createImage {
		logs.Error("createImage 异常")
		return
	}

	oldImage := serviceData.ImageTag
	serviceStatus, serviceData = updateServiceStart(jobName, user, serviceData, jobHistoryData)
	logs.Info("获取到serviceStatus", util.ObjToString(serviceData))
	if ! serviceStatus {
		logs.Error("serviceStatus 2 异常")
		return
	}

	serviceDatas := make([]app2.CloudAppService, 0)
	serviceDatas = append(serviceDatas, serviceData)
	count = 0
	var updateStatus bool
	for {
		if count > 300 {
			serviceData.ImageTag = oldImage
			app.ExecUpdate(serviceData, "image", user)
			pipelog.Messages = "获取服务状态超时,10分钟"
			updatePipeLogStatus("update_service_end_time",
				"update_service_status",
				"失败",
				"update_service_errormsg",
				"获取服务状态超时,10分钟",
				jobName)
			break
		}
		go app.GoServerThread(serviceDatas)
		data := app.GetServiceRunData(serviceDatas, "")
		for _, v := range data {
			logs.Info("更新服务中", serviceData.ImageTag, v.Image, v.Status, v.AvailableReplicas, serviceData.Replicas, v.Image == serviceData.ImageTag, v.Status == "True", v.AvailableReplicas == int32(serviceData.Replicas))
			if v.Image == serviceData.ImageTag && v.Status == "True" && v.AvailableReplicas == int32(serviceData.Replicas) {
				pipelog.Messages = strings.Join([]string{v.Image, v.Status, strconv.Itoa(int(v.AvailableReplicas))}, " ")
				updateStatus = true

				updatePipeLogStatus("update_service_end_time",
					"update_service_status",
					"成功",
					"update_service_errormsg",
					pipelog.Messages,
					jobName)
				break
			}
		}
		time.Sleep(time.Second * 5)
		count += 1
		if updateStatus {
			break
		}
	}
	updatePipeLogTime("update_service_end_time", jobName)
	pipelog.Messages = ""
	updatePipelineLog(pipelog, start, jobName, user)
	updatePipeLogTime("end_time", jobName)
}

// 更新拉取服务错误状态
func updateServiceErrorStatus(serviceData app2.CloudAppService, jobName string) bool {
	if serviceData.ServiceId == 0 {
		updatePipeLogStatus(
			"update_service_end_time",
			"update_service_status",
			"失败",
			"update_service_errormsg",
			"拉取服务失败了",
			jobName)
		updatePipeLogTime("end_time", jobName)
		return false
	}
	return true
}

// 2018-02-05 14:42
// 更新服务启动
func updateServiceStart(jobName string, user string, serviceData app2.CloudAppService, jobHistoryData ci2.CloudBuildJobHistory) (bool, app2.CloudAppService) {
	serviceData.ImageTag = getImageTag(jobHistoryData)
	logs.Info("更新服务启动imageTag", serviceData.ImageTag, serviceData.ImageRegistry)
	if serviceData.MinReady == 0 {
		serviceData.MinReady = 50
	}

	updatePipeLogTime("update_service_start_time", jobName)
	err := app.ExecUpdate(serviceData, "image", user)

	if err != nil {
		updatePipeLogStatus(
			"update_service_end_time",
			"update_service_status",
			"失败",
			"update_service_errormsg",
			"更新服务失败"+err.Error(),
			jobName)
		updatePipeLogTime("end_time", jobName)
		return false, serviceData
	}
	return true, serviceData
}

// 2018-02-05 14:51
// 更新构建失败状态
func updateBuildFaild(jobName string, status string, msg string) {
	updatePipeLogStatus(
		"end_job_time",
		"build_job_status",
		status,
		"build_job_errormsg",
		msg,
		jobName)
	updatePipeLogTime("end_time", jobName)
}

// 2018-02-05 14;39
// 获取构建任务数据,和更新流水线状态
func getHistoryDataStatus(jobName string) (bool, ci2.CloudBuildJobHistory) {
	jobHistoryData := ci.GetHistoryData(jobName)
	if jobHistoryData.BuildStatus != "构建成功" {
		updateBuildFaild(jobName, "失败", "构建任务执行失败")
		return false, jobHistoryData
	} else {
		updateBuildFaild(jobName, "成功", "构建任务成功")
	}
	return true, jobHistoryData
}

// 2018-02-06 13:04
// 获取镜像tag
func getImageTag(jobHistoryData ci2.CloudBuildJobHistory) string {
	imageTag := strings.Join([]string{
		jobHistoryData.RegistryServer,
		jobHistoryData.RegistryGroup,
		jobHistoryData.ItemName},
		"/") + ":" + jobHistoryData.ImageTag
	return imageTag
}

// 2018-02-05 12:59
// 更新日志各个表的时间点
func updatePipeLogTime(column string, jobName string) {
	q := "update cloud_pipeline_log set " + column + `="` + util.GetDate() + `" where job_name="` + jobName + `"`
	sql.Raw(q).Exec()
}

// 2018-02-05 14:09
// 更新日志状态数据
func updatePipeLogStatus(timeColumn string, statusColumn string, status string, statusMsgColumn string, statusMsg string, jobName string) {
	q := "update cloud_pipeline_log set " +
		timeColumn + `="` + util.GetDate() +
		`",` + statusColumn + `="` + status + `", ` +
		statusMsgColumn + `="` + statusMsg + `"
		 where job_name="` + jobName + `"`
	sql.Raw(q).Exec()
}

// 更新流水线结果
func updatePipelineLog(pipelog pipeline.CloudPipelineLog, start int64, jobName string, user string) {
	pipelog.Status = "执行失败"
	if pipelog.Messages == "" {
		pipelog.Messages = "更新成功"
		pipelog.Status = "执行成功"
	}
	pipelog.RunTime = time.Now().Unix() - start
	searchMap := sql.GetSearchMapV("JobName", jobName, "CreateUser", user, "CreateTime", pipelog.CreateTime)
	q := sql.UpdateSql(pipelog, pipeline.UpdateCloudPipelineLog, searchMap, pipeline.UpdateCloudPipelineLogExclude)
	sql.Raw(q).Exec()
}

// 执行流水线任务
// 2018-02-03 22;01
// @router /api/pipeline/exec/:id:int [get]
func (this *ControllerPipeline) PipelineExec() {
	pipeData, _ := getPipeData(this)
	user := util.GetUser(this.GetSession("username"))
	go startPipeline(user, pipeData)
	setPipelineJson(this, util.ApiResponse(true, "执行中,成功"))
}
