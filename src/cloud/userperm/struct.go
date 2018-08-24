package userperm

//2018-08-23 08:38:08.9273685 +0800 CST
type CloudUserPerm struct {
	//
	ResourceType string
	//
	LastModifyTime string
	//
	LastModifyUser string
	//
	Ent string
	//
	Description string
	//
	ParentUser string
	//
	Name string
	//
	UserName string
	//
	ResourceName string
	//
	SubUser string
	//
	CreateUser string
	//
	CreateTime string
	//
	ClusterName string
	//
	PermId int64
	//
	GroupName string
}

const SelectCloudUserPerm = "select create_time,cluster_name,perm_id,group_name,sub_user,create_user,resource_type,last_modify_time,description,parent_user,name,last_modify_user,ent,user_name,resource_name from cloud_user_perm"

