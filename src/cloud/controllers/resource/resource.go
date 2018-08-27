package resource

import (
	"cloud/sql"
	"cloud/util"
	"github.com/astaxie/beego"
	"strings"
	"cloud/models/app"
	"cloud/models/registry"
	"cloud/models/ci"
	"cloud/models/pipeline"
	"cloud/k8s"
)

type ControllerResource struct {
	beego.Controller
}


// 获取select选项
// 2018-08-23  15:32
func getServiceName(searchMap sql.SearchMap, clusterName string, entName string, q string, html []string) string {
	searchMap.Put("ClusterName", clusterName)
	searchMap.Put("Entname", entName)
	data := make([]app.CloudAppService, 0)
	q = sql.SearchSql(app.CloudAppService{}, q, searchMap)
	sql.Raw(q).QueryRows(&data)
	for _, v := range data {
		if v.ServiceName == "" {
			html = append(html, util.GetSelectOptionName(v.AppName))
		}else {
			html = append(html, util.GetSelectOptionName(v.AppName + ";" + v.ResourceName + ";" + v.ServiceName))
		}
	}
	return strings.Join(html, "\n")
}


// 获取select选项
// 2018-08-23  10:19
func getAppTemplateName(searchMap sql.SearchMap, clusterName string, html []string) string {
	searchMap.Put("ClusterName", clusterName)
	data := make([]app.CloudAppTemplate, 0)
	q := sql.SearchSql(app.CloudAppTemplate{}, app.SelectCloudAppTemplate, searchMap)
	sql.Raw(q).QueryRows(&data)
	for _, v := range data {
		html = append(html, util.GetSelectOptionName(v.TemplateName))
	}
	return strings.Join(html, "\n")
}

// 获取select选项
// 2018-08-23  10:19
func getAppConfigName(searchMap sql.SearchMap, clusterName string, entName string, html []string) string {
	searchMap.Put("ClusterName", clusterName)
	searchMap.Put("Entname", entName)
	data := make([]app.CloudAppConfigure, 0)
	q := sql.SearchSql(app.CloudAppConfigure{}, app.SelectCloudAppConfigure, searchMap)
	sql.Raw(q).QueryRows(&data)
	for _, v := range data {
		html = append(html, util.GetSelectOptionName(v.ClusterName))
	}
	return strings.Join(html, "\n")
}

// 获取仓库组
// 2018-08-23  10:49
func getRegistryGroupName(searchMap sql.SearchMap, clusterName string, html []string) string {
	searchMap.Put("ClusterName", clusterName)
	data := make([]registry.CloudRegistryGroup, 0)
	q := sql.SearchSql(registry.CloudRegistryGroup{}, registry.SelectCloudRegistryGroup, searchMap)
	sql.Raw(q).QueryRows(&data)
	for _, v := range data {
		html = append(html, util.GetSelectOptionName(v.GroupName))
	}
	return strings.Join(html, "\n")
}


// 获取构建项目
// 2018-08-23  10:55
func getJobName(searchMap sql.SearchMap, clusterName string, html []string) string {
	searchMap.Put("ClusterName", clusterName)
	data := make([]ci.CloudBuildJob, 0)
	q := sql.SearchSql(ci.CloudBuildJob{}, ci.SelectCloudBuildJob, searchMap)
	sql.Raw(q).QueryRows(&data)
	for _, v := range data {
		html = append(html, util.GetSelectOptionName(v.ItemName))
	}
	return strings.Join(html, "\n")
}


// 获取构建项目
// 2018-08-23  10:55
func getDockerFile(searchMap sql.SearchMap, clusterName string, html []string) string {
	searchMap.Put("ClusterName", clusterName)
	data := make([]ci.CloudCiDockerfile, 0)
	q := sql.SearchSql(ci.CloudCiDockerfile{}, ci.SelectCloudCiDockerfile, searchMap)
	sql.Raw(q).QueryRows(&data)
	for _, v := range data {
		html = append(html, util.GetSelectOptionName(v.Name))
	}
	return strings.Join(html, "\n")
}

