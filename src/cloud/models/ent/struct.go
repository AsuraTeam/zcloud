package ent

//2018-02-05 18:01:19.0774817 +0800 CST
type CloudEnt struct {
    //创建时间
    CreateTime string
    //环境名称
    Entname string
    //
    EntId int64
    //创建人
    CreateUser string
    //最近修改人
    LastModifyUser string
    //最近修改时间
    LastModifyTime string
    //备注信息
    Description string
    // 拥有集群
    Clusters string
}
