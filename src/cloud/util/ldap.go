package util

import (
	"gopkg.in/ldap.v2"
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego"
	"golang.org/x/crypto/openpgp/errors"
	"strings"
)

// ldap登录验证
// 2018-01-19 09:59
func LdapLoginAuth(username string, passowrd string) (bool,error) {
	server := beego.AppConfig.String("ldap.server")
	port, err := beego.AppConfig.Int("ldap.port")
	prefix := beego.AppConfig.String("ldap.prefix")
	username = prefix +"\\"+ username
	if server == "" || prefix == "" {
		return false,errors.ErrUnknownIssuer
	}

	if err != nil {
		return false,err
	}

	l, err := ldap.Dial("tcp", fmt.Sprintf("%s:%d", server, port))
	if err != nil {
		logs.Error(err)
		return false,err
	}

	defer l.Close()
	controls := []ldap.Control{}
	controls = append(controls, ldap.NewControlBeheraPasswordPolicy())
	//bindRequest := ldap.NewSimpleBindRequest(strings.TrimSpace(prefix+username), strings.TrimSpace(passowrd), controls)

	bindRequest := ldap.NewSimpleBindRequest(username, strings.TrimSpace(passowrd), controls)
	_, err = l.SimpleBind(bindRequest)
	if err != nil {
		return false,err
	}
	return true,err
}