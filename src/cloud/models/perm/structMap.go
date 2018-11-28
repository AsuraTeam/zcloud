package perm

const SelectCloudPermRole = "select create_time,is_del,last_modify_time,last_modify_user,description,role_id,role_name,permissions,create_user from cloud_perm_role"
const UpdateCloudPermRole = "update cloud_perm_role"
const InsertCloudPermRole = "insert into cloud_perm_role" 
const DeleteCloudPermRole = "delete from cloud_perm_role" 

const SelectCloudPerm = "select user,perm_id,create_user,description,groups,roles,create_time,last_modify_time,last_modify_user from cloud_perm"
const UpdateCloudPerm = "update cloud_perm"
const InsertCloudPerm = "insert into cloud_perm" 
const DeleteCloudPerm = "delete from cloud_perm" 

const SelectCloudApiResource = "select last_modify_time,method,parent,last_modify_user,description,api_url,name,api_type,resource_id,create_time,create_user from cloud_api_resource"
const UpdateCloudApiResource = "update cloud_api_resource"
const InsertCloudApiResource = "insert into cloud_api_resource" 
const DeleteCloudApiResource = "delete from cloud_api_resource"
const SelectPerm3 = `select name,api_type,parent,api_url  from cloud_api_resource where api_type in  (select  distinct parent from cloud_api_resource where parent  in (select api_type  from cloud_api_resource where api_type is not null) ) `
const SelectPerm4 = `select name,api_type,parent,api_url  from cloud_api_resource where  api_type in  (select name from cloud_api_resource where api_type  in  (select  distinct parent from cloud_api_resource where parent  in (select api_type  from cloud_api_resource where api_type is not null ) )) `
const SelectPerm5 = `select name, api_type, parent ,api_url from cloud_api_resource where api_type in  (select name  from cloud_api_resource where   api_type in  (select name from cloud_api_resource where api_type  in  (select  distinct parent from cloud_api_resource where parent  in (select api_type  from cloud_api_resource where api_type is not null) )))`
const SelectCloudUserPerm = "select create_time,cluster_name,perm_id,group_name,sub_user,create_user,resource_type,last_modify_time,description,parent_user,name,last_modify_user,ent,user_name,resource_name from cloud_user_perm"
const UpdateCloudUserPerm = "update cloud_user_perm"
const InsertCloudUserPerm = "insert into cloud_user_perm" 
const DeleteCloudUserPerm = "delete from cloud_user_perm" 

const SelectCloudPermRolePerm = "select role_id,perm_name,create_user,create_time from cloud_perm_role_perm"
const UpdateCloudPermRolePerm = "update cloud_perm_role_perm"
const InsertCloudPermRolePerm = "insert into cloud_perm_role_perm" 
const DeleteCloudPermRolePerm = "delete from cloud_perm_role_perm" 

const SelectCloudPermRoleUser = "select role_id,user_name,group_name,create_user,create_time from cloud_perm_role_user"
const UpdateCloudPermRoleUser = "update cloud_perm_role_user"
const InsertCloudPermRoleUser = "insert into cloud_perm_role_user" 
const DeleteCloudPermRoleUser = "delete from cloud_perm_role_user" 


