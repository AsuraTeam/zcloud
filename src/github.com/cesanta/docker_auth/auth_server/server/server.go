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

package server

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"regexp"
	"sort"
	"strings"
	"time"

	"github.com/cesanta/docker_auth/auth_server/authn"
	"github.com/cesanta/docker_auth/auth_server/authz"
	"github.com/cesanta/glog"
	"github.com/docker/distribution/registry/auth/token"
	"github.com/astaxie/beego/logs"
	"cloud/models/registry"
	"cloud/sql"
	"cloud/util"
	"github.com/docker/distribution/registry/api/errcode"
	"cloud/models/groups"
	_ "github.com/astaxie/beego/cache/redis"
	"strconv"
	"github.com/garyburd/redigo/redis"
	"cloud/cache"
)

var (
	hostPortRegex = regexp.MustCompile(`\[?(.+?)\]?:\d+$`)
)

type AuthServer struct {
	config         *Config
	authenticators []authn.Authenticator
	authorizers    []authz.Authorizer
}

func NewAuthServer(c *Config) (*AuthServer, error) {
	as := &AuthServer{
		config:      c,
		authorizers: []authz.Authorizer{},
	}
	u := map[string]*authn.Requirements{}
	c.Users = u
	as.authenticators = append(as.authenticators, authn.NewStaticUserAuth(c.Users))
	logs.Info("c.users", c.Users)
	if c.Users != nil {
		as.authenticators = append(as.authenticators, authn.NewStaticUserAuth(c.Users))
		logs.Info("as.authenticators", as.authenticators)
	}

	return as, nil
}

type authRequest struct {
	RemoteConnAddr string
	RemoteAddr     string
	RemoteIP       net.IP
	User           string
	Password       authn.PasswordString
	Account        string
	Service        string
	Scopes         []authScope
	Labels         authn.Labels
	UserPass       string
	AuthType       string
}

type authScope struct {
	Type    string
	Name    string
	Actions []string
}

type authzResult struct {
	scope            authScope
	autorizedActions []string
}

func (ar authRequest) String() string {
	return fmt.Sprintf("{%s:%s@%s %s}", ar.User, ar.Password, ar.RemoteAddr, ar.Scopes)
}

func parseRemoteAddr(ra string) net.IP {
	hp := hostPortRegex.FindStringSubmatch(ra)
	if hp != nil {
		ra = string(hp[1])
	}
	res := net.ParseIP(ra)
	return res
}

func (as *AuthServer) ParseRequest(req *http.Request) (*authRequest, error) {

	ar := &authRequest{RemoteConnAddr: req.RemoteAddr, RemoteAddr: req.RemoteAddr}
	if as.config.Server.RealIPHeader != "" {
		hv := req.Header.Get(as.config.Server.RealIPHeader)
		ips := strings.Split(hv, ",")

		realIPPos := as.config.Server.RealIPPos
		if realIPPos < 0 {
			realIPPos = len(ips) + realIPPos
			if realIPPos < 0 {
				realIPPos = 0
			}
		}

		ar.RemoteAddr = strings.TrimSpace(ips[realIPPos])
		glog.V(3).Infof("Conn ip %s, %s: %s, addr: %s", ar.RemoteAddr, as.config.Server.RealIPHeader, hv, ar.RemoteAddr)
		if ar.RemoteAddr == "" {
			return nil, fmt.Errorf("client address not provided")
		}
	}
	ar.RemoteIP = parseRemoteAddr(ar.RemoteAddr)
	if ar.RemoteIP == nil {
		return nil, fmt.Errorf("unable to parse remote addr %s", ar.RemoteAddr)
	}
	user, password, haveBasicAuth := req.BasicAuth()
	if haveBasicAuth {
		ar.User = user
		ar.Password = authn.PasswordString(password)
		ar.UserPass = password
	}
	ar.Account = req.FormValue("account")
	if ar.Account == "" {
		ar.Account = ar.User
	} else if haveBasicAuth && ar.Account != ar.User {
		return nil, fmt.Errorf("user and account are not the same (%q vs %q)", ar.User, ar.Account)
	}
	ar.Service = req.FormValue("service")
	if err := req.ParseForm(); err != nil {
		return nil, fmt.Errorf("invalid form value")
	}
	// https://github.com/docker/distribution/blob/1b9ab303a477ded9bdd3fc97e9119fa8f9e58fca/docs/spec/auth/scope.md#resource-scope-grammar
	if req.FormValue("scope") != "" {
		for _, scopeStr := range req.Form["scope"] {
			parts := strings.Split(scopeStr, ":")
			var scope authScope
			switch len(parts) {
			case 3:
				scope = authScope{
					Type:    parts[0],
					Name:    parts[1],
					Actions: strings.Split(parts[2], ","),
				}
			case 4:
				scope = authScope{
					Type:    parts[0],
					Name:    parts[1] + ":" + parts[2],
					Actions: strings.Split(parts[3], ","),
				}
			default:
				return nil, fmt.Errorf("invalid scope: %q", scopeStr)
			}
			sort.Strings(scope.Actions)
			ar.Scopes = append(ar.Scopes, scope)
		}
	}
	return ar, nil
}

