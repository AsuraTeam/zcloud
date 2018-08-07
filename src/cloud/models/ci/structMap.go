package ci

const SelectCloudCodeRepostitory = "select create_time,create_user,description,code_url,code_source,username,gitlab_token,type,repostitory_id,last_modify_time,last_modify_user,password from cloud_code_repostitory"
const UpdateCloudCodeRepostitory = "update cloud_code_repostitory"
const InsertCloudCodeRepostitory = "insert into cloud_code_repostitory"
const DeleteCloudCodeRepostitory = "delete from cloud_code_repostitory"
const SelectCloudCodeRepostitoryWhere  = ` where 1=1 and (code_url like "%?%" or description like "%?%" or create_user like "%?%")`

const SelectCloudCiDockerfile = "select last_modify_user,description,is_del,file_id,content,script,create_time,create_user,name,last_modify_time from cloud_ci_dockerfile"
const UpdateCloudCiDockerfile = "update cloud_ci_dockerfile"
const UpdateDockerfileExclude  = "CreateTime,CreateDockerFile"
const InsertCloudCiDockerfile = "insert into cloud_ci_dockerfile" 
const DeleteCloudCiDockerfile = "delete from cloud_ci_dockerfile"
const SelectDockerfiles  = "select file_id from cloud_ci_dockerfile where create_user in (?)"
const SelectDockerfileWhere = ` and (name like "%?%" or description like "%?" or content like "%?%" or script like "%?%" )`




const SelectCloudBuildJob = "select base_image,description,last_tag,content,time_out,last_modify_user,script,job_code,docker_file,create_user,build_status,job_name,create_time,build_id,last_build_time,job_id,registry_server,last_modify_time,item_name,image_tag,cluster_name from cloud_build_job"
const UpdateCloudBuildJob = "update cloud_build_job"
const UpdateCloudBuildJobExclude2 = "CreateTime,CreateUser,RegistryServer,ClusterName"
const InsertCloudBuildJob = "insert into cloud_build_job" 
const DeleteCloudBuildJob = "delete from cloud_build_job"
const UpdateCloudBuildJobExclude  = "ImageTag,CreateTime,CreateUser,DockerFile,Content,TimeOut,BaseImage,Script"
const SelectCloudBuildJobWhere  =  ` and (item_name like "%?%" or description like "%?%" or create_user like "%?%")`
const SelectUserJobs  = "select job_id from cloud_build_job where create_user in (?)"

const SelectCloudBuildJobHistoryDockerfile  = "select docker_file from cloud_build_job_history"
const SelectCloudBuildJobHistory = "select base_image,registry_server,history_id,registry_group,job_name,create_time,build_time,script,item_name,image_tag,cluster_name,create_user,build_status,build_logs,job_id,docker_file from cloud_build_job_history"
const UpdateCloudBuildJobHistory = "update cloud_build_job_history"
const InsertCloudBuildJobHistory = "insert into cloud_build_job_history"
const SelectBuildJobToApp  = "select history_id, registry_server, registry_group, item_name,image_tag from cloud_build_job_history"
const ExcludeUpdateHistoryColumn = "JobId,JobName,CreateTime,CreateUser,ImageTag,ClusterName,ImageTag,DockerFile,RegistryServer,ItemName,RegistryGroup,Script"
const SelectBuildHistoryWhere =` and (job_name like "%?%" or build_logs like "%?%" or item_name like "%?%" or create_user like "%?%")`
const SelectJobTimeout = "select  history_id,a.create_time, b.time_out as build_time from cloud_build_job_history a, cloud_build_job b where a.build_status='构建中' and a.job_id=b.job_id"
const UpdateCloudBuildJobTimeout = "update cloud_build_job_history set build_status='构建超时' where history_id="


const SelectCloudBuildJobHistoryLast = "select history_id, job_id,build_logs,build_status, job_name,create_time from cloud_build_job_history where job_id={0} order by history_id desc limit 1"


const SelectCloudCiService = "select service_name,new_version,percent,lb_version,current_version,entname,group_name,create_time,create_user,release_type,service_id,cluster_name,status,domain,last_modify_user,app_name,description,image_name,last_modify_time from cloud_ci_service"
const UpdateCloudCiService = "update cloud_ci_service"
const InsertCloudCiService = "insert into cloud_ci_service" 
const DeleteCloudCiService = "delete from cloud_ci_service"
const SelectCloudCiServiceWhere  = `where 1=1  and (service_name like "%?%" or description like "%?%" or domain like "%?%" or app_name like "%?%")`
const UpdateCiServicePercent = UpdateCloudCiService + ` set percent={0} where domain="{1}"`


const SelectCloudCiReleaseHistory = "select release_type,old_images,create_user,domain,app_name,release_bug_description,release_demand_description,entname,description,release_item_description,release_job_pm_link,release_test_user,release_bug_pm_link,image_name,status,history_id,create_time,service_name,cluster_name,release_online_type from cloud_ci_release_history"
const UpdateCloudCiReleaseHistory = "update cloud_ci_release_history"
const InsertCloudCiReleaseHistory = "insert into cloud_ci_release_history" 
const DeleteCloudCiReleaseHistory = "delete from cloud_ci_release_history"
const SelectReleaseHistoryWhere  = `  and (service_name like "%?%" or description like "%?%" or domain like "%?%" or app_name like "%?%" or release_demand_description like "%?%" or release_item_description like "%?%")`
const UpdateHistoryExclude  = `release_type,old_images,domain,service_name,app_name,cluster_name,entname,image_name,create_time,status,release_type,status,release_version,action,service_id,lb_version`
const SelectCloudCiReleaseLog = "select log_id,ip,domain,app_name,entname,messages,service_name,cluster_name,action,create_time,create_user from cloud_ci_release_log"
const InsertCloudCiReleaseLog = "insert into cloud_ci_release_log"
const DeleteCloudCiReleaseLog = "delete from cloud_ci_release_log"
const SelectCiReleaseLogWhere = ` and (service_name like "%?%" or domain like "%?%" or app_name like "%?%"  or messages like "%?%" )`

const SelectCloudCiPerm = "select perm_id,username,groups_name,create_user,datas,create_time,last_modify_time,last_modify_user from cloud_ci_perm"
const SelectCloudCiPermWhere = ` where dates like "%?%"`
const UpdateCloudCiPerm = "update cloud_ci_perm"
const InsertCloudCiPerm = "insert into cloud_ci_perm" 
const DeleteCloudCiPerm = "delete from cloud_ci_perm" 