// 获取构建项目
// 2018-08-23  11:03
func getFlowName(searchMap sql.SearchMap, clusterName string, html []string) string {
	searchMap.Put("ClusterName", clusterName)
	data := make([]pipeline.CloudPipeline, 0)
	q := sql.SearchSql(pipeline.CloudPipeline{}, pipeline.SelectCloudPipeline, searchMap)
	sql.Raw(q).QueryRows(&data)
	for _, v := range data {
		html = append(html, util.GetSelectOptionName(v.PipelineName))
	}
	return strings.Join(html, "\n")
}

// 获取构建项目
// 2018-08-23  11:03
func getLbServiceName(searchMap sql.SearchMap, clusterName string, entName string,html []string) string {
	searchMap.Put("ClusterName", clusterName)
	searchMap.Put("Entname", entName)
	data := make([]k8s.CloudLbService, 0)
	q := sql.SearchSql(k8s.CloudLbService{}, k8s.SelectCloudLbService, searchMap)
	sql.Raw(q).QueryRows(&data)
	for _, v := range data {
		html = append(html, util.GetSelectOptionName(v.Domain))
	}
	return strings.Join(html, "\n")
}


// 获取构建项目
// 2018-08-23  11:18
func getCiServiceName(searchMap sql.SearchMap, clusterName string, entName string,html []string) string {
	searchMap.Put("ClusterName", clusterName)
	searchMap.Put("Entname", entName)
	data := make([]ci.CloudCiService, 0)
	q := sql.SearchSql(ci.CloudCiService{}, ci.SelectCloudCiService, searchMap)
	sql.Raw(q).QueryRows(&data)
	for _, v := range data {
		html = append(html, util.GetSelectOptionName(v.Domain))
	}
	return strings.Join(html, "\n")
}

// 获取批量部署
// 2018-08-23  11:18
func getCiBatchJobName(searchMap sql.SearchMap, clusterName string,html []string) string {
	searchMap.Put("ClusterName", clusterName)
	data := make([]ci.CloudCiBatchJob, 0)
	q := sql.SearchSql(ci.CloudCiBatchJob{}, ci.SelectCloudCiBatchJob, searchMap)
	sql.Raw(q).QueryRows(&data)
	for _, v := range data {
		html = append(html, util.GetSelectOptionName(v.BatchName))
	}
	return strings.Join(html, "\n")
}

// 获取环境数据
// 2018-01-20 17:45
// router /api/resource/name [get]
func (this *ControllerResource) GetResourceSelect() {
	html := make([]string, 0)
	clusterName := this.GetString("ClusterName")
	entName := this.GetString("EntName")
	Type := this.GetString("Type")
	user := util.GetUser(this.GetSession("username"))
	searchMap := sql.SearchMap{}
	searchMap.Put("CreateUser", user)
	var result string
	switch Type {
	case "服务":
		result = getServiceName(searchMap, clusterName, entName, app.SelectCloudAppService, html)
		break
	case "应用":
		result = getServiceName(searchMap, clusterName, entName, app.SelectCloudApp, html)
		break
	case "应用模板":
		result = getAppTemplateName(searchMap, clusterName, html)
		break
	case "配置管理":
		result = getAppConfigName(searchMap, clusterName, entName, html)
		break
	case "镜像仓库组":
		result = getRegistryGroupName(searchMap, clusterName, html)
		break
	case "构建项目":
		result = getJobName(searchMap, clusterName, html)
		break
	case "服务发布":
		result = getCiServiceName(searchMap, clusterName, entName, html)
		break
	case "流水线":
		result = getFlowName(searchMap, clusterName, html)
		break
	case "批量部署":
		result = getCiBatchJobName(searchMap, clusterName, html)
		break
	case "DockerFile":
		result = getDockerFile(searchMap, clusterName, html)
		break
	case "负载均衡":
		result = getLbServiceName(searchMap, clusterName, entName, html)
		break
	default:
		break
	}
	if result == "" {
		result = util.GetSelectOptionName("无数据")
	}
	this.Ctx.WriteString(result)
}
