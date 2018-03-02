package groups

const SelectCloudUserGroups = "select groups_name,users,groups_id,create_time,create_user,last_modify_time,last_modify_user,description from cloud_user_groups"
const SelectCloudUserGroupsWhere = ` where 1=1 and (groups_name like "%?%" or users like "%?%")`
const SelectCloudUserDept = "select groups_name from cloud_user_groups where users like '?,%' or users like '%,?' or users like '%,?,%' or users='?'"
const SelectGroupUsers = "select users from cloud_user_groups where groups_name in (?)"

const UpdateCloudUserGroups = "update cloud_user_groups"
const UpdateCloudUserGroupsExclude = "CreateTime,CreateUser,GroupsName"
const InsertCloudUserGroups = "insert into cloud_user_groups" 
const DeleteCloudUserGroups = "delete from cloud_user_groups"
const UserGroupsLike = "select groups_id,groups_name from  cloud_user_groups where users like concat(\"NAME\",',%') or users=\"NAME\" or users like concat(\"%,\", \"NAME\") or users like concat(\"%,\",\"NAME\",\",%\")"

