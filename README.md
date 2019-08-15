# elton-basic-auth

[![Build Status](https://img.shields.io/travis/vicanso/elton-basic-auth.svg?label=linux+build)](https://travis-ci.org/vicanso/elton-basic-auth)


Basic auth middleware for elton.

```go
package main

import (
	"bytes"
	"fmt"

	"github.com/vicanso/elton"
	basicauth "github.com/vicanso/elton-basic-auth"
	"github.com/vicanso/hes"
)

func main() {
	d := elton.New()

	d.Use(basicauth.New(basicauth.Config{
		Validate: func(account, pwd string, c *elton.Context) (bool, error) {
			if account == "tree.xie" && pwd == "password" {
				return true, nil
			}
			if account == "n" {
				return false, hes.New("account is invalid")
			}
			return false, nil
		},
	}))

	d.GET("/", func(c *elton.Context) (err error) {
		c.BodyBuffer = bytes.NewBufferString("hello world")
		return
	})

	d.ListenAndServe(":7001")
}
```