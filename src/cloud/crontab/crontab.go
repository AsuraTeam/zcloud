package crontab

import (
	"github.com/jakecoffman/cron"
	"cloud/controllers/docker/application/app"
	"cloud/controllers/base/hosts"
	"cloud/controllers/image"
	"cloud/controllers/base/cluster"
	"cloud/controllers/monitor"
	"github.com/astaxie/beego"
)

// 自动刷新数据任务计划
func CronStart() {

	// 如果配置文件没有开启就不跑
	if beego.AppConfig.String("cron") != "true" {
		return
	}

	cron := cron.New()
	// 应用缓存,每2分钟
	cron.AddFunc("1 */2 * * * ?", func() {
		app.CacheAppData()
	}, "CacheAppData")
	// 容器
	cron.AddFunc("30 */3 * * * ?", func() {
		app.MakeContainerData("")
	}, "MakeContainerData")
	// node状态写入到缓存
	cron.AddFunc("20 */10 * * * ?", func() {
		hosts.CronCache()
	}, "NodeStatusCache")
	// 服务数据写入到缓存
	cron.AddFunc("40 */2 * * * ?", func() {
		app.CronServiceCache()
	}, "CronServiceCache")
	// 仓库镜像写入缓存
	cron.AddFunc("1 */3 * * * ?", func() {
		registry.UpdateGroupImageInfo()
	}, "UpdateGroupImageInfo")
	// 集群数据写入缓存
	cron.AddFunc("1 */5 * * * ?", func() {
		cluster.CacheClusterData()
	}, "CacheClusterData")
	// 监控自动扩容
	cron.AddFunc("*/30 * * * * ?", func() {
		monitor.CronAutoScale()
	}, "CronAutoScale")
	cron.Start()
}