// 获取操作权限
// 2018-01-27 16:39
func getScops(ar *authRequest) []authn.AuthScope {
	scops := make([]authn.AuthScope, 0)
	for _, v := range ar.Scopes {
		t := authn.AuthScope{}
		t.Name = v.Name
		t.Actions = v.Actions
		t.Type = v.Type
		scops = append(scops, t)
	}
	return scops
}

func (as *AuthServer) Authenticate(ar *authRequest) (bool, authn.Labels, string, error) {
	logs.Info("Authenticate as.authenticators", len(as.authenticators))
	for i, a := range as.authenticators {
		logs.Info("Authenticate", i)
		scops := getScops(ar)
		result, labels, authtype, err := a.Authenticate(ar.Account, ar.Service, ar.UserPass, scops)
		logs.Info("Authenticate result", result, labels, err)
		//logs.Info("Authn %s %s -> %t, %+v, %v", a.Name(), ar.Account, result, labels, err)
		if err != nil {
			if err == authn.NoMatch {
				logs.Info("authn.NoMatch", err)
				continue
			} else if err == authn.WrongPass {
				logs.Warn("Failed authentication with %s: %s", err, ar.Account)
				return false, nil, "", nil
			}
			err = fmt.Errorf("authn #%d returned error: %s", i+1, err)
			logs.Error("%s: %s", ar, err)
			return false, nil, "", err
		}
		logs.Info("Auth result", result, labels, nil)
		return result, labels, authtype, nil
	}
	// Deny by default.
	logs.Info("%s did not match any authn rule", ar)
	return false, nil, "", nil
}

// 替换镜像名称,用于查询
// 2018-01-20 06:11
func replaceImage(image string) string {
	return strings.Replace(PROJECT_LIKE, "NAME", sql.Replace(image), -1)
}

// 镜像名称拼接
// 2018-01-20 6:14
func joinImage(name ...string) string {
	n := make([]string, 0)
	for _, v := range name {
		n = append(n, v)
	}
	return strings.Join(n, "/")
}

const PROJECT_LIKE = ` or project like "NAME" `

// 获取模糊匹配权限的sql
// 写死匹配最多3级目录
// 2018-01-20 06:01
func GetLikeProjectSql(ai *authz.AuthRequestInfo) string {
	images := strings.Split(ai.Name, "/")
	q := ""
	for id := range images {
		q += replaceImage(joinImage(images[0:id+1]...))
	}
	return " and (" + strings.Replace(q, "or", "", 1) + " or project='*') "
}

// 获取用户组的权限
// 2018-01-20 07:20
func getUserGroups(user string) string {
	g := make([]groups.CloudUserGroups, 0)
	q := strings.Replace(groups.UserGroupsLike, "NAME", sql.Replace(user), -1)
	sql.GetOrm().Raw(q).QueryRows(&g)
	if len(g) > 0 {
		gl := make([]string, 0)
		for _, v := range g {
			gl = append(gl, strconv.FormatInt(v.GroupsId, 10))
		}
		r := " or groups_name in (" + strings.Join(gl, ",") + " ) "
		return r
	}
	return ""
}

// 查询用户权限
// 2018-01-22 17:13
func getUsernameSql(username string) string {
	r := ` and ((user_name like "` + sql.Replace(username) + `,%" or user_name like "%,` + sql.Replace(username) + `,%" or user_name like "%,` + sql.Replace(username) + `%" or user_name="` + sql.Replace(username) + `") or (user_name is null and groups_name is null) GROUPS )`
	return r
}

// 获取用户权限
// 2018-01-20 8:08
func getPermissions(ai *authz.AuthRequestInfo) []registry.CloudRegistryPermissions {
	permission := make([]registry.CloudRegistryPermissions, 0)
	likesql := GetLikeProjectSql(ai)
	// 查询对象和用户的权限
	services := strings.Split(ai.Service, ".")
	if len(services) < 2 {
		services = append(services, "")
	}
	searchMap := sql.GetSearchMapV("ServiceName", services[0], "ClusterName", services[1])
	q := sql.SearchSql(registry.CloudRegistryPermissions{}, registry.SelectCloudRegistryPermissions, searchMap)
	q += likesql

	var cachePermission = cache.RegistryPermCache.Get(ai.Account)
	if cachePermission == nil {
		cachePermission = nil
		groupsP := getUserGroups(ai.Account)

		if groupsP != "" {
			q += strings.Replace(getUsernameSql(ai.Account), "GROUPS", groupsP, -1)
		} else {
			q += strings.Replace(getUsernameSql(ai.Account), "GROUPS", "", -1)
		}
	}
	cacheKey := ai.Account + ai.Service
	if cachePermission != nil {
		cachePermission = cache.RegistryPermCache.Get(cacheKey)
		r, _ := redis.String(cachePermission, nil)
		json.Unmarshal([]byte(r), &permission)
	} else {
		sql.GetOrm().Raw(q).QueryRows(&permission)
		if len(permission) > 0 {
			cache.RegistryPermCache.Put(cacheKey, util.ObjToString(permission), time.Minute*10)
		}
	}
	return permission
}

