package permdata

import (
	"cloud/sql"
	"github.com/astaxie/beego/logs"
	"cloud/userperm"
	"cloud/util"
	"fmt"
	"cloud/models/perm"
	"strings"
)

const (
	userKey     = "-user-服务"
	groupKey    = "-group-服务"
	userAppKey  = "-user-应用"
	groupAppKey = "-group-应用"
	service = "服务"
	app = "应用"
)

// 缓存用户拥有的权限数据
func getUserPerm(user string, lock util.Lock) {
	serviceKey := fmt.Sprintf("%v-service", user)
	if _, ok := lock.Get(serviceKey); ok {
		return
	}
	perm := userperm.GetResourceName(service, user)
	permApp := userperm.GetResourceName(app, user)
	lock.Put(serviceKey, perm)
	lock.Put(fmt.Sprintf("%v-app", user), permApp)
	lock.Put(fmt.Sprintf("%v-group", user), userperm.GetUserGroups(user))
}

// 获取分配资源权限的用户和组资源数据
func getUserResourcePerm(tp string, lock util.Lock) {
	data := make([]perm.CloudUserPerm, 0)
	sql.GetOrm().Raw(SelectUserResourcePerm, tp).QueryRows(&data)
	for _, v := range data {
		key := strings.Split(fmt.Sprintf("%v-user-,%v-group-", v.UserName, v.GroupName), ",")
		for _, k := range key {
			k += tp + "," + v.Ent + "," + v.ClusterName
			obj, ok := lock.Get(k)
			names := strings.Split(v.Name, ",")
			if ok {
				names = append(names, obj.([]string)...)
			}
			lock.Put(k, names)
		}
	}
}

func getDefaultString(str string) string {
	if str == "" {
		return "-"
	}
	return str
}

func writePermDetail(user, name, cluster, ent, tp, group string) {
	user = getDefaultString(user)
	group = getDefaultString(group)
	name = getDefaultString(name)
	i := CloudUserPermDetail{
		DetailId:     0,
		Username:     user,
		Name:         name,
		ClusterName:  cluster,
		Ent:          ent,
		ResourceType: tp,
		GroupName:    group,
	}
	q := sql.InsertSql(i, InsertCloudUserPermDetail)
	sql.Exec(q)
}

// 写入用户自己建立的数据
func writeCreateUserContainer(v CloudContainer, lock util.Lock) {
	k := fmt.Sprintf("%v-group", v.CreateUser)
	groupsData, ok := lock.Get(k)
	if ok {
		groups := groupsData.([]string)
		for _, g := range groups {
			writePermDetail(v.CreateUser, v.ContainerName, v.ClusterName, v.Entname, "container", g)
		}
	} else {
		writePermDetail(v.CreateUser, v.ContainerName, v.ClusterName, v.Entname, "container", "")
	}
}

// 写入配置的权限数据
func writeConfigurePermContainer(d interface{}, v CloudContainer, meta []string, userKey string, isUser bool, isGroup bool, name string) {
	names := d.([]string)
	for _, n := range names {
		if name == n && v.Entname == meta[1] && v.ClusterName == meta[2] {
			u := strings.Split(meta[0], userKey)[0]
			if isUser {
				writePermDetail(u, v.ContainerName, v.ClusterName, v.Entname, "container", "")
			}
			if isGroup {
				writePermDetail("", v.ContainerName, v.ClusterName, v.Entname, "container", u)
			}
		}
	}
}

/**
生成用户容器权限数据
 */
func MakeContainerData() {
	sql.Exec(DeleteExpireContainerPerm)
	lock := util.Lock{}
	data := make([]CloudContainer, 0)
	sql.Raw(SelectContainerData).QueryRows(&data)
	logs.Info(data)
	for _, v := range data {
		getUserPerm(v.CreateUser, lock)
	}
	getUserResourcePerm(service, lock)
	getUserResourcePerm(app, lock)

	for _, v := range data {
		writeCreateUserContainer(v, lock)
		for k, d := range lock.GetData() {
			meta := strings.Split(k, ",")
			isUser := strings.Contains(k, userKey)
			isGroup := strings.Contains(k, groupKey)
			if isUser || isGroup {
				writeConfigurePermContainer(d, v, meta, userKey, isUser, isGroup, v.ServiceName)
			}
			isAppUser := strings.Contains(k, userAppKey)
			isAppGroup := strings.Contains(k, groupAppKey)
			if isAppGroup || isAppUser {
				writeConfigurePermContainer(d, v, meta, userKey, isUser, isGroup, v.AppName)
			}
		}
	}
}
