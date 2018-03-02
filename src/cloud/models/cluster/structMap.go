package cluster
const SelectCloudCluster= "select cert_data,ca_data,key_data,network_cart,cluster_id,cluster_alias,last_modify_time,docker_version,docker_install_dir,cluster_type,cluster_name,create_time,create_user from cloud_cluster"
const SelectCloudClusterWhere = `  where 1=1 and (cluster_name like "?" or cluster_alias like "?")`
const UpdateCloudCluster= "update cloud_cluster"
const InsertCloudCluster= "insert into cloud_cluster" 
const DeleteCloudCluster= "delete from cloud_cluster" 
