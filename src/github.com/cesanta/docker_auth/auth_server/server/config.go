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
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/cesanta/docker_auth/auth_server/authn"
	"github.com/docker/libtrust"
	"path/filepath"
)

type Config struct {
	Server     ServerConfig                   `yaml:"server"`
	Token      TokenConfig                    `yaml:"token"`
	Users      map[string]*authn.Requirements `yaml:"users,omitempty"`
}

type ServerConfig struct {
	ListenAddress string            `yaml:"addr,omitempty"`
	PathPrefix    string            `yaml:"path_prefix,omitempty"`
	RealIPHeader  string            `yaml:"real_ip_header,omitempty"`
	RealIPPos     int               `yaml:"real_ip_pos,omitempty"`
	CertFile      string            `yaml:"certificate,omitempty"`
	KeyFile       string            `yaml:"key,omitempty"`
	LetsEncrypt   LetsEncryptConfig `yaml:"letsencrypt,omitempty"`

	publicKey  libtrust.PublicKey
	privateKey libtrust.PrivateKey
}

type LetsEncryptConfig struct {
	Host     string `yaml:"host,omitempty"`
	Email    string `yaml:"email,omitempty"`
	CacheDir string `yaml:"cache_dir,omitempty"`
}

type TokenConfig struct {
	Issuer     string `yaml:"issuer,omitempty"`
	CertFile   string `yaml:"certificate,omitempty"`
	KeyFile    string `yaml:"key,omitempty"`
	Expiration int64  `yaml:"expiration,omitempty"`

	publicKey  libtrust.PublicKey
	privateKey libtrust.PrivateKey
}

func validate(c *Config) error {
	if c.Server.ListenAddress == "" {
		return errors.New("server.addr is required")
	}
	if c.Server.PathPrefix != "" && !strings.HasPrefix(c.Server.PathPrefix, "/") {
		return errors.New("server.path_prefix must be an absolute path")
	}

	if c.Token.Issuer == "" {
		return errors.New("token.issuer is required")
	}
	if c.Token.Expiration <= 0 {
		return fmt.Errorf("expiration must be positive, got %d", c.Token.Expiration)
	}
	//if c.Users == nil {
	//	return errors.New("no auth methods are configured, this is probably a mistake. Use an empty user map if you really want to deny everyone.")
	//}

	return nil
}

func loadCertAndKey(certFile, keyFile string) (pk libtrust.PublicKey, prk libtrust.PrivateKey, err error) {
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return
	}
	x509Cert, err := x509.ParseCertificate(cert.Certificate[0])
	if err != nil {
		return
	}
	pk, err = libtrust.FromCryptoPublicKey(x509Cert.PublicKey)
	if err != nil {
		return
	}
	prk, err = libtrust.FromCryptoPrivateKey(cert.PrivateKey)
	return
}

func LoadConfig() (*Config, error) {
	//contents, err := ioutil.ReadFile(fileName)
	//if err != nil {
	//	return nil, fmt.Errorf("could not read %s: %s", fileName, err)
	//}
	var err error
	c := &Config{}
	c.Token = TokenConfig{Issuer:"Acme auth server",Expiration:900}
	c.Server = ServerConfig{ListenAddress:":5001"}
	pwd,_ := os.Getwd()
	key := filepath.Join(pwd,"conf", "key", "server.key")
	pem := filepath.Join(pwd,"conf", "key", "server.pem")
	c.Server.CertFile = pem
	c.Server.KeyFile = key
	//if err = yaml.Unmarshal(contents, c); err != nil {
	//	return nil, fmt.Errorf("could not parse config: %s", err)
	//}
	if err := validate(c); err != nil {
		return nil, fmt.Errorf("invalid config: %s", err)
	}
	serverConfigured := false
	if c.Server.CertFile != "" || c.Server.KeyFile != "" {

		// Check for partial configuration.
		if c.Server.CertFile == "" || c.Server.KeyFile == "" {
			return nil, fmt.Errorf("failed to load server cert and key: both were not provided")
		}
		c.Server.publicKey, c.Server.privateKey, err = loadCertAndKey(c.Server.CertFile, c.Server.KeyFile)
		if err != nil {
			return nil, fmt.Errorf("failed to load server cert and key: %s", err)
		}
		serverConfigured = true
	}
	tokenConfigured := false
	if c.Token.CertFile != "" || c.Token.KeyFile != "" {
		// Check for partial configuration.
		if c.Token.CertFile == "" || c.Token.KeyFile == "" {
			return nil, fmt.Errorf("failed to load token cert and key: both were not provided")
		}
		c.Token.publicKey, c.Token.privateKey, err = loadCertAndKey(c.Token.CertFile, c.Token.KeyFile)
		if err != nil {
			return nil, fmt.Errorf("failed to load token cert and key: %s", err)
		}
		tokenConfigured = true
	}

	if serverConfigured && !tokenConfigured {
		c.Token.publicKey, c.Token.privateKey = c.Server.publicKey, c.Server.privateKey
		tokenConfigured = true
	}

	if !tokenConfigured {
		return nil, fmt.Errorf("failed to load token cert and key: none provided")
	}

	if !serverConfigured && c.Server.LetsEncrypt.Email != "" {
		if c.Server.LetsEncrypt.CacheDir == "" {
			return nil, fmt.Errorf("server.letsencrypt.cache_dir is required")
		}
		// We require that LetsEncrypt is an existing directory, because we really don't want it
		// to be misconfigured and obtained certificates to be lost.
		fi, err := os.Stat(c.Server.LetsEncrypt.CacheDir)
		if err != nil || !fi.IsDir() {
			return nil, fmt.Errorf("server.letsencrypt.cache_dir (%s) does not exist or is not a directory", c.Server.LetsEncrypt.CacheDir)
		}
	}
	return c, nil
}
