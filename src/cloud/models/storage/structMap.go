package storage

const SelectCloudStorageServer = "select entname,host_path,server_id,last_modify_time,description,cluster_name,storage_type,used_type,server_address,last_modify_user,create_time,create_user from cloud_storage_server"
const UpdateCloudStorageServer = "update cloud_storage_server"
const InsertCloudStorageServer = "insert into cloud_storage_server" 
const DeleteCloudStorageServer = "delete from cloud_storage_server"
const UpdateStorageServerWhere = "Entname,ClusterName,StorageType,HostPath"

const SelectCloudStorage = "select shared_type,entname,storage_server,last_modify_user,name,cluster_name,storage_type,description,storage_size,storage_format,storage_id,last_modify_time,create_time,create_user from cloud_storage"
const UpdateCloudStorage = "update cloud_storage"
const InsertCloudStorage = "insert into cloud_storage" 
const DeleteCloudStorage = "delete from cloud_storage"
const UpdateStorageExclude = "CreateTime,CreateUser,ClusterName,StorageType,StorageName"

const SelectCloudStorageMountInfo = "select status,storage_name,create_user,storage_server,service_name,app_name,cluster_name,mount_path,model,storage_type,mount_id,create_time from cloud_storage_mount_info"
const UpdateCloudStorageMountInfo = "update cloud_storage_mount_info"
const InsertCloudStorageMountInfo = "insert into cloud_storage_mount_info" 
const DeleteCloudStorageMountInfo = "delete from cloud_storage_mount_info" 
