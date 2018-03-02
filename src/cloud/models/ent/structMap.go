package ent

const SelectCloudEnt = "select create_time,clusters,entname,ent_id,create_user,last_modify_user,last_modify_time,description from cloud_ent"
const FindByIdCloudEnt = SelectCloudEnt + " where ent_id={1}"
const UpdateCloudEnt = "update cloud_ent"
const InsertCloudEnt = "insert into cloud_ent" 
const DeleteCloudEnt = "delete from cloud_ent"
const SelectCloudEntWhere  =  `where 1=1 and (ent_name like "%?%" or description like "%?%")`
