package pipeline

const SelectCloudPipeline = "select fail_action,job_id,description,pipeline_name,last_modify_user,job_name,create_user,exec_time,resource_name,pipeline_id,app_name,cluster_name,service_name,last_modify_time,create_time from cloud_pipeline"
const SelectCloudPipelineWhere = `and (pipeline_name like "%?%" or service_name like "%?%")`
const UpdateCloudPipeline = "update cloud_pipeline"
const InsertCloudPipeline = "insert into cloud_pipeline" 
const DeleteCloudPipeline = "delete from cloud_pipeline"
const SelectUserPipeline  = "select pipeline_id from cloud_pipeline where create_user in (?)"

const SelectCloudPipelineLog = "select messages,update_service_status,update_service_errormsg,build_job_errormsg,build_job_status,log_id,pipeline_name,job_name,status,update_service_end_time,app_name,cluster_name,create_user,end_job_time,push_image_start_time,update_service_start_time,service_name,resource_name,start_time,job_status,start_job_time,push_image_end_time,end_time,create_time,exec_time,job_id,run_time from cloud_pipeline_log"
const UpdateCloudPipelineLog = "update cloud_pipeline_log"
const UpdateCloudPipelineLogExclude  = "ResourceName,EndJobTime,UpdateServiceEndTime,PushImageStartTime,BuildJobErrormsg,LogId,StartJobTime,JobStatus,EndTime,UpdateServiceStatus,PushImageEndTime,UpdateServiceErrormsg,UpdateServiceStartTime,BuildJobStatus"
const InsertCloudPipelineLog = "insert into cloud_pipeline_log"
const OrderByLogIdLimt1  = " order by log_id desc limit 1"
