package monitor

const SelectCloudAutoScale = "select description,entname,metric_name,data_source, es, metric_type,scale_id,create_time,lt,gt,increase_step,create_user,app_name,resource_name,service_name,last_modify_time,last_modify_user,end,reduce_step,namespace,service_version,action_interval,cluster_name,last_count,step,start,query,msg_group from cloud_auto_scale"
const UpdateCloudAutoScale = "update cloud_auto_scale"
const UpdateAutoScaleExclude = "CreateTime,CreateUser"
const InsertCloudAutoScale = "insert into cloud_auto_scale" 
const DeleteCloudAutoScale = "delete from cloud_auto_scale"
const SelectCloudAutoScaleWhere = ` where app_name like "%?%" or service_name like "%?%"  or create_user like "%?%" or query like "%?%"`

const SelectCloudAutoScaleLog = "select action_interval,gt,step,service_name,entname,increase_step,replicas_max,cluster_name,last_count,query,app_name,metric_type,reduce_step,replicas,log_id,create_time,status,monitor_value,replicas_min,metric_name,es from cloud_auto_scale_log"
const UpdateCloudAutoScaleLog = "update cloud_auto_scale_log"
const SelectAutoScaleLogWhere = ` where app_name like "%?%" or service_name like "%?%"  or create_user like "%?%" or query like "%?%"`
const DeleteCloudAutoScaleLog = "delete from cloud_auto_scale_log" 
