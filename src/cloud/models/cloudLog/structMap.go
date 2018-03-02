package cloudLog
const SelectCloudOperLog= "select log_id,time,user,messages,cluster,ip from cloud_oper_log"
const FindByIdCloudOperLog= SelectCloudOperLog + " where log_id={1}"
const UpdateCloudOperLog= "update cloud_oper_log"
const InsertCloudOperLog= "insert into cloud_oper_log" 
const DeleteCloudOperLog= "delete from cloud_oper_log" 
