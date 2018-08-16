// @APIVersion 1.0.0
// @Title Zcloud API
// @Description api
// @Contact 270851812@qq.com
// @TermsOfServiceUrl http://github.com/Asura/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"github.com/astaxie/beego"
	"cloud/controllers/index"
	"cloud/controllers/base/cluster"
	"cloud/controllers/docker/application/app"
	"github.com/astaxie/beego/context"
	"strings"
	"github.com/astaxie/beego/logs"
	"cloud/util"
	"cloud/controllers/base/hosts"
	"cloud/controllers/base/quota"
	"cloud/controllers/base/lb"
	"cloud/controllers/users"
	"cloud/controllers/image"
	"cloud/controllers/ci"
	"cloud/controllers/base/storage"
	"cloud/controllers/pipeline"
	"cloud/controllers/operlog"
	"cloud/controllers/ent"
	"cloud/controllers/perm"
	"cloud/controllers/monitor"
)

func init() {

	// 公共入口
	beego.Router("/", &index.IndexController{}, "get:Index")
	beego.Router("/index", &index.IndexController{}, "get:Index")
	beego.Router("/index/detail/:hi(.*)", &index.IndexController{}, "get:IndexDetail")
	beego.Router("/login", &index.IndexController{}, "get:LoginPage")
	beego.Router("/api/user", &index.IndexController{}, "get:GetUser")
	beego.Router("/api/user/login", &index.IndexController{}, "post:Login")
	beego.Router("/api/user/logout", &index.IndexController{}, "*:OutLogin")
	beego.Router("/webtty/:id:int", &index.IndexController{}, "*:WebTty")

	applicationNs :=
		beego.NewNamespace("/application",
			// 应用交互中心,
			// 容器治理,
			// 应用管理页面,
			beego.NSNamespace("/app",
				beego.NSRouter("/list", &app.AppController{}, "get:AppList"),
				// 添加应用页面,
				beego.NSRouter("/add", &app.AppController{}, "get:AppAdd"),
				// 应用详情页面,
				beego.NSRouter("/detail/:id:int", &app.AppController{}, "get:AppDetail"),
			),
			beego.NSNamespace("/configure",
				beego.NSRouter("/list", &app.ConfigureController{}, "get:ConfigureList"),
				// 添加配置页面
				beego.NSRouter("/add", &app.ConfigureController{}, "get:ConfigureAdd"),
				// 应用详情页面,
				beego.NSRouter("/detail/:hi:*", &app.ConfigureController{}, "get:DetailPage"),
				// 应用配置添加页面
				beego.NSRouter("/data/add", &app.DataController{}, "post:ConfigDataAdd"),
			),
			// 容器管理,
			// 容器入口,
			beego.NSRouter("/container/list", &app.AppController{}, "get:ContainerList"),
			// 容器详情,
			beego.NSRouter("/container/detail/:hi(.*)", &app.AppController{}, "get:ContainerDetail"),
			// Service管理,
			// 应用Service入口页面,
			beego.NSNamespace("/service",
				beego.NSRouter("/list", &app.ServiceController{}, "get:ServiceList"),
				// 应用Service添加页面,
				beego.NSRouter("/add", &app.ServiceController{}, "get:ServiceAdd"),
				// 应用Service 手动扩展页面,
				beego.NSRouter("/scale/add/:id:int", &app.ServiceController{}, "get:ScaleAdd"),
				// 应用Service 修改配置页面,
				beego.NSRouter("/config/add/:id:int", &app.ServiceController{}, "get:ConfigAdd"),
				// 应用Service 修改滚动升级页面,
				beego.NSRouter("/image/add/:id:int", &app.ServiceController{}, "get:ImageAdd"),
				// 应用Service 修改环境变量页面,
				beego.NSRouter("/env/add/:id:int", &app.ServiceController{}, "get:EnvAdd"),
				// 应用Service 修改健康检查页面,
				beego.NSRouter("/health/add/:id:int", &app.ServiceController{}, "get:HealthChange"),
				// 应用Service 修改端口页面,
				beego.NSRouter("/port/add/:id:int", &app.ServiceController{}, "get:PortChange"),
				// 创建服务添加存储页面,
				beego.NSRouter("/storage/add", &app.ServiceController{}, "get:StorageAdd"),
				// 创建服务添加健康检查页面,
				beego.NSRouter("/health/add", &app.ServiceController{}, "get:HealthAdd"),
				// 创建服务添加健康检查页面,
				beego.NSRouter("/configure/add", &app.ServiceController{}, "get:ConfigureAdd"),
			),
			beego.NSNamespace("/template",
				// 模板管理,
				// 应用模板入口页面,
				beego.NSRouter("/list", &app.AppController{}, "get:TemplateList"),
				// 应用模板入口页面,
				beego.NSRouter("/deploy/history", &app.AppController{}, "get:HistoryList"),
				// 应用模板添加页面,
				beego.NSRouter("/add", &app.AppController{}, "post:TemplateAdd"),
				// 应用模板添加页面,
				beego.NSRouter("/update/add", &app.AppController{}, "post:TemplateUpdateAdd"),
				// 应用模板拉起页面,
				beego.NSRouter("/deploy/add", &app.AppController{}, "post:TemplateDeployAdd"),
			),
			// 环境配置文件,
			beego.NSRouter("/evnfile", &app.AppController{}, "*:EnvfileList"),
		)

	pipelineNs :=
		beego.NewNamespace("/pipeline",
			// 流水线列表入口,
			beego.NSRouter("/history/list", &pipeline.ControllerPipeline{}, "get:PipelineHistoryList"),
			// 流水线列表入口,
			beego.NSRouter("/list", &pipeline.ControllerPipeline{}, "get:PipelineList"),
			// 流水线添加页面,
			beego.NSRouter("/add", &pipeline.ControllerPipeline{}, "get:PipelineAdd"),
			// 流水线详情页面,
			beego.NSRouter("/detail/:hi(.*)", &pipeline.ControllerPipeline{}, "get:PipelineDetail"),
		)

	registryNs :=
		beego.NewNamespace("/image",
			beego.NSNamespace("/sync",
				// 同步历史
				beego.NSRouter("/history", &registry.SyncController{}, "get:HistoryList"),
				// 镜像同步页面
				beego.NSRouter("/list", &registry.SyncController{}, "get:SyncList"),
				// 镜像同步添加
				beego.NSRouter("/add", &registry.SyncController{}, "get:SyncAdd"),
			),
			beego.NSNamespace("/registry",
				// 镜像中心,
				// 仓库配置入口页面,
				beego.NSRouter("/list", &registry.ImageController{}, "get:RegistryServerList"),
				// 仓库配置页面,
				beego.NSRouter("/add", &registry.ImageController{}, "get:RegistryServerAdd"),
				// 镜像中心,
				// 仓库权限入口页面,
				beego.NSNamespace("/perm",
					beego.NSRouter("/list", &registry.RegistryPermController{}, "get:RegistryPermList"),
					// 仓库权限配置页面,
					beego.NSRouter("/add", &registry.RegistryPermController{}, "get:RegistryPermAdd"),
				),
				// 基础镜像
				beego.NSNamespace("/base",
					beego.NSRouter("/list", &registry.BaseController{}, "get:BaseList"),
					// 仓库权限配置页面,
					beego.NSRouter("/add", &registry.BaseController{}, "get:BaseAdd"),
				),
				beego.NSNamespace("/group",
					// 仓库分组
					beego.NSRouter("/list", &registry.RegistryGroupController{}, "get:RegistryGroupList"),
					// 仓库分组添加
					beego.NSRouter("/add", &registry.RegistryGroupController{}, "get:RegistryGroupAdd"),
					// 仓库分组详情页面
					beego.NSRouter("/detail/:id:int", &registry.RegistryGroupController{}, "get:GroupDetailPage"),
				),
			),
		)

	// ci
	ciNs :=
		beego.NewNamespace("/ci",
			// dockerfile
			beego.NSNamespace("/dockerfile",
				// docker file详情入口
				beego.NSRouter("/detail/:hi(.*)", &ci.DockerFileController{}, "get:DockerFileDetail"),
				// docker file列表入口
				beego.NSRouter("/list", &ci.DockerFileController{}, "get:DockerFileList"),
				// docker file添加页面
				beego.NSRouter("/add", &ci.DockerFileController{}, "get:DockerFileAdd"),
			),
			// 代码仓库
			beego.NSNamespace("/code",
				// 代码仓库列表入口,
				beego.NSRouter("/list", &ci.CodeController{}, "get:CodeList"),
				// 代码仓库添加页面,
				beego.NSRouter("/add", &ci.CodeController{}, "get:CodeAdd"),
			),
			// 服务发布
			beego.NSNamespace("/service",
				// 服务发布列表入口,
				beego.NSRouter("/list", &ci.ServiceController{}, "get:ServiceList"),
				// 服务发布添加页面,/ci/service/release
				beego.NSRouter("/add", &ci.ServiceController{}, "get:ServiceAdd"),
				// 服务发布添加页面,/ci/service/release
				beego.NSRouter("/top/:id:int", &ci.ServiceController{}, "get:ServiceTop"),
				// 服务发布弹出页面
				beego.NSRouter("/release", &ci.ServiceController{}, "get:ServiceRelease"),
				// 服务发布弹出页面
				beego.NSRouter("/release/history", &ci.ServiceController{}, "get:HistoryList"),
				// 服务发布日志
				beego.NSRouter("/release/logs", &ci.ServiceController{}, "get:ServiceLog"),
				// 流量切入页面
				beego.NSRouter("/flow/:id:int", &ci.ServiceController{}, "get:StartFlow"),
				// 发布蓝色服务,滚动更新页面
				beego.NSRouter("/rolling/:id:int", &ci.ServiceController{}, "get:RollingUpdate"),
				beego.NSNamespace("/perm",
					beego.NSRouter("/list", &ci.CiPermController{}, "get:CiPermList"),
					// 仓库权限配置页面,
					beego.NSRouter("/add", &ci.CiPermController{}, "get:CiPermAdd"),
				),
			),
			// 构建任务
			beego.NSNamespace("/job",
				// 构建任务列表入口,
				beego.NSRouter("/history/list", &ci.JobController{}, "get:JobHistoryList"),
				// 构建任务列表入口,
				beego.NSRouter("/list", &ci.JobController{}, "get:JobList"),
				// 构建任务添加页面,
				beego.NSRouter("/add", &ci.JobController{}, "get:JobAdd"),
				// 构建日志页面
				beego.NSRouter("/logs/:id:int", &ci.JobController{}, "get:JobLogsPage"),
				// 构建任务详情页面,
				beego.NSRouter("/detail/:hi(.*)", &ci.JobController{}, "get:JobDetail"),
			),

		)

	// application
	applicationApi :=
		beego.NewNamespace("/api",
			// 获取容器数据
			beego.NSRouter("/container", &app.AppController{}, "get:ContainerData"),
			// 删除容器,容器删除后会重建
			beego.NSRouter("/container/:id:int", &app.AppController{}, "delete:ContainerDelete"),
			beego.NSNamespace("/pipeline",
				// 流水线保存,
				beego.NSRouter("", &pipeline.ControllerPipeline{}, "post:PipelineSave"),
				// 获取流水线数据所有数据,
				beego.NSRouter("", &pipeline.ControllerPipeline{}, "get:PipelineData"),
				// 删除配额,
				beego.NSRouter("/:id:int", &pipeline.ControllerPipeline{}, "delete:PipelineDelete"),
				// 获取流水线数据单条数据,
				beego.NSRouter("/:id:int", &pipeline.ControllerPipeline{}, "get:PipelineData"),
				// 执行流水线
				beego.NSRouter("/exec/:id:int", &pipeline.ControllerPipeline{}, "get:PipelineExec"),
				// 流水线历史日子
				beego.NSRouter("/history", &pipeline.ControllerPipeline{}, "get:PipelineHistoryData"),
			),
			beego.NSNamespace("/app",
				// 获取应用名称,
				beego.NSRouter("/name", &app.AppController{}, "get:GetAppName"),
				// 重新部署应用,
				beego.NSRouter("/redeploy", &app.AppController{}, "get:RedeployApp"),
				// 删除应用,
				beego.NSRouter("/:id:int", &app.AppController{}, "delete:AppDelete"),
				// 启停应用接口,
				beego.NSRouter("/scale/:id:int", &app.AppController{}, "*:AppScale"),
				// 应用数据,
				beego.NSRouter("", &app.AppController{}, "get:AppData"),
				beego.NSRouter("/:id:int", &app.AppController{}, "get:AppData"),
			),
			beego.NSNamespace("/template",
				// 应用模板保存,
				beego.NSRouter("", &app.AppController{}, "post:TemplateSave"),
				// 更细模板yaml
				beego.NSRouter("/update", &app.AppController{}, "post:TemplateUpdate"),
				// 拉起应用
				beego.NSRouter("/deploy/:id:int", &app.AppController{}, "post:StartDeploy"),
				// 拉起应用历史
				beego.NSRouter("/deploy/history", &app.AppController{}, "post:HistoryData"),
				// 获取应用模板数据所有数据,
				beego.NSRouter("", &app.AppController{}, "get:TemplateData"),
				// 获取应用模板数据所有数据的名称和ID,
				beego.NSRouter("/name", &app.AppController{}, "get:GetTemplateName"),
				// 删除模板,
				beego.NSRouter("/:id:int", &app.AppController{}, "delete:TemplateDelete"),
				// 获取应用模板数据单条数据,
				beego.NSRouter("/:id:int", &app.AppController{}, "get:TemplateData"),
				// 检查yaml文件是否可以转换到json,
				beego.NSRouter("/yaml/check", &app.AppController{}, "post:YamlCheck"),
			),
			beego.NSNamespace("/service",
				// 获取应用Service数据所有数据,
				beego.NSRouter("", &app.ServiceController{}, "get:ServiceData"),
				// 获取应用Service数据所有数据的名称和ID,
				beego.NSRouter("/name", &app.ServiceController{}, "get:GetServiceName"),
				// 删除Service,
				beego.NSRouter("/:id:int", &app.ServiceController{}, "delete:ServiceDelete"),
				// 获取应用Service数据单条数据,
				beego.NSRouter("/:id:int", &app.ServiceController{}, "get:ServiceData"),
				// 扩容或缩容, 传参为 服务的 id  ?replicas=1,
				beego.NSRouter("/scale/:id:int", &app.ServiceController{}, "*:ServiceScale"),
				// 应用Service保存,
				beego.NSRouter("", &app.ServiceController{}, "post:ServiceSave"),
				// 更新服务信息,cpu,内存,环境变量等可更新的数据,
				beego.NSRouter("/update/:id:int", &app.ServiceController{}, "post:ServiceUpdate"),
			),
			beego.NSNamespace("/configure",
				// 应用配置保存,
				beego.NSRouter("", &app.ConfigureController{}, "post:ConfigureSave"),
				// 获取应用配置数据所有数据,
				beego.NSRouter("", &app.ConfigureController{}, "get:ConfigureData"),
				// 获取应用配置数据所有数据的名称和ID,
				beego.NSRouter("/name", &app.ConfigureController{}, "get:GetConfigureName"),
				// 删除配置,
				beego.NSRouter("/:id:int", &app.ConfigureController{}, "delete:ConfigureDelete"),
				// 获取应用配置数据单条数据,
				beego.NSRouter("/:id:int", &app.ConfigureController{}, "get:ConfigureData"),
				// 应用配置保存,
				beego.NSRouter("/data", &app.DataController{}, "post:ConfigDataSave"),
				// 获取应用配置数据所有数据,
				beego.NSRouter("/data", &app.DataController{}, "get:ConfigData"),
				// 获取应用配置数据所有数据的名称和ID,
				beego.NSRouter("/data/name", &app.DataController{}, "get:GetConfigDataName"),
				// 删除配置,
				beego.NSRouter("/data/:id:int", &app.DataController{}, "delete:ConfigDataDelete"),
				// 获取应用配置数据单条数据,
				beego.NSRouter("/data/:id:int", &app.DataController{}, "get:ConfigData"),
			),
		)

	ciApi :=
		beego.NewNamespace("/api/ci",
			beego.NSNamespace("/dockerfile",
				beego.NSRouter("", &ci.DockerFileController{}, "get:DockerFileDatas"),
				// 获取docker file信息
				beego.NSRouter("/name", &ci.DockerFileController{}, "get:DockerFileDataName"),
				// 获取docker file数据
				beego.NSRouter("", &ci.DockerFileController{}, "get:DockerFileDatas"),
				// docker file保存
				beego.NSRouter("", &ci.DockerFileController{}, "post:DockerFileSave"),
				// 获取docker file数据所有数据
				beego.NSRouter("", &ci.DockerFileController{}, "get:DockerFileData"),
				// 删除部门团队
				beego.NSRouter("/:id:int", &ci.DockerFileController{}, "delete:DockerFileDelete"),
				// 获取docker file数据单条数据
				beego.NSRouter("/:id:int", &ci.DockerFileController{}, "get:DockerFileData"),

			),
			beego.NSNamespace("/code",
				// 代码仓库保存,
				beego.NSRouter("", &ci.CodeController{}, "post:CodeSave"),
				// 获取代码仓库数据所有数据,
				beego.NSRouter("", &ci.CodeController{}, "get:CodeDatas"),
				// 删除代码仓库,
				beego.NSRouter("/:id:int", &ci.CodeController{}, "delete:CodeDelete"),
				// 获取代码仓库数据单条数据,
				beego.NSRouter("/:id:int", &ci.CodeController{}, "get:CodeData"),
				// 获取代码仓库信息,
				beego.NSRouter("/name", &ci.CodeController{}, "get:CodeDataName"),
				// 获取代码仓库数据,
				beego.NSRouter("", &ci.CodeController{}, "get:CodeDatas"),
			),
			beego.NSNamespace("/service",
				// 服务发布保存,
				beego.NSRouter("", &ci.ServiceController{}, "post:ServiceSave"),
				// 获取服务发布数据所有数据,
				beego.NSRouter("", &ci.ServiceController{}, "get:ServiceDatas"),
				// 删除服务发布,
				beego.NSRouter("/:id:int", &ci.ServiceController{}, "delete:ServiceDelete"),
				// 获取服务发布信息,
				beego.NSRouter("/name", &ci.ServiceController{}, "get:ServiceDataName"),
				// 获取服务发布数据,
				beego.NSRouter("", &ci.ServiceController{}, "get:ServiceDatas"),
				// 发布历史数据
				beego.NSRouter("/history", &ci.ServiceController{}, "get:ReleaseHistory"),
				// 获取服务发布操作日志
				beego.NSRouter("/logs", &ci.ServiceController{}, "get:ServiceLogs"),
				// 执行服务发布
				beego.NSRouter("/release/:id:int", &ci.ServiceController{}, "post:ServiceReleaseExec"),
				// 修改发布历史 history
				beego.NSRouter("/history/:id:int", &ci.ServiceController{}, "post:UpdateHistory"),
				// 执行服务下线 online
				beego.NSRouter("/release/:id:int", &ci.ServiceController{}, "delete:ServiceOffline"),
				// 上线服务
				beego.NSRouter("/online/:id:int", &ci.ServiceController{}, "post:ServiceOnline"),
				// 回滚服务
				beego.NSRouter("/rollback/:id:int", &ci.ServiceController{}, "post:ServiceRollback"),
				// 发布蓝色服务,将老的应用更新到新的镜像
				beego.NSRouter("/blue/:id:int", &ci.ServiceController{}, "post:UpdateBlueService"),
				// 发布蓝色服务,将老的应用更新到新的镜像
				beego.NSRouter("/flow/:id:int", &ci.ServiceController{}, "post:StartFlowExec"),
				// 滚动更新服务,蓝版的
				beego.NSRouter("/rolling/:id:int", &ci.ServiceController{}, "post:RollingUpdateExec"),
				// 发布权限
				beego.NSNamespace("/perm",
					// 发布权限配置保存,
					beego.NSRouter("", &ci.CiPermController{}, "post:CiPermSave"),
					// 获取权限仓库配置数据所有数据,
					beego.NSRouter("", &ci.CiPermController{}, "get:CiPerm"),
					// 删除权限配置,
					beego.NSRouter("/:id:int", &ci.CiPermController{}, "delete:CiPermDelete"),
					// 获取权限仓库配置数据单条数据,
					beego.NSRouter("/:id:int", &ci.CiPermController{}, "get:CiPerm"),
				),
			),
			beego.NSNamespace("/job",
				// 构建任务保存,
				beego.NSRouter("", &ci.JobController{}, "post:JobSave"),
				// 获取构建任务数据所有数据,
				beego.NSRouter("", &ci.JobController{}, "get:JobDatas"),
				// 删除构建任务,
				beego.NSRouter("/:id:int", &ci.JobController{}, "delete:JobDelete"),
				// 获取构建任务数据单条数据,
				beego.NSRouter("/:id:int", &ci.JobController{}, "get:JobData"),
				// 获取构建任务信息,
				beego.NSRouter("/name", &ci.JobController{}, "get:JobDataName"),
				// 获取构建任务数据,
				beego.NSRouter("/history", &ci.JobController{}, "get:JobHistoryDatas"),
				// 获取历史数据
				beego.NSRouter("", &ci.JobController{}, "get:JobDatas"),
				// 执行构建任务
				beego.NSRouter("/exec/:id:int", &ci.JobController{}, "get:JobExec"),
				// 执行构建任务
				beego.NSRouter("/logs/:id:int", &ci.JobController{}, "get:JobLogs"),
				// 获取构建的dockerfile
				beego.NSRouter("/dockerfile/:id:int", &ci.JobController{}, "get:JobDockerfile"),
				// 获取docker file数据单条数据
				beego.NSRouter("/dockerfile/:hi(.*)", &ci.JobController{}, "get:JobDockerfile"),
			),
		)

	// 仓库管理
	registryApi :=
		beego.NewNamespace("/api/",
			beego.NSNamespace("/image",
				beego.NSNamespace("/sync",
					// 审批通过
					beego.NSRouter("/approved/:id:int", &registry.SyncController{}, "post:ApprovedSave"),
					// 保存镜像同步请求
					beego.NSRouter("", &registry.SyncController{}, "post:SyncSave"),
					// 保存镜像同步请求
					beego.NSRouter("", &registry.SyncController{}, "get:SyncDatas"),
					// 执行镜像同步
					beego.NSRouter("/:id:int", &registry.SyncController{}, "get:SyncExec"),
					// 保存镜像同步请求
					beego.NSRouter("/:id:int", &registry.SyncController{}, "delete:SyncDelete"),
					beego.NSNamespace("/history",
						beego.NSRouter("", &registry.SyncController{}, "get:HistorDatas"),
					),
				),
			),
			beego.NSNamespace("/registry",
				// 仓库配置保存,
				beego.NSRouter("", &registry.ImageController{}, "post:RegistryServerSave"),
				// 获取仓库配置数据所有数据,
				beego.NSRouter("", &registry.ImageController{}, "get:RegistryServer"),
				// 删除配置,
				beego.NSRouter("/:id:int", &registry.ImageController{}, "delete:RegistryServerDelete"),
				// 获取仓库配置数据单条数据,
				beego.NSRouter("/:id:int", &registry.ImageController{}, "get:RegistryServer"),
				// 重新部署仓库服务器
				beego.NSRouter("/recreate", &registry.ImageController{}, "post:RecreateRegistry"),
				// 权限配置
				beego.NSNamespace("/perm",
					// 仓库权限配置保存,
					beego.NSRouter("", &registry.RegistryPermController{}, "post:RegistryPermSave"),
					// 获取权限仓库配置数据所有数据,
					beego.NSRouter("", &registry.RegistryPermController{}, "get:RegistryPerm"),
					// 删除权限配置,
					beego.NSRouter("/:id:int", &registry.RegistryPermController{}, "delete:RegistryPermDelete"),
					// 获取权限仓库配置数据单条数据,
					beego.NSRouter("/:id:int", &registry.RegistryPermController{}, "get:RegistryPerm"),
				),
				// 基础镜像管理
				beego.NSNamespace("/base",
					// 仓库基础镜像配置保存,
					beego.NSRouter("", &registry.BaseController{}, "post:BaseSave"),
					// 获取基础镜像仓库配置数据所有数据,
					beego.NSRouter("", &registry.BaseController{}, "get:Base"),
					// 删除基础镜像配置,
					beego.NSRouter("/:id:int", &registry.BaseController{}, "delete:BaseDelete"),
					// 获取基础镜像仓库配置数据单条数据,
					beego.NSRouter("/:id:int", &registry.BaseController{}, "get:Base"),
				),
				beego.NSNamespace("/group",
					// 仓库配置保存,
					beego.NSRouter("", &registry.RegistryGroupController{}, "post:SaveRegistryGroup"),
					// 获取仓库配置数据所有数据,
					beego.NSRouter("", &registry.RegistryGroupController{}, "get:RegistryGroup"),
					// 删除配置,
					beego.NSRouter("/:id:int", &registry.RegistryGroupController{}, "delete:DeleteRegistryGroup"),
					// 获取仓库配置数据单条数据,
					beego.NSRouter("/:id:int", &registry.RegistryGroupController{}, "get:RegistryGroup"),
					// 获取仓库配置数据所有数据,
					beego.NSRouter("/images", &registry.RegistryGroupController{}, "get:RegistryGroupImages"),
					// 获取用户操作镜像日志
					beego.NSRouter("/images/log", &registry.RegistryGroupController{}, "get:RegistryImagesLog"),
					// 获取单个镜像
					beego.NSRouter("/images/:id:int", &registry.RegistryGroupController{}, "get:GetRegistryGroupImage"),
					// 获取单个镜像
					beego.NSRouter("/images/:hi(.*)", &registry.RegistryGroupController{}, "get:GetRegistryGroupImage"),
					// 删除镜像
					beego.NSRouter("/images/:id:int", &registry.RegistryGroupController{}, "delete:DeleteRegistryGroupImage"),
				),
				// 08-2-7 11:20
				// 在安装应用时候选择的镜像数据
				beego.NSNamespace("/deploy",
					beego.NSRouter("/images", &registry.RegistryGroupController{}, "get:GetDeployImage"),
				),
			),
		)

	// 基础设施配置
	clusterApi :=
		beego.NewNamespace("/api/",
			beego.NSNamespace("/cluster",
				// 获取集群的节点,
				beego.NSRouter("/nodes", &cluster.ClusterController{}, "*:NodesData"),
				// 保存集群,
				beego.NSRouter("", &cluster.ClusterController{}, "post:Save"),
				// 删除集群,
				beego.NSRouter("/:id:int", &cluster.ClusterController{}, "delete:Delete"),
				// 集群数据,
				beego.NSRouter("", &cluster.ClusterController{}, "get:ClusterData"),
				// 集群数据,直返回,集群名称和id的数据,
				beego.NSRouter("/name", &cluster.ClusterController{}, "get:ClusterName"),
				// 单条数据,
				beego.NSRouter("/:id:int", &cluster.ClusterController{}, "get:ClusterData"),
				// 保存标签,
				beego.NSRouter("/label", &hosts.HostsController{}, "post:LabelSave"),
				beego.NSNamespace("/hosts",
					// 保存主机,
					beego.NSRouter("", &hosts.HostsController{}, "post:Save"),
					// 删除主机,
					beego.NSRouter("/:id:int", &hosts.HostsController{}, "delete:Delete"),
					// 主机数据,
					beego.NSRouter("", &hosts.HostsController{}, "get:HostsData"),
					// 单条主机数据,
					beego.NSRouter("/:id:int", &hosts.HostsController{}, "get:HostsData"),
					// 调度维护设置
					beego.NSRouter("/:id:int", &hosts.HostsController{}, "post:Schedulable"),
					// 调度维护设置
					beego.NSRouter("/images/:id:int", &hosts.HostsController{}, "get:GetHostImages"),
				),
			),
			beego.NSNamespace("/storage",
				// 保存存储卷,
				beego.NSRouter("", &storage.StorageController{}, "post:StorageSave"),
				// 删除存储卷,
				beego.NSRouter("/:id:int", &storage.StorageController{}, "delete:StorageDelete"),
				// 存储卷数据,
				beego.NSRouter("", &storage.StorageController{}, "get:StorageData"),
				// 单条数据,
				beego.NSRouter("/:id:int", &storage.StorageController{}, "get:StorageData"),
				beego.NSNamespace("/server",
					// 保存存储服务器
					beego.NSRouter("", &storage.StorageServerController{}, "post:StorageServerSave"),
					// 删除存储服务器
					beego.NSRouter("/:id:int", &storage.StorageServerController{}, "delete:StorageServerDelete"),
					// 存储服务器数据
					beego.NSRouter("", &storage.StorageServerController{}, "get:StorageServerData"),
					// 单条数据,
					beego.NSRouter("/:id:int", &storage.StorageServerController{}, "get:StorageServerData"),
				),
			),
		)

	//// 存储管理
	//// 存储卷
	//beego.Router("/storage/volume/list", &cluster.ClusterController{}, "get:List")
	//// 快照
	//beego.Router("/storage/snapshot/list", &cluster.ClusterController{}, "get:List")

	// 基础设施
	baseNs :=
		beego.NewNamespace("/base",
			beego.NSNamespace("/quota",
				// 资源配额入口页面,
				beego.NSRouter("/list", &quota.ControllerQuota{}, "get:QuotaList"),
				// 资源配额添加页面,
				beego.NSRouter("/add", &quota.ControllerQuota{}, "post:QuotaAdd"),
				// 集群配额详情页面,
				beego.NSRouter("/detail/:id:int", &quota.ControllerQuota{}, "get:QuotaDetailPage"),
				beego.NSRouter("/detail/:hi(.*)", &quota.ControllerQuota{}, "get:QuotaDetailPage"),
			),
			beego.NSNamespace("/storage",
				// 基础设施,
				// 存储管理,
				beego.NSRouter("/list", &storage.StorageController{}, "get:StorageList"),
				// 存储添加页面,
				beego.NSRouter("/add", &storage.StorageController{}, "get:StorageAdd"),
				beego.NSNamespace("/server",
					// 存储服务器配置
					beego.NSRouter("/list", &storage.StorageServerController{}, "get:StorageServerList"),
					// 存储服务器添加
					beego.NSRouter("/add", &storage.StorageServerController{}, "get:StorageServerAdd"),
				),
			),
			beego.NSNamespace("/network",
				// 基础设施,
				// 网络管理,
				beego.NSNamespace("/lb",
					beego.NSRouter("/list", &lb.LbController{}, "get:LbList"),
					// 资源负载均衡器添加页面,
					beego.NSRouter("/add", &lb.LbController{}, "get:LbAdd"),
					// 集群负载均衡器详情页面,
					beego.NSRouter("/detail/:id:int", &lb.LbController{}, "get:LbDetailPage"),
					// 负载均衡服务添加页面,
					beego.NSRouter("/service/add", &lb.ServiceController{}, "get:ServiceAdd"),
				),
				// 证书管理
				beego.NSNamespace("/cert",
					// 证书入口
					beego.NSRouter("/list", &lb.CertController{}, "get:CertList"),
					// 证书添加页面,
					beego.NSRouter("/add", &lb.CertController{}, "get:CertAdd"),
				),
			),
			// 服务使用量,
			beego.NSRouter("/service/list", &cluster.ClusterController{}, "get:List"),
			beego.NSNamespace("/cluster",
				// 基础设施管理,
				// 集群管理,
				beego.NSRouter("/list", &cluster.ClusterController{}, "get:List"),
				// 主机镜像详情弹出页面
				beego.NSRouter("/image/:id:int", &cluster.ClusterController{}, "get:Images"),
				// 添加集群页面,
				beego.NSRouter("/add", &cluster.ClusterController{}, "*:Add"),
				// 集群详情数据页面,
				beego.NSRouter("/detail/:hi(.*)", &cluster.ClusterController{}, "get:DetailPage"),
				// 基础设施管理,
				// 集群主机管理,
				beego.NSRouter("/hosts/list", &hosts.HostsController{}, "get:List"),
				beego.NSRouter("/hosts/add", &hosts.HostsController{}, "*:Add"),
				// 添加主机标签,
				beego.NSRouter("/label/add", &hosts.HostsController{}, "*:LabelAdd"),
			),
		)

	baseApi :=
		beego.NewNamespace("/api",
			beego.NSNamespace("/quota",
				// 资源配额保存,
				beego.NSRouter("", &quota.ControllerQuota{}, "post:QuotaSave"),
				// 获取资源配额数据所有数据,
				beego.NSRouter("", &quota.ControllerQuota{}, "get:QuotaData"),
				// 删除配额,
				beego.NSRouter("/:id:int", &quota.ControllerQuota{}, "delete:QuotaDelete"),
				// 获取资源配额数据单条数据,
				beego.NSRouter("/:id:int", &quota.ControllerQuota{}, "get:QuotaData"),
				// 获取配额名称,
				beego.NSRouter("/name", &quota.ControllerQuota{}, "get:GetQuotaName"),
			),
			beego.NSNamespace("/lb",
				// 资源负载均衡器保存,
				beego.NSRouter("", &lb.LbController{}, "post:LbSave"),
				// 获取资源负载均衡器数据所有数据,
				beego.NSRouter("", &lb.LbController{}, "get:LbData"),
				// 删除负载均衡器,
				beego.NSRouter("/:id:int", &lb.LbController{}, "delete:LbDelete"),
				// 获取资源负载均衡器数据单条数据,
				beego.NSRouter("/:id:int", &lb.LbController{}, "get:LbData"),
				// 返回负载均衡服务所有数据,
				beego.NSRouter("/service/:hi(.*)", &lb.ServiceController{}, "get:ServiceData"),
			),
			beego.NSNamespace("/network",
				beego.NSNamespace("/lb",
					beego.NSNamespace("/service",
						// 负载均衡服务管理,
						beego.NSRouter("/:id:int", &lb.ServiceController{}, "delete:ServiceDelete"),
						// 负载均衡服务保存,
						beego.NSRouter("", &lb.ServiceController{}, "post:ServiceSave"),
					),
					beego.NSNamespace("/nginx",
						// 获取nginx配置
						beego.NSRouter("/:id:int", &lb.ServiceController{}, "get:GetNginxConf"),
						// 保存nginx配置
						beego.NSRouter("/:id:int", &lb.ServiceController{}, "post:SaveNginxConf"),
					),
					beego.NSRouter("/domain", &lb.ServiceController{}, "get:GetLbDomain"),
				),
				beego.NSNamespace("/cert",
					// 负载均衡服务管理,
					beego.NSRouter("/:id:int", &lb.CertController{}, "delete:CertDelete"),
					// 负载均衡服务保存,
					beego.NSRouter("", &lb.CertController{}, "post:CertSave"),
					// 获取资源负载均衡器数据所有数据,
					beego.NSRouter("", &lb.CertController{}, "get:CertData"),
				),
			),
		)

	systemNs :=
		beego.NewNamespace("/system",
			beego.NSNamespace("/users",
				// 用户管理,
				// 部门团队管理,
				beego.NSNamespace("/groups",
					// 应用部门团队入口页面,
					beego.NSRouter("/list", &users.GroupsController{}, "get:GroupsList"),
					// 应用部门团队添加页面,
					beego.NSRouter("/add", &users.GroupsController{}, "get:GroupsAdd"),
				),
				// 用户管理,
				beego.NSNamespace("/user",
					// 用户列表入口,
					beego.NSRouter("/list", &users.UserController{}, "get:UserList"),
					// 用户添加页面,
					beego.NSRouter("/add", &users.UserController{}, "get:UserAdd"),
				),
			),
			beego.NSNamespace("/ent",
				// 环境入口列表
				beego.NSRouter("/list", &ent.EntController{}, "get:EntList"),
				// 环境添加页面,
				beego.NSRouter("/add", &ent.EntController{}, "get:EntAdd"),
			),
			beego.NSNamespace("/operlog",
				// 日志入口
				beego.NSRouter("/list", &operlog.LogController{}, "get:OperlogList"),
			),
			beego.NSNamespace("/perm",
				// 权限入口
				beego.NSRouter("/list", &perm.PermController{}, "get:PermList"),
				// 权限添加
				beego.NSRouter("/add", &perm.PermController{}, "get:PermAdd"),
				beego.NSNamespace("/role",
					// 权限入口
					beego.NSRouter("/list", &perm.PermRoleController{}, "get:PermRoleList"),
					// 权限添加
					beego.NSRouter("/add", &perm.PermRoleController{}, "get:PermRoleAdd"),
				),
				// api资源配置
				beego.NSNamespace("/resource",
					// 权限入口
					beego.NSRouter("/list", &perm.ResourceController{}, "get:ResourceList"),
					// 权限添加
					beego.NSRouter("/add", &perm.ResourceController{}, "get:ResourceAdd"),
				),
			),
		)

	// 系统设置
	systemApi :=
		beego.NewNamespace("/api",
			beego.NSNamespace("/groups",
				// 应用部门团队保存,
				beego.NSRouter("", &users.GroupsController{}, "post:GroupsSave"),
				// 获取应用部门团队数据所有数据,
				beego.NSRouter("", &users.GroupsController{}, "get:GroupsData"),
				// 获取应用部门团队数据所有数据的名称和ID,
				beego.NSRouter("/name", &users.GroupsController{}, "get:GetGroupsName"),
				// 获取组名map,
				beego.NSRouter("/map", &users.GroupsController{}, "get:GetGroupsMap"),
				// 删除部门团队,
				beego.NSRouter("/:id:int", &users.GroupsController{}, "delete:GroupsDelete"),
				// 获取应用部门团队数据单条数据,
				beego.NSRouter("/:id:int", &users.GroupsController{}, "get:GroupsData"),
			),
			beego.NSNamespace("/users",
				// 用户保存,
				beego.NSRouter("", &users.UserController{}, "post:UserSave"),
				//// 获取用户数据所有数据,
				//beego.NSRouter("", &users.UserController{}, "get:UserData"),
				// 删除部门团队,
				beego.NSRouter("/:id:int", &users.UserController{}, "delete:UserDelete"),
				// 获取用户数据单条数据,
				beego.NSRouter("/:id:int", &users.UserController{}, "get:UserData"),
				// 获取用户信息,
				beego.NSRouter("/name", &users.UserController{}, "get:UserDataName"),
				// 获取用户数据,
				beego.NSRouter("", &users.UserController{}, "get:UserDatas"),
			),
			beego.NSNamespace("/operlog",
				beego.NSRouter("", &operlog.LogController{}, "*:OperlogDatas"),
			),
			beego.NSNamespace("/ent",
				// 环境保存,
				beego.NSRouter("", &ent.EntController{}, "post:EntSave"),
				// 删除环境
				beego.NSRouter("/:id:int", &ent.EntController{}, "delete:EntDelete"),
				// 获取环境数据单条数据,
				beego.NSRouter("/:id:int", &ent.EntController{}, "get:EntData"),
				// 获取环境信息,
				beego.NSRouter("/name", &ent.EntController{}, "get:EntDataName"),
				// 获取环境数据,
				beego.NSRouter("", &ent.EntController{}, "get:EntDatas"),
			),
			beego.NSNamespace("/perm",
				// 权限保存,
				beego.NSRouter("", &perm.PermController{}, "post:PermSave"),
				// 删除权限,
				beego.NSRouter("/:id:int", &perm.PermController{}, "delete:PermDelete"),
				// 获取权限数据单条数据,
				beego.NSRouter("/:id:int", &perm.PermController{}, "get:PermData"),
				// 获取权限信息,
				beego.NSRouter("/name", &perm.PermController{}, "get:PermDataName"),
				// 获取权限数据,
				beego.NSRouter("", &perm.PermController{}, "get:PermDatas"),
				beego.NSNamespace("/role",
					// 角色保存,
					beego.NSRouter("", &perm.PermRoleController{}, "post:PermRoleSave"),
					// 删除角色,
					beego.NSRouter("/:id:int", &perm.PermRoleController{}, "delete:PermRoleDelete"),
					// 获取角色数据单条数据,
					beego.NSRouter("/:id:int", &perm.PermRoleController{}, "get:PermRoleData"),
					// 获取角色信息,
					beego.NSRouter("/name", &perm.PermRoleController{}, "get:PermRoleDataName"),
					// 获取角色数据,
					beego.NSRouter("", &perm.PermRoleController{}, "get:PermRoleDatas"),
				),
				beego.NSNamespace("/resource",
					// api资源保存,
					beego.NSRouter("", &perm.ResourceController{}, "post:ResourceSave"),
					// 删除api资源,
					beego.NSRouter("/:id:int", &perm.ResourceController{}, "delete:ResourceDelete"),
					// 获取api资源数据单条数据,
					beego.NSRouter("/:id:int", &perm.ResourceController{}, "get:ResourceData"),
					// 获取api资源信息,
					beego.NSRouter("/name", &perm.ResourceController{}, "get:ResourceDataName"),
					// 获取api资源数据,
					beego.NSRouter("", &perm.ResourceController{}, "get:ResourceDatas"),
				),
			),
		)

	// 2018-02-19 18:30
	// 监控中心
	monitorNs :=
		beego.NewNamespace("/monitor",
			beego.NSNamespace("/scale",
				// 环境入口列表
				beego.NSRouter("/list", &monitor.AutoScaleController{}, "get:AutoScaleList"),
				// 环境添加页面,
				beego.NSRouter("/add", &monitor.AutoScaleController{}, "get:AutoScaleAdd"),
				// 发布日志
				beego.NSRouter("/logs", &monitor.AutoScaleController{}, "get:AutoScaleLogs"),
			),
		)

	// 2018-02-19 18:32
	// 监控中心api
	monitorApi :=
		beego.NewNamespace("/api/",
			beego.NSNamespace("/monitor",
				beego.NSNamespace("/scale",
					// 自动伸缩保存,
					beego.NSRouter("", &monitor.AutoScaleController{}, "post:AutoScaleSave"),
					// 删除自动伸缩
					beego.NSRouter("/:id:int", &monitor.AutoScaleController{}, "delete:AutoScaleDelete"),
					// 获取自动伸缩数据单条数据,
					beego.NSRouter("/:id:int", &monitor.AutoScaleController{}, "get:AutoScaleData"),
					// 获取自动伸缩数据,
					beego.NSRouter("", &monitor.AutoScaleController{}, "get:AutoScaleDatas"),
					// 获取自动伸缩数据,
					beego.NSRouter("/logs", &monitor.AutoScaleController{}, "get:AutoScaleLogsData"),
				),
			),
		)

	beego.AddNamespace(monitorApi)
	beego.AddNamespace(monitorNs)
	beego.AddNamespace(pipelineNs)
	beego.AddNamespace(ciNs)
	beego.AddNamespace(registryNs)
	beego.AddNamespace(systemNs)
	beego.AddNamespace(baseNs)
	beego.AddNamespace(ciApi)
	beego.AddNamespace(applicationNs)
	beego.AddNamespace(applicationApi)
	beego.AddNamespace(registryApi)
	beego.AddNamespace(clusterApi)
	beego.AddNamespace(baseApi)
	beego.AddNamespace(systemApi)

	// 过滤器功能实现,拦截未登陆请求
	var FilterUser = func(ctx *context.Context) {
		uri := ctx.Request.RequestURI
		if !strings.Contains(uri, "/static/") && !strings.Contains(uri, "/api/user/login") {
			_, ok := ctx.Input.Session("username").(string)
			uris := strings.Split(uri, "?referer=/")
			if !ok && uri != "/login" && uris[0] != "/login" {
				logs.Error("用户未登陆,请求URL为", uri, ctx.Request.RemoteAddr)
				url := util.GetUri(*ctx)
				ctx.Redirect(302, "/login?referer="+url)
			}
		}
	}
	beego.InsertFilter("/*", beego.BeforeRouter, FilterUser)
}