// 检查用户权限
// 2018-01-19 22:39
func SelectUserPermissions(ai *authz.AuthRequestInfo) []string {
	actions := make([]string, 0)
	// 查询用户是否有权限
	set := util.Lock{}
	// 先检查管理员
	key := ai.Account + "_admin" + "_" + ai.Service
	r := cache.RegistryPermCache.Get(key)
	logs.Info("获取到管理员cache", r)
	if len(util.ObjToString(r)) > 5 {
		logs.Info("获取到管理员操作", ai.Service, ai.Account, ai.IP)
		if strings.Join(ai.Actions, "") == "*" {
			return []string{"*"}
		}
		return []string{"pull", "push"}
	}
	permission := getPermissions(ai)
	if len(permission) > 0 {
		for _, v := range permission {
			logs.Info("获取到权限", v.Action)
			for _, p := range strings.Split(v.Action, ",") {
				if p != "" {
					logs.Info("写入权限", p)
					set.Put(p, "1")
				}
			}
		}
	}
	for k := range set.GetData() {
		actions = append(actions, k)
	}
	return actions
}

// 验证用户权限
func (as *AuthServer) authorizeScope(ai *authz.AuthRequestInfo) ([]string, error) {
	if len(ai.Actions) == 1 && ai.Actions[0] == "pull" {
		logs.Info("通过公共权限获取")
		// 验证公告权限
		scopes := make([]authn.AuthScope, 0)
		scopes = append(scopes, authn.AuthScope{Name: ai.Name, Type: ai.Type, Actions: ai.Actions})
		r := authn.VerfiyPublicUser(scopes, ai.Service)
		if r {
			return []string{"pull"}, nil
		}
	}
	//imageName := ai.Name
	r := SelectUserPermissions(ai)
	if len(r) > 0 {
		return r, nil
	}
	logs.Info("authorizeScope ", as.authorizers)
	return make([]string, 0), errcode.ErrorCodeUnauthorized
}

func (as *AuthServer) Authorize(ar *authRequest) ([]authzResult, error) {

	ares := []authzResult{}
	for _, scope := range ar.Scopes {
		ai := &authz.AuthRequestInfo{
			Account: ar.Account,
			Type:    scope.Type,
			Name:    scope.Name,
			Service: ar.Service,
			IP:      ar.RemoteIP,
			Actions: scope.Actions,
			Labels:  ar.Labels,
		}
		t, _ := json.Marshal(ai)
		logs.Info("scope对象", string(t))
		actions, err := as.authorizeScope(ai)
		if err != nil {
			logs.Error("authorizeScope", err)
			return nil, err
		}
		ares = append(ares, authzResult{scope: scope, autorizedActions: actions})
	}
	return ares, nil
}

// https://github.com/docker/distribution/blob/master/docs/spec/auth/token.md#example
func (as *AuthServer) CreateToken(ar *authRequest, ares []authzResult) (string, error) {
	now := time.Now().Unix()
	tc := &as.config.Token

	// Sign something dummy to find out which algorithm is used.
	_, sigAlg, err := tc.privateKey.Sign(strings.NewReader("dummy"), 0)
	if err != nil {
		return "", fmt.Errorf("failed to sign: %s", err)
	}
	header := token.Header{
		Type:       "JWT",
		SigningAlg: sigAlg,
		KeyID:      tc.publicKey.KeyID(),
	}
	headerJSON, err := json.Marshal(header)
	if err != nil {
		return "", fmt.Errorf("failed to marshal header: %s", err)
	}

	claims := token.ClaimSet{
		Issuer:     tc.Issuer,
		Subject:    ar.Account,
		Audience:   ar.Service,
		NotBefore:  now - 10,
		IssuedAt:   now,
		Expiration: now + tc.Expiration,
		JWTID:      fmt.Sprintf("%d", rand.Int63()),
		Access:     []*token.ResourceActions{},
	}
	for _, a := range ares {
		ra := &token.ResourceActions{
			Type:    a.scope.Type,
			Name:    a.scope.Name,
			Actions: a.autorizedActions,
		}
		if ra.Actions == nil {
			ra.Actions = []string{}
		}
		sort.Strings(ra.Actions)
		claims.Access = append(claims.Access, ra)
	}
	claimsJSON, err := json.Marshal(claims)
	if err != nil {
		return "", fmt.Errorf("failed to marshal claims: %s", err)
	}

	payload := fmt.Sprintf("%s%s%s", joseBase64UrlEncode(headerJSON), token.TokenSeparator, joseBase64UrlEncode(claimsJSON))

	sig, sigAlg2, err := tc.privateKey.Sign(strings.NewReader(payload), 0)
	if err != nil || sigAlg2 != sigAlg {
		return "", fmt.Errorf("failed to sign token: %s", err)
	}
	logs.Info("New token for %s %+v: %s", *ar, ar.Labels, claimsJSON)
	return fmt.Sprintf("%s%s%s", payload, token.TokenSeparator, joseBase64UrlEncode(sig)), nil
}

