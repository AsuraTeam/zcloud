package quota

const UpdateCloudQuotaExclude = "CreateTime,CreateUser,QuotaName"
const SelectCloudQuotaWhere = ` where 1=1 and (quota_name like "%?%" or description like "%?%")`

const SelectCloudQuota = "select job_number,registry_group_number,docker_file_number,pipeline_number,user_name,quota_memory,status,app_number,lb_number,pod_number,service_number,description,quota_cpu,resource_name,last_modify_user,quota_name,create_user,last_modify_time,quota_id,create_time,group_name from cloud_quota"
const UpdateCloudQuota = "update cloud_quota"
const InsertCloudQuota = "insert into cloud_quota" 
const DeleteCloudQuota = "delete from cloud_quota"

const SelectAppQuotaUsed = `select count(*) as cnt,  resource_name from cloud_app ? group by resource_name`
const SelectAppQuotaUsedWhere = `where resource_name="?"`
