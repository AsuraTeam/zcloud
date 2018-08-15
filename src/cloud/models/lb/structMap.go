package lb
const SelectCloudLb = "select service_number,entname,cpu,memory,host_log_path,lb_ip,lb_type,lb_id,description,create_time,cluster_name,resource_name,last_modify_time,lb_name,lb_domain_prefix,lb_domain_suffix,create_user,last_modify_user,status from cloud_lb"
const UpdateCloudLb = "update cloud_lb"
const UpdateLbExclude  = "CreateTime,CreateUser,LbName"
const InsertCloudLb = "insert into cloud_lb" 
const DeleteCloudLb = "delete from cloud_lb"
const SelectUserLbs = "select lb_id from cloud_lb where create_user in (?)"

const SelectCloudLbWhere = 	`where 1=1 and (lb_name like "%?%" or description like "%?%")`
const SelectCloudLbService = "select app_name,flow_service_name,percent,lb_type,entname,service_version,protocol,lb_method,lb_id,default_domain,domain,lb_service_id,cluster_name,last_modify_user,service_id,last_modify_time,create_time,create_user,service_name,lb_name,cert_file,description,listen_port,container_port from cloud_lb_service"
const UpdateCloudLbService = "update cloud_lb_service"
const UpdateLbServiceExclude  = "CreateTime,CreateUser,LbName"
const InsertCloudLbService = "insert into cloud_lb_service" 
const DeleteCloudLbService = "delete from cloud_lb_service"
const SelectCloudLbServiceWhere  = ` and (lb_name like "%?%" or service_name like "%?%")`
const SelectLbDomainData  =  SelectCloudLbService + ` where domain = "{0}"`
const SelectLbDomain  = "select domain, service_version from cloud_lb_service "
const SelectLbServiceVersion = SelectLbDomain + " where domain in (?)"
const UpdateLbServiceServiceVersion =  UpdateCloudLbService + ` set service_version={0} where domain="{1}"`
const UpdateLbServicePercent = UpdateCloudLbService + ` set percent={0},flow_service_name="{2}" where domain="{1}"`

const UpdateCloudLbNginxConf = "update cloud_lb_nginx_conf"

const DeleteCloudLbNginxConf = "delete from cloud_lb_nginx_conf" 


const UpdateCloudLbCert = "update cloud_lb_cert"
const InsertCloudLbCert = "insert into cloud_lb_cert" 
const DeleteCloudLbCert = "delete from cloud_lb_cert" 