func (as *AuthServer) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	glog.V(3).Infof("Request: %+v", req)
	path_prefix := as.config.Server.PathPrefix
	switch {
	case req.URL.Path == path_prefix+"/":
		as.doIndex(rw, req)
	case req.URL.Path == path_prefix+"/auth":
		as.doAuth(rw, req)
	default:
		http.Error(rw, "Not found", http.StatusNotFound)
		return
	}
}

// https://developers.google.com/identity/sign-in/web/server-side-flow
func (as *AuthServer) doIndex(rw http.ResponseWriter, req *http.Request) {
	rw.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(rw, "<h1>%s</h1>\n", as.config.Token.Issuer)

}

func (as *AuthServer) doAuth(rw http.ResponseWriter, req *http.Request) {
	ar, err := as.ParseRequest(req)

	ares := make([]authzResult, 0)
	if err != nil {
		logs.Warn("Bad request: %s", err)
		http.Error(rw, fmt.Sprintf("Bad request: %s", err), http.StatusBadRequest)
		return
	}
	logs.Info("Auth request: %+v", ar)
	{
		authnResult, labels, authtype, err := as.Authenticate(ar)
		if err != nil {
			http.Error(rw, fmt.Sprintf("Authentication failed (%s)", err), http.StatusInternalServerError)
			return
		}
		if !authnResult {
			logs.Error("Auth failed: %s", *ar)
			rw.Header()["WWW-Authenticate"] = []string{fmt.Sprintf(`Basic realm="%s"`, as.config.Token.Issuer)}
			http.Error(rw, "Auth failed.", http.StatusUnauthorized)
			return
		}
		ar.Labels = labels
		ar.AuthType = authtype
	}
	if len(ar.Scopes) > 0 {
		logs.Info("Authorize start", ar)
		ares, err = as.Authorize(ar)
		if err != nil {
			http.Error(rw, fmt.Sprintf("Authorization failed (%s)", err), http.StatusInternalServerError)
			return
		}
	} else {
		// Authentication-only request ("docker login"), pass through.
	}
	token, err := as.CreateToken(ar, ares)
	if err != nil {
		msg := fmt.Sprintf("Failed to generate token %s", err)
		http.Error(rw, msg, http.StatusInternalServerError)
		logs.Error("%s: %s", ar, msg)
		return
	}
	// 排除管理员和公开账号下载
	if len(ar.Scopes) > 0 && ar.AuthType == "static" {
		// 查询对象和用户的权限
		services := strings.Split(ar.Service, ".")
		if len(services) < 2 {
			services = append(services, "")
		}
		if ar.Account == "" {
			ar.Account = "public"
		}
		for _, v := range ar.Scopes {
			imglog := registry.CloudImageLog{}
			imglog.CreateUser = ar.Account
			imglog.Ip = ar.RemoteAddr
			imglog.CreateTime = util.GetDate()
			imglog.ClusterName = services[1]
			imglog.Repositories = services[0]
			imglog.Name = v.Name
			imglog.OperType = strings.Join(v.Actions, ",")
			imglog.RepositoriesGroup = strings.Split(v.Name, "/")[0]
			q := sql.InsertSql(imglog, registry.InsertCloudImageLog)
			sql.GetOrm().Raw(q).Exec()
		}
	}
	result, _ := json.Marshal(&map[string]string{"token": token})
	logs.Info("%s", result)
	rw.Header().Set("Content-Type", "application/json")
	rw.Write(result)
}

//func (as *AuthServer) Stop() {
//	for _, an := range as.authenticators {
//		an.Stop()
//	}
//	for _, az := range as.authorizers {
//		az.Stop()
//	}
//	glog.Infof("Server stopped")
//}

// Copy-pasted from libtrust where it is private.
func joseBase64UrlEncode(b []byte) string {
	return strings.TrimRight(base64.URLEncoding.EncodeToString(b), "=")
}
