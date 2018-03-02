package perm

const SelectCloudPermRole = "select create_time,last_modify_time,last_modify_user,description,role_id,role_name,permissions,create_user from cloud_perm_role"
const UpdateCloudPermRole = "update cloud_perm_role"
const InsertCloudPermRole = "insert into cloud_perm_role" 
const DeleteCloudPermRole = "delete from cloud_perm_role" 

const SelectCloudPerm = "select user,perm_id,create_user,description,groups,roles,create_time,last_modify_time,last_modify_user from cloud_perm"
const UpdateCloudPerm = "update cloud_perm"
const InsertCloudPerm = "insert into cloud_perm" 
const DeleteCloudPerm = "delete from cloud_perm" 

const SelectCloudApiResource = "select last_modify_time,last_modify_user,description,api_url,name,api_type,resource_id,create_time,create_user from cloud_api_resource"
const UpdateCloudApiResource = "update cloud_api_resource"
const InsertCloudApiResource = "insert into cloud_api_resource" 
const DeleteCloudApiResource = "delete from cloud_api_resource" 
