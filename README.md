# cod-basic-auth

[![Build Status](https://img.shields.io/travis/vicanso/cod-basic-auth.svg?label=linux+build)](https://travis-ci.org/vicanso/cod-basic-auth)


Basic auth middleware for cod.

```go
package main

import (
	"bytes"
	"fmt"

	"github.com/vicanso/cod"
	basicauth "github.com/vicanso/cod-basic-auth"
	"github.com/vicanso/hes"
)

func main() {
	d := cod.New()

	d.Use(basicauth.NewBasicAuth(basicauth.Config{
		Validate: func(account, pwd string, c *cod.Context) (bool, error) {
			if account == "tree.xie" && pwd == "password" {
				return true, nil
			}
			if account == "n" {
				return false, hes.New("account is invalid")
			}
			return false, nil
		},
	}))

	d.GET("/", func(c *cod.Context) (err error) {
		c.BodyBuffer = bytes.NewBufferString("hello world")
		return
	})

	d.ListenAndServe(":7001")
}
```