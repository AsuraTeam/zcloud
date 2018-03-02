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
	"encoding/json"
	"github.com/astaxie/beego/logs"
	"cloud/controllers/index"
	"time"
	"cloud/cache"
)

type Requirements struct {
	Password *PasswordString `yaml:"password,omitempty" json:"password,omitempty"`
	Labels   Labels          `yaml:"labels,omitempty" json:"labels,omitempty"`
}

type staticUsersAuth struct {
	users map[string]*Requirements
}

func (r Requirements) String() string {
	p := r.Password
	if p != nil {
		pm := PasswordString("***")
		r.Password = &pm
	}
	b, _ := json.Marshal(r)
	r.Password = p
	return string(b)
}

func NewStaticUserAuth(users map[string]*Requirements) *staticUsersAuth {
	return &staticUsersAuth{users: users}
}

// 2018-01-19 17:09
func (sua *staticUsersAuth) Authenticate(user string,  service string, userpass string, scopes []AuthScope) (bool, Labels, string, error) {
	// 验证管理员面膜
	r := index.VerifyUser(user, userpass, service)
	if r {
		key := user + "_admin" + "_" + service
		if cache.RedisUserCache != nil {
			logs.Info("写入管理员cache", key)
			cache.RegistryPermCache.Put(key, `["pull","push"]`, time.Minute * 30)
		}
		logs.Info("管理员用户认证成功", user, make([]string,0))
		return true, make(map[string][]string, 0), "admin", nil
	}
	logs.Info("通过静态方式验证", user)
	r, _ = index.RecordLoginUser(user, userpass)
	if r {
		logs.Info("用户认证成功", user, make([]string,0))
		return true, make(map[string][]string, 0), "static", nil
	}
	logs.Info("验证公开权限")
	r = VerfiyPublicUser(scopes, service)
	if r {
		logs.Info("公开账号验证成功", user, make([]string,0))
		return true, make(map[string][]string, 0), "public", nil
	}

	return false, nil, "", nil
}

func (sua *staticUsersAuth) Stop() {
}

func (sua *staticUsersAuth) Name() string {
	return "static"
}
