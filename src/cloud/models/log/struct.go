package log

type LogShowFilter struct {

    //
    Env string
    //
    Click int64
    //
    Id int64
    //
    Hostname string
    //
    Ip string
    //
    CreateTime string
    //
    CreateUser string
    //
    Query string
    //
    Appname string
}

//2018-09-14 09:04:13.4109859 +0800 CST
type LogShowIp struct {
    //
    Ip string
    //
    CreateTime string
    //
    Id int64
    //
    AppName string
}

type LogShowHostname struct {
    //
    Id int64
    //
    Hostname string
    //
    CreateTime string
}

type LogShowAppname struct {
    //
    Id int64
    //
    Appname string
    //
    CreateTime string
}

type LogShowHistory struct {
    //
    Ip string
    //
    CreateUser string
    //
    Query string
    //
    Appname string
    //
    Id int64
    //
    Hostname string

    //
    Env string
    //
    CreateTime string
}

// 搜索数据架构
type Search struct {
    Value string
}