package hosts
const SelectCloudClusterHosts= "select image_num,container_num,cpu_free,host_ip,host_label,status,last_modify_user,mem_free,pod_num,mem_size,create_user,host_type,cpu_percent,create_time,last_modify_time,is_valid,mem_percent,host_id,create_method,cpu_num,cluster_name,api_port from cloud_cluster_hosts"
const FindByIdCloudClusterHosts= SelectCloudClusterHosts + " where host_id={1}"
const UpdateCloudClusterHosts= "update cloud_cluster_hosts"
const InsertCloudClusterHosts= "insert into cloud_cluster_hosts" 
const DeleteCloudClusterHosts= "delete from cloud_cluster_hosts"
const UpdateExclude = "HostIp,CreateTime,ClusterName"