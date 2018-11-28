package permdata
const SelectCloudUserPermDetail = "select name,cluster_name,ent,detail_id,username,resource_type,group_name from cloud_user_perm_detail"
const UpdateCloudUserPermDetail = "update cloud_user_perm_detail"
const InsertCloudUserPermDetail = "insert into cloud_user_perm_detail"
const DeleteCloudUserPermDetail = "delete from cloud_user_perm_detail"

const SelectContainerData  = `select container_name,app_name,service_name,create_user,cluster_name,resource_name,entname from cloud_container`
const SelectUserResourcePerm = `select user_name, group_name, name,cluster_name,ent from cloud_user_perm where resource_type=?`

const DeleteExpireContainerPerm = `delete from cloud_user_perm_detail where name in (select name from (select name from cloud_user_perm_detail  where resource_type='container' and name not in (select container_name from cloud_container)) as temp) `
