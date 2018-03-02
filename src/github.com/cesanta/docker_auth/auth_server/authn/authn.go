/*
   Copyright 2015 Cesanta Software Ltd.

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       https://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package authn

import (
	"errors"
	"strings"
	"cloud/sql"
)

type Labels map[string][]string

type AuthScope struct {
	Type    string
	Name    string
	Actions []string
}

//2018-01-27 15:08:25.7829086 +0800 CST
type CloudRegistryGroup struct {
	//镜像类型,分为共有和私有
	GroupType string
	// 组名称
	GroupName string
	// 镜像服务器域名
	ServerDomain string
	// 集群名称
	ClusterName string
}

// 验证公告账号
// 2018-01-27 16:27
func VerfiyPublicUser(scops []AuthScope, service string)  bool {
	// 查询对象和用户的权限
	services := strings.Split(service, ".")
	if len(services) < 2 {
		services = append(services, "")
	}
	for _,v := range scops {
		r := CloudRegistryGroup{}
		q := `select group_type from cloud_registry_group where group_type="公开" and cluster_name="`+sql.Replace(services[1])+`" and server_domain="`+sql.Replace(services[0])+`" and group_name="`+sql.Replace(v.Name)+`"`
		sql.GetOrm().Raw(q).QueryRow(&r)
		if r.GroupType == "公开" {
			return true
		}
	}
	return false
}
// Authentication plugin interface.
type Authenticator interface {
	// Given a user name and a password (plain text), responds with the result or an error.
	// Error should only be reported if request could not be serviced, not if it should be denied.
	// A special NoMatch error is returned if the authorizer could not reach a decision,
	// e.g. none of the rules matched.
	// Another special WrongPass error is returned if the authorizer failed to authenticate.
	// Implementations must be goroutine-safe.
	Authenticate(user string, service string, userpass string, scopse []AuthScope) (bool, Labels,string, error)

	// Finalize resources in preparation for shutdown.
	// When this call is made there are guaranteed to be no Authenticate requests in flight
	// and there will be no more calls made to this instance.
	Stop()

	// Human-readable name of the authenticator.
	Name() string
}

var NoMatch = errors.New("did not match any rule")
var WrongPass = errors.New("wrong password for user")

//go:generate go-bindata -pkg authn -modtime 1 -mode 420 -nocompress data/

type PasswordString string

func (ps PasswordString) String() string {
	if len(ps) == 0 {
		return ""
	}
	return "***"
}
