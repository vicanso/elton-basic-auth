// Copyright 2018 tree xie
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package basicauth

import (
	"errors"
	"net/http"

	"github.com/vicanso/cod"
	"github.com/vicanso/hes"
)

const (
	defaultRealm = "basic auth tips"
	// ErrCategory basic auth error category
	ErrCategory = "cod-basic-auth"
)

type (
	// Validate validate function
	Validate func(username string, password string, c *cod.Context) (bool, error)
	// Config basic auth config
	Config struct {
		Realm    string
		Validate Validate
		Skipper  cod.Skipper
	}
)

var (
	// errUnauthorized unauthorized err
	errUnauthorized = getBasicAuthError(errors.New("unAuthorized"), http.StatusUnauthorized)
)

func getBasicAuthError(err error, statusCode int) *hes.Error {
	he := hes.Wrap(err)
	he.StatusCode = statusCode
	he.Category = ErrCategory
	return he
}

// New new basic auth
func New(config Config) cod.Handler {
	if config.Validate == nil {
		panic("require validate function")
	}
	basic := "basic"
	realm := defaultRealm
	if config.Realm != "" {
		realm = config.Realm
	}
	wwwAuthenticate := basic + " realm=" + realm
	skipper := config.Skipper
	if skipper == nil {
		skipper = cod.DefaultSkipper
	}
	return func(c *cod.Context) (err error) {
		if skipper(c) {
			return c.Next()
		}
		user, password, hasAuth := c.Request.BasicAuth()
		// 如果请求头无认证头，则返回出错
		if !hasAuth {
			c.SetHeader(cod.HeaderWWWAuthenticate, wwwAuthenticate)
			err = errUnauthorized
			return
		}

		valid, e := config.Validate(user, password, c)

		// 如果返回出错，则输出出错信息
		if e != nil {
			err, ok := e.(*hes.Error)
			if !ok {
				err = getBasicAuthError(e, http.StatusBadRequest)
			}
			return err
		}

		// 如果校验失败，设置认证头，客户重新输入
		if !valid {
			c.SetHeader(cod.HeaderWWWAuthenticate, wwwAuthenticate)
			err = errUnauthorized
			return
		}
		return c.Next()
	}
}
