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
)

type PipelineController struct {
	beego.Controller
}

// 2018-02-03 17:30
// 流水线入口页面
// @router /pipeline/list [get]
func (this *PipelineController) PipelineList() {
	this.TplName = "pipeline/list.html"
}

// 流水线历史入口页面
// @router /pipeline/history/list [get]
func (this *PipelineController) PipelineHistoryList() {
	this.TplName = "pipeline/history.html"
}

// 流水线详情页面
// @router /pipeline/detail/:hi(.*) [get]
func (this *PipelineController) PipelineDetail() {
	data,_ := getPipeData(this)
	jobName := this.Ctx.Input.Param("JobName")
	searchMap := sql.GetSearchMapV("CreateUser", getUser(this))
	if jobName != ""{
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
func (this *PipelineController) PipelineAdd() {
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
func (this *PipelineController) PipelineSave() {
	d := pipeline.CloudPipeline{}
	err := this.ParseForm(&d)
	if err != nil {
		this.Ctx.WriteString("参数错误" + err.Error())
		return
	}

	user := getUser(this)
	util.SetPublicData(d, user, &d)
	jobData := ci.GetJobName(user, d.ClusterName, d.JobName)

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
	}else{
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
func checkPipelineQuota(username string) (bool,string) {
	quotaDatas := quota.GetUserQuotaData(username, "")
	for _, v := range quotaDatas {
		if v.PipelineUsed + 1 > v.PipelineNumber {
			return false, "流水线数量超过配额限制"
		}
	}
	return true, ""
}

// 2018-02-04 13:29
// 流水线运行历史
// @router /api/pipeline/history [get]
func (this *PipelineController) PipelineHistoryData() {
	data := []pipeline.CloudPipelineLog{}
	searchMap := sql.SearchMap{}
	key := this.GetString("key")
	user := getUser(this)
	searchMap.Put("CreateUser", user)

	searchSql := sql.SearchSql(pipeline.CloudPipelineLog{}, pipeline.SelectCloudPipelineLog, searchMap)
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
		sql.CountSearchMap("cloud_pipeline_log", sql.GetSearchMapV("CreateUser", user), len(data), key),
		this.GetString("draw"))
	setPipelineJson(this, r)
}

// 流水线数据
// @router /api/pipeline [get]
func (this *PipelineController) PipelineData() {
	data := []pipeline.CloudPipeline{}
	searchMap := sql.SearchMap{}
	key := this.GetString("key")
	user := getUser(this)
	searchMap.Put("CreateUser", user)
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

	appDataMap := app.GetAppServiceDataMap()
	result := []pipeline.CloudPipeline{}
	for _, v := range data{
		tk := v.ClusterName + v.AppName + v.ServiceName
		if _,ok := appDataMap.Get(tk); ! ok {
			v.Status = "false"
		}
		result = append(result, v)
	}

	r := util.ResponseMap(result,
		sql.CountSearchMap("cloud_pipeline", sql.GetSearchMapV("CreateUser", user), len(data), key),
		this.GetString("draw"))
	setPipelineJson(this, r)
}

// 2018-02-03 22:03
// 获取流水线数据
func getPipeData(this *PipelineController) (pipeline.CloudPipeline, sql.SearchMap) {
	searchMap := sql.SearchMap{}
	id := this.Ctx.Input.Param(":id")
	if id != "" {
		searchMap.Put("PipelineId", id)
	}
	name := this.Ctx.Input.Param(":hi")
	if name != "" {
		searchMap.Put("PipelineName", name)
	}
	searchMap.Put("CreateUser", getUser(this))
	data := pipeline.CloudPipeline{}
	q := sql.SearchSql(data, pipeline.SelectCloudPipeline, searchMap)
	sql.Raw(q).QueryRow(&data)
	return data, searchMap
}

// 2018-02-05 12:56
// 获取用户名
func getUser(this *PipelineController) string {
	return util.GetUser(this.GetSession("username"))
}

// json
// 删除流水线
// @router /api/network/lb/service/:id:int [delete]
func (this *PipelineController) PipelineDelete() {
	pipedata, searchMap := getPipeData(this)
	q := sql.DeleteSql(pipeline.DeleteCloudPipeline, searchMap)
	r, err := sql.Raw(q).Exec()
	data := util.DeleteResponse(err,
		*this.Ctx, "删除流水线"+pipedata.PipelineName,
		this.GetSession("username"),
		pipedata.CreateUser,
		r)
	q = sql.DeleteSql(lb.DeleteCloudLbNginxConf, searchMap)
	sql.Raw(q).Exec()
	setPipelineJson(this, data)

}

func setPipelineJson(this *PipelineController, data interface{}) {
	this.Data["json"] = data
	this.ServeJSON(false)
}

// 2018-02-04 21:22
// 创建镜像安全数据
func createImagePullSecret(jobHistroyData ci2.CloudBuildJobHistory, start int64, jobName string, pipelog pipeline.CloudPipelineLog, serviceData app2.CloudAppService) bool {
	registryServer := strings.Split(jobHistroyData.RegistryServer, ":")
	servers := registry.GetRegistryServer(registryServer[0])
	user := jobHistroyData.CreateUser
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

	if jobHistroyData.RegistryServer == "" {
		pipelog.Messages = "获取仓库服务器错误"
		updatePipelineLog(pipelog, start, jobName, user)
		return false
	}

	serviceParam := k8s.ServiceParam{
		RegistryAuth: servers[0].Admin + ":" + util.Base64Decoding(servers[0].Password),
		Registry:     jobHistroyData.RegistryServer,
		Namespace:    namespace,
		Cl3:          cl,
	}
	k8s.CreateImagePullSecret(serviceParam)
	return true
}

// 2018-02-04 08:32
// 后台执行流水线程序
func startPipeline(user string, pipedata pipeline.CloudPipeline) {
	jobData := ci.GetJobName(user, pipedata.ClusterName, pipedata.JobName)
	if len(jobData) == 0 {
		logs.Error("获取构建程序失败", jobData)
		return
	}

	job := jobData[0]
	jobName := "job-" + util.Md5Uuid()
	go ci.JobExecStart(job, user, jobName, "")
	pipelog := pipeline.CloudPipelineLog{}

	temp, _ := json.Marshal(pipedata)
	json.Unmarshal(temp, &pipelog)
	job.JobName = jobName

	pipelog.JobId = job.JobId
	pipelog.CreateUser = user
	pipelog.CreateTime = pipedata.CreateTime
	pipelog.JobName = jobName
	pipelog.StartTime = util.GetDate()

	q := sql.InsertSql(pipelog, pipeline.InsertCloudPipelineLog)
	sql.Raw(q).Exec()

	start := time.Now().Unix()
	var logsr string
	var count int
	updatePipelogTime("start_job_time", jobName)
	for {
		if count > 500 || count > job.TimeOut {
			break
		}
		logsr = ci.GetJobLogs(job)
		if strings.Contains(logsr, "构建失败") {
			updateBuildFaild(jobName, "构建失败" , "执行构建失败")
			return
		}
		if strings.Contains(logsr, "构建完成") || strings.Contains(logsr, "构建失败") {
			updateBuildFaild(jobName, "构建成功" , "执行构建成功")
			break
		}
		time.Sleep(time.Second * 5)
		count += 1
	}

	updatePipelogTime("end_job_time", jobName)
	serviceData := app.GetServiceData(pipedata.ServiceName, pipedata.ClusterName, pipedata.AppName)
	serviceStatus := updateServiceErrorStatus(serviceData, jobName)
	if !serviceStatus{
		return
	}

	jobstatus, jobHistroyData := getHistoryDataStatus(jobName)
	if !jobstatus {
		return
	}

	logs.Info("getHistoryDataStatus", util.ObjToString(jobHistroyData))
	// 创建拉取镜像时的安全=数据
	createImage := createImagePullSecret(jobHistroyData, start, jobName, pipelog, serviceData)
	if !createImage {
		return
	}

	oldImage := serviceData.ImageTag
	serviceStatus, serviceData = updateServiceStart(jobName, user, serviceData, jobHistroyData)
	logs.Info("获取到serviceStatus", util.ObjToString(serviceData))
	if ! serviceStatus {
		return
	}

	serviceDatas := []app2.CloudAppService{}
	serviceDatas = append(serviceDatas, serviceData)
	count = 0
	var updateStatus bool
	for {
		if count > 300 {
			serviceData.ImageTag = oldImage
			app.ExecUpdate(serviceData, "image", user)
			pipelog.Messages = "获取服务状态超时,10分钟"
			updatePipelogStatus("update_service_end_time",
				"update_service_status",
				"失败",
				"update_service_errormsg",
				"获取服务状态超时,10分钟",
				jobName)
			break
		}
		go app.GoServerThread(serviceDatas)
		data := app.GetServiceRunData(serviceDatas)
		for _, v := range data {
			if v.Image == serviceData.ImageTag && v.Status == "True" && v.AvailableReplicas == int32(serviceData.Replicas) {
				pipelog.Messages = strings.Join([]string{v.Image, v.Status, strconv.Itoa(int(v.AvailableReplicas))}, " ")
				updateStatus = true

				updatePipelogStatus("update_service_end_time",
					"update_service_status",
					"成功",
					"update_service_errormsg",
					pipelog.Messages ,
					jobName)
				break
			}
			logs.Info("更新服务中",serviceData.ImageTag, v.Image, v.Status, v.AvailableReplicas, serviceData.Replicas,v.Image == serviceData.ImageTag , v.Status == "True" , v.AvailableReplicas == int32(serviceData.Replicas))
		}
		time.Sleep(time.Second * 5)
		count += 1
		if updateStatus {
			break
		}
	}
	updatePipelogTime("update_service_end_time", jobName)
	pipelog.Messages = ""
	updatePipelineLog(pipelog, start, jobName, user)
	updatePipelogTime("end_time", jobName)
}

// 更新拉取服务错误状态
func updateServiceErrorStatus(serviceData app2.CloudAppService,jobName string) bool  {
	if serviceData.ServiceId == 0 {
		updatePipelogStatus(
			"update_service_end_time",
			"update_service_status",
			"失败",
			"update_service_errormsg",
			"拉取服务失败了" ,
			jobName)
		updatePipelogTime("end_time", jobName)
		return false
	}
	return true
}

// 2018-02-05 14:42
// 更新服务启动
func updateServiceStart(jobName string, user string, serviceData app2.CloudAppService,jobHistroyData ci2.CloudBuildJobHistory) (bool,app2.CloudAppService) {
	serviceData.ImageTag = getImageTag(jobHistroyData)
	logs.Info("更新服务启动imageTag", serviceData.ImageTag, serviceData.ImageRegistry)
	if serviceData.MinReady == 0 {
		serviceData.MinReady = 50
	}

	updatePipelogTime("update_service_start_time", jobName)
	err := app.ExecUpdate(serviceData, "image", user)

	if err != nil {
		updatePipelogStatus(
			"update_service_end_time",
			"update_service_status",
			"失败",
			"update_service_errormsg",
			"更新服务失败" + err.Error(),
			jobName)
		updatePipelogTime("end_time", jobName)
		return false,serviceData
	}
	return true,serviceData
}

// 2018-02-05 14:51
// 更新构建失败状态
func updateBuildFaild(jobName string, status string , msg string)  {
	updatePipelogStatus(
		"end_job_time",
		"build_job_status",
		status,
		"build_job_errormsg",
		msg ,
		jobName)
	updatePipelogTime("end_time", jobName)
}

// 2018-02-05 14;39
// 获取构建任务数据,和更新流水线状态
func getHistoryDataStatus(jobName string) (bool,ci2.CloudBuildJobHistory) {
	jobHistroyData := ci.GetHistoryData(jobName)
	if jobHistroyData.BuildStatus != "构建成功" {
		updateBuildFaild(jobName, "失败", "构建任务执行失败")
		return false,jobHistroyData
	}else{
		updateBuildFaild(jobName, "成功", "构建任务成功")
	}
	return true,jobHistroyData
}

// 2018-02-06 13:04
// 获取镜像tag
func getImageTag(jobHistroyData ci2.CloudBuildJobHistory) string {
	imageTag := strings.Join([]string{
		jobHistroyData.RegistryServer,
		jobHistroyData.RegistryGroup,
		jobHistroyData.ItemName},
		"/") + ":" + jobHistroyData.ImageTag
	return imageTag
}

// 2018-02-05 12:59
// 更新日志各个表的时间点
func updatePipelogTime(columnt string, jobname string){
	q := "update cloud_pipeline_log set " + columnt + `="` + util.GetDate() + `" where job_name="`+jobname+`"`
	sql.Raw(q).Exec()
}

// 2018-02-05 14:09
// 更新日志状态数据
func updatePipelogStatus(timeColumnt string, statusColumnt string,  status string,statusMsgColumn string, statusmsg string, jobname string)  {
	q := "update cloud_pipeline_log set " +
		timeColumnt + `="` + util.GetDate() +
		`",`+statusColumnt+`="`+status+`", `+
		statusMsgColumn+`="`+statusmsg+`"
		 where job_name="`+jobname+`"`
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
	q := sql.UpdateSql(pipelog, pipeline.UpdateCloudPipelineLog, searchMap, "LogId")
	sql.Raw(q).Exec()
}

// 执行流水线任务
// 2018-02-03 22;01
// @router /api/pipeline/exec/:id:int [get]
func (this *PipelineController) PipelineExec() {
	pipedata, _ := getPipeData(this)
	user := util.GetUser(this.GetSession("username"))
	go startPipeline(user, pipedata)
	setPipelineJson(this, util.ApiResponse(true, "执行中,成功"))
}
