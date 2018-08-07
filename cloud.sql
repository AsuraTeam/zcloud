-- MySQL dump 10.14  Distrib 5.5.52-MariaDB, for Linux (x86_64)
--
-- Host: localhost    Database: cloud
-- ------------------------------------------------------
-- Server version	5.5.52-MariaDB

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Current Database: `cloud`
--

CREATE DATABASE /*!32312 IF NOT EXISTS*/ `cloud` /*!40100 DEFAULT CHARACTER SET latin1 */;

USE `cloud`;

--
-- Table structure for table `cloud_api_resource`
--

DROP TABLE IF EXISTS `cloud_api_resource`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `cloud_api_resource` (
  `resource_id` int(11) NOT NULL AUTO_INCREMENT,
  `create_time` varchar(32) DEFAULT NULL COMMENT '创建时间',
  `create_user` varchar(32) DEFAULT NULL COMMENT '创建用户',
  `last_modify_time` varchar(32) DEFAULT NULL COMMENT '最近修改时间',
  `last_modify_user` varchar(32) DEFAULT NULL COMMENT '最近修改用户',
  `description` varchar(300) DEFAULT NULL COMMENT '服务描述信息',
  `api_url` varchar(200) DEFAULT NULL COMMENT 'apiurl地址',
  `name` varchar(100) DEFAULT NULL COMMENT 'api名称',
  `api_type` varchar(10) DEFAULT NULL COMMENT '是否是公开的,公开的将不受权限控制',
  PRIMARY KEY (`resource_id`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `cloud_api_resource`
--

LOCK TABLES `cloud_api_resource` WRITE;
/*!40000 ALTER TABLE `cloud_api_resource` DISABLE KEYS */;
INSERT INTO `cloud_api_resource` VALUES (3,'2018-02-08 17:50:14','zhaozq14','2018-02-08 17:50:14','zhaozq14','保存应用','/api/app','应用保存','on');
/*!40000 ALTER TABLE `cloud_api_resource` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `cloud_app`
--

DROP TABLE IF EXISTS `cloud_app`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `cloud_app` (
  `app_id` int(11) NOT NULL AUTO_INCREMENT,
  `app_name` varchar(36) DEFAULT NULL COMMENT '应用名称',
  `app_type` varchar(32) DEFAULT NULL COMMENT '应用类型',
  `status` varchar(20) DEFAULT NULL COMMENT '运行状态',
  `resource_name` varchar(32) DEFAULT NULL COMMENT '资源空间',
  `create_time` varchar(32) DEFAULT NULL COMMENT '创建时间',
  `create_user` varchar(32) DEFAULT NULL COMMENT '创建用户',
  `last_modify_time` varchar(32) DEFAULT NULL COMMENT '最近修改时间',
  `last_modify_user` varchar(32) DEFAULT NULL COMMENT '最近修改用户',
  `app_labels` text,
  `json_data` text COMMENT '其他非固定数据存储',
  `last_update_time` varchar(32) DEFAULT NULL COMMENT '最近更新时间',
  `yaml` text COMMENT '编排畏惧',
  `cluster_name` varchar(40) DEFAULT NULL COMMENT '集群名称',
  `is_service` varchar(3) DEFAULT NULL,
  `uuid` varchar(35) DEFAULT NULL,
  `service_yaml` text COMMENT '创建服务的yaml内容',
  `node_port` varchar(500) DEFAULT NULL COMMENT '存放服务的端口,json格式 {name:xxxx,port:8080}',
  `entname` varchar(100) DEFAULT NULL COMMENT '环境名称',
  PRIMARY KEY (`app_id`),
  UNIQUE KEY `uidx_app_name_cluster_name` (`app_name`,`cluster_name`),
  KEY `idx_cloud_app_resource_name_create_user` (`resource_name`,`create_user`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `cloud_app`
--

LOCK TABLES `cloud_app` WRITE;
/*!40000 ALTER TABLE `cloud_app` DISABLE KEYS */;
/*!40000 ALTER TABLE `cloud_app` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `cloud_app_configure`
--

DROP TABLE IF EXISTS `cloud_app_configure`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `cloud_app_configure` (
  `configure_id` int(11) NOT NULL AUTO_INCREMENT,
  `configure_name` varchar(32) DEFAULT NULL COMMENT '模板名称',
  `create_time` varchar(32) DEFAULT NULL COMMENT '创建时间',
  `create_user` varchar(32) DEFAULT NULL COMMENT '创建用户',
  `last_modify_time` varchar(32) DEFAULT NULL COMMENT '最近修改时间',
  `last_modify_user` varchar(32) DEFAULT NULL COMMENT '最近修改用户',
  `cluster_name` varchar(32) DEFAULT NULL COMMENT '集群名称',
  `description` varchar(50) DEFAULT NULL COMMENT '描述信息',
  `entname` varchar(32) DEFAULT NULL COMMENT '环境名次',
  PRIMARY KEY (`configure_id`),
  UNIQUE KEY `uidx_template_name` (`configure_name`,`cluster_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `cloud_app_configure`
--

LOCK TABLES `cloud_app_configure` WRITE;
/*!40000 ALTER TABLE `cloud_app_configure` DISABLE KEYS */;
/*!40000 ALTER TABLE `cloud_app_configure` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `cloud_app_service`
--

DROP TABLE IF EXISTS `cloud_app_service`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `cloud_app_service` (
  `service_id` int(11) NOT NULL AUTO_INCREMENT,
  `service_name` varchar(36) DEFAULT NULL COMMENT 'service名称',
  `service_type` varchar(32) DEFAULT NULL COMMENT 'service类型,有状态和无状态',
  `status` varchar(20) DEFAULT NULL COMMENT '运行状态',
  `resource_name` varchar(32) DEFAULT NULL COMMENT '资源空间',
  `cluster_name` varchar(40) DEFAULT NULL COMMENT '集群名称',
  `create_time` varchar(32) DEFAULT NULL COMMENT '创建时间',
  `create_user` varchar(32) DEFAULT NULL COMMENT '创建用户',
  `last_modify_time` varchar(32) DEFAULT NULL COMMENT '最近修改时间',
  `last_modify_user` varchar(32) DEFAULT NULL COMMENT '最近修改用户',
  `last_update_time` varchar(32) DEFAULT NULL COMMENT '最近更新时间',
  `app_labels` text,
  `service_labels` text COMMENT '服务标签,用map标识',
  `cpu` int(11) DEFAULT NULL COMMENT 'cpu核数',
  `memory` int(11) DEFAULT NULL COMMENT '内存数',
  `json_data` text COMMENT '其他非固定数据存储',
  `yaml` text COMMENT '编排文件',
  `network_mode` varchar(10) DEFAULT NULL COMMENT '网络模式 flannel host',
  `container_port` varchar(100) DEFAULT NULL COMMENT '容器端口,多个逗号分隔',
  `lb_name` varchar(400) DEFAULT NULL COMMENT '负载均衡名称',
  `lb_data` text COMMENT '负载均衡数据',
  `image_tag` varchar(300) DEFAULT NULL,
  `deploy_type` varchar(50) DEFAULT NULL COMMENT '部署模式, deployment daemonset statefulset',
  `replicas` int(11) DEFAULT NULL COMMENT '容器副本数量',
  `replicas_max` int(11) DEFAULT NULL COMMENT '容器最多数量',
  `replicas_min` int(11) DEFAULT NULL COMMENT '容器最小数量',
  `env_file` varchar(30) DEFAULT NULL COMMENT '参考环境配置的信息',
  `config` text COMMENT '手动配置文件的内容',
  `health_data` text COMMENT '健康检查数据',
  `service_lables_data` text COMMENT '服务器标签数据',
  `storage_data` text COMMENT '存储配置数据',
  `configure_data` text COMMENT '挂载配置文件数据',
  `envs` text COMMENT 'æ·»åŠ çš„çŽ¯å¢ƒå˜é‡ä¿¡æ¯',
  `app_name` varchar(40) DEFAULT NULL COMMENT '应用名称',
  `max_surge` int(11) DEFAULT '1' COMMENT '滚动升级时候,会优先启动的pod数量',
  `max_unavailable` int(11) DEFAULT '1' COMMENT '滚动升级时候,最大的unavailable数量',
  `min_ready` int(11) DEFAULT NULL COMMENT '指定没有任何容器crash的Pod并被认为是可用状态的最小秒数',
  `image_registry` varchar(100) DEFAULT NULL COMMENT '镜像仓库地址',
  `entname` varchar(100) DEFAULT NULL COMMENT '环境名称',
  `service_version` varchar(3) DEFAULT '1' COMMENT '做蓝绿,灰度部署标签,有1和2,如果1存在那么就部署一个2,如果2存在就部署一个1,当确认发布完成 ,删除一个未使用的部署',
  PRIMARY KEY (`service_id`),
  KEY `idx_cloud_app_service_resource_name_create_user` (`resource_name`,`create_user`),
  KEY `uidx_app_name_service_name_cluster_name` (`app_name`,`service_name`,`cluster_name`,`service_version`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `cloud_app_service`
--

LOCK TABLES `cloud_app_service` WRITE;
/*!40000 ALTER TABLE `cloud_app_service` DISABLE KEYS */;
/*!40000 ALTER TABLE `cloud_app_service` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `cloud_app_template`
--

DROP TABLE IF EXISTS `cloud_app_template`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `cloud_app_template` (
  `template_id` int(11) NOT NULL AUTO_INCREMENT,
  `template_name` varchar(32) DEFAULT NULL COMMENT '模板名称',
  `create_time` varchar(32) DEFAULT NULL COMMENT '创建时间',
  `create_user` varchar(32) DEFAULT NULL COMMENT '创建用户',
  `last_modify_time` varchar(32) DEFAULT NULL COMMENT '最近修改时间',
  `last_modify_user` varchar(32) DEFAULT NULL COMMENT '最近修改用户',
  `resource_name` varchar(32) DEFAULT NULL COMMENT '资源名称',
  `description` varchar(50) DEFAULT NULL COMMENT '描述信息',
  `yaml` text COMMENT 'yaml编排文件',
  PRIMARY KEY (`template_id`),
  UNIQUE KEY `uidx_template_name` (`template_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `cloud_app_template`
--

LOCK TABLES `cloud_app_template` WRITE;
/*!40000 ALTER TABLE `cloud_app_template` DISABLE KEYS */;
/*!40000 ALTER TABLE `cloud_app_template` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `cloud_authority_user`
--

DROP TABLE IF EXISTS `cloud_authority_user`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `cloud_authority_user` (
  `user_id` int(11) NOT NULL AUTO_INCREMENT COMMENT '用户id',
  `user_name` varchar(64) DEFAULT NULL COMMENT '用户名称',
  `user_pic` varchar(255) DEFAULT NULL COMMENT '用户头像',
  `third_id` varchar(40) DEFAULT NULL COMMENT '第三方id',
  `third_true_name` varchar(64) DEFAULT NULL COMMENT '真实姓名',
  `user_email` varchar(255) DEFAULT NULL COMMENT '用户邮箱',
  `user_mobile` varchar(20) DEFAULT NULL COMMENT '用户电话',
  `is_valid` int(1) DEFAULT NULL COMMENT '是否启用（0无效，1有效）',
  `is_del` int(1) DEFAULT NULL COMMENT '是否删除（0未删除，1删除）',
  `pwd` varchar(32) DEFAULT NULL,
  `real_name` varchar(32) DEFAULT NULL COMMENT '真实姓名',
  `description` varchar(32) DEFAULT NULL COMMENT '描叙信息',
  `dept` varchar(132) DEFAULT NULL COMMENT '所属部门',
  `last_modify_time` varchar(32) DEFAULT NULL COMMENT '最近修改时间',
  `last_modify_user` varchar(32) DEFAULT NULL COMMENT '最近修改用户',
  `create_time` varchar(32) DEFAULT NULL COMMENT '创建时间',
  `create_user` varchar(32) DEFAULT NULL COMMENT '创建用户',
  PRIMARY KEY (`user_id`),
  UNIQUE KEY `uidx_cloud_authority_user_username` (`user_name`)
) ENGINE=InnoDB AUTO_INCREMENT=58 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `cloud_authority_user`
--

--
-- Table structure for table `cloud_auto_scale`
--

DROP TABLE IF EXISTS `cloud_auto_scale`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `cloud_auto_scale` (
  `scale_id` int(11) NOT NULL AUTO_INCREMENT,
  `create_time` varchar(32) DEFAULT NULL COMMENT '创建时间',
  `create_user` varchar(32) DEFAULT NULL COMMENT '创建用户',
  `cluster_name` varchar(36) DEFAULT NULL COMMENT '集群名称',
  `last_modify_time` varchar(32) DEFAULT NULL COMMENT '最近修改时间',
  `last_modify_user` varchar(32) DEFAULT NULL COMMENT '最近修改用户',
  `app_name` varchar(100) DEFAULT NULL COMMENT '应用名称',
  `resource_name` varchar(100) DEFAULT NULL COMMENT '资源名称',
  `lt` int(11) DEFAULT NULL COMMENT '阈值小于多少',
  `gt` int(11) DEFAULT NULL COMMENT '阈值大于多少',
  `last_count` int(11) DEFAULT NULL COMMENT '最近几次超过阈值',
  `step` varchar(32) DEFAULT NULL COMMENT '查询监控时间步长',
  `start` varchar(32) DEFAULT NULL COMMENT '开始时间',
  `end` varchar(32) DEFAULT NULL COMMENT '结束时间',
  `query` varchar(500) DEFAULT NULL COMMENT '查询参数',
  `namespace` varchar(50) DEFAULT NULL COMMENT '命名空间',
  `service_version` varchar(3) DEFAULT '1' COMMENT '服务版本号',
  `service_name` varchar(36) DEFAULT NULL COMMENT '服务的名称',
  `increase_step` int(11) DEFAULT NULL COMMENT '扩容步长',
  `reduce_step` int(11) DEFAULT NULL COMMENT '缩容步长',
  `action_interval` int(11) DEFAULT NULL COMMENT '扩容或缩容间隔',
  `msg_group` varchar(100) DEFAULT NULL COMMENT '扩容或缩容进行时,发送通知组',
  `entname` varchar(32) DEFAULT NULL COMMENT '环境名称',
  `description` varchar(100) DEFAULT NULL COMMENT '描述信息',
  `metric_type` varchar(10) DEFAULT NULL COMMENT '指标类型',
  `metric_name` varchar(10) DEFAULT NULL COMMENT '指标名称',
  `es` varchar(100) DEFAULT NULL COMMENT 'es连接地址',
  `data_source` varchar(100) DEFAULT NULL COMMENT '数据源,prometheus,es',
  PRIMARY KEY (`scale_id`)
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `cloud_auto_scale`
--

LOCK TABLES `cloud_auto_scale` WRITE;
/*!40000 ALTER TABLE `cloud_auto_scale` DISABLE KEYS */;
/*!40000 ALTER TABLE `cloud_auto_scale` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `cloud_auto_scale_log`
--

DROP TABLE IF EXISTS `cloud_auto_scale_log`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `cloud_auto_scale_log` (
  `log_id` int(11) NOT NULL AUTO_INCREMENT,
  `create_time` varchar(32) DEFAULT NULL COMMENT '创建时间',
  `query` varchar(500) DEFAULT NULL COMMENT '查询参数',
  `app_name` varchar(100) DEFAULT NULL COMMENT '应用名称',
  `cluster_name` varchar(36) DEFAULT NULL COMMENT '集群名称',
  `metric_type` varchar(10) DEFAULT NULL COMMENT '指标类型',
  `metric_name` varchar(10) DEFAULT NULL COMMENT '指标名称',
  `es` varchar(100) DEFAULT NULL COMMENT 'es连接地址',
  `last_count` int(11) DEFAULT NULL COMMENT '最近几次超过阈值',
  `increase_step` int(11) DEFAULT NULL COMMENT '扩容步长',
  `reduce_step` int(11) DEFAULT NULL COMMENT '缩容步长',
  `status` varchar(32) DEFAULT NULL COMMENT '扩容状态,成功失败',
  `monitor_value` double DEFAULT NULL,
  `replicas_max` int(11) DEFAULT NULL COMMENT '最大值',
  `replicas_min` int(11) DEFAULT NULL COMMENT '最小值',
  `replicas` int(11) DEFAULT NULL COMMENT '扩展到',
  `action_interval` int(11) DEFAULT NULL COMMENT '扩容或缩容间隔',
  `gt` int(11) DEFAULT NULL,
  `step` varchar(32) DEFAULT NULL,
  `service_name` varchar(132) DEFAULT NULL,
  `entname` varchar(100) DEFAULT NULL,
  PRIMARY KEY (`log_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `cloud_auto_scale_log`
--

LOCK TABLES `cloud_auto_scale_log` WRITE;
/*!40000 ALTER TABLE `cloud_auto_scale_log` DISABLE KEYS */;
/*!40000 ALTER TABLE `cloud_auto_scale_log` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `cloud_build_job`
--

DROP TABLE IF EXISTS `cloud_build_job`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `cloud_build_job` (
  `job_id` int(11) NOT NULL AUTO_INCREMENT,
  `job_name` varchar(100) DEFAULT NULL COMMENT '任务计划名称',
  `job_code` int(11) DEFAULT NULL COMMENT '参考code代码仓库',
  `docker_file` varchar(100) DEFAULT NULL,
  `registry_server` varchar(100) DEFAULT NULL COMMENT '注册服务器',
  `item_name` varchar(100) DEFAULT NULL COMMENT '项目名称',
  `image_tag` varchar(100) DEFAULT NULL COMMENT '镜像tag',
  `cluster_name` varchar(100) DEFAULT NULL COMMENT '集群名称',
  `create_time` varchar(32) DEFAULT NULL COMMENT '创建时间',
  `create_user` varchar(32) DEFAULT NULL COMMENT '创建用户',
  `last_modify_time` varchar(32) DEFAULT NULL COMMENT '最近修改时间',
  `last_modify_user` varchar(32) DEFAULT NULL COMMENT '最近修改用户',
  `build_id` varchar(100) DEFAULT NULL COMMENT 'build时k8s.job名称',
  `last_build_time` varchar(32) DEFAULT NULL COMMENT '最近构建时间',
  `build_status` varchar(32) DEFAULT NULL COMMENT '构建状态',
  `description` varchar(100) DEFAULT NULL,
  `content` text COMMENT 'dockerfile数据',
  `script` text COMMENT '构建脚本',
  `time_out` int(11) DEFAULT NULL COMMENT '构建超时时间,最大3600秒,最小10秒',
  `last_tag` varchar(50) DEFAULT NULL COMMENT '最新tag',
  `base_image` varchar(200) DEFAULT NULL COMMENT '基础镜像,里面应该运行docker服务',
  PRIMARY KEY (`job_id`)
) ENGINE=InnoDB AUTO_INCREMENT=18 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `cloud_build_job`
--


--
-- Table structure for table `cloud_build_job_history`
--

DROP TABLE IF EXISTS `cloud_build_job_history`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `cloud_build_job_history` (
  `history_id` int(11) NOT NULL AUTO_INCREMENT,
  `job_name` varchar(100) DEFAULT NULL COMMENT '任务计划名称',
  `job_id` int(11) DEFAULT NULL COMMENT '参考job表ID',
  `docker_file` text,
  `registry_server` varchar(100) DEFAULT NULL COMMENT '注册服务器',
  `item_name` varchar(100) DEFAULT NULL COMMENT '项目名称',
  `image_tag` varchar(100) DEFAULT NULL COMMENT '镜像tag',
  `cluster_name` varchar(100) DEFAULT NULL COMMENT '集群名称',
  `create_time` varchar(32) DEFAULT NULL COMMENT '创建时间',
  `create_user` varchar(32) DEFAULT NULL COMMENT '创建用户',
  `build_status` varchar(32) DEFAULT NULL COMMENT '构建状态',
  `build_time` int(11) DEFAULT NULL COMMENT '构建时间',
  `build_logs` text COMMENT '构建日志',
    `script` text COMMENT '构建脚本',
  `registry_group` varchar(100) DEFAULT NULL COMMENT '仓库组',
  `base_image` varchar(200) DEFAULT NULL COMMENT '基础镜像,里面应该运行docker服务',
  PRIMARY KEY (`history_id`),
  KEY `idx_cloud_build_job_history_job_name` (`job_name`,`job_id`,`item_name`)
) ENGINE=InnoDB AUTO_INCREMENT=446 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `cloud_build_job_history`
--



--
-- Table structure for table `cloud_ci_dockerfile`
--

DROP TABLE IF EXISTS `cloud_ci_dockerfile`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `cloud_ci_dockerfile` (
  `file_id` int(11) NOT NULL AUTO_INCREMENT,
  `content` text COMMENT 'dockerfile内容',
  `script` text COMMENT '构建脚本',
  `name` varchar(100) DEFAULT NULL COMMENT '文件名称',
  `create_time` varchar(32) DEFAULT NULL COMMENT '创建时间',
  `create_user` varchar(32) DEFAULT NULL COMMENT '创建用户',
  `last_modify_time` varchar(32) DEFAULT NULL COMMENT '最近修改时间',
  `last_modify_user` varchar(32) DEFAULT NULL COMMENT '最近修改用户',
  `description` varchar(50) DEFAULT NULL COMMENT '描述信息',
  `is_del` int(11) DEFAULT NULL COMMENT '是否删除,1删除',
  PRIMARY KEY (`file_id`),
  UNIQUE KEY `udx_cloud_ci_dockerfile_name` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=16 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `cloud_ci_dockerfile`
--


--
-- Table structure for table `cloud_ci_perm`
--

DROP TABLE IF EXISTS `cloud_ci_perm`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `cloud_ci_perm` (
  `perm_id` int(11) NOT NULL AUTO_INCREMENT,
  `username` varchar(32) DEFAULT NULL COMMENT '权限用户',
  `groups_name` varchar(32) DEFAULT NULL COMMENT '组用户',
  `datas` text COMMENT '拥有权限',
  `create_time` varchar(32) DEFAULT NULL COMMENT '创建时间',
  `create_user` varchar(32) DEFAULT NULL COMMENT '创建用户',
  `last_modify_time` varchar(32) DEFAULT NULL COMMENT '最近修改时间',
  `last_modify_user` varchar(20) DEFAULT NULL COMMENT '最近修改用户',
  PRIMARY KEY (`perm_id`),
  UNIQUE KEY `uidx_username_groups_name` (`username`,`groups_name`)
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `cloud_ci_perm`
--
--
-- Table structure for table `cloud_ci_release_history`
--

DROP TABLE IF EXISTS `cloud_ci_release_history`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `cloud_ci_release_history` (
  `history_id` int(11) NOT NULL AUTO_INCREMENT,
  `create_time` varchar(32) DEFAULT NULL COMMENT '创建时间',
  `create_user` varchar(32) DEFAULT NULL COMMENT '创建用户',
  `domain` varchar(100) DEFAULT NULL COMMENT '域名',
  `service_name` varchar(100) DEFAULT NULL COMMENT '服务名称',
  `cluster_name` varchar(100) DEFAULT NULL COMMENT '集群名称',
  `app_name` varchar(100) DEFAULT NULL COMMENT '应用名称',
  `entname` varchar(100) DEFAULT NULL COMMENT '环境名称',
  `description` varchar(400) DEFAULT NULL,
  `release_online_type` varchar(32) DEFAULT NULL COMMENT '发布类型',
  `release_bug_description` varchar(232) DEFAULT NULL COMMENT 'Bug修复功能描述',
  `release_demand_description` varchar(232) DEFAULT NULL COMMENT '需求名称',
  `release_item_description` varchar(232) DEFAULT NULL COMMENT '项目名称',
  `release_bug_pm_link` varchar(232) DEFAULT NULL COMMENT '禅道Bug链接',
  `release_job_pm_link` varchar(232) DEFAULT NULL COMMENT '禅道任务链接',
  `release_test_user` varchar(32) DEFAULT NULL COMMENT '测试人员',
  `release_type` varchar(100) DEFAULT NULL COMMENT '发布类型,金丝雀,蓝绿,滚动',
  `image_name` varchar(200) DEFAULT NULL COMMENT '镜像名称',
  `status` varchar(32) DEFAULT NULL COMMENT '发布状态',
  `service_id` int(11) DEFAULT NULL,
  `auto_switch` varchar(5) DEFAULT NULL COMMENT '自动切换',
  `lb_version` varchar(10) DEFAULT NULL,
  `release_version` varchar(3) DEFAULT NULL COMMENT '发布版本',
  `action` varchar(13) DEFAULT NULL COMMENT '发布或回滚,update,rollback',
  `percent` int(11) DEFAULT NULL COMMENT '流量切入百分比',
  `old_images` varchar(200) DEFAULT NULL COMMENT '旧的版本',
  PRIMARY KEY (`history_id`)
) ENGINE=InnoDB AUTO_INCREMENT=18 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `cloud_ci_release_history`
--

--
-- Table structure for table `cloud_ci_release_log`
--

DROP TABLE IF EXISTS `cloud_ci_release_log`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `cloud_ci_release_log` (
  `log_id` int(11) NOT NULL AUTO_INCREMENT,
  `domain` varchar(100) DEFAULT NULL COMMENT '域名',
  `service_name` varchar(100) DEFAULT NULL COMMENT '服务名称',
  `cluster_name` varchar(100) DEFAULT NULL COMMENT '集群名称',
  `app_name` varchar(100) DEFAULT NULL COMMENT '应用名称',
  `entname` varchar(100) DEFAULT NULL COMMENT '环境名称',
  `action` varchar(13) DEFAULT NULL COMMENT '发布或回滚,update,rollback',
  `create_time` varchar(32) DEFAULT NULL COMMENT '创建时间',
  `create_user` varchar(32) DEFAULT NULL COMMENT '创建用户',
  `messages` varchar(200) DEFAULT NULL,
  `ip` varchar(32) DEFAULT NULL COMMENT '操作IP',
  PRIMARY KEY (`log_id`)
) ENGINE=InnoDB AUTO_INCREMENT=108 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `cloud_ci_release_log`
--


--
-- Table structure for table `cloud_ci_service`
--

DROP TABLE IF EXISTS `cloud_ci_service`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `cloud_ci_service` (
  `service_id` int(11) NOT NULL AUTO_INCREMENT,
  `domain` varchar(100) DEFAULT NULL COMMENT '域名',
  `service_name` varchar(100) DEFAULT NULL COMMENT '服务名称',
  `entname` varchar(100) DEFAULT NULL COMMENT '环境名称',
  `group_name` varchar(100) DEFAULT NULL COMMENT '组名称',
  `app_name` varchar(100) DEFAULT NULL COMMENT '应用名称',
  `cluster_name` varchar(100) DEFAULT NULL COMMENT '集群名称',
  `create_time` varchar(32) DEFAULT NULL COMMENT '创建时间',
  `create_user` varchar(32) DEFAULT NULL COMMENT '创建用户',
  `last_modify_time` varchar(32) DEFAULT NULL COMMENT '最近修改时间',
  `last_modify_user` varchar(32) DEFAULT NULL COMMENT '最近修改用户',
  `description` varchar(300) DEFAULT NULL COMMENT '服务描述信息',
  `release_type` varchar(100) DEFAULT NULL COMMENT '发布类型,金丝雀,蓝绿,滚动',
  `image_name` varchar(200) DEFAULT NULL COMMENT '镜像名称',
  `status` varchar(32) DEFAULT NULL COMMENT '发布状态',
  `current_version` varchar(5) DEFAULT NULL,
  `lb_version` varchar(32) DEFAULT NULL COMMENT '负载均衡服务使用的版本',
  `image_info_blue` varchar(32) DEFAULT NULL,
  `image_info_green` varchar(32) DEFAULT NULL,
  `blue_access` varchar(32) DEFAULT NULL,
  `green_access` varchar(32) DEFAULT NULL,
  `blue_pod` varchar(32) DEFAULT NULL,
  `percent` int(11) DEFAULT NULL COMMENT 'æµé‡åˆ‡å…¥ç™¾åˆ†æ¯”',
  `lb_service` varchar(32) DEFAULT NULL,
  `new_version` varchar(32) DEFAULT NULL COMMENT '最新版本名称',
  PRIMARY KEY (`service_id`),
  UNIQUE KEY `uidx_cloud_ci_service_domain` (`domain`)
) ENGINE=InnoDB AUTO_INCREMENT=10 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `cloud_ci_service`
--


--
-- Table structure for table `cloud_cluster`
--

DROP TABLE IF EXISTS `cloud_cluster`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `cloud_cluster` (
  `cluster_id` int(11) NOT NULL AUTO_INCREMENT,
  `cluster_type` varchar(15) DEFAULT NULL COMMENT '集群类型',
  `cluster_name` varchar(32) DEFAULT NULL COMMENT '集群名称,必须英文',
  `cluster_alias` varchar(32) DEFAULT NULL COMMENT '集群显示名称',
  `create_time` varchar(32) DEFAULT NULL COMMENT '创建时间',
  `create_user` varchar(32) DEFAULT NULL COMMENT '创建用户',
  `last_modify_time` varchar(32) DEFAULT NULL COMMENT '最近修改时间',
  `docker_version` varchar(15) DEFAULT NULL COMMENT 'docker版本',
  `docker_install_dir` varchar(100) DEFAULT NULL COMMENT 'docker安装路径',
  `network_cart` varchar(10) DEFAULT NULL COMMENT '内网网卡名称',
  `last_modify_user` varchar(32) DEFAULT NULL COMMENT '最近修改时间',
  `ca_data` text COMMENT 'ca证书文件内容ca.pem文件',
  `cert_data` text COMMENT 'node证书文件内容ca.pem文件',
  `key_data` text COMMENT 'node证书文件内容worker-key.pem文件',
  `api_address` varchar(50) DEFAULT NULL COMMENT '集群APi地址',
  PRIMARY KEY (`cluster_id`),
  UNIQUE KEY `uidx_cloud_cluster_cluser_name` (`cluster_name`)
) ENGINE=InnoDB AUTO_INCREMENT=57 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `cloud_cluster`
--

--
-- Table structure for table `cloud_cluster_hosts`
--

DROP TABLE IF EXISTS `cloud_cluster_hosts`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `cloud_cluster_hosts` (
  `host_id` int(11) NOT NULL AUTO_INCREMENT,
  `host_ip` varchar(32) DEFAULT NULL COMMENT '主机IP',
  `host_label` varchar(100) DEFAULT NULL COMMENT '主机标签',
  `status` varchar(10) DEFAULT NULL COMMENT '状态',
  `create_method` varchar(32) DEFAULT NULL COMMENT '创建方法',
  `pod_num` int(11) DEFAULT NULL COMMENT 'pod数量',
  `cpu_num` int(11) DEFAULT NULL COMMENT 'cpu数量',
  `mem_size` varchar(32) DEFAULT NULL COMMENT '内存大小',
  `host_type` varchar(32) DEFAULT NULL COMMENT '主机类型',
  `create_time` varchar(32) DEFAULT NULL COMMENT '创建时间',
  `create_user` varchar(32) DEFAULT NULL COMMENT '创建用户',
  `last_modify_time` varchar(32) DEFAULT NULL COMMENT '最近修改时间',
  `last_modify_user` varchar(32) DEFAULT NULL COMMENT '最近修改时间',
  `is_valid` int(11) DEFAULT NULL COMMENT '是否有效',
  `container_num` int(11) DEFAULT NULL COMMENT '容器数量',
  `mem_free` varchar(32) DEFAULT NULL COMMENT '内存剩余量',
  `cpu_free` varchar(32) DEFAULT NULL COMMENT 'cpu剩余量',
  `mem_percent` varchar(10) DEFAULT NULL COMMENT '内使用百分比',
  `cpu_percent` varchar(10) DEFAULT NULL COMMENT 'cpu使用百分比',
  `cluster_name` varchar(36) DEFAULT NULL COMMENT '所属集群',
  `api_port` varchar(8) DEFAULT NULL COMMENT 'k8sAPiç«¯å£,åªéœ€è¦masteræœ‰å°±è¡Œäº†',
  `image_num` int(11) DEFAULT NULL,
  PRIMARY KEY (`host_id`)
) ENGINE=InnoDB AUTO_INCREMENT=31 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `cloud_cluster_hosts`
--

--
-- Table structure for table `cloud_code_repostitory`
--

DROP TABLE IF EXISTS `cloud_code_repostitory`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `cloud_code_repostitory` (
  `repostitory_id` int(11) NOT NULL AUTO_INCREMENT,
  `create_time` varchar(32) DEFAULT NULL COMMENT '创建时间',
  `create_user` varchar(32) DEFAULT NULL COMMENT '创建用户',
  `last_modify_time` varchar(32) DEFAULT NULL COMMENT '最近修改时间',
  `last_modify_user` varchar(32) DEFAULT NULL COMMENT '最近修改用户',
  `code_source` varchar(30) DEFAULT NULL COMMENT '代码来源, gitlab,github,svn ',
  `username` varchar(100) DEFAULT NULL COMMENT '用户名',
  `password` varchar(100) DEFAULT NULL COMMENT '密码,base64存储',
  `gitlab_token` varchar(100) DEFAULT NULL COMMENT 'gitlab token',
  `type` varchar(32) DEFAULT NULL COMMENT '代码类型, public private 共有和私有',
  `code_url` varchar(200) DEFAULT NULL COMMENT '代码路径',
  `description` varchar(100) DEFAULT NULL COMMENT '描述信息',
  PRIMARY KEY (`repostitory_id`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `cloud_code_repostitory`
--


--
-- Table structure for table `cloud_config_data`
--

DROP TABLE IF EXISTS `cloud_config_data`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `cloud_config_data` (
  `data_id` int(11) NOT NULL AUTO_INCREMENT,
  `create_time` varchar(32) DEFAULT NULL COMMENT '创建时间',
  `create_user` varchar(32) DEFAULT NULL COMMENT '创建用户',
  `last_modify_time` varchar(32) DEFAULT NULL COMMENT '最近修改时间',
  `last_modify_user` varchar(32) DEFAULT NULL COMMENT '最近修改用户',
  `configure_id` int(11) DEFAULT NULL COMMENT '参考config的id',
  `configure_name` varchar(32) DEFAULT NULL COMMENT '配置名称,参考配置名称',
  `data` text COMMENT '配置文件数据',
  `data_name` varchar(100) DEFAULT NULL,
  PRIMARY KEY (`data_id`),
  UNIQUE KEY `uidx_config_name_data_name` (`configure_id`,`data_name`)
) ENGINE=InnoDB AUTO_INCREMENT=28 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `cloud_config_data`
--



--
-- Table structure for table `cloud_configure_mount`
--

DROP TABLE IF EXISTS `cloud_configure_mount`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `cloud_configure_mount` (
  `mount_id` int(11) NOT NULL AUTO_INCREMENT,
  `configure_name` varchar(100) DEFAULT NULL COMMENT '配置文件名称',
  `namespace` varchar(200) DEFAULT NULL COMMENT '命名空间',
  `cluster_name` varchar(100) DEFAULT NULL COMMENT '集群名称',
  `last_update_time` varchar(32) DEFAULT NULL COMMENT '最近更新时间',
  `create_time` varchar(32) DEFAULT NULL COMMENT '首次挂载时间',
  `data_name` varchar(100) DEFAULT NULL COMMENT '数据key名称',
  `mount_path` varchar(300) DEFAULT NULL COMMENT '挂载路径',
  `service_name` varchar(300) DEFAULT NULL COMMENT '服务名称',
  PRIMARY KEY (`mount_id`),
  UNIQUE KEY `uidx_configure_mount_name_cluster_data_key` (`configure_name`,`data_name`,`cluster_name`,`namespace`)
) ENGINE=InnoDB AUTO_INCREMENT=726 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `cloud_configure_mount`
--


--
-- Table structure for table `cloud_container`
--

DROP TABLE IF EXISTS `cloud_container`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `cloud_container` (
  `container_id` int(11) NOT NULL AUTO_INCREMENT,
  `container_name` varchar(200) DEFAULT NULL COMMENT '容器名称',
  `create_time` varchar(32) DEFAULT NULL COMMENT '创建时间',
  `service_name` varchar(100) DEFAULT NULL COMMENT '服务名称',
  `cluster_name` varchar(100) DEFAULT NULL COMMENT '集群名称',
  `server_address` varchar(32) DEFAULT NULL COMMENT '宿主机地址',
  `container_ip` varchar(32) DEFAULT NULL COMMENT '容器ip',
  `image` varchar(300) DEFAULT NULL COMMENT '镜像名称',
  `app_name` varchar(132) DEFAULT NULL COMMENT '应用名称',
  `status` varchar(32) DEFAULT NULL COMMENT '运行状态',
  `resource_name` varchar(100) DEFAULT NULL COMMENT '资源空间',
  `cpu` int(11) DEFAULT NULL,
  `memory` int(11) DEFAULT NULL,
  `env` text,
  `process` text,
  `storage_data` text,
  `waiting_messages` text COMMENT '容器等待时的信息',
  `waiting_reason` varchar(100) DEFAULT NULL COMMENT '容器等待原因',
  `terminated_reason` varchar(100) DEFAULT NULL COMMENT '容器停止原因',
  `terminated_messages` varchar(100) DEFAULT NULL COMMENT '容器停止信息',
  `create_user` varchar(100) DEFAULT NULL,
  `entname` varchar(100) DEFAULT NULL COMMENT '环境名称',
  `events` varchar(32) DEFAULT NULL,
  PRIMARY KEY (`container_id`),
  UNIQUE KEY `uidx_container_name` (`container_name`),
  KEY `idx_name` (`container_name`),
  KEY `idx_cloud_container_resource_name_create_user` (`resource_name`,`create_user`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `cloud_container`
--

LOCK TABLES `cloud_container` WRITE;
/*!40000 ALTER TABLE `cloud_container` DISABLE KEYS */;
/*!40000 ALTER TABLE `cloud_container` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `cloud_ent`
--

DROP TABLE IF EXISTS `cloud_ent`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `cloud_ent` (
  `ent_id` int(11) NOT NULL AUTO_INCREMENT,
  `create_user` varchar(32) DEFAULT NULL COMMENT '创建人',
  `last_modify_user` varchar(32) DEFAULT NULL COMMENT '最近修改人',
  `last_modify_time` varchar(32) DEFAULT NULL COMMENT '最近修改时间',
  `description` varchar(50) DEFAULT NULL COMMENT '备注信息',
  `create_time` varchar(32) DEFAULT NULL COMMENT '创建时间',
  `entname` varchar(32) DEFAULT NULL COMMENT '环境名称',
  `clusters` text COMMENT '拥有集群',
  PRIMARY KEY (`ent_id`),
  UNIQUE KEY `uidx_ci_entname` (`entname`)
) ENGINE=InnoDB AUTO_INCREMENT=10 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `cloud_ent`
--


--
-- Table structure for table `cloud_image`
--

DROP TABLE IF EXISTS `cloud_image`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `cloud_image` (
  `image_id` int(11) NOT NULL AUTO_INCREMENT,
  `create_time` varchar(32) DEFAULT NULL COMMENT '创建时间',
  `create_user` varchar(32) DEFAULT NULL COMMENT '创建用户',
  `size` int(11) DEFAULT NULL COMMENT '镜像大小',
  `name` varchar(100) DEFAULT NULL COMMENT '镜像名称',
  `repositories` varchar(100) DEFAULT NULL COMMENT '所属仓库',
  `repositories_group` varchar(100) DEFAULT NULL COMMENT '镜像仓库组',
  `image_type` varchar(20) DEFAULT NULL COMMENT '镜像类型,分为共有和私有',
  `tag_number` int(11) DEFAULT NULL COMMENT '版本数量',
  `access` varchar(100) DEFAULT NULL COMMENT '访问方式',
  `layers_number` int(11) DEFAULT NULL COMMENT '镜像层数',
  `tags` text,
  `download` int(11) DEFAULT NULL COMMENT '下载次数',
  PRIMARY KEY (`image_id`),
  UNIQUE KEY `cloud_image_name_repositories_group` (`repositories_group`,`name`)
) ENGINE=InnoDB AUTO_INCREMENT=90 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `cloud_image`
--


--
-- Table structure for table `cloud_image_base`
--

DROP TABLE IF EXISTS `cloud_image_base`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `cloud_image_base` (
  `base_id` int(11) NOT NULL AUTO_INCREMENT,
  `image_name` varchar(300) DEFAULT NULL COMMENT '镜像名称',
  `registry_server` varchar(100) DEFAULT NULL COMMENT '镜像仓库地址',
  `create_time` varchar(32) DEFAULT NULL COMMENT '创建时间',
  `create_user` varchar(32) DEFAULT NULL COMMENT '创建用户',
  `last_modify_time` varchar(32) DEFAULT NULL COMMENT '最近修改时间',
  `last_modify_user` varchar(32) DEFAULT NULL COMMENT '最近修改用户',
  `icon` varchar(100) DEFAULT NULL COMMENT '镜像图标',
  `description` varchar(200) DEFAULT NULL COMMENT '镜像描述信息',
  `image_type` varchar(20) DEFAULT NULL COMMENT '镜像类型',
  PRIMARY KEY (`base_id`)
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `cloud_image_base`
--



--
-- Table structure for table `cloud_image_log`
--

DROP TABLE IF EXISTS `cloud_image_log`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `cloud_image_log` (
  `log_id` int(11) NOT NULL AUTO_INCREMENT,
  `create_time` varchar(32) DEFAULT NULL COMMENT '操作时间',
  `create_user` varchar(32) DEFAULT NULL COMMENT '操作人',
  `name` varchar(100) DEFAULT NULL COMMENT '镜像名称',
  `repositories` varchar(100) DEFAULT NULL COMMENT '所属仓库',
  `cluster_name` varchar(100) DEFAULT NULL COMMENT '所属集群',
  `repositories_group` varchar(100) DEFAULT NULL COMMENT '镜像仓库组',
  `oper_type` varchar(20) DEFAULT NULL COMMENT '镜像获取类型,pull,create,push',
  `label` varchar(120) DEFAULT NULL COMMENT '标签名称',
  `ip` varchar(32) DEFAULT NULL COMMENT '操作Ip',
  PRIMARY KEY (`log_id`),
  KEY `cloud_image_log_name_repositories` (`name`,`repositories`,`repositories_group`)
) ENGINE=InnoDB AUTO_INCREMENT=19 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `cloud_image_log`
--

--
-- Table structure for table `cloud_image_sync`
--

DROP TABLE IF EXISTS `cloud_image_sync`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `cloud_image_sync` (
  `sync_id` int(11) NOT NULL AUTO_INCREMENT,
  `cluster_name` varchar(100) NOT NULL,
  `target_cluster` varchar(100) NOT NULL,
  `target_registry` varchar(100) NOT NULL,
  `registry_group` varchar(100) NOT NULL,
  `version` varchar(100) DEFAULT NULL COMMENT '版本号',
  `image_name` varchar(200) DEFAULT NULL COMMENT '镜像名称',
  `target_entname` varchar(100) NOT NULL,
  `entname` varchar(100) NOT NULL,
  `create_time` varchar(32) DEFAULT NULL COMMENT '创建时间',
  `create_user` varchar(32) DEFAULT NULL COMMENT '创建用户',
  `last_modify_time` varchar(32) DEFAULT NULL COMMENT '最近修改时间',
  `last_modify_user` varchar(32) DEFAULT NULL COMMENT '最近修改用户',
  `description` varchar(100) DEFAULT NULL COMMENT '备注信息',
  `approved_by` varchar(100) DEFAULT NULL COMMENT '审批人',
  `approved_time` varchar(32) DEFAULT NULL COMMENT '审批时间',
  `status` varchar(10) DEFAULT NULL COMMENT 'åŒæ­¥çŠ¶æ€',
  `registry` varchar(100) NOT NULL,
  PRIMARY KEY (`sync_id`)
) ENGINE=InnoDB AUTO_INCREMENT=12 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `cloud_image_sync`
--



--
-- Table structure for table `cloud_image_sync_log`
--

DROP TABLE IF EXISTS `cloud_image_sync_log`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `cloud_image_sync_log` (
  `log_id` int(11) NOT NULL AUTO_INCREMENT,
  `create_time` varchar(32) DEFAULT NULL COMMENT '创建时间',
  `messages` text COMMENT '执行内容',
  `create_user` varchar(32) DEFAULT NULL COMMENT '创建用户',
  `runtime` int(11) DEFAULT NULL COMMENT '程序运行时间',
  `registry_group` varchar(32) DEFAULT NULL COMMENT '镜像仓库组',
  `registry_server_1` varchar(100) DEFAULT NULL,
  `registry_server_2` varchar(100) DEFAULT NULL,
  `item_name` varchar(100) DEFAULT NULL COMMENT '项目名称',
  `version` varchar(100) DEFAULT NULL COMMENT '镜像同步日志',
  `status` varchar(10) DEFAULT NULL COMMENT '同步状态',
  PRIMARY KEY (`log_id`)
) ENGINE=InnoDB AUTO_INCREMENT=72 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `cloud_image_sync_log`
--


--
-- Table structure for table `cloud_lb`
--

DROP TABLE IF EXISTS `cloud_lb`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `cloud_lb` (
  `lb_id` int(11) NOT NULL AUTO_INCREMENT,
  `lb_name` varchar(36) DEFAULT NULL COMMENT '负载均衡名称',
  `lb_ip` text COMMENT 'IP地址',
  `lb_type` varchar(36) DEFAULT NULL COMMENT '负载均衡类型,nginx,haproxy',
  `lb_domain_prefix` varchar(36) DEFAULT NULL COMMENT '域名前缀',
  `lb_domain_suffix` varchar(36) DEFAULT NULL COMMENT '域名后缀',
  `description` varchar(300) DEFAULT NULL COMMENT '配额描述信息',
  `create_time` varchar(32) DEFAULT NULL COMMENT '创建时间',
  `create_user` varchar(32) DEFAULT NULL COMMENT '创建用户',
  `cluster_name` varchar(36) DEFAULT NULL COMMENT '集群名称',
  `resource_name` varchar(36) DEFAULT NULL COMMENT '资源空间',
  `last_modify_time` varchar(32) DEFAULT NULL COMMENT '最近修改时间',
  `last_modify_user` varchar(32) DEFAULT NULL COMMENT '最近修改用户',
  `status` varchar(32) DEFAULT NULL,
  `service_number` int(11) DEFAULT NULL,
  `entname` varchar(32) DEFAULT NULL COMMENT '环境名称',
  PRIMARY KEY (`lb_id`),
  UNIQUE KEY `uidx_cloud_lb_lb_name` (`lb_name`)
) ENGINE=InnoDB AUTO_INCREMENT=18 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `cloud_lb`
--


--
-- Table structure for table `cloud_lb_cert`
--

DROP TABLE IF EXISTS `cloud_lb_cert`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `cloud_lb_cert` (
  `cert_id` int(11) NOT NULL AUTO_INCREMENT,
  `cert_key` varchar(100) DEFAULT NULL COMMENT '证书名称',
  `cert_value` text COMMENT '证书内容',
  `create_time` varchar(32) DEFAULT NULL COMMENT '创建时间',
  `create_user` varchar(32) DEFAULT NULL COMMENT '创建用户',
  `last_modify_time` varchar(32) DEFAULT NULL COMMENT '最近修改时间',
  `last_modify_user` varchar(32) DEFAULT NULL COMMENT '最近修改用户',
  `description` varchar(100) DEFAULT NULL COMMENT '描述信息',
  `pem_value` text COMMENT '证书pem文件',
  PRIMARY KEY (`cert_id`)
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `cloud_lb_cert`
--



--
-- Table structure for table `cloud_lb_nginx_conf`
--

DROP TABLE IF EXISTS `cloud_lb_nginx_conf`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `cloud_lb_nginx_conf` (
  `conf_id` int(11) NOT NULL AUTO_INCREMENT,
  `domain` varchar(200) DEFAULT NULL COMMENT '域名',
  `vhost` text COMMENT 'vhost数据',
  `create_user` varchar(100) DEFAULT NULL COMMENT '创建用户',
  `create_time` varchar(32) DEFAULT NULL COMMENT '创建时间',
  `lb_service_id` varchar(11) DEFAULT NULL COMMENT '参考lb服务id',
  `resource_name` varchar(32) DEFAULT NULL COMMENT '资源空间',
  `app_name` varchar(100) DEFAULT NULL COMMENT '应用名称',
  `cluster_name` varchar(32) DEFAULT NULL COMMENT '集群名称',
  `service_name` varchar(32) DEFAULT NULL COMMENT '服务名称',
  `last_modify_time` varchar(32) DEFAULT NULL COMMENT '最近修改时间',
  `last_modify_user` varchar(32) DEFAULT NULL COMMENT ' 最近修改人',
  `service_id` int(11) DEFAULT NULL,
  `cert_file` varchar(100) DEFAULT NULL COMMENT '证书文件路径',
  PRIMARY KEY (`conf_id`),
  UNIQUE KEY `uidx_cloud_lb_nginx_conf_domain_service_name_cluster_name` (`domain`,`cluster_name`)
) ENGINE=InnoDB AUTO_INCREMENT=45 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `cloud_lb_nginx_conf`
--

--
-- Table structure for table `cloud_lb_service`
--

DROP TABLE IF EXISTS `cloud_lb_service`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `cloud_lb_service` (
  `service_id` int(11) NOT NULL AUTO_INCREMENT,
  `service_name` varchar(36) DEFAULT NULL COMMENT '要负载的服务的名称',
  `lb_name` varchar(36) DEFAULT NULL COMMENT '负载均衡名称',
  `cert_file` varchar(36) DEFAULT NULL COMMENT '证书文件',
  `lb_type` varchar(36) DEFAULT NULL COMMENT '负载均衡类型,tcp,http,https',
  `description` varchar(300) DEFAULT NULL COMMENT '服务描述信息',
  `listen_port` varchar(300) DEFAULT NULL COMMENT '监听端口',
  `container_port` varchar(300) DEFAULT NULL COMMENT '转到容器的端口',
  `create_time` varchar(32) DEFAULT NULL COMMENT '创建时间',
  `create_user` varchar(32) DEFAULT NULL COMMENT '创建用户',
  `cluster_name` varchar(36) DEFAULT NULL COMMENT '集群名称',
  `last_modify_time` varchar(32) DEFAULT NULL COMMENT '最近修改时间',
  `last_modify_user` varchar(32) DEFAULT NULL COMMENT '最近修改用户',
  `domain` varchar(200) DEFAULT NULL COMMENT '访问地址',
  `app_name` varchar(100) DEFAULT NULL COMMENT 'appåç§°',
  `resource_name` varchar(100) DEFAULT NULL COMMENT 'èµ„æºç©ºé—´',
  `lb_service_id` varchar(10) DEFAULT NULL COMMENT '参考应用服务ID,生成对应的数据',
  `lb_id` int(11) DEFAULT NULL COMMENT '参考lb服务ID',
  `default_domain` varchar(5) DEFAULT NULL COMMENT '是否配置默认域名,默认为空',
  `lb_method` varchar(10) DEFAULT NULL COMMENT '负载方式,分为node和pod俩种模式',
  `protocol` varchar(10) DEFAULT NULL COMMENT '协议',
  `service_version` varchar(3) DEFAULT '1' COMMENT '服务版本号,做灰度或蓝绿发布使用',
  `entname` varchar(32) DEFAULT NULL COMMENT '环境名称',
  `percent` int(11) DEFAULT NULL COMMENT 'æµé‡åˆ‡å…¥ç™¾åˆ†æ¯”',
  `flow_service_name` varchar(100) DEFAULT NULL,
  PRIMARY KEY (`service_id`),
  UNIQUE KEY `uidx_cloud_lb_service_domain` (`domain`)
) ENGINE=InnoDB AUTO_INCREMENT=60 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `cloud_lb_service`
--

LOCK TABLES `cloud_lb_service` WRITE;
/*!40000 ALTER TABLE `cloud_lb_service` DISABLE KEYS */;
INSERT INTO `cloud_lb_service` VALUES (55,'blue-service','nginx-lb','','Nginx','','','5000','2018-02-14 15:43:51','admin','glusterfs-cluster','2018-02-16 21:52:27','admin','www.gg.com','aaaaaaaaaa','admin-quota','215',11,'0','service','HTTP','1','生产环境',0,'blue-service--2'),(57,'hd1-service-app','nginx-lb',NULL,'Nginx',NULL,NULL,'5000','2018-02-17 19:10:58','admin','glusterfs-cluster','2018-02-17 19:10:58','admin','www.xx.com','hdapp','admin-quota','229',11,'0','service','HTTP','1','生产环境',29,'hd1-service-app--1'),(59,'web-service','nginx-lb',NULL,'Nginx',NULL,NULL,'5000','2018-02-17 19:21:28','admin','glusterfs-cluster','2018-02-17 19:21:28','admin','www.web.com','web','admin-quota','233',11,'0','service','HTTP','2','生产环境',0,'web-service--2');
/*!40000 ALTER TABLE `cloud_lb_service` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `cloud_login_record`
--

DROP TABLE IF EXISTS `cloud_login_record`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `cloud_login_record` (
  `record_id` int(11) NOT NULL AUTO_INCREMENT,
  `login_time` varchar(32) DEFAULT NULL COMMENT '登录时间',
  `login_ip` varchar(32) DEFAULT NULL COMMENT '登录IP',
  `login_user` varchar(32) DEFAULT NULL COMMENT '登录用户名',
  `login_Status` int(11) DEFAULT NULL COMMENT '登录状态',
  PRIMARY KEY (`record_id`)
) ENGINE=InnoDB AUTO_INCREMENT=1039 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `cloud_login_record`
--

--
-- Table structure for table `cloud_oper_log`
--

DROP TABLE IF EXISTS `cloud_oper_log`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `cloud_oper_log` (
  `log_id` int(11) NOT NULL AUTO_INCREMENT,
  `time` varchar(32) DEFAULT NULL COMMENT '操作时间',
  `user` varchar(32) DEFAULT NULL COMMENT '操作用户',
  `messages` text COMMENT '操作信息',
  `cluster` varchar(200) DEFAULT NULL COMMENT '在哪个集群操作的',
  `ip` varchar(32) DEFAULT NULL COMMENT '操作IP地址',
  PRIMARY KEY (`log_id`)
) ENGINE=InnoDB AUTO_INCREMENT=2586 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `cloud_oper_log`
--

LOCK TABLES `cloud_oper_log` WRITE;
/*!40000 ALTER TABLE `cloud_oper_log` DISABLE KEYS */;
/*!40000 ALTER TABLE `cloud_oper_log` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `cloud_perm`
--

DROP TABLE IF EXISTS `cloud_perm`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `cloud_perm` (
  `perm_id` int(11) NOT NULL AUTO_INCREMENT,
  `create_time` varchar(32) DEFAULT NULL COMMENT '创建时间',
  `create_user` varchar(32) DEFAULT NULL COMMENT '创建用户',
  `last_modify_time` varchar(32) DEFAULT NULL COMMENT '最近修改时间',
  `last_modify_user` varchar(32) DEFAULT NULL COMMENT '最近修改用户',
  `description` varchar(300) DEFAULT NULL COMMENT '服务描述信息',
  `user` text COMMENT '拥有用户',
  `groups` text COMMENT '拥有团队',
  `roles` text COMMENT '拥有权限角色',
  PRIMARY KEY (`perm_id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `cloud_perm`
--

LOCK TABLES `cloud_perm` WRITE;
/*!40000 ALTER TABLE `cloud_perm` DISABLE KEYS */;
INSERT INTO `cloud_perm` VALUES (1,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL);
/*!40000 ALTER TABLE `cloud_perm` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `cloud_perm_role`
--

DROP TABLE IF EXISTS `cloud_perm_role`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `cloud_perm_role` (
  `role_id` int(11) NOT NULL AUTO_INCREMENT,
  `role_name` varchar(50) DEFAULT NULL COMMENT '角色名称',
  `permissions` text COMMENT '拥有权限',
  `create_time` varchar(32) DEFAULT NULL COMMENT '创建时间',
  `create_user` varchar(32) DEFAULT NULL COMMENT '创建用户',
  `last_modify_time` varchar(32) DEFAULT NULL COMMENT '最近修改时间',
  `last_modify_user` varchar(32) DEFAULT NULL COMMENT '最近修改用户',
  `description` varchar(300) DEFAULT NULL COMMENT '服务描述信息',
  PRIMARY KEY (`role_id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `cloud_perm_role`
--

LOCK TABLES `cloud_perm_role` WRITE;
/*!40000 ALTER TABLE `cloud_perm_role` DISABLE KEYS */;
/*!40000 ALTER TABLE `cloud_perm_role` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `cloud_pipeline`
--

DROP TABLE IF EXISTS `cloud_pipeline`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `cloud_pipeline` (
  `pipeline_id` int(11) NOT NULL AUTO_INCREMENT,
  `pipeline_name` varchar(100) DEFAULT NULL COMMENT '流水线名称',
  `app_name` varchar(100) DEFAULT NULL COMMENT '应用名称',
  `cluster_name` varchar(32) DEFAULT NULL COMMENT '集群名称',
  `service_name` varchar(32) DEFAULT NULL COMMENT '服务名称',
  `last_modify_time` varchar(32) DEFAULT NULL COMMENT '最近修改时间',
  `last_modify_user` varchar(32) DEFAULT NULL COMMENT ' 最近修改人',
  `job_name` varchar(100) DEFAULT NULL COMMENT '构建任务内容,关联构建任务',
  `create_user` varchar(100) DEFAULT NULL COMMENT '创建用户',
  `create_time` varchar(32) DEFAULT NULL COMMENT '创建时间',
  `exec_time` varchar(32) DEFAULT NULL COMMENT '执行时间',
  `resource_name` varchar(32) DEFAULT NULL COMMENT '资源空间',
  `description` varchar(132) DEFAULT NULL COMMENT '描述信息',
  `fail_action` varchar(32) DEFAULT NULL COMMENT '继续,或暂停',
  `job_id` int(11) DEFAULT NULL COMMENT '参考构建任务ID，用来查看及时构建日志使用',
  `status` varchar(32) DEFAULT NULL COMMENT '显示应用或服务的状态',
  PRIMARY KEY (`pipeline_id`)
) ENGINE=InnoDB AUTO_INCREMENT=18 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `cloud_pipeline`
--

LOCK TABLES `cloud_pipeline` WRITE;
/*!40000 ALTER TABLE `cloud_pipeline` DISABLE KEYS */;
INSERT INTO `cloud_pipeline` VALUES (15,'pipeline1','pipeline','glusterfs-cluster','pipeline-1','2018-02-04 19:38:34','zhaozq14','zmc','zhaozq14','2018-02-04 10:48:31','000',NULL,'啦啦啦啦啦啦啦啦啦啦啦啦啦啦啦','pause',7,NULL),(17,'pipeline-1','test-cloud','glusterfs-cluster','test-cloud-service','2018-02-09 16:05:38','admin','hd2','admin','2018-02-09 09:02:51','*','','发布app-scale','pause',11,'');
/*!40000 ALTER TABLE `cloud_pipeline` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `cloud_pipeline_log`
--

DROP TABLE IF EXISTS `cloud_pipeline_log`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `cloud_pipeline_log` (
  `log_id` int(11) NOT NULL AUTO_INCREMENT,
  `pipeline_name` varchar(100) DEFAULT NULL COMMENT '流水线名称',
  `app_name` varchar(100) DEFAULT NULL COMMENT '应用名称',
  `cluster_name` varchar(32) DEFAULT NULL COMMENT '集群名称',
  `service_name` varchar(32) DEFAULT NULL COMMENT '服务名称',
  `job_name` varchar(100) DEFAULT NULL COMMENT '构建任务内容,关联构建任务',
  `create_user` varchar(100) DEFAULT NULL COMMENT '创建用户',
  `create_time` varchar(32) DEFAULT NULL COMMENT '创建时间',
  `exec_time` varchar(32) DEFAULT NULL COMMENT '执行时间',
  `resource_name` varchar(32) DEFAULT NULL COMMENT '资源空间',
  `status` varchar(32) DEFAULT NULL COMMENT '执行状态,成功或失败',
  `messages` text COMMENT '执行日志',
  `job_id` int(11) DEFAULT NULL,
  `run_time` int(11) DEFAULT NULL COMMENT '运行时间',
  `start_time` varchar(32) DEFAULT NULL COMMENT '启动执行时间',
  `start_job_time` varchar(32) DEFAULT NULL COMMENT '启动构建任务时间时间',
  `end_job_time` varchar(32) DEFAULT NULL COMMENT '结束构建任务时间时间',
  `job_status` varchar(32) DEFAULT NULL COMMENT '构建任务是否成功',
  `push_image_start_time` varchar(32) DEFAULT NULL COMMENT '提交镜像时间',
  `push_image_end_time` varchar(32) DEFAULT NULL COMMENT '提交镜像完成时间',
  `update_service_start_time` varchar(32) DEFAULT NULL COMMENT '更新服务启动时间',
  `update_service_end_time` varchar(32) DEFAULT NULL COMMENT '更新服务结束时间',
  `end_time` varchar(32) DEFAULT NULL COMMENT '流程结束时间',
  `update_service_status` varchar(100) DEFAULT NULL COMMENT '更新服务状态',
  `update_service_errormsg` varchar(200) DEFAULT NULL COMMENT '更新失败状态信息',
  `build_job_errormsg` varchar(200) DEFAULT NULL COMMENT '更新失败状态信息',
  `build_job_status` varchar(200) DEFAULT NULL COMMENT '更新失败状态信息',
  PRIMARY KEY (`log_id`),
  KEY `idx_cloud_pipeline_log_all` (`pipeline_name`,`job_name`,`create_user`)
) ENGINE=InnoDB AUTO_INCREMENT=126 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `cloud_pipeline_log`
--

LOCK TABLES `cloud_pipeline_log` WRITE;
/*!40000 ALTER TABLE `cloud_pipeline_log` DISABLE KEYS */;
INSERT INTO `cloud_pipeline_log` VALUES (0,'pipeline1','pipeline','glusterfs-cluster','pipeline-1','job-567e4af553df794a2ac209193dc9a53d','zhaozq14','2018-02-04 10:48:31','000',NULL,'执行失败','获取服务状态超时,10分钟',7,1546,NULL,NULL,'2018-02-05 14:28:38',NULL,NULL,NULL,'2018-02-05 14:28:38','2018-02-05 14:53:45','2018-02-05 14:53:46','失败','获取服务状态超时,10分钟',NULL,NULL),(1,NULL,NULL,NULL,NULL,'zmc',NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,NULL,'2018-02-05 13:57:29','2018-02-05 13:58:21',NULL,NULL,NULL,'2018-02-05 13:58:23',NULL,NULL,NULL,NULL,NULL,NULL),(3,'pipeline1','storage','glusterfs-cluster','storage-1','zmc','zhaozq14','2018-02-03 22:26:10','000',NULL,NULL,NULL,0,NULL,NULL,'2018-02-05 13:57:29','2018-02-05 13:58:21',NULL,NULL,NULL,'2018-02-05 13:58:23',NULL,NULL,NULL,NULL,NULL,NULL),(15,'pipeline1','storage','glusterfs-cluster','storage-1','zmc','zhaozq14','2018-02-04 10:49:53','000',NULL,NULL,NULL,7,NULL,NULL,'2018-02-05 13:57:29','2018-02-05 13:58:21',NULL,NULL,NULL,'2018-02-05 13:58:23',NULL,NULL,NULL,NULL,NULL,NULL),(17,'pipeline1','app-1','glusterfs-cluster','app','job-638f82dff37407bb8564f6f55749b17c','zhaozq14','2018-02-04 18:31:41','000',NULL,NULL,NULL,7,0,NULL,'2018-02-05 13:57:29','2018-02-05 13:58:21',NULL,NULL,NULL,'2018-02-05 13:58:23',NULL,NULL,NULL,NULL,NULL,NULL),(19,'pipeline1','app-1','glusterfs-cluster','app','job-54ff8ffedc9c9941cd120eb75b115f27','zhaozq14','2018-02-04 18:36:15','000',NULL,'执行失败','获取服务状态超时,10分钟',7,964,NULL,'2018-02-05 13:57:29','2018-02-05 13:58:21',NULL,NULL,NULL,'2018-02-05 13:58:23',NULL,NULL,NULL,NULL,NULL,NULL),(21,'pipeline1','app-1','glusterfs-cluster','app','job-3bf78f0da551a1bee108a36fc8c21b22','zhaozq14','2018-02-04 18:45:00','000',NULL,NULL,NULL,7,0,NULL,'2018-02-05 13:57:29','2018-02-05 13:58:21',NULL,NULL,NULL,'2018-02-05 13:58:23',NULL,NULL,NULL,NULL,NULL,NULL),(23,'pipeline1','app-1','glusterfs-cluster','app','job-f6badc7c93deb3bcf4ba3a5524806d9d','zhaozq14','2018-02-04 18:53:08','000',NULL,NULL,NULL,7,0,NULL,'2018-02-05 13:57:29','2018-02-05 13:58:21',NULL,NULL,NULL,'2018-02-05 13:58:23',NULL,NULL,NULL,NULL,NULL,NULL),(25,'pipeline1','app-1','glusterfs-cluster','app','job-3a92cc2f8e2a62dd0866fe6605234ef6','zhaozq14','2018-02-04 19:01:11','000',NULL,'执行失败','构建失败',7,5,NULL,'2018-02-05 13:57:29','2018-02-05 13:58:21',NULL,NULL,NULL,'2018-02-05 13:58:23',NULL,NULL,NULL,NULL,NULL,NULL),(27,'pipeline1','app-1','glusterfs-cluster','app','job-4279b9dfd2a37126b585d1a2eabc5327','zhaozq14','2018-02-04 19:03:05','000',NULL,'执行失败','构建失败',7,14,NULL,'2018-02-05 13:57:29','2018-02-05 13:58:21',NULL,NULL,NULL,'2018-02-05 13:58:23',NULL,NULL,NULL,NULL,NULL,NULL),(29,'pipeline1','pipeline','glusterfs-cluster','pipeline-1','job-9b02b8462096a8e321c881708ba086bd','zhaozq14','2018-02-04 19:40:12','000',NULL,NULL,NULL,7,0,NULL,'2018-02-05 13:57:29','2018-02-05 13:58:21',NULL,NULL,NULL,'2018-02-05 13:58:23',NULL,NULL,NULL,NULL,NULL,NULL),(31,'pipeline1','pipeline','glusterfs-cluster','pipeline-1','job-9b55e52143139d0476445bef18ecc498','zhaozq14','2018-02-04 19:56:26','000',NULL,NULL,NULL,7,0,NULL,'2018-02-05 13:57:29','2018-02-05 13:58:21',NULL,NULL,NULL,'2018-02-05 13:58:23',NULL,NULL,NULL,NULL,NULL,NULL),(33,'pipeline1','pipeline','glusterfs-cluster','pipeline-1','job-2f839ab7e37c7a6b43812357aa85282a','zhaozq14','2018-02-04 21:40:01','000',NULL,NULL,NULL,7,0,NULL,'2018-02-05 13:57:29','2018-02-05 13:58:21',NULL,NULL,NULL,'2018-02-05 13:58:23',NULL,NULL,NULL,NULL,NULL,NULL),(35,'pipeline1','pipeline','glusterfs-cluster','pipeline-1','job-c35b185903ef443755c32bd06cf2ca67','zhaozq14','2018-02-04 21:51:30','000',NULL,'执行失败','构建失败',7,36,NULL,'2018-02-05 13:57:29','2018-02-05 13:58:21',NULL,NULL,NULL,'2018-02-05 13:58:23',NULL,NULL,NULL,NULL,NULL,NULL),(37,'pipeline1','pipeline','glusterfs-cluster','pipeline-1','job-9899b5d66854969147e7649a086c9a81','zhaozq14','2018-02-04 21:52:51','000',NULL,'执行失败','构建失败',7,16,NULL,'2018-02-05 13:57:29','2018-02-05 13:58:21',NULL,NULL,NULL,'2018-02-05 13:58:23',NULL,NULL,NULL,NULL,NULL,NULL),(39,'pipeline1','pipeline','glusterfs-cluster','pipeline-1','job-c57c5ca8238a508fcf0542035c2334dc','zhaozq14','2018-02-04 21:54:29','000',NULL,'执行失败','构建失败',7,18,NULL,'2018-02-05 13:57:29','2018-02-05 13:58:21',NULL,NULL,NULL,'2018-02-05 13:58:23',NULL,NULL,NULL,NULL,NULL,NULL),(41,'pipeline1','pipeline','glusterfs-cluster','pipeline-1','job-9e9f7556574cbe8d626fddcd3728dae5','zhaozq14','2018-02-04 21:55:57','000',NULL,'执行失败','构建失败',7,0,NULL,'2018-02-05 13:57:29','2018-02-05 13:58:21',NULL,NULL,NULL,'2018-02-05 13:58:23',NULL,NULL,NULL,NULL,NULL,NULL),(43,'pipeline1','pipeline','glusterfs-cluster','pipeline-1','job-48997098c671e075799e974b3c47456a','zhaozq14','2018-02-04 21:58:00','000',NULL,'执行失败','构建失败',7,1,NULL,'2018-02-05 13:57:29','2018-02-05 13:58:21',NULL,NULL,NULL,'2018-02-05 13:58:23',NULL,NULL,NULL,NULL,NULL,NULL),(45,'pipeline1','pipeline','glusterfs-cluster','pipeline-1','job-85bc85072115db3f5b3a4a81736e4586','zhaozq14','2018-02-04 21:59:13','000',NULL,'执行失败','构建失败',7,16,NULL,'2018-02-05 13:57:29','2018-02-05 13:58:21',NULL,NULL,NULL,'2018-02-05 13:58:23',NULL,NULL,NULL,NULL,NULL,NULL),(47,'pipeline1','pipeline','glusterfs-cluster','pipeline-1','job-a61e752703e2dedfae49c90fa4b1df6d','zhaozq14','2018-02-04 22:01:47','000',NULL,'执行失败','构建失败',7,0,NULL,'2018-02-05 13:57:29','2018-02-05 13:58:21',NULL,NULL,NULL,'2018-02-05 13:58:23',NULL,NULL,NULL,NULL,NULL,NULL),(49,'pipeline1','pipeline','glusterfs-cluster','pipeline-1','job-f9e0948e2f1763719a1f1b33b3c400a9','zhaozq14','2018-02-04 22:03:57','000',NULL,NULL,NULL,7,0,NULL,'2018-02-05 13:57:29','2018-02-05 13:58:21',NULL,NULL,NULL,'2018-02-05 13:58:23',NULL,NULL,NULL,NULL,NULL,NULL),(51,'pipeline1','pipeline','glusterfs-cluster','pipeline-1','job-d734dc91a22072d0fd48a24a0a35fe16','zhaozq14','2018-02-04 22:11:10','000',NULL,'执行失败','构建任务执行失败',7,0,NULL,'2018-02-05 13:57:29','2018-02-05 13:58:21',NULL,NULL,NULL,'2018-02-05 13:58:23',NULL,NULL,NULL,NULL,NULL,NULL),(53,'pipeline1','pipeline','glusterfs-cluster','pipeline-1','job-8b3b9dd4b6df2c622e05562f58cb8b88','zhaozq14','2018-02-04 22:13:36','000',NULL,NULL,NULL,7,0,NULL,'2018-02-05 13:57:29','2018-02-05 13:58:21',NULL,NULL,NULL,'2018-02-05 13:58:23',NULL,NULL,NULL,NULL,NULL,NULL),(55,'pipeline1','pipeline','glusterfs-cluster','pipeline-1','job-12835e79547c64cd5105cab726d12388','zhaozq14','2018-02-04 22:27:23','000',NULL,'执行失败','获取服务状态超时,10分钟',7,33176,NULL,'2018-02-05 13:57:29','2018-02-05 13:58:21',NULL,NULL,NULL,'2018-02-05 13:58:23',NULL,NULL,NULL,NULL,NULL,NULL),(57,'pipeline1','pipeline','glusterfs-cluster','pipeline-1','job-2020c75d456a08d6d3cbe676cbb00d3c','zhaozq14','2018-02-05 08:48:45','000',NULL,NULL,NULL,7,0,NULL,'2018-02-05 13:57:29','2018-02-05 13:58:21',NULL,NULL,NULL,'2018-02-05 13:58:23',NULL,NULL,NULL,NULL,NULL,NULL),(59,'pipeline1','pipeline','glusterfs-cluster','pipeline-1','job-a67f150e92982d029e095333bcafe40e','zhaozq14','2018-02-05 10:56:18','000',NULL,NULL,NULL,7,0,NULL,'2018-02-05 13:57:29','2018-02-05 13:58:21',NULL,NULL,NULL,'2018-02-05 13:58:23',NULL,NULL,NULL,NULL,NULL,NULL),(61,'pipeline1','pipeline','glusterfs-cluster','pipeline-1','job-131096ede8a06cecd0bac2a5eeefe8d8','zhaozq14','2018-02-04 10:48:31','000',NULL,NULL,NULL,7,0,NULL,NULL,'2018-02-05 13:58:21',NULL,NULL,NULL,'2018-02-05 13:58:23',NULL,NULL,NULL,NULL,NULL,NULL),(65,'pipeline1','pipeline','glusterfs-cluster','pipeline-1','job-039786d97541fec5f76c882c9f3362b2','zhaozq14','2018-02-04 10:48:31','000',NULL,NULL,NULL,7,0,NULL,'2018-02-05 14:57:41','2018-02-05 14:58:29',NULL,NULL,NULL,'2018-02-05 14:58:29',NULL,'2018-02-05 14:58:29',NULL,NULL,'构建任务成功','成功'),(67,'pipeline1','pipeline','glusterfs-cluster','pipeline-1','job-ca8f5914d1f4221ecd7a69712eaed51c','zhaozq14','2018-02-04 10:48:31','000',NULL,NULL,NULL,7,0,NULL,'2018-02-05 15:01:26','2018-02-05 15:02:17',NULL,NULL,NULL,'2018-02-05 15:02:17',NULL,'2018-02-05 15:02:17',NULL,NULL,'构建任务成功','成功'),(69,'pipeline1','pipeline','glusterfs-cluster','pipeline-1','job-dba20231497f599c87fdd4fbef490db4','zhaozq14','2018-02-04 10:48:31','000',NULL,NULL,NULL,7,0,NULL,'2018-02-05 15:04:11','2018-02-05 15:04:47',NULL,NULL,NULL,'2018-02-05 15:04:49',NULL,'2018-02-05 15:04:48',NULL,NULL,'构建任务成功','成功'),(71,'pipeline1','pipeline','glusterfs-cluster','pipeline-1','job-7a13ab7e04fb304d9acfe606ac888ff2','zhaozq14','2018-02-04 10:48:31','000',NULL,NULL,NULL,7,0,NULL,'2018-02-05 15:11:29','2018-02-05 15:12:18',NULL,NULL,NULL,'2018-02-05 15:12:19',NULL,'2018-02-05 15:12:18',NULL,NULL,'构建任务成功','成功'),(73,'pipeline1','pipeline','glusterfs-cluster','pipeline-1','job-731053ae2f1f881b6ea9509ad607b3e9','zhaozq14','2018-02-04 10:48:31','000',NULL,NULL,NULL,7,0,NULL,'2018-02-05 15:15:17','2018-02-05 15:16:25',NULL,NULL,NULL,'2018-02-05 15:16:26',NULL,'2018-02-05 15:16:26',NULL,NULL,'构建任务成功','成功'),(75,'pipeline1','pipeline','glusterfs-cluster','pipeline-1','job-ca3e3bcdfc80910e77d205ff8b74af9d','zhaozq14','2018-02-04 10:48:31','000',NULL,NULL,NULL,7,0,NULL,'2018-02-05 15:19:53','2018-02-05 15:20:32',NULL,NULL,NULL,'2018-02-05 15:20:32',NULL,'2018-02-05 15:20:32',NULL,NULL,'构建任务成功','成功'),(77,'pipeline1','pipeline','glusterfs-cluster','pipeline-1','job-5059b4f84f3b9c69b813f1d35bd79288','zhaozq14','2018-02-04 10:48:31','000',NULL,NULL,NULL,7,0,NULL,'2018-02-05 15:30:20','2018-02-05 15:31:06',NULL,NULL,NULL,'2018-02-05 15:31:06','2018-02-05 15:37:30','2018-02-05 15:31:06','成功','reg2.asura.com:49000/zmc/zmc:20180205-153019 True 1','构建任务成功','成功'),(79,'pipeline1','pipeline','glusterfs-cluster','pipeline-1','job-c5a2c0f86d6359f6ec657ab14ce81410','zhaozq14','2018-02-04 10:48:31','000',NULL,NULL,NULL,7,0,NULL,'2018-02-05 15:38:18','2018-02-05 15:39:18',NULL,NULL,NULL,'2018-02-05 15:39:19','2018-02-05 15:41:55','2018-02-05 15:39:18','成功','reg2.asura.com:49000/zmc/zmc:20180205-153818 True 1','构建任务成功','成功'),(81,'pipeline1','pipeline','glusterfs-cluster','pipeline-1','job-35a0f91ec2ab222d539a393f7dc91098','zhaozq14','2018-02-04 10:48:31','000',NULL,NULL,NULL,7,0,NULL,'2018-02-05 15:42:39','2018-02-05 15:43:20',NULL,NULL,NULL,'2018-02-05 15:43:20','2018-02-05 15:43:31','2018-02-05 15:43:34','成功','reg2.asura.com:49000/zmc/zmc:20180205-154236 True 1','构建任务成功','成功'),(83,'pipeline1','pipeline','glusterfs-cluster','pipeline-1','job-3a91a9c34fa883b474d592e9f675b7dc','zhaozq14','2018-02-04 10:48:31','000',NULL,NULL,NULL,7,0,'2018-02-05 15:56:27','2018-02-05 15:56:27','2018-02-05 15:56:28',NULL,NULL,NULL,NULL,NULL,'2018-02-05 15:56:28',NULL,NULL,'构建任务执行失败','失败'),(85,'pipeline1','pipeline','glusterfs-cluster','pipeline-1','job-54142bc384548cbc374e667001cec6d1','zhaozq14','2018-02-04 10:48:31','000',NULL,'执行失败','reg2.asura.com:49000/zmc/zmc:20180205-160551 True 1',7,51,'2018-02-05 16:05:51','2018-02-05 16:05:52','2018-02-05 16:06:32',NULL,NULL,NULL,'2018-02-05 16:06:32','2018-02-05 16:06:43','2018-02-05 16:06:43','成功','reg2.asura.com:49000/zmc/zmc:20180205-160551 True 1','构建任务成功','成功'),(87,'pipeline1','pipeline','glusterfs-cluster','pipeline-1','job-bb8009a89e032870fb0b092080ff9a62','zhaozq14','2018-02-04 10:48:31','000',NULL,'执行失败','reg2.asura.com:49000/zmc/zmc:20180205-161432 True 1',7,64,'2018-02-05 16:14:32','2018-02-05 16:14:38','2018-02-05 16:15:30',NULL,NULL,NULL,'2018-02-05 16:15:31','2018-02-05 16:15:41','2018-02-05 16:15:42','成功','reg2.asura.com:49000/zmc/zmc:20180205-161432 True 1','构建任务成功','成功'),(89,'pipeline1','pipeline','glusterfs-cluster','pipeline-1','job-6337ebc22fbc78b899944616cbf521eb','zhaozq14','2018-02-04 10:48:31','000',NULL,'执行成功','更新成功',7,71,'2018-02-05 16:59:04','2018-02-05 16:59:04','2018-02-05 16:59:56',NULL,NULL,NULL,'2018-02-05 16:59:56','2018-02-05 17:00:07','2018-02-05 17:00:17','成功','reg2.asura.com:49000/zmc/zmc:20180205-165904 True 1','构建任务成功','成功'),(91,'pipeline1','pipeline','glusterfs-cluster','pipeline-1','job-2c9a56725d0a4945762c3b6afc760390','zhaozq14','2018-02-04 10:48:31','000',NULL,NULL,NULL,7,0,'2018-02-06 20:45:45','2018-02-06 20:45:51','2018-02-06 20:47:06',NULL,NULL,NULL,NULL,NULL,'2018-02-06 20:47:06',NULL,NULL,'构建任务执行失败','失败'),(93,'pipeline-1','app-scale','glusterfs-cluster','app-scale','job-2c25e61772aaf5acd99346fa388b1af7','admin','2018-02-09 09:02:51','*',NULL,NULL,NULL,11,0,'2018-02-09 09:02:57','2018-02-09 09:02:57','2018-02-09 09:02:58',NULL,NULL,NULL,NULL,NULL,'2018-02-09 09:02:58',NULL,NULL,'执行构建失败','构建失败'),(95,'pipeline-1','app-scale','glusterfs-cluster','app-scale','job-c4c7301f042e2f2eaa2774f636038b80','admin','2018-02-09 09:02:51','*',NULL,NULL,NULL,11,0,'2018-02-09 09:04:41','2018-02-09 09:04:43','2018-02-09 09:05:21',NULL,NULL,NULL,'2018-02-09 09:05:21','2018-02-09 09:05:31','2018-02-09 09:05:32','成功','//:20180209-090441 True 4','构建任务成功','成功'),(97,'pipeline-1','app-scale','glusterfs-cluster','app-scale','job-eb3bee4f917cdae097a2cc4851144dad','admin','2018-02-09 09:02:51','*',NULL,NULL,NULL,11,0,'2018-02-09 09:13:27','2018-02-09 09:13:27','2018-02-09 09:13:28',NULL,NULL,NULL,NULL,NULL,'2018-02-09 09:13:28',NULL,NULL,'构建任务执行失败','失败'),(99,'pipeline-1','app-scale','glusterfs-cluster','app-scale','job-0bb718028661f739f5764d8b6d149f49','admin','2018-02-09 09:02:51','*',NULL,NULL,NULL,11,0,'2018-02-09 09:19:46','2018-02-09 09:19:48','2018-02-09 09:20:22',NULL,NULL,NULL,'2018-02-09 09:20:23',NULL,'2018-02-09 09:20:22',NULL,NULL,'构建任务成功','成功'),(101,'pipeline-1','app-scale','glusterfs-cluster','app-scale','job-e6d0c5692bbead107a4a7db286c01037','admin','2018-02-09 09:02:51','*',NULL,NULL,NULL,11,0,'2018-02-09 09:32:30','2018-02-09 09:32:31','2018-02-09 09:35:06',NULL,NULL,NULL,NULL,'2018-02-09 09:35:06','2018-02-09 09:35:06','失败','拉取服务失败了','执行构建成功','构建成功'),(103,'pipeline-1','maxcpu','glusterfs-cluster','maxcpu','job-726709feac4b279719f15ea07a248faf','admin','2018-02-09 09:02:51','*',NULL,NULL,NULL,11,0,'2018-02-09 09:37:32','2018-02-09 09:37:33','2018-02-09 09:37:33',NULL,NULL,NULL,NULL,NULL,'2018-02-09 09:37:33',NULL,NULL,'构建任务执行失败','失败'),(105,'pipeline-1','maxcpu','glusterfs-cluster','maxcpu','job-3dfd6105b71633962452f508418437bc','admin','2018-02-09 09:02:51','*',NULL,NULL,NULL,11,0,'2018-02-09 09:39:40','2018-02-09 09:39:40','2018-02-09 09:39:42',NULL,NULL,NULL,NULL,NULL,'2018-02-09 09:39:42',NULL,NULL,'构建任务执行失败','失败'),(107,'pipeline-1','maxcpu','glusterfs-cluster','maxcpu','job-0afd79e82107c533484650176c466aff','admin','2018-02-09 09:02:51','*',NULL,NULL,NULL,11,0,'2018-02-09 14:07:56','2018-02-09 14:07:56','2018-02-09 14:08:59',NULL,NULL,NULL,'2018-02-09 14:08:59',NULL,'2018-02-09 14:08:59',NULL,NULL,'构建任务成功','成功'),(109,'pipeline-1','maxcpu','glusterfs-cluster','maxcpu','job-3de3f7916ec8c2ce409097a3bf1cd672','admin','2018-02-09 09:02:51','*',NULL,NULL,NULL,11,0,'2018-02-09 14:28:37','2018-02-09 14:28:37','2018-02-09 14:28:37',NULL,NULL,NULL,NULL,NULL,'2018-02-09 14:28:37',NULL,NULL,'构建任务执行失败','失败'),(111,'pipeline-1','maxcpu','glusterfs-cluster','maxcpu','job-7eb7db2603613a3dd3f121aebbaca69c','admin','2018-02-09 09:02:51','*',NULL,NULL,NULL,11,0,'2018-02-09 14:39:16','2018-02-09 14:39:16','2018-02-09 14:39:59',NULL,NULL,NULL,NULL,NULL,'2018-02-09 14:40:01',NULL,NULL,'构建任务成功','成功'),(113,'pipeline-1','maxcpu','glusterfs-cluster','maxcpu','job-73eb47ecc1daaa46e2956a16a3ce2272','admin','2018-02-09 09:02:51','*','','执行失败','获取仓库服务器错误',11,38,'2018-02-09 14:51:23','','','','','','','','','','','',''),(115,'pipeline-1','maxcpu','glusterfs-cluster','maxcpu','job-f43492e22aa4b4f731a6f624f776f96b','admin','2018-02-09 09:02:51','*',NULL,NULL,NULL,11,0,'2018-02-09 14:58:49','2018-02-09 14:58:49','2018-02-09 14:59:25',NULL,NULL,NULL,'2018-02-09 14:59:26',NULL,'2018-02-09 14:59:26',NULL,NULL,'构建任务成功','成功'),(117,'pipeline-1','maxcpu','glusterfs-cluster','maxcpu','job-407214be775d5080523eb53a3cbcc4aa','admin','2018-02-09 09:02:51','*',NULL,NULL,NULL,11,0,'2018-02-09 15:01:10','2018-02-09 15:01:11','2018-02-09 15:01:46',NULL,NULL,NULL,'2018-02-09 15:01:46',NULL,'2018-02-09 15:01:46',NULL,NULL,'构建任务成功','成功'),(119,'pipeline-1','maxcpu','glusterfs-cluster','maxcpu','job-26d5f977f5f4a0479d9f9535c0412457','admin','2018-02-09 09:02:51','*',NULL,NULL,NULL,11,0,'2018-02-09 15:03:20','2018-02-09 15:03:20','2018-02-09 15:03:20',NULL,NULL,NULL,NULL,NULL,'2018-02-09 15:03:20',NULL,NULL,'构建任务执行失败','失败'),(121,'pipeline-1','maxcpu','glusterfs-cluster','maxcpu','job-5cb1d18f4cad87d82678a8bcefb6896b','admin','2018-02-09 09:02:51','*',NULL,NULL,NULL,11,0,'2018-02-09 15:12:55','2018-02-09 15:12:55','2018-02-09 15:13:45',NULL,NULL,NULL,'2018-02-09 15:13:45',NULL,'2018-02-09 15:13:45',NULL,NULL,'构建任务成功','成功'),(123,'pipeline-1','test-cloud','glusterfs-cluster','test-cloud-service','job-634e102baeb42aea2ca3cbeca0b0770f','admin','2018-02-09 09:02:51','*',NULL,NULL,NULL,11,0,'2018-02-09 16:06:04','2018-02-09 16:06:04','2018-02-09 16:06:05',NULL,NULL,NULL,NULL,NULL,'2018-02-09 16:06:06',NULL,NULL,'构建任务执行失败','失败'),(125,'pipeline-1','test-cloud','glusterfs-cluster','test-cloud-service','job-8a699d9f8237be5e8851493bf7e730bf','admin','2018-02-09 09:02:51','*','','执行成功','更新成功',11,59,'2018-02-09 16:10:09','','','','','','','','2018-02-09 16:11:08','','','','');
/*!40000 ALTER TABLE `cloud_pipeline_log` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `cloud_quota`
--

DROP TABLE IF EXISTS `cloud_quota`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `cloud_quota` (
  `quota_id` int(11) NOT NULL AUTO_INCREMENT,
  `quota_name` varchar(36) DEFAULT NULL COMMENT '配额名称',
  `description` varchar(300) DEFAULT NULL COMMENT '配额描述信息',
  `quota_cpu` varchar(20) DEFAULT NULL COMMENT 'cpu配额多少核心',
  `quota_memory` varchar(20) DEFAULT NULL COMMENT '内存配额多少MB',
  `create_time` varchar(32) DEFAULT NULL COMMENT '创建时间',
  `create_user` varchar(32) DEFAULT NULL COMMENT '创建用户',
  `resource_name` varchar(32) DEFAULT NULL COMMENT '资源名称',
  `last_modify_time` varchar(32) DEFAULT NULL COMMENT '最近修改时间',
  `last_modify_user` varchar(32) DEFAULT NULL COMMENT '最近修改用户',
  `status` varchar(32) DEFAULT NULL,
  `pod_number` int(11) DEFAULT NULL COMMENT '限制pod数量',
  `service_number` int(11) DEFAULT NULL COMMENT '限制应用数量',
  `app_number` int(11) DEFAULT NULL COMMENT '限制应用数量',
  `lb_number` int(11) DEFAULT NULL COMMENT '限制负载均衡数量',
  `job_number` int(11) DEFAULT NULL COMMENT '限制发布任务数量',
  `pipeline_number` int(11) DEFAULT NULL COMMENT '限制流水线数量',
  `user_name` varchar(100) DEFAULT NULL COMMENT '受限人名称',
  `group_name` varchar(100) DEFAULT NULL COMMENT '受限业务线名称',
  `registry_group_number` int(11) DEFAULT NULL COMMENT '镜像仓库组数量',
  `docker_file_number` int(11) DEFAULT NULL COMMENT '镜像仓库组数量',
  PRIMARY KEY (`quota_id`),
  UNIQUE KEY `uidx_cloud_quota_quota_name` (`quota_name`),
  UNIQUE KEY `uidx_cloud_quota_user_name_group_name` (`user_name`,`quota_name`,`group_name`)
) ENGINE=InnoDB AUTO_INCREMENT=18 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `cloud_quota`
--

LOCK TABLES `cloud_quota` WRITE;
/*!40000 ALTER TABLE `cloud_quota` DISABLE KEYS */;
INSERT INTO `cloud_quota` VALUES (5,'dfsad','赵云配额','8','10249','2018-01-05 14:04:43','admin','','2018-02-12 07:07:48','admin','',30,10,10,10,10,10,'zhaozq14','',NULL,NULL),(15,'admin-quota','admin','32','40512','2018-02-12 07:05:14','admin','','2018-02-27 10:24:42','admin','',24,14,11,3,4,4,'admin','',1,1),(17,'crmtm-group-quota','crm团队资源配额','11','3512','2018-02-12 07:15:47','admin','','2018-02-12 11:31:27','zhaozq14','',1,1,1,1,4,1,'','CRM研发部',12,4);
/*!40000 ALTER TABLE `cloud_quota` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `cloud_registry_group`
--

DROP TABLE IF EXISTS `cloud_registry_group`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `cloud_registry_group` (
  `group_id` int(11) NOT NULL AUTO_INCREMENT,
  `create_time` varchar(32) DEFAULT NULL COMMENT '创建时间',
  `create_user` varchar(32) DEFAULT NULL COMMENT '创建用户',
  `image_number` int(11) DEFAULT NULL COMMENT '镜像数量',
  `group_type` varchar(20) DEFAULT NULL COMMENT '镜像类型,分为共有和私有',
  `last_modify_time` varchar(32) DEFAULT NULL COMMENT '最近修改时间',
  `last_modify_user` varchar(20) DEFAULT NULL COMMENT '最近修改用户',
  `size_totle` int(11) DEFAULT NULL COMMENT '镜像总大小',
  `group_name` varchar(100) DEFAULT NULL COMMENT '镜像组名称',
  `server_domain` varchar(100) DEFAULT NULL COMMENT '所在镜像仓库域名',
  `cluster_name` varchar(100) DEFAULT NULL COMMENT '所在集群名称',
  `tag_number` int(11) DEFAULT NULL,
  PRIMARY KEY (`group_id`),
  UNIQUE KEY `uidx_cloud_registry_group_name` (`group_name`,`server_domain`,`cluster_name`)
) ENGINE=InnoDB AUTO_INCREMENT=22 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `cloud_registry_group`
--

LOCK TABLES `cloud_registry_group` WRITE;
/*!40000 ALTER TABLE `cloud_registry_group` DISABLE KEYS */;
INSERT INTO `cloud_registry_group` VALUES (7,'2018-01-27 18:56:25','zhaozq14',2,'公开','2018-02-07 11:47:42','zhaozq14',0,'zmc','reg2.asura.com','glusterfs-cluster',2),(9,'2018-01-28 20:53:52','zhaozq14',0,'公开','2018-02-08 13:46:17','admin',0,'zmc/a/b','reg2.asura.com','glusterfs-cluster',0),(13,'2018-02-07 11:37:10','zhaozq14',0,'公开','2018-02-07 11:47:54','zhaozq14',0,'crm','reg2.asura.com','glusterfs-cluster',0),(17,'2018-02-08 10:19:28','admin',1,'公开','2018-02-08 13:47:22','admin',0,'hd','reg2.asura.com','glusterfs-cluster',0),(19,'2018-02-12 08:17:39','test',1,'公开','2018-02-12 08:17:39','test',0,'test','reg2.asura.com','glusterfs-cluster',0),(21,'2018-02-12 08:36:26','test',0,'公开','2018-02-12 08:36:26','test',0,'quota2','reg2.asura.com','glusterfs-cluster',0);
/*!40000 ALTER TABLE `cloud_registry_group` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `cloud_registry_permissions`
--

DROP TABLE IF EXISTS `cloud_registry_permissions`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `cloud_registry_permissions` (
  `permissions_id` int(11) NOT NULL AUTO_INCREMENT,
  `user_type` varchar(10) DEFAULT NULL COMMENT '用户类型,用户或组',
  `registry_server` varchar(32) DEFAULT NULL COMMENT '注册中心地址',
  `service_name` varchar(100) DEFAULT NULL COMMENT '注册中心名称',
  `user_name` varchar(32) DEFAULT NULL COMMENT '用户名称',
  `groups_name` varchar(100) DEFAULT NULL COMMENT '用户组名称',
  `project` varchar(100) DEFAULT NULL COMMENT '项目名称',
  `image_name` varchar(100) DEFAULT NULL COMMENT '镜像名称',
  `create_time` varchar(32) DEFAULT NULL COMMENT '创建时间',
  `create_user` varchar(32) DEFAULT NULL COMMENT '创建用户',
  `last_modify_time` varchar(32) DEFAULT NULL COMMENT '最近修改时间',
  `last_modify_user` varchar(32) DEFAULT NULL COMMENT '最近修改用户',
  `action` varchar(30) DEFAULT NULL COMMENT '操作权限,pull,push',
  `cluster_name` varchar(100) DEFAULT NULL COMMENT '集群名称',
  `description` varchar(100) DEFAULT NULL COMMENT '描述信息',
  PRIMARY KEY (`permissions_id`),
  KEY `cloud_registry_permissions_p_u_g` (`project`,`user_name`,`groups_name`)
) ENGINE=InnoDB AUTO_INCREMENT=28 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `cloud_registry_permissions`
--

LOCK TABLES `cloud_registry_permissions` WRITE;
/*!40000 ALTER TABLE `cloud_registry_permissions` DISABLE KEYS */;
INSERT INTO `cloud_registry_permissions` VALUES (11,NULL,NULL,'registry',NULL,'21','asms',NULL,'2018-01-22 16:58:35','zhaozq14','2018-01-22 17:32:50',NULL,'pull,push','asdfasdfdasf','所有用户权限测试'),(25,NULL,NULL,'registry','zhaozq14',NULL,'zmc',NULL,'2018-01-28 17:59:24','zhaozq14','2018-01-28 17:59:24',NULL,'pull','glusterfs-cluster','web负载均衡器');
/*!40000 ALTER TABLE `cloud_registry_permissions` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `cloud_registry_server`
--

DROP TABLE IF EXISTS `cloud_registry_server`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `cloud_registry_server` (
  `server_id` int(11) NOT NULL AUTO_INCREMENT,
  `server_address` varchar(32) DEFAULT NULL,
  `server_domain` varchar(32) DEFAULT NULL,
  `last_modify_time` varchar(32) DEFAULT NULL COMMENT '最近修改时间',
  `last_modify_user` varchar(20) DEFAULT NULL COMMENT '最近修改用户',
  `create_time` varchar(32) DEFAULT NULL COMMENT '创建时间',
  `create_user` varchar(32) DEFAULT NULL COMMENT '创建用户',
  `groups` varchar(32) DEFAULT NULL COMMENT '组名称,不同的服务属于不同的组',
  `images_number` int(11) DEFAULT '0' COMMENT '镜像数量',
  `description` varchar(32) DEFAULT NULL COMMENT '描述信息',
  `groups_id` int(11) DEFAULT NULL,
  `prefix` varchar(32) DEFAULT NULL COMMENT '镜像前缀 如 online test develop ${user}',
  `username` varchar(32) DEFAULT NULL COMMENT '镜像用户名名',
  `password` varchar(32) DEFAULT NULL COMMENT '镜像密码',
  `cluster_name` varchar(200) DEFAULT NULL COMMENT '集群名称',
  `registry_type` varchar(32) DEFAULT NULL COMMENT 'public,private,仓库类型,公共仓库所有人都可以上传',
  `name` varchar(32) DEFAULT NULL COMMENT '仓库名称',
  `auth_server` varchar(200) DEFAULT NULL,
  `admin` varchar(32) DEFAULT NULL COMMENT '管理员',
  `access` varchar(300) DEFAULT NULL COMMENT '访问信息',
  `entname` varchar(100) DEFAULT NULL COMMENT '环境名称',
  `mount_path` varchar(1000) DEFAULT NULL COMMENT '仓库挂载路径',
  `labels` varchar(30) DEFAULT NULL COMMENT '安装集群的标签',
  `replicas` int(11) DEFAULT NULL COMMENT '副本数量',
  PRIMARY KEY (`server_id`),
  UNIQUE KEY `uidx_cloud_registry_server_clsuter_name` (`cluster_name`,`name`),
  UNIQUE KEY `uidx_cloud_registry_server_server_domain` (`server_domain`)
) ENGINE=InnoDB AUTO_INCREMENT=29 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `cloud_registry_server`
--

LOCK TABLES `cloud_registry_server` WRITE;
/*!40000 ALTER TABLE `cloud_registry_server` DISABLE KEYS */;
INSERT INTO `cloud_registry_server` VALUES (27,'reg2.asura.com:49000','reg2.asura.com','2018-03-02 11:34:11','admin','2018-01-31 09:36:25','zhaozq14','',0,'存储服务',0,'','','V1ZkU2RHRlhORDA9','glusterfs-cluster',NULL,'registry','https://registry.asura.com:5001/auth','admin','容器内&nbsp;<br>registry.registryv2--registryv2:49000<br>集群外&nbsp;<br><a target=\'_blank\' href=\'https://reg2.asura.com:49000/v2/\'>reg2.asura.com:49000</a>','生产环境',NULL,NULL,NULL);
/*!40000 ALTER TABLE `cloud_registry_server` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `cloud_storage`
--

DROP TABLE IF EXISTS `cloud_storage`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `cloud_storage` (
  `storage_id` int(11) NOT NULL AUTO_INCREMENT,
  `last_modify_time` varchar(32) DEFAULT NULL COMMENT '最近修改时间',
  `last_modify_user` varchar(20) DEFAULT NULL COMMENT '最近修改用户',
  `create_time` varchar(32) DEFAULT NULL COMMENT '创建时间',
  `create_user` varchar(32) DEFAULT NULL COMMENT '创建用户',
  `description` varchar(32) DEFAULT NULL COMMENT '描述信息',
  `cluster_name` varchar(200) DEFAULT NULL COMMENT '集群名称',
  `storage_type` varchar(200) DEFAULT NULL COMMENT 'glusterfs, nfs, host',
  `storage_size` varchar(10) DEFAULT NULL COMMENT '存储大小,单位GB',
  `storage_format` varchar(10) DEFAULT NULL COMMENT '存储格式',
  `name` varchar(100) DEFAULT NULL COMMENT '存储名称',
  `storage_server` varchar(100) DEFAULT NULL COMMENT '存储服务器地址',
  `shared_type` varchar(10) DEFAULT NULL COMMENT '存储共享类型,分独享还是共享的',
  `entname` varchar(100) DEFAULT NULL COMMENT '环境名称',
  `status` varchar(100) DEFAULT NULL COMMENT '使用状态',
  PRIMARY KEY (`storage_id`),
  UNIQUE KEY `uidx_cloud_storage_name` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=28 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `cloud_storage`
--

LOCK TABLES `cloud_storage` WRITE;
/*!40000 ALTER TABLE `cloud_storage` DISABLE KEYS */;
INSERT INTO `cloud_storage` VALUES (9,'2018-01-31 14:43:21',NULL,'2018-01-31 14:43:21',NULL,'zmc存储','glusterfs-cluster','Nfs','512',NULL,'zmc-nfs',NULL,'0','生产环境',NULL),(13,'2018-01-31 14:44:02',NULL,'2018-01-31 14:44:02',NULL,'php-nfs','glusterfs-cluster','Nfs','512',NULL,'php-nfs',NULL,'0','生产环境',NULL),(23,'2018-01-31 16:41:21',NULL,'2018-01-31 16:41:21',NULL,'web负载均衡器','glusterfs-cluster','Nfs','512',NULL,'crm-nfs',NULL,'0','生产环境',NULL),(25,'2018-02-03 06:48:29',NULL,'2018-02-03 06:48:29',NULL,'存储卷描述信息','glusterfs-cluster','Nfs','512',NULL,'cmd-nfs',NULL,'0','生产环境',NULL),(27,'2018-02-22 17:41:38','admin','2018-02-22 17:08:33','admin','glusterfs给crm的1gb存储','glusterfs-cluster','Glusterfs','1024','','glusterfs-crm-1g','','1','生产环境','');
/*!40000 ALTER TABLE `cloud_storage` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `cloud_storage_mount_info`
--

DROP TABLE IF EXISTS `cloud_storage_mount_info`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `cloud_storage_mount_info` (
  `mount_id` int(11) NOT NULL AUTO_INCREMENT,
  `service_name` varchar(100) DEFAULT NULL COMMENT '服务名称',
  `app_name` varchar(100) DEFAULT NULL COMMENT '应用名称',
  `create_time` varchar(32) DEFAULT NULL COMMENT '创建时间',
  `create_user` varchar(32) DEFAULT NULL COMMENT '创建用户',
  `cluster_name` varchar(100) DEFAULT NULL COMMENT '集群名称',
  `mount_path` varchar(300) DEFAULT NULL COMMENT '容器挂载路径',
  `model` varchar(32) DEFAULT NULL COMMENT '读写权限',
  `storage_server` varchar(100) DEFAULT NULL COMMENT '存储服务器',
  `storage_type` varchar(100) DEFAULT NULL COMMENT '存储类型',
  `status` varchar(32) DEFAULT NULL,
  `storage_name` varchar(100) DEFAULT NULL COMMENT '存储卷名称',
  `resource_name` varchar(100) DEFAULT NULL,
  PRIMARY KEY (`mount_id`),
  KEY `idx_name_cluster` (`storage_name`,`cluster_name`,`service_name`)
) ENGINE=InnoDB AUTO_INCREMENT=14 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `cloud_storage_mount_info`
--

LOCK TABLES `cloud_storage_mount_info` WRITE;
/*!40000 ALTER TABLE `cloud_storage_mount_info` DISABLE KEYS */;
INSERT INTO `cloud_storage_mount_info` VALUES (3,'storage-1',NULL,'2018-01-31 15:36:27','zhaozq14','glusterfs-cluster','/mnt',NULL,NULL,NULL,'1','crm-nfs','dfsad'),(5,'storage-1',NULL,'2018-02-04 16:57:44','zhaozq14','glusterfs-cluster','/mnt',NULL,NULL,'共享型','1','crm-nfs','dfsad'),(7,'glusterfs-app',NULL,'2018-02-22 21:31:37','admin','glusterfs-cluster','/mnt',NULL,NULL,NULL,'1','glusterfs-crm-1g','admin-quota'),(9,'glusterfs-app',NULL,'2018-02-22 21:35:16','admin','glusterfs-cluster','/mnt',NULL,NULL,'Glusterfs','1','glusterfs-crm-1g','admin-quota'),(11,'glusterfs-app',NULL,'2018-02-22 21:42:24','admin','glusterfs-cluster','/mnt',NULL,NULL,'Glusterfs','1','glusterfs-crm-1g','admin-quota'),(13,'config-test',NULL,'2018-02-27 10:29:28','admin','glusterfs-cluster','/mnt',NULL,NULL,'Nfs','1','zmc-nfs','admin-quota');
/*!40000 ALTER TABLE `cloud_storage_mount_info` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `cloud_storage_server`
--

DROP TABLE IF EXISTS `cloud_storage_server`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `cloud_storage_server` (
  `server_id` int(11) NOT NULL AUTO_INCREMENT,
  `server_address` varchar(32) DEFAULT NULL,
  `last_modify_time` varchar(32) DEFAULT NULL COMMENT '最近修改时间',
  `last_modify_user` varchar(20) DEFAULT NULL COMMENT '最近修改用户',
  `create_time` varchar(32) DEFAULT NULL COMMENT '创建时间',
  `create_user` varchar(32) DEFAULT NULL COMMENT '创建用户',
  `description` varchar(32) DEFAULT NULL COMMENT '描述信息',
  `cluster_name` varchar(200) DEFAULT NULL COMMENT '集群名称',
  `storage_type` varchar(200) DEFAULT NULL COMMENT 'glusterfs, nfs, host',
  `used_type` varchar(200) DEFAULT NULL COMMENT '独享型，共享型',
  `entname` varchar(100) DEFAULT NULL COMMENT '环境名称',
  `host_path` varchar(300) DEFAULT NULL COMMENT '磁盘路径或目录路径',
  PRIMARY KEY (`server_id`),
  UNIQUE KEY `uidx_cloud_storage_server_storage_type_cluster_name` (`storage_type`,`cluster_name`)
) ENGINE=InnoDB AUTO_INCREMENT=10 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `cloud_storage_server`
--

LOCK TABLES `cloud_storage_server` WRITE;
/*!40000 ALTER TABLE `cloud_storage_server` DISABLE KEYS */;
INSERT INTO `cloud_storage_server` VALUES (5,NULL,'2018-02-08 09:58:00','zhaozq14','2018-02-08 09:58:00','zhaozq14','nfs存储提供者','glusterfs-cluster','Nfs',NULL,'生产环境',NULL),(9,NULL,'2018-02-22 14:29:44','admin','2018-02-22 14:29:44','admin','glusterfs服务提供','glusterfs-cluster','Glusterfs',NULL,'生产环境','/dev/vdb');
/*!40000 ALTER TABLE `cloud_storage_server` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `cloud_user_groups`
--

DROP TABLE IF EXISTS `cloud_user_groups`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `cloud_user_groups` (
  `groups_id` int(11) NOT NULL AUTO_INCREMENT,
  `create_time` varchar(32) DEFAULT NULL COMMENT '创建时间',
  `create_user` varchar(32) DEFAULT NULL COMMENT '创建用户',
  `last_modify_time` varchar(32) DEFAULT NULL COMMENT '最近修改时间',
  `last_modify_user` varchar(32) DEFAULT NULL COMMENT '最近修改用户',
  `groups_name` varchar(100) DEFAULT NULL COMMENT '组名称',
  `users` text COMMENT '组成员,用逗号分隔',
  `description` varchar(100) DEFAULT NULL COMMENT '描述信息',
  PRIMARY KEY (`groups_id`)
) ENGINE=InnoDB AUTO_INCREMENT=24 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `cloud_user_groups`
--

LOCK TABLES `cloud_user_groups` WRITE;
/*!40000 ALTER TABLE `cloud_user_groups` DISABLE KEYS */;
INSERT INTO `cloud_user_groups` VALUES (21,'2018-01-20 13:49:37','zhaozq14','2018-02-12 08:38:43','test','CRM研发部','zhaozq14,test','crm研发部'),(23,'2018-01-20 13:49:53','zhaozq14','2018-01-25 21:47:31','zhaozq14','生活服务家修研发部','zhaozq14','生活服务家修研发部');
/*!40000 ALTER TABLE `cloud_user_groups` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2018-03-09 12:30:13
alter table cloud_cluster_hosts add unique index uidx_cloud_cluster_hosts_ip_api_port_host_type(host_ip,api_port,host_type);