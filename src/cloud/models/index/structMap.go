package index

const SelectDockerCloudCiEnv = "select create_user,service_name,auto_compile," +
	"build_status,env_id,last_modify_time,code_branch,images_tp,cpu,memory,description,item_tp" +
	",groups_id,auto_compile_tp,auto_compile_status,images," +
	"approve_time,gson_data,code_path,groups_ids,last_modify_user,create_time,entname,approve,dept," +
	"domain,groups_name,auto_compile_time,container_max_number,build_script,container_number,health_script," +
	"approve_user from docker_cloud_ci_env"
const FindByIdDockerCloudCiEnv = SelectDockerCloudCiEnv + " where env_id={1}"
const UpdateDockerCloudCiEnv = "update docker_cloud_ci_env"
const InsertDockerCloudCiEnv = "insert into docker_cloud_ci_env"
const DeleteDockerCloudCiEnv = "delete from docker_cloud_ci_env"

const SelectDockerCloudAuthorityUser= "select create_time,token,create_user,last_modify_time,last_modify_user,dept,description,real_name,is_del,user_name,third_id,third_true_name,user_email,is_valid,user_id,user_pic,user_mobile,pwd from cloud_authority_user"
const SelectUserWhere = ` where 1=1 and (user_name like "%?%" or description like "%?%" or user_email like "%?%")`

const UpdateDockerCloudAuthorityUser= "update cloud_authority_user"
const InsertDockerCloudAuthorityUser= "insert into cloud_authority_user" 
const DeleteDockerCloudAuthorityUser= "delete from cloud_authority_user"

const SelectCloudLoginRecord= "select record_id,login_time,login_ip,login_user,login_Status from cloud_login_record"
const FindByIdCloudLoginRecord= SelectCloudLoginRecord + " where record_id={1}"
const UpdateCloudLoginRecord= "update cloud_login_record"
const InsertCloudLoginRecord= "insert into cloud_login_record" 
const DeleteCloudLoginRecord= "delete from cloud_login_record" 

