package registry

const SelectCloudRegistryServer = "select entname,access,auth_server, name,admin,groups_id,last_modify_time,groups,server_address,cluster_name,create_time,create_user,images_number,username,server_id,last_modify_user,prefix,password,server_domain,description from cloud_registry_server"
const UpdateCloudRegistryServer = "update cloud_registry_server"
const UpdateRegistryServerExclude = "AuthServer,Name,CreateTime,CreateUser"
const InsertCloudRegistryServer = "insert into cloud_registry_server" 
const DeleteCloudRegistryServer = "delete from cloud_registry_server"
const SelectRegistryServerWhere =` where 1=1 and (server_address like "%?%" or description like "%?%" or server_domain like "%?%")`
const UpdateRegistryServerExcludePass  = "ServerId,Password,ClusterName,ServerDomain"
const SelectRegistryAccess  = "集群外&nbsp;<br><a target='_blank' href='https://?/v2/'>?</a>"

const SelectCloudRegistryPermissions = "select cluster_name,description,action,create_time,last_modify_user,service_name,user_name,project,image_name,create_user,last_modify_time,permissions_id,user_type,registry_server,groups_name from cloud_registry_permissions"
const UpdateCloudRegistryPermissions = "update cloud_registry_permissions"
const InsertCloudRegistryPermissions = "insert into cloud_registry_permissions" 
const DeleteCloudRegistryPermissions = "delete from cloud_registry_permissions"
const UpdateCloudRedisPermExclude  = "AuthServer,Name,CreateTime,CreateUser"
const SelectCloudRegistryPermWhere  = ` and (service_name like "?" or user_name like "?" or project like "?")`

const SelectCloudRegistryGroup = "select tag_number,cluster_name,server_domain,group_name,last_modify_time,image_number,create_time,create_user,group_type,last_modify_user,size_totle,group_id from cloud_registry_group"
const SelectCloudRegistryGroupWhere  = `and (group_name like "%?%"  or create_user like "?")`
const UpdateCloudRegistryGroup = "update cloud_registry_group"
const UpdateGroupExclude = "GroupName,CreateTime"
const UpdateCloudRegistryGroupExclude  =  "GroupName,CreateUser,CreateTime,LastModifyUser,ServerDomain,ClusterName"
const InsertCloudRegistryGroup = "insert into cloud_registry_group" 
const DeleteCloudRegistryGroup = "delete from cloud_registry_group"
const SelectRegistryServerGroup  = `select b.server_address as server_address,b.server_domain as server_domain, b.auth_server as auth_server from cloud_registry_group a , cloud_registry_server b where a.server_domain=b.server_domain and a.cluster_name=b.cluster_name and a.group_name="{0}" and a.cluster_name="{1}"`
const SelectUserRegistryGroups = "select group_id from cloud_registry_group where create_user in (?)"

const SelectCloudImage = "select access,download,tags,tag_number,layers_number,create_user,name,repositories,image_type,image_id,create_time,size,repositories_group from cloud_image"
const SelectCloudImageWhere = ` and (name like "%?%"  or create_time like "%?%")`
const UpdateCloudImage = "update cloud_image"
const UpdateCloudImageExclude  = "CreateTime,CreateUser,ImageId,Repositories"
const InsertCloudImage = "insert into cloud_image" 
const DeleteCloudImage = "delete from cloud_image"
const SelectCloudImageExists  = "select name,repositories_group  from cloud_image"
const SelectImageTgs  = "select image_id,tags,access,name from cloud_image"
const SelectDeployImage  = `select a.cluster_name as cluster_name, c.tags,c.name as name, a.server_domain as server_domain,a.server_address  as server_address from
 cloud_registry_server a ,
  cloud_registry_group b ,
  cloud_image c
where a.server_domain = b.server_domain and
      a.cluster_name=b.cluster_name and
      c.repositories_group=b.group_name and
      (b.create_user=? or b.group_type="公开")
      and (c.name like "%{0}%") order by b.create_time`

const SelectImageDownload = "select count(*) as log_id from cloud_image_log"
const SelectCloudImageLog = "select ip,cluster_name,oper_type,log_id,create_user,repositories,repositories_group,label,create_time,name from cloud_image_log"
const InsertCloudImageLog = "insert into cloud_image_log"
const DeleteCloudImageLog = "delete from cloud_image_log"
const SelectImageLogWhere  = ` and (name like "%?%" or create_user like "%?%")`

const SelectCloudImageSyncLog = "select status,version,registry_server_1,create_time,runtime,create_user,registry_group,registry_server_2,item_name,log_id,messages from cloud_image_sync_log"
const UpdateCloudImageSyncLog = "update cloud_image_sync_log"
const InsertCloudImageSyncLog = "insert into cloud_image_sync_log" 
const DeleteCloudImageSyncLog = "delete from cloud_image_sync_log"
const SelectImageSyncLogWhere = ` and (version like "%?%" or item_name like "%?%")`

const SelectCloudImageSync = "select last_modify_user,registry,status,approved_by,sync_id,cluster_name,target_cluster,description,registry_group,image_name,entname,create_time,create_user,target_registry,version,target_entname,last_modify_time,approved_time from cloud_image_sync"
const SelectCloudImageSyncWhere = ` and (description like "?" or  version like "%?%" or item_name like "%?%")`
const UpdateCloudImageSync = "update cloud_image_sync"
const InsertCloudImageSync = "insert into cloud_image_sync" 
const DeleteCloudImageSync = "delete from cloud_image_sync"

const SelectCloudImageBase = "select image_type,image_name,last_modify_time,last_modify_user,icon,description,base_id,registry_server,create_time,create_user from cloud_image_base"
const SelectCloudBaseWhere  = ` where 1=1 and (image_name like "%?%" or description like "%?%" or registry_server like "%?%")`
const UpdateCloudImageBase = "update cloud_image_base"
const InsertCloudImageBase = "insert into cloud_image_base" 
const DeleteCloudImageBase = "delete from cloud_image_base" 
