package app


const SelectCloudApp = "select service_yaml,entname,yaml,node_port,uuid,is_service,cluster_name,app_id,app_name,status,last_modify_time,app_labels,json_data,app_type,resource_name,create_time,create_user,last_modify_user,last_update_time from cloud_app"
const UpdateCloudApp = "update cloud_app"
const InsertCloudApp = "insert into cloud_app"
const DeleteCloudApp = "delete from cloud_app"
const GetAppName  = "select app_id, app_name,create_user,cluster_name,entname from cloud_app"
const SelectUserApp = `select app_id from cloud_app where create_user in (?) `


const SelectCloudAppTemplate = "select yaml,service_name,ent,cluster,template_name,create_time,create_user,last_modify_time,last_modify_user,resource_name,description,template_id from cloud_app_template"
const UpdateCloudAppTemplate = "update cloud_app_template"
const InsertCloudAppTemplate = "insert into cloud_app_template"
const DeleteCloudAppTemplate = "delete from cloud_app_template"
const SelectServiceYaml  = "select yaml from cloud_app_service where service_name in (%s)"

const SelectCloudAppConfigure = "select create_user,entname,last_modify_time,cluster_name,description,configure_id,configure_name,create_time,last_modify_user from cloud_app_configure"
const UpdateCloudAppConfigure = "update cloud_app_configure"
const UpdateCloudAppConfigureExclude = "CreateTime,CreateUser,ConfigureName"
const InsertCloudAppConfigure = "insert into cloud_app_configure"
const DeleteCloudAppConfigure = "delete from cloud_app_configure"
const SelectCloudAppConfigSearch  = ` and (configure_name like "%?%" or description like "%?%")`

const SelectCloudConfigData = "select data_name,data_id,create_time,last_modify_user,create_user,last_modify_time,configure_id,configure_name,data from cloud_config_data"
const UpdateCloudConfigData = "update cloud_config_data"
const InsertCloudConfigData = "insert into cloud_config_data"
const DeleteCloudConfigData = "delete from cloud_config_data"
const UpdateConfigDataExclude = "CreateTime,CreateUser,DataName"
const SelectConfigDataWhere = ` where 1=1 and (configure_name like "%?%" or description like "%?%")`

const SelectServiceNameSpace = `select distinct concat(app_name,"--",resource_name) as service_name, entname,cluster_name from cloud_app_service`
const SelectCloudAppService = "select termination_seconds,log_path,max_surge,domain,service_version,entname,max_unavailable,min_ready,envs,app_name,service_lables_data,cluster_name,network_mode,config,replicas_min,health_data,service_id,status,image_tag,replicas_max,create_user,service_labels,lb_data,service_type,json_data,container_port,env_file,storage_data,resource_name,cpu,lb_name,app_labels,yaml,configure_data,create_time,last_modify_user,last_update_time,deploy_type,replicas,service_name,last_modify_time,memory from cloud_app_service"
const UpdateCloudAppService = "update cloud_app_service"
const UpdateCloudAppServiceWhere = "CreateTime,CreateUser"
const SelectAppServiceName  = "select app_name,service_name,cluster_name,entname from cloud_app_service"
const SelectServiceName = "select distinct service_name from cloud_app_service"
const SelectUserServices  = `select service_id from cloud_app_service where create_user in (?) `
const InsertCloudAppService = "insert into cloud_app_service"
const DeleteCloudAppService = "delete from cloud_app_service"
const ServiceSearchKey  = "ClusterName,Entname,ServiceName,AppName"
const SelectUsersMemory  = `select sum(memory * replicas) as memory from cloud_app_service where create_user in (?) `
const SelectUsersCpu = `select sum(cpu * replicas) as cpu from cloud_app_service where create_user in (?) `
const SelectCurrentVersion = "select service_version,image_tag,service_id from cloud_app_service"
const SelectServiceInfo = `select concat(app_name, "--",resource_name) as namespace,  concat(service_name, "--",service_version) as service_name from cloud_app_service a, cloud_ent b where (a.entname=b.entname or a.entname=b.description) and (b.description=? or b.entname=?) and a.service_name=?`

const SelectCloudContainer = "select waiting_reason, service, restart, waiting_messages,terminated_messages,terminated_reason,process,storage_data,cpu,memory,env,resource_name,create_time,create_user,service_name,cluster_name,image,status,container_id,container_name,server_address,container_ip,app_name from cloud_container"
const InsertCloudContainer = "insert into cloud_container"
const DeleteCloudContainer = "delete from cloud_container"
const SelectUserContainer = `select container_id from cloud_container where create_user in (?)`

const SelectCloudTemplateDeployHistory = "select cluster_name,domain,template_name,history_id,service_name,app_name,resource_name,create_time,create_user,entname from cloud_template_deploy_history"
const InsertCloudTemplateDeployHistory = "insert into cloud_template_deploy_history"

const UpdateServiceDomain = "update cloud_app_service set domain='%s' where app_name='%s' and cluster_name='%s' and service_name='%s' and resource_name='%s'"
