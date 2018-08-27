package userperm

import (
	"cloud/controllers/users"
	"cloud/sql"
	"strings"
	"cloud/util"
	"fmt"
)


// 2018-08-24 14:47
// 获取用户组
func getUserGroups(username string) []string {
	depts := users.GetUserDept(username)
	userDepts := make([]string, 0)
	for _,v := range depts {
		userDepts = append(userDepts, `"` + v +`"`)
	}
	if len(userDepts) > 0 {
		return userDepts
	}
	return make([]string, 0)
}

const QUERY = ` and (user_name like "%?" or user_name like "%?," or user_name like ",%?," or user_name like ",%?" GROUP)`

//  2018-08-24 16:00
// 获取资源名称
func GetResourceName(tp string, user string) util.Lock {
	data := make([]CloudUserPerm, 0)
	searchMap := sql.SearchMap{}
	searchMap.Put("ResourceType", tp)
	q := sql.SearchSql(CloudUserPerm{}, SelectCloudUserPerm, searchMap)
	//q += fmt.Sprintf(` and ( user_name in (%v) or group_name in (%v))`,`"`+user+`"`, getUserGroups(user))
	q += strings.Replace(QUERY, "?", user, -1)
	groups := make([]string, 0)
	for _, v := range getUserGroups(user) {
		g := strings.Replace(QUERY, "user_name", "group_name", -1)
		g = strings.Replace(g, "?", v, -1)
		groups = append(groups, g)
	}
	q = strings.Replace(q, "GROUP", strings.Join(groups, "or"), -1)
	sql.Raw(q).QueryRows(&data)
	lock := util.Lock{}
	for _, v := range data{
		vs := strings.Split(v.Name, ",")
		for _, n := range vs {
			lock.Put(fmt.Sprintf("%v;%v;%v", n , v.ClusterName,v.Ent), "1")
			lock.Put(fmt.Sprintf("%v;%v;", n , v.ClusterName), "1")
			lock.Put(fmt.Sprintf("%v;", n), "1")
		}
	}
	return lock
}

// 2018-08-24 16:08
// 检查用户拥有权限的资源是存存在
func CheckPerm(name string,cluster string, ent string, data util.Lock) bool {
	// 都检查
	v := fmt.Sprintf("%v;%v;%v", name , cluster, ent)
	if _, ok := data.Get(v) ; ok {
		return true
	}

	// 只检查集群和环境
	if len(ent) == 0 {
		v = fmt.Sprintf("%v;%v;", name , cluster)
		if _, ok := data.Get(v) ; ok {
			return true
		}
	}
	// 只检查集群
	if len(cluster) == 0 && len(ent) == 0 {
		v = fmt.Sprintf("%v;", name)
		if _, ok := data.Get(v) ; ok {
			return true
		}
	}
	return false
}